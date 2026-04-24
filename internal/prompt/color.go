// Package prompt provides PS1 prompt string generation.
package prompt

// Color represents a named ANSI foreground color.
type Color string

// colorCodes maps named colors to their ANSI SGR foreground parameters.
var colorCodes = map[Color]string{
	"black":            "30",
	"red":              "31",
	"green":            "32",
	"yellow":           "33",
	"blue":             "34",
	"magenta":          "35",
	"cyan":             "36",
	"white":            "37",
	"bright_black":     "90",
	"bright_red":       "91",
	"bright_green":     "92",
	"bright_yellow":    "93",
	"bright_blue":      "94",
	"bright_magenta":   "95",
	"bright_cyan":      "96",
	"bright_white":     "97",
}

// resetSequence is the ANSI SGR reset code.
const resetSequence = "\x1b[0m"

// toANSI returns the SGR escape sequence for this color, or an empty string
// if the color name is unrecognized.
func (c Color) toANSI() string {
	code, ok := colorCodes[c]
	if !ok {
		return ""
	}
	return "\x1b[" + code + "m"
}
