# Command Reference

This page describes every `prst` command, flag, and environment variable.

## Commands

### `prst prompt [0|1|2|3|4]`

Print a prompt string for the given PS level. Only PS1 is currently implemented; PS0, PS2, PS3, and PS4 are reserved for future expansion and currently print nothing.

```bash
$ prst prompt 1
user@host:/current/path$
```

### `prst init <shell> [0] [1] [2] [3] [4]`

Print shell-specific initialization code. Outputs a script that defines wrapper functions and sets the requested PS variables for the target shell. Supported shells: `bash`, `zsh`.

```bash
# bash — wraps output in \[ \] non-printing markers
$ prst init bash 1
prst_ps1() {
    local raw
    raw="$(prst prompt 1)"
    printf '\[%s\]' "$raw"
}
PS1='$(prst_ps1)'

# zsh — wraps output in %{ %} and enables promptsubst
$ prst init zsh 1
prst_ps1() {
    local raw
    raw="$(prst prompt 1)"
    printf '%{%s%}' "$raw"
}
setopt promptsubst
PS1='$(prst_ps1)'
```

If no prompt numbers are given, it defaults to `1`.

### `prst install [0] [1] [2] [3] [4]`

Automatically installs `prst` into your shell configuration. Detects the current shell from `$SHELL`, appends `eval "$(prst init <shell> ...)"` to the appropriate rc file (e.g. `~/.bashrc`, `~/.zshrc`), and is idempotent — running it again replaces the previous block instead of duplicating it.

```bash
# Auto-detect shell and install PS1
$ prst install 1

# Explicitly target zsh, install PS1 and PS2
$ prst install --shell zsh 1 2

# Preview changes without writing
$ prst install --dry-run 1

# Remove prst from your shell configuration
$ prst install --remove
```

| Flag | Type | Default | Description |
|---|---|---|---|
| `--shell` | `string` | auto-detect | Target shell (`bash` or `zsh`). |
| `--dry-run` | `bool` | `false` | Print what would be written without modifying files. |
| `--remove` | `bool` | `false` | Remove the prst initialization block from the rc file. |

### `prst version`

Print the version of `prst`.

```bash
$ prst version
v0.0.0-dev
```

## Global flags

These flags are available on every command.

| Flag | Type | Default | Description |
|---|---|---|---|
| `--log-level` | `string` | `info` | Set the logging level (`debug`, `info`, `warn`, `error`). |
| `--color` | `string` | `auto` | Color output mode: `auto`, `always`, or `never`. |

## Environment variables

All environment variables use the `PRST_` prefix.

| Variable | Maps to | Description |
|---|---|---|
| `PRST_LOG_LEVEL` | `--log-level` | Override the log level. |
| `PRST_COLOR` | `--color` | Override the color mode (`auto`, `always`, `never`). |
| `PRST_COLOR_ENABLED` | `color.enabled` | Force colors on or off via config key mapping. |
| `NO_COLOR` | — | Set to any value to disable colors (honored independently of `--color`). |

Because Viper replaces `.` with `_` for nested keys, any config file value can be overridden by an environment variable. For example, a hypothetical `style.name` in the config file would map to `PRST_STYLE_NAME`.

## Exit codes

| Code | Meaning |
|---|---|
| `0` | Success. |
| `1` | Error during startup (DI graph construction failed) or command execution failed. |
