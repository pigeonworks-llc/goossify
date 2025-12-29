package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/pigeonworks-llc/goossify/internal/analyzer"
	"github.com/pigeonworks-llc/goossify/internal/ossify"
	"github.com/spf13/cobra"
)

var (
	statusFormat      string
	statusFix         bool
	statusGitHub      bool
	statusGitHubToken string
)

// statusCmd represents the status command
var statusCmd = &cobra.Command{
	Use:   "status [path]",
	Short: "Check project health and OSS readiness",
	Long: `Analyze project health and identify missing elements required for OSS.

This command checks the following:
• Basic project structure
• Documentation status
• GitHub integration (CI/CD, templates, etc.)
• Quality tools (Linter, tests, etc.)
• Dependency management
• License information

Displays scores and recommendations, suggesting improvements.`,
	Args: cobra.MaximumNArgs(1),
	RunE: runStatus,
}

func init() {
	rootCmd.AddCommand(statusCmd)

	statusCmd.Flags().StringVarP(&statusFormat, "format", "f", "human", "Output format (human, json)")
	statusCmd.Flags().BoolVar(&statusFix, "fix", false, "Execute automatic fixes for missing items")
	statusCmd.Flags().BoolVar(&statusGitHub, "github", false, "Include GitHub settings check (requires TOKEN)")
	statusCmd.Flags().StringVar(&statusGitHubToken, "github-token", "", "GitHub Personal Access Token")
}

func runStatus(cmd *cobra.Command, args []string) error {
	// 分析対象パスの決定
	var targetPath string
	if len(args) == 0 {
		targetPath = "."
	} else {
		targetPath = args[0]
	}

	// Convert to absolute path
	absPath, err := filepath.Abs(targetPath)
	if err != nil {
		return fmt.Errorf("path resolution error: %w", err)
	}

	// Check directory existence
	if _, err := os.Stat(absPath); os.IsNotExist(err) {
		return fmt.Errorf("specified path does not exist: %s", absPath)
	}

	fmt.Printf("🔍 Starting project analysis: %s\n\n", absPath)

	// Execute project analysis
	analyzer := analyzer.New(absPath)
	result, err := analyzer.Analyze()
	if err != nil {
		return fmt.Errorf("error occurred during analysis: %w", err)
	}

	// GitHub settings check (optional)
	var githubCheck interface{}
	if statusGitHub {
		token := statusGitHubToken
		if token == "" {
			token = os.Getenv("GITHUB_TOKEN")
		}
		if token == "" {
			fmt.Println("⚠️  GitHub settings check requires a token (--github-token or GITHUB_TOKEN)")
		} else {
			fmt.Println("🔍 Analyzing GitHub settings...")
			// GitHub analysis feature temporarily disabled
			fmt.Println("GitHub analysis feature is under development")
		}
	}

	// Output results
	switch statusFormat {
	case "json":
		return outputJSON(result, githubCheck)
	default:
		return outputHuman(result, githubCheck)
	}
}

func outputJSON(result *analyzer.AnalysisResult, githubCheck interface{}) error {
	output := map[string]interface{}{
		"analysis": result,
	}
	if githubCheck != nil {
		output["github"] = githubCheck
	}

	jsonData, err := json.MarshalIndent(output, "", "  ")
	if err != nil {
		return fmt.Errorf("JSON conversion error: %w", err)
	}
	fmt.Println(string(jsonData))
	return nil
}

func outputHuman(result *analyzer.AnalysisResult, githubCheck interface{}) error {
	// Header information
	fmt.Printf("📊 Project: %s (%s)\n", result.ProjectName, result.ProjectType)
	fmt.Printf("📈 Overall Score: %d/100 (%s)\n\n", result.OverallScore, getScoreEmoji(result.OverallScore))

	// Category results
	fmt.Println("📋 Category Results:")
	for _, category := range result.Categories {
		statusIcon := getStatusIcon(category.Status)
		fmt.Printf("  %s %s: %d/100\n", statusIcon, category.Name, category.Score)

		// Detailed display (when score is low)
		if category.Score < 80 {
			for _, item := range category.Items {
				if item.Status == "missing" && item.Required {
					fmt.Printf("    ❌ %s (required)\n", item.Name)
				} else if item.Status == "missing" {
					fmt.Printf("    ⚠️  %s (recommended)\n", item.Name)
				}
			}
		}
	}

	// Missing items
	if len(result.Missing) > 0 {
		fmt.Printf("\n🚨 Missing Items (%d items):\n", len(result.Missing))
		for _, missing := range result.Missing {
			priorityIcon := getPriorityIcon(missing.Priority)
			fmt.Printf("  %s %s - %s\n", priorityIcon, missing.Name, missing.Description)
		}
	}

	// Recommendations
	if len(result.Recommendations) > 0 {
		fmt.Printf("\n💡 Recommendations (%d items):\n", len(result.Recommendations))
		for _, rec := range result.Recommendations {
			priorityIcon := getPriorityIcon(rec.Priority)
			fmt.Printf("  %s %s\n", priorityIcon, rec.Title)
			fmt.Printf("     %s\n", rec.Description)
			if rec.Command != "" {
				fmt.Printf("     Run: %s\n", rec.Command)
			}
		}
	}

	// GitHub settings (optional)
	if githubCheck != nil {
		fmt.Printf("\n🐙 GitHub Settings: Under development\n")
	}

	// Summary
	fmt.Printf("\n📝 %s\n", result.Summary)

	// Auto-fix suggestions
	if statusFix {
		fmt.Printf("\n🔧 Running automatic fixes...\n")
		return runAutoFix(result)
	} else if len(result.Missing) > 0 {
		fmt.Printf("\n💡 Many issues can be auto-fixed with 'goossify ossify .'\n")
		fmt.Printf("   To run auto-fixes: goossify status --fix\n")
	}

	return nil
}

func runAutoFix(result *analyzer.AnalysisResult) error {
	fmt.Println("🔧 Running automatic fixes...")

	// Get current working directory
	projectPath, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %w", err)
	}

	// Create and execute Ossifier
	ossifier := ossify.New(projectPath)
	ossifier.SetInteractive(false)
	ossifier.SetDryRun(false)

	if err := ossifier.Execute(); err != nil {
		return fmt.Errorf("auto-fix failed: %w", err)
	}

	fmt.Println("\n✅ Auto-fix completed!")
	fmt.Println("Run 'goossify status' again to verify improvements")

	return nil
}

func getScoreEmoji(score int) string {
	if score >= 90 {
		return "🏆 Excellent"
	} else if score >= 80 {
		return "✅ Good"
	} else if score >= 60 {
		return "⚠️ Fair"
	} else if score >= 40 {
		return "❌ Needs Improvement"
	}
	return "🆘 Critical"
}

func getStatusIcon(status string) string {
	switch status {
	case "good":
		return "✅"
	case "warning":
		return "⚠️"
	case "error":
		return "❌"
	default:
		return "❓"
	}
}

func getPriorityIcon(priority string) string {
	switch priority {
	case "high":
		return "🔴"
	case "medium":
		return "🟡"
	case "low":
		return "🟢"
	default:
		return "⚪"
	}
}
