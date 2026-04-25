package commands

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/andrewhowdencom/prst/internal/shell"
)

func TestParsePromptNumbers(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		want    []int
		wantErr bool
	}{
		{"empty", []string{}, []int{}, false},
		{"single 0", []string{"0"}, []int{0}, false},
		{"single 4", []string{"4"}, []int{4}, false},
		{"multiple", []string{"0", "1", "2"}, []int{0, 1, 2}, false},
		{"negative", []string{"-1"}, nil, true},
		{"too high", []string{"5"}, nil, true},
		{"invalid text", []string{"abc"}, nil, true},
		{"mixed invalid", []string{"1", "abc"}, nil, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parsePromptNumbers(tt.args)
			if (err != nil) != tt.wantErr {
				t.Fatalf("parsePromptNumbers(%v) error = %v, wantErr %v", tt.args, err, tt.wantErr)
			}
			if !tt.wantErr {
				if len(got) != len(tt.want) {
					t.Fatalf("parsePromptNumbers(%v) = %v, want %v", tt.args, got, tt.want)
				}
				for i := range got {
					if got[i] != tt.want[i] {
						t.Fatalf("parsePromptNumbers(%v) = %v, want %v", tt.args, got, tt.want)
					}
				}
			}
		})
	}
}

func TestNewPromptCommand(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		wantErr bool
	}{
		{"valid 0", []string{"0"}, false},
		{"valid 4", []string{"4"}, false},
		{"no args", []string{}, true},
		{"too many args", []string{"1", "2"}, true},
		{"negative", []string{"-1"}, true},
		{"too high", []string{"5"}, true},
		{"non-numeric", []string{"abc"}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := NewPromptCommand(nil, nil)
			buf := new(bytes.Buffer)
			cmd.SetOut(buf)
			cmd.SetErr(buf)
			cmd.SetArgs(tt.args)

			err := cmd.Execute()
			if (err != nil) != tt.wantErr {
				t.Fatalf("Execute() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNewInitCommand(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		wantErr bool
		wantSub string
	}{
		{"bash default", []string{"bash"}, false, "PS1='$(prst_ps1)'"},
		{"bash with numbers", []string{"bash", "0", "1", "2"}, false, "PS0='$(prst_ps0)'"},
		{"zsh default", []string{"zsh"}, false, "setopt promptsubst"},
		{"zsh with numbers", []string{"zsh", "1", "4"}, false, "PS4='$(prst_ps4)'"},
		{"unsupported shell", []string{"fish"}, true, ""},
		{"missing shell", []string{}, true, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := NewInitCommand()
			buf := new(bytes.Buffer)
			cmd.SetOut(buf)
			cmd.SetErr(buf)
			cmd.SetArgs(tt.args)

			err := cmd.Execute()
			if (err != nil) != tt.wantErr {
				t.Fatalf("Execute() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && tt.wantSub != "" {
				if !strings.Contains(buf.String(), tt.wantSub) {
					t.Errorf("output missing %q, got:\n%s", tt.wantSub, buf.String())
				}
			}
		})
	}
}

func TestInitScriptBash(t *testing.T) {
	script := shell.Bash.InitScript([]int{1})
	if !strings.Contains(script, "printf '\\001%s\\002' \"$raw\"") {
		t.Errorf("bash init script missing bash non-printing wrapper (\\001...\\002), got:\n%s", script)
	}
	if strings.Contains(script, "printf '\\[%s\\]'") {
		t.Errorf("bash init script incorrectly uses \\[ \\] instead of \\001 \\002, got:\n%s", script)
	}
	if !strings.Contains(script, "PS1='$(prst_ps1)'") {
		t.Errorf("bash init script missing PS1 assignment, got:\n%s", script)
	}
	if !strings.Contains(script, "prst prompt --color=always 1") {
		t.Errorf("bash init script missing prst prompt --color=always call, got:\n%s", script)
	}
}

func TestInitScriptZsh(t *testing.T) {
	script := shell.Zsh.InitScript([]int{1})
	if !strings.Contains(script, "setopt promptsubst") {
		t.Errorf("zsh init script missing promptsubst, got:\n%s", script)
	}
	if !strings.Contains(script, "printf '%{%s%}' \"$raw\"") {
		t.Errorf("zsh init script missing zsh non-printing wrapper, got:\n%s", script)
	}
	if !strings.Contains(script, "PS1='$(prst_ps1)'") {
		t.Errorf("zsh init script missing PS1 assignment, got:\n%s", script)
	}
	if !strings.Contains(script, "prst prompt --color=always 1") {
		t.Errorf("zsh init script missing prst prompt --color=always call, got:\n%s", script)
	}
}

func TestInstallCommandDryRun(t *testing.T) {
	cmd := NewInstallCommand()
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs([]string{"--shell", "bash", "--dry-run", "1"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "Dry-run") {
		t.Errorf("dry-run output missing expected text, got:\n%s", out)
	}
	if !strings.Contains(out, "prst init bash 1") {
		t.Errorf("dry-run output missing init command, got:\n%s", out)
	}
}

func TestInstallCommandAppendAndRemove(t *testing.T) {
	rcFile := filepath.Join(t.TempDir(), ".bashrc")

	// Override the rc file path by monkey-patching via a custom shell type isn't
	// easy, so we test the internal block helpers directly.
	block := generateBlock(shell.Bash, []int{1, 2})
	if !strings.Contains(block, prstBlockStart) {
		t.Fatalf("generated block missing start marker")
	}
	if !strings.Contains(block, prstBlockEnd) {
		t.Fatalf("generated block missing end marker")
	}
	if !strings.Contains(block, "prst init bash 1 2") {
		t.Fatalf("generated block missing init command")
	}

	// Test appendBlock idempotency.
	if err := appendBlock(rcFile, shell.Bash, []int{1}); err != nil {
		t.Fatalf("appendBlock error = %v", err)
	}

	content, err := os.ReadFile(rcFile)
	if err != nil {
		t.Fatalf("ReadFile error = %v", err)
	}
	if !strings.Contains(string(content), prstBlockStart) {
		t.Fatalf("rc file missing prst block after append")
	}

	// Append again — should not duplicate.
	if err := appendBlock(rcFile, shell.Bash, []int{1}); err != nil {
		t.Fatalf("appendBlock second call error = %v", err)
	}

	content, err = os.ReadFile(rcFile)
	if err != nil {
		t.Fatalf("ReadFile error = %v", err)
	}
	if strings.Count(string(content), prstBlockStart) != 1 {
		t.Fatalf("rc file has duplicate prst blocks")
	}

	// Test removeBlock.
	if err := removeBlock(rcFile); err != nil {
		t.Fatalf("removeBlock error = %v", err)
	}

	content, err = os.ReadFile(rcFile)
	if err != nil {
		t.Fatalf("ReadFile error = %v", err)
	}
	if strings.Contains(string(content), prstBlockStart) {
		t.Fatalf("rc file still has prst block after remove")
	}
}

func TestStripBlock(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    string
	}{
		{
			name:  "no block",
			input: "echo hello\n",
			want:  "echo hello\n",
		},
		{
			name: "single block",
			input: fmt.Sprintf("echo hello\n%s\neval init\n%s\necho world\n", prstBlockStart, prstBlockEnd),
			want:  "echo hello\necho world\n",
		},
		{
			name: "block at start",
			input: fmt.Sprintf("%s\neval init\n%s\necho world\n", prstBlockStart, prstBlockEnd),
			want:  "echo world\n",
		},
		{
			name: "block at end",
			input: fmt.Sprintf("echo hello\n%s\neval init\n%s\n", prstBlockStart, prstBlockEnd),
			want:  "echo hello\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := string(stripBlock([]byte(tt.input)))
			if got != tt.want {
				t.Errorf("stripBlock() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestDetectShell(t *testing.T) {
	t.Run("from SHELL env", func(t *testing.T) {
		t.Setenv("SHELL", "/bin/zsh")
		if got := shell.Detect(); got != shell.Zsh {
			t.Errorf("Detect() = %v, want %v", got, shell.Zsh)
		}
	})

	t.Run("fallback to bash", func(t *testing.T) {
		t.Setenv("SHELL", "")
		if got := shell.Detect(); got != shell.Bash {
			t.Errorf("Detect() = %v, want %v", got, shell.Bash)
		}
	})
}

func TestParseShell(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    shell.Shell
		wantErr bool
	}{
		{"bash", "bash", shell.Bash, false},
		{"zsh", "zsh", shell.Zsh, false},
		{"fish", "fish", "", true},
		{"empty", "", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := shell.ParseShell(tt.input)
			if (err != nil) != tt.wantErr {
				t.Fatalf("ParseShell(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
			}
			if !tt.wantErr && got != tt.want {
				t.Errorf("ParseShell(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestRCFile(t *testing.T) {
	home, _ := os.UserHomeDir()
	tests := []struct {
		name string
		sh   shell.Shell
		want string
	}{
		{"bash", shell.Bash, filepath.Join(home, ".bashrc")},
		{"zsh", shell.Zsh, filepath.Join(home, ".zshrc")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.sh.RCFile(); got != tt.want {
				t.Errorf("RCFile() = %q, want %q", got, tt.want)
			}
		})
	}
}
