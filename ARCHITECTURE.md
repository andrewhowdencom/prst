# Architecture

`prst` is a single-purpose CLI tool that generates Bash prompt strings (PS0–PS4) from a declarative YAML configuration. This document describes the core technologies, package layout, data flow, and key design decisions.

## Core Technologies

| Technology | Purpose |
|---|---|
| [Cobra](https://github.com/spf13/cobra) | CLI framework: root command, subcommands, global/local flags. |
| [Viper](https://github.com/spf13/viper) | Configuration management: YAML, env-var, and flag precedence. |
| [adrg/xdg](https://github.com/adrg/xdg) | XDG Base Directory resolution for the config file path. |
| [Google Wire](https://github.com/google/wire) | Compile-time dependency injection. |
| `log/slog` | Structured logging (Go standard library). |
| `golang.org/x/term` | Terminal detection for color capability. |

## Package Layout

```
cmd/prst/main.go              # Minimal entry point: delegates to DI graph.
internal/
  app/app.go                  # Cobra root command + global flags (--log-level, --no-color).
  commands/commands.go        # Subcommands: 0, 1, 2, 3, 4, version.
  configuration/
    configuration.go          # Viper setup: env prefix PRST_, XDG config path, YAML read.
  di/
    wire.go                   # Wire provider set.
    wire_gen.go               # Generated DI graph (do not edit manually).
  prompt/
    capability.go             # Terminal color capability detection (none, 16, 256, truecolor).
    color.go                  # Color parsing: named, 256, rgb, hex → ANSI SGR codes.
    generator.go              # PS1Generator: assembles segments into the final prompt string.
    segment.go                # Segment resolution: runtime values (user, host, cwd, time, …).
```

## Data Flow

A single invocation of `prst 1` follows this path:

1. **Entry** (`cmd/prst/main.go`)
   - Calls `di.NewApplication()` to build the Cobra command tree.

2. **Dependency Injection** (`internal/di`)
   - Wire resolves providers in order:
     - `configuration.NewViper` → `*viper.Viper`
     - `prompt.NewPS1Config`  → `PS1Config`
     - `prompt.NewPS1Generator` → `*PS1Generator`
     - `commands.NewCommands` → `[]*cobra.Command`
     - `app.NewRootCommand` → `*cobra.Command`

3. **Configuration Load** (`internal/configuration`)
   - Viper is configured with:
     - Env prefix `PRST_`
     - Key replacer `.` → `_` (so `color.enabled` maps to `PRST_COLOR_ENABLED`)
     - Config path resolved via `xdg.ConfigFile("prst/config.yaml")`
   - If the file is missing, Viper continues with defaults/env/flags only.

4. **Command Dispatch** (`internal/commands`)
   - `prst 1` invokes `NewCommand1`, which:
     - Reads the `--no-color` flag.
     - Calls `prompt.DefaultColorCapability(noColor, v)` to determine terminal color support.
     - Invokes `g.Generate(cap)` and writes the result to stdout.

5. **Prompt Generation** (`internal/prompt`)
   - `PS1Generator.Generate(cap)` iterates over the configured segments.
   - For each segment:
     - `segmentContent()` resolves runtime values (e.g., `os.Getenv("USER")`, `os.Getwd()`, `time.Now()`).
     - `Color.toANSI(cap)` converts the segment's color specification into an ANSI escape sequence.
     - ANSI codes are wrapped in Bash non-printing byte markers (`\x01` / `\x02`, equivalent to `\[` / `\]`) so Bash calculates cursor position correctly.
   - If no segments are configured, the generator falls back to a plain default prompt.

## Key Design Decisions

### Runtime Resolution Instead of Template Expansion

Prompt values (username, hostname, cwd) are resolved at **runtime** on every invocation rather than being baked into a static template. This means the prompt is always accurate even when you `cd`, `su`, or `ssh` within the same shell session. The trade-off is a tiny per-prompt fork/exec cost, which is negligible for interactive shells.

### YAML over Shell Variables

Using a YAML config file (instead of exporting `PS1` from `.bashrc`) makes the prompt:

- **Declarative**: segments are an ordered list with explicit types and colors.
- **Portable**: the same config works across machines without shell-specific syntax.
- **Type-safe**: Viper unmarshals into typed structs; invalid keys or colors are caught and warned.

### XDG Base Directories

The configuration lives at `$XDG_CONFIG_HOME/prst/config.yaml` (falling back to `~/.config/prst/config.yaml`). This follows the [XDG Base Directory Specification](https://specifications.freedesktop.org/basedir-spec/basedir-spec-latest.html), keeping user config out of `$HOME` clutter and making it easy to version-control via dotfiles managers.

### Compile-Time Dependency Injection

Google Wire generates the object graph at build time rather than using reflection-based DI frameworks. This keeps startup fast, provides compile-time safety, and makes the dependency graph explicit in `internal/di/wire.go`.

### Color Capability Detection

`prst` auto-detects the richest color format the terminal can safely display. The detection chain (see `internal/prompt/capability.go`) honors explicit user overrides first (`--no-color`, `color.enabled: false`, `$NO_COLOR`) before inspecting `$TERM`, `$COLORTERM`, and TTY state. This follows the [NO_COLOR convention](https://no-color.org/) while still supporting truecolor for modern terminals.
