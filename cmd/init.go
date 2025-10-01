package cmd

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"github.com/pigeonworks-llc/goossify/internal/generator"
)

const (
	projectTypeCLI     = "cli-tool"
	projectTypeLibrary = "library"
	projectTypeWebAPI  = "web-api"
	projectTypeService = "service"
)

var (
	interactive    bool
	projectType    string
	templateName   string
	author         string
	email          string
	license        string
	githubUsername string
)

// initCmd represents the init command.
var initCmd = &cobra.Command{
	Use:   "init [project-name]",
	Short: "新しいGo OSSプロジェクトを初期化",
	Long: `新しいGo言語のOSSプロジェクトを完全自動化で初期化します。

このコマンドは以下を自動生成します：
🏗️  最適化されたディレクトリ構造
📄  必須ファイル群 (README, LICENSE, .gitignore等)
🔧  開発ツール設定 (golangci-lint, GoReleaser等)
🤖  CI/CD パイプライン (GitHub Actions)
📊  品質管理ツール統合
👥  コミュニティファイル

利用可能なプロジェクトタイプ：
• cli-tool  - CLIアプリケーション (Cobra使用)
• library   - Go言語ライブラリ・パッケージ
• web-api   - REST API / GraphQL サーバー
• service   - マイクロサービス・デーモン

使用例:
  goossify init my-awesome-project
  goossify init --type cli-tool my-cli-app
  goossify init --interactive my-project`,
	Args: cobra.MaximumNArgs(1),
	RunE: runInit,
}

func init() {
	rootCmd.AddCommand(initCmd)

	initCmd.Flags().BoolVarP(&interactive, "interactive", "i", false, "対話的モードで設定")
	initCmd.Flags().StringVarP(&projectType, "type", "t", "", "プロジェクトタイプ (cli-tool|library|web-api|service)")
	initCmd.Flags().StringVar(&templateName, "template", "", "使用するテンプレート名")
	initCmd.Flags().StringVarP(&author, "author", "a", "", "作成者名")
	initCmd.Flags().StringVarP(&email, "email", "e", "", "作成者メールアドレス")
	initCmd.Flags().StringVarP(&license, "license", "l", "Apache-2.0", "ライセンス")
	initCmd.Flags().StringVarP(&githubUsername, "github", "g", "", "GitHubユーザー名")
}

func runInit(cmd *cobra.Command, args []string) error {
	var projectName string
	if len(args) > 0 {
		projectName = args[0]
	}

	// プロジェクト設定を収集
	config, err := collectProjectConfig(projectName)
	if err != nil {
		return fmt.Errorf("プロジェクト設定の収集に失敗: %w", err)
	}

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

	fmt.Printf("🎉 Go OSSプロジェクト '%s' が正常に作成されました！\n\n", config.Name)
	fmt.Println("次の手順:")
	fmt.Printf("  cd %s\n", config.Name)
	fmt.Println("  go mod tidy")
	fmt.Println("  git init")
	fmt.Println("  git add .")
	fmt.Println("  git commit -m \"🎉 Initial commit\"")
	fmt.Println()
	fmt.Println("プロジェクト管理:")
	fmt.Println("  goossify status     # プロジェクト健全性確認")

	return nil
}

func collectProjectConfig(projectName string) (*generator.ProjectConfig, error) {
	config := &generator.ProjectConfig{
		Name:           projectName,
		Type:           projectType,
		Author:         author,
		Email:          email,
		License:        license,
		GitHubUsername: githubUsername,
	}

	if interactive || projectName == "" || config.Type == "" {
		if err := collectConfigInteractively(config); err != nil {
			return nil, err
		}
	}

	// デフォルト値の設定
	if config.Name == "" {
		return nil, fmt.Errorf("プロジェクト名は必須です")
	}

	if config.Type == "" {
		config.Type = projectTypeCLI
	}

	if config.License == "" {
		config.License = "Apache-2.0"
	}

	if config.GitHubUsername == "" {
		config.GitHubUsername = "your-username"
	}

	// 説明の自動生成
	if config.Description == "" {
		config.Description = generateDescription(config.Type, config.Name)
	}

	return config, nil
}

func collectConfigInteractively(config *generator.ProjectConfig) error {
	reader := bufio.NewReader(os.Stdin)

	if err := promptProjectName(reader, config); err != nil {
		return err
	}
	if err := promptProjectType(reader, config); err != nil {
		return err
	}
	if err := promptDescription(reader, config); err != nil {
		return err
	}
	if err := promptAuthor(reader, config); err != nil {
		return err
	}
	if err := promptEmail(reader, config); err != nil {
		return err
	}
	if err := promptGitHubUsername(reader, config); err != nil {
		return err
	}
	if err := promptLicense(reader, config); err != nil {
		return err
	}

	return nil
}

func promptProjectName(reader *bufio.Reader, config *generator.ProjectConfig) error {
	if config.Name == "" {
		fmt.Print("プロジェクト名: ")
		name, err := reader.ReadString('\n')
		if err != nil {
			return err
		}
		config.Name = strings.TrimSpace(name)
	}
	return nil
}

func promptProjectType(reader *bufio.Reader, config *generator.ProjectConfig) error {
	if config.Type == "" {
		fmt.Println("\nプロジェクトタイプを選択してください:")
		fmt.Println("  1. cli-tool  - CLIアプリケーション")
		fmt.Println("  2. library   - Go言語ライブラリ")
		fmt.Println("  3. web-api   - REST API / GraphQL サーバー")
		fmt.Println("  4. service   - マイクロサービス・デーモン")
		fmt.Print("選択 [1-4] (1): ")

		choice, err := reader.ReadString('\n')
		if err != nil {
			return err
		}

		switch strings.TrimSpace(choice) {
		case "1", "":
			config.Type = projectTypeCLI
		case "2":
			config.Type = projectTypeLibrary
		case "3":
			config.Type = projectTypeWebAPI
		case "4":
			config.Type = projectTypeService
		default:
			config.Type = projectTypeCLI
		}
	}
	return nil
}

func promptDescription(reader *bufio.Reader, config *generator.ProjectConfig) error {
	fmt.Print("プロジェクトの説明: ")
	desc, err := reader.ReadString('\n')
	if err != nil {
		return err
	}
	config.Description = strings.TrimSpace(desc)
	return nil
}

func promptAuthor(reader *bufio.Reader, config *generator.ProjectConfig) error {
	if config.Author == "" {
		fmt.Print("作成者名: ")
		authorName, err := reader.ReadString('\n')
		if err != nil {
			return err
		}
		config.Author = strings.TrimSpace(authorName)
	}
	return nil
}

func promptEmail(reader *bufio.Reader, config *generator.ProjectConfig) error {
	if config.Email == "" {
		fmt.Print("メールアドレス: ")
		emailAddr, err := reader.ReadString('\n')
		if err != nil {
			return err
		}
		config.Email = strings.TrimSpace(emailAddr)
	}
	return nil
}

func promptGitHubUsername(reader *bufio.Reader, config *generator.ProjectConfig) error {
	if config.GitHubUsername == "" {
		fmt.Print("GitHubユーザー名: ")
		username, err := reader.ReadString('\n')
		if err != nil {
			return err
		}
		config.GitHubUsername = strings.TrimSpace(username)
	}
	return nil
}

func promptLicense(reader *bufio.Reader, config *generator.ProjectConfig) error {
	if config.License == "" {
		fmt.Print("ライセンス (Apache-2.0): ")
		licenseType, err := reader.ReadString('\n')
		if err != nil {
			return err
		}
		license := strings.TrimSpace(licenseType)
		if license == "" {
			license = "Apache-2.0"
		}
		config.License = license
	}
	return nil
}

func generateDescription(projectType, projectName string) string {
	switch projectType {
	case projectTypeCLI:
		return fmt.Sprintf("A powerful CLI tool built with Go - %s", projectName)
	case projectTypeLibrary:
		return fmt.Sprintf("A Go language library - %s", projectName)
	case projectTypeWebAPI:
		return fmt.Sprintf("A REST API service built with Go - %s", projectName)
	case projectTypeService:
		return fmt.Sprintf("A microservice built with Go - %s", projectName)
	default:
		return fmt.Sprintf("An awesome Go project - %s", projectName)
	}
}

func createProjectDirectory(projectName string) (string, error) {
	if projectName == "" {
		return "", fmt.Errorf("プロジェクト名は必須です")
	}

	wd, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("作業ディレクトリ取得失敗: %w", err)
	}

	projectPath := filepath.Join(wd, projectName)

	if info, err := os.Stat(projectPath); err == nil {
		if info.IsDir() {
			return "", fmt.Errorf("ディレクトリ '%s' は既に存在します", projectPath)
		}
		return "", fmt.Errorf("'%s' は既存のファイルです", projectPath)
	} else if !errors.Is(err, os.ErrNotExist) {
		return "", fmt.Errorf("ディレクトリ確認失敗: %w", err)
	}

	if err := os.MkdirAll(projectPath, 0o755); err != nil {
		return "", fmt.Errorf("ディレクトリ作成失敗: %w", err)
	}

	return projectPath, nil
}
