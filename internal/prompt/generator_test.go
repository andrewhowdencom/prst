package prompt

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

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
	for _, cap := range []ColorCapability{ColorNone, ColorBasic16, Color256, ColorTrueColor} {
		t.Run(fmt.Sprintf("cap_%d", cap), func(t *testing.T) {
			if got := g.Generate(cap); got != want {
				t.Errorf("Generate(%v) = %q, want %q", cap, got, want)
			}
		})
	}
}

func TestPS1GeneratorSegments(t *testing.T) {
	tests := []struct {
		name   string
		config PS1Config
		cap    ColorCapability
		wantFn func() string
	}{
		{
			name:   "single user segment no color",
			config: PS1Config{Segments: []SegmentConfig{{Type: "user"}}},
			cap:    ColorNone,
			wantFn: func() string { return resolveUser() },
		},
		{
			name:   "single colored user segment basic",
			config: PS1Config{Segments: []SegmentConfig{{Type: "user", Color: "green"}}},
			cap:    ColorBasic16,
			wantFn: func() string {
				return wrapNonPrinting(NewColor("green").toANSI(ColorBasic16)) + resolveUser() + wrapNonPrinting(resetSequence)
			},
		},
		{
			name:   "colored user segment on none",
			config: PS1Config{Segments: []SegmentConfig{{Type: "user", Color: "green"}}},
			cap:    ColorNone,
			wantFn: func() string { return resolveUser() },
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
					{Type: "literal", Text: " "},
					{Type: "prompt_char"},
				},
			},
			cap: ColorBasic16,
			wantFn: func() string {
				return wrapNonPrinting(NewColor("green").toANSI(ColorBasic16)) + resolveUser() + wrapNonPrinting(resetSequence) +
					"@" +
					wrapNonPrinting(NewColor("cyan").toANSI(ColorBasic16)) + resolveHostShort() + wrapNonPrinting(resetSequence) +
					":" +
					wrapNonPrinting(NewColor("blue").toANSI(ColorBasic16)) + resolveCWD() + wrapNonPrinting(resetSequence) +
					" " + resolvePromptChar()
			},
		},
		{
			name: "256-color on capable terminal",
			config: PS1Config{
				Segments: []SegmentConfig{{Type: "user", Color: "256:82"}},
			},
			cap: Color256,
			wantFn: func() string {
				return wrapNonPrinting(NewColor("256:82").toANSI(Color256)) + resolveUser() + wrapNonPrinting(resetSequence)
			},
		},
		{
			name: "256-color on basic terminal",
			config: PS1Config{
				Segments: []SegmentConfig{{Type: "user", Color: "256:82"}},
			},
			cap: ColorBasic16,
			wantFn: func() string { return resolveUser() },
		},
		{
			name: "truecolor on capable terminal",
			config: PS1Config{
				Segments: []SegmentConfig{{Type: "user", Color: "rgb:255,128,0"}},
			},
			cap: ColorTrueColor,
			wantFn: func() string {
				return wrapNonPrinting(NewColor("rgb:255,128,0").toANSI(ColorTrueColor)) + resolveUser() + wrapNonPrinting(resetSequence)
			},
		},
		{
			name: "hex color on capable terminal",
			config: PS1Config{
				Segments: []SegmentConfig{{Type: "user", Color: "#ff8000"}},
			},
			cap: ColorTrueColor,
			wantFn: func() string {
				return wrapNonPrinting(NewColor("#ff8000").toANSI(ColorTrueColor)) + resolveUser() + wrapNonPrinting(resetSequence)
			},
		},
		{
			name: "literal with backslash escaping",
			config: PS1Config{
				Segments: []SegmentConfig{{Type: "literal", Text: `\n`}},
			},
			cap:    ColorNone,
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
			cap:    ColorNone,
			wantFn: func() string { return resolveUser() },
		},
		{
			name: "segment with unknown color",
			config: PS1Config{
				Segments: []SegmentConfig{
					{Type: "user", Color: "chartreuse"},
				},
			},
			cap:    ColorBasic16,
			wantFn: func() string { return resolveUser() },
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := NewPS1Generator(tt.config)
			if got := g.Generate(tt.cap); got != tt.wantFn() {
				t.Errorf("Generate(%v) = %q, want %q", tt.cap, got, tt.wantFn())
			}
		})
	}
}
