package analyzer

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// Secret patterns to detect in source code.
var secretPatterns = []*regexp.Regexp{
	regexp.MustCompile(`(?i)AKIA[0-9A-Z]{16}`),                              // AWS Access Key
	regexp.MustCompile(`(?i)["']sk[-_]live[-_][a-zA-Z0-9]{20,}["']`),        // Stripe Secret Key
	regexp.MustCompile(`(?i)["']ghp_[a-zA-Z0-9]{36}["']`),                   // GitHub PAT
	regexp.MustCompile(`(?i)["']gho_[a-zA-Z0-9]{36}["']`),                   // GitHub OAuth
	regexp.MustCompile(`(?i)password\s*[:=]\s*["'][^"']{8,}["']`),           // Hardcoded password
	regexp.MustCompile(`(?i)secret\s*[:=]\s*["'][^"']{8,}["']`),             // Hardcoded secret
	regexp.MustCompile(`(?i)api[_-]?key\s*[:=]\s*["'][a-zA-Z0-9]{16,}["']`), // API key assignment
	regexp.MustCompile(`-----BEGIN (RSA |EC )?PRIVATE KEY-----`),            // Private key
}

// Internal reference patterns to detect in source code.
var internalRefPatterns = []*regexp.Regexp{
	regexp.MustCompile(`\b10\.\d{1,3}\.\d{1,3}\.\d{1,3}\b`),                // 10.x.x.x
	regexp.MustCompile(`\b172\.(1[6-9]|2[0-9]|3[01])\.\d{1,3}\.\d{1,3}\b`), // 172.16-31.x.x
	regexp.MustCompile(`\b192\.168\.\d{1,3}\.\d{1,3}\b`),                   // 192.168.x.x
	regexp.MustCompile(`\b127\.0\.0\.1\b`),                                 // localhost IPv4
	regexp.MustCompile(`(?i)\.localhost\b`),                                // *.localhost
	regexp.MustCompile(`(?i)\.internal\b`),                                 // *.internal
	regexp.MustCompile(`(?i)\.local\b`),                                    // *.local
	regexp.MustCompile(`(?i)\.corp\b`),                                     // *.corp
	regexp.MustCompile(`(?i)forgejo\.localhost`),                           // Forgejo local
}

type finding struct {
	file    string
	line    int
	match   string
	pattern string
}

// analyzeCredentialScanning checks for hardcoded secrets in source code.
func (a *ProjectAnalyzer) analyzeCredentialScanning() CategoryResult {
	findings := a.scanFiles(secretPatterns, false)

	items := []Item{
		a.checkFile(".gitleaks.toml", "Gitleaks設定", false),
	}

	if len(findings) > 0 {
		desc := fmt.Sprintf("%d件の疑わしいパターンを検出", len(findings))
		for _, f := range findings {
			desc += fmt.Sprintf("\n  %s:%d %s", f.file, f.line, f.pattern)
		}
		items = append(items, Item{
			Name:        "Hardcoded secrets",
			Status:      "missing",
			Required:    true,
			Description: desc,
		})
	} else {
		items = append(items, Item{
			Name:        "Hardcoded secrets",
			Status:      "present",
			Required:    true,
			Description: "ハードコードされたシークレットなし",
		})
	}

	score := a.calculateCategoryScore(items)
	status := a.getStatusFromScore(score)

	return CategoryResult{
		Name:        "Credential Scanning",
		Score:       score,
		Status:      status,
		Description: "Hardcoded credential detection",
		Items:       items,
	}
}

// analyzeInternalReferences checks for private IPs and internal domains.
func (a *ProjectAnalyzer) analyzeInternalReferences() CategoryResult {
	findings := a.scanFiles(internalRefPatterns, true)

	items := []Item{}

	if len(findings) > 0 {
		desc := fmt.Sprintf("%d件の内部参照を検出", len(findings))
		for _, f := range findings {
			desc += fmt.Sprintf("\n  %s:%d %s", f.file, f.line, f.match)
		}
		items = append(items, Item{
			Name:        "Internal references",
			Status:      "missing",
			Required:    true,
			Description: desc,
		})
	} else {
		items = append(items, Item{
			Name:        "Internal references",
			Status:      "present",
			Required:    true,
			Description: "内部参照 (プライベートIP, 内部ドメイン) なし",
		})
	}

	score := a.calculateCategoryScore(items)
	status := a.getStatusFromScore(score)

	return CategoryResult{
		Name:        "Internal References",
		Score:       score,
		Status:      status,
		Description: "Private IP and internal domain detection",
		Items:       items,
	}
}

// scanFiles walks .go files and checks for pattern matches.
// If skipTests is true, *_test.go files are excluded.
func (a *ProjectAnalyzer) scanFiles(patterns []*regexp.Regexp, skipTests bool) []finding {
	var results []finding

	err := filepath.Walk(a.projectPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil // skip inaccessible files
		}
		if info.IsDir() {
			if info.Name() == "vendor" || info.Name() == "node_modules" || info.Name() == ".git" {
				return filepath.SkipDir
			}
			return nil
		}
		if !strings.HasSuffix(path, ".go") {
			return nil
		}
		if skipTests && strings.HasSuffix(path, "_test.go") {
			return nil
		}

		relPath, relErr := filepath.Rel(a.projectPath, path)
		if relErr != nil {
			relPath = path
		}

		fileFindings, scanErr := scanFile(path, relPath, patterns)
		if scanErr != nil {
			return nil // skip unreadable files
		}
		results = append(results, fileFindings...)

		return nil
	})
	if err != nil {
		return results
	}

	return results
}

func scanFile(path, relPath string, patterns []*regexp.Regexp) ([]finding, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open %s: %w", path, err)
	}
	defer f.Close()

	var results []finding
	scanner := bufio.NewScanner(f)
	lineNum := 0

	for scanner.Scan() {
		lineNum++
		line := scanner.Text()
		for _, p := range patterns {
			if p.MatchString(line) {
				results = append(results, finding{
					file:    relPath,
					line:    lineNum,
					match:   strings.TrimSpace(line),
					pattern: p.String(),
				})
				break // one match per line is enough
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return results, fmt.Errorf("scan %s: %w", path, err)
	}

	return results, nil
}
