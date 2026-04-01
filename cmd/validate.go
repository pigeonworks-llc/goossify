package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/pigeonworks-llc/goossify/internal/analyzer"
	"github.com/spf13/cobra"
)

var (
	validateFormat string
	validateStrict bool
)

// validateCmd はプロジェクトのベストプラクティス準拠を検証
var validateCmd = &cobra.Command{
	Use:   "validate [path]",
	Short: "Validate OSS best practices compliance",
	Long: `Validate if the project complies with OSS best practices.

This command performs comprehensive checks including:
- OpenSSF Best Practices recommendations
- Linux Foundation OSS guidelines
- Community health files
- Security configurations
- Automation setup`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		// プロジェクトパス決定
		projectPath := "."
		if len(args) > 0 {
			projectPath = args[0]
		}

		absPath, err := filepath.Abs(projectPath)
		if err != nil {
			return fmt.Errorf("パス解決失敗: %w", err)
		}

		// プロジェクト分析
		fmt.Println("🔍 OSS Best Practices検証を開始...")
		fmt.Printf("📁 プロジェクト: %s\n\n", absPath)

		projectAnalyzer := analyzer.New(absPath)
		result, err := projectAnalyzer.Analyze()
		if err != nil {
			return fmt.Errorf("分析失敗: %w", err)
		}

		// OpenSSFレベル判定
		openSSFLevel := getOpenSSFLevel(result.OverallScore)

		// 出力形式による分岐
		switch validateFormat {
		case "json":
			return outputValidationJSON(result, openSSFLevel)
		default:
			return outputValidationText(result, openSSFLevel)
		}
	},
}

func init() {
	rootCmd.AddCommand(validateCmd)

	validateCmd.Flags().StringVar(&validateFormat, "format", "text", "出力形式 (text, json)")
	validateCmd.Flags().BoolVar(&validateStrict, "strict", false, "厳格モード（100点未満でエラー）")
}

// getOpenSSFLevel はスコアからOpenSSFレベルを判定
func getOpenSSFLevel(score int) string {
	if score >= 95 {
		return "🥇 Gold (Excellent)"
	} else if score >= 85 {
		return "🥈 Silver (Good)"
	} else if score >= 75 {
		return "🥉 Bronze (Passing)"
	}
	return "❌ Not Passing"
}

// outputValidationText はテキスト形式で出力
func outputValidationText(result *analyzer.AnalysisResult, openSSFLevel string) error {
	fmt.Println("=" + string(make([]byte, 70)) + "=")
	fmt.Printf("📊 OSS Best Practices Validation Report\n")
	fmt.Println("=" + string(make([]byte, 70)) + "=")
	fmt.Println()

	// 総合スコア
	fmt.Printf("🎯 Overall Score: %d/100\n", result.OverallScore)
	fmt.Printf("🏆 OpenSSF Level: %s\n", openSSFLevel)
	fmt.Println()

	// カテゴリ別結果
	fmt.Println("📋 Category Breakdown:")
	fmt.Println()

	for _, category := range result.Categories {
		statusIcon := getStatusIcon(category.Status)
		fmt.Printf("  %s %s: %d/100 (%s)\n", statusIcon, category.Name, category.Score, category.Status)

		// 詳細項目（オプション）
		if verbose {
			for _, item := range category.Items {
				var itemIcon string
				switch item.Status {
				case "missing":
					itemIcon = "❌"
				case "warning":
					itemIcon = "⚠️"
				default:
					itemIcon = "✅"
				}
				fmt.Printf("    %s %s\n", itemIcon, item.Name)
			}
			fmt.Println()
		}
	}

	// 不足項目
	if len(result.Missing) > 0 {
		fmt.Println("\n❌ Missing Items:")
		fmt.Println()
		for _, missing := range result.Missing {
			priorityIcon := getPriorityIcon(missing.Priority)
			fmt.Printf("  %s %s (%s)\n", priorityIcon, missing.Name, missing.Category)
			fmt.Printf("     └─ %s\n", missing.Description)
		}
	}

	// 推奨事項
	if len(result.Recommendations) > 0 {
		fmt.Println("\n💡 Recommendations:")
		fmt.Println()
		for i, rec := range result.Recommendations {
			priorityIcon := getPriorityIcon(rec.Priority)
			fmt.Printf("  %d. %s %s\n", i+1, priorityIcon, rec.Title)
			fmt.Printf("     %s\n", rec.Description)
			if rec.Command != "" {
				fmt.Printf("     💻 Command: %s\n", rec.Command)
			}
			fmt.Println()
		}
	}

	// 結論
	fmt.Println("\n" + string(make([]byte, 70)))
	if result.OverallScore >= 90 {
		fmt.Println("✅ 素晴らしい！プロジェクトはOSSベストプラクティスに準拠しています。")
	} else if result.OverallScore >= 75 {
		fmt.Println("⚠️  良好です。いくつかの改善でさらに良くなります。")
	} else {
		fmt.Println("❌ 改善が必要です。不足項目を確認してください。")
	}

	// 厳格モード
	if validateStrict && result.OverallScore < 100 {
		fmt.Println("\n⚠️  厳格モードが有効です。100点未満のためエラーを返します。")
		os.Exit(1)
	}

	return nil
}

// outputValidationJSON はJSON形式で出力
func outputValidationJSON(result *analyzer.AnalysisResult, openSSFLevel string) error {
	output := map[string]interface{}{
		"overall_score":   result.OverallScore,
		"openssf_level":   openSSFLevel,
		"project_name":    result.ProjectName,
		"project_type":    result.ProjectType,
		"categories":      result.Categories,
		"missing":         result.Missing,
		"recommendations": result.Recommendations,
		"summary":         result.Summary,
	}

	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	return encoder.Encode(output)
}
