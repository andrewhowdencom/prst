# prst

`prst` is a tool for managing Bash prompt strings (PS0–PS4).

See [Bash/Prompt customization](https://wiki.archlinux.org/title/Bash/Prompt_customization) for background.

## Usage

```bash
# No-op commands for each prompt string
prst 0
prst 1
prst 2
prst 3
prst 4

# Show version
prst version

# Adjust log level
prst --log-level debug 1
```

## Development

This project uses [Task](https://taskfile.dev/) as its task runner.

```bash
task setup      # Install tools
task generate   # Run code generation (Wire)
task build      # Compile the binary
task test       # Run tests
task lint       # Run linters
task validate   # Run all of the above
```

## Configuration

Configuration is read from (in order of precedence):

1. Command-line flags (`--log-level`)
2. Environment variables (`PRST_LOG_LEVEL`)
3. Config file at `$XDG_CONFIG_HOME/prst/config.yaml`
