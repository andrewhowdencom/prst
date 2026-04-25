// Package prompt provides PS1 prompt string generation.
package prompt

import (
	"fmt"
	"log/slog"
	"strconv"
	"strings"
)

// Color holds a parsed color specification.
// Supported formats:
//   - Named basic colors: "green", "bright_red", …
//   - 256-color index:    "256:82"
//   - RGB decimal:        "rgb:255,0,0"
//   - RGB hex:            "#ff0000"
type Color struct {
	raw string
}

// NewColor creates a Color from a raw string.
func NewColor(raw string) Color {
	return Color{raw: raw}
}

// basicColorCodes maps named colors to their ANSI SGR foreground parameters.
var basicColorCodes = map[string]string{
	"black":          "30",
	"red":            "31",
	"green":          "32",
	"yellow":         "33",
	"blue":           "34",
	"magenta":        "35",
	"cyan":           "36",
	"white":          "37",
	"bright_black":   "90",
	"bright_red":     "91",
	"bright_green":   "92",
	"bright_yellow":  "93",
	"bright_blue":    "94",
	"bright_magenta": "95",
	"bright_cyan":    "96",
	"bright_white":   "97",
}

// toANSI returns the ANSI SGR escape sequence for this color, taking the
// terminal's color capability into account. If the capability is too low
// for the requested color format, or the color string is invalid/empty,
// an empty string is returned.
func (c Color) toANSI(cap ColorCapability) string {
	if cap == ColorNone || c.raw == "" {
		return ""
	}

	if code, ok := basicColorCodes[c.raw]; ok {
		return "\x1b[" + code + "m"
	}

	if strings.HasPrefix(c.raw, "256:") {
		if cap < Color256 {
			return ""
		}
		idx, err := strconv.Atoi(strings.TrimPrefix(c.raw, "256:"))
		if err != nil || idx < 0 || idx > 255 {
			slog.Warn("invalid 256-color index", "color", c.raw)
			return ""
		}
		return fmt.Sprintf("\x1b[38;5;%dm", idx)
	}

	if strings.HasPrefix(c.raw, "rgb:") {
		if cap < ColorTrueColor {
			return ""
		}
		parts := strings.Split(strings.TrimPrefix(c.raw, "rgb:"), ",")
		if len(parts) != 3 {
			slog.Warn("invalid rgb color", "color", c.raw)
			return ""
		}
		r, err1 := strconv.Atoi(strings.TrimSpace(parts[0]))
		g, err2 := strconv.Atoi(strings.TrimSpace(parts[1]))
		b, err3 := strconv.Atoi(strings.TrimSpace(parts[2]))
		if err1 != nil || err2 != nil || err3 != nil {
			slog.Warn("invalid rgb color values", "color", c.raw)
			return ""
		}
		return fmt.Sprintf("\x1b[38;2;%d;%d;%dm", r, g, b)
	}

	if strings.HasPrefix(c.raw, "#") {
		if cap < ColorTrueColor {
			return ""
		}
		rgb, err := hexToRGB(c.raw)
		if err != nil {
			slog.Warn("invalid hex color", "color", c.raw, "error", err)
			return ""
		}
		return fmt.Sprintf("\x1b[38;2;%d;%d;%dm", rgb[0], rgb[1], rgb[2])
	}

	slog.Warn("unknown color", "color", c.raw)
	return ""
}

// resetSequence is the ANSI SGR reset code.
const resetSequence = "\x1b[0m"

// nonPrintStart is the raw SOH byte Bash uses to mark the start of a
// non-printing sequence in PS1 (equivalent to \[).
const nonPrintStart = "\x01"

// nonPrintEnd is the raw STX byte Bash uses to mark the end of a non-printing
// sequence in PS1 (equivalent to \]).
const nonPrintEnd = "\x02"

// wrapNonPrinting wraps s in Bash non-printing byte markers.
func wrapNonPrinting(s string) string {
	return nonPrintStart + s + nonPrintEnd
}

// hexToRGB parses a #RRGGBB hex string into decimal RGB values.
func hexToRGB(hex string) ([3]int, error) {
	hex = strings.TrimPrefix(hex, "#")
	if len(hex) != 6 {
		return [3]int{}, fmt.Errorf("invalid hex length %d", len(hex))
	}
	r, err1 := strconv.ParseInt(hex[0:2], 16, 0)
	g, err2 := strconv.ParseInt(hex[2:4], 16, 0)
	b, err3 := strconv.ParseInt(hex[4:6], 16, 0)
	if err1 != nil || err2 != nil || err3 != nil {
		return [3]int{}, fmt.Errorf("invalid hex value")
	}
	return [3]int{int(r), int(g), int(b)}, nil
}
