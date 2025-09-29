package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile string
	verbose bool
)

// rootCmd represents the base command when called without any subcommands.
var rootCmd = &cobra.Command{
	Use:   "goossify",
	Short: "🚀 Complete automation tool for Go OSS project management",
	Long: `goossify - Go Language OSS Project Management Automation

Next-generation boilerplate generator & project management tool that completely
automates the creation, management, and maintenance of Go OSS projects.

Key Features:
🏗️  Fully automated project initialization
🤖  Continuous maintenance automation
📊  Quality & performance monitoring
👥  Community management support
🔄  Complete release management automation

Usage Examples:
  goossify init my-project      # Create new project
  goossify create --template cli-tool my-cli
  goossify ossify               # Convert existing project to OSS
  goossify status               # Check project health
  goossify ready                # Check release readiness
  goossify release              # Execute automated release

Documentation: https://github.com/pigeonworks-llc/goossify`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "Config file (default: $HOME/.goossify.yaml)")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Enable verbose output")

	_ = viper.BindPFlag("verbose", rootCmd.PersistentFlags().Lookup("verbose"))
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".goossify" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".goossify")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil && verbose {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}
