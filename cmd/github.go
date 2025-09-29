package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/pigeonworks-llc/goossify/internal/github"
	"github.com/spf13/cobra"
)

var (
	githubToken string
	dryRun      bool
)

var githubCmd = &cobra.Command{
	Use:   "github",
	Short: "Automate GitHub repository settings",
	Long: `Automate repository settings using the GitHub API.

This command configures:
• Branch protection rules
• Label settings
• General repository settings

Requires a GitHub Personal Access Token.`,
}

var githubSetupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Execute basic GitHub repository setup",
	RunE:  runGitHubSetup,
}

func runGitHubSetup(cmd *cobra.Command, args []string) error {
	// Check GitHub token
	token := githubToken
	if token == "" {
		token = os.Getenv("GITHUB_TOKEN")
	}
	if token == "" {
		return fmt.Errorf("GitHub token is required. Set --token flag or GITHUB_TOKEN environment variable")
	}

	// Get owner/repo from Git remote URL
	owner, repo, err := getRepositoryInfo()
	if err != nil {
		return fmt.Errorf("failed to get repository information: %w", err)
	}

	fmt.Printf("🔧 Starting GitHub repository setup: %s/%s\n", owner, repo)

	if dryRun {
		fmt.Println("🔍 Dry-run mode: No actual changes will be made")
		return nil
	}

	// Create GitHub client
	client, err := github.NewClient(github.Config{
		Token: token,
		Owner: owner,
		Repo:  repo,
	})
	if err != nil {
		return fmt.Errorf("failed to create GitHub client: %w", err)
	}

	// Execute with default settings
	settings := github.RepositorySettings{
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

	if err := client.SetupRepository(settings); err != nil {
		return fmt.Errorf("failed to setup repository: %w", err)
	}

	fmt.Println("✅ GitHub repository setup completed")
	return nil
}

func getRepositoryInfo() (owner, repo string, err error) {
	// git remote get-url origin
	cmd := exec.Command("git", "remote", "get-url", "origin")
	output, err := cmd.Output()
	if err != nil {
		return "", "", fmt.Errorf("failed to get git remote URL: %w", err)
	}

	remoteURL := strings.TrimSpace(string(output))
	return github.ParseRepositoryURL(remoteURL)
}

func init() {
	rootCmd.AddCommand(githubCmd)
	githubCmd.AddCommand(githubSetupCmd)

	githubSetupCmd.Flags().StringVar(&githubToken, "token", "", "GitHub Personal Access Token")
	githubSetupCmd.Flags().BoolVar(&dryRun, "dry-run", false, "Show settings without making actual changes")
}
