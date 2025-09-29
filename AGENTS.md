# Repository Guidelines

## Project Structure & Module Organization
The CLI entrypoint lives in `main.go`, delegating to Cobra commands under `cmd/` (init, create, status, release). Reusable business logic sits inside `internal/` (configuration, generator, template helpers) while exportable APIs belong in `pkg/`. Scaffolding blueprints and CI assets are defined in `templates/`, and runnable showcase projects are kept in `examples/`. Generated binaries such as `goossify` should not be committed; treat them as build artifacts.

## Build, Test, and Development Commands
Use `go run . --help` during local development to exercise the CLI. `go build ./...` produces a platform-native binary, while `go install github.com/pigeonworks-llc/goossify@latest` validates module metadata matches remote expectations. Run `go test ./...` for quick feedback and `go test -race ./...` before opening a pull request. Lint with `golangci-lint run` once the optional `.golangci.yml` is present, and refresh dependencies via `go mod tidy` when templates add new imports.

## Coding Style & Naming Conventions
All Go code must pass `gofmt`; use tabs for indentation and keep line length within Go defaults. Adopt `goimports` to maintain import grouping. Exported identifiers follow PascalCase and include doc comments that start with the identifier name; internal helpers stay lowerCamelCase. Cobra command files should expose a single `*cobra.Command` and register in `init()`, mirroring the existing command layout.

## Testing Guidelines
Place table-driven `_test.go` files beside the code they verify (e.g., `internal/generator/generator_test.go`). Aim for the 80% coverage target enforced by generated templates, covering edge cases like missing config files or invalid template inputs. Prefer `testing` + `testify/require` style assertions if dependencies are added, and exercise CLI surfaces through small integration tests invoking `cmd.Execute()` with stubbed I/O. Always run `go test -race ./...` to catch data races in concurrent scaffolding routines.

## Commit & Pull Request Guidelines
Follow Conventional Commits (`type(scope): summary`), keeping the scope aligned with package or command names (`feat(cmd/init):`, `fix(internal/config):`). Each pull request should explain the motivation, reference a GitHub issue when available, and include CLI output or screenshots for user-facing changes. Before requesting review, ensure linting and tests pass, regressions are covered by new cases, and documentation or templates remain synchronized.

## Security & Configuration Tips
Keep `.goossify.yml` free of secrets; use environment variables or local overrides instead. When adding templates, double-check that generated workflows pin third-party actions to specific versions. Avoid committing developer-specific paths or credentials inside sample configs, and review new commands for file-system writes or network calls that require user confirmation.
