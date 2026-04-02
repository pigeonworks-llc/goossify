package ossify

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

// Ossifier はプロジェクトのOSS化を行う
type Ossifier struct {
	projectPath string
	projectName string
	interactive bool
	dryRun      bool
}

// New は新しいOssifierを作成
func New(projectPath string) *Ossifier {
	projectName := filepath.Base(projectPath)
	return &Ossifier{
		projectPath: projectPath,
		projectName: projectName,
	}
}

// SetInteractive は対話モードを設定
func (o *Ossifier) SetInteractive(interactive bool) {
	o.interactive = interactive
}

// SetDryRun はdry-runモードを設定
func (o *Ossifier) SetDryRun(dryRun bool) {
	o.dryRun = dryRun
}

// Execute はOSS化処理を実行
func (o *Ossifier) Execute() error {
	fmt.Println("📋 現在の状況を分析中...")

	// 1. 必須ファイルを生成
	if err := o.ensureLicense(); err != nil {
		return err
	}

	if err := o.ensureGitHubFiles(); err != nil {
		return err
	}

	if err := o.ensureIssueTemplates(); err != nil {
		return err
	}

	if err := o.ensureCommunityFiles(); err != nil {
		return err
	}

	if err := o.ensureEnhancedCommunityFiles(); err != nil {
		return err
	}

	// 2. Git初期化
	if err := o.initGitIfNeeded(); err != nil {
		return err
	}

	// 3. 基本ディレクトリ生成
	if err := o.ensureDirectories(); err != nil {
		return err
	}

	// 4. 設定ファイル生成
	if err := o.ensureConfigFiles(); err != nil {
		return err
	}

	// 5. 依存関係整理
	if err := o.tidyDependencies(); err != nil {
		return err
	}

	return nil
}

func (o *Ossifier) ensureLicense() error {
	licensePath := filepath.Join(o.projectPath, "LICENSE")
	if _, err := os.Stat(licensePath); err == nil {
		fmt.Println("📄 LICENSE ファイルは既に存在します")
		return nil
	}

	fmt.Println("📄 LICENSE ファイルを生成中...")
	licenseContent := fmt.Sprintf(`Apache License
Version 2.0, January 2004
http://www.apache.org/licenses/

Copyright %d Pigeon Works LLC

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.`, time.Now().Year())

	return os.WriteFile(licensePath, []byte(licenseContent), 0600)
}

func (o *Ossifier) ensureGitHubFiles() error {
	githubDir := filepath.Join(o.projectPath, ".github")
	workflowsDir := filepath.Join(githubDir, "workflows")

	// .github/workflows ディレクトリ作成
	if err := os.MkdirAll(workflowsDir, 0o750); err != nil {
		return fmt.Errorf("gitHub ディレクトリ作成失敗: %w", err)
	}

	// CI workflow生成
	ciPath := filepath.Join(workflowsDir, "ci.yml")
	if _, err := os.Stat(ciPath); err != nil {
		fmt.Println("🔧 CI/CD workflow を生成中...")
		ciContent := `name: CI

on:
  push:
    branches: [ main, develop ]
  pull_request:
    branches: [ main ]

jobs:
  test:
    runs-on: ubuntu-latest
    strategy:
      matrix:
        go-version: [1.21, 1.22, 1.23]

    steps:
    - uses: actions/checkout@v4

    - name: Set up Go
      uses: actions/setup-go@v4
      with:
        go-version: ${{ matrix.go-version }}

    - name: Verify dependencies
      run: go mod verify

    - name: Build
      run: go build -v ./...

    - name: Run tests
      run: go test -race -coverprofile=coverage.out -covermode=atomic ./...

    - name: Run vet
      run: go vet ./...

  lint:
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@v4
    - uses: actions/setup-go@v4
      with:
        go-version: 1.23
    - name: golangci-lint
      uses: golangci/golangci-lint-action@v3
      with:
        version: latest`

		if err := os.WriteFile(ciPath, []byte(ciContent), 0600); err != nil {
			return fmt.Errorf("ci設定ファイル作成失敗: %w", err)
		}
	}

	// Release workflow生成
	releasePath := filepath.Join(workflowsDir, "release.yml")
	if _, err := os.Stat(releasePath); err != nil {
		fmt.Println("🚀 Release workflow を生成中...")
		releaseContent := `name: Release

on:
  push:
    tags:
      - 'v*'

permissions:
  contents: write

jobs:
  release:
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@v4
      with:
        fetch-depth: 0

    - name: Set up Go
      uses: actions/setup-go@v4
      with:
        go-version: stable

    - name: Run GoReleaser
      uses: goreleaser/goreleaser-action@v5
      with:
        distribution: goreleaser
        version: latest
        args: release --clean
      env:
        GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
        GITHUB_OWNER: ${{ github.repository_owner }}`

		if err := os.WriteFile(releasePath, []byte(releaseContent), 0600); err != nil {
			return fmt.Errorf("release workflow作成失敗: %w", err)
		}
	}

	return nil
}

func (o *Ossifier) ensureCommunityFiles() error {
	// CONTRIBUTING.md
	contributingPath := filepath.Join(o.projectPath, "CONTRIBUTING.md")
	if _, err := os.Stat(contributingPath); err != nil {
		fmt.Println("📖 CONTRIBUTING.md を生成中...")
		contributingContent := fmt.Sprintf(`# Contributing to %s

We love your input! We want to make contributing to this project as easy and transparent as possible.

## Development Process

1. Fork the repo and create your branch from main.
2. If you've added code that should be tested, add tests.
3. If you've changed APIs, update the documentation.
4. Ensure the test suite passes.
5. Make sure your code lints.
6. Issue that pull request!

## Getting Started

### Prerequisites

- Go 1.21 or later
- Git

### Setup

`+"```bash"+`
git clone https://github.com/pigeonworks-llc/%s.git
cd %s
go mod download
`+"```"+`

### Running Tests

`+"```bash"+`
go test ./...
`+"```"+`

### Running the Application

`+"```bash"+`
go run main.go
`+"```"+`

## Pull Request Process

1. Update the README.md with details of changes to the interface.
2. Increase the version numbers in any examples files and the README.md to the new version that this Pull Request would represent.
3. You may merge the Pull Request in once you have the sign-off of two other developers.

## Code of Conduct

By participating, you are expected to uphold this code. Please report unacceptable behavior to [support@pigeonworks.llc](mailto:support@pigeonworks.llc).

## License

By contributing, you agree that your contributions will be licensed under the Apache License 2.0.`, o.projectName, o.projectName, o.projectName)

		if err := os.WriteFile(contributingPath, []byte(contributingContent), 0600); err != nil {
			return fmt.Errorf("cONTRIBUTING.md作成失敗: %w", err)
		}
	}

	// SECURITY.md（ルートディレクトリに生成）
	securityPath := filepath.Join(o.projectPath, "SECURITY.md")
	if _, err := os.Stat(securityPath); err != nil {
		fmt.Println("🔒 SECURITY.md を生成中...")
		securityContent := fmt.Sprintf(`# Security Policy

## Supported Versions

| Version | Supported          |
| ------- | ------------------ |
| 1.0.x   | :white_check_mark: |
| < 1.0   | :x:                |

## Reporting a Vulnerability

We take the security of %s seriously. If you believe you have found a security vulnerability, please report it to us as described below.

**Please do not report security vulnerabilities through public GitHub issues.**

Instead, please report them via email to [security@pigeonworks.llc](mailto:security@pigeonworks.llc).

You should receive a response within 48 hours. If for some reason you do not, please follow up via email to ensure we received your original message.

Please include the requested information listed below (as much as you can provide) to help us better understand the nature and scope of the possible issue:

* Type of issue (e.g. buffer overflow, SQL injection, cross-site scripting, etc.)
* Full paths of source file(s) related to the manifestation of the issue
* The location of the affected source code (tag/branch/commit or direct URL)
* Any special configuration required to reproduce the issue
* Step-by-step instructions to reproduce the issue
* Proof-of-concept or exploit code (if possible)
* Impact of the issue, including how an attacker might exploit the issue

This information will help us triage your report more quickly.

## Preferred Languages

We prefer all communications to be in English or Japanese.`, o.projectName)

		if err := os.WriteFile(securityPath, []byte(securityContent), 0600); err != nil {
			return fmt.Errorf("sECURITY.md作成失敗: %w", err)
		}
	}

	return nil
}

func (o *Ossifier) initGitIfNeeded() error {
	gitDir := filepath.Join(o.projectPath, ".git")
	if _, err := os.Stat(gitDir); err == nil {
		fmt.Println("📦 Git リポジトリは既に初期化済みです")
		return nil
	}

	fmt.Println("📦 Git リポジトリを初期化中...")
	cmd := exec.Command("git", "init")
	cmd.Dir = o.projectPath
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("git初期化失敗: %w", err)
	}

	return nil
}

func (o *Ossifier) tidyDependencies() error {
	fmt.Println("📦 依存関係を整理中...")
	cmd := exec.Command("go", "mod", "tidy")
	cmd.Dir = o.projectPath
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("go mod tidy失敗: %w", err)
	}

	return nil
}

func (o *Ossifier) ensureIssueTemplates() error {
	issueTemplateDir := filepath.Join(o.projectPath, ".github", "ISSUE_TEMPLATE")

	// ISSUE_TEMPLATE ディレクトリ作成
	if err := os.MkdirAll(issueTemplateDir, 0o750); err != nil {
		return fmt.Errorf("issue template ディレクトリ作成失敗: %w", err)
	}

	// Bug report template
	bugReportPath := filepath.Join(issueTemplateDir, "bug_report.md")
	if _, err := os.Stat(bugReportPath); err != nil {
		fmt.Println("🐛 Bug report テンプレートを生成中...")
		bugReportContent := `---
name: Bug report
about: Create a report to help us improve
title: '[BUG] '
labels: 'bug'
assignees: ''
---

**Describe the bug**
A clear and concise description of what the bug is.

**To Reproduce**
Steps to reproduce the behavior:
1. Go to '...'
2. Click on '....'
3. Scroll down to '....'
4. See error

**Expected behavior**
A clear and concise description of what you expected to happen.

**Screenshots**
If applicable, add screenshots to help explain your problem.

**Environment (please complete the following information):**
 - OS: [e.g. macOS, Linux, Windows]
 - Go version: [e.g. 1.21]
 - Version: [e.g. v1.0.0]

**Additional context**
Add any other context about the problem here.`

		if err := os.WriteFile(bugReportPath, []byte(bugReportContent), 0600); err != nil {
			return fmt.Errorf("bug report テンプレート作成失敗: %w", err)
		}
	}

	// Feature request template
	featureReqPath := filepath.Join(issueTemplateDir, "feature_request.md")
	if _, err := os.Stat(featureReqPath); err != nil {
		fmt.Println("✨ Feature request テンプレートを生成中...")
		featureReqContent := `---
name: Feature request
about: Suggest an idea for this project
title: '[FEATURE] '
labels: 'enhancement'
assignees: ''
---

**Is your feature request related to a problem? Please describe.**
A clear and concise description of what the problem is. Ex. I'm always frustrated when [...]

**Describe the solution you'd like**
A clear and concise description of what you want to happen.

**Describe alternatives you've considered**
A clear and concise description of any alternative solutions or features you've considered.

**Additional context**
Add any other context or screenshots about the feature request here.`

		if err := os.WriteFile(featureReqPath, []byte(featureReqContent), 0600); err != nil {
			return fmt.Errorf("feature request テンプレート作成失敗: %w", err)
		}
	}

	// Pull Request template
	prTemplatePath := filepath.Join(o.projectPath, ".github", "PULL_REQUEST_TEMPLATE.md")
	if _, err := os.Stat(prTemplatePath); err != nil {
		fmt.Println("📝 Pull Request テンプレートを生成中...")
		prTemplateContent := `## Description
Brief description of the changes made in this PR.

## Type of Change
- [ ] Bug fix (non-breaking change which fixes an issue)
- [ ] New feature (non-breaking change which adds functionality)
- [ ] Breaking change (fix or feature that would cause existing functionality to not work as expected)
- [ ] Documentation update
- [ ] Code refactoring
- [ ] Other (please describe):

## How Has This Been Tested?
Please describe the tests that you ran to verify your changes.

- [ ] Unit tests
- [ ] Integration tests
- [ ] Manual testing

## Checklist
- [ ] My code follows the style guidelines of this project
- [ ] I have performed a self-review of my own code
- [ ] I have commented my code, particularly in hard-to-understand areas
- [ ] I have made corresponding changes to the documentation
- [ ] My changes generate no new warnings
- [ ] I have added tests that prove my fix is effective or that my feature works
- [ ] New and existing unit tests pass locally with my changes

## Screenshots (if appropriate)
Add screenshots to help explain your changes.

## Additional Notes
Any additional information, configuration, or data that might be helpful in reviewing this PR.`

		if err := os.WriteFile(prTemplatePath, []byte(prTemplateContent), 0600); err != nil {
			return fmt.Errorf("pR テンプレート作成失敗: %w", err)
		}
	}

	return nil
}

func (o *Ossifier) ensureDirectories() error {
	// docsディレクトリ作成
	docsDir := filepath.Join(o.projectPath, "docs")
	if _, err := os.Stat(docsDir); err != nil {
		fmt.Println("📚 docs ディレクトリを生成中...")
		if err := os.MkdirAll(docsDir, 0o750); err != nil {
			return fmt.Errorf("docs ディレクトリ作成失敗: %w", err)
		}

		// .gitkeepファイルを作成してディレクトリを保持
		gitkeepPath := filepath.Join(docsDir, ".gitkeep")
		gitkeepContent := "# このファイルはdocsディレクトリを保持するためのものです\n# ドキュメントファイルを追加したら削除してください\n"
		if err := os.WriteFile(gitkeepPath, []byte(gitkeepContent), 0600); err != nil {
			return fmt.Errorf(".gitkeep作成失敗: %w", err)
		}
	}

	return nil
}

func (o *Ossifier) ensureConfigFiles() error {
	// Renovate設定
	renovatePath := filepath.Join(o.projectPath, "renovate.json")
	if _, err := os.Stat(renovatePath); err != nil {
		fmt.Println("🔄 Renovate設定ファイルを生成中...")
		renovateContent := `{
  "$schema": "https://docs.renovatebot.com/renovate-schema.json",
  "extends": [
    "config:base",
    ":semanticCommitTypeAll(deps)"
  ],
  "schedule": ["before 9am on monday"],
  "timezone": "Asia/Tokyo",
  "labels": ["dependencies"],
  "commitMessagePrefix": "deps:",
  "prTitle": "deps: {{#if isSingleVersion}}{{depName}}{{else}}{{depName}} packages{{/if}}",
  "packageRules": [
    {
      "description": "Golang dependencies",
      "matchFileNames": ["go.mod", "go.sum"],
      "matchDepTypes": ["require"],
      "semanticCommitType": "deps",
      "commitMessageTopic": "Go module {{depName}}",
      "groupName": "Go dependencies"
    },
    {
      "description": "GitHub Actions",
      "matchFileNames": [".github/workflows/**"],
      "matchDepTypes": ["action"],
      "semanticCommitType": "ci",
      "commitMessageTopic": "GitHub Action {{depName}}",
      "groupName": "GitHub Actions"
    }
  ],
  "golang": {
    "minimumReleaseAge": "3 days"
  },
  "lockFileMaintenance": {
    "enabled": true,
    "schedule": ["before 9am on the first day of the month"]
  },
  "vulnerabilityAlerts": {
    "enabled": true,
    "schedule": ["at any time"],
    "dependencyDashboardApproval": false
  },
  "dependencyDashboard": true,
  "automerge": false,
  "rebaseWhen": "conflicted"
}`

		if err := os.WriteFile(renovatePath, []byte(renovateContent), 0600); err != nil {
			return fmt.Errorf("renovate設定ファイル作成失敗: %w", err)
		}
	}

	// GoReleaser設定
	goreleaserPath := filepath.Join(o.projectPath, ".goreleaser.yml")
	if _, err := os.Stat(goreleaserPath); err != nil {
		fmt.Println("🚀 GoReleaser設定ファイルを生成中...")
		goreleaserContent := `# GoReleaser configuration
# https://goreleaser.com/customization/

project_name: ` + o.projectName + `

before:
  hooks:
    - go mod tidy
    - go generate ./...

builds:
  - env:
      - CGO_ENABLED=0
    goos:
      - linux
      - windows
      - darwin
    goarch:
      - amd64
      - arm64
    binary: ` + o.projectName + `
    main: ./main.go

archives:
  - format: tar.gz
    name_template: >-
      {{ .ProjectName }}_
      {{- title .Os }}_
      {{- if eq .Arch "amd64" }}x86_64
      {{- else if eq .Arch "386" }}i386
      {{- else }}{{ .Arch }}{{ end }}
      {{- if .Arm }}v{{ .Arm }}{{ end }}
    format_overrides:
      - goos: windows
        format: zip

checksum:
  name_template: 'checksums.txt'

release:
  github:
    owner: ` + "{{ .Env.GITHUB_OWNER }}" + `
    name: ` + o.projectName + `
  name_template: "Release {{ .Tag }}"
  draft: false
  prerelease: auto

changelog:
  sort: asc
  use: github
  filters:
    exclude:
      - '^docs:'
      - '^test:'
      - '^ci:'
      - '^build:'
      - Merge pull request
      - Merge branch
  groups:
    - title: Features
      regexp: "^.*feat[(\\w)]*:+.*$"
      order: 0
    - title: 'Bug fixes'
      regexp: "^.*fix[(\\w)]*:+.*$"
      order: 1
    - title: Others
      order: 999`

		if err := os.WriteFile(goreleaserPath, []byte(goreleaserContent), 0600); err != nil {
			return fmt.Errorf("goReleaser設定ファイル作成失敗: %w", err)
		}
	}

	// golangci-lint設定
	golangciPath := filepath.Join(o.projectPath, ".golangci.yml")
	if _, err := os.Stat(golangciPath); err != nil {
		fmt.Println("🔍 golangci-lint設定ファイルを生成中...")
		golangciContent := `# golangci-lint configuration
# https://golangci-lint.run/usage/configuration/

run:
  timeout: 5m
  modules-download-mode: readonly

linters:
  enable:
    - errcheck      # checking for unchecked errors
    - gosimple      # simplify code
    - govet         # reports suspicious constructs
    - ineffassign   # detects unused assignments
    - staticcheck   # go vet on steroids
    - typecheck     # type-checks Go code
    - unused        # checks for unused constants, variables, functions and types
    - goimports     # fix imports and format code
    - misspell      # finds misspelled English words in comments
    - gocritic      # provides diagnostics that check for bugs, performance and style issues
    - revive        # fast, configurable, extensible, flexible replacement for golint
    - gosec         # inspects source code for security problems
    - unconvert     # remove unnecessary type conversions
    - unparam       # reports unused function parameters
    - gocyclo       # computes cyclomatic complexities
    - gofmt         # checks whether code was gofmt-ed
    - whitespace    # tool for detection of leading and trailing whitespace

linters-settings:
  gocyclo:
    min-complexity: 15

  gocritic:
    enabled-tags:
      - diagnostic
      - experimental
      - opinionated
      - performance
      - style
    disabled-checks:
      - dupImport
      - ifElseChain
      - octalLiteral
      - whyNoLint
      - wrapperFunc

  revive:
    rules:
      - name: exported
        arguments: [checkPrivateReceivers]

  gosec:
    excludes:
      - G204 # subprocess launched with variable
      - G304 # file path provided as tainted input

issues:
  exclude-rules:
    # Exclude some linters from running on tests files
    - path: _test\.go
      linters:
        - gocyclo
        - gosec
        - errcheck
        - revive

    # Exclude known linters from partially hard-to-fix issues
    - path: internal/.*
      text: "unexported-return"
      linters:
        - revive

    # Exclude shadow checking on err variables in some cases
    - text: "shadow: declaration of \"err\""
      linters:
        - govet

  max-issues-per-linter: 50
  max-same-issues: 3`

		if err := os.WriteFile(golangciPath, []byte(golangciContent), 0600); err != nil {
			return fmt.Errorf("golangci-lint設定ファイル作成失敗: %w", err)
		}
	}

	return nil
}

func (o *Ossifier) ensureEnhancedCommunityFiles() error {
	// SUPPORT.md
	supportPath := filepath.Join(o.projectPath, "SUPPORT.md")
	supportContent := fmt.Sprintf(`# Support

%sをお使いいただき、ありがとうございます！質問やサポートが必要な場合は、以下のリソースをご利用ください。

## 📚 Documentation

まず、以下のドキュメントをご確認ください：

- [README](README.md) - プロジェクトの概要と基本的な使用方法
- [Contributing Guide](CONTRIBUTING.md) - 貢献方法
- [Examples](examples/) - 使用例

## 🤔 Getting Help

### 💬 GitHub Discussions

質問や議論は[GitHub Discussions](https://github.com/pigeonworks-llc/%s/discussions)をご利用ください：

- **Q&A**: 使い方に関する質問
- **Ideas**: 新機能のアイデア
- **Show and Tell**: あなたの作品を共有
- **General**: その他の議論

### 🐛 Bug Reports

バグを発見した場合は[Issues](https://github.com/pigeonworks-llc/%s/issues)で報告してください：

1. Bug Report Templateを使用
2. 再現手順と環境情報を含める
3. エラーメッセージやログを添付

### ✨ Feature Requests

新機能の提案は[Issues](https://github.com/pigeonworks-llc/%s/issues)で行ってください：

1. Feature Request Templateを使用
2. 具体的な使用例を含める
3. 既存の代替案との比較を説明

## 📧 Direct Contact

緊急の問題やプライベートな質問がある場合：

- **Email**: support@pigeonworks.llc
- **Response Time**: 通常48時間以内

## 💡 Self-Help Resources

### 🔧 Troubleshooting

よくある問題と解決方法については、ドキュメントをご覧ください。

### 📖 Learning Resources

- [Go Documentation](https://golang.org/doc/)
- [Go by Example](https://gobyexample.com/)
- [Effective Go](https://golang.org/doc/effective_go.html)

## 🙏 Thank You

%sコミュニティの一員になっていただき、ありがとうございます！

あなたの質問、フィードバック、貢献がプロジェクトの改善に役立ちます。🚀

---

**注意**: セキュリティに関する問題は[Security Policy](SECURITY.md)に従って報告してください。`, o.projectName, o.projectName, o.projectName, o.projectName, o.projectName)

	if _, err := os.Stat(supportPath); err != nil {
		fmt.Println("🆘 SUPPORT.md を生成中...")
		if err := o.writeFileInteractive(supportPath, supportContent, "SUPPORT.md - サポート情報"); err != nil {
			return fmt.Errorf("sUPPORT.md作成失敗: %w", err)
		}
	}

	return nil
}

// writeFileInteractive はファイルを対話的に書き込む
func (o *Ossifier) writeFileInteractive(filePath, content, description string) error {
	// 既存ファイル確認
	exists := false
	if _, err := os.ReadFile(filePath); err == nil {
		exists = true
	}

	// 対話モードまたはdry-runモードの処理
	if o.interactive || o.dryRun {
		fmt.Printf("\n📄 %s\n", description)
		fmt.Printf("   ファイル: %s\n", filePath)

		if exists {
			fmt.Println("   状態: 既存ファイル（変更なし）")
			if o.interactive {
				fmt.Print("\n   スキップしますか? [Y/n]: ")
				if !o.askConfirmation(true) {
					return nil // ユーザーがスキップを選択
				}
			}
			return nil // 既存ファイルは変更しない
		} else {
			fmt.Println("   状態: 新規作成")
			fmt.Println("\n--- 内容プレビュー ---")
			lines := strings.Split(content, "\n")
			previewLines := 15
			for i, line := range lines {
				if i >= previewLines {
					fmt.Printf("... (%d行省略)\n", len(lines)-previewLines)
					break
				}
				fmt.Printf("   %s\n", line)
			}
			fmt.Println("--- プレビュー終了 ---")

			if o.dryRun {
				fmt.Println("   [DRY-RUN] 実際には作成しません")
				return nil
			}

			if o.interactive {
				fmt.Print("\n   このファイルを作成しますか? [Y/n]: ")
				if !o.askConfirmation(true) {
					fmt.Println("   ⏭️  スキップしました")
					return nil
				}
			}
		}
	}

	// ファイル書き込み
	if !o.dryRun {
		if err := os.WriteFile(filePath, []byte(content), 0600); err != nil {
			return err
		}
	}

	return nil
}

// askConfirmation はユーザーに確認を求める
func (o *Ossifier) askConfirmation(defaultYes bool) bool {
	reader := bufio.NewReader(os.Stdin)
	response, err := reader.ReadString('\n')
	if err != nil {
		return defaultYes
	}

	response = strings.ToLower(strings.TrimSpace(response))
	if response == "" {
		return defaultYes
	}

	return response == "y" || response == "yes"
}
