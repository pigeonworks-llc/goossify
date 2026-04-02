// Package generator provides project generation functionality for goossify.
package generator

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"time"
	"unicode"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// ProjectConfig represents the configuration for project generation.
type ProjectConfig struct {
	Name           string
	Description    string
	Type           string
	Author         string
	Email          string
	License        string
	GitHubUsername string
	ModulePath     string
	Year           int
	Version        string
	PackageName    string
	StructName     string
}

// Generator generates Go OSS projects from
type Generator struct {
	basePath string
	config   *ProjectConfig
}

// New creates a new project generator.
func New(basePath string, config *ProjectConfig) *Generator {
	if config.Year == 0 {
		config.Year = time.Now().Year()
	}
	if config.Version == "" {
		config.Version = "0.1.0"
	}

	if config.ModulePath == "" {
		config.ModulePath = buildModulePath(config.GitHubUsername, config.Name)
	}

	if config.PackageName == "" {
		config.PackageName = sanitizePackageName(config.Name)
	}

	if config.StructName == "" {
		config.StructName = toExportedName(config.Name)
	}

	return &Generator{basePath: basePath, config: config}
}

// Generate generates the complete project structure.
func (g *Generator) Generate() error {
	// 基本ディレクトリ構造を作成
	if err := g.createDirectoryStructure(); err != nil {
		return fmt.Errorf("ディレクトリ構造作成失敗: %w", err)
	}

	// 基本ファイルを生成
	if err := g.generateBaseFiles(); err != nil {
		return fmt.Errorf("基本ファイル生成失敗: %w", err)
	}

	// プロジェクトタイプ固有のファイルを生成
	if err := g.generateTypeSpecificFiles(); err != nil {
		return fmt.Errorf("タイプ固有ファイル生成失敗: %w", err)
	}

	// 設定ファイルを生成
	if err := g.generateConfigFiles(); err != nil {
		return fmt.Errorf("設定ファイル生成失敗: %w", err)
	}

	// GitHub関連ファイルを生成
	if err := g.generateGitHubFiles(); err != nil {
		return fmt.Errorf("GitHubファイル生成失敗: %w", err)
	}

	return nil
}

func (g *Generator) createDirectoryStructure() error {
	dirs := []string{
		"cmd",
		"internal",
		"pkg",
		"docs",
		"examples",
		"tests",
		".github/workflows",
		".github/ISSUE_TEMPLATE",
	}

	// プロジェクトタイプ固有のディレクトリ
	switch g.config.Type {
	case "cli-tool":
		dirs = append(dirs, "cmd/"+g.config.Name)
	case "web-api":
		dirs = append(dirs, "internal/handler", "internal/middleware", "internal/model")
	case "service":
		dirs = append(dirs, "internal/service", "internal/repository", "deployments")
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(g.absPath(dir), 0o750); err != nil {
			return fmt.Errorf("ディレクトリ %s 作成失敗: %w", dir, err)
		}
	}

	return nil
}

func (g *Generator) generateBaseFiles() error {
	baseFiles := map[string]string{
		"README.md":    readmeTemplate,
		"LICENSE":      g.getLicenseTemplate(),
		".gitignore":   gitignoreTemplate,
		"go.mod":       goModTemplate,
		"main.go":      g.getMainTemplate(),
		"Makefile":     makefileTemplate,
		"CHANGELOG.md": changelogTemplate,
	}

	for filename, templateContent := range baseFiles {
		if err := g.writeFileFromTemplate(filename, templateContent); err != nil {
			return fmt.Errorf("ファイル %s 生成失敗: %w", filename, err)
		}
	}

	return nil
}

func (g *Generator) generateTypeSpecificFiles() error {
	switch g.config.Type {
	case "cli-tool":
		return g.generateCLIFiles()
	case "library":
		return g.generateLibraryFiles()
	case "web-api":
		return g.generateWebAPIFiles()
	case "service":
		return g.generateServiceFiles()
	default:
		return fmt.Errorf("未知のプロジェクトタイプ: %s", g.config.Type)
	}
}

func (g *Generator) generateCLIFiles() error {
	cliFiles := map[string]string{
		filepath.Join("cmd", g.config.Name, "main.go"):     cliMainTemplate,
		filepath.Join("internal", "cmd", "root.go"):        cliRootTemplate,
		filepath.Join("internal", "cmd", "version.go"):     cliVersionTemplate,
		filepath.Join("internal", "version", "version.go"): versionTemplate,
	}

	for filename, templateContent := range cliFiles {
		if err := os.MkdirAll(filepath.Dir(g.absPath(filename)), 0o750); err != nil {
			return err
		}
		if err := g.writeFileFromTemplate(filename, templateContent); err != nil {
			return err
		}
	}

	return nil
}

func (g *Generator) generateLibraryFiles() error {
	libraryFiles := map[string]string{
		filepath.Join("pkg", g.config.PackageName, g.config.PackageName+".go"):      libraryMainTemplate,
		filepath.Join("pkg", g.config.PackageName, g.config.PackageName+"_test.go"): libraryTestTemplate,
		filepath.Join("examples", "basic", "main.go"):                               libraryExampleTemplate,
		"doc.go": libraryDocTemplate,
	}

	for filename, templateContent := range libraryFiles {
		if err := os.MkdirAll(filepath.Dir(g.absPath(filename)), 0o750); err != nil {
			return err
		}
		if err := g.writeFileFromTemplate(filename, templateContent); err != nil {
			return err
		}
	}

	return nil
}

func (g *Generator) generateWebAPIFiles() error {
	webAPIFiles := map[string]string{
		filepath.Join("internal", "handler", "handler.go"): webAPIHandlerTemplate,
		filepath.Join("internal", "middleware", "cors.go"): webAPICORSTemplate,
		filepath.Join("internal", "model", "response.go"):  webAPIModelTemplate,
		filepath.Join("internal", "server", "server.go"):   webAPIServerTemplate,
		"api/openapi.yaml": openAPITemplate,
	}

	for filename, templateContent := range webAPIFiles {
		if err := os.MkdirAll(filepath.Dir(g.absPath(filename)), 0o750); err != nil {
			return err
		}
		if err := g.writeFileFromTemplate(filename, templateContent); err != nil {
			return err
		}
	}

	return nil
}

func (g *Generator) generateServiceFiles() error {
	serviceFiles := map[string]string{
		filepath.Join("internal", "service", "service.go"):            serviceMainTemplate,
		filepath.Join("internal", "repository", "repository.go"):      serviceRepositoryTemplate,
		filepath.Join("internal", "config", "config.go"):              serviceConfigTemplate,
		filepath.Join("deployments", "docker", "Dockerfile"):          dockerfileTemplate,
		filepath.Join("deployments", "kubernetes", "deployment.yaml"): kubernetesTemplate,
	}

	for filename, templateContent := range serviceFiles {
		if err := os.MkdirAll(filepath.Dir(g.absPath(filename)), 0o750); err != nil {
			return err
		}
		if err := g.writeFileFromTemplate(filename, templateContent); err != nil {
			return err
		}
	}

	return nil
}

func (g *Generator) generateConfigFiles() error {
	configFiles := map[string]string{
		".golangci.yml":   golangCITemplate,
		".goreleaser.yml": goReleaserTemplate,
		".goossify.yml":   goossifyConfigTemplate,
	}

	for filename, templateContent := range configFiles {
		if err := g.writeFileFromTemplate(filename, templateContent); err != nil {
			return fmt.Errorf("設定ファイル %s 生成失敗: %w", filename, err)
		}
	}

	return nil
}

func (g *Generator) generateGitHubFiles() error {
	githubFiles := map[string]string{
		// Workflows
		".github/workflows/ci.yml":                 gitHubCITemplate,
		".github/workflows/release.yml":            gitHubReleaseTemplate,
		".github/workflows/auto-label.yml":         autoLabelerTemplate,
		".github/workflows/codeql.yml":             codeQLTemplate,
		".github/workflows/project-management.yml": projectManagementTemplate,

		// Issue and PR templates
		".github/ISSUE_TEMPLATE/bug_report.md":      bugReportTemplate,
		".github/ISSUE_TEMPLATE/feature_request.md": featureRequestTemplate,
		".github/ISSUE_TEMPLATE/question.md":        questionTemplate,
		".github/ISSUE_TEMPLATE/config.yml":         issueFormsConfigTemplate,
		".github/PULL_REQUEST_TEMPLATE.md":          pRTemplate,

		// Community files
		".github/CONTRIBUTING.md":    contributingTemplate,
		".github/CODE_OF_CONDUCT.md": codeOfConductTemplate,
		".github/SECURITY.md":        securityTemplate,
		".github/SUPPORT.md":         supportTemplate,
		".github/FUNDING.yml":        fundingTemplate,

		// Configuration files
		".github/dependabot.yml":      dependabotTemplate,
		".github/labeler.yml":         labelerConfigTemplate,
		".github/labels.yml":          gitHubLabelsTemplate,
		".github/auto-assign.yml":     autoAssignTemplate,
		".github/CODEOWNERS":          codeOwnersTemplate,
		".github/REPOSITORY_SETUP.md": gitHubSettingsTemplate,
	}

	for filename, templateContent := range githubFiles {
		if err := g.writeFileFromTemplate(filename, templateContent); err != nil {
			return fmt.Errorf("GitHubファイル %s 生成失敗: %w", filename, err)
		}
	}

	return nil
}

func (g *Generator) writeFileFromTemplate(filename, templateContent string) error {
	// テンプレート関数を定義
	titleCaser := cases.Title(language.Und, cases.NoLower)
	funcMap := template.FuncMap{
		"title": titleCaser.String,
		"upper": strings.ToUpper,
		"lower": strings.ToLower,
	}

	tmpl, err := template.New(filename).Funcs(funcMap).Parse(templateContent)
	if err != nil {
		return fmt.Errorf("テンプレート解析失敗: %w", err)
	}

	fullPath := g.absPath(filename)
	if err := os.MkdirAll(filepath.Dir(fullPath), 0o750); err != nil {
		return fmt.Errorf("ディレクトリ作成失敗: %w", err)
	}

	file, err := os.Create(fullPath)
	if err != nil {
		return fmt.Errorf("ファイル作成失敗: %w", err)
	}
	defer func() { _ = file.Close() }()

	if err := tmpl.Execute(file, g.config); err != nil {
		return fmt.Errorf("テンプレート実行失敗: %w", err)
	}

	return nil
}

func (g *Generator) getLicenseTemplate() string {
	switch g.config.License {
	case "MIT":
		return mITLicenseTemplate
	case "Apache-2.0":
		return apache2LicenseTemplate
	case "BSD-3-Clause":
		return bSD3LicenseTemplate
	default:
		return apache2LicenseTemplate
	}
}

func (g *Generator) getMainTemplate() string {
	switch g.config.Type {
	case "cli-tool":
		return cliMainEntryTemplate
	case "web-api":
		return webAPIMainTemplate
	case "service":
		return serviceMainEntryTemplate
	case "library":
		return "" // ライブラリにはmain.goは不要
	default:
		return defaultMainTemplate
	}
}

func (g *Generator) absPath(rel string) string {
	if filepath.IsAbs(rel) {
		return rel
	}
	return filepath.Join(g.basePath, rel)
}

func sanitizePackageName(name string) string {
	name = strings.TrimSpace(name)
	if name == "" {
		return "app"
	}

	var b strings.Builder
	lastUnderscore := false

	for _, r := range name {
		switch {
		case unicode.IsLetter(r):
			b.WriteRune(unicode.ToLower(r))
			lastUnderscore = false
		case unicode.IsDigit(r):
			if b.Len() == 0 {
				b.WriteRune('_')
			}
			b.WriteRune(r)
			lastUnderscore = false
		default:
			if !lastUnderscore && b.Len() > 0 {
				b.WriteRune('_')
				lastUnderscore = true
			}
		}
	}

	result := strings.Trim(b.String(), "_")
	if result == "" {
		return "app"
	}

	if !unicode.IsLetter(rune(result[0])) && result[0] != '_' {
		result = "_" + result
	}

	return result
}

func toExportedName(name string) string {
	name = strings.TrimSpace(name)
	if name == "" {
		return "App"
	}

	segments := strings.FieldsFunc(name, func(r rune) bool {
		return !unicode.IsLetter(r) && !unicode.IsDigit(r)
	})

	if len(segments) == 0 {
		return "App"
	}

	titleCaser := cases.Title(language.Und, cases.NoLower)
	var b strings.Builder

	for _, segment := range segments {
		s := titleCaser.String(segment)
		if s == "" {
			continue
		}
		b.WriteString(s)
	}

	if b.Len() == 0 {
		return "App"
	}

	return b.String()
}

func buildModulePath(githubUser, projectName string) string {
	user := strings.TrimSpace(githubUser)
	if user == "" {
		user = "your-username"
	}

	segment := strings.TrimSpace(projectName)
	if segment == "" {
		segment = "app"
	}
	segment = strings.ToLower(segment)
	segment = strings.Map(func(r rune) rune {
		if unicode.IsLetter(r) || unicode.IsDigit(r) || r == '-' || r == '/' {
			return r
		}
		return '-'
	}, segment)
	segment = strings.Trim(segment, "-")
	if segment == "" {
		segment = "app"
	}

	return fmt.Sprintf("github.com/%s/%s", user, segment)
}
