package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/pigeonworks-llc/goossify/internal/ossify"
	"github.com/spf13/cobra"
)

var (
	interactiveMode bool
	dryRunMode      bool
)

var ossifyCmd = &cobra.Command{
	Use:   "ossify [path]",
	Short: "既存プロジェクトをOSS対応に変換",
	Long: `既存のGoプロジェクトをOSS公開に必要なファイルとセットアップで拡張します。

このコマンドは以下を自動生成・設定します：
• LICENSE ファイル
• .github/workflows/ci.yml (CI/CD設定)
• .github/ 内のコミュニティファイル
• CONTRIBUTING.md, SECURITY.md等
• Git初期化（未初期化の場合）`,
	Args: cobra.MaximumNArgs(1),
	RunE: runOssify,
}

func runOssify(cmd *cobra.Command, args []string) error {
	var targetPath string
	if len(args) == 0 {
		targetPath = "."
	} else {
		targetPath = args[0]
	}

	// パスを絶対パスに変換
	absPath, err := filepath.Abs(targetPath)
	if err != nil {
		return fmt.Errorf("パス解決エラー: %w", err)
	}

	// ディレクトリの存在確認
	if _, err := os.Stat(absPath); os.IsNotExist(err) {
		return fmt.Errorf("指定されたパスが存在しません: %s", absPath)
	}

	fmt.Printf("🚀 OSS化を開始します: %s\n", absPath)

	// Ossifierを初期化して実行
	ossifier := ossify.New(absPath)
	ossifier.SetInteractive(interactiveMode)
	ossifier.SetDryRun(dryRunMode)

	if err := ossifier.Execute(); err != nil {
		return fmt.Errorf("OSS化処理中にエラーが発生しました: %w", err)
	}

	if dryRunMode {
		fmt.Println("✅ Dry-runモード: 実際のファイル変更は行われませんでした")
	} else {
		fmt.Println("✅ OSS化が完了しました！")
	}
	return nil
}

func init() {
	rootCmd.AddCommand(ossifyCmd)
	ossifyCmd.Flags().BoolVarP(&interactiveMode, "interactive", "i", false, "対話的に確認しながら適用")
	ossifyCmd.Flags().BoolVar(&dryRunMode, "dry-run", false, "実際には適用せず、変更内容のみ表示")
}
