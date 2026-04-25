package prompt

import (
	"os"
	"path/filepath"
	"strings"
	"time"
)

// SegmentConfig describes a single segment of a PS1 prompt.
// The Type field determines which other fields are relevant:
//
//   - "user":       no extra fields
//   - "host":       Mode ("short"|"full", default "short")
//   - "cwd":        Mode ("full"|"basename", default "full")
//   - "prompt":     Style ("char", default "char")
//   - "time":       Format ("short"|"full"|"date", default "short")
//   - "newline":    no extra fields
//   - "literal":    Text (free-form text)
type SegmentConfig struct {
	Type   string `mapstructure:"type"`
	Color  string `mapstructure:"color"`
	Mode   string `mapstructure:"mode"`   // host, cwd
	Format string `mapstructure:"format"` // time
	Style  string `mapstructure:"style"`  // prompt
	Text   string `mapstructure:"text"`   // literal
}

// Content returns the runtime-resolved string for this segment.
// If the segment type is unknown, it returns an empty string.
func (s SegmentConfig) Content() string {
	switch s.Type {
	case "user":
		return resolveUser()
	case "host":
		if s.Mode == "full" {
			return resolveHostFull()
		}
		return resolveHostShort()
	case "cwd":
		if s.Mode == "basename" {
			return resolveCWDBasename()
		}
		return resolveCWD()
	case "prompt":
		return resolvePromptChar()
	case "time":
		switch s.Format {
		case "full":
			return time.Now().Format("15:04:05")
		case "date":
			return time.Now().Format("Mon Jan 2")
		default:
			return time.Now().Format("15:04")
		}
	case "newline":
		return "\n"
	case "literal":
		return literalEscapes(s.Text)
	default:
		return ""
	}
}

func resolveUser() string {
	if u := os.Getenv("USER"); u != "" {
		return u
	}
	return "?"
}

func resolveHostShort() string {
	h, err := os.Hostname()
	if err != nil {
		return "?"
	}
	if i := strings.IndexByte(h, '.'); i >= 0 {
		return h[:i]
	}
	return h
}

func resolveHostFull() string {
	h, err := os.Hostname()
	if err != nil {
		return "?"
	}
	return h
}

func resolveCWD() string {
	wd, err := os.Getwd()
	if err != nil {
		return "?"
	}
	home := os.Getenv("HOME")
	if home != "" && strings.HasPrefix(wd, home) {
		return "~" + strings.TrimPrefix(wd, home)
	}
	return wd
}

func resolveCWDBasename() string {
	wd, err := os.Getwd()
	if err != nil {
		return "?"
	}
	return filepath.Base(wd)
}

func resolvePromptChar() string {
	if os.Geteuid() == 0 {
		return "#"
	}
	return "$"
}

// literalEscapes replaces backslashes with double-backslashes so that Bash
// does not interpret them as prompt escape sequences.
func literalEscapes(text string) string {
	return strings.ReplaceAll(text, `\`, `\\`)
}
