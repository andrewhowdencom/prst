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
//
//   1. color flag == "never" → ColorNone
//   2. color.enabled: false in config → ColorNone
//   3. $NO_COLOR env var present → ColorNone
//   4. $TERM == "dumb" → ColorNone
//   5. color flag == "auto" and non-TTY and not explicitly enabled → ColorNone
//   6. color flag == "auto" and non-TTY and explicitly enabled → ColorTrueColor
//   7. $COLORTERM == "truecolor" || "24bit" → ColorTrueColor
//   8. $TERM contains "256color" → Color256
//   9. Default → ColorBasic16
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

	// These explicit disable signals are always respected.
	if colorEnabledSet && !colorEnabled {
		return ColorNone
	}
	if noColorEnv != "" {
		return ColorNone
	}
	if termEnv == "dumb" {
		return ColorNone
	}

	// In auto mode, suppress colors for non-TTYs unless the user
	// has explicitly enabled them via configuration.
	if colorFlag == "auto" && !isTerminal {
		if colorEnabledSet && colorEnabled {
			return ColorTrueColor
		}
		return ColorNone
	}

	// At this point:
	//   - colorFlag is "auto" and we are in a TTY, OR
	//   - colorFlag is "always" (isTerminal is irrelevant).
	// Detect capability from the environment.

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
func DefaultColorCapability(colorFlag string, v *viper.Viper) ColorCapability {
	colorEnabledSet := v.IsSet("color.enabled")
	colorEnabled := v.GetBool("color.enabled")
	return DetectColorCapability(
		colorFlag,
		colorEnabledSet,
		colorEnabled,
		os.Getenv("NO_COLOR"),
		os.Getenv("TERM"),
		os.Getenv("COLORTERM"),
		term.IsTerminal(int(os.Stdout.Fd())),
	)
}
