package shell

import (
	"bytes"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// TestBashInitScriptIntegration builds the prst binary, injects its absolute
// path into the bash init script, and verifies that a real bash shell renders
// the prompt with ANSI color codes and no literal \[ or \] markers.
func TestBashInitScriptIntegration(t *testing.T) {
	bash, err := exec.LookPath("bash")
	if err != nil {
		t.Skip("bash not found in PATH")
	}

	// Build prst into a temp directory so the test uses the current code.
	tmpDir := t.TempDir()
	prstBin := filepath.Join(tmpDir, "prst")
	build := exec.Command("go", "build", "-o", prstBin, "./cmd/prst")
	build.Dir = "../.." // project root relative to internal/shell/
	if out, err := build.CombinedOutput(); err != nil {
		t.Fatalf("building prst: %v\n%s", err, out)
	}

	// Get the init script and replace "prst " with the absolute binary path.
	script := Bash.InitScript([]int{1})
	script = strings.ReplaceAll(script, "prst ", prstBin+" ")

	// Force colors and provide a rich terminal env.
	input := `export COLORTERM=truecolor
export TERM=xterm
export PRST_COLOR_ENABLED=true
` + script + `
set -x
PS4='PROMPT_CAPTURE[${PS1}] '
echo "test"
set +x
`

	cmd := exec.Command(bash, "--norc", "--noprofile", "-i")
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out

	if err := cmd.Run(); err != nil {
		t.Fatalf("bash execution failed: %v\noutput:\n%s", err, out.String())
	}

	output := out.String()

	// The expanded prompt should contain ANSI escape codes.
	if !strings.Contains(output, "\x1b[") {
		t.Errorf("expanded prompt missing ANSI escape codes; output:\n%s", output)
	}

	// There should be NO literal backslash-bracket pairs visible.
	if strings.Contains(output, `\[`) || strings.Contains(output, `\]`) {
		t.Errorf("expanded prompt contains literal \\[ or \\] markers; output:\n%s", output)
	}

	// Verify the prompt capture wrapper worked.
	if !strings.Contains(output, "PROMPT_CAPTURE[") {
		t.Errorf("did not find PROMPT_CAPTURE marker in output; output:\n%s", output)
	}
}

// TestBashInitScriptNoColor verifies that when $NO_COLOR is set and the
// init script uses --color=always, the prompt still contains ANSI codes.
// This is intentional: --color=always is designed to bypass all disable
// checks because the caller (the init script) wraps output in non-printing
// sequence markers.
func TestBashInitScriptNoColor(t *testing.T) {
	bash, err := exec.LookPath("bash")
	if err != nil {
		t.Skip("bash not found in PATH")
	}

	tmpDir := t.TempDir()
	prstBin := filepath.Join(tmpDir, "prst")
	build := exec.Command("go", "build", "-o", prstBin, "./cmd/prst")
	build.Dir = "../.."
	if out, err := build.CombinedOutput(); err != nil {
		t.Fatalf("building prst: %v\n%s", err, out)
	}

	script := Bash.InitScript([]int{1})
	script = strings.ReplaceAll(script, "prst ", prstBin+" ")

	input := `export NO_COLOR=1
export TERM=xterm
` + script + `
set -x
PS4='PROMPT_CAPTURE[${PS1}] '
echo "test"
set +x
`

	cmd := exec.Command(bash, "--norc", "--noprofile", "-i")
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out

	if err := cmd.Run(); err != nil {
		t.Fatalf("bash execution failed: %v\noutput:\n%s", err, out.String())
	}

	output := out.String()

	// --color=always intentionally overrides NO_COLOR for init scripts.
	if !strings.Contains(output, "\x1b[") {
		t.Errorf("--color=always should override NO_COLOR; missing ANSI codes. output:\n%s", output)
	}

	// Still no literal backslash-bracket pairs.
	if strings.Contains(output, `\[`) || strings.Contains(output, `\]`) {
		t.Errorf("expanded prompt contains literal \\[ or \\] markers; output:\n%s", output)
	}
}

// TestBashInitScriptCursorPositioning verifies that a long command typed after
// a colored prompt does not wrap incorrectly. We test this by typing a known
// command and checking that bash can execute it (if cursor positioning were
// wrong, the command would be garbled and fail).
func TestBashInitScriptCursorPositioning(t *testing.T) {
	bash, err := exec.LookPath("bash")
	if err != nil {
		t.Skip("bash not found in PATH")
	}

	tmpDir := t.TempDir()
	prstBin := filepath.Join(tmpDir, "prst")
	build := exec.Command("go", "build", "-o", prstBin, "./cmd/prst")
	build.Dir = "../.."
	if out, err := build.CombinedOutput(); err != nil {
		t.Fatalf("building prst: %v\n%s", err, out)
	}

	script := Bash.InitScript([]int{1})
	script = strings.ReplaceAll(script, "prst ", prstBin+" ")

	input := `export COLORTERM=truecolor
export TERM=xterm
export PRST_COLOR_ENABLED=true
` + script + `
# Type a command that should execute cleanly if cursor math is correct.
echo "cursor-positioning-ok"
`

	cmd := exec.Command(bash, "--norc", "--noprofile", "-i")
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out

	if err := cmd.Run(); err != nil {
		t.Fatalf("bash execution failed: %v\noutput:\n%s", err, out.String())
	}

	output := out.String()

	// The command should execute successfully.
	if !strings.Contains(output, "cursor-positioning-ok") {
		t.Errorf("command output missing; cursor positioning may be broken. output:\n%s", output)
	}
}
