package prompt

import (
	"fmt"
	"log/slog"
	"strings"

	"github.com/spf13/viper"
)

// Generator produces a PS1 prompt string.
type Generator interface {
	Generate(cap ColorCapability) string
}

// PS1Config holds the user-defined configuration for a PS1 prompt.
type PS1Config struct {
	Segments []SegmentConfig `mapstructure:"segments"`
}

// PS1Generator generates PS1 strings from configured segments.
type PS1Generator struct {
	config PS1Config
}

// NewPS1Config extracts PS1 configuration from Viper.
// If no ps1 key is present, it returns an empty config, which causes the
// generator to fall back to the default plain prompt.
func NewPS1Config(v *viper.Viper) PS1Config {
	var cfg PS1Config
	if err := v.UnmarshalKey("ps1", &cfg); err != nil {
		slog.Warn("failed to unmarshal ps1 config, using defaults", "error", err)
	}
	return cfg
}

// NewPS1Generator creates a generator from the provided config.
func NewPS1Generator(cfg PS1Config) *PS1Generator {
	return &PS1Generator{config: cfg}
}

// Generate produces the PS1 string. If no segments are configured, it
// returns the plain uncolored default prompt with runtime-resolved values.
func (g *PS1Generator) Generate(cap ColorCapability) string {
	if len(g.config.Segments) == 0 {
		return fmt.Sprintf("%s@%s:%s %s ",
			resolveUser(), resolveHostShort(), resolveCWD(), resolvePromptChar())
	}

	var b strings.Builder
	for _, seg := range g.config.Segments {
		content := g.segmentContent(seg)
		if content == "" {
			continue
		}

		color := NewColor(seg.Color)
		ansi := color.toANSI(cap)
		if ansi != "" {
			b.WriteString(wrapNonPrinting(ansi))
		}

		b.WriteString(content)

		if ansi != "" {
			b.WriteString(wrapNonPrinting(resetSequence))
		}
	}

	return b.String()
}

func (g *PS1Generator) segmentContent(seg SegmentConfig) string {
	if seg.Type == "literal" {
		return literalEscapes(seg.Text)
	}

	content := resolveSegment(seg.Type)
	if content == "" {
		slog.Warn("unknown segment type", "type", seg.Type)
	}

	return content
}
