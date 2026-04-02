package cmd

import (
	"fmt"
	"sync"

	"github.com/pigeonworks-llc/goossify/internal/analyzer"
	"github.com/pigeonworks-llc/goossify/internal/exitcode"
	"github.com/pigeonworks-llc/goossify/internal/ossify"
	"github.com/pigeonworks-llc/goossify/internal/output"
	"github.com/spf13/cobra"
)

var (
	batchFormat   string
	batchParallel int
)

// BatchResult represents the result of batch processing.
type BatchResult struct {
	TotalCount   int             `json:"total_count"`
	SuccessCount int             `json:"success_count"`
	FailCount    int             `json:"fail_count"`
	Results      []PackageResult `json:"results"`
}

// PackageResult represents the result of processing a single package.
type PackageResult struct {
	Path    string `json:"path"`
	Success bool   `json:"success"`
	Score   int    `json:"score,omitempty"`
	Ready   bool   `json:"ready,omitempty"`
	Error   string `json:"error,omitempty"`
}

var batchCmd = &cobra.Command{
	Use:   "batch <command> <paths...>",
	Short: "Run a command on multiple packages",
	Long: `Execute a goossify command on multiple packages.

Supported commands: status, ready, ossify, pipeline

Examples:
  goossify batch status ./pkg1 ./pkg2 ./pkg3
  goossify batch status --format json ./packages/*
  goossify batch pipeline --parallel 4 ./oss-packages/*
  goossify batch ready ./oss-packages/*`,
	Args: cobra.MinimumNArgs(2),
	RunE: runBatch,
}

func init() {
	rootCmd.AddCommand(batchCmd)

	batchCmd.Flags().StringVarP(&batchFormat, "format", "f", "human", "Output format (human, json)")
	batchCmd.Flags().IntVar(&batchParallel, "parallel", 1, "Number of parallel executions")
}

func runBatch(cmd *cobra.Command, args []string) error {
	command := args[0]
	paths := args[1:]

	// Validate command
	validCommands := map[string]bool{
		"status":   true,
		"ready":    true,
		"ossify":   true,
		"pipeline": true,
	}
	if !validCommands[command] {
		return fmt.Errorf("invalid command: %s (valid: status, ready, ossify, pipeline)", command)
	}

	isJSON := batchFormat == "json"
	result := &BatchResult{
		TotalCount: len(paths),
		Results:    make([]PackageResult, 0, len(paths)),
	}

	if !isJSON {
		fmt.Printf("🚀 Running batch %s on %d packages...\n\n", command, len(paths))
	}

	var mu sync.Mutex
	var wg sync.WaitGroup
	semaphore := make(chan struct{}, batchParallel)

	for _, path := range paths {
		wg.Add(1)
		go func(p string) {
			defer wg.Done()
			semaphore <- struct{}{}
			defer func() { <-semaphore }()

			pkgResult := executeOnPackage(command, p, isJSON)

			mu.Lock()
			result.Results = append(result.Results, pkgResult)
			if pkgResult.Success {
				result.SuccessCount++
			} else {
				result.FailCount++
			}
			mu.Unlock()
		}(path)
	}

	wg.Wait()

	// Set exit code
	if result.FailCount > 0 {
		if result.SuccessCount > 0 {
			ExitCode = exitcode.Warning // Partial success
		} else {
			ExitCode = exitcode.Error // All failed
		}
	}

	// Output
	if isJSON {
		formatter := output.New("json")
		return formatter.JSON(result)
	}

	// Human-readable summary
	fmt.Println("\n" + "─────────────────────────────────────")
	fmt.Printf("📊 Batch Results:\n")
	fmt.Printf("   Total:   %d\n", result.TotalCount)
	fmt.Printf("   Success: %d\n", result.SuccessCount)
	fmt.Printf("   Failed:  %d\n", result.FailCount)

	if result.FailCount > 0 {
		fmt.Println("\n❌ Failed packages:")
		for _, r := range result.Results {
			if !r.Success {
				fmt.Printf("   - %s: %s\n", r.Path, r.Error)
			}
		}
	}

	return nil
}

func executeOnPackage(command, path string, silent bool) PackageResult {
	result := PackageResult{Path: path, Success: true}

	switch command {
	case "status":
		projectAnalyzer := analyzer.New(path)
		analysisResult, err := projectAnalyzer.Analyze()
		if err != nil {
			result.Success = false
			result.Error = err.Error()
		} else {
			result.Score = analysisResult.OverallScore
			if !silent {
				fmt.Printf("  ✅ %s: %d/100\n", path, result.Score)
			}
		}

	case "ready":
		projectAnalyzer := analyzer.New(path)
		analysisResult, err := projectAnalyzer.Analyze()
		if err != nil {
			result.Success = false
			result.Error = err.Error()
		} else {
			result.Score = analysisResult.OverallScore
			result.Ready = analysisResult.OverallScore >= 90
			if !result.Ready {
				result.Success = false
				result.Error = fmt.Sprintf("score %d < 90", result.Score)
			}
			if !silent {
				if result.Ready {
					fmt.Printf("  ✅ %s: Ready (score: %d)\n", path, result.Score)
				} else {
					fmt.Printf("  ❌ %s: Not ready (score: %d)\n", path, result.Score)
				}
			}
		}

	case "ossify":
		ossifier := ossify.New(path)
		ossifier.SetInteractive(false)
		ossifier.SetDryRun(false)
		if err := ossifier.Execute(); err != nil {
			result.Success = false
			result.Error = err.Error()
			if !silent {
				fmt.Printf("  ❌ %s: %v\n", path, err)
			}
		} else if !silent {
			fmt.Printf("  ✅ %s: Ossified\n", path)
		}

	case "pipeline":
		// Run pipeline steps
		ossifier := ossify.New(path)
		ossifier.SetInteractive(false)
		ossifier.SetDryRun(false)
		if err := ossifier.Execute(); err != nil {
			result.Success = false
			result.Error = fmt.Sprintf("ossify failed: %v", err)
			if !silent {
				fmt.Printf("  ❌ %s: %v\n", path, err)
			}
			return result
		}

		projectAnalyzer := analyzer.New(path)
		analysisResult, err := projectAnalyzer.Analyze()
		if err != nil {
			result.Success = false
			result.Error = fmt.Sprintf("status failed: %v", err)
			if !silent {
				fmt.Printf("  ❌ %s: %v\n", path, err)
			}
		} else {
			result.Score = analysisResult.OverallScore
			result.Ready = analysisResult.OverallScore >= 90
			if !silent {
				if result.Ready {
					fmt.Printf("  ✅ %s: Pipeline complete (score: %d, ready: yes)\n", path, result.Score)
				} else {
					fmt.Printf("  ⚠️  %s: Pipeline complete (score: %d, ready: no)\n", path, result.Score)
				}
			}
		}
	}

	return result
}
