package prompt

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
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

func TestWrapNonPrinting(t *testing.T) {
	got := wrapNonPrinting("\x1b[32m")
	want := "\x01\x1b[32m\x02"
	if got != want {
		t.Errorf("wrapNonPrinting() = %q, want %q", got, want)
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

func TestResolveUser(t *testing.T) {
	want := os.Getenv("USER")
	if want == "" {
		want = "?"
	}
	if got := resolveUser(); got != want {
		t.Errorf("resolveUser() = %q, want %q", got, want)
	}
}

func TestResolveHostShort(t *testing.T) {
	got := resolveHostShort()
	if got == "" || got == "?" {
		t.Errorf("resolveHostShort() = %q, want non-empty hostname", got)
	}
	if strings.Contains(got, ".") {
		t.Errorf("resolveHostShort() = %q, should not contain domain part", got)
	}
}

func TestResolveHostFull(t *testing.T) {
	got := resolveHostFull()
	if got == "" || got == "?" {
		t.Errorf("resolveHostFull() = %q, want non-empty hostname", got)
	}
}

func TestResolveCWD(t *testing.T) {
	got := resolveCWD()
	if got == "" || got == "?" {
		t.Errorf("resolveCWD() = %q, want non-empty path", got)
	}
	// If HOME is set and we're inside it, should start with ~
	home := os.Getenv("HOME")
	if home != "" && strings.HasPrefix(got, home) {
		t.Errorf("resolveCWD() = %q, should use ~ prefix when inside HOME", got)
	}
}

func TestResolveCWDBasename(t *testing.T) {
	want := filepath.Base(resolveCWD())
	if got := resolveCWDBasename(); got != want {
		t.Errorf("resolveCWDBasename() = %q, want %q", got, want)
	}
}

func TestResolvePromptChar(t *testing.T) {
	want := "$"
	if os.Geteuid() == 0 {
		want = "#"
	}
	if got := resolvePromptChar(); got != want {
		t.Errorf("resolvePromptChar() = %q, want %q", got, want)
	}
}

func TestResolveSegment(t *testing.T) {
	tests := []struct {
		name string
		typ  string
		want string
	}{
		{"user", "user", resolveUser()},
		{"host", "host", resolveHostShort()},
		{"host_full", "host_full", resolveHostFull()},
		{"cwd", "cwd", resolveCWD()},
		{"cwd_basename", "cwd_basename", resolveCWDBasename()},
		{"prompt_char", "prompt_char", resolvePromptChar()},
		{"time_short", "time_short", time.Now().Format("15:04")},
		{"time_full", "time_full", time.Now().Format("15:04:05")},
		{"date", "date", time.Now().Format("Mon Jan 2")},
		{"newline", "newline", "\n"},
		{"unknown", "unknown_thing", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := resolveSegment(tt.typ)
			if tt.typ == "time_short" || tt.typ == "time_full" || tt.typ == "date" {
				// Time-based segments may differ by a second; just check non-empty
				if got == "" {
					t.Errorf("resolveSegment(%q) = empty, want non-empty", tt.typ)
				}
				return
			}
			if got != tt.want {
				t.Errorf("resolveSegment(%q) = %q, want %q", tt.typ, got, tt.want)
			}
		})
	}
}

func TestPS1GeneratorDefault(t *testing.T) {
	g := NewPS1Generator(PS1Config{})
	want := fmt.Sprintf("%s@%s:%s %s ",
		resolveUser(), resolveHostShort(), resolveCWD(), resolvePromptChar())
	if got := g.Generate(); got != want {
		t.Errorf("Generate() = %q, want %q", got, want)
	}
}

func TestPS1GeneratorSegments(t *testing.T) {
	tests := []struct {
		name   string
		config PS1Config
		wantFn func() string
	}{
		{
			name: "single user segment",
			config: PS1Config{
				Segments: []SegmentConfig{{Type: "user"}},
			},
			wantFn: func() string { return resolveUser() },
		},
		{
			name: "single colored user segment",
			config: PS1Config{
				Segments: []SegmentConfig{{Type: "user", Color: "green"}},
			},
			wantFn: func() string {
				return wrapNonPrinting(Color("green").toANSI()) + resolveUser() + wrapNonPrinting(resetSequence)
			},
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
			wantFn: func() string {
				return wrapNonPrinting(Color("green").toANSI()) + resolveUser() + wrapNonPrinting(resetSequence) +
					"@" +
					wrapNonPrinting(Color("cyan").toANSI()) + resolveHostShort() + wrapNonPrinting(resetSequence) +
					":" +
					wrapNonPrinting(Color("blue").toANSI()) + resolveCWD() + wrapNonPrinting(resetSequence) +
					" $ " + resolvePromptChar()
			},
		},
		{
			name: "literal with backslash escaping",
			config: PS1Config{
				Segments: []SegmentConfig{
					{Type: "literal", Text: `\n`},
				},
			},
			wantFn: func() string { return `\\n` },
		},
		{
			name: "unknown segment type skipped",
			config: PS1Config{
				Segments: []SegmentConfig{
					{Type: "unknown_thing"},
					{Type: "user"},
				},
			},
			wantFn: func() string { return resolveUser() },
		},
		{
			name: "segment with unknown color",
			config: PS1Config{
				Segments: []SegmentConfig{
					{Type: "user", Color: "chartreuse"},
				},
			},
			wantFn: func() string { return resolveUser() },
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := NewPS1Generator(tt.config)
			if got := g.Generate(); got != tt.wantFn() {
				t.Errorf("Generate() = %q, want %q", got, tt.wantFn())
			}
		})
	}
}
