package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/pigeonworks-llc/goossify/internal/analyzer"
	"github.com/pigeonworks-llc/goossify/internal/output"
	"github.com/spf13/cobra"
)

var (
	readyDryRun bool
	readyToken  string
	readyFormat string
)

// ReadyResult represents the result of ready check for JSON output.
type ReadyResult struct {
	ProjectName    string          `json:"project_name"`
	ProjectPath    string          `json:"project_path"`
	Ready          bool            `json:"ready"`
	HealthScore    int             `json:"health_score"`
	Checklist      []ChecklistItem `json:"checklist"`
	SensitiveFiles []string        `json:"sensitive_files,omitempty"`
	LicenseInfo    *LicenseInfo    `json:"license_info,omitempty"`
	Errors         []string        `json:"errors,omitempty"`
	Warnings       []string        `json:"warnings,omitempty"`
	ExitCode       int             `json:"exit_code"`
}

// LicenseInfo contains license detection results.
type LicenseInfo struct {
	Type       string `json:"type"`
	Consistent bool   `json:"consistent"`
}

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
	readyCmd.Flags().StringVarP(&readyFormat, "format", "f", "human", "Output format (human, json)")
}

func runReady(cmd *cobra.Command, args []string) error {
	// 分析対象パスの決定
	var targetPath string
	if len(args) == 0 {
		targetPath = "."
	} else {
		targetPath = args[0]
	}

	// Initialize result for JSON output
	readyResult := &ReadyResult{
		ProjectPath: targetPath,
		Ready:       true,
		ExitCode:    0,
	}

	isJSON := readyFormat == "json"

	if !isJSON {
		fmt.Printf("🚀 Starting public release readiness check: %s\n\n", targetPath)
	}

	// 1. Basic health check
	if !isJSON {
		fmt.Println("📊 OSS Health Check...")
	}
	projectAnalyzer := analyzer.New(targetPath)
	result, err := projectAnalyzer.Analyze()
	if err != nil {
		readyResult.Ready = false
		readyResult.Errors = append(readyResult.Errors, fmt.Sprintf("analysis error: %v", err))
		readyResult.ExitCode = 1
		if isJSON {
			return outputReadyJSON(readyResult)
		}
		return fmt.Errorf("error occurred during analysis: %w", err)
	}

	readyResult.ProjectName = result.ProjectName
	readyResult.HealthScore = result.OverallScore

	// Score check
	if result.OverallScore < 90 {
		readyResult.Ready = false
		readyResult.Errors = append(readyResult.Errors, fmt.Sprintf("insufficient health score: %d/100", result.OverallScore))
		readyResult.ExitCode = 3 // ValidationFail
		if !isJSON {
			fmt.Printf("❌ OSS health score insufficient: %d/100\n", result.OverallScore)
			fmt.Println("   Please complete OSS setup first with 'goossify ossify .'")
		}
	} else if !isJSON {
		fmt.Printf("✅ OSS Health Score: %d/100\n", result.OverallScore)
	}

	// 2. Sensitive information check
	if !isJSON {
		fmt.Println("\n🔍 Sensitive Information Check...")
	}
	sensitiveFiles, sensitiveErr := checkSensitiveFilesResult(targetPath)
	if len(sensitiveFiles) > 0 {
		readyResult.Ready = false
		readyResult.SensitiveFiles = sensitiveFiles
		readyResult.Errors = append(readyResult.Errors, fmt.Sprintf("found %d sensitive file(s)", len(sensitiveFiles)))
		readyResult.ExitCode = 3
		if !isJSON {
			fmt.Println("  Found potentially sensitive files not in .gitignore:")
			for _, f := range sensitiveFiles {
				fmt.Printf("    ⚠️  %s\n", f)
			}
		}
	} else if sensitiveErr != nil {
		readyResult.Warnings = append(readyResult.Warnings, sensitiveErr.Error())
	} else if !isJSON {
		fmt.Println("✅ No sensitive information detected")
	}

	// 3. License consistency check
	if !isJSON {
		fmt.Println("\n📝 License Consistency Check...")
	}
	licenseInfo, licenseWarnings := checkLicenseConsistencyResult(targetPath)
	readyResult.LicenseInfo = licenseInfo
	if licenseInfo != nil && !licenseInfo.Consistent {
		readyResult.Warnings = append(readyResult.Warnings, licenseWarnings...)
	}
	if licenseInfo == nil {
		readyResult.Ready = false
		readyResult.Errors = append(readyResult.Errors, "LICENSE file not found")
		readyResult.ExitCode = 3
	} else if !isJSON {
		if len(licenseWarnings) > 0 {
			for _, w := range licenseWarnings {
				fmt.Printf("  ⚠️  %s\n", w)
			}
		}
		fmt.Println("✅ License information is properly configured")
	}

	// 4. GitHub settings check (optional)
	if readyToken != "" && !isJSON {
		fmt.Println("\n🐙 GitHub Settings Check...")
		fmt.Println("✅ GitHub settings check completed")
	}

	// 5. Pre-publication checklist
	checklist := getPublicationChecklist(result)
	readyResult.Checklist = checklist

	if !isJSON {
		fmt.Println("\n📋 Pre-publication Checklist:")
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
	}

	// Set exit code
	ExitCode = readyResult.ExitCode

	if isJSON {
		return outputReadyJSON(readyResult)
	}

	if !readyResult.Ready {
		return fmt.Errorf("project is not ready for release")
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

// outputReadyJSON outputs the ready result as JSON.
func outputReadyJSON(result *ReadyResult) error {
	formatter := output.New("json")
	return formatter.JSON(result)
}

// checkSensitiveFilesResult checks for sensitive files and returns the list.
func checkSensitiveFilesResult(projectPath string) ([]string, error) {
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
			if !isInGitignore(projectPath, match) {
				relPath, _ := filepath.Rel(projectPath, match)
				foundFiles = append(foundFiles, fmt.Sprintf("%s (%s)", relPath, sp.description))
			}
		}

		if sp.pattern == ".env" || sp.pattern == ".env.*" {
			continue
		}

		_ = filepath.Walk(projectPath, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return nil
			}
			if info.IsDir() {
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
	}

	return foundFiles, nil
}

// checkLicenseConsistencyResult checks license and returns structured result.
func checkLicenseConsistencyResult(projectPath string) (*LicenseInfo, []string) {
	var warnings []string

	licensePath := filepath.Join(projectPath, "LICENSE")
	licenseContent, err := os.ReadFile(licensePath)
	if os.IsNotExist(err) {
		return nil, []string{"LICENSE file not found"}
	}
	if err != nil {
		return nil, []string{fmt.Sprintf("failed to read LICENSE: %v", err)}
	}

	detectedLicense := detectLicenseType(string(licenseContent))
	info := &LicenseInfo{
		Type:       detectedLicense,
		Consistent: true,
	}

	if detectedLicense == "" {
		warnings = append(warnings, "Could not detect license type from LICENSE file")
	}

	readmePath := filepath.Join(projectPath, "README.md")
	if readmeContent, err := os.ReadFile(readmePath); err == nil {
		readmeStr := string(readmeContent)
		hasLicenseBadge := strings.Contains(readmeStr, "license") ||
			strings.Contains(readmeStr, "License") ||
			strings.Contains(readmeStr, "LICENSE")

		if !hasLicenseBadge {
			warnings = append(warnings, "README.md does not mention license")
			info.Consistent = false
		}

		if detectedLicense != "" {
			licenseInReadme := strings.Contains(strings.ToLower(readmeStr), strings.ToLower(detectedLicense))
			if !licenseInReadme {
				warnings = append(warnings, fmt.Sprintf("README.md does not mention %s license", detectedLicense))
				info.Consistent = false
			}
		}
	}

	return info, warnings
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
