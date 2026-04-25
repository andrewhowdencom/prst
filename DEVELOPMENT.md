# Development

This document explains how to set up a development environment, run tests, and contribute to `prst`.

## Prerequisites

- [Go](https://go.dev/) (the project targets Go 1.26.2, but works with any recent version).
- [Task](https://taskfile.dev/) (task runner).
- [golangci-lint](https://golangci-lint.run/) (optional — installed automatically by `task setup`).

## Quick Start

```bash
git clone https://github.com/andrewhowdencom/prst.git
cd prst
task setup   # Install tools and dependencies
```

## Common Tasks

| Task | Command | Description |
|---|---|---|
| Generate | `task generate` | Run code generation (Wire DI graph). |
| Build | `task build` | Compile the `prst` binary to `./prst`. |
| Test | `task test` | Run the full test suite with the race detector (`go test -race ./...`). |
| Lint | `task lint` | Run `golangci-lint` with `errcheck`, `govet`, `staticcheck`, and `gosec`. |
| Validate | `task validate` | Run lint, test, and build in one go. Run this before committing. |

## Code Generation

`prst` uses [Google Wire](https://github.com/google/wire) for compile-time dependency injection.

- The provider set is defined in `internal/di/wire.go`.
- The generated graph lives in `internal/di/wire_gen.go`.

**Never edit `wire_gen.go` by hand.** After changing constructors or adding new providers, run:

```bash
wire ./internal/di
```

Or use the Task target:

```bash
task generate
```

## Testing Conventions

- **Always** run tests with the race detector: `go test -race ./...`
- Use **table-driven tests** for unit tests.
- Name test files `*_test.go` in the same package as the code under test.
- Assert expected behavior, not implementation details.

Example pattern:

```go
func TestFoo(t *testing.T) {
    cases := []struct {
        name     string
        input    string
        expected string
    }{
        {"basic", "hello", "HELLO"},
        {"empty", "", ""},
    }

    for _, tc := range cases {
        t.Run(tc.name, func(t *testing.T) {
            got := Foo(tc.input)
            if got != tc.expected {
                t.Errorf("Foo(%q) = %q; want %q", tc.input, got, tc.expected)
            }
        })
    }
}
```

## Coding Style

- **Formatting**: `gofmt` — strict, no exceptions.
- **Linting**: `golangci-lint` with the following enabled linters:
  - `errcheck`
  - `govet`
  - `staticcheck`
  - `gosec`
- **Error wrapping**: always wrap errors with context using `fmt.Errorf("...: %w", err)`.
- **Logging**: use `log/slog` (configured via `--log-level`).
- **Configuration keys**: use nested objects (not underscores), e.g. `style.name`.

## Dependency Injection

When adding a new component:

1. Define its constructor function in its package.
2. Add the constructor to `internal/di/wire.go`.
3. Run `wire ./internal/di` to regenerate `wire_gen.go`.
4. If the component is needed by a command, update the command's constructor signature and `commands.NewCommands`.

## Submitting Changes

1. Open an issue or discussion to describe the bug/feature before large refactors.
2. Create a feature branch from `main`.
3. Make your changes, ensuring tests and lints pass (`task validate`).
4. Update relevant documentation:
   - Code changes → update `ARCHITECTURE.md` if architectural impact.
   - Build/test changes → update `DEVELOPMENT.md`.
   - User-facing changes → update the `docs/` tree (Diátaxis quadrants).
5. Open a Pull Request with a clear description and link to the issue.
