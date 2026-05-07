package prompt

import (
	"log/slog"
	"os/exec"
	"strconv"
	"strings"
	"text/template"
)

// gitStatus holds the runtime git state exposed to the git segment template.
type gitStatus struct {
	Branch   string
	Dirty    bool
	Ahead    int
	Behind   int
	ShortSHA string
}

// resolveGit collects git repository state and renders the user-provided
// template string. It returns an empty string when not inside a git worktree,
// when git is unavailable, or when template execution fails.
func resolveGit(templateStr string) string {
	// Check if we're inside a git worktree.
	if out, err := gitOutput("rev-parse", "--is-inside-work-tree"); err != nil || out != "true" {
		return ""
	}

	gs := gitStatus{}

	// Branch name.
	gs.Branch = gitBranch()

	// Dirty status.
	gs.Dirty = gitDirty()

	// Ahead/behind.
	gs.Ahead, gs.Behind = gitAheadBehind()

	// Short SHA.
	gs.ShortSHA = gitShortSHA()

	// If we have no branch and no SHA, something went wrong; bail out.
	if gs.Branch == "" && gs.ShortSHA == "" {
		return ""
	}

	tmpl, err := template.New("git").Parse(templateStr)
	if err != nil {
		slog.Warn("invalid git segment template", "error", err, "template", templateStr)
		return ""
	}

	data := map[string]interface{}{
		"Git": gs,
	}

	var b strings.Builder
	if err := tmpl.Execute(&b, data); err != nil {
		slog.Warn("failed to execute git segment template", "error", err, "template", templateStr)
		return ""
	}

	return b.String()
}

// gitOutput runs a git command and returns trimmed stdout. Errors are returned
// to the caller.
func gitOutput(args ...string) (string, error) {
	cmd := exec.Command("git", args...)
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}

// gitBranch returns the current branch name. On detached HEAD, it falls back
// to the short commit SHA.
func gitBranch() string {
	if out, err := gitOutput("symbolic-ref", "--short", "HEAD"); err == nil {
		return out
	}
	return gitShortSHA()
}

// gitDirty reports whether the working tree or index has uncommitted changes.
func gitDirty() bool {
	out, err := gitOutput("status", "--porcelain")
	if err != nil {
		return false
	}
	return isDirtyFromPorcelain(out)
}

// isDirtyFromPorcelain determines whether git status --porcelain indicates
// uncommitted changes (including staged, unstaged, or untracked files).
func isDirtyFromPorcelain(output string) bool {
	return strings.TrimSpace(output) != ""
}

// gitAheadBehind returns the number of commits the current branch is ahead
// and behind its upstream. If no upstream is configured, both values are 0.
func gitAheadBehind() (ahead int, behind int) {
	out, err := gitOutput("rev-list", "--left-right", "--count", "HEAD...@{upstream}")
	if err != nil {
		return 0, 0
	}
	return parseAheadBehind(out)
}

// parseAheadBehind parses the tab-separated output of
// git rev-list --left-right --count into ahead and behind counts.
func parseAheadBehind(output string) (ahead int, behind int) {
	parts := strings.Split(strings.TrimSpace(output), "\t")
	if len(parts) != 2 {
		return 0, 0
	}
	ahead, _ = strconv.Atoi(parts[0])
	behind, _ = strconv.Atoi(parts[1])
	return
}

// gitShortSHA returns the short commit hash of HEAD.
func gitShortSHA() string {
	out, err := gitOutput("rev-parse", "--short", "HEAD")
	if err != nil {
		return ""
	}
	return out
}
