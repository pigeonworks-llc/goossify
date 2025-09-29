package templates

// .golangci.yml template
const GolangCITemplate = `# golangci-lint configuration for {{.Name}}

run:
  timeout: 5m
  issues-exit-code: 1
  tests: true
  modules-download-mode: readonly
  go: "1.21"

output:
  formats:
    - format: colored-line-number
  print-issued-lines: true
  print-linter-name: true
  sort-results: true

issues:
  max-issues-per-linter: 0
  max-same-issues: 0
  exclude-use-default: false
  exclude-dirs:
    - vendor
    - third_party
    - testdata
    - examples
    - .git
  exclude-files:
    - ".*\\.pb\\.go$"
    - ".*_mock\\.go$"

linters-settings:
  errcheck:
    check-type-assertions: true
    check-blank: true

  gocyclo:
    min-complexity: 15

  funlen:
    lines: 100
    statements: 50

  gocognit:
    min-complexity: 20

  lll:
    line-length: 120

  revive:
    severity: warning
    rules:
      - name: exported
        severity: error
      - name: package-comments
        severity: warning

  goimports:
    local-prefixes: {{.ModulePath}}

  gosec:
    severity: medium
    confidence: medium

linters:
  enable:
    - errcheck
    - errorlint
    - gocyclo
    - gocognit
    - funlen
    - nestif
    - godox
    - gofmt
    - gofumpt
    - goimports
    - revive
    - stylecheck
    - unconvert
    - unparam
    - unused
    - varnamelen
    - whitespace
    - prealloc
    - bodyclose
    - rowserrcheck
    - gosec
    - copyloopvar
    - nilerr
    - nilnil
    - noctx
    - nolintlint
    - ineffassign
    - staticcheck
    - typecheck
    - goconst
    - gocritic
    - godot
    - misspell
    - predeclared
    - thelper
    - tparallel
    - wastedassign
    - testpackage
    - usetesting
    - imports
    - lll

  disable:
    - exhaustruct
    - forbidigo
    - forcetypeassert
    - cyclop
    - gci
    - ireturn
    - paralleltest
    - wsl

severity:
  default-severity: error
  case-sensitive: false
  rules:
    - linters:
        - revive
        - stylecheck
        - godot
      severity: warning
    - linters:
        - lll
        - varnamelen
        - godox
      severity: info
`

// .goreleaser.yml template
const GoReleaserTemplate = `# GoReleaser configuration for {{.Name}}

project_name: {{.Name}}

before:
  hooks:
    - go mod tidy

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
    ignore:
      - goos: windows
        goarch: arm64
    flags:
      - -trimpath
    ldflags:
      - "-s -w"

archives:
  - format: tar.gz
    format_overrides:
      - goos: windows
        format: zip

checksum:
  name_template: 'checksums.txt'

snapshot:
  name_template: "SNAPSHOT"

changelog:
  sort: asc
  filters:
    exclude:
      - "^test:"
      - "^chore"
      - "merge conflict"
      - Merge pull request
      - Merge remote-tracking branch
      - Merge branch
      - go mod tidy

release:
  draft: false
  prerelease: auto
`

// .goossify.yml template
const GoossifyConfigTemplate = `# goossify configuration for {{.Name}}

project:
  name: "{{.Name}}"
  description: "{{.Description}}"
  type: "{{.Type}}"
  license: "{{.License}}"
  author: "{{.Author}}"
  email: "{{.Email}}"
  github_username: "{{.GitHubUsername}}"

automation:
  release:
    strategy: "auto"              # auto, manual, scheduled
    versioning: "semantic"       # semantic, calendar
    changelog: true
    goreleaser: true
    auto_tag: true

  dependencies:
    auto_update: true
    security_check: true
    schedule: "weekly"           # daily, weekly, monthly
    exclude_list: []
    max_updates: 5

  quality:
    coverage_threshold: 80
    benchmarks: true
    performance_regression: 10   # % threshold
    static_analysis: true

  community:
    issue_templates: true
    pr_template: true
    contributing_guide: true
    code_of_conduct: true
    security_policy: true
    support_file: true

integrations:
  github:
    branch_protection: true
    required_reviews: 2
    status_checks: ["test", "lint", "security"]
    auto_merge: false
    delete_branch_on_merge: true

  ci_cd:
    provider: "github-actions"
    go_versions: ["1.21", "1.22", "1.23"]
    platforms: ["linux", "darwin", "windows"]
    coverage_reporting: true

monitoring:
  health_checks: true
  dependency_updates: "weekly"
  security_scans: "daily"
  performance_tracking: true

templates:
  update_policy: "auto"         # auto, manual, prompt
  custom_templates: []
  template_version: "latest"
`
