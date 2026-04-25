# prst

`prst` is a tool for managing Bash prompt strings (PS0–PS4).

See [Bash/Prompt customization](https://wiki.archlinux.org/title/Bash/Prompt_customization) for background.

## Quick Start

```bash
# Automatic installation into your shell
prst install 1

# Or manual installation in ~/.bashrc or ~/.zshrc
eval "$(prst init bash 1)"
```

`prst` reads its configuration from `$XDG_CONFIG_HOME/prst/config.yaml` and prints a prompt string to stdout. All values (username, hostname, path, etc.) are resolved at runtime, and color codes are handled by shell-specific init scripts so cursor positioning stays correct.

By default, with no configuration, it emits a plain classic prompt:

```
user@host:/full/path$
```

## Documentation

Full documentation is available at the links below, structured according to the [Diátaxis framework](https://diataxis.fr).

| | |
|---|---|
| 🎓 [Getting Started](docs/tutorials/getting-started.md) | Step-by-step tutorial for new users. |
| 👐 [How-to Guides](docs/how-to/customize-prompt.md) | Recipes for common customization tasks. |
| 📖 [Configuration Reference](docs/reference/configuration.md) | Full config schema, segment types, and color formats. |
| 📖 [Command Reference](docs/reference/commands.md) | CLI commands, flags, and environment variables. |
| 💡 [Design Rationale](docs/explanation/design.md) | Why `prst` works the way it does. |

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

See [DEVELOPMENT.md](DEVELOPMENT.md) for the full contributor guide and [ARCHITECTURE.md](ARCHITECTURE.md) for system design details.
