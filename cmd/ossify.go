package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/pigeonworks-llc/goossify/internal/ossify"
)

var ossifyCmd = &cobra.Command{
	Use:   "ossify [path]",
	Short: "Convert existing project to OSS-ready",
	Long: `Enhance an existing Go project with files and setup required for OSS publication.

This command automatically generates and configures:
• LICENSE file
• .github/workflows/ci.yml (CI/CD configuration)
• Community files in .github/
• CONTRIBUTING.md, SECURITY.md, etc.
• Git initialization (if not initialized)`,
	Args: cobra.MaximumNArgs(1),
	RunE: runOssify,
}

func runOssify(cmd *cobra.Command, args []string) error {
	var targetPath string
	if len(args) == 0 {
		targetPath = "."
	} else {
		targetPath = args[0]
	}

	// Convert to absolute path
	absPath, err := filepath.Abs(targetPath)
	if err != nil {
		return fmt.Errorf("path resolution error: %w", err)
	}

	// Check directory existence
	if _, err := os.Stat(absPath); os.IsNotExist(err) {
		return fmt.Errorf("specified path does not exist: %s", absPath)
	}

	fmt.Printf("🚀 Starting OSS conversion: %s\n", absPath)

	// Initialize and execute Ossifier
	ossifier := ossify.New(absPath)
	if err := ossifier.Execute(); err != nil {
		return fmt.Errorf("error during OSS conversion: %w", err)
	}

	fmt.Println("✅ OSS conversion completed!")
	return nil
}

func init() {
	rootCmd.AddCommand(ossifyCmd)
}