# Command Reference

This page describes every `prst` command, flag, and environment variable.

## Commands

### `prst 0`

Print the PS0 prompt string (pre-command). Currently a no-op; reserved for future use.

### `prst 1`

Print the PS1 prompt string (primary prompt). This is the command you typically use in your `.bashrc`.

```bash
PS1='$(prst 1)'
```

### `prst 2`

Print the PS2 prompt string (continuation prompt). Currently a no-op; reserved for future use.

### `prst 3`

Print the PS3 prompt string (select prompt). Currently a no-op; reserved for future use.

### `prst 4`

Print the PS4 prompt string (debug prefix). Currently a no-op; reserved for future use.

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
| `--no-color` | `bool` | `false` | Disable colored output for this run. |

## Environment variables

All environment variables use the `PRST_` prefix.

| Variable | Maps to | Description |
|---|---|---|
| `PRST_LOG_LEVEL` | `--log-level` | Override the log level. |
| `PRST_NO_COLOR` | `--no-color` | Set to any value to disable colors (note: this is the env-var mapping of the flag; the standard `$NO_COLOR` convention is also honored independently). |
| `PRST_COLOR_ENABLED` | `color.enabled` | Force colors on or off via config key mapping. |

Because Viper replaces `.` with `_` for nested keys, any config file value can be overridden by an environment variable. For example, a hypothetical `style.name` in the config file would map to `PRST_STYLE_NAME`.

## Exit codes

| Code | Meaning |
|---|---|
| `0` | Success. |
| `1` | Error during startup (DI graph construction failed) or command execution failed. |
