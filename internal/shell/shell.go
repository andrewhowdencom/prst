// Package shell provides shell detection, rc file resolution, and shell-specific
// prompt initialization script generation.
package shell

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

// Shell represents a supported shell.
type Shell string

const (
	// Bash is the Bourne-again shell.
	Bash Shell = "bash"
	// Zsh is the Z shell.
	Zsh Shell = "zsh"
)

// ValidShells returns the list of supported shell names.
func ValidShells() []string {
	return []string{string(Bash), string(Zsh)}
}

// ParseShell validates a shell name and returns the corresponding Shell value.
func ParseShell(s string) (Shell, error) {
	switch Shell(s) {
	case Bash, Zsh:
		return Shell(s), nil
	default:
		return "", fmt.Errorf("unsupported shell %q: supported shells are %s", s, strings.Join(ValidShells(), ", "))
	}
}

// Detect attempts to detect the user's shell from the environment.
// It first checks the SHELL environment variable, falling back to Bash
// if detection fails or the shell is unsupported.
func Detect() Shell {
	if sh := os.Getenv("SHELL"); sh != "" {
		if base := filepath.Base(sh); base != "" {
			if s, err := ParseShell(base); err == nil {
				return s
			}
		}
	}
	return Bash
}

// RCFile returns the path to the shell's primary runtime configuration file.
func (s Shell) RCFile() string {
	home, _ := os.UserHomeDir()
	switch s {
	case Bash:
		return filepath.Join(home, ".bashrc")
	case Zsh:
		return filepath.Join(home, ".zshrc")
	default:
		return ""
	}
}

// InitScript generates a shell-specific initialization script for the given
// prompt numbers. Each number N corresponds to prst command N and the shell
// variable PSN. The output is suitable for direct evaluation (e.g. via eval).
func (s Shell) InitScript(numbers []int) string {
	sort.Ints(numbers)

	var b strings.Builder

	if s == Zsh {
		b.WriteString("setopt promptsubst\n")
	}

	for _, n := range numbers {
		b.WriteString(s.initFunction(n))
		b.WriteByte('\n')
	}

	for _, n := range numbers {
		b.WriteString(s.initVariable(n))
		b.WriteByte('\n')
	}

	return b.String()
}

// initFunction generates the shell function that wraps prst N output with
// the appropriate non-printing sequence markers for the target shell.
func (s Shell) initFunction(n int) string {
	name := "prst_ps" + strconv.Itoa(n)
	switch s {
	case Bash:
		return fmt.Sprintf("%s() {\n    local raw\n    raw=\"$(prst prompt --color=always %d)\"\n    printf '\\001%%s\\002' \"$raw\"\n}", name, n)
	case Zsh:
		return fmt.Sprintf("%s() {\n    local raw\n    raw=\"$(prst prompt --color=always %d)\"\n    printf '%%{%%s%%}' \"$raw\"\n}", name, n)
	default:
		return ""
	}
}

// initVariable generates the PSN assignment line using the wrapper function.
func (s Shell) initVariable(n int) string {
	name := "prst_ps" + strconv.Itoa(n)
	return fmt.Sprintf("PS%d='$(%s)'", n, name)
}
