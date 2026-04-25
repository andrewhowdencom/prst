package prompt

import (
	"os"

	"golang.org/x/term"
)

// EnvReader abstracts environment state so that DetectColorCapability can be
// tested without mutating global environment variables.
type EnvReader interface {
	// Getenv returns the value of the named environment variable.
	Getenv(key string) string
	// IsTerminal reports whether the given file descriptor is a terminal.
	IsTerminal(fd int) bool
}

// osEnvReader is the production EnvReader backed by the real OS.
type osEnvReader struct{}

func (osEnvReader) Getenv(key string) string { return os.Getenv(key) }
func (osEnvReader) IsTerminal(fd int) bool   { return term.IsTerminal(fd) }

// defaultEnvReader is the singleton production EnvReader.
var defaultEnvReader EnvReader = osEnvReader{}

// DefaultColorCapabilityWithEnv reads environment state through the provided
// EnvReader and returns the detected color capability for stdout.
func DefaultColorCapabilityWithEnv(colorFlag string, v interface{ IsSet(string) bool; GetBool(string) bool }, env EnvReader) ColorCapability {
	colorEnabledSet := v.IsSet("color.enabled")
	colorEnabled := v.GetBool("color.enabled")
	return DetectColorCapability(
		colorFlag,
		colorEnabledSet,
		colorEnabled,
		env.Getenv("NO_COLOR"),
		env.Getenv("TERM"),
		env.Getenv("COLORTERM"),
		env.IsTerminal(int(os.Stdout.Fd())),
	)
}

// DefaultColorCapability is a convenience wrapper that uses the production OS
// environment.
func DefaultColorCapability(colorFlag string, v interface{ IsSet(string) bool; GetBool(string) bool }) ColorCapability {
	return DefaultColorCapabilityWithEnv(colorFlag, v, defaultEnvReader)
}
