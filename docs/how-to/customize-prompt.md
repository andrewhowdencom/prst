# Customize Your Prompt

These recipes show how to solve common prompt customization problems with `prst`.

## How do I add colors?

Use the `color` field on any colored segment. `prst` auto-detects the richest format your terminal supports.

```yaml
ps1:
  segments:
    - type: user        color: bright_green
    - type: host         color: "#ff8000"   # hex true-color
```

See the [Reference](../reference/configuration.md) for all supported color formats.

## How do I show only the current directory name?

Use the `cwd_basename` segment instead of `cwd`:

```yaml
ps1:
  segments:
    - type: cwd_basename  color: blue
    - type: literal       text: " $"
```

## How do I add the current time?

```yaml
ps1:
  segments:
    - type: time_short    color: yellow
    - type: literal       text: " "
    - type: cwd           color: blue
```

Use `time_full` for seconds (`HH:MM:SS`) or `date` for `Weekday Month Day`.

## How do I force colors even when piped?

By default, `prst` disables colors when stdout is not a TTY. To force colors (useful if your terminal emulator handles piped output):

```yaml
color:
  enabled: true
```

Or for a single run:

```bash
prst prompt --color=always 1
```

## How do I add a newline before the prompt?

Use the `newline` segment:

```yaml
ps1:
  segments:
    - type: newline
    - type: user         color: green
    - type: literal       text: "@"
    - type: host          color: cyan
    - type: literal       text: ":"
    - type: cwd           color: blue
    - type: literal       text: "\$ "
```

## How do I use a custom separator?

Use the `literal` segment with the text you want:

```yaml
ps1:
  segments:
    - type: user         color: magenta
    - type: literal       text: " in "
    - type: cwd           color: magenta
    - type: literal       text: " \u276f "   # Unicode ❯
```

Note: `literal` segments automatically escape backslashes so Bash does not interpret them as prompt escape sequences.

## How do I show the full hostname?

Use the `host_full` segment instead of `host`:

```yaml
ps1:
  segments:
    - type: user         color: green
    - type: literal       text: "@"
    - type: host_full     color: cyan
```

## How do I make the prompt look different for root?

The `prompt_char` segment already adapts: it prints `#` when running as root and `$` otherwise. You can pair it with a colored segment for extra visibility:

```yaml
ps1:
  segments:
    - type: user         color: red
    - type: literal       text: "# "
    - type: prompt_char
```

## How do I debug my configuration?

Run `prst` with debug logging to see which segments are resolved and what color capability is detected:

```bash
prst prompt --log-level debug 1
```
