# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

goossify is a Go CLI application that automates the creation and management of Go OSS projects. It generates complete project structures with proper boilerplate code, configuration files, and GitHub workflows.

## Commands

### Build and Development
- `go build` - Build the main binary
- `go build -o goossify` - Build with specific output name
- `go run main.go [command]` - Run directly without building
- `go install` - Install the binary to GOPATH/bin

### Testing
- `go test ./...` - Run all tests recursively
- `go test -v ./...` - Run tests with verbose output
- `go test ./cmd` - Run tests in specific package
- `go test -race ./...` - Run tests with race detection

### Quality Checks
- `go vet ./...` - Run Go's built-in static analysis
- `go fmt ./...` - Format all Go code
- `gofumpt -w .` - More strict formatting (if available)

### Project Generation (Tool Usage)
- `goossify init [project-name]` - Interactive project creation
- `goossify create --template cli-tool --author "Name" --email "email@example.com" --github "username" project-name` - Create from template
- `goossify status` - Check project health

## Architecture

### Core Components

1. **CLI Framework**: Built with Cobra (`github.com/spf13/cobra`)
   - `cmd/root.go` - Root command and global configuration
   - `cmd/init.go` - Interactive project initialization
   - `cmd/create.go` - Template-based project creation
   - `cmd/status.go` - Project health checking

2. **Project Generator** (`internal/generator/generator.go`)
   - Main generation engine that creates complete project structures
   - Handles different project types: cli-tool, library, web-api, service
   - Uses Go templates to generate files with proper substitution

3. **Template System** (`internal/template/templates/`)
   - Contains all template definitions for different file types
   - Supports various project configurations and licenses
   - Generates GitHub workflows, documentation, and boilerplate code

### Project Structure Generated
The tool generates projects with this standard structure:
```
project/
├── cmd/                    # Entry points
├── internal/              # Private packages
├── pkg/                   # Public packages
├── docs/                  # Documentation
├── examples/              # Usage examples
├── tests/                 # Test files
├── .github/               # GitHub workflows and templates
├── go.mod                 # Go module definition
├── main.go                # Main entry point
└── various config files   # .golangci.yml, .goreleaser.yml, etc.
```

### Key Features
- **Multi-template support**: CLI tools, libraries, web APIs, and services
- **Complete automation**: Generates all necessary files for OSS projects
- **GitHub integration**: Creates workflows, issue templates, and community files
- **Configuration management**: Uses Viper for flexible configuration
- **Japanese language support**: Includes Japanese text and documentation

### Dependencies
- `github.com/spf13/cobra` - CLI framework
- `github.com/spf13/viper` - Configuration management
- `golang.org/x/text` - Text processing and internationalization
- `gopkg.in/yaml.v3` - YAML parsing

### Development Environment
- Uses `mise.toml` for Go version management (Go 1.25.1)
- Module requires Go 1.21+
- No external build tools required beyond standard Go toolchain