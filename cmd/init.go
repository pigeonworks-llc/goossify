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
	Short: "Initialize a new Go OSS project",
	Long: `Initialize a new Go OSS project with full automation.

This command automatically generates:
🏗️  Optimized directory structure
📄  Essential files (README, LICENSE, .gitignore, etc.)
🔧  Development tool configurations (golangci-lint, GoReleaser, etc.)
🤖  CI/CD pipeline (GitHub Actions)
📊  Quality management tool integration
👥  Community files

Available project types:
• cli-tool  - CLI application (using Cobra)
• library   - Go library/package
• web-api   - REST API / GraphQL server
• service   - Microservice/daemon

Examples:
  goossify init my-awesome-project
  goossify init --type cli-tool my-cli-app
  goossify init --interactive my-project`,
	Args: cobra.MaximumNArgs(1),
	RunE: runInit,
}

func init() {
	rootCmd.AddCommand(initCmd)

	initCmd.Flags().BoolVarP(&interactive, "interactive", "i", false, "Interactive mode")
	initCmd.Flags().StringVarP(&projectType, "type", "t", "", "Project type (cli-tool|library|web-api|service)")
	initCmd.Flags().StringVar(&templateName, "template", "", "Template name to use")
	initCmd.Flags().StringVarP(&author, "author", "a", "", "Author name")
	initCmd.Flags().StringVarP(&email, "email", "e", "", "Author email address")
	initCmd.Flags().StringVarP(&license, "license", "l", "MIT", "License type")
	initCmd.Flags().StringVarP(&githubUsername, "github", "g", "", "GitHub username")
}

func runInit(cmd *cobra.Command, args []string) error {
	var projectName string
	if len(args) > 0 {
		projectName = args[0]
	}

	// Collect project configuration
	config, err := collectProjectConfig(projectName)
	if err != nil {
		return fmt.Errorf("failed to collect project configuration: %w", err)
	}

	// Create project directory
	projectPath, err := createProjectDirectory(config.Name)
	if err != nil {
		return fmt.Errorf("failed to create project directory: %w", err)
	}

	// Generate project
	gen := generator.New(projectPath, config)
	if err := gen.Generate(); err != nil {
		return fmt.Errorf("failed to generate project: %w", err)
	}

	fmt.Printf("🎉 Successfully created Go OSS project '%s'!\n\n", config.Name)
	fmt.Println("Next steps:")
	fmt.Printf("  cd %s\n", config.Name)
	fmt.Println("  go mod tidy")
	fmt.Println("  git init")
	fmt.Println("  git add .")
	fmt.Println("  git commit -m \"🎉 Initial commit\"")
	fmt.Println()
	fmt.Println("Project management:")
	fmt.Println("  goossify status     # Check project health")

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

	// Set default values
	if config.Name == "" {
		return nil, fmt.Errorf("project name is required")
	}

	if config.Type == "" {
		config.Type = projectTypeCLI
	}

	if config.License == "" {
		config.License = "MIT"
	}

	if config.GitHubUsername == "" {
		config.GitHubUsername = "your-username"
	}

	// Auto-generate description
	if config.Description == "" {
		config.Description = generateDescription(config.Type, config.Name)
	}

	return config, nil
}

func collectConfigInteractively(config *generator.ProjectConfig) error {
	reader := bufio.NewReader(os.Stdin)

	// Project name
	if config.Name == "" {
		fmt.Print("Project name: ")
		name, err := reader.ReadString('\n')
		if err != nil {
			return err
		}
		config.Name = strings.TrimSpace(name)
	}

	// Project type
	if config.Type == "" {
		fmt.Println("\nSelect project type:")
		fmt.Println("  1. cli-tool  - CLI application")
		fmt.Println("  2. library   - Go library")
		fmt.Println("  3. web-api   - REST API / GraphQL server")
		fmt.Println("  4. service   - Microservice/daemon")
		fmt.Print("Choice [1-4] (1): ")

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

	// Description
	fmt.Print("Project description: ")
	desc, err := reader.ReadString('\n')
	if err != nil {
		return err
	}
	config.Description = strings.TrimSpace(desc)

	// Author name
	if config.Author == "" {
		fmt.Print("Author name: ")
		authorName, err := reader.ReadString('\n')
		if err != nil {
			return err
		}
		config.Author = strings.TrimSpace(authorName)
	}

	// Email address
	if config.Email == "" {
		fmt.Print("Email address: ")
		emailAddr, err := reader.ReadString('\n')
		if err != nil {
			return err
		}
		config.Email = strings.TrimSpace(emailAddr)
	}

	// GitHub username
	if config.GitHubUsername == "" {
		fmt.Print("GitHub username: ")
		username, err := reader.ReadString('\n')
		if err != nil {
			return err
		}
		config.GitHubUsername = strings.TrimSpace(username)
	}

	// License
	if config.License == "" {
		fmt.Print("License (MIT): ")
		licenseType, err := reader.ReadString('\n')
		if err != nil {
			return err
		}
		license := strings.TrimSpace(licenseType)
		if license == "" {
			license = "MIT"
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
		return "", fmt.Errorf("project name is required")
	}

	wd, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("failed to get working directory: %w", err)
	}

	projectPath := filepath.Join(wd, projectName)

	if info, err := os.Stat(projectPath); err == nil {
		if info.IsDir() {
			return "", fmt.Errorf("directory '%s' already exists", projectPath)
		}
		return "", fmt.Errorf("'%s' is an existing file", projectPath)
	} else if !errors.Is(err, os.ErrNotExist) {
		return "", fmt.Errorf("failed to check directory: %w", err)
	}

	if err := os.MkdirAll(projectPath, 0o755); err != nil {
		return "", fmt.Errorf("failed to create directory: %w", err)
	}

	return projectPath, nil
}
