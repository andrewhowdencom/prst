package prompt

import (
	"os"
	"strings"

	"github.com/spf13/viper"
	"golang.org/x/term"
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
//   1. --no-color flag → ColorNone
//   2. color.enabled: false in config → ColorNone
//   3. $NO_COLOR env var present → ColorNone
//   4. $TERM == "dumb" → ColorNone
//   5. Non-TTY and not explicitly enabled → ColorNone
//   6. Non-TTY and explicitly enabled → ColorTrueColor
//   7. $COLORTERM == "truecolor" || "24bit" → ColorTrueColor
//   8. $TERM contains "256color" → Color256
//   9. Default → ColorBasic16
func DetectColorCapability(
	noColorFlag bool,
	colorEnabledSet bool,
	colorEnabled bool,
	noColorEnv string,
	termEnv string,
	colorTermEnv string,
	isTerminal bool,
) ColorCapability {
	if noColorFlag {
		return ColorNone
	}
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
	if colorTermEnv == "truecolor" || colorTermEnv == "24bit" {
		return ColorTrueColor
	}
	if strings.Contains(termEnv, "256color") {
		return Color256
	}
	return ColorBasic16
}

// DefaultColorCapability reads the runtime environment and returns the
// detected color capability for stdout.
func DefaultColorCapability(noColorFlag bool, v *viper.Viper) ColorCapability {
	colorEnabledSet := v.IsSet("color.enabled")
	colorEnabled := v.GetBool("color.enabled")
	return DetectColorCapability(
		noColorFlag,
		colorEnabledSet,
		colorEnabled,
		os.Getenv("NO_COLOR"),
		os.Getenv("TERM"),
		os.Getenv("COLORTERM"),
		term.IsTerminal(int(os.Stdout.Fd())),
	)
}
