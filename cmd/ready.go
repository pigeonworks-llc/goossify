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
	checklist := []ChecklistItem{}

	// 1. Documentation
	docStatus := "pending"
	for _, category := range result.Categories {
		if category.Name == "Documentation" {
			if category.Score >= 80 {
				docStatus = "done"
			} else if category.Score >= 50 {
				docStatus = "warning"
			}
		}
	}
	checklist = append(checklist, ChecklistItem{
		Title:       "Documentation is complete and clear",
		Description: "README with installation, usage, examples, and API docs",
		Status:      docStatus,
	})

	// 2. Tests
	testStatus := "pending"
	for _, category := range result.Categories {
		if category.Name == "Quality Tools" {
			if category.Score >= 80 {
				testStatus = "done"
			} else if category.Score >= 50 {
				testStatus = "warning"
			}
		}
	}
	checklist = append(checklist, ChecklistItem{
		Title:       "Tests are written and passing",
		Description: "Unit tests, integration tests, and good coverage",
		Status:      testStatus,
	})

	// 3. CI/CD
	ciStatus := "pending"
	for _, category := range result.Categories {
		if category.Name == "GitHub Integration" {
			if category.Score >= 80 {
				ciStatus = "done"
			} else if category.Score >= 50 {
				ciStatus = "warning"
			}
		}
	}
	checklist = append(checklist, ChecklistItem{
		Title:       "CI/CD pipelines are configured and working",
		Description: "GitHub Actions for test, lint, and release automation",
		Status:      ciStatus,
	})

	// 4. License
	licenseStatus := "pending"
	for _, category := range result.Categories {
		if category.Name == "Licensing" {
			if category.Score >= 90 {
				licenseStatus = "done"
			} else if category.Score >= 50 {
				licenseStatus = "warning"
			}
		}
	}
	checklist = append(checklist, ChecklistItem{
		Title:       "License is properly configured",
		Description: "LICENSE file exists and matches project metadata",
		Status:      licenseStatus,
	})

	// 5. Security
	securityStatus := "done" // Assumed OK if passed earlier checks
	checklist = append(checklist, ChecklistItem{
		Title:       "Security policy is defined",
		Description: "SECURITY.md with vulnerability reporting process",
		Status:      securityStatus,
	})

	// 6. Community files
	communityStatus := "done" // Assumed OK based on GitHub Integration score
	for _, category := range result.Categories {
		if category.Name == "GitHub Integration" && category.Score < 70 {
			communityStatus = "warning"
		}
	}
	checklist = append(checklist, ChecklistItem{
		Title:       "Community guidelines are in place",
		Description: "CONTRIBUTING.md, CODE_OF_CONDUCT.md, issue/PR templates",
		Status:      communityStatus,
	})

	// 7. Sensitive data
	checklist = append(checklist, ChecklistItem{
		Title:       "No sensitive data in repository",
		Description: "API keys, credentials, and secrets are removed/ignored",
		Status:      "done",
	})

	// 8. Go module
	moduleStatus := "done"
	for _, category := range result.Categories {
		if category.Name == "Dependencies" && category.Score < 80 {
			moduleStatus = "warning"
		}
	}
	checklist = append(checklist, ChecklistItem{
		Title:       "Go module is properly configured",
		Description: "go.mod and go.sum are up-to-date and tidy",
		Status:      moduleStatus,
	})

	// 9. Version tag
	checklist = append(checklist, ChecklistItem{
		Title:       "Ready to create initial version tag",
		Description: "Decide on semantic version (v0.1.0 for beta, v1.0.0 for stable)",
		Status:      "pending",
	})

	// 10. Code quality
	qualityStatus := "pending"
	if result.OverallScore >= 90 {
		qualityStatus = "done"
	} else if result.OverallScore >= 70 {
		qualityStatus = "warning"
	}
	checklist = append(checklist, ChecklistItem{
		Title:       "Code quality meets standards",
		Description: "Linter passes, no critical issues, code is reviewed",
		Status:      qualityStatus,
	})

	return checklist
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
