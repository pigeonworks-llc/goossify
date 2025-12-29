// Package templates provides all template content for goossify project generation.
package templates

// README.md template
const ReadmeTemplate = `# {{.Name}}

{{.Description}}

[![Go Report Card](https://goreportcard.com/badge/{{.ModulePath}})](https://goreportcard.com/report/{{.ModulePath}})
[![MIT License](https://img.shields.io/badge/license-{{.License}}-blue.svg)](LICENSE)
[![Release](https://img.shields.io/github/release/{{.GitHubUsername}}/{{.Name}}.svg)](https://github.com/{{.GitHubUsername}}/{{.Name}}/releases)

## 🚀 Features

- ✨ Feature 1
- 🔧 Feature 2
- 📊 Feature 3

## 📦 Installation

### Using go install

` + "```bash" + `
go install {{.ModulePath}}@latest
` + "```" + `

### Using releases

Download the latest release from [GitHub Releases](https://github.com/{{.GitHubUsername}}/{{.Name}}/releases).

## 🔧 Usage

` + "```bash" + `
{{.Name}} --help
` + "```" + `

## 🤝 Contributing

Contributions are welcome! Please see [CONTRIBUTING.md](.github/CONTRIBUTING.md) for details.

## 📄 License

This project is licensed under the {{.License}} License - see the [LICENSE](LICENSE) file for details.

## 👤 Author

**{{.Author}}**

- GitHub: [@{{.GitHubUsername}}](https://github.com/{{.GitHubUsername}})
- Email: {{.Email}}
`

// .gitignore template
const GitignoreTemplate = `# Binaries for programs and plugins
*.exe
*.exe~
*.dll
*.so
*.dylib

# Test binary, built with go test -c
*.test

# Output of the go coverage tool, specifically when used with LiteIDE
*.out

# Dependency directories (remove the comment below to include it)
# vendor/

# Go workspace file
go.work

# IDE files
.vscode/
.idea/
*.swp
*.swo
*~

# OS generated files
.DS_Store
.DS_Store?
._*
.Spotlight-V100
.Trashes
ehthumbs.db
Thumbs.db

# Temporary files
tmp/
temp/

# Log files
*.log

# Coverage files
coverage.txt
coverage.html

# Binary output
dist/
{{.Name}}

# Configuration files (if they contain secrets)
config.yaml
config.json
.env
.env.local

# GoReleaser
.goreleaser.yml.bak
`

// go.mod template
const GoModTemplate = `module {{.ModulePath}}

go 1.22

require (
	github.com/spf13/cobra v1.8.0
	github.com/spf13/viper v1.18.2
)
`

// MIT License template
const MITLicenseTemplate = `MIT License

Copyright (c) {{.Year}} {{.Author}}

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
`

// Apache 2.0 License template
const Apache2LicenseTemplate = `Apache License
Version 2.0, January 2004
http://www.apache.org/licenses/

TERMS AND CONDITIONS FOR USE, REPRODUCTION, AND DISTRIBUTION

1. Definitions.

   "License" shall mean the terms and conditions for use, reproduction,
   and distribution as defined by Sections 1 through 9 of this document.

   [Full Apache 2.0 license text would go here]

Copyright {{.Year}} {{.Author}}

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
`

// BSD 3-Clause License template
const BSD3LicenseTemplate = `BSD 3-Clause License

Copyright (c) {{.Year}}, {{.Author}}
All rights reserved.

Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions are met:

1. Redistributions of source code must retain the above copyright notice, this
   list of conditions and the following disclaimer.

2. Redistributions in binary form must reproduce the above copyright notice,
   this list of conditions and the following disclaimer in the documentation
   and/or other materials provided with the distribution.

3. Neither the name of the copyright holder nor the names of its
   contributors may be used to endorse or promote products derived from
   this software without specific prior written permission.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE
FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL
DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR
SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER
CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY,
OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
`

// Makefile template
const MakefileTemplate = `.PHONY: build test clean lint fmt vet coverage bench security ci release release-snapshot install tidy deps dev help

# Build variables
BINARY_NAME={{.Name}}
BUILD_DIR=dist
VERSION?=$(shell git describe --tags --always --dirty)
LDFLAGS=-ldflags "-X main.version=${VERSION}"

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
GOFMT=gofmt
GOVET=$(GOCMD) vet

# Default target
all: clean lint test build

# Build the binary
build:
	@echo "Building ${BINARY_NAME}..."
	@mkdir -p ${BUILD_DIR}
	${GOBUILD} ${LDFLAGS} -o ${BUILD_DIR}/${BINARY_NAME} .

# Run tests
test:
	@echo "Running tests..."
	${GOTEST} -v -race -coverprofile=coverage.out ./...

# Run tests with coverage
coverage: test
	@echo "Generating coverage report..."
	${GOCMD} tool cover -html=coverage.out -o coverage.html

# Run benchmarks
bench:
	@echo "Running benchmarks..."
	${GOTEST} -bench=. -benchmem ./...

# Run security checks
security:
	@echo "Running security checks..."
	@command -v govulncheck >/dev/null 2>&1 || { echo "govulncheck not installed. Run: go install golang.org/x/vuln/cmd/govulncheck@latest"; exit 1; }
	govulncheck ./...

# Clean build artifacts
clean:
	@echo "Cleaning..."
	${GOCLEAN}
	rm -rf ${BUILD_DIR}
	rm -f coverage.out coverage.html

# Format code
fmt:
	@echo "Formatting code..."
	${GOFMT} -s -w .

# Vet code
vet:
	@echo "Vetting code..."
	${GOVET} ./...

# Run linter
lint:
	@echo "Running linter..."
	golangci-lint run

# Tidy dependencies
tidy:
	@echo "Tidying dependencies..."
	${GOMOD} tidy
	${GOMOD} verify

# Install dependencies
deps:
	@echo "Installing dependencies..."
	${GOMOD} download
	${GOMOD} verify

# Development workflow
dev: fmt vet lint test

# CI workflow
ci: fmt vet lint test build
	@echo "All CI checks passed!"

# Release build (cross-platform)
release:
	@echo "Building release..."
	@command -v goreleaser >/dev/null 2>&1 || { echo "goreleaser not installed. See: https://goreleaser.com/install/"; exit 1; }
	goreleaser release --clean

# Release snapshot (local testing)
release-snapshot:
	@echo "Building release snapshot..."
	@command -v goreleaser >/dev/null 2>&1 || { echo "goreleaser not installed. See: https://goreleaser.com/install/"; exit 1; }
	goreleaser release --snapshot --clean

# Install the binary
install:
	@echo "Installing ${BINARY_NAME}..."
	${GOBUILD} ${LDFLAGS} -o ${GOPATH}/bin/${BINARY_NAME} .

# Help
help:
	@echo "Available targets:"
	@echo "  build            - Build the binary"
	@echo "  test             - Run tests"
	@echo "  coverage         - Run tests with coverage"
	@echo "  bench            - Run benchmarks"
	@echo "  security         - Run security checks (govulncheck)"
	@echo "  clean            - Clean build artifacts"
	@echo "  fmt              - Format code"
	@echo "  vet              - Vet code"
	@echo "  lint             - Run linter"
	@echo "  tidy             - Tidy dependencies"
	@echo "  deps             - Install dependencies"
	@echo "  dev              - Run development workflow"
	@echo "  ci               - Run all CI checks"
	@echo "  release          - Build release with goreleaser"
	@echo "  release-snapshot - Build snapshot release (local)"
	@echo "  install          - Install binary to GOPATH"
	@echo "  help             - Show this help"
`

// CHANGELOG.md template
const ChangelogTemplate = `# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- Initial project setup

### Changed

### Deprecated

### Removed

### Fixed

### Security

## [0.1.0] - {{.Year}}-01-01

### Added
- Initial release of {{.Name}}
- Basic project structure
- Core functionality

[Unreleased]: https://github.com/{{.GitHubUsername}}/{{.Name}}/compare/v0.1.0...HEAD
[0.1.0]: https://github.com/{{.GitHubUsername}}/{{.Name}}/releases/tag/v0.1.0
`

// Default main.go template
const DefaultMainTemplate = `package main

import (
	"fmt"
	"log"
)

func main() {
	fmt.Println("Hello from {{.Name}}!")
}
`

// CLI main entry template
const CLIMainEntryTemplate = `package main

import (
	"fmt"
	"os"

	"{{.ModulePath}}/internal/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
`

// CLI main template
const CLIMainTemplate = `package main

import (
	"{{.ModulePath}}/internal/cmd"
)

func main() {
	cmd.Execute()
}
`

// CLI root command template
const CLIRootTemplate = `package cmd

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

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "{{.Name}}",
	Short: "{{.Description}}",
	Long: "{{.Description}}",
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.{{.Name}}.yaml)")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose output")

	viper.BindPFlag("verbose", rootCmd.PersistentFlags().Lookup("verbose"))
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".{{.Name}}")
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil && verbose {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}
`

// CLI version command template
const CLIVersionTemplate = `package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"{{.ModulePath}}/internal/version"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show version information",
	Long:  "Display the current version of {{.Name}}",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("{{.Name}} version %s\n", version.Version)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
`

// Version template
const VersionTemplate = `package version

// Version represents the current version of {{.Name}}
const Version = "0.1.0"
`

// Library main template
const LibraryMainTemplate = `package {{.PackageName}}

// {{.StructName}} represents the main functionality
type {{.StructName}} struct {
	// Add your fields here
}

// New creates a new {{.StructName}} instance
func New() *{{.StructName}} {
	return &{{.StructName}}{}
}

// Example method
func (c *{{.StructName}}) DoSomething() string {
	return "Hello from {{.StructName}}!"
}
`

// Library test template
const LibraryTestTemplate = `package {{.PackageName}}

import (
	"testing"
)

func TestNew(t *testing.T) {
	c := New()
	if c == nil {
		t.Error("New() should not return nil")
	}
}

func TestDoSomething(t *testing.T) {
	c := New()
	result := c.DoSomething()
	expected := "Hello from {{.StructName}}!"

	if result != expected {
		t.Errorf("DoSomething() = %v, want %v", result, expected)
	}
}

func BenchmarkDoSomething(b *testing.B) {
	c := New()
	for i := 0; i < b.N; i++ {
		c.DoSomething()
	}
}
`

// Library example template
const LibraryExampleTemplate = `package main

import (
	"fmt"

	"{{.ModulePath}}/pkg/{{.PackageName}}"
)

func main() {
	client := {{.PackageName}}.New()
	result := client.DoSomething()
	fmt.Println(result)
}
`

// Library doc template
const LibraryDocTemplate = `// Package {{.PackageName}} {{.Description}}
//
// This package provides the core functionality for {{.StructName}}.
//
// Example usage:
//
//	client := {{.PackageName}}.New()
//	result := client.DoSomething()
//	fmt.Println(result)
//
package {{.PackageName}}
`

// Web API main template
const WebAPIMainTemplate = `package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"{{.ModulePath}}/internal/server"
)

func main() {
	srv := server.New()

	// Start server
	go func() {
		log.Printf("Starting server on :8080")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed to start: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited")
}
`

// Web API handler template
const WebAPIHandlerTemplate = `package handler

import (
	"encoding/json"
	"net/http"

	"{{.ModulePath}}/internal/model"
)

// Handler represents the HTTP handler
type Handler struct{}

// New creates a new handler
func New() *Handler {
	return &Handler{}
}

// Health handles health check requests
func (h *Handler) Health(w http.ResponseWriter, r *http.Request) {
	response := model.Response{
		Status:  "ok",
		Message: "{{.Name}} is running",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// Hello handles hello requests
func (h *Handler) Hello(w http.ResponseWriter, r *http.Request) {
	response := model.Response{
		Status:  "success",
		Message: "Hello from {{.Name}}!",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
`

// Web API CORS template
const WebAPICORSTemplate = `package middleware

import (
	"net/http"
)

// CORS handles Cross-Origin Resource Sharing
func CORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
`

// Web API model template
const WebAPIModelTemplate = `package model

// Response represents a standard API response
type Response struct {
	Status  string      ` + "`json:\"status\"`" + `
	Message string      ` + "`json:\"message\"`" + `
	Data    interface{} ` + "`json:\"data,omitempty\"`" + `
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Status  string ` + "`json:\"status\"`" + `
	Error   string ` + "`json:\"error\"`" + `
	Code    int    ` + "`json:\"code\"`" + `
}
`

// Web API server template
const WebAPIServerTemplate = `package server

import (
	"net/http"

	"{{.ModulePath}}/internal/handler"
	"{{.ModulePath}}/internal/middleware"
)

// New creates a new HTTP server
func New() *http.Server {
	h := handler.New()

	mux := http.NewServeMux()

	// Routes
	mux.HandleFunc("/health", h.Health)
	mux.HandleFunc("/hello", h.Hello)

	// Apply middleware
	handler := middleware.CORS(mux)

	return &http.Server{
		Addr:    ":8080",
		Handler: handler,
	}
}
`

// OpenAPI template
const OpenAPITemplate = `openapi: 3.0.3
info:
  title: {{.Name}} API
  description: {{.Description}}
  version: 1.0.0
  contact:
    name: {{.Author}}
    email: {{.Email}}
  license:
    name: {{.License}}

servers:
  - url: http://localhost:8080
    description: Development server

paths:
  /health:
    get:
      summary: Health check
      description: Returns the health status of the service
      responses:
        '200':
          description: Service is healthy
          content:
            application/json:
              schema:
                type: object
                properties:
                  status:
                    type: string
                    example: ok
                  message:
                    type: string
                    example: {{.Name}} is running

  /hello:
    get:
      summary: Hello endpoint
      description: Returns a hello message
      responses:
        '200':
          description: Hello message
          content:
            application/json:
              schema:
                type: object
                properties:
                  status:
                    type: string
                    example: success
                  message:
                    type: string
                    example: Hello from {{.Name}}!
`

// Service templates will be added in the next part...
const ServiceMainTemplate = `package service

import (
	"context"
	"log"
)

// Service represents the main service
type Service struct {
	// Add your dependencies here
}

// New creates a new service instance
func New() *Service {
	return &Service{}
}

// Start starts the service
func (s *Service) Start(ctx context.Context) error {
	log.Println("Starting {{.Name}} service...")

	// Service logic here

	return nil
}

// Stop stops the service
func (s *Service) Stop(ctx context.Context) error {
	log.Println("Stopping {{.Name}} service...")

	// Cleanup logic here

	return nil
}
`

const ServiceRepositoryTemplate = `package repository

import (
	"context"
)

// Repository interface defines the data access methods
type Repository interface {
	// Add your repository methods here
	Get(ctx context.Context, id string) (interface{}, error)
	Save(ctx context.Context, data interface{}) error
}

// repository implements the Repository interface
type repository struct {
	// Add your database connections or other dependencies here
}

// New creates a new repository instance
func New() Repository {
	return &repository{}
}

// Get retrieves data by ID
func (r *repository) Get(ctx context.Context, id string) (interface{}, error) {
	// Implementation here
	return nil, nil
}

// Save saves data
func (r *repository) Save(ctx context.Context, data interface{}) error {
	// Implementation here
	return nil
}
`

const ServiceConfigTemplate = `package config

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

// Config represents the service configuration
type Config struct {
	Port     int    ` + "`mapstructure:\"port\"`" + `
	LogLevel string ` + "`mapstructure:\"log_level\"`" + `
	Database DatabaseConfig ` + "`mapstructure:\"database\"`" + `
}

// DatabaseConfig represents database configuration
type DatabaseConfig struct {
	Host     string ` + "`mapstructure:\"host\"`" + `
	Port     int    ` + "`mapstructure:\"port\"`" + `
	Username string ` + "`mapstructure:\"username\"`" + `
	Password string ` + "`mapstructure:\"password\"`" + `
	Database string ` + "`mapstructure:\"database\"`" + `
}

// Load loads configuration from file and environment
func Load() (*Config, error) {
	viper.SetDefault("port", 8080)
	viper.SetDefault("log_level", "info")

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./config")

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("error reading config file: %w", err)
		}
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("error unmarshaling config: %w", err)
	}

	return &config, nil
}
`

const ServiceMainEntryTemplate = `package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"{{.ModulePath}}/internal/service"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Create service
	svc := service.New()

	// Start service
	go func() {
		if err := svc.Start(ctx); err != nil {
			log.Fatalf("Service failed to start: %v", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down service...")

	// Stop service
	if err := svc.Stop(ctx); err != nil {
		log.Fatalf("Service failed to stop gracefully: %v", err)
	}

	log.Println("Service exited")
}
`

const DockerfileTemplate = `FROM golang:1.21-alpine AS builder

WORKDIR /app

# Install dependencies
RUN apk add --no-cache git

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o {{.Name}} .

FROM alpine:latest

RUN apk --no-cache add ca-certificates
WORKDIR /root/

# Copy the binary from builder
COPY --from=builder /app/{{.Name}} .

# Expose port
EXPOSE 8080

# Run the binary
CMD ["./{{.Name}}"]
`

const KubernetesTemplate = `apiVersion: apps/v1
kind: Deployment
metadata:
  name: {{.Name}}
  labels:
    app: {{.Name}}
spec:
  replicas: 3
  selector:
    matchLabels:
      app: {{.Name}}
  template:
    metadata:
      labels:
        app: {{.Name}}
    spec:
      containers:
      - name: {{.Name}}
        image: {{.Name}}:latest
        ports:
        - containerPort: 8080
        env:
        - name: PORT
          value: "8080"
        livenessProbe:
          httpGet:
            path: /health
            port: 8080
          initialDelaySeconds: 30
          periodSeconds: 10
        readinessProbe:
          httpGet:
            path: /health
            port: 8080
          initialDelaySeconds: 5
          periodSeconds: 5
---
apiVersion: v1
kind: Service
metadata:
  name: {{.Name}}-service
spec:
  selector:
    app: {{.Name}}
  ports:
    - protocol: TCP
      port: 80
      targetPort: 8080
  type: LoadBalancer
`

// GitHub Actions CI template
const GitHubCITemplate = `name: CI

on:
  push:
    branches: [ main, develop ]
  pull_request:
    branches: [ main ]

env:
  GO_VERSION: "1.21"

jobs:
  test:
    name: Test
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@v4

    - name: Set up Go
      uses: actions/setup-go@v4
      with:
        go-version: "1.22"

    - name: Cache dependencies
      uses: actions/cache@v3
      with:
        path: |
          ~/.cache/go-build
          ~/go/pkg/mod
        key: ubuntu-latest-go-1.21-go-sum
        restore-keys: |
          ubuntu-latest-go-1.21-

    - name: Download dependencies
      run: go mod download

    - name: Verify dependencies
      run: go mod verify

    - name: Run tests
      run: go test -v -race -coverprofile=coverage.out ./...

    - name: Generate coverage report
      run: go tool cover -html=coverage.out -o coverage.html

    - name: Upload coverage to Codecov
      uses: codecov/codecov-action@v3
      with:
        file: ./coverage.out
        flags: unittests
        name: codecov-umbrella

  lint:
    name: Lint
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@v4

    - name: Set up Go
      uses: actions/setup-go@v4
      with:
        go-version: "1.22"

    - name: golangci-lint
      uses: golangci/golangci-lint-action@v3
      with:
        version: latest
        args: --timeout=5m

  security:
    name: Security Scan
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@v4

    - name: Set up Go
      uses: actions/setup-go@v4
      with:
        go-version: "1.22"

    - name: Run Gosec Security Scanner
      uses: securecodewarrior/github-action-gosec@master
      with:
        args: '-no-fail -fmt sarif -out results.sarif ./...'

    - name: Upload SARIF file
      uses: github/codeql-action/upload-sarif@v2
      with:
        sarif_file: results.sarif

  build:
    name: Build
    runs-on: ubuntu-latest
    needs: [test, lint]
    steps:
    - uses: actions/checkout@v4

    - name: Set up Go
      uses: actions/setup-go@v4
      with:
        go-version: "1.22"

    - name: Build
      run: go build -v ./...

    - name: Build for multiple platforms
      run: |
        GOOS=linux GOARCH=amd64 go build -o dist/{{.Name}}-linux-amd64 .
        GOOS=windows GOARCH=amd64 go build -o dist/{{.Name}}-windows-amd64.exe .
        GOOS=darwin GOARCH=amd64 go build -o dist/{{.Name}}-darwin-amd64 .
        GOOS=darwin GOARCH=arm64 go build -o dist/{{.Name}}-darwin-arm64 .

    - name: Upload build artifacts
      uses: actions/upload-artifact@v3
      with:
        name: binaries
        path: dist/
`

// GitHub Actions Release template
const GitHubReleaseTemplate = `name: Release

on:
  push:
    tags:
      - 'v*'

permissions:
  contents: write
  packages: write
  issues: write

jobs:
  goreleaser:
    runs-on: ubuntu-latest
    steps:
    - name: Checkout
      uses: actions/checkout@v4
      with:
        fetch-depth: 0

    - name: Set up Go
      uses: actions/setup-go@v4
      with:
        go-version: "1.22"

    - name: Run GoReleaser
      uses: goreleaser/goreleaser-action@v5
      with:
        distribution: goreleaser
        version: latest
        args: release --clean
      env:
        GITHUB_TOKEN: "SET_YOUR_GITHUB_TOKEN"

  docker:
    runs-on: ubuntu-latest
    steps:
    - name: Checkout
      uses: actions/checkout@v4

    - name: Set up Docker Buildx
      uses: docker/setup-buildx-action@v3

    - name: Log in to Docker Hub
      uses: docker/login-action@v3
      with:
        username: SET_YOUR_DOCKER_USERNAME
        password: SET_YOUR_DOCKER_PASSWORD

    - name: Log in to GitHub Container Registry
      uses: docker/login-action@v3
      with:
        registry: ghcr.io
        username: SET_YOUR_GITHUB_USERNAME
        password: SET_YOUR_GITHUB_TOKEN

    - name: Build and push
      uses: docker/build-push-action@v5
      with:
        context: .
        platforms: linux/amd64,linux/arm64
        push: true
        tags: {{.GitHubUsername}}/{{.Name}}:latest
`

// Auto Labeler workflow template
const AutoLabelerTemplate = `name: Auto Label

on:
  pull_request:
    types: [opened, synchronize]
  issues:
    types: [opened]

jobs:
  label:
    runs-on: ubuntu-latest
    steps:
    - name: Apply Labels
      uses: actions/labeler@v4
      with:
        repo-token: SET_YOUR_GITHUB_TOKEN
        configuration-path: .github/labeler.yml
        sync-labels: true
`

// Labeler configuration template
const LabelerConfigTemplate = `# Labeler configuration for {{.Name}}

# Component labels based on file changes
"component/core":
  - "internal/**/*"
  - "pkg/**/*"

"component/cli":
  - "cmd/**/*"
  - "main.go"

"component/docs":
  - "docs/**/*"
  - "*.md"
  - "**/*.md"

"component/tests":
  - "**/*_test.go"
  - "test/**/*"
  - "tests/**/*"

"component/ci-cd":
  - ".github/**/*"
  - ".golangci.yml"
  - ".goreleaser.yml"
  - "Makefile"
  - "Dockerfile"

"type/documentation":
  - "*.md"
  - "docs/**/*"
  - "**/*.md"

"dependencies":
  - "go.mod"
  - "go.sum"

# Size-based labels
"effort/xs":
  - any: ["*.md"]
    all: ["!docs/**/*"]

"effort/s":
  - "*_test.go"
  - "docs/**/*"

"effort/m":
  - "cmd/**/*"

"effort/l":
  - "internal/**/*"
  - "pkg/**/*"
`

// Dependabot configuration template
const DependabotTemplate = `# Dependabot configuration for {{.Name}}

version: 2
updates:
  # Go modules
  - package-ecosystem: "gomod"
    directory: "/"
    schedule:
      interval: "weekly"
      day: "monday"
      time: "09:00"
    open-pull-requests-limit: 5
    reviewers:
      - "{{.GitHubUsername}}"
    assignees:
      - "{{.GitHubUsername}}"
    commit-message:
      prefix: "deps"
      include: "scope"
    labels:
      - "dependencies"
      - "type/maintenance"

  # GitHub Actions
  - package-ecosystem: "github-actions"
    directory: "/"
    schedule:
      interval: "weekly"
      day: "monday"
      time: "09:00"
    open-pull-requests-limit: 3
    reviewers:
      - "{{.GitHubUsername}}"
    assignees:
      - "{{.GitHubUsername}}"
    commit-message:
      prefix: "ci"
      include: "scope"
    labels:
      - "component/ci-cd"
      - "type/maintenance"

  # Docker (if Dockerfile exists)
  - package-ecosystem: "docker"
    directory: "/"
    schedule:
      interval: "weekly"
    open-pull-requests-limit: 2
    labels:
      - "dependencies"
      - "component/docker"
`

// CodeQL workflow template
const CodeQLTemplate = `name: CodeQL Security Scan

on:
  push:
    branches: [ main, develop ]
  pull_request:
    branches: [ main ]
  schedule:
    - cron: '0 6 * * 1'  # Weekly on Mondays

jobs:
  analyze:
    name: Analyze
    runs-on: ubuntu-latest
    permissions:
      actions: read
      contents: read
      security-events: write

    strategy:
      fail-fast: false
      matrix:
        language: [ 'go' ]

    steps:
    - name: Checkout repository
      uses: actions/checkout@v4

    - name: Initialize CodeQL
      uses: github/codeql-action/init@v2
      with:
        languages: golang
        queries: security-extended,security-and-quality

    - name: Set up Go
      uses: actions/setup-go@v4
      with:
        go-version: "1.22"

    - name: Build
      run: |
        go mod download
        go build -v ./...

    - name: Perform CodeQL Analysis
      uses: github/codeql-action/analyze@v2
      with:
        category: "/language:go"
`

// CODEOWNERS template
const CodeOwnersTemplate = `# Code owners for {{.Name}}

# Global owners
* @{{.GitHubUsername}}

# Core functionality
/internal/ @{{.GitHubUsername}}
/pkg/ @{{.GitHubUsername}}

# CLI commands
/cmd/ @{{.GitHubUsername}}

# Documentation
*.md @{{.GitHubUsername}}
/docs/ @{{.GitHubUsername}}

# CI/CD and project configuration
/.github/ @{{.GitHubUsername}}
/.golangci.yml @{{.GitHubUsername}}
/.goreleaser.yml @{{.GitHubUsername}}
/Makefile @{{.GitHubUsername}}
/Dockerfile @{{.GitHubUsername}}

# Go modules
/go.mod @{{.GitHubUsername}}
/go.sum @{{.GitHubUsername}}

# Add more specific owners as your team grows:
# /internal/api/ @api-team
# /internal/database/ @database-team
# /docs/ @docs-team
`

// Funding configuration template
const FundingTemplate = `# Funding configuration for {{.Name}}

# GitHub Sponsors
github: [{{.GitHubUsername}}]

# Open Collective
# open_collective: your-collective-name

# Ko-fi
# ko_fi: your-kofi-username

# Tidelift
# tidelift: npm/{{.Name}}

# Community Bridge
# community_bridge: your-project-name

# Liberapay
# liberapay: your-liberapay-username

# IssueHunt
# issuehunt: your-issuehunt-username

# Buy Me A Coffee
# buy_me_a_coffee: your-bmac-username

# Patreon
# patreon: your-patreon-username

# Custom URL (uncomment and add your donation link)
# custom: ["https://paypal.me/your-paypal", "https://example.com/donate"]
`

// Project management workflow template
const ProjectManagementTemplate = `name: Project Management

on:
  issues:
    types: [opened, reopened, closed, labeled]
  pull_request:
    types: [opened, reopened, closed, merged, labeled]

jobs:
  manage-project:
    runs-on: ubuntu-latest
    steps:
    - name: Add to Project Board
      if: github.event.action == 'opened'
      uses: actions/add-to-project@v0.5.0
      with:
        project-url: https://github.com/users/{{.GitHubUsername}}/projects/1
        github-token: SET_YOUR_GITHUB_TOKEN

    - name: Move completed issues
      if: github.event.action == 'closed' && github.event.issue.state_reason == 'completed'
      uses: actions/add-to-project@v0.5.0
      with:
        project-url: https://github.com/users/{{.GitHubUsername}}/projects/1
        github-token: SET_YOUR_GITHUB_TOKEN
        labeled: status/done

  auto-assign:
    runs-on: ubuntu-latest
    if: github.event.action == 'opened'
    steps:
    - name: Auto-assign issue to author
      if: github.event_name == 'issues'
      uses: pozil/auto-assign-issue@v1
      with:
        assignees: {{.GitHubUsername}}
        numOfAssignee: 1

    - name: Auto-assign PR to author
      if: github.event_name == 'pull_request'
      uses: kentaro-m/auto-assign-action@v1.2.5
      with:
        configuration-path: '.github/auto-assign.yml'

  stale-management:
    runs-on: ubuntu-latest
    if: github.event.action == 'schedule'
    steps:
    - name: Close stale issues and PRs
      uses: actions/stale@v8
      with:
        repo-token: SET_YOUR_GITHUB_TOKEN
        stale-issue-message: |
          This issue has been automatically marked as stale because it has not had recent activity.
          It will be closed if no further activity occurs. Thank you for your contributions.
        stale-pr-message: |
          This pull request has been automatically marked as stale because it has not had recent activity.
          It will be closed if no further activity occurs. Thank you for your contributions.
        close-issue-message: |
          This issue was closed because it has been stale for 30 days with no activity.
        close-pr-message: |
          This pull request was closed because it has been stale for 30 days with no activity.
        days-before-stale: 60
        days-before-close: 30
        stale-issue-label: status/stale
        stale-pr-label: status/stale
        exempt-issue-labels: status/accepted,priority/high,priority/critical
        exempt-pr-labels: status/in-progress,priority/high,priority/critical
`

// Auto-assign configuration template
const AutoAssignTemplate = `# Auto-assign configuration for {{.Name}}

# Set to true to add reviewers to pull requests
addReviewers: true

# Set to true to add assignees to pull requests
addAssignees: false

# A list of reviewers to be added to pull requests (GitHub user name)
reviewers:
  - {{.GitHubUsername}}

# A number of reviewers added to the pull request
numberOfReviewers: 1

# A list of assignees, overrides reviewers if set
# assignees:
#   - {{.GitHubUsername}}

# A number of assignees to add to the pull request
# numberOfAssignees: 1

# A list of keywords to be skipped the process that add reviewers if pull requests include it
skipKeywords:
  - wip
  - draft
  - "[WIP]"
  - "[Draft]"
`

// Issue forms configuration template
const IssueFormsConfigTemplate = `# Issue forms configuration for {{.Name}}

blank_issues_enabled: true

contact_links:
  - name: "💬 GitHub Discussions"
    url: https://github.com/{{.GitHubUsername}}/{{.Name}}/discussions
    about: "Ask questions and discuss ideas with the community"

  - name: "📖 Documentation"
    url: https://github.com/{{.GitHubUsername}}/{{.Name}}/wiki
    about: "Check our documentation for help and guides"

  - name: "🔒 Security Issue"
    url: https://github.com/{{.GitHubUsername}}/{{.Name}}/security/advisories/new
    about: "Report a security vulnerability (private disclosure)"

  - name: "📧 Direct Contact"
    url: mailto:{{.Email}}
    about: "For sensitive issues or private discussions"
`

// Bug report template
const BugReportTemplate = `---
name: 🐛 Bug Report
about: バグを報告する
title: '[BUG] '
labels: bug
assignees: ''

---

## 🐛 Bug Description
バグの詳細な説明をお書きください。

## 🔄 To Reproduce
バグを再現する手順:
1. Go to '...'
2. Click on '....'
3. Scroll down to '....'
4. See error

## ✅ Expected Behavior
期待される動作の説明をお書きください。

## 📷 Screenshots
可能であれば、スクリーンショットを添付してください。

## 🌍 Environment
**Desktop (please complete the following information):**
 - OS: [e.g. macOS, Linux, Windows]
 - Version: [e.g. {{.Name}} v1.0.0]
 - Go Version: [e.g. 1.21.0]

## 📄 Additional Context
その他の関連情報があればお書きください。

## 📝 Error Logs
` + "```" + `
エラーログがある場合はここに貼り付けてください
` + "```" + `
`

// Feature request template
const FeatureRequestTemplate = `---
name: ✨ Feature Request
about: 新機能を提案する
title: '[FEATURE] '
labels: enhancement
assignees: ''

---

## ✨ Feature Description
実装したい機能の詳細な説明をお書きください。

## 🎯 Motivation
この機能が必要な理由や解決したい問題について説明してください。

## 💡 Proposed Solution
どのように実装すべきかのアイデアがあれば説明してください。

## 🔄 Alternatives Considered
検討した代替案があれば説明してください。

## 📷 Mockups/Examples
モックアップや例があれば添付してください。

## 📄 Additional Context
その他の関連情報があればお書きください。

## ✅ Acceptance Criteria
この機能が完成したと判断するための基準:
- [ ] Criterion 1
- [ ] Criterion 2
- [ ] Criterion 3
`

// Pull request template
const PRTemplate = `<!--
Thank you for contributing to {{.Name}}!
Please take a moment to fill out this template to help us review your changes efficiently.
-->

## 📋 Description

<!-- Provide a clear and concise description of what this PR does -->

### Summary
<!-- Brief summary of changes -->

### Motivation and Context
<!-- Why is this change required? What problem does it solve? -->
<!-- If it fixes an open issue, please link to the issue here -->

**Related Issue(s):**
- Closes #(issue number)
- Fixes #(issue number)
- Related to #(issue number)

## 🎯 Type of Change

**What type of change does this PR introduce?** *(Check all that apply)*

### Core Changes
- [ ] 🐛 **Bug fix** - Non-breaking change that fixes an issue
- [ ] ✨ **New feature** - Non-breaking change that adds functionality
- [ ] 💥 **Breaking change** - Change that would cause existing functionality to not work as expected
- [ ] ⚡ **Performance improvement** - Change that improves performance

### Maintenance
- [ ] 🔧 **Refactoring** - Code change that neither fixes a bug nor adds a feature
- [ ] 📝 **Documentation** - Changes to documentation only
- [ ] 🧪 **Tests** - Adding missing tests or correcting existing tests
- [ ] 🎨 **Code style** - Changes that do not affect the meaning of the code
- [ ] 🔨 **Build/CI** - Changes to build system or CI configuration

### Dependencies
- [ ] ⬆️ **Dependency upgrade** - Upgrading a dependency
- [ ] ➕ **New dependency** - Adding a new dependency
- [ ] ➖ **Remove dependency** - Removing a dependency

## 🧪 Testing

### Test Strategy
<!-- How did you test these changes? -->

**Automated Testing:**
- [ ] Unit tests added/updated and passing
- [ ] Integration tests added/updated and passing
- [ ] End-to-end tests added/updated and passing
- [ ] Performance tests added/updated (if applicable)

**Manual Testing:**
- [ ] Tested locally with different scenarios
- [ ] Tested on different platforms/environments
- [ ] Verified backward compatibility

### Test Evidence
<!-- Paste output, screenshots, or provide evidence of testing -->

` + "```bash" + `
# Example: Test output
$ make test
PASS
coverage: 85.4% of statements
` + "```" + `

## 📊 Performance Impact

<!-- If this change affects performance, provide benchmarks -->

**Before/After Comparison:**
- [ ] No performance impact
- [ ] Performance improved (provide benchmarks)
- [ ] Performance decreased (justify why)

<!--
` + "```bash" + `
# Example benchmark
$ go test -bench=. -benchmem
BenchmarkOldFunction-8    1000000    1234 ns/op    456 B/op    7 allocs/op
BenchmarkNewFunction-8    2000000     678 ns/op    234 B/op    3 allocs/op
` + "```" + `
-->

## 💼 Breaking Changes

<!-- If this is a breaking change, describe the impact and migration path -->

**Impact:**
- [ ] API changes
- [ ] Configuration changes
- [ ] Behavior changes
- [ ] CLI interface changes

**Migration Guide:**
<!-- Provide clear steps for users to migrate -->

## 📋 Pre-Submission Checklist

### Code Quality
- [ ] My code follows the project's [style guidelines](.github/CONTRIBUTING.md)
- [ ] I have performed a self-review of my own code
- [ ] I have commented my code, particularly in hard-to-understand areas
- [ ] My changes generate no new compiler warnings
- [ ] I have added error handling for all new code paths

### Testing
- [ ] I have added tests that prove my fix is effective or that my feature works
- [ ] New and existing unit tests pass locally with my changes
- [ ] I have tested my changes with realistic data/scenarios

### Documentation
- [ ] I have made corresponding changes to the documentation
- [ ] My changes don't require documentation updates
- [ ] I have updated the CHANGELOG.md (if applicable)
- [ ] I have updated command help text (if applicable)

### Dependencies
- [ ] I have updated go.mod and go.sum appropriately
- [ ] I have verified that new dependencies are necessary and well-maintained
- [ ] I have checked for license compatibility

### Security
- [ ] My changes don't introduce security vulnerabilities
- [ ] I have considered the security implications of my changes
- [ ] I have updated security documentation if needed

## 📸 Screenshots/Recordings

<!-- If your changes affect the UI or CLI output, provide before/after screenshots -->

### Before
<!-- Screenshots showing current behavior -->

### After
<!-- Screenshots showing new behavior -->

## 📄 Additional Context

### Implementation Details
<!-- Technical details about your implementation approach -->

### Alternatives Considered
<!-- Other approaches you considered and why you chose this one -->

### Future Work
<!-- Any follow-up work this PR enables or requires -->

### Notes for Reviewers
<!-- Anything specific you want reviewers to focus on -->

---

## 📝 Reviewer Guidelines

**For reviewers, please verify:**

- [ ] **Code Quality**: Code is readable, maintainable, and follows project conventions
- [ ] **Testing**: Adequate test coverage and tests are meaningful
- [ ] **Performance**: No unnecessary performance regressions
- [ ] **Security**: No security vulnerabilities introduced
- [ ] **Documentation**: Changes are documented appropriately
- [ ] **Breaking Changes**: Breaking changes are justified and migration path is clear

---

<!--
🚀 Thank you for contributing to {{.Name}}!
Your contribution helps make this project better for everyone.
-->
`

// Contributing guide template
const ContributingTemplate = `# Contributing to {{.Name}}

{{.Name}}への貢献をお考えいただき、ありがとうございます！🎉

## 🤝 How to Contribute

### 🐛 Reporting Bugs

バグを発見した場合は、[Bug Report](.github/ISSUE_TEMPLATE/bug_report.md)テンプレートを使用してIssueを作成してください。

### ✨ Suggesting Features

新機能のアイデアがある場合は、[Feature Request](.github/ISSUE_TEMPLATE/feature_request.md)テンプレートを使用してIssueを作成してください。

### 💻 Code Contributions

1. **Fork & Clone**
   ` + "```bash" + `
   git clone https://github.com/{{.GitHubUsername}}/{{.Name}}.git
   cd {{.Name}}
   ` + "```" + `

2. **Set up development environment**
   ` + "```bash" + `
   go mod download
   go mod verify
   ` + "```" + `

3. **Create a branch**
   ` + "```bash" + `
   git checkout -b feature/your-feature-name
   ` + "```" + `

4. **Make your changes**
   - コードスタイルガイドに従ってください
   - テストを追加/更新してください
   - ドキュメントを更新してください

5. **Test your changes**
   ` + "```bash" + `
   make test
   make lint
   ` + "```" + `

6. **Commit your changes**
   ` + "```bash" + `
   git add .
   git commit -m "feat: add your feature description"
   ` + "```" + `

7. **Push and create PR**
   ` + "```bash" + `
   git push origin feature/your-feature-name
   ` + "```" + `

## 📋 Development Guidelines

### 🏗️ Code Style

- [gofmt](https://golang.org/cmd/gofmt/)でコードをフォーマットしてください
- [golangci-lint](https://golangci-lint.run/)を使用してリンターチェックを通してください
- 関数とパッケージにドキュメントコメントを追加してください

### 🧪 Testing

- 新機能にはテストを追加してください
- テストカバレッジは80%以上を維持してください
- ` + "`go test -race ./...`" + `ですべてのテストが通ることを確認してください

### 📝 Commit Message Convention

[Conventional Commits](https://www.conventionalcommits.org/)に従ってコミットメッセージを記述してください：

` + "```" + `
type(scope): description

body

footer
` + "```" + `

**Types:**
- ` + "`feat`" + `: 新機能
- ` + "`fix`" + `: バグ修正
- ` + "`docs`" + `: ドキュメント変更
- ` + "`style`" + `: コードフォーマット変更
- ` + "`refactor`" + `: リファクタリング
- ` + "`test`" + `: テスト追加・修正
- ` + "`chore`" + `: その他の変更

**Example:**
` + "```" + `
feat(cli): add new status command

Add status command to check project health and goossify configuration.

Closes #123
` + "```" + `

### 🔧 Local Development

#### Prerequisites

- Go 1.21 or later
- Git
- Make

#### Setup

` + "```bash" + `
# Clone the repository
git clone https://github.com/{{.GitHubUsername}}/{{.Name}}.git
cd {{.Name}}

# Install dependencies
go mod download

# Run tests
make test

# Run linter
make lint

# Build
make build
` + "```" + `

#### Available Make Targets

- ` + "`make test`" + `: Run all tests
- ` + "`make lint`" + `: Run linter
- ` + "`make build`" + `: Build the binary
- ` + "`make clean`" + `: Clean build artifacts
- ` + "`make coverage`" + `: Generate coverage report

## 📄 Code of Conduct

このプロジェクトは[Code of Conduct](CODE_OF_CONDUCT.md)に従います。参加することで、この規約を守ることに同意したものとみなされます。

## 🆘 Getting Help

質問がある場合は：

1. [Discussions](https://github.com/{{.GitHubUsername}}/{{.Name}}/discussions)で質問してください
2. [Issues](https://github.com/{{.GitHubUsername}}/{{.Name}}/issues)で報告してください
3. メールでお問い合わせください: {{.Email}}

## 🙏 Recognition

貢献者の皆様は[CONTRIBUTORS.md](CONTRIBUTORS.md)で認識されます。

---

再度、{{.Name}}への貢献をご検討いただき、ありがとうございます！🚀
`

// Code of conduct template
const CodeOfConductTemplate = `# Contributor Covenant Code of Conduct

## Our Pledge

私たち{{.Name}}コミュニティのメンバー、貢献者、およびリーダーは、年齢、体型、目に見える、または目に見えない障害、民族性、性的特徴、性同一性と表現、経験レベル、教育、社会経済的地位、国籍、外見、人種、宗教、または性的アイデンティティと指向に関係なく、すべての人にとってハラスメントのない参加を保証することを誓います。

私たちは、オープンで歓迎的で、多様で包括的で健全なコミュニティに貢献する方法で行動し、交流することを誓います。

## Our Standards

ポジティブな環境を作ることに貢献する行動の例：

* 他の人々に対する共感と親切を示すこと
* 異なる意見、視点、経験を尊重すること
* 建設的なフィードバックを与え、それを受け入れること
* 私たちの間違いの影響を受けた人々に責任を受け入れ、謝罪し、その経験から学ぶこと
* 個人としてだけでなく、コミュニティ全体にとって最善のことに焦点を当てること

受け入れられない行動の例：

* 性的な言葉や画像の使用、およびあらゆる種類の性的注意または進歩
* トローリング、侮辱的または軽蔑的なコメント、個人的または政治的攻撃
* 公的または私的なハラスメント
* 明示的な許可なしに、物理的または電子メールアドレスなどの他の人の個人情報を公開すること
* プロフェッショナルな環境で合理的に不適切と見なされる可能性のあるその他の行為

## Enforcement Responsibilities

コミュニティリーダーは、受け入れられる行動の基準を明確にし、実施する責任があり、不適切、脅迫的、攻撃的、または有害と考える行動に対して適切で公正な是正措置を講じます。

コミュニティリーダーは、この行動規範に沿わないコメント、コミット、コード、wiki編集、問題、およびその他の貢献を削除、編集、または拒否する権利と責任を持ち、モデレーションの決定の理由を適切な場合に伝達します。

## Scope

この行動規範は、すべてのコミュニティスペース内で適用され、個人がパブリックスペースでコミュニティを公式に代表している場合にも適用されます。私たちのコミュニティを代表する例には、公式の電子メールアドレスの使用、公式のソーシャルメディアアカウントでの投稿、またはオンラインまたはオフラインイベントで指定された代表者として行動することが含まれます。

## Enforcement

虐待的、嫌がらせ、またはその他の受け入れられない行動の例は、実施を担当するコミュニティリーダーに{{.Email}}で報告できます。すべての苦情は迅速かつ公正にレビューおよび調査されます。

すべてのコミュニティリーダーは、インシデントの報告者のプライバシーとセキュリティを尊重する義務があります。

## Enforcement Guidelines

コミュニティリーダーは、この行動規範に違反すると判断した行動に対する結果を決定する際に、これらのコミュニティ影響ガイドラインに従います：

### 1. Correction

**コミュニティへの影響**: 不適切な言葉の使用、またはコミュニティで非専門的または歓迎されないと思われるその他の行動。

**結果**: コミュニティリーダーからの非公開の書面による警告で、違反の性質と行動が不適切であった理由の明確化。公的な謝罪が要求される場合があります。

### 2. Warning

**コミュニティへの影響**: 単一のインシデントまたは一連の行動による違反。

**結果**: 継続的な行動の結果を伴う警告。指定された期間、行動規範の実施者との一方的な対話を含む、関係者との相互作用なし。これには、コミュニティスペースおよびソーシャルメディアなどの外部チャネルでの相互作用の回避が含まれます。これらの条件に違反すると、一時的または永続的な禁止につながる可能性があります。

### 3. Temporary Ban

**コミュニティへの影響**: 持続的な不適切な行動を含む、コミュニティ基準の深刻な違反。

**結果**: 指定された期間、コミュニティとのあらゆる種類の相互作用または公的コミュニケーションの一時的な禁止。この期間中、行動規範の実施者との一方的な対話を含む、関係者との公的または私的な相互作用は許可されません。これらの条件に違反すると、永続的な禁止につながる可能性があります。

### 4. Permanent Ban

**コミュニティへの影響**: 持続的な不適切な行動、個人への嫌がらせ、または個人のクラスに対する攻撃性または軽蔑のパターンを示すこと。

**結果**: コミュニティ内でのあらゆる種類の公的相互作用の永続的な禁止。

## Attribution

この行動規範は[Contributor Covenant](https://www.contributor-covenant.org)、バージョン2.0から適応されており、https://www.contributor-covenant.org/version/2/0/code_of_conduct.htmlで入手できます。

コミュニティ影響ガイドラインは[Mozilla's code of conduct enforcement ladder](https://github.com/mozilla/diversity)にインスパイアされました。

この行動規範に関する一般的な質問の回答については、https://www.contributor-covenant.org/faqのFAQを参照してください。翻訳は、https://www.contributor-covenant.org/translationsで入手できます。
`

// Security policy template
const SecurityTemplate = `# Security Policy for {{.Name}}

We take the security of {{.Name}} seriously. This document outlines our security practices and how to report vulnerabilities.

## 📋 Supported Versions

Security updates are provided for the following versions of {{.Name}}:

| Version | Supported          | End of Life        |
| ------- | ------------------ | ------------------ |
| 2.x.x   | :white_check_mark: | TBD                |
| 1.x.x   | :white_check_mark: | 2025-12-31         |
| < 1.0   | :x:                | 2024-12-31         |

### Version Support Policy

- **Major versions**: Supported for at least 18 months after release
- **Minor versions**: Latest minor version receives security updates
- **Patch versions**: Critical security patches are backported to supported versions

## 🔒 Reporting a Vulnerability

**Please DO NOT report security vulnerabilities through public GitHub issues.**

### Preferred Reporting Methods

1. **GitHub Security Advisories** (Recommended)
   - Go to: https://github.com/{{.GitHubUsername}}/{{.Name}}/security/advisories/new
   - This ensures private disclosure and proper tracking

2. **Email**
   - Send to: {{.Email}}
   - Subject: [SECURITY] {{.Name}} - Brief description
   - Use GPG encryption if possible: [GPG Key](https://github.com/{{.GitHubUsername}}.gpg)

3. **Security Contact**
   - For critical vulnerabilities, contact us directly

### 📝 What to Include in Your Report

Please provide as much information as possible:

#### Required Information
- **Vulnerability Type**: What kind of vulnerability is this?
- **Location**: Where is the vulnerability located in the code?
- **Description**: Clear description of the issue
- **Impact**: What could an attacker potentially achieve?

#### Helpful Additional Information
- **Proof of Concept**: Steps to reproduce or PoC code
- **Affected Versions**: Which versions are affected?
- **Environment**: OS, Go version, deployment details
- **Mitigation**: Any temporary workarounds you've identified
- **CVSS Score**: If you've calculated one

#### Example Report Template

` + "```" + `
**Vulnerability Type**: [e.g., SQL Injection, XSS, Command Injection]
**Affected Component**: [e.g., cmd/server, internal/auth]
**Affected Versions**: [e.g., v1.0.0 - v1.2.3]
**Severity**: [e.g., High, Critical]

**Description**:
[Clear description of the vulnerability]

**Impact**:
[What can an attacker do with this vulnerability?]

**Steps to Reproduce**:
1. [First step]
2. [Second step]
3. [Continue...]

**Proof of Concept**:
[Include code, commands, or screenshots]

**Suggested Fix**:
[If you have ideas for how to fix it]
` + "```" + `

### 🕐 Response Timeline

We are committed to responding quickly to security issues:

| Timeline | Action |
|----------|--------|
| **< 24 hours** | Initial acknowledgment of your report |
| **< 72 hours** | Preliminary assessment and severity classification |
| **< 1 week** | Detailed analysis and remediation plan |
| **< 30 days** | Security patch released (for high/critical issues) |
| **< 90 days** | Public disclosure (coordinated with reporter) |

### 🎯 Severity Classification

We use the CVSS 3.1 scoring system with the following impact guidelines:

| Severity | CVSS Score | Response Time | Examples |
|----------|------------|---------------|----------|
| **Critical** | 9.0-10.0 | 24-48 hours | RCE, Auth bypass |
| **High** | 7.0-8.9 | 1 week | Privilege escalation |
| **Medium** | 4.0-6.9 | 2 weeks | Information disclosure |
| **Low** | 0.1-3.9 | 30 days | Minor information leaks |

## 🛡️ Security Measures

### Current Security Practices

- **Code Review**: All code changes require review
- **Automated Scanning**:
  - CodeQL for static analysis
  - Gosec for Go-specific vulnerabilities
  - Dependabot for dependency updates
- **CI/CD Security**:
  - All builds run in isolated environments
  - Secrets are properly managed
  - Dependencies are regularly updated

### For Users

#### Secure Configuration
- Always use the latest stable version
- Enable all recommended security features
- Regularly update dependencies
- Use strong authentication methods

#### Environment Security
- Run with minimal required permissions
- Use containers or sandboxing when possible
- Monitor logs for suspicious activity
- Implement proper network security

#### Deployment Best Practices
- Use HTTPS/TLS for all communications
- Implement proper input validation
- Use secure defaults
- Regular security audits

## 🎖️ Security Hall of Fame

We gratefully acknowledge security researchers who have helped improve {{.Name}}:

<!--
Format:
- **[Researcher Name](https://github.com/username)** - Vulnerability type (YYYY-MM-DD)
-->

### 2024
- *Be the first to contribute!*

### Recognition Policy

Security researchers who report valid vulnerabilities will be:

1. **Credited** in the security advisory (with permission)
2. **Listed** in this Hall of Fame
3. **Thanked** in release notes
4. **Invited** to test fixes before release

We do not currently offer monetary rewards but greatly appreciate contributions to our security.

## 🚫 Out of Scope

The following are generally considered out of scope:

- **Theoretical vulnerabilities** without proof of exploitability
- **Issues in third-party dependencies** (report to the respective maintainers)
- **Social engineering attacks**
- **Physical access attacks**
- **Denial of Service** through resource exhaustion
- **Issues affecting only outdated/unsupported versions**

## 🔍 Security Resources

### For Developers
- [Go Security Guidelines](https://golang.org/doc/security.html)
- [OWASP Go Secure Coding Practices](https://owasp.org/www-project-go-secure-coding-practices-guide/)
- [Our Contributing Guidelines](.github/CONTRIBUTING.md#security)

### Tools We Use
- [gosec](https://github.com/securecodewarrior/gosec) - Go AST scanner
- [CodeQL](https://codeql.github.com/) - Semantic code analysis
- [Dependabot](https://github.com/dependabot) - Dependency updates
- [nancy](https://github.com/sonatypecommunity/nancy) - Vulnerability scanner

## 📞 Contact Information

- **Security Email**: {{.Email}}
- **GPG Key**: [{{.GitHubUsername}}.gpg](https://github.com/{{.GitHubUsername}}.gpg)
- **Security Advisories**: https://github.com/{{.GitHubUsername}}/{{.Name}}/security/advisories
- **Project Maintainer**: [@{{.GitHubUsername}}](https://github.com/{{.GitHubUsername}})

---

**Last Updated**: 2024-01-01
**Policy Version**: 1.0

Thank you for helping keep {{.Name}} and the community safe! 🙏

<sub>This security policy is adapted from industry best practices and is regularly updated.</sub>
`

// Support template
const SupportTemplate = `# Support

{{.Name}}をお使いいただき、ありがとうございます！質問やサポートが必要な場合は、以下のリソースをご利用ください。

## 📚 Documentation

まず、以下のドキュメントをご確認ください：

- [README](../README.md) - プロジェクトの概要と基本的な使用方法
- [Contributing Guide](CONTRIBUTING.md) - 貢献方法
- [Examples](../examples/) - 使用例
- [Wiki](https://github.com/{{.GitHubUsername}}/{{.Name}}/wiki) - 詳細なドキュメント

## 🤔 Getting Help

### 💬 GitHub Discussions

質問や議論は[GitHub Discussions](https://github.com/{{.GitHubUsername}}/{{.Name}}/discussions)をご利用ください：

- **Q&A**: 使い方に関する質問
- **Ideas**: 新機能のアイデア
- **Show and Tell**: あなたの作品を共有
- **General**: その他の議論

### 🐛 Bug Reports

バグを発見した場合は[Issues](https://github.com/{{.GitHubUsername}}/{{.Name}}/issues)で報告してください：

1. [Bug Report Template](.github/ISSUE_TEMPLATE/bug_report.md)を使用
2. 再現手順と環境情報を含める
3. エラーメッセージやログを添付

### ✨ Feature Requests

新機能の提案は[Issues](https://github.com/{{.GitHubUsername}}/{{.Name}}/issues)で行ってください：

1. [Feature Request Template](.github/ISSUE_TEMPLATE/feature_request.md)を使用
2. 具体的な使用例を含める
3. 既存の代替案との比較を説明

## 📧 Direct Contact

緊急の問題やプライベートな質問がある場合：

- **Email**: {{.Email}}
- **Response Time**: 通常48時間以内

## 💡 Self-Help Resources

### 🔧 Troubleshooting

よくある問題と解決方法：

#### Installation Issues

` + "```bash" + `
# Go環境の確認
go version

# モジュールの再インストール
go clean -modcache
go mod download
` + "```" + `

#### Build Issues

` + "```bash" + `
# 依存関係の更新
go mod tidy

# クリーンビルド
make clean
make build
` + "```" + `

#### Runtime Issues

` + "```bash" + `
# 詳細ログの有効化
{{.Name}} --verbose

# 設定の確認
{{.Name}} config show
` + "```" + `

### 📖 Learning Resources

- [Go Documentation](https://golang.org/doc/)
- [Go by Example](https://gobyexample.com/)
- [Effective Go](https://golang.org/doc/effective_go.html)

## 🤝 Community

### 💬 Chat

- [Discord](https://discord.gg/your-discord) (if available)
- [Slack](https://your-slack.slack.com) (if available)

### 🐦 Social Media

- [Twitter](https://twitter.com/{{.GitHubUsername}}) (if available)
- [LinkedIn](https://linkedin.com/in/{{.GitHubUsername}}) (if available)

## 🏷️ Issue Labels

Issueを作成する際の参考として：

- ` + "`bug`" + `: バグ報告
- ` + "`enhancement`" + `: 新機能・改善
- ` + "`question`" + `: 質問
- ` + "`documentation`" + `: ドキュメント関連
- ` + "`good first issue`" + `: 初心者向け
- ` + "`help wanted`" + `: コミュニティからの支援募集

## 🙏 Thank You

{{.Name}}コミュニティの一員になっていただき、ありがとうございます！

あなたの質問、フィードバック、貢献がプロジェクトの改善に役立ちます。🚀

---

**注意**: セキュリティに関する問題は[Security Policy](SECURITY.md)に従って報告してください。
`

// GitHub Repository Settings template
const GitHubSettingsTemplate = `# GitHub Repository Settings for {{.Name}}

This document outlines the recommended settings for your {{.Name}} repository.

## Repository Settings

### General
- **Repository name**: {{.Name}}
- **Description**: {{.Description}}
- **Website**: https://{{.GitHubUsername}}.github.io/{{.Name}} (optional)
- **Topics**: Add relevant topics like ` + "`golang`" + `, ` + "`cli`" + `, ` + "`oss`" + `, ` + "`{{.Type}}`" + `

### Features
- ✅ Issues
- ✅ Projects
- ✅ Wiki (if documentation is extensive)
- ✅ Discussions (for community Q&A)
- ✅ Sponsorships (if accepting donations)

### Pull Requests
- ✅ Allow merge commits
- ✅ Allow squash merging (recommended)
- ✅ Allow rebase merging
- ✅ Always suggest updating pull request branches
- ✅ Automatically delete head branches

### Branch Protection Rules (main/master)

#### Required Settings:
- ✅ Require a pull request before merging
  - Required approvals: 1+ (adjust based on team size)
  - ✅ Dismiss stale PR approvals when new commits are pushed
  - ✅ Require review from code owners
- ✅ Require status checks to pass before merging
  - Required checks: ` + "`test`" + `, ` + "`lint`" + `, ` + "`security`" + `
- ✅ Require branches to be up to date before merging
- ✅ Require linear history (recommended for cleaner history)
- ❌ Allow force pushes (disabled for main branch security)
- ❌ Allow deletions (disabled for main branch security)

#### Additional Settings:
- ✅ Restrict pushes that create files larger than 100MB
- ✅ Include administrators (apply rules to admins too)

## Automated Setup Commands

Run these commands to configure your repository settings:

` + "```bash" + `
# Enable required features
gh repo edit {{.GitHubUsername}}/{{.Name}} --enable-issues=true
gh repo edit {{.GitHubUsername}}/{{.Name}} --enable-projects=true
gh repo edit {{.GitHubUsername}}/{{.Name}} --enable-wiki=true

# Set up branch protection
gh api repos/{{.GitHubUsername}}/{{.Name}}/branches/main/protection \\
  --method PUT \\
  --field required_status_checks='{"strict":true,"contexts":["test","lint","security"]}' \\
  --field enforce_admins=true \\
  --field required_pull_request_reviews='{"required_approving_review_count":1,"dismiss_stale_reviews":true,"require_code_owner_reviews":true}' \\
  --field restrictions=null

# Add repository topics
gh repo edit {{.GitHubUsername}}/{{.Name}} --add-topic golang,{{.Type}},oss,cli
` + "```" + `

## Environment Secrets

Set up the following secrets for CI/CD:

### Required Secrets:
- ` + "`GITHUB_TOKEN`" + `: Automatically provided by GitHub Actions
- ` + "`CODECOV_TOKEN`" + `: For code coverage reporting (if using Codecov)

### Optional Secrets (for releases):
- ` + "`DOCKER_USERNAME`" + `: Docker Hub username
- ` + "`DOCKER_PASSWORD`" + `: Docker Hub password/token
- ` + "`HOMEBREW_TOKEN`" + `: For Homebrew formula publishing

## Issue and PR Management

### Labels Setup
See [.github/labels.yml](.github/labels.yml) for comprehensive label configuration.

### Templates
- Bug Report: [.github/ISSUE_TEMPLATE/bug_report.md](.github/ISSUE_TEMPLATE/bug_report.md)
- Feature Request: [.github/ISSUE_TEMPLATE/feature_request.md](.github/ISSUE_TEMPLATE/feature_request.md)
- Question: [.github/ISSUE_TEMPLATE/question.md](.github/ISSUE_TEMPLATE/question.md)
- PR Template: [.github/PULL_REQUEST_TEMPLATE.md](.github/PULL_REQUEST_TEMPLATE.md)

### Automation
GitHub Actions workflows will:
- Auto-assign labels based on file changes
- Run tests on all PRs
- Auto-merge dependabot PRs (after tests pass)
- Create releases when tags are pushed
- Update changelog automatically

## Community Health

The repository includes:
- [x] Code of Conduct
- [x] Contributing Guidelines
- [x] Security Policy
- [x] Support Documentation
- [x] Issue Templates
- [x] PR Template
- [x] License

## Monitoring and Analytics

Consider enabling:
- GitHub Insights for contributor analytics
- Dependency graph and security alerts
- CodeQL code scanning
- Secret scanning
- GitHub Sponsors (if applicable)

---

This setup ensures your {{.Name}} project follows OSS best practices and provides a welcoming environment for contributors.
`

// Labels configuration template
const GitHubLabelsTemplate = `# GitHub Labels Configuration for {{.Name}}

# This file defines labels for better issue and PR management.
# Apply with: gh label list --repo {{.GitHubUsername}}/{{.Name}}

# Priority Labels
- name: "priority/critical"
  color: "d73a4a"
  description: "Critical priority, needs immediate attention"

- name: "priority/high"
  color: "ff6b6b"
  description: "High priority"

- name: "priority/medium"
  color: "ffab00"
  description: "Medium priority"

- name: "priority/low"
  color: "28a745"
  description: "Low priority"

# Type Labels
- name: "type/bug"
  color: "d73a4a"
  description: "Something isn't working"

- name: "type/feature"
  color: "0366d6"
  description: "New feature or request"

- name: "type/enhancement"
  color: "7057ff"
  description: "Improvement to existing feature"

- name: "type/documentation"
  color: "0075ca"
  description: "Improvements or additions to documentation"

- name: "type/question"
  color: "d876e3"
  description: "Further information is requested"

- name: "type/discussion"
  color: "c5def5"
  description: "Discussion or brainstorming"

# Status Labels
- name: "status/triage"
  color: "fbca04"
  description: "Needs triage and prioritization"

- name: "status/accepted"
  color: "28a745"
  description: "Accepted for development"

- name: "status/in-progress"
  color: "0052cc"
  description: "Currently being worked on"

- name: "status/blocked"
  color: "d73a4a"
  description: "Blocked by dependency or external factor"

- name: "status/needs-review"
  color: "fbca04"
  description: "Ready for review"

- name: "status/ready-to-merge"
  color: "28a745"
  description: "Approved and ready to merge"

# Effort Labels
- name: "effort/xs"
  color: "e6e6fa"
  description: "Extra small effort (< 1 hour)"

- name: "effort/s"
  color: "d8bfd8"
  description: "Small effort (1-4 hours)"

- name: "effort/m"
  color: "da70d6"
  description: "Medium effort (0.5-1 day)"

- name: "effort/l"
  color: "ba55d3"
  description: "Large effort (1-3 days)"

- name: "effort/xl"
  color: "9370db"
  description: "Extra large effort (3-5 days)"

- name: "effort/xxl"
  color: "663399"
  description: "Huge effort (1+ weeks)"

# Component Labels (adjust based on your project structure)
- name: "component/cli"
  color: "1d76db"
  description: "Command line interface"

- name: "component/core"
  color: "0e8a16"
  description: "Core functionality"

- name: "component/api"
  color: "0366d6"
  description: "API related"

- name: "component/docs"
  color: "0075ca"
  description: "Documentation"

- name: "component/tests"
  color: "ffd33d"
  description: "Testing related"

- name: "component/ci-cd"
  color: "5319e7"
  description: "CI/CD pipeline"

# Special Labels
- name: "good first issue"
  color: "7057ff"
  description: "Good for newcomers"

- name: "help wanted"
  color: "008672"
  description: "Extra attention is needed"

- name: "hacktoberfest"
  color: "ff6600"
  description: "Hacktoberfest eligible"

- name: "wontfix"
  color: "ffffff"
  description: "This will not be worked on"

- name: "invalid"
  color: "e4e669"
  description: "This doesn't seem right"

- name: "duplicate"
  color: "cfd3d7"
  description: "This issue or pull request already exists"

# Dependencies
- name: "dependencies"
  color: "0366d6"
  description: "Pull requests that update a dependency file"

- name: "security"
  color: "d73a4a"
  description: "Security related"

# Release Labels
- name: "release/patch"
  color: "28a745"
  description: "Patch release (bug fixes)"

- name: "release/minor"
  color: "0366d6"
  description: "Minor release (new features)"

- name: "release/major"
  color: "d73a4a"
  description: "Major release (breaking changes)"
`

// Question template
const QuestionTemplate = `---
name: ❓ Question
about: 使い方やプロジェクトについて質問する
title: '[QUESTION] '
labels: type/question
assignees: ''

---

## ❓ Question

質問内容を詳しくお書きください。

## 🔍 What I've Tried

試したことや調べたことがあれば教えてください：
- [ ] READMEを読みました
- [ ] ドキュメントを確認しました
- [ ] 既存のIssueを検索しました
- [ ] Discussionsを確認しました

## 💻 Environment

**実行環境:**
- OS: [e.g. macOS 12.0, Ubuntu 20.04, Windows 11]
- Go Version: [e.g. 1.21.0]
- {{.Name}} Version: [e.g. v1.0.0]

## 📄 Additional Context

その他の関連情報があればお書きください。

---

**💡 ヒント**:
- より迅速な回答を得るため、[Discussions](../../discussions) の利用もご検討ください
- コードに関する質問の場合は、最小限の再現可能な例を提供してください
`
