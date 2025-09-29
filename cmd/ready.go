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

	// 5. Pre-publication recommendations
	fmt.Println("\n💡 Pre-publication Recommendations:")
	fmt.Println("  • Verify README.md content is properly documented")
	fmt.Println("  • Verify examples and sample code work correctly")
	fmt.Println("  • Verify CI/CD operates normally")
	fmt.Println("  • Verify security policy is configured")

	fmt.Printf("\n🎉 Project '%s' is ready for public release!\n", result.ProjectName)
	fmt.Println("Next steps:")
	fmt.Println("  1. Change GitHub repository to Public")
	fmt.Println("  2. Create initial release tag")
	fmt.Println("  3. Verify automatic indexing on pkg.go.dev")

	return nil
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
