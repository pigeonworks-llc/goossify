# 🚀 goossify

**Next-generation boilerplate generator for Go OSS projects with complete automation**

[![Go Report Card](https://goreportcard.com/badge/github.com/pigeonworks-llc/goossify)](https://goreportcard.com/report/github.com/pigeonworks-llc/goossify)
[![MIT License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)
[![Release](https://img.shields.io/github/release/pigeonworks-llc/goossify.svg)](https://github.com/pigeonworks-llc/goossify/releases)
[![CI](https://github.com/pigeonworks-llc/goossify/workflows/CI/badge.svg)](https://github.com/pigeonworks-llc/goossify/actions)

> ⚠️ **Project in Development** - Currently under active development. Some features are planned for implementation. Check [Implementation Status](#-implementation-status) for details.

[🇯🇵 日本語版 README](README-ja.md)

## 🎯 Concept

goossify is a tool that **completely automates** the launch of Go OSS projects.
Beyond traditional boilerplate generation, it manages the entire project lifecycle.

## ✨ Features

### 🏗️ Fully Automated Project Initialization
- 📁 Optimal directory structure generation
- 📄 Auto-creation of essential files (README, LICENSE, .gitignore, etc.)
- 🔧 Development tool configuration (golangci-lint, GoReleaser, GitHub Actions, etc.)
- 📊 Quality management tool integration (coverage, benchmarks, etc.)

### 🤖 Continuous Maintenance Automation
- 🔄 Automatic dependency updates and vulnerability monitoring
- 📈 Complete release management automation (semantic versioning)
- 📝 Automatic changelog and documentation generation
- 👥 Community file management

### 🌟 Go Ecosystem Optimization
- 🐹 Full Go Modules support
- 📦 Automatic indexing on pkg.go.dev
- 🏆 Go language best practices application
- 🌍 International development environment consideration

## 🚀 Quick Start

### Installation

```bash
# Build development version (recommended)
git clone https://github.com/pigeonworks-llc/goossify.git
cd goossify
go build -o goossify

# Or install from releases (future)
go install github.com/pigeonworks-llc/goossify@latest
```

### Creating New OSS Projects

```bash
# 🚧 Planned: Interactive mode for new project creation
goossify init my-awesome-project

# 🚧 Planned: Create from templates
goossify create --template cli-tool my-cli-app
goossify create --template library my-go-lib
```

### Converting Existing Projects to OSS ✅

```bash
# Convert existing project to OSS-ready
cd existing-project
goossify ossify .

# Ossify current directory
goossify ossify

# This goossify project itself was dog-fooded
# Generated files: LICENSE, CONTRIBUTING.md, SECURITY.md, .github/workflows/ci.yml
```

## 📈 Recommended Workflow: Private → Public

### 1. **Start Development in Private Repository** 🔒

```bash
# Create private repository on GitHub
git clone git@github.com:yourusername/your-project.git
cd your-project

# Generate OSS preparation files
goossify ossify .
```

### 2. **Generate GitHub Personal Access Token** 🔑

Generate a Personal Access Token on GitHub:
- GitHub Settings → Developer settings → Personal access tokens → Generate new token
- Required permissions: `repo`, `read:org` (for branch protection checks)

### 3. **Gradual Quality Improvement** 📊

```bash
# Basic health check
goossify status .

# Complete check including GitHub settings
goossify status --github --github-token ghp_xxxxxxxxxxxx .

# Aim for 100/100 score
# If items are missing, fix with ossify again
goossify ossify .
```

### 4. **Public Release Readiness Check** 🚀

```bash
# Final confirmation of release readiness
goossify ready --github-token ghp_xxxxxxxxxxxx .

# Basic check is possible without token
goossify ready .
```

### 5. **Go Public** 🌍

- Change GitHub repository to **Public**
- Create initial release tag
- Verify automatic indexing on pkg.go.dev

### 💡 Benefits of This Workflow

- **🔒 Security**: Prevents accidental exposure of confidential information
- **📈 Quality Assurance**: Go public only when 100% ready
- **🔄 Iterative Improvement**: Safe trial and error in private environment
- **🤖 CI/CD Verification**: Pre-verify GitHub Actions work correctly

## 📁 Generated Project Structure

```
my-project/
├── 📄 README.md                 # Comprehensive project description
├── 📜 LICENSE                   # License file
├── 🐹 go.mod                    # Go Modules configuration
├── 🔧 .golangci.yml             # Linter configuration
├── 🚀 .goreleaser.yml           # Release automation configuration
├── 📋 .goossify.yml             # goossify management configuration
├── 🔄 .github/
│   ├── workflows/
│   │   ├── ci.yml               # Continuous Integration
│   │   ├── release.yml          # Automated Release
│   │   └── security.yml         # Security Scan
│   ├── ISSUE_TEMPLATE/
│   │   ├── bug_report.md        # Bug Report Template
│   │   └── feature_request.md   # Feature Request Template
│   └── PULL_REQUEST_TEMPLATE.md # Pull Request Template
├── 📚 docs/                     # Documentation
├── 🎯 examples/                 # Usage Examples
├── 🔧 internal/                 # Internal Packages
├── 📦 pkg/                      # Public Packages
├── 🚀 cmd/                      # Entry Points
└── 📋 renovate.json            # Dependency Update Automation
```

## 📊 Commands

### `goossify ossify` ✅
Convert existing projects to OSS-ready state

```bash
goossify ossify .                # Convert current directory
goossify ossify /path/to/project # Convert specified path
```

**Generated Files:**
- 📄 README.md (if not exists)
- 📜 LICENSE (MIT License)
- 📋 CONTRIBUTING.md (Contribution Guidelines)
- 🛡️ SECURITY.md (Security Policy)
- 🚫 .gitignore (Go-optimized)
- 🔧 .golangci.yml (Linter Configuration)
- 🚀 .goreleaser.yml (Release Automation)
- 🔄 .github/workflows/ (CI/CD Workflows)
- 📋 renovate.json (Dependency Management)

### `goossify status` ✅
Comprehensive project health analysis

```bash
goossify status .                           # Basic health check
goossify status --github --github-token TOKEN .  # Include GitHub settings
goossify status --format json .            # JSON output
```

**Analysis Categories:**
- 🏗️ **Basic Structure** (go.mod, README, .gitignore, etc.)
- 📚 **Documentation** (README, CONTRIBUTING, docs/, examples/)
- 🐙 **GitHub Integration** (CI/CD, issue templates, security policy)
- 🔧 **Quality Tools** (linter, tests, release automation)
- 📦 **Dependencies** (go.mod consistency, vulnerability check)
- 📜 **Licensing** (LICENSE file, go.mod license info)

### `goossify ready` ✅
Check if project is ready for public release

```bash
goossify ready .                           # Basic readiness check
goossify ready --github-token TOKEN .     # Include GitHub settings check
```

**Checks Performed:**
- ✅ OSS Health Score (must be 100/100)
- 🔍 Sensitive Information Detection
- 📜 License Consistency
- 🐙 GitHub Configuration (optional)

### `goossify create` 🚧
Create new projects from templates (Planned)

```bash
goossify create --template cli-tool my-app
goossify create --template library my-lib
```

## 🔧 Implementation Status

| Feature | Status | Description |
|---------|--------|-------------|
| 🏗️ `ossify` | ✅ **Complete** | Convert existing projects to OSS |
| 📊 `status` | ✅ **Complete** | Health analysis & scoring |
| 🚀 `ready` | ✅ **Complete** | Public release readiness check |
| 🐙 GitHub Integration | ✅ **Complete** | Branch protection, settings analysis |
| 🔧 `create` | 🚧 **Planned** | Template-based project creation |
| 🎯 `init` | 🚧 **Planned** | Interactive project initialization |
| 📈 `release` | 🚧 **Planned** | Automated release management |

## 🛠️ Development

### Prerequisites
- Go 1.21+
- Git
- GitHub CLI (optional, for enhanced GitHub integration)

### Building from Source

```bash
git clone https://github.com/pigeonworks-llc/goossify.git
cd goossify
go mod tidy
go build -o goossify .
```

### Running Tests

```bash
go test ./...
```

### Code Quality

```bash
go vet ./...
golangci-lint run
```

## 🤝 Contributing

We welcome contributions! Please see [CONTRIBUTING.md](CONTRIBUTING.md) for details.

### Quick Contributing Guide
1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🔗 Links

- [Documentation](https://github.com/pigeonworks-llc/goossify/tree/main/docs)
- [Examples](https://github.com/pigeonworks-llc/goossify/tree/main/examples)
- [Issue Tracker](https://github.com/pigeonworks-llc/goossify/issues)
- [Discussions](https://github.com/pigeonworks-llc/goossify/discussions)

## 🙏 Acknowledgments

- Go team for the excellent language and ecosystem
- All contributors who make this project better
- The Go OSS community for inspiration and best practices

---

**goossify** - Making Go OSS project creation and management effortless 🚀