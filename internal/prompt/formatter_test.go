package prompt

import "testing"

func TestFormatForShell(t *testing.T) {
	tests := []struct {
		name   string
		raw    string
		shell  string
		want   string
	}{
		{
			name:  "bash single color",
			raw:   "\x1b[33muser\x1b[0m:",
			shell: "bash",
			want:  "\x01\x1b[33m\x02user\x01\x1b[0m\x02:",
		},
		{
			name:  "bash multiple colors",
			raw:   "\x1b[33muser\x1b[0m:\x1b[34m/path\x1b[0m$",
			shell: "bash",
			want:  "\x01\x1b[33m\x02user\x01\x1b[0m\x02:\x01\x1b[34m\x02/path\x01\x1b[0m\x02$",
		},
		{
			name:  "bash 256-color",
			raw:   "\x1b[38;5;82mtext\x1b[0m",
			shell: "bash",
			want:  "\x01\x1b[38;5;82m\x02text\x01\x1b[0m\x02",
		},
		{
			name:  "bash truecolor",
			raw:   "\x1b[38;2;255;128;0mtext\x1b[0m",
			shell: "bash",
			want:  "\x01\x1b[38;2;255;128;0m\x02text\x01\x1b[0m\x02",
		},
		{
			name:  "bash no ansi",
			raw:   "user@host:/path$",
			shell: "bash",
			want:  "user@host:/path$",
		},
		{
			name:  "zsh single color",
			raw:   "\x1b[33muser\x1b[0m:",
			shell: "zsh",
			want:  "%{\x1b[33m%}user%{\x1b[0m%}:",
		},
		{
			name:  "zsh multiple colors",
			raw:   "\x1b[33muser\x1b[0m:\x1b[34m/path\x1b[0m$",
			shell: "zsh",
			want:  "%{\x1b[33m%}user%{\x1b[0m%}:%{\x1b[34m%}/path%{\x1b[0m%}$",
		},
		{
			name:  "raw no shell",
			raw:   "\x1b[33muser\x1b[0m:",
			shell: "",
			want:  "\x1b[33muser\x1b[0m:",
		},
		{
			name:  "unknown shell returns raw",
			raw:   "\x1b[33muser\x1b[0m:",
			shell: "fish",
			want:  "\x1b[33muser\x1b[0m:",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FormatForShell(tt.raw, tt.shell)
			if got != tt.want {
				t.Errorf("FormatForShell(%q, %q) = %q, want %q", tt.raw, tt.shell, got, tt.want)
			}
		})
	}
}
