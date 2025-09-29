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
	Short: "Create a new Go OSS project from a template",
	Long: `Create a new Go OSS project from a predefined template.

Available templates:
🔧 cli-tool  - CLI application (using Cobra)
📚 library   - Go library/package
🌐 web-api   - REST API / GraphQL server
⚙️  service   - Microservice/daemon

Examples:
  goossify create --template cli-tool my-cli-app
  goossify create --template library my-go-lib
  goossify create --template web-api my-api-server
  goossify create --template service my-service`,
	Args: cobra.ExactArgs(1),
	RunE: runCreate,
}

func init() {
	rootCmd.AddCommand(createCmd)

	createCmd.Flags().StringVarP(&createTemplate, "template", "t", "", "Template to use (cli-tool|library|web-api|service)")
	createCmd.Flags().StringVarP(&createAuthor, "author", "a", "", "Author name")
	createCmd.Flags().StringVarP(&createEmail, "email", "e", "", "Author email address")
	createCmd.Flags().StringVarP(&createLicense, "license", "l", "MIT", "License type")
	createCmd.Flags().StringVarP(&createGithub, "github", "g", "", "GitHub username")
	createCmd.Flags().BoolVarP(&createInteractive, "interactive", "i", false, "Interactive mode")

	_ = createCmd.MarkFlagRequired("template")
}

func runCreate(cmd *cobra.Command, args []string) error {
	projectName := args[0]

	// Validate template
	if !isValidTemplate(createTemplate) {
		return fmt.Errorf("invalid template: %s\nAvailable: cli-tool, library, web-api, service", createTemplate)
	}

	// Collect project configuration
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
			return fmt.Errorf("failed to collect configuration: %w", err)
		}
	}

	// Set default values
	setDefaultValues(config)

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

	fmt.Printf("🎉 Successfully created Go OSS project '%s' (%s)!\n\n", config.Name, config.Type)
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
		config.License = "MIT"
	}
	if config.GitHubUsername == "" {
		config.GitHubUsername = "your-username"
	}
	if config.Description == "" {
		config.Description = generateDescription(config.Type, config.Name)
	}
}

func printNextSteps(projectName string) {
	fmt.Println("Next steps:")
	fmt.Printf("  cd %s\n", projectName)
	fmt.Println("  go mod tidy")
	fmt.Println("  git init")
	fmt.Println("  git add .")
	fmt.Println("  git commit -m \"🎉 Initial commit\"")
	fmt.Println()
	fmt.Println("Create and push to GitHub repository:")
	fmt.Println("  gh repo create --public")
	fmt.Println("  git push -u origin main")
	fmt.Println()
	fmt.Println("Project management:")
	fmt.Println("  goossify status     # Check project health")
}
