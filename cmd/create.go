package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/pigeonworks-llc/goossify/internal/generator"
)

var (
	createTemplate    string
	createAuthor      string
	createEmail       string
	createLicense     string
	createGithub      string
	createInteractive bool
)

// createCmd represents the create command.
var createCmd = &cobra.Command{
	Use:   "create [project-name]",
	Short: "テンプレートから新しいGo OSSプロジェクトを作成",
	Long: `事前定義されたテンプレートから新しいGo言語のOSSプロジェクトを作成します。

利用可能なテンプレート：
🔧 cli-tool  - CLIアプリケーション (Cobra使用)
📚 library   - Go言語ライブラリ・パッケージ
🌐 web-api   - REST API / GraphQL サーバー
⚙️  service   - マイクロサービス・デーモン

使用例:
  goossify create --template cli-tool my-cli-app
  goossify create --template library my-go-lib
  goossify create --template web-api my-api-server
  goossify create --template service my-service`,
	Args: cobra.ExactArgs(1),
	RunE: runCreate,
}

func init() {
	rootCmd.AddCommand(createCmd)

	createCmd.Flags().StringVarP(&createTemplate, "template", "t", "", "使用するテンプレート (cli-tool|library|web-api|service)")
	createCmd.Flags().StringVarP(&createAuthor, "author", "a", "", "作成者名")
	createCmd.Flags().StringVarP(&createEmail, "email", "e", "", "作成者メールアドレス")
	createCmd.Flags().StringVarP(&createLicense, "license", "l", "Apache-2.0", "ライセンス")
	createCmd.Flags().StringVarP(&createGithub, "github", "g", "", "GitHubユーザー名")
	createCmd.Flags().BoolVarP(&createInteractive, "interactive", "i", false, "対話的モードで設定")

	_ = createCmd.MarkFlagRequired("template")
}

func runCreate(cmd *cobra.Command, args []string) error {
	projectName := args[0]

	// テンプレートの有効性チェック
	if !isValidTemplate(createTemplate) {
		return fmt.Errorf("無効なテンプレート: %s\n利用可能: cli-tool, library, web-api, service", createTemplate)
	}

	// プロジェクト設定を収集
	config := &generator.ProjectConfig{
		Name:           projectName,
		Type:           createTemplate,
		Author:         createAuthor,
		Email:          createEmail,
		License:        createLicense,
		GitHubUsername: createGithub,
	}

	if createInteractive || needsInteractiveInput(config) {
		if err := collectConfigInteractively(config); err != nil {
			return fmt.Errorf("設定収集に失敗: %w", err)
		}
	}

	// デフォルト値設定
	setDefaultValues(config)

	// プロジェクトディレクトリ作成
	projectPath, err := createProjectDirectory(config.Name)
	if err != nil {
		return fmt.Errorf("プロジェクトディレクトリ作成に失敗: %w", err)
	}

	// プロジェクト生成
	gen := generator.New(projectPath, config)
	if err := gen.Generate(); err != nil {
		return fmt.Errorf("プロジェクト生成に失敗: %w", err)
	}

	fmt.Printf("🎉 Go OSSプロジェクト '%s' (%s) が正常に作成されました！\n\n", config.Name, config.Type)
	printNextSteps(config.Name)

	return nil
}

func isValidTemplate(template string) bool {
	validTemplates := []string{"cli-tool", "library", "web-api", "service"}
	for _, valid := range validTemplates {
		if template == valid {
			return true
		}
	}
	return false
}

func needsInteractiveInput(config *generator.ProjectConfig) bool {
	return config.Author == "" || config.Email == "" || config.GitHubUsername == ""
}

func setDefaultValues(config *generator.ProjectConfig) {
	if config.License == "" {
		config.License = "Apache-2.0"
	}
	if config.GitHubUsername == "" {
		config.GitHubUsername = "your-username"
	}
	if config.Description == "" {
		config.Description = generateDescription(config.Type, config.Name)
	}
}

func printNextSteps(projectName string) {
	fmt.Println("次の手順:")
	fmt.Printf("  cd %s\n", projectName)
	fmt.Println("  go mod tidy")
	fmt.Println("  git init")
	fmt.Println("  git add .")
	fmt.Println("  git commit -m \"🎉 Initial commit\"")
	fmt.Println()
	fmt.Println("GitHubリポジトリを作成してプッシュ:")
	fmt.Println("  gh repo create --public")
	fmt.Println("  git push -u origin main")
	fmt.Println()
	fmt.Println("プロジェクト管理:")
	fmt.Println("  goossify status     # プロジェクト健全性確認")
}
