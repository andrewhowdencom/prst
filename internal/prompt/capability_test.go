package prompt

import "testing"

func TestDetectColorCapability(t *testing.T) {
	tests := []struct {
		name            string
		noColorFlag     bool
		colorEnabledSet bool
		colorEnabled    bool
		noColorEnv      string
		termEnv         string
		colorTermEnv    string
		isTerminal      bool
		want            ColorCapability
	}{
		// Explicit disables
		{"no-color flag", true, false, false, "", "", "", true, ColorNone},
		{"config disabled", false, true, false, "", "", "", true, ColorNone},
		{"NO_COLOR env", false, false, false, "1", "", "", true, ColorNone},
		{"TERM dumb", false, false, false, "", "dumb", "", true, ColorNone},

		// Non-TTY without explicit enable
		{"non-tty auto", false, false, false, "", "xterm", "", false, ColorNone},
		{"non-tty explicit enabled", false, true, true, "", "xterm", "", false, ColorTrueColor},

		// Detection from COLORTERM
		{"truecolor explicit", false, false, false, "", "xterm", "truecolor", true, ColorTrueColor},
		{"24bit explicit", false, false, false, "", "xterm", "24bit", true, ColorTrueColor},
		{"COLORTERM wins over TERM", false, false, false, "", "xterm-256color", "truecolor", true, ColorTrueColor},

		// Detection from TERM
		{"256color in TERM", false, false, false, "", "xterm-256color", "", true, Color256},
		{"plain TERM", false, false, false, "", "xterm", "", true, ColorBasic16},
		{"vt100 TERM", false, false, false, "", "vt100", "", true, ColorBasic16},
		{"linux TERM", false, false, false, "", "linux", "", true, ColorBasic16},

		// Priority ordering
		{"flag over COLORTERM", true, false, false, "", "xterm", "truecolor", true, ColorNone},
		{"config over COLORTERM", false, true, false, "", "xterm", "truecolor", true, ColorNone},
		{"NO_COLOR over COLORTERM", false, false, false, "1", "xterm", "truecolor", true, ColorNone},
		{"dumb over everything", false, false, false, "", "dumb", "truecolor", true, ColorNone},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := DetectColorCapability(
				tt.noColorFlag,
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
