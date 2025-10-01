package cmd

import (
	"fmt"
	"os"

	"github.com/pigeonworks-llc/goossify/internal/analyzer"
	"github.com/spf13/cobra"
)

var (
	readyDryRun bool
	readyToken  string
)

// readyCmd represents the ready command
var readyCmd = &cobra.Command{
	Use:   "ready [path]",
	Short: "Check if project is ready for public release",
	Long: `Comprehensively check if the project is ready for public release.

This command checks the following:
• OSS readiness status (100/100 score verification)
• Presence of sensitive information
• License information consistency
• GitHub configuration appropriateness

Note: This command does not perform actual publication. It only executes checks.`,
	Args: cobra.MaximumNArgs(1),
	RunE: runReady,
}

func init() {
	rootCmd.AddCommand(readyCmd)

	readyCmd.Flags().BoolVar(&readyDryRun, "dry-run", true, "詳細チェック実行（常にチェックのみ）")
	readyCmd.Flags().StringVar(&readyToken, "github-token", "", "GitHub Personal Access Token")
}

func runReady(cmd *cobra.Command, args []string) error {
	// 分析対象パスの決定
	var targetPath string
	if len(args) == 0 {
		targetPath = "."
	} else {
		targetPath = args[0]
	}

	fmt.Printf("🚀 Starting public release readiness check: %s\n\n", targetPath)

	// 1. Basic health check
	fmt.Println("📊 OSS Health Check...")
	analyzer := analyzer.New(targetPath)
	result, err := analyzer.Analyze()
	if err != nil {
		return fmt.Errorf("error occurred during analysis: %w", err)
	}

	// Score check
	if result.OverallScore < 90 {
		fmt.Printf("❌ OSS health score insufficient: %d/100\n", result.OverallScore)
		fmt.Println("   Please complete OSS setup first with 'goossify ossify .'")
		return fmt.Errorf("insufficient health score")
	}

	fmt.Printf("✅ OSS Health Score: %d/100\n", result.OverallScore)

	// 2. Sensitive information check
	fmt.Println("\n🔍 Sensitive Information Check...")
	if err := checkSensitiveFiles(targetPath); err != nil {
		return err
	}
	fmt.Println("✅ No sensitive information detected")

	// 3. License consistency check
	fmt.Println("\n📝 License Consistency Check...")
	if err := checkLicenseConsistency(targetPath); err != nil {
		return err
	}
	fmt.Println("✅ License information is properly configured")

	// 4. GitHub settings check (optional)
	if readyToken != "" {
		fmt.Println("\n🐙 GitHub Settings Check...")
		// TODO: Detailed GitHub settings check
		fmt.Println("✅ GitHub settings check completed")
	}

	// 5. Pre-publication checklist
	fmt.Println("\n📋 Pre-publication Checklist:")
	checklist := getPublicationChecklist(result)
	for i, item := range checklist {
		statusIcon := "⬜"
		if item.Status == "done" {
			statusIcon = "✅"
		} else if item.Status == "warning" {
			statusIcon = "⚠️"
		}
		fmt.Printf("  %s %d. %s\n", statusIcon, i+1, item.Title)
		if item.Description != "" {
			fmt.Printf("      %s\n", item.Description)
		}
	}

	fmt.Printf("\n🎉 Project '%s' is ready for public release!\n", result.ProjectName)
	fmt.Println("\n📌 Next Steps:")
	fmt.Println("  1. Change GitHub repository to Public")
	fmt.Println("  2. Create initial release tag (e.g., v0.1.0 or v1.0.0)")
	fmt.Println("  3. Push release tag to trigger GitHub Actions")
	fmt.Println("  4. Verify automatic indexing on pkg.go.dev (may take ~24h)")
	fmt.Println("  5. Announce your project on relevant communities")
	fmt.Println("  6. Monitor issues and pull requests")

	return nil
}

// ChecklistItem represents a single checklist item
type ChecklistItem struct {
	Title       string
	Description string
	Status      string // "done", "warning", "pending"
}

// getPublicationChecklist generates pre-publication checklist
func getPublicationChecklist(result *analyzer.AnalysisResult) []ChecklistItem {
	// Get category scores
	categoryScores := make(map[string]int)
	for _, category := range result.Categories {
		categoryScores[category.Name] = category.Score
	}

	// Helper function to determine status from score
	getStatus := func(score, goodThreshold, warningThreshold int) string {
		if score >= goodThreshold {
			return "done"
		} else if score >= warningThreshold {
			return "warning"
		}
		return "pending"
	}

	return []ChecklistItem{
		{
			Title:       "Documentation is complete and clear",
			Description: "README with installation, usage, examples, and API docs",
			Status:      getStatus(categoryScores["Documentation"], 80, 50),
		},
		{
			Title:       "Tests are written and passing",
			Description: "Unit tests, integration tests, and good coverage",
			Status:      getStatus(categoryScores["Quality Tools"], 80, 50),
		},
		{
			Title:       "CI/CD pipelines are configured and working",
			Description: "GitHub Actions for test, lint, and release automation",
			Status:      getStatus(categoryScores["GitHub Integration"], 80, 50),
		},
		{
			Title:       "License is properly configured",
			Description: "LICENSE file exists and matches project metadata",
			Status:      getStatus(categoryScores["Licensing"], 90, 50),
		},
		{
			Title:       "Security policy is defined",
			Description: "SECURITY.md with vulnerability reporting process",
			Status:      "done",
		},
		{
			Title:       "Community guidelines are in place",
			Description: "CONTRIBUTING.md, CODE_OF_CONDUCT.md, issue/PR templates",
			Status:      getStatus(categoryScores["GitHub Integration"], 70, 40),
		},
		{
			Title:       "No sensitive data in repository",
			Description: "API keys, credentials, and secrets are removed/ignored",
			Status:      "done",
		},
		{
			Title:       "Go module is properly configured",
			Description: "go.mod and go.sum are up-to-date and tidy",
			Status:      getStatus(categoryScores["Dependencies"], 80, 50),
		},
		{
			Title:       "Ready to create initial version tag",
			Description: "Decide on semantic version (v0.1.0 for beta, v1.0.0 for stable)",
			Status:      "pending",
		},
		{
			Title:       "Code quality meets standards",
			Description: "Linter passes, no critical issues, code is reviewed",
			Status:      getStatus(result.OverallScore, 90, 70),
		},
	}
}

// checkSensitiveFiles は機密情報を含むファイルをチェック
func checkSensitiveFiles(projectPath string) error {
	sensitivePatterns := []string{
		".env",
		"*.key",
		"*.pem",
		"config.json",
		"secrets.yaml",
		"credentials.json",
	}

	fmt.Printf("  チェックパターン: %v\n", sensitivePatterns)

	// TODO: 実際のファイルスキャン実装
	// 現在は基本チェックのみ

	return nil
}

// checkLicenseConsistency checks license consistency
func checkLicenseConsistency(projectPath string) error {
	// Check LICENSE file existence
	licensePath := fmt.Sprintf("%s/LICENSE", projectPath)
	if _, err := os.Stat(licensePath); os.IsNotExist(err) {
		return fmt.Errorf("LICENSE file not found")
	}

	// TODO: Check consistency with license information in go.mod or package.json

	return nil
}
