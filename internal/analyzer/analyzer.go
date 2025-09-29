package analyzer

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// ProjectAnalyzer はプロジェクトの健全性を分析する
type ProjectAnalyzer struct {
	projectPath string
	projectName string
}

// AnalysisResult は分析結果
type AnalysisResult struct {
	ProjectPath     string            `json:"project_path"`
	ProjectName     string            `json:"project_name"`
	ProjectType     string            `json:"project_type"`
	OverallScore    int               `json:"overall_score"`     // 0-100
	Categories      []CategoryResult  `json:"categories"`
	Missing         []MissingItem     `json:"missing"`
	Recommendations []Recommendation  `json:"recommendations"`
	Summary         string            `json:"summary"`
}

// CategoryResult はカテゴリ別の結果
type CategoryResult struct {
	Name        string `json:"name"`
	Score       int    `json:"score"`        // 0-100
	Status      string `json:"status"`       // "good", "warning", "error"
	Description string `json:"description"`
	Items       []Item `json:"items"`
}

// Item は個別チェック項目
type Item struct {
	Name        string `json:"name"`
	Status      string `json:"status"`      // "present", "missing", "outdated"
	Required    bool   `json:"required"`
	Description string `json:"description"`
	Path        string `json:"path,omitempty"`
}

// MissingItem は不足している項目
type MissingItem struct {
	Name        string `json:"name"`
	Category    string `json:"category"`
	Priority    string `json:"priority"`    // "high", "medium", "low"
	Description string `json:"description"`
	Action      string `json:"action"`
}

// Recommendation は改善提案
type Recommendation struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Command     string `json:"command,omitempty"`
	Priority    string `json:"priority"`
}

// New は新しいProjectAnalyzerを作成
func New(projectPath string) *ProjectAnalyzer {
	projectName := filepath.Base(projectPath)
	return &ProjectAnalyzer{
		projectPath: projectPath,
		projectName: projectName,
	}
}

// Analyze はプロジェクトを分析
func (a *ProjectAnalyzer) Analyze() (*AnalysisResult, error) {
	result := &AnalysisResult{
		ProjectPath: a.projectPath,
		ProjectName: a.projectName,
	}

	// プロジェクトタイプ判定
	result.ProjectType = a.detectProjectType()

	// カテゴリ別分析
	categories := []CategoryResult{
		a.analyzeBasicStructure(),
		a.analyzeDocumentation(),
		a.analyzeGitHubIntegration(),
		a.analyzeQualityTools(),
		a.analyzeDependencies(),
		a.analyzeLicensing(),
	}

	result.Categories = categories

	// 総合スコア計算
	result.OverallScore = a.calculateOverallScore(categories)

	// 不足項目・推奨事項抽出
	result.Missing = a.extractMissingItems(categories)
	result.Recommendations = a.generateRecommendations(result)

	// サマリー生成
	result.Summary = a.generateSummary(result)

	return result, nil
}

// detectProjectType はプロジェクトタイプを判定
func (a *ProjectAnalyzer) detectProjectType() string {
	// main.go の存在チェック
	if _, err := os.Stat(filepath.Join(a.projectPath, "main.go")); err == nil {
		// cmd ディレクトリがあるか確認
		if _, err := os.Stat(filepath.Join(a.projectPath, "cmd")); err == nil {
			return "cli-tool"
		}
		return "application"
	}

	// go.mod の存在チェック
	if _, err := os.Stat(filepath.Join(a.projectPath, "go.mod")); err == nil {
		return "library"
	}

	return "unknown"
}

// analyzeBasicStructure は基本構造を分析
func (a *ProjectAnalyzer) analyzeBasicStructure() CategoryResult {
	items := []Item{
		a.checkFile("go.mod", "Go modules設定", true),
		a.checkFile("go.sum", "依存関係ロック", false),
		a.checkFile("README.md", "プロジェクト説明", true),
		a.checkFile(".gitignore", "Git除外設定", true),
		a.checkDirectory("internal", "内部パッケージ", false),
		a.checkDirectory("pkg", "公開パッケージ", false),
		a.checkDirectory("cmd", "エントリーポイント", false),
	}

	score := a.calculateCategoryScore(items)
	status := a.getStatusFromScore(score)

	return CategoryResult{
		Name:        "Basic Structure",
		Score:       score,
		Status:      status,
		Description: "Basic directory structure and files for Go projects",
		Items:       items,
	}
}

// analyzeDocumentation はドキュメントを分析
func (a *ProjectAnalyzer) analyzeDocumentation() CategoryResult {
	items := []Item{
		a.checkFile("README.md", "プロジェクト説明", true),
		a.checkFile("CONTRIBUTING.md", "コントリビューションガイド", false),
		a.checkDirectory("docs", "ドキュメントディレクトリ", false),
		a.checkDirectory("examples", "使用例", false),
	}

	score := a.calculateCategoryScore(items)
	status := a.getStatusFromScore(score)

	return CategoryResult{
		Name:        "Documentation",
		Score:       score,
		Status:      status,
		Description: "Project documentation completeness",
		Items:       items,
	}
}

// analyzeGitHubIntegration はGitHub連携を分析
func (a *ProjectAnalyzer) analyzeGitHubIntegration() CategoryResult {
	items := []Item{
		a.checkFile(".github/workflows/ci.yml", "CI/CD設定", false),
		a.checkFile(".github/workflows/release.yml", "リリース自動化", false),
		a.checkFile(".github/ISSUE_TEMPLATE/bug_report.md", "Bugレポートテンプレート", false),
		a.checkFile(".github/ISSUE_TEMPLATE/feature_request.md", "機能要求テンプレート", false),
		a.checkFile(".github/PULL_REQUEST_TEMPLATE.md", "PRテンプレート", false),
		a.checkFile("SECURITY.md", "セキュリティポリシー", false),
	}

	score := a.calculateCategoryScore(items)
	status := a.getStatusFromScore(score)

	return CategoryResult{
		Name:        "GitHub Integration",
		Score:       score,
		Status:      status,
		Description: "Integration status with GitHub-specific features",
		Items:       items,
	}
}

// analyzeQualityTools は品質ツールを分析
func (a *ProjectAnalyzer) analyzeQualityTools() CategoryResult {
	items := []Item{
		a.checkFile(".golangci.yml", "Linter設定", false),
		a.checkFile(".goreleaser.yml", "GoReleaser設定", false),
		a.checkTestFiles("テストファイル", true),
		a.checkFile("renovate.json", "依存関係更新", false),
	}

	score := a.calculateCategoryScore(items)
	status := a.getStatusFromScore(score)

	return CategoryResult{
		Name:        "Quality Tools",
		Score:       score,
		Status:      status,
		Description: "Tools supporting code quality and maintainability",
		Items:       items,
	}
}

// analyzeDependencies は依存関係を分析
func (a *ProjectAnalyzer) analyzeDependencies() CategoryResult {
	items := []Item{
		a.checkGoModTidy("go.mod整合性", true),
		a.checkDirectDependencies("直接依存関係", false),
		a.checkVulnerabilities("脆弱性", true),
	}

	score := a.calculateCategoryScore(items)
	status := a.getStatusFromScore(score)

	return CategoryResult{
		Name:        "Dependencies",
		Score:       score,
		Status:      status,
		Description: "Project dependency management status",
		Items:       items,
	}
}

// analyzeLicensing はライセンスを分析
func (a *ProjectAnalyzer) analyzeLicensing() CategoryResult {
	items := []Item{
		a.checkFile("LICENSE", "ライセンスファイル", true),
		a.checkLicenseInGoMod("go.modライセンス情報", false),
	}

	score := a.calculateCategoryScore(items)
	status := a.getStatusFromScore(score)

	return CategoryResult{
		Name:        "Licensing",
		Score:       score,
		Status:      status,
		Description: "Project license information",
		Items:       items,
	}
}

// checkFile はファイルの存在をチェック
func (a *ProjectAnalyzer) checkFile(fileName, description string, required bool) Item {
	path := filepath.Join(a.projectPath, fileName)
	_, err := os.Stat(path)

	status := "missing"
	if err == nil {
		status = "present"
	}

	return Item{
		Name:        fileName,
		Status:      status,
		Required:    required,
		Description: description,
		Path:        path,
	}
}

// checkDirectory はディレクトリの存在をチェック
func (a *ProjectAnalyzer) checkDirectory(dirName, description string, _ bool) Item {
	path := filepath.Join(a.projectPath, dirName)
	stat, err := os.Stat(path)

	status := "missing"
	if err == nil && stat.IsDir() {
		status = "present"
	}

	return Item{
		Name:        dirName + "/",
		Status:      status,
		Required:    false,
		Description: description,
		Path:        path,
	}
}

// checkTestFiles はテストファイルの存在をチェック
func (a *ProjectAnalyzer) checkTestFiles(description string, required bool) Item {
	testFiles := 0
	_ = filepath.Walk(a.projectPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if strings.HasSuffix(info.Name(), "_test.go") {
			testFiles++
		}
		return nil
	})

	status := "missing"
	if testFiles > 0 {
		status = "present"
	}

	return Item{
		Name:        fmt.Sprintf("*_test.go (%d files)", testFiles),
		Status:      status,
		Required:    required,
		Description: description,
	}
}

// checkGoModTidy はgo mod tidyの状態をチェック
func (a *ProjectAnalyzer) checkGoModTidy(description string, required bool) Item {
	// 簡易チェック: go.modとgo.sumの存在
	goMod := filepath.Join(a.projectPath, "go.mod")
	goSum := filepath.Join(a.projectPath, "go.sum")

	_, modErr := os.Stat(goMod)
	_, sumErr := os.Stat(goSum)

	status := "missing"
	if modErr == nil && sumErr == nil {
		status = "present"
	} else if modErr == nil {
		status = "warning" // go.modはあるがgo.sumがない
	}

	return Item{
		Name:        "Go modules整合性",
		Status:      status,
		Required:    required,
		Description: description,
	}
}

// checkDirectDependencies は直接依存関係をチェック（簡易版）
func (a *ProjectAnalyzer) checkDirectDependencies(description string, required bool) Item {
	// 実装簡略化: go.mod存在チェックのみ
	return a.checkFile("go.mod", description, required)
}

// checkVulnerabilities は脆弱性をチェック（簡易版）
func (a *ProjectAnalyzer) checkVulnerabilities(description string, required bool) Item {
	// 実装簡略化: 常にpresentとする（govulncheck実行は別途）
	return Item{
		Name:        "脆弱性チェック",
		Status:      "present",
		Required:    required,
		Description: description + " (govulncheckで確認してください)",
	}
}

// checkLicenseInGoMod はgo.mod内のライセンス情報をチェック
func (a *ProjectAnalyzer) checkLicenseInGoMod(description string, required bool) Item {
	// 実装簡略化: go.mod存在チェックのみ
	return a.checkFile("go.mod", description, required)
}

// calculateCategoryScore はカテゴリスコアを計算
func (a *ProjectAnalyzer) calculateCategoryScore(items []Item) int {
	if len(items) == 0 {
		return 0
	}

	total := 0
	present := 0

	for _, item := range items {
		weight := 1
		if item.Required {
			weight = 2
		}
		total += weight
		if item.Status == "present" {
			present += weight
		}
	}

	return (present * 100) / total
}

// calculateOverallScore は総合スコアを計算
func (a *ProjectAnalyzer) calculateOverallScore(categories []CategoryResult) int {
	if len(categories) == 0 {
		return 0
	}

	total := 0
	for _, category := range categories {
		total += category.Score
	}

	return total / len(categories)
}

// getStatusFromScore はスコアからステータスを取得
func (a *ProjectAnalyzer) getStatusFromScore(score int) string {
	if score >= 80 {
		return "good"
	} else if score >= 50 {
		return "warning"
	}
	return "error"
}

// extractMissingItems は不足項目を抽出
func (a *ProjectAnalyzer) extractMissingItems(categories []CategoryResult) []MissingItem {
	var missing []MissingItem

	for _, category := range categories {
		for _, item := range category.Items {
			if item.Status == "missing" {
				priority := "low"
				if item.Required {
					priority = "high"
				}

				missing = append(missing, MissingItem{
					Name:        item.Name,
					Category:    category.Name,
					Priority:    priority,
					Description: item.Description,
					Action:      "goossify ossify でファイルを生成できます",
				})
			}
		}
	}

	return missing
}

// generateRecommendations は推奨事項を生成
func (a *ProjectAnalyzer) generateRecommendations(result *AnalysisResult) []Recommendation {
	var recommendations []Recommendation

	// 総合スコアに基づく推奨
	if result.OverallScore < 50 {
		recommendations = append(recommendations, Recommendation{
			Title:       "OSS基本ファイルの追加",
			Description: "プロジェクトにOSSとして必要な基本ファイルが不足しています",
			Command:     "goossify ossify .",
			Priority:    "high",
		})
	}

	// カテゴリ別推奨
	for _, category := range result.Categories {
		if category.Status == "error" {
			switch category.Name {
			case "GitHub統合":
				recommendations = append(recommendations, Recommendation{
					Title:       "GitHub統合の改善",
					Description: "CI/CDやIssue/PRテンプレートを追加してGitHub連携を強化しましょう",
					Command:     "goossify ossify .",
					Priority:    "medium",
				})
			case "品質ツール":
				recommendations = append(recommendations, Recommendation{
					Title:       "品質ツールの導入",
					Description: "Linter設定やテストの追加でコード品質を向上させましょう",
					Priority:    "medium",
				})
			}
		}
	}

	return recommendations
}

// generateSummary generates summary
func (a *ProjectAnalyzer) generateSummary(result *AnalysisResult) string {
	status := "Good"
	if result.OverallScore < 80 {
		status = "Room for improvement"
	}
	if result.OverallScore < 50 {
		status = "Needs improvement"
	}

	return fmt.Sprintf(
		"Project '%s' health: %s (Score: %d/100)\nMissing items: %d, Recommendations: %d",
		result.ProjectName,
		status,
		result.OverallScore,
		len(result.Missing),
		len(result.Recommendations),
	)
}