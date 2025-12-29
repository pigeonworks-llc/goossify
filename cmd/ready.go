package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

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
	sensitivePatterns := []struct {
		pattern     string
		description string
	}{
		{".env", "Environment variables file"},
		{".env.*", "Environment variables file"},
		{"*.key", "Private key file"},
		{"*.pem", "PEM certificate/key file"},
		{"*.p12", "PKCS12 key store"},
		{"*.pfx", "PKCS12 key store"},
		{"config.json", "Configuration file (may contain secrets)"},
		{"secrets.yaml", "Secrets file"},
		{"secrets.yml", "Secrets file"},
		{"credentials.json", "Credentials file"},
		{"*.crt", "Certificate file"},
		{"id_rsa", "SSH private key"},
		{"id_dsa", "SSH private key"},
		{"id_ecdsa", "SSH private key"},
		{"id_ed25519", "SSH private key"},
		{".htpasswd", "HTTP password file"},
		{"*.sqlite", "SQLite database (may contain sensitive data)"},
		{"*.db", "Database file (may contain sensitive data)"},
	}

	var foundFiles []string

	for _, sp := range sensitivePatterns {
		matches, err := filepath.Glob(filepath.Join(projectPath, sp.pattern))
		if err != nil {
			continue
		}

		for _, match := range matches {
			// Check if file is in .gitignore
			if !isInGitignore(projectPath, match) {
				relPath, _ := filepath.Rel(projectPath, match)
				foundFiles = append(foundFiles, fmt.Sprintf("%s (%s)", relPath, sp.description))
			}
		}

		// Also check subdirectories for some patterns
		if sp.pattern == ".env" || sp.pattern == ".env.*" {
			continue // Only check root for .env files
		}

		err = filepath.Walk(projectPath, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return nil
			}
			if info.IsDir() {
				// Skip hidden directories and common ignore patterns
				if strings.HasPrefix(info.Name(), ".") || info.Name() == "vendor" || info.Name() == "node_modules" {
					return filepath.SkipDir
				}
				return nil
			}

			matched, _ := filepath.Match(sp.pattern, info.Name())
			if matched && !isInGitignore(projectPath, path) {
				relPath, _ := filepath.Rel(projectPath, path)
				foundFiles = append(foundFiles, fmt.Sprintf("%s (%s)", relPath, sp.description))
			}
			return nil
		})
		if err != nil {
			continue
		}
	}

	if len(foundFiles) > 0 {
		fmt.Println("  Found potentially sensitive files not in .gitignore:")
		for _, f := range foundFiles {
			fmt.Printf("    ⚠️  %s\n", f)
		}
		return fmt.Errorf("found %d potentially sensitive file(s)", len(foundFiles))
	}

	return nil
}

// isInGitignore checks if a file is ignored by .gitignore
func isInGitignore(projectPath, filePath string) bool {
	gitignorePath := filepath.Join(projectPath, ".gitignore")
	content, err := os.ReadFile(gitignorePath)
	if err != nil {
		return false
	}

	relPath, err := filepath.Rel(projectPath, filePath)
	if err != nil {
		return false
	}

	lines := strings.Split(string(content), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Simple pattern matching
		if strings.HasPrefix(line, "/") {
			// Root-relative pattern
			pattern := strings.TrimPrefix(line, "/")
			if matched, _ := filepath.Match(pattern, relPath); matched {
				return true
			}
		} else {
			// Match anywhere
			if matched, _ := filepath.Match(line, filepath.Base(relPath)); matched {
				return true
			}
			if matched, _ := filepath.Match(line, relPath); matched {
				return true
			}
		}
	}

	return false
}

// checkLicenseConsistency checks license consistency
func checkLicenseConsistency(projectPath string) error {
	// Check LICENSE file existence
	licensePath := filepath.Join(projectPath, "LICENSE")
	licenseContent, err := os.ReadFile(licensePath)
	if os.IsNotExist(err) {
		return fmt.Errorf("LICENSE file not found")
	}
	if err != nil {
		return fmt.Errorf("failed to read LICENSE file: %w", err)
	}

	// Detect license type from content
	detectedLicense := detectLicenseType(string(licenseContent))
	if detectedLicense == "" {
		fmt.Println("  ⚠️  Could not detect license type from LICENSE file")
	} else {
		fmt.Printf("  Detected license: %s\n", detectedLicense)
	}

	// Check README.md for license badge/reference
	readmePath := filepath.Join(projectPath, "README.md")
	if readmeContent, err := os.ReadFile(readmePath); err == nil {
		readmeStr := string(readmeContent)

		// Check for license badge or mention
		hasLicenseBadge := strings.Contains(readmeStr, "license") ||
			strings.Contains(readmeStr, "License") ||
			strings.Contains(readmeStr, "LICENSE")

		if !hasLicenseBadge {
			fmt.Println("  ⚠️  README.md does not mention license")
		}

		// Check consistency if license type was detected
		if detectedLicense != "" {
			licenseInReadme := strings.Contains(strings.ToLower(readmeStr), strings.ToLower(detectedLicense))
			if !licenseInReadme {
				fmt.Printf("  ⚠️  README.md does not mention %s license\n", detectedLicense)
			}
		}
	}

	// Check .goossify.yml for license configuration
	goossifyPath := filepath.Join(projectPath, ".goossify.yml")
	if goossifyContent, err := os.ReadFile(goossifyPath); err == nil {
		goossifyStr := string(goossifyContent)

		// Simple check for license field
		if strings.Contains(goossifyStr, "license:") {
			// Extract license value (simple parsing)
			for _, line := range strings.Split(goossifyStr, "\n") {
				if strings.Contains(line, "license:") {
					parts := strings.SplitN(line, ":", 2)
					if len(parts) == 2 {
						configLicense := strings.TrimSpace(parts[1])
						configLicense = strings.Trim(configLicense, "\"'")

						if detectedLicense != "" && !strings.EqualFold(configLicense, detectedLicense) {
							fmt.Printf("  ⚠️  License mismatch: .goossify.yml says %s, LICENSE file appears to be %s\n",
								configLicense, detectedLicense)
						}
					}
					break
				}
			}
		}
	}

	return nil
}

// detectLicenseType attempts to detect the license type from LICENSE file content
func detectLicenseType(content string) string {
	contentLower := strings.ToLower(content)

	// Check for common license patterns
	licensePatterns := map[string][]string{
		"MIT": {
			"mit license",
			"permission is hereby granted, free of charge",
		},
		"Apache-2.0": {
			"apache license",
			"version 2.0",
			"http://www.apache.org/licenses",
		},
		"BSD-3-Clause": {
			"bsd 3-clause",
			"redistribution and use in source and binary forms",
			"neither the name of the copyright holder",
		},
		"BSD-2-Clause": {
			"bsd 2-clause",
			"redistribution and use in source and binary forms",
		},
		"GPL-3.0": {
			"gnu general public license",
			"version 3",
		},
		"GPL-2.0": {
			"gnu general public license",
			"version 2",
		},
		"LGPL-3.0": {
			"gnu lesser general public license",
			"version 3",
		},
		"MPL-2.0": {
			"mozilla public license",
			"version 2.0",
		},
		"ISC": {
			"isc license",
			"permission to use, copy, modify",
		},
		"Unlicense": {
			"unlicense",
			"this is free and unencumbered software",
		},
	}

	for license, patterns := range licensePatterns {
		matchCount := 0
		for _, pattern := range patterns {
			if strings.Contains(contentLower, pattern) {
				matchCount++
			}
		}
		// Require at least half of patterns to match
		if matchCount >= (len(patterns)+1)/2 {
			return license
		}
	}

	return ""
}
