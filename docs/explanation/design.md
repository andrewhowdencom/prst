# Design Rationale

This page explains the "why" behind `prst`'s key design choices.

## Why PS0–PS4?

POSIX-compatible shells define five prompt variables:

| Variable | When displayed |
|---|---|
| `PS0` | After a command is read, before it is executed. |
| `PS1` | The primary prompt (what you see most of the time). |
| `PS2` | The continuation prompt (when a command spans multiple lines). |
| `PS3` | The prompt for the `select` built-in. |
| `PS4` | The prefix for execution trace output (`set -x`). |

Most prompt tools only handle `PS1`. `prst` names its commands after the shell variable numbers (`0`–`4`) so the mapping is obvious and future-proof. Even though PS0, PS2, PS3, and PS4 are currently no-ops in `prst`, the command surface is reserved for future expansion without breaking changes.

## Why runtime resolution instead of template expansion?

Many prompt tools generate a static string once and embed shell escape sequences (like `\u` for user or `\w` for cwd in Bash) directly into `PS1`. This is fast but inflexible: the prompt cannot easily react to runtime conditions beyond what the shell provides natively.

`prst` resolves every value at runtime on each invocation. This has several benefits:

- **Accurate cwd**: the path updates correctly even when `$PWD` is manipulated by external tools.
- **Accurate time**: `time_short` and `time_full` show the exact moment the prompt is rendered, not when the rc file was sourced.
- **Portability**: the same YAML config works on any system without relying on shell-specific prompt escape sequences.

The trade-off is a small per-prompt fork/exec overhead. For interactive shells, this is negligible (well under a millisecond on modern hardware).

## Why YAML over shell variables?

A YAML configuration file has advantages over inline shell scripting:

- **Declarative**: segments are an explicit ordered list with named types and colors. There is no string concatenation or escaping to get wrong.
- **Type-safe**: Viper unmarshals the YAML into typed structs. Invalid segment types or malformed colors are caught at runtime and warned about, rather than silently producing garbage output.
- **Portable**: the same config can be shared across machines, checked into dotfiles repositories, and understood without reading shell syntax.
- **Tooling-friendly**: YAML is easy to validate, generate, and edit with standard editors.

## Why XDG Base Directories?

Placing the config at `$XDG_CONFIG_HOME/prst/config.yaml` (falling back to `~/.config/prst/config.yaml`) follows the [XDG Base Directory Specification](https://specifications.freedesktop.org/basedir-spec/basedir-spec-latest.html). This keeps user-specific configuration out of `$HOME` clutter, makes it easy to back up and version-control with standard dotfiles managers, and respects platform conventions on macOS and Linux.

## Why compile-time dependency injection?

`prst` uses [Google Wire](https://github.com/google/wire) to generate the object graph at build time rather than using a reflection-based runtime DI container. This decision was made for three reasons:

1. **Performance**: no reflection overhead at startup.
2. **Safety**: wiring errors are caught by the Go compiler, not at runtime.
3. **Clarity**: the dependency graph is explicit in `internal/di/wire.go`, making it easy for new contributors to trace how components connect.

## Why auto-detect color capability?

Terminal color support varies widely: some terminals support only 16 colors, others 256, and modern terminals support 24-bit RGB. `prst` auto-detects the terminal's capability so users do not need to configure it manually in most cases.

The detection chain explicitly honors user overrides first (`--color=never`, `color.enabled: false`, `$NO_COLOR`) before inspecting `$TERM`, `$COLORTERM`, and TTY state. This follows the [NO_COLOR convention](https://no-color.org/) and respects accessibility preferences while still taking advantage of rich color support when available.

## Why shell-specific init scripts?

ANSI escape sequences are invisible characters, but every shell counts them differently when calculating line length for cursor positioning. If color codes are not marked as non-printing with the correct shell-specific syntax, long command lines wrap incorrectly and editing becomes broken.

Rather than hard-coding Bash-specific `\x01` / `\x02` markers into the prompt output (which would break zsh and any future shell), `prst` separates concerns:

1. `prst N` emits raw ANSI codes.
2. `prst init <shell>` generates a wrapper that applies the correct non-printing markers for the target shell (`\[` `\]` for Bash, `%{...%}` for zsh).
3. `prst install` appends that init script to the correct rc file.

This keeps `prst` universal: the same core generator works for any shell, and only the thin init layer changes.
