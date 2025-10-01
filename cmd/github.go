package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/pigeonworks-llc/goossify/internal/github"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	githubToken string
	dryRun      bool
)

var githubCmd = &cobra.Command{
	Use:   "github",
	Short: "GitHub リポジトリ設定の自動化",
	Long: `GitHub API を使用してリポジトリの設定を自動化します。

このコマンドは以下を設定します：
• ブランチ保護ルール
• ラベル設定
• リポジトリ一般設定

GitHub Personal Access Token が必要です。`,
}

var githubSetupCmd = &cobra.Command{
	Use:   "setup",
	Short: "GitHub リポジトリの基本設定を実行",
	RunE:  runGitHubSetup,
}

func runGitHubSetup(cmd *cobra.Command, args []string) error {
	// Git remote URLからowner/repo取得
	owner, repo, err := getRepositoryInfo()
	if err != nil {
		return fmt.Errorf("リポジトリ情報取得失敗: %w", err)
	}

	fmt.Printf("🔧 GitHub リポジトリ設定を開始します: %s/%s\n", owner, repo)

	// .goossify.yml から設定を読み込み
	settings, err := loadGitHubSettings()
	if err != nil {
		fmt.Printf("⚠️  .goossify.yml の読み込みに失敗しました。デフォルト設定を使用します: %v\n", err)
		settings = getDefaultSettings()
	} else {
		fmt.Println("✅ .goossify.yml から設定を読み込みました")
	}

	if dryRun {
		fmt.Println("\n🔍 ドライランモード: 以下の設定を適用します（実際の変更は行いません）")
		printSettings(settings)
		return nil
	}

	// GitHub token確認（dry-runでない場合のみ）
	token := githubToken
	if token == "" {
		token = os.Getenv("GITHUB_TOKEN")
	}
	if token == "" {
		return fmt.Errorf("GitHub token が必要です。--token フラグか GITHUB_TOKEN 環境変数を設定してください")
	}

	// GitHub クライアント作成
	client, err := github.NewClient(github.Config{
		Token: token,
		Owner: owner,
		Repo:  repo,
	})
	if err != nil {
		return fmt.Errorf("GitHub クライアント作成失敗: %w", err)
	}

	// 設定実行
	fmt.Println("\n📋 設定を適用中...")
	if err := client.SetupRepository(settings); err != nil {
		return fmt.Errorf("リポジトリ設定失敗: %w", err)
	}

	fmt.Println("✅ GitHub リポジトリ設定が完了しました")
	fmt.Println("\n設定内容:")
	printSettings(settings)
	return nil
}

// loadGitHubSettings は .goossify.yml から GitHub 設定を読み込む
func loadGitHubSettings() (*github.RepositorySettings, error) {
	configPath := filepath.Join(".", ".goossify.yml")
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return nil, fmt.Errorf(".goossify.yml が見つかりません")
	}

	v := viper.New()
	v.SetConfigFile(configPath)
	v.SetConfigType("yaml")

	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("設定ファイル読み込み失敗: %w", err)
	}

	// GitHub設定を読み込み
	branchProtection := v.GetBool("integrations.github.branch_protection")
	requiredReviews := v.GetInt("integrations.github.required_reviews")
	statusChecks := v.GetStringSlice("integrations.github.status_checks")
	deleteBranchOnMerge := v.GetBool("integrations.github.delete_branch_on_merge")

	// デフォルト値
	if requiredReviews == 0 {
		requiredReviews = 1
	}
	if len(statusChecks) == 0 {
		statusChecks = []string{"test", "lint"}
	}

	settings := &github.RepositorySettings{
		Labels:              github.GetDefaultLabels(),
		DeleteBranchOnMerge: deleteBranchOnMerge,
	}

	if branchProtection {
		settings.BranchProtection = github.BranchProtectionSettings{
			Branch:                  "main",
			RequiredStatusChecks:    statusChecks,
			RequiredReviews:         requiredReviews,
			DismissStaleReviews:     true,
			RequireCodeOwnerReviews: false,
			RestrictPushes:          false,
		}
	}

	return settings, nil
}

// getDefaultSettings はデフォルトの GitHub 設定を返す
func getDefaultSettings() *github.RepositorySettings {
	return &github.RepositorySettings{
		BranchProtection: github.BranchProtectionSettings{
			Branch:                  "main",
			RequiredStatusChecks:    []string{"test", "lint"},
			RequiredReviews:         1,
			DismissStaleReviews:     true,
			RequireCodeOwnerReviews: false,
			RestrictPushes:          false,
		},
		Labels:              github.GetDefaultLabels(),
		DeleteBranchOnMerge: true,
	}
}

// printSettings は設定内容を表示
func printSettings(settings *github.RepositorySettings) {
	fmt.Println("  📌 ブランチ保護:")
	if settings.BranchProtection.Branch != "" {
		fmt.Printf("     - ブランチ: %s\n", settings.BranchProtection.Branch)
		fmt.Printf("     - 必須レビュー数: %d\n", settings.BranchProtection.RequiredReviews)
		fmt.Printf("     - 必須ステータスチェック: %v\n", settings.BranchProtection.RequiredStatusChecks)
		fmt.Printf("     - 古いレビューの却下: %v\n", settings.BranchProtection.DismissStaleReviews)
	} else {
		fmt.Println("     - 無効")
	}
	fmt.Printf("  🏷️  ラベル数: %d\n", len(settings.Labels))
	fmt.Printf("  🗑️  マージ後のブランチ削除: %v\n", settings.DeleteBranchOnMerge)
}

func getRepositoryInfo() (owner, repo string, err error) {
	// git remote get-url origin
	cmd := exec.Command("git", "remote", "get-url", "origin")
	output, err := cmd.Output()
	if err != nil {
		return "", "", fmt.Errorf("git remote URL 取得失敗: %w", err)
	}

	remoteURL := strings.TrimSpace(string(output))
	return github.ParseRepositoryURL(remoteURL)
}

func init() {
	rootCmd.AddCommand(githubCmd)
	githubCmd.AddCommand(githubSetupCmd)

	githubSetupCmd.Flags().StringVar(&githubToken, "token", "", "GitHub Personal Access Token")
	githubSetupCmd.Flags().BoolVar(&dryRun, "dry-run", false, "設定内容を表示するのみで実際の変更は行わない")
}
