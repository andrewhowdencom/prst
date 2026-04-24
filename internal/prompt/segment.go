package prompt

import "strings"

// SegmentConfig describes a single segment of a PS1 prompt.
type SegmentConfig struct {
	Type  string `mapstructure:"type"`
	Color string `mapstructure:"color"`
	Text  string `mapstructure:"text"`
}

// bashEscapes maps segment types to Bash prompt escape sequences.
// See bash(1) § PROMPTING for the full list.
var bashEscapes = map[string]string{
	"user":         `\u`,
	"host":         `\h`,
	"host_full":    `\H`,
	"cwd":          `\w`,
	"cwd_basename": `\W`,
	"prompt_char":  `\$`,
	"time_short":   `\A`,
	"time_full":    `\t`,
	"date":         `\d`,
	"newline":      `\n`,
}

// literalEscapes replaces backslashes with double-backslashes so that Bash
// does not interpret them as prompt escape sequences.
func literalEscapes(text string) string {
	return strings.ReplaceAll(text, `\`, `\\`)
}
