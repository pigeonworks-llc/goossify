package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
	"github.com/pigeonworks-llc/goossify/internal/github"
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
	// GitHub token確認
	token := githubToken
	if token == "" {
		token = os.Getenv("GITHUB_TOKEN")
	}
	if token == "" {
		return fmt.Errorf("GitHub token が必要です。--token フラグか GITHUB_TOKEN 環境変数を設定してください")
	}

	// Git remote URLからowner/repo取得
	owner, repo, err := getRepositoryInfo()
	if err != nil {
		return fmt.Errorf("リポジトリ情報取得失敗: %w", err)
	}

	fmt.Printf("🔧 GitHub リポジトリ設定を開始します: %s/%s\n", owner, repo)

	if dryRun {
		fmt.Println("🔍 ドライランモード: 実際の設定は行いません")
		return nil
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

	// デフォルト設定で実行
	settings := github.RepositorySettings{
		BranchProtection: github.BranchProtectionSettings{
			Branch:                  "main",
			RequiredStatusChecks:    []string{"test", "lint"},
			RequiredReviews:         1,
			DismissStaleReviews:     true,
			RequireCodeOwnerReviews: false,
			RestrictPushes:          false,
		},
		Labels:                 github.GetDefaultLabels(),
		DeleteBranchOnMerge:    true,
	}

	if err := client.SetupRepository(&settings); err != nil {
		return fmt.Errorf("リポジトリ設定失敗: %w", err)
	}

	fmt.Println("✅ GitHub リポジトリ設定が完了しました")
	return nil
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