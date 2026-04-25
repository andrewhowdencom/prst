package prompt

import "testing"

func TestDetectColorCapability(t *testing.T) {
	tests := []struct {
		name            string
		colorFlag       string
		colorEnabledSet bool
		colorEnabled    bool
		noColorEnv      string
		termEnv         string
		colorTermEnv    string
		isTerminal      bool
		want            ColorCapability
	}{
		// Explicit disables
		{"flag never", "never", false, false, "", "", "", true, ColorNone},
		{"flag never over everything", "never", true, true, "", "xterm", "truecolor", true, ColorNone},
		{"config disabled", "auto", true, false, "", "", "", true, ColorNone},
		{"NO_COLOR env", "auto", false, false, "1", "", "", true, ColorNone},
		{"TERM dumb", "auto", false, false, "", "dumb", "", true, ColorNone},

		// Non-TTY without explicit enable
		{"non-tty auto", "auto", false, false, "", "xterm", "", false, ColorNone},
		{"non-tty explicit enabled", "auto", true, true, "", "xterm", "", false, ColorTrueColor},

		// always mode
		{"always non-tty", "always", false, false, "", "xterm", "", false, ColorBasic16},
		{"always truecolor", "always", false, false, "", "xterm", "truecolor", false, ColorTrueColor},
		{"always 256", "always", false, false, "", "xterm-256color", "", false, Color256},
		{"always respects NO_COLOR", "always", false, false, "1", "xterm", "truecolor", true, ColorNone},
		{"always respects dumb", "always", false, false, "", "dumb", "truecolor", true, ColorNone},

		// Detection from COLORTERM
		{"truecolor explicit", "auto", false, false, "", "xterm", "truecolor", true, ColorTrueColor},
		{"24bit explicit", "auto", false, false, "", "xterm", "24bit", true, ColorTrueColor},
		{"COLORTERM wins over TERM", "auto", false, false, "", "xterm-256color", "truecolor", true, ColorTrueColor},

		// Detection from TERM
		{"256color in TERM", "auto", false, false, "", "xterm-256color", "", true, Color256},
		{"plain TERM", "auto", false, false, "", "xterm", "", true, ColorBasic16},
		{"vt100 TERM", "auto", false, false, "", "vt100", "", true, ColorBasic16},
		{"linux TERM", "auto", false, false, "", "linux", "", true, ColorBasic16},

		// Priority ordering
		{"flag never over COLORTERM", "never", false, false, "", "xterm", "truecolor", true, ColorNone},
		{"config over COLORTERM", "auto", true, false, "", "xterm", "truecolor", true, ColorNone},
		{"NO_COLOR over COLORTERM", "auto", false, false, "1", "xterm", "truecolor", true, ColorNone},
		{"dumb over everything", "auto", false, false, "", "dumb", "truecolor", true, ColorNone},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := DetectColorCapability(
				tt.colorFlag,
				tt.colorEnabledSet,
				tt.colorEnabled,
				tt.noColorEnv,
				tt.termEnv,
				tt.colorTermEnv,
				tt.isTerminal,
			)
			if got != tt.want {
				t.Errorf("DetectColorCapability() = %v, want %v", got, tt.want)
			}
		})
	}
}
