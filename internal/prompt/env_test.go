package prompt

import "testing"

// mockEnvReader is a test-double for EnvReader.
type mockEnvReader struct {
	env        map[string]string
	isTerminal bool
}

func (m *mockEnvReader) Getenv(key string) string {
	return m.env[key]
}

func (m *mockEnvReader) IsTerminal(fd int) bool {
	return m.isTerminal
}

// mockViper is a minimal test-double that satisfies the Viper interface
// expected by DefaultColorCapabilityWithEnv.
type mockViper struct {
	isSet map[string]bool
	bool  map[string]bool
}

func (m *mockViper) IsSet(key string) bool { return m.isSet[key] }
func (m *mockViper) GetBool(key string) bool { return m.bool[key] }

func TestDefaultColorCapabilityWithEnv(t *testing.T) {
	tests := []struct {
		name      string
		colorFlag string
		v         *mockViper
		env       *mockEnvReader
		want      ColorCapability
	}{
		{
			name:      "always truecolor",
			colorFlag: "always",
			v:         &mockViper{},
			env: &mockEnvReader{
				env: map[string]string{
					"COLORTERM": "truecolor",
					"TERM":      "xterm",
				},
				isTerminal: false,
			},
			want: ColorTrueColor,
		},
		{
			name:      "always overrides NO_COLOR",
			colorFlag: "always",
			v:         &mockViper{},
			env: &mockEnvReader{
				env: map[string]string{
					"NO_COLOR":  "1",
					"COLORTERM": "truecolor",
					"TERM":      "xterm",
				},
				isTerminal: false,
			},
			want: ColorTrueColor,
		},
		{
			name:      "always 256",
			colorFlag: "always",
			v:         &mockViper{},
			env: &mockEnvReader{
				env: map[string]string{
					"TERM": "xterm-256color",
				},
				isTerminal: false,
			},
			want: Color256,
		},
		{
			name:      "always basic16",
			colorFlag: "always",
			v:         &mockViper{},
			env: &mockEnvReader{
				env: map[string]string{
					"TERM": "xterm",
				},
				isTerminal: false,
			},
			want: ColorBasic16,
		},
		{
			name:      "never always none",
			colorFlag: "never",
			v:         &mockViper{},
			env: &mockEnvReader{
				env: map[string]string{
					"COLORTERM": "truecolor",
				},
				isTerminal: true,
			},
			want: ColorNone,
		},
		{
			name:      "auto non-tty no config",
			colorFlag: "auto",
			v:         &mockViper{},
			env: &mockEnvReader{
				env: map[string]string{
					"TERM": "xterm",
				},
				isTerminal: false,
			},
			want: ColorNone,
		},
		{
			name:      "auto non-tty explicit enabled",
			colorFlag: "auto",
			v: &mockViper{
				isSet: map[string]bool{"color.enabled": true},
				bool:  map[string]bool{"color.enabled": true},
			},
			env: &mockEnvReader{
				env: map[string]string{
					"TERM": "xterm",
				},
				isTerminal: false,
			},
			want: ColorTrueColor,
		},
		{
			name:      "auto tty COLORTERM truecolor",
			colorFlag: "auto",
			v:         &mockViper{},
			env: &mockEnvReader{
				env: map[string]string{
					"COLORTERM": "truecolor",
					"TERM":      "xterm",
				},
				isTerminal: true,
			},
			want: ColorTrueColor,
		},
		{
			name:      "auto NO_COLOR",
			colorFlag: "auto",
			v:         &mockViper{},
			env: &mockEnvReader{
				env: map[string]string{
					"NO_COLOR":  "1",
					"COLORTERM": "truecolor",
					"TERM":      "xterm",
				},
				isTerminal: true,
			},
			want: ColorNone,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := DefaultColorCapabilityWithEnv(tt.colorFlag, tt.v, tt.env)
			if got != tt.want {
				t.Errorf("DefaultColorCapabilityWithEnv() = %v, want %v", got, tt.want)
			}
		})
	}
}
