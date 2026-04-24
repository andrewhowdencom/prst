# prst

`prst` is a tool for managing Bash prompt strings (PS0–PS4).

See [Bash/Prompt customization](https://wiki.archlinux.org/title/Bash/Prompt_customization) for background.

## Usage

Add `prst` to your `.bashrc` so Bash evaluates it before every prompt:

```bash
PS1='$(prst 1)'
```

`prst` reads its configuration from `$XDG_CONFIG_HOME/prst/config.yaml` and prints a prompt string to stdout. All values (username, hostname, path, etc.) are resolved at runtime, and color codes are wrapped in Bash non-printing byte markers so cursor positioning stays correct.

By default, with no configuration, it emits a plain classic prompt:

```
user@host:/full/path$ 
```

### Available commands

```bash
prst 0   # Print PS0 string (pre-command)
prst 1   # Print PS1 string (primary prompt)
prst 2   # Print PS2 string (continuation prompt)
prst 3   # Print PS3 string (select prompt)
prst 4   # Print PS4 string (debug prefix)
prst version          # Show version
prst --log-level debug 1   # Adjust log level for a single run
```

### Configuration

Configuration is read from (in order of precedence):

1. Command-line flags (`--log-level`)
2. Environment variables (`PRST_LOG_LEVEL`)
3. Config file at `$XDG_CONFIG_HOME/prst/config.yaml`

#### PS1 segments

The `ps1` key defines the primary prompt as an ordered list of segments:

```yaml
ps1:
  segments:
    - type: user        color: green
    - type: literal     text: "@"
    - type: host         color: cyan
    - type: literal     text: ":"
    - type: cwd          color: blue
    - type: literal     text: " $"
    - type: prompt_char
```

Each colored segment is automatically wrapped in Bash non-printing byte markers, so Bash calculates cursor position correctly.

#### Segment types

| Type | Description |
|---|---|
| `user` | Current username (`$USER`) |
| `host` | Short hostname (first component) |
| `host_full` | Full hostname (FQDN) |
| `cwd` | Current working directory (`~` substituted for `$HOME`) |
| `cwd_basename` | Basename of cwd |
| `prompt_char` | `#` for root, `$` otherwise |
| `time_short` | Current time as `HH:MM` |
| `time_full` | Current time as `HH:MM:SS` |
| `date` | Current date as `Weekday Month Day` |
| `newline` | Line break |
| `literal` | Free-form text (backslashes are escaped automatically) |

#### Colors

Named colors from the 16-color ANSI palette: `black`, `red`, `green`, `yellow`, `blue`, `magenta`, `cyan`, `white`, and `bright_*` variants (e.g. `bright_blue`). Unknown colors are ignored (segment renders uncolored).

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
