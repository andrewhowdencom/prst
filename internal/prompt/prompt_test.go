package prompt

import (
	"testing"
)

func TestColorToANSI(t *testing.T) {
	tests := []struct {
		name string
		c    Color
		want string
	}{
		{"black", "black", "\x1b[30m"},
		{"red", "red", "\x1b[31m"},
		{"green", "green", "\x1b[32m"},
		{"yellow", "yellow", "\x1b[33m"},
		{"blue", "blue", "\x1b[34m"},
		{"magenta", "magenta", "\x1b[35m"},
		{"cyan", "cyan", "\x1b[36m"},
		{"white", "white", "\x1b[37m"},
		{"bright_black", "bright_black", "\x1b[90m"},
		{"bright_red", "bright_red", "\x1b[91m"},
		{"bright_green", "bright_green", "\x1b[92m"},
		{"bright_yellow", "bright_yellow", "\x1b[93m"},
		{"bright_blue", "bright_blue", "\x1b[94m"},
		{"bright_magenta", "bright_magenta", "\x1b[95m"},
		{"bright_cyan", "bright_cyan", "\x1b[96m"},
		{"bright_white", "bright_white", "\x1b[97m"},
		{"unknown", "chartreuse", ""},
		{"empty", "", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.c.toANSI(); got != tt.want {
				t.Errorf("toANSI() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestLiteralEscapes(t *testing.T) {
	tests := []struct {
		name string
		text string
		want string
	}{
		{"no backslash", "hello", "hello"},
		{"single backslash", `hello\world`, `hello\\world`},
		{"multiple backslashes", `a\\b`, `a\\\\b`},
		{"bash escape", `\n\u`, `\\n\\u`},
		{"empty", "", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := literalEscapes(tt.text); got != tt.want {
				t.Errorf("literalEscapes(%q) = %q, want %q", tt.text, got, tt.want)
			}
		})
	}
}

func TestPS1GeneratorGenerate(t *testing.T) {
	tests := []struct {
		name   string
		config PS1Config
		want   string
	}{
		{
			name:   "default empty config",
			config: PS1Config{},
			want:   `\u@\h:\w\$ `,
		},
		{
			name: "default with explicit empty segments",
			config: PS1Config{
				Segments: []SegmentConfig{},
			},
			want: `\u@\h:\w\$ `,
		},
		{
			name: "single user segment",
			config: PS1Config{
				Segments: []SegmentConfig{{Type: "user"}},
			},
			want: `\u`,
		},
		{
			name: "single colored user segment",
			config: PS1Config{
				Segments: []SegmentConfig{{Type: "user", Color: "green"}},
			},
			want: `\[` + "\x1b[32m" + `\]\u\[` + "\x1b[0m" + `\]`,
		},
		{
			name: "multiple segments mixed colors",
			config: PS1Config{
				Segments: []SegmentConfig{
					{Type: "user", Color: "green"},
					{Type: "literal", Text: "@"},
					{Type: "host", Color: "cyan"},
					{Type: "literal", Text: ":"},
					{Type: "cwd", Color: "blue"},
					{Type: "literal", Text: " $ "},
					{Type: "prompt_char"},
				},
			},
			want: `\[` + "\x1b[32m" + `\]\u\[` + "\x1b[0m" + `\]@` +
				`\[` + "\x1b[36m" + `\]\h\[` + "\x1b[0m" + `\]:` +
				`\[` + "\x1b[34m" + `\]\w\[` + "\x1b[0m" + `\] $ \$`,
		},
		{
			name: "literal with backslash escaping",
			config: PS1Config{
				Segments: []SegmentConfig{
					{Type: "literal", Text: `\n`},
				},
			},
			want: `\\n`,
		},
		{
			name: "unknown segment type skipped",
			config: PS1Config{
				Segments: []SegmentConfig{
					{Type: "unknown_thing"},
					{Type: "user"},
				},
			},
			want: `\u`,
		},
		{
			name: "all bash escapes",
			config: PS1Config{
				Segments: []SegmentConfig{
					{Type: "user"},
					{Type: "host"},
					{Type: "host_full"},
					{Type: "cwd"},
					{Type: "cwd_basename"},
					{Type: "prompt_char"},
					{Type: "time_short"},
					{Type: "time_full"},
					{Type: "date"},
					{Type: "newline"},
				},
			},
			want: `\u\h\H\w\W\$\A\t\d\n`,
		},
		{
			name: "segment with unknown color",
			config: PS1Config{
				Segments: []SegmentConfig{
					{Type: "user", Color: "chartreuse"},
				},
			},
			want: `\u`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := NewPS1Generator(tt.config)
			if got := g.Generate(); got != tt.want {
				t.Errorf("Generate() = %q, want %q", got, tt.want)
			}
		})
	}
}
