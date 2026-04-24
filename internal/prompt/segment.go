package prompt

import (
	"os"
	"path/filepath"
	"strings"
	"time"
)

// SegmentConfig describes a single segment of a PS1 prompt.
type SegmentConfig struct {
	Type  string `mapstructure:"type"`
	Color string `mapstructure:"color"`
	Text  string `mapstructure:"text"`
}

// resolveSegment returns the runtime-resolved string for a segment type.
func resolveSegment(segType string) string {
	switch segType {
	case "user":
		return resolveUser()
	case "host":
		return resolveHostShort()
	case "host_full":
		return resolveHostFull()
	case "cwd":
		return resolveCWD()
	case "cwd_basename":
		return resolveCWDBasename()
	case "prompt_char":
		return resolvePromptChar()
	case "time_short":
		return time.Now().Format("15:04")
	case "time_full":
		return time.Now().Format("15:04:05")
	case "date":
		return time.Now().Format("Mon Jan 2")
	case "newline":
		return "\n"
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
