package prompt

import "regexp"

// ansiSGR matches ANSI Select Graphic Rendition (SGR) escape sequences
// (e.g. \x1b[33m, \x1b[0m, \x1b[38;5;82m, \x1b[38;2;255;128;0m).
var ansiSGR = regexp.MustCompile(`\x1b\[[0-9;]*m`)

// FormatForShell post-processes a raw prompt string, wrapping each ANSI SGR
// escape sequence in the appropriate non-printing sequence markers for the
// target shell.
//
// Supported shells:
//   - "bash": wraps each ANSI sequence in \x01 (SOH) / \x02 (STX)
//   - "zsh":  wraps each ANSI sequence in %{...%}
//   - "":     returns the input unchanged (raw ANSI)
func FormatForShell(raw string, shell string) string {
	switch shell {
	case "bash":
		return ansiSGR.ReplaceAllString(raw, "\x01$0\x02")
	case "zsh":
		return ansiSGR.ReplaceAllString(raw, "%{$0%}")
	default:
		return raw
	}
}
