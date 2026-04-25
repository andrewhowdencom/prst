package prompt

import (
	"testing"
)

func TestColorToANSI_Basic(t *testing.T) {
	tests := []struct {
		name string
		raw  string
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
			c := NewColor(tt.raw)
			if got := c.toANSI(ColorBasic16); got != tt.want {
				t.Errorf("toANSI(ColorBasic16) = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestColorToANSI_256(t *testing.T) {
	tests := []struct {
		name string
		raw  string
		cap  ColorCapability
		want string
	}{
		{"valid 256 on 256-capable", "256:82", Color256, "\x1b[38;5;82m"},
		{"valid 256 on truecolor", "256:82", ColorTrueColor, "\x1b[38;5;82m"},
		{"valid 256 on basic16", "256:82", ColorBasic16, ""},
		{"valid 256 on none", "256:82", ColorNone, ""},
		{"out of range high", "256:256", Color256, ""},
		{"out of range low", "256:-1", Color256, ""},
		{"invalid format", "256:abc", Color256, ""},
		{"named falls through", "green", Color256, "\x1b[32m"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewColor(tt.raw)
			if got := c.toANSI(tt.cap); got != tt.want {
				t.Errorf("toANSI(%v) = %q, want %q", tt.cap, got, tt.want)
			}
		})
	}
}

func TestColorToANSI_RGB(t *testing.T) {
	tests := []struct {
		name string
		raw  string
		cap  ColorCapability
		want string
	}{
		{"rgb on truecolor", "rgb:255,128,0", ColorTrueColor, "\x1b[38;2;255;128;0m"},
		{"rgb on 256", "rgb:255,128,0", Color256, ""},
		{"rgb on basic16", "rgb:255,128,0", ColorBasic16, ""},
		{"rgb on none", "rgb:255,128,0", ColorNone, ""},
		{"rgb with spaces", "rgb: 255 , 128 , 0 ", ColorTrueColor, "\x1b[38;2;255;128;0m"},
		{"rgb too few parts", "rgb:255,0", ColorTrueColor, ""},
		{"rgb too many parts", "rgb:255,0,0,0", ColorTrueColor, ""},
		{"rgb non-numeric", "rgb:abc,def,ghi", ColorTrueColor, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewColor(tt.raw)
			if got := c.toANSI(tt.cap); got != tt.want {
				t.Errorf("toANSI(%v) = %q, want %q", tt.cap, got, tt.want)
			}
		})
	}
}

func TestColorToANSI_Hex(t *testing.T) {
	tests := []struct {
		name string
		raw  string
		cap  ColorCapability
		want string
	}{
		{"hex on truecolor", "#ff8000", ColorTrueColor, "\x1b[38;2;255;128;0m"},
		{"hex on 256", "#ff8000", Color256, ""},
		{"hex on basic16", "#ff8000", ColorBasic16, ""},
		{"hex on none", "#ff8000", ColorNone, ""},
		{"hex lowercase", "#ff8000", ColorTrueColor, "\x1b[38;2;255;128;0m"},
		{"hex uppercase", "#FF8000", ColorTrueColor, "\x1b[38;2;255;128;0m"},
		{"hex too short", "#fff", ColorTrueColor, ""},
		{"hex too long", "#fffff00", ColorTrueColor, ""},
		{"hex invalid chars", "#ggg", ColorTrueColor, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewColor(tt.raw)
			if got := c.toANSI(tt.cap); got != tt.want {
				t.Errorf("toANSI(%v) = %q, want %q", tt.cap, got, tt.want)
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

func TestHexToRGB(t *testing.T) {
	tests := []struct {
		name    string
		hex     string
		want    [3]int
		wantErr bool
	}{
		{"valid lowercase", "#ff8000", [3]int{255, 128, 0}, false},
		{"valid uppercase", "#FF8000", [3]int{255, 128, 0}, false},
		{"no hash", "ff8000", [3]int{255, 128, 0}, false},
		{"too short", "#fff", [3]int{}, true},
		{"too long", "#ffffff0", [3]int{}, true},
		{"invalid chars", "#ggg", [3]int{}, true},
		{"empty", "", [3]int{}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := hexToRGB(tt.hex)
			if (err != nil) != tt.wantErr {
				t.Fatalf("hexToRGB(%q) error = %v, wantErr %v", tt.hex, err, tt.wantErr)
			}
			if err == nil && got != tt.want {
				t.Errorf("hexToRGB(%q) = %v, want %v", tt.hex, got, tt.want)
			}
		})
	}
}
