// Package release provides release management functionality for Go OSS projects.
package release

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"
)

// Config holds the release configuration
type Config struct {
	ProjectPath string
	Version     string
	DryRun      bool
	SkipTests   bool
	SkipLint    bool
	SkipBuild   bool
}

// Manager handles release operations
type Manager struct {
	config *Config
}

// New creates a new release manager
func New(config *Config) *Manager {
	return &Manager{config: config}
}

// Execute performs the release process
func (m *Manager) Execute() error {
	fmt.Printf("🚀 Starting release process for version %s\n\n", m.config.Version)

	// Validate version format
	if !isValidSemver(m.config.Version) {
		return fmt.Errorf("invalid version format: %s (expected: vX.Y.Z)", m.config.Version)
	}

	// Check prerequisites
	if err := m.checkPrerequisites(); err != nil {
		return err
	}

	// Run pre-release checks
	if !m.config.SkipTests {
		if err := m.runTests(); err != nil {
			return err
		}
	}

	if !m.config.SkipLint {
		if err := m.runLint(); err != nil {
			return err
		}
	}

	if !m.config.SkipBuild {
		if err := m.runBuild(); err != nil {
			return err
		}
	}

	// Update CHANGELOG.md
	if err := m.updateChangelog(); err != nil {
		return err
	}

	// Update version in source files
	if err := m.updateVersionFiles(); err != nil {
		return err
	}

	// Create git tag
	if err := m.createGitTag(); err != nil {
		return err
	}

	if m.config.DryRun {
		fmt.Println("\n📋 Dry-run completed. No changes were made.")
		return nil
	}

	fmt.Printf("\n🎉 Release %s prepared successfully!\n", m.config.Version)
	fmt.Println("\nNext steps:")
	fmt.Println("  1. Review the changes: git diff HEAD~1")
	fmt.Println("  2. Push changes: git push && git push --tags")
	fmt.Println("  3. Create GitHub release or wait for CI/CD automation")

	return nil
}

func (m *Manager) checkPrerequisites() error {
	fmt.Println("📋 Checking prerequisites...")

	// Check if git is available
	if _, err := exec.LookPath("git"); err != nil {
		return fmt.Errorf("git is not installed or not in PATH")
	}

	// Check if in git repository
	cmd := exec.Command("git", "rev-parse", "--git-dir")
	cmd.Dir = m.config.ProjectPath
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("not a git repository")
	}

	// Check for uncommitted changes
	cmd = exec.Command("git", "status", "--porcelain")
	cmd.Dir = m.config.ProjectPath
	output, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("failed to check git status: %w", err)
	}

	if strings.TrimSpace(string(output)) != "" && !m.config.DryRun {
		return fmt.Errorf("there are uncommitted changes. Please commit or stash them first")
	}

	// Check if tag already exists
	cmd = exec.Command("git", "tag", "-l", m.config.Version)
	cmd.Dir = m.config.ProjectPath
	output, err = cmd.Output()
	if err != nil {
		return fmt.Errorf("failed to check existing tags: %w", err)
	}

	if strings.TrimSpace(string(output)) == m.config.Version {
		return fmt.Errorf("tag %s already exists", m.config.Version)
	}

	fmt.Println("  ✅ All prerequisites met")
	return nil
}

func (m *Manager) runTests() error {
	fmt.Println("\n🧪 Running tests...")

	cmd := exec.Command("go", "test", "-race", "-cover", "./...")
	cmd.Dir = m.config.ProjectPath
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("tests failed: %w", err)
	}

	fmt.Println("  ✅ All tests passed")
	return nil
}

func (m *Manager) runLint() error {
	fmt.Println("\n🔍 Running linter...")

	// Check if golangci-lint is available
	if _, err := exec.LookPath("golangci-lint"); err != nil {
		fmt.Println("  ⚠️  golangci-lint not found, skipping lint check")
		return nil
	}

	cmd := exec.Command("golangci-lint", "run", "./...")
	cmd.Dir = m.config.ProjectPath
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("lint check failed: %w", err)
	}

	fmt.Println("  ✅ Lint check passed")
	return nil
}

func (m *Manager) runBuild() error {
	fmt.Println("\n🔨 Building project...")

	cmd := exec.Command("go", "build", "-v", "./...")
	cmd.Dir = m.config.ProjectPath
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("build failed: %w", err)
	}

	fmt.Println("  ✅ Build successful")
	return nil
}

func (m *Manager) updateChangelog() error {
	fmt.Println("\n📝 Updating CHANGELOG.md...")

	changelogPath := filepath.Join(m.config.ProjectPath, "CHANGELOG.md")

	// Read existing changelog
	content, err := os.ReadFile(changelogPath)
	if err != nil {
		if os.IsNotExist(err) {
			// Create new changelog
			return m.createChangelog()
		}
		return fmt.Errorf("failed to read CHANGELOG.md: %w", err)
	}

	// Get commits since last tag
	commits, err := m.getCommitsSinceLastTag()
	if err != nil {
		fmt.Printf("  ⚠️  Could not get commits: %v\n", err)
		commits = []string{}
	}

	// Generate new release section
	releaseSection := m.generateReleaseSection(commits)

	// Insert new release section after ## [Unreleased]
	newContent := insertReleaseSection(string(content), m.config.Version, releaseSection)

	if m.config.DryRun {
		fmt.Println("  [DRY-RUN] Would update CHANGELOG.md with:")
		fmt.Println(releaseSection)
		return nil
	}

	if err := os.WriteFile(filepath.Clean(changelogPath), []byte(newContent), 0644); err != nil {
		return fmt.Errorf("failed to write CHANGELOG.md: %w", err)
	}

	fmt.Println("  ✅ CHANGELOG.md updated")
	return nil
}

func (m *Manager) createChangelog() error {
	changelogPath := filepath.Join(m.config.ProjectPath, "CHANGELOG.md")

	content := fmt.Sprintf(`# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [%s] - %s

### Added
- Initial release

`, strings.TrimPrefix(m.config.Version, "v"), time.Now().Format("2006-01-02"))

	if m.config.DryRun {
		fmt.Println("  [DRY-RUN] Would create CHANGELOG.md")
		return nil
	}

	return os.WriteFile(filepath.Clean(changelogPath), []byte(content), 0644)
}

func (m *Manager) getCommitsSinceLastTag() ([]string, error) {
	// Get last tag
	cmd := exec.Command("git", "describe", "--tags", "--abbrev=0")
	cmd.Dir = m.config.ProjectPath
	output, err := cmd.Output()
	if err != nil {
		// No tags yet, get all commits
		cmd = exec.Command("git", "log", "--oneline", "--no-decorate", "-50")
		cmd.Dir = m.config.ProjectPath
		output, err = cmd.Output()
		if err != nil {
			return nil, err
		}
		return parseCommitLines(string(output)), nil
	}

	lastTag := strings.TrimSpace(string(output))

	// Get commits since last tag
	cmd = exec.Command("git", "log", fmt.Sprintf("%s..HEAD", lastTag), "--oneline", "--no-decorate")
	cmd.Dir = m.config.ProjectPath
	output, err = cmd.Output()
	if err != nil {
		return nil, err
	}

	return parseCommitLines(string(output)), nil
}

func parseCommitLines(output string) []string {
	var commits []string
	scanner := bufio.NewScanner(strings.NewReader(output))
	for scanner.Scan() {
		line := scanner.Text()
		if line != "" {
			// Remove commit hash prefix
			parts := strings.SplitN(line, " ", 2)
			if len(parts) == 2 {
				commits = append(commits, parts[1])
			}
		}
	}
	return commits
}

func (m *Manager) generateReleaseSection(commits []string) string {
	var features, fixes, others []string

	for _, commit := range commits {
		lower := strings.ToLower(commit)
		if strings.HasPrefix(lower, "feat") || strings.HasPrefix(lower, "add") {
			features = append(features, fmt.Sprintf("- %s", commit))
		} else if strings.HasPrefix(lower, "fix") {
			fixes = append(fixes, fmt.Sprintf("- %s", commit))
		} else if !strings.HasPrefix(lower, "merge") && !strings.HasPrefix(lower, "chore") {
			others = append(others, fmt.Sprintf("- %s", commit))
		}
	}

	var sb strings.Builder
	fmt.Fprintf(&sb, "## [%s] - %s\n\n",
		strings.TrimPrefix(m.config.Version, "v"),
		time.Now().Format("2006-01-02"))

	if len(features) > 0 {
		sb.WriteString("### Added\n")
		sort.Strings(features)
		for _, f := range features {
			sb.WriteString(f + "\n")
		}
		sb.WriteString("\n")
	}

	if len(fixes) > 0 {
		sb.WriteString("### Fixed\n")
		sort.Strings(fixes)
		for _, f := range fixes {
			sb.WriteString(f + "\n")
		}
		sb.WriteString("\n")
	}

	if len(others) > 0 {
		sb.WriteString("### Changed\n")
		sort.Strings(others)
		for _, o := range others {
			sb.WriteString(o + "\n")
		}
		sb.WriteString("\n")
	}

	return sb.String()
}

func insertReleaseSection(content, _, releaseSection string) string {
	// Find ## [Unreleased] and insert after it
	unreleasedPattern := regexp.MustCompile(`(?m)^## \[Unreleased\].*\n`)
	loc := unreleasedPattern.FindStringIndex(content)

	if loc != nil {
		// Insert after ## [Unreleased] section
		insertPos := loc[1]

		// Skip any content until next ## heading or end
		nextSectionPattern := regexp.MustCompile(`(?m)^## \[`)
		remaining := content[insertPos:]
		nextLoc := nextSectionPattern.FindStringIndex(remaining)

		if nextLoc != nil {
			insertPos += nextLoc[0]
		}

		return content[:insertPos] + "\n" + releaseSection + content[insertPos:]
	}

	// No [Unreleased] section, prepend after header
	headerPattern := regexp.MustCompile(`(?m)^# .*\n`)
	loc = headerPattern.FindStringIndex(content)
	if loc != nil {
		insertPos := loc[1]
		return content[:insertPos] + "\n## [Unreleased]\n\n" + releaseSection + content[insertPos:]
	}

	// No header, prepend
	return "# Changelog\n\n## [Unreleased]\n\n" + releaseSection + content
}

func (m *Manager) updateVersionFiles() error {
	fmt.Println("\n📌 Updating version files...")

	// Common version file patterns
	versionFiles := []struct {
		pattern string
		updater func(path, version string) error
	}{
		{"internal/version/version.go", m.updateGoVersionFile},
		{"version.go", m.updateGoVersionFile},
		{"cmd/version.go", m.updateGoVersionFile},
	}

	updated := false
	for _, vf := range versionFiles {
		path := filepath.Join(m.config.ProjectPath, vf.pattern)
		if _, err := os.Stat(path); err == nil {
			if err := vf.updater(path, m.config.Version); err != nil {
				fmt.Printf("  ⚠️  Failed to update %s: %v\n", vf.pattern, err)
			} else {
				fmt.Printf("  ✅ Updated %s\n", vf.pattern)
				updated = true
			}
		}
	}

	if !updated {
		fmt.Println("  ℹ️  No version files found to update")
	}

	return nil
}

func (m *Manager) updateGoVersionFile(path, version string) error {
	content, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	// Match patterns like: Version = "0.1.0" or const Version = "0.1.0"
	versionPattern := regexp.MustCompile(`(Version\s*=\s*")[^"]+"`)
	cleanVersion := strings.TrimPrefix(version, "v")

	newContent := versionPattern.ReplaceAllString(string(content), fmt.Sprintf("${1}%s\"", cleanVersion))

	if m.config.DryRun {
		if string(content) != newContent {
			fmt.Printf("  [DRY-RUN] Would update version to %s in %s\n", cleanVersion, path)
		}
		return nil
	}

	return os.WriteFile(filepath.Clean(path), []byte(newContent), 0644)
}

func (m *Manager) createGitTag() error {
	fmt.Printf("\n🏷️  Creating git tag %s...\n", m.config.Version)

	if m.config.DryRun {
		fmt.Println("  [DRY-RUN] Would create git tag")
		return nil
	}

	// Stage changes
	cmd := exec.Command("git", "add", "-A")
	cmd.Dir = m.config.ProjectPath
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to stage changes: %w", err)
	}

	// Commit changes
	commitMsg := fmt.Sprintf("chore(release): prepare release %s", m.config.Version)
	cmd = exec.Command("git", "commit", "-m", commitMsg)
	cmd.Dir = m.config.ProjectPath
	if err := cmd.Run(); err != nil {
		// Ignore error if nothing to commit
		fmt.Println("  ℹ️  No changes to commit")
	}

	// Create annotated tag
	tagMsg := fmt.Sprintf("Release %s", m.config.Version)
	cmd = exec.Command("git", "tag", "-a", m.config.Version, "-m", tagMsg)
	cmd.Dir = m.config.ProjectPath
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to create tag: %w", err)
	}

	fmt.Printf("  ✅ Created tag %s\n", m.config.Version)
	return nil
}

func isValidSemver(version string) bool {
	// Match v1.0.0, v1.0.0-alpha, v1.0.0-beta.1, etc.
	pattern := regexp.MustCompile(`^v?(\d+)\.(\d+)\.(\d+)(-[a-zA-Z0-9.-]+)?$`)
	return pattern.MatchString(version)
}

// SuggestNextVersion suggests the next version based on the current tag
func SuggestNextVersion(projectPath string) (string, error) {
	cmd := exec.Command("git", "describe", "--tags", "--abbrev=0")
	cmd.Dir = projectPath
	output, err := cmd.Output()
	if err != nil {
		return "v0.1.0", nil // Default initial version
	}

	currentVersion := strings.TrimSpace(string(output))
	return bumpVersion(currentVersion, "patch"), nil
}

func bumpVersion(version, bumpType string) string {
	version = strings.TrimPrefix(version, "v")
	parts := strings.Split(version, ".")

	if len(parts) != 3 {
		return "v0.1.0"
	}

	var major, minor, patch int
	_, _ = fmt.Sscanf(parts[0], "%d", &major)
	_, _ = fmt.Sscanf(parts[1], "%d", &minor)
	// Handle pre-release suffix
	patchStr := parts[2]
	if idx := strings.Index(patchStr, "-"); idx != -1 {
		patchStr = patchStr[:idx]
	}
	_, _ = fmt.Sscanf(patchStr, "%d", &patch)

	switch bumpType {
	case "major":
		major++
		minor = 0
		patch = 0
	case "minor":
		minor++
		patch = 0
	default: // patch
		patch++
	}

	return fmt.Sprintf("v%d.%d.%d", major, minor, patch)
}
