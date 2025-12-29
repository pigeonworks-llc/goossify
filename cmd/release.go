package cmd

import (
	"fmt"
	"os"

	"github.com/pigeonworks-llc/goossify/internal/release"
	"github.com/spf13/cobra"
)

var (
	releaseVersion   string
	releaseDryRun    bool
	releaseSkipTests bool
	releaseSkipLint  bool
	releaseSkipBuild bool
	releaseBump      string
)

// releaseCmd represents the release command
var releaseCmd = &cobra.Command{
	Use:   "release [--tag version]",
	Short: "Prepare a new release",
	Long: `Prepare a new release for your Go OSS project.

This command performs the following:
1. Validates the version format (semantic versioning)
2. Runs tests (can be skipped with --skip-tests)
3. Runs linter (can be skipped with --skip-lint)
4. Builds the project (can be skipped with --skip-build)
5. Updates CHANGELOG.md with commits since last tag
6. Updates version files (if found)
7. Creates an annotated git tag

Examples:
  goossify release --tag v1.0.0
  goossify release --tag v1.0.0 --dry-run
  goossify release --bump patch
  goossify release --bump minor
  goossify release --bump major`,
	RunE: runRelease,
}

func init() {
	rootCmd.AddCommand(releaseCmd)

	releaseCmd.Flags().StringVar(&releaseVersion, "tag", "", "Version tag (e.g., v1.0.0)")
	releaseCmd.Flags().BoolVar(&releaseDryRun, "dry-run", false, "Preview changes without making them")
	releaseCmd.Flags().BoolVar(&releaseSkipTests, "skip-tests", false, "Skip running tests")
	releaseCmd.Flags().BoolVar(&releaseSkipLint, "skip-lint", false, "Skip running linter")
	releaseCmd.Flags().BoolVar(&releaseSkipBuild, "skip-build", false, "Skip building the project")
	releaseCmd.Flags().StringVar(&releaseBump, "bump", "", "Auto-bump version (patch, minor, major)")
}

func runRelease(cmd *cobra.Command, args []string) error {
	// Get project path
	projectPath, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %w", err)
	}

	// Determine version
	version := releaseVersion
	if version == "" && releaseBump != "" {
		// Auto-bump version
		suggested, err := release.SuggestNextVersion(projectPath)
		if err != nil {
			return fmt.Errorf("failed to suggest version: %w", err)
		}

		switch releaseBump {
		case "patch", "minor", "major":
			version = bumpVersionFromSuggested(suggested, releaseBump)
		default:
			return fmt.Errorf("invalid bump type: %s (use patch, minor, or major)", releaseBump)
		}

		fmt.Printf("Auto-bumping to version: %s\n", version)
	}

	if version == "" {
		// Suggest next version
		suggested, _ := release.SuggestNextVersion(projectPath)
		fmt.Println("No version specified.")
		fmt.Printf("Suggested next version: %s\n", suggested)
		fmt.Println("\nUsage:")
		fmt.Println("  goossify release --tag v1.0.0")
		fmt.Println("  goossify release --bump patch")
		return nil
	}

	// Ensure version starts with 'v'
	if version[0] != 'v' {
		version = "v" + version
	}

	// Create release config
	config := &release.Config{
		ProjectPath: projectPath,
		Version:     version,
		DryRun:      releaseDryRun,
		SkipTests:   releaseSkipTests,
		SkipLint:    releaseSkipLint,
		SkipBuild:   releaseSkipBuild,
	}

	// Execute release
	manager := release.New(config)
	return manager.Execute()
}

func bumpVersionFromSuggested(suggested, bumpType string) string {
	// The suggested version is already the next patch version
	// For minor/major, we need to adjust
	switch bumpType {
	case "minor":
		return adjustBump(suggested, "minor")
	case "major":
		return adjustBump(suggested, "major")
	default:
		return suggested
	}
}

func adjustBump(version, bumpType string) string {
	// Parse version
	var major, minor, patch int
	version = version[1:] // Remove 'v' prefix
	fmt.Sscanf(version, "%d.%d.%d", &major, &minor, &patch)

	switch bumpType {
	case "minor":
		// Reset patch from suggested
		if patch > 0 {
			patch = 0
			minor++
		}
	case "major":
		// Reset minor and patch from suggested
		patch = 0
		minor = 0
		major++
	}

	return fmt.Sprintf("v%d.%d.%d", major, minor, patch)
}
