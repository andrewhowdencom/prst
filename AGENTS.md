# Agent Instructions: prst

## Project

- **Name**: prst
- **Module**: `github.com/andrewhowdencom/prst`
- **Purpose**: CLI tool for managing Bash prompt strings (PS0–PS4)

## Build & Development

- **Task Runner**: [Task](https://taskfile.dev/) (`Taskfile.yml`)
- **Standard Tasks**: `setup`, `generate`, `build`, `test`, `lint`, `validate`
- **Go Version**: 1.26.2 (use whatever is installed on the system)
- **Testing**: Always run with `go test -race ./...`
- **Linting**: `golangci-lint` with `errcheck`, `govet`, `staticcheck`, `gosec`
- **Formatting**: Strict `gofmt`

## Architecture

- **cmd/prst/main.go**: Minimal entry point — delegates to Wire-generated DI graph.
- **internal/app/**: Cobra root command and global flags (`--log-level`).
- **internal/commands/**: Subcommands (`0`, `1`, `2`, `3`, `4`, `version`).
- **internal/configuration/**: Viper + `adrg/xdg` configuration infrastructure.
- **internal/di/**: Google Wire provider set and generated `wire_gen.go`.

## Conventions

- Use `log/slog` for logging (configured via `--log-level`).
- Configuration keys use nested objects (not underscores), e.g. `style.name`.
- Environment variable prefix: `PRST_`.
- Config file location: `$XDG_CONFIG_HOME/prst/config.yaml`.
- Use the functional options pattern for object construction.
- Table-driven tests for unit tests.
- Wrap errors with `fmt.Errorf("...: %w", err)`.
