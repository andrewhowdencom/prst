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

func TestSegmentConfigContent(t *testing.T) {
	tests := []struct {
		name string
		seg  SegmentConfig
		want string
	}{
		{"user", SegmentConfig{Type: "user"}, resolveUser()},
		{"host default", SegmentConfig{Type: "host"}, resolveHostShort()},
		{"host short", SegmentConfig{Type: "host", Mode: "short"}, resolveHostShort()},
		{"host full", SegmentConfig{Type: "host", Mode: "full"}, resolveHostFull()},
		{"cwd default", SegmentConfig{Type: "cwd"}, resolveCWD()},
		{"cwd full", SegmentConfig{Type: "cwd", Mode: "full"}, resolveCWD()},
		{"cwd basename", SegmentConfig{Type: "cwd", Mode: "basename"}, resolveCWDBasename()},
		{"prompt", SegmentConfig{Type: "prompt"}, resolvePromptChar()},
		{"time short default", SegmentConfig{Type: "time"}, time.Now().Format("15:04")},
		{"time short explicit", SegmentConfig{Type: "time", Format: "short"}, time.Now().Format("15:04")},
		{"time full", SegmentConfig{Type: "time", Format: "full"}, time.Now().Format("15:04:05")},
		{"date", SegmentConfig{Type: "time", Format: "date"}, time.Now().Format("Mon Jan 2")},
		{"newline", SegmentConfig{Type: "newline"}, "\n"},
		{"literal", SegmentConfig{Type: "literal", Text: "hello"}, "hello"},
		{"literal escaped", SegmentConfig{Type: "literal", Text: `\n`}, `\\n`},
		{"unknown", SegmentConfig{Type: "unknown_thing"}, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.seg.Content()
			if tt.seg.Type == "time" || (tt.seg.Type == "time" && tt.seg.Format == "date") {
				if got == "" {
					t.Errorf("Content() = empty, want non-empty")
				}
				return
			}
			if got != tt.want {
				t.Errorf("Content() = %q, want %q", got, tt.want)
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
				return NewColor("green").toANSI(ColorBasic16) + resolveUser() + resetSequence
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
					{Type: "prompt"},
				},
			},
			cap: ColorBasic16,
			wantFn: func() string {
				return NewColor("green").toANSI(ColorBasic16) + resolveUser() + resetSequence +
					"@" +
					NewColor("cyan").toANSI(ColorBasic16) + resolveHostShort() + resetSequence +
					":" +
					NewColor("blue").toANSI(ColorBasic16) + resolveCWD() + resetSequence +
					" " + resolvePromptChar()
			},
		},
		{
			name: "host full mode",
			config: PS1Config{
				Segments: []SegmentConfig{{Type: "host", Mode: "full"}},
			},
			cap:    ColorNone,
			wantFn: func() string { return resolveHostFull() },
		},
		{
			name: "cwd basename mode",
			config: PS1Config{
				Segments: []SegmentConfig{{Type: "cwd", Mode: "basename"}},
			},
			cap:    ColorNone,
			wantFn: func() string { return resolveCWDBasename() },
		},
		{
			name: "time full format",
			config: PS1Config{
				Segments: []SegmentConfig{{Type: "time", Format: "full"}},
			},
			cap: ColorNone,
			wantFn: func() string {
				return time.Now().Format("15:04:05")
			},
		},
		{
			name: "time date format",
			config: PS1Config{
				Segments: []SegmentConfig{{Type: "time", Format: "date"}},
			},
			cap: ColorNone,
			wantFn: func() string {
				return time.Now().Format("Mon Jan 2")
			},
		},
		{
			name: "256-color on capable terminal",
			config: PS1Config{
				Segments: []SegmentConfig{{Type: "user", Color: "256:82"}},
			},
			cap: Color256,
			wantFn: func() string {
				return NewColor("256:82").toANSI(Color256) + resolveUser() + resetSequence
			},
		},
		{
			name: "256-color on basic terminal",
			config: PS1Config{
				Segments: []SegmentConfig{{Type: "user", Color: "256:82"}},
			},
			cap:    ColorBasic16,
			wantFn: func() string { return resolveUser() },
		},
		{
			name: "truecolor on capable terminal",
			config: PS1Config{
				Segments: []SegmentConfig{{Type: "user", Color: "rgb:255,128,0"}},
			},
			cap: ColorTrueColor,
			wantFn: func() string {
				return NewColor("rgb:255,128,0").toANSI(ColorTrueColor) + resolveUser() + resetSequence
			},
		},
		{
			name: "hex color on capable terminal",
			config: PS1Config{
				Segments: []SegmentConfig{{Type: "user", Color: "#ff8000"}},
			},
			cap: ColorTrueColor,
			wantFn: func() string {
				return NewColor("#ff8000").toANSI(ColorTrueColor) + resolveUser() + resetSequence
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
