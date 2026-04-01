package cmd

import (
	"fmt"

	"github.com/pigeonworks-llc/goossify/internal/analyzer"
	"github.com/pigeonworks-llc/goossify/internal/exitcode"
	"github.com/pigeonworks-llc/goossify/internal/ossify"
	"github.com/pigeonworks-llc/goossify/internal/output"
	"github.com/spf13/cobra"
)

var (
	pipelineFormat    string
	pipelineSkipReady bool
	pipelineDryRun    bool
	pipelineThreshold int
)

// PipelineResult represents the result of a pipeline execution.
type PipelineResult struct {
	Path        string               `json:"path"`
	OssifyDone  bool                 `json:"ossify_done"`
	StatusScore int                  `json:"status_score"`
	Ready       bool                 `json:"ready"`
	Steps       []PipelineStepResult `json:"steps"`
	ExitCode    int                  `json:"exit_code"`
}

// PipelineStepResult represents the result of a single pipeline step.
type PipelineStepResult struct {
	Name    string `json:"name"`
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
	Error   string `json:"error,omitempty"`
}

var pipelineCmd = &cobra.Command{
	Use:   "pipeline [path]",
	Short: "Run full OSS pipeline: ossify -> status -> ready",
	Long: `Execute the complete OSS preparation pipeline in sequence.

This command runs:
1. ossify - Generate missing OSS files
2. status - Analyze project health
3. ready  - Check release readiness (optional)

Use --skip-ready to stop after status check.

Examples:
  goossify pipeline .
  goossify pipeline --skip-ready .
  goossify pipeline --format json .
  goossify pipeline --dry-run .
  goossify pipeline --threshold 80 .`,
	Args: cobra.MaximumNArgs(1),
	RunE: runPipeline,
}

func init() {
	rootCmd.AddCommand(pipelineCmd)

	pipelineCmd.Flags().StringVarP(&pipelineFormat, "format", "f", "human", "Output format (human, json)")
	pipelineCmd.Flags().BoolVar(&pipelineSkipReady, "skip-ready", false, "Skip ready check")
	pipelineCmd.Flags().BoolVar(&pipelineDryRun, "dry-run", false, "Preview changes without making them")
	pipelineCmd.Flags().IntVar(&pipelineThreshold, "threshold", 80, "Minimum score threshold for success")
}

func runPipeline(cmd *cobra.Command, args []string) error {
	targetPath := "."
	if len(args) > 0 {
		targetPath = args[0]
	}

	isJSON := pipelineFormat == "json"
	result := &PipelineResult{
		Path:     targetPath,
		ExitCode: exitcode.Success,
	}

	if !isJSON {
		fmt.Printf("🚀 Running OSS pipeline for: %s\n\n", targetPath)
	}

	// Step 1: Ossify
	ossifyStep := PipelineStepResult{Name: "ossify", Success: true}
	if !isJSON {
		fmt.Println("📦 Step 1: Running ossify...")
	}

	ossifier := ossify.New(targetPath)
	ossifier.SetInteractive(false)
	ossifier.SetDryRun(pipelineDryRun)

	if err := ossifier.Execute(); err != nil {
		ossifyStep.Success = false
		ossifyStep.Error = err.Error()
		result.ExitCode = exitcode.Error
		if !isJSON {
			fmt.Printf("   ❌ Ossify failed: %v\n", err)
		}
	} else {
		result.OssifyDone = true
		ossifyStep.Message = "OSS files generated"
		if !isJSON {
			fmt.Println("   ✅ Ossify completed")
		}
	}
	result.Steps = append(result.Steps, ossifyStep)

	// Step 2: Status
	statusStep := PipelineStepResult{Name: "status", Success: true}
	if !isJSON {
		fmt.Println("\n📊 Step 2: Checking status...")
	}

	projectAnalyzer := analyzer.New(targetPath)
	analysisResult, err := projectAnalyzer.Analyze()
	if err != nil {
		statusStep.Success = false
		statusStep.Error = err.Error()
		result.ExitCode = exitcode.Error
		if !isJSON {
			fmt.Printf("   ❌ Status check failed: %v\n", err)
		}
	} else {
		result.StatusScore = analysisResult.OverallScore
		statusStep.Message = fmt.Sprintf("Score: %d/100", analysisResult.OverallScore)

		if analysisResult.OverallScore < pipelineThreshold {
			statusStep.Success = false
			statusStep.Error = fmt.Sprintf("score %d below threshold %d", analysisResult.OverallScore, pipelineThreshold)
			if result.ExitCode == exitcode.Success {
				result.ExitCode = exitcode.ValidationFail
			}
			if !isJSON {
				fmt.Printf("   ⚠️  Score %d/100 (below threshold %d)\n", analysisResult.OverallScore, pipelineThreshold)
			}
		} else if !isJSON {
			fmt.Printf("   ✅ Score: %d/100\n", analysisResult.OverallScore)
		}
	}
	result.Steps = append(result.Steps, statusStep)

	// Step 3: Ready (optional)
	if !pipelineSkipReady {
		readyStep := PipelineStepResult{Name: "ready", Success: true}
		if !isJSON {
			fmt.Println("\n🎯 Step 3: Checking release readiness...")
		}

		// Check if score is high enough for ready
		if result.StatusScore >= 90 {
			result.Ready = true
			readyStep.Message = "Project is ready for release"
			if !isJSON {
				fmt.Println("   ✅ Project is ready for public release")
			}
		} else {
			readyStep.Success = false
			readyStep.Error = fmt.Sprintf("score %d is below 90 required for release", result.StatusScore)
			if result.ExitCode == exitcode.Success {
				result.ExitCode = exitcode.ValidationFail
			}
			if !isJSON {
				fmt.Printf("   ⚠️  Score %d/100 is below 90 required for release\n", result.StatusScore)
			}
		}
		result.Steps = append(result.Steps, readyStep)
	}

	// Set global exit code
	ExitCode = result.ExitCode

	// Output
	if isJSON {
		formatter := output.New("json")
		return formatter.JSON(result)
	}

	// Summary
	fmt.Println("\n" + "─────────────────────────────────────")
	switch result.ExitCode {
	case exitcode.Success:
		fmt.Println("✅ Pipeline completed successfully!")
	case exitcode.ValidationFail:
		fmt.Println("⚠️  Pipeline completed with validation failures")
	default:
		fmt.Println("❌ Pipeline failed")
	}

	return nil
}
