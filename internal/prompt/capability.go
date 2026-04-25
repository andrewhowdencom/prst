package prompt

import (
	"strings"
)

// ColorCapability represents the color depth a terminal supports.
type ColorCapability int

const (
	ColorNone      ColorCapability = iota // No color support
	ColorBasic16                          // 16 standard ANSI colors
	Color256                              // 256-color palette
	ColorTrueColor                        // 24-bit RGB (16 million colors)
)

// DetectColorCapability determines the terminal's color capability from
// explicit flags, environment variables, and terminal state.
//
// Resolution order (first match wins):
//
//   1. color flag == "never" → ColorNone
//   2. color flag == "always" → skip all disable checks, detect capability
//   3. color.enabled: false in config → ColorNone
//   4. $NO_COLOR env var present → ColorNone
//   5. $TERM == "dumb" → ColorNone
//   6. color flag == "auto" and non-TTY and not explicitly enabled → ColorNone
//   7. color flag == "auto" and non-TTY and explicitly enabled → ColorTrueColor
//   8. $COLORTERM == "truecolor" || "24bit" → ColorTrueColor
//   9. $TERM contains "256color" → Color256
//  10. Default → ColorBasic16
func DetectColorCapability(
	colorFlag string,
	colorEnabledSet bool,
	colorEnabled bool,
	noColorEnv string,
	termEnv string,
	colorTermEnv string,
	isTerminal bool,
) ColorCapability {
	// Explicit "never" flag overrides everything.
	if colorFlag == "never" {
		return ColorNone
	}

	// "always" means the caller knows colors should be emitted (e.g. init
	// scripts wrapping the output in non-printing sequence markers). Skip
	// all disable checks and detect capability from the environment only.
	if colorFlag == "always" {
		return detectFromEnvironment(colorTermEnv, termEnv)
	}

	// colorFlag == "auto" from here on.

	if colorEnabledSet && !colorEnabled {
		return ColorNone
	}
	if noColorEnv != "" {
		return ColorNone
	}
	if termEnv == "dumb" {
		return ColorNone
	}

	if !isTerminal {
		if colorEnabledSet && colorEnabled {
			return ColorTrueColor
		}
		return ColorNone
	}

	return detectFromEnvironment(colorTermEnv, termEnv)
}

// detectFromEnvironment returns the color capability based on COLORTERM
// and TERM environment variables.
func detectFromEnvironment(colorTermEnv, termEnv string) ColorCapability {
	if colorTermEnv == "truecolor" || colorTermEnv == "24bit" {
		return ColorTrueColor
	}
	if strings.Contains(termEnv, "256color") {
		return Color256
	}
	return ColorBasic16
}
