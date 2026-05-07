package prompt

import (
	"os"
	"os/exec"
	"strings"
	"testing"
)

func skipIfNoGit(t *testing.T) {
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not found in PATH, skipping integration tests")
	}
}

func TestParseAheadBehind(t *testing.T) {
	tests := []struct {
		name         string
		input        string
		wantAhead    int
		wantBehind   int
	}{
		{"ahead two behind zero", "2\t0", 2, 0},
		{"ahead zero behind one", "0\t1", 0, 1},
		{"diverged", "3\t2", 3, 2},
		{"up to date", "0\t0", 0, 0},
		{"empty", "", 0, 0},
		{"single value", "5", 0, 0},
		{"non-numeric", "a\tb", 0, 0},
		{"leading and trailing whitespace", " 2\t0 ", 2, 0},
		{"spaces instead of tabs", "1 2", 0, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotAhead, gotBehind := parseAheadBehind(tt.input)
			if gotAhead != tt.wantAhead || gotBehind != tt.wantBehind {
				t.Errorf("parseAheadBehind(%q) = (%d, %d), want (%d, %d)",
					tt.input, gotAhead, gotBehind, tt.wantAhead, tt.wantBehind)
			}
		})
	}
}

func TestIsDirtyFromPorcelain(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  bool
	}{
		{"clean", "", false},
		{"whitespace only", "   \n  ", false},
		{"modified file", " M file.go\n", true},
		{"staged file", "M  file.go\n", true},
		{"untracked file", "?? file.go\n", true},
		{"multiple changes", " M a.go\nM  b.go\n?? c.go\n", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isDirtyFromPorcelain(tt.input); got != tt.want {
				t.Errorf("isDirtyFromPorcelain(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestResolveGitNotInRepo(t *testing.T) {
	dir := t.TempDir()
	t.Chdir(dir)
	got := resolveGit("{{.Git.Branch}}")
	if got != "" {
		t.Errorf("resolveGit outside repo = %q, want empty", got)
	}
}

func TestResolveGit(t *testing.T) {
	skipIfNoGit(t)

	t.Run("clean repo on branch", func(t *testing.T) {
		dir := t.TempDir()
		t.Chdir(dir)
		if err := exec.Command("git", "init").Run(); err != nil {
			t.Fatalf("git init: %v", err)
		}
		if err := exec.Command("git", "config", "user.email", "test@example.com").Run(); err != nil {
			t.Fatalf("git config user.email: %v", err)
		}
		if err := exec.Command("git", "config", "user.name", "Test").Run(); err != nil {
			t.Fatalf("git config user.name: %v", err)
		}
		if err := os.WriteFile("a.txt", []byte("a"), 0o644); err != nil {
			t.Fatalf("write file: %v", err)
		}
		if err := exec.Command("git", "add", "a.txt").Run(); err != nil {
			t.Fatalf("git add: %v", err)
		}
		if err := exec.Command("git", "commit", "-m", "first").Run(); err != nil {
			t.Fatalf("git commit: %v", err)
		}

		out, err := exec.Command("git", "branch", "--show-current").Output()
		if err != nil {
			t.Fatalf("git branch: %v", err)
		}
		branch := strings.TrimSpace(string(out))

		got := resolveGit("{{.Git.Branch}}")
		if got != branch {
			t.Errorf("resolveGit = %q, want %q", got, branch)
		}
	})

	t.Run("clean repo not dirty", func(t *testing.T) {
		dir := t.TempDir()
		t.Chdir(dir)
		if err := exec.Command("git", "init").Run(); err != nil {
			t.Fatalf("git init: %v", err)
		}
		if err := exec.Command("git", "config", "user.email", "test@example.com").Run(); err != nil {
			t.Fatalf("git config user.email: %v", err)
		}
		if err := exec.Command("git", "config", "user.name", "Test").Run(); err != nil {
			t.Fatalf("git config user.name: %v", err)
		}
		if err := os.WriteFile("a.txt", []byte("a"), 0o644); err != nil {
			t.Fatalf("write file: %v", err)
		}
		if err := exec.Command("git", "add", "a.txt").Run(); err != nil {
			t.Fatalf("git add: %v", err)
		}
		if err := exec.Command("git", "commit", "-m", "first").Run(); err != nil {
			t.Fatalf("git commit: %v", err)
		}

		got := resolveGit("{{.Git.Dirty}}")
		if got != "false" {
			t.Errorf("resolveGit = %q, want false", got)
		}
	})

	t.Run("dirty repo", func(t *testing.T) {
		dir := t.TempDir()
		t.Chdir(dir)
		if err := exec.Command("git", "init").Run(); err != nil {
			t.Fatalf("git init: %v", err)
		}
		if err := exec.Command("git", "config", "user.email", "test@example.com").Run(); err != nil {
			t.Fatalf("git config user.email: %v", err)
		}
		if err := exec.Command("git", "config", "user.name", "Test").Run(); err != nil {
			t.Fatalf("git config user.name: %v", err)
		}
		if err := os.WriteFile("a.txt", []byte("a"), 0o644); err != nil {
			t.Fatalf("write file: %v", err)
		}
		if err := exec.Command("git", "add", "a.txt").Run(); err != nil {
			t.Fatalf("git add: %v", err)
		}
		if err := exec.Command("git", "commit", "-m", "first").Run(); err != nil {
			t.Fatalf("git commit: %v", err)
		}
		if err := os.WriteFile("b.txt", []byte("b"), 0o644); err != nil {
			t.Fatalf("write file: %v", err)
		}

		got := resolveGit("{{.Git.Dirty}}")
		if got != "true" {
			t.Errorf("resolveGit = %q, want true", got)
		}
	})

	t.Run("dirty marker template", func(t *testing.T) {
		dir := t.TempDir()
		t.Chdir(dir)
		if err := exec.Command("git", "init").Run(); err != nil {
			t.Fatalf("git init: %v", err)
		}
		if err := exec.Command("git", "config", "user.email", "test@example.com").Run(); err != nil {
			t.Fatalf("git config user.email: %v", err)
		}
		if err := exec.Command("git", "config", "user.name", "Test").Run(); err != nil {
			t.Fatalf("git config user.name: %v", err)
		}
		if err := os.WriteFile("a.txt", []byte("a"), 0o644); err != nil {
			t.Fatalf("write file: %v", err)
		}
		if err := exec.Command("git", "add", "a.txt").Run(); err != nil {
			t.Fatalf("git add: %v", err)
		}
		if err := exec.Command("git", "commit", "-m", "first").Run(); err != nil {
			t.Fatalf("git commit: %v", err)
		}
		if err := os.WriteFile("b.txt", []byte("b"), 0o644); err != nil {
			t.Fatalf("write file: %v", err)
		}

		got := resolveGit("{{if .Git.Dirty}}*{{end}}")
		if got != "*" {
			t.Errorf("resolveGit = %q, want *", got)
		}
	})

	t.Run("no upstream", func(t *testing.T) {
		dir := t.TempDir()
		t.Chdir(dir)
		if err := exec.Command("git", "init").Run(); err != nil {
			t.Fatalf("git init: %v", err)
		}
		if err := exec.Command("git", "config", "user.email", "test@example.com").Run(); err != nil {
			t.Fatalf("git config user.email: %v", err)
		}
		if err := exec.Command("git", "config", "user.name", "Test").Run(); err != nil {
			t.Fatalf("git config user.name: %v", err)
		}
		if err := os.WriteFile("a.txt", []byte("a"), 0o644); err != nil {
			t.Fatalf("write file: %v", err)
		}
		if err := exec.Command("git", "add", "a.txt").Run(); err != nil {
			t.Fatalf("git add: %v", err)
		}
		if err := exec.Command("git", "commit", "-m", "first").Run(); err != nil {
			t.Fatalf("git commit: %v", err)
		}

		got := resolveGit("{{.Git.Ahead}}{{.Git.Behind}}")
		if got != "00" {
			t.Errorf("resolveGit = %q, want 00", got)
		}
	})

	t.Run("ahead of upstream", func(t *testing.T) {
		dir := t.TempDir()
		t.Chdir(dir)
		if err := exec.Command("git", "init").Run(); err != nil {
			t.Fatalf("git init: %v", err)
		}
		if err := exec.Command("git", "config", "user.email", "test@example.com").Run(); err != nil {
			t.Fatalf("git config user.email: %v", err)
		}
		if err := exec.Command("git", "config", "user.name", "Test").Run(); err != nil {
			t.Fatalf("git config user.name: %v", err)
		}
		if err := os.WriteFile("a.txt", []byte("a"), 0o644); err != nil {
			t.Fatalf("write file: %v", err)
		}
		if err := exec.Command("git", "add", "a.txt").Run(); err != nil {
			t.Fatalf("git add: %v", err)
		}
		if err := exec.Command("git", "commit", "-m", "first").Run(); err != nil {
			t.Fatalf("git commit: %v", err)
		}
		if err := os.WriteFile("b.txt", []byte("b"), 0o644); err != nil {
			t.Fatalf("write file: %v", err)
		}
		if err := exec.Command("git", "add", "b.txt").Run(); err != nil {
			t.Fatalf("git add: %v", err)
		}
		if err := exec.Command("git", "commit", "-m", "second").Run(); err != nil {
			t.Fatalf("git commit: %v", err)
		}

		out, err := exec.Command("git", "branch", "--show-current").Output()
		if err != nil {
			t.Fatalf("git branch: %v", err)
		}
		originalBranch := strings.TrimSpace(string(out))

		if err := exec.Command("git", "checkout", "-b", "feature").Run(); err != nil {
			t.Fatalf("git checkout: %v", err)
		}
		if err := os.WriteFile("c.txt", []byte("c"), 0o644); err != nil {
			t.Fatalf("write file: %v", err)
		}
		if err := exec.Command("git", "add", "c.txt").Run(); err != nil {
			t.Fatalf("git add: %v", err)
		}
		if err := exec.Command("git", "commit", "-m", "third").Run(); err != nil {
			t.Fatalf("git commit: %v", err)
		}
		if err := exec.Command("git", "branch", "--set-upstream-to", originalBranch).Run(); err != nil {
			t.Fatalf("git branch --set-upstream-to: %v", err)
		}

		got := resolveGit("{{.Git.Ahead}}")
		if got != "1" {
			t.Errorf("resolveGit = %q, want 1", got)
		}
	})

	t.Run("behind upstream", func(t *testing.T) {
		dir := t.TempDir()
		t.Chdir(dir)
		if err := exec.Command("git", "init").Run(); err != nil {
			t.Fatalf("git init: %v", err)
		}
		if err := exec.Command("git", "config", "user.email", "test@example.com").Run(); err != nil {
			t.Fatalf("git config user.email: %v", err)
		}
		if err := exec.Command("git", "config", "user.name", "Test").Run(); err != nil {
			t.Fatalf("git config user.name: %v", err)
		}
		if err := os.WriteFile("a.txt", []byte("a"), 0o644); err != nil {
			t.Fatalf("write file: %v", err)
		}
		if err := exec.Command("git", "add", "a.txt").Run(); err != nil {
			t.Fatalf("git add: %v", err)
		}
		if err := exec.Command("git", "commit", "-m", "first").Run(); err != nil {
			t.Fatalf("git commit: %v", err)
		}
		if err := os.WriteFile("b.txt", []byte("b"), 0o644); err != nil {
			t.Fatalf("write file: %v", err)
		}
		if err := exec.Command("git", "add", "b.txt").Run(); err != nil {
			t.Fatalf("git add: %v", err)
		}
		if err := exec.Command("git", "commit", "-m", "second").Run(); err != nil {
			t.Fatalf("git commit: %v", err)
		}

		out, err := exec.Command("git", "branch", "--show-current").Output()
		if err != nil {
			t.Fatalf("git branch: %v", err)
		}
		originalBranch := strings.TrimSpace(string(out))

		if err := exec.Command("git", "checkout", "-b", "feature").Run(); err != nil {
			t.Fatalf("git checkout: %v", err)
		}
		if err := exec.Command("git", "reset", "--hard", "HEAD~1").Run(); err != nil {
			t.Fatalf("git reset: %v", err)
		}
		if err := exec.Command("git", "branch", "--set-upstream-to", originalBranch).Run(); err != nil {
			t.Fatalf("git branch --set-upstream-to: %v", err)
		}

		got := resolveGit("{{.Git.Behind}}")
		if got != "1" {
			t.Errorf("resolveGit = %q, want 1", got)
		}
	})

	t.Run("detached HEAD", func(t *testing.T) {
		dir := t.TempDir()
		t.Chdir(dir)
		if err := exec.Command("git", "init").Run(); err != nil {
			t.Fatalf("git init: %v", err)
		}
		if err := exec.Command("git", "config", "user.email", "test@example.com").Run(); err != nil {
			t.Fatalf("git config user.email: %v", err)
		}
		if err := exec.Command("git", "config", "user.name", "Test").Run(); err != nil {
			t.Fatalf("git config user.name: %v", err)
		}
		if err := os.WriteFile("a.txt", []byte("a"), 0o644); err != nil {
			t.Fatalf("write file: %v", err)
		}
		if err := exec.Command("git", "add", "a.txt").Run(); err != nil {
			t.Fatalf("git add: %v", err)
		}
		if err := exec.Command("git", "commit", "-m", "first").Run(); err != nil {
			t.Fatalf("git commit: %v", err)
		}

		shaOut, err := exec.Command("git", "rev-parse", "HEAD").Output()
		if err != nil {
			t.Fatalf("git rev-parse HEAD: %v", err)
		}
		sha := strings.TrimSpace(string(shaOut))

		if err := exec.Command("git", "checkout", sha).Run(); err != nil {
			t.Fatalf("git checkout: %v", err)
		}

		got := resolveGit("{{.Git.Branch}}")
		if got != sha[:7] {
			t.Errorf("resolveGit = %q, want %q", got, sha[:7])
		}
	})

	t.Run("short SHA", func(t *testing.T) {
		dir := t.TempDir()
		t.Chdir(dir)
		if err := exec.Command("git", "init").Run(); err != nil {
			t.Fatalf("git init: %v", err)
		}
		if err := exec.Command("git", "config", "user.email", "test@example.com").Run(); err != nil {
			t.Fatalf("git config user.email: %v", err)
		}
		if err := exec.Command("git", "config", "user.name", "Test").Run(); err != nil {
			t.Fatalf("git config user.name: %v", err)
		}
		if err := os.WriteFile("a.txt", []byte("a"), 0o644); err != nil {
			t.Fatalf("write file: %v", err)
		}
		if err := exec.Command("git", "add", "a.txt").Run(); err != nil {
			t.Fatalf("git add: %v", err)
		}
		if err := exec.Command("git", "commit", "-m", "first").Run(); err != nil {
			t.Fatalf("git commit: %v", err)
		}

		shaOut, err := exec.Command("git", "rev-parse", "--short", "HEAD").Output()
		if err != nil {
			t.Fatalf("git rev-parse --short HEAD: %v", err)
		}
		sha := strings.TrimSpace(string(shaOut))

		got := resolveGit("{{.Git.ShortSHA}}")
		if got != sha {
			t.Errorf("resolveGit = %q, want %q", got, sha)
		}
	})

	t.Run("invalid template syntax", func(t *testing.T) {
		dir := t.TempDir()
		t.Chdir(dir)
		if err := exec.Command("git", "init").Run(); err != nil {
			t.Fatalf("git init: %v", err)
		}
		if err := exec.Command("git", "config", "user.email", "test@example.com").Run(); err != nil {
			t.Fatalf("git config user.email: %v", err)
		}
		if err := exec.Command("git", "config", "user.name", "Test").Run(); err != nil {
			t.Fatalf("git config user.name: %v", err)
		}
		if err := os.WriteFile("a.txt", []byte("a"), 0o644); err != nil {
			t.Fatalf("write file: %v", err)
		}
		if err := exec.Command("git", "add", "a.txt").Run(); err != nil {
			t.Fatalf("git add: %v", err)
		}
		if err := exec.Command("git", "commit", "-m", "first").Run(); err != nil {
			t.Fatalf("git commit: %v", err)
		}

		got := resolveGit("{{.BadSyntax")
		if got != "" {
			t.Errorf("resolveGit with invalid template = %q, want empty", got)
		}
	})

	t.Run("undefined template variable", func(t *testing.T) {
		dir := t.TempDir()
		t.Chdir(dir)
		if err := exec.Command("git", "init").Run(); err != nil {
			t.Fatalf("git init: %v", err)
		}
		if err := exec.Command("git", "config", "user.email", "test@example.com").Run(); err != nil {
			t.Fatalf("git config user.email: %v", err)
		}
		if err := exec.Command("git", "config", "user.name", "Test").Run(); err != nil {
			t.Fatalf("git config user.name: %v", err)
		}
		if err := os.WriteFile("a.txt", []byte("a"), 0o644); err != nil {
			t.Fatalf("write file: %v", err)
		}
		if err := exec.Command("git", "add", "a.txt").Run(); err != nil {
			t.Fatalf("git add: %v", err)
		}
		if err := exec.Command("git", "commit", "-m", "first").Run(); err != nil {
			t.Fatalf("git commit: %v", err)
		}

		got := resolveGit("{{.Git.NonExistent}}")
		if got != "" {
			t.Errorf("resolveGit with undefined var = %q, want empty", got)
		}
	})

	t.Run("template with literal text", func(t *testing.T) {
		dir := t.TempDir()
		t.Chdir(dir)
		if err := exec.Command("git", "init").Run(); err != nil {
			t.Fatalf("git init: %v", err)
		}
		if err := exec.Command("git", "config", "user.email", "test@example.com").Run(); err != nil {
			t.Fatalf("git config user.email: %v", err)
		}
		if err := exec.Command("git", "config", "user.name", "Test").Run(); err != nil {
			t.Fatalf("git config user.name: %v", err)
		}
		if err := os.WriteFile("a.txt", []byte("a"), 0o644); err != nil {
			t.Fatalf("write file: %v", err)
		}
		if err := exec.Command("git", "add", "a.txt").Run(); err != nil {
			t.Fatalf("git add: %v", err)
		}
		if err := exec.Command("git", "commit", "-m", "first").Run(); err != nil {
			t.Fatalf("git commit: %v", err)
		}

		out, err := exec.Command("git", "branch", "--show-current").Output()
		if err != nil {
			t.Fatalf("git branch: %v", err)
		}
		branch := strings.TrimSpace(string(out))

		got := resolveGit("[{{.Git.Branch}}]")
		want := "[" + branch + "]"
		if got != want {
			t.Errorf("resolveGit = %q, want %q", got, want)
		}
	})
}

func TestGitOutputError(t *testing.T) {
	_, err := gitOutput("not-a-valid-git-command")
	if err == nil {
		t.Error("gitOutput expected error for invalid command, got nil")
	}
}

func TestResolveGitMissingBinary(t *testing.T) {
	t.Setenv("PATH", "/nonexistent")
	dir := t.TempDir()
	t.Chdir(dir)

	got := resolveGit("{{.Git.Branch}}")
	if got != "" {
		t.Errorf("resolveGit with missing git binary = %q, want empty", got)
	}
}

func TestResolveGitEmptyRepo(t *testing.T) {
	skipIfNoGit(t)
	dir := t.TempDir()
	t.Chdir(dir)
	if err := exec.Command("git", "init").Run(); err != nil {
		t.Fatalf("git init: %v", err)
	}

	out, err := exec.Command("git", "branch", "--show-current").Output()
	if err != nil {
		t.Fatalf("git branch: %v", err)
	}
	branch := strings.TrimSpace(string(out))

	got := resolveGit("{{.Git.Branch}}")
	if got != branch {
		t.Errorf("resolveGit in empty repo = %q, want %q", got, branch)
	}
}

func TestResolveGitEmptyTemplate(t *testing.T) {
	skipIfNoGit(t)
	dir := t.TempDir()
	t.Chdir(dir)
	if err := exec.Command("git", "init").Run(); err != nil {
		t.Fatalf("git init: %v", err)
	}
	if err := exec.Command("git", "config", "user.email", "test@example.com").Run(); err != nil {
		t.Fatalf("git config user.email: %v", err)
	}
	if err := exec.Command("git", "config", "user.name", "Test").Run(); err != nil {
		t.Fatalf("git config user.name: %v", err)
	}
	if err := os.WriteFile("a.txt", []byte("a"), 0o644); err != nil {
		t.Fatalf("write file: %v", err)
	}
	if err := exec.Command("git", "add", "a.txt").Run(); err != nil {
		t.Fatalf("git add: %v", err)
	}
	if err := exec.Command("git", "commit", "-m", "first").Run(); err != nil {
		t.Fatalf("git commit: %v", err)
	}

	got := resolveGit("")
	if got != "" {
		t.Errorf("resolveGit with empty template = %q, want empty", got)
	}
}
