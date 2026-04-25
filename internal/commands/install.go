package commands

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/andrewhowdencom/prst/internal/shell"
	"github.com/spf13/cobra"
)

const (
	prstBlockStart = "# >>> prst init >>>"
	prstBlockEnd   = "# <<< prst init <<<"
)

// NewInstallCommand returns the prst install command.
func NewInstallCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "install [0] [1] [2] [3] [4]",
		Short: "Install prst into your shell configuration",
		Long: `Detects your shell and appends prst initialization to the
appropriate runtime configuration file (e.g. ~/.bashrc, ~/.zshrc).
Pass one or more prompt numbers (0-4) to configure which PS variables
are managed. If no numbers are given, defaults to 1.`,
		RunE: runInstall,
	}

	cmd.Flags().String("shell", "", fmt.Sprintf("Target shell (default: auto-detect; supported: %s)", strings.Join(shell.ValidShells(), ", ")))
	cmd.Flags().Bool("dry-run", false, "Print what would be written without modifying files")
	cmd.Flags().Bool("remove", false, "Remove prst initialization from the shell configuration")

	return cmd
}

func runInstall(cmd *cobra.Command, args []string) error {
	shFlag, _ := cmd.Flags().GetString("shell")
	dryRun, _ := cmd.Flags().GetBool("dry-run")
	remove, _ := cmd.Flags().GetBool("remove")

	var sh shell.Shell
	if shFlag != "" {
		var err error
		sh, err = shell.ParseShell(shFlag)
		if err != nil {
			return err
		}
	} else {
		sh = shell.Detect()
	}

	numbers, err := parsePromptNumbers(args)
	if err != nil {
		return err
	}
	if len(numbers) == 0 {
		numbers = []int{1}
	}

	rcFile := sh.RCFile()
	if rcFile == "" {
		return fmt.Errorf("could not determine rc file for shell %q", sh)
	}

	if dryRun {
		fmt.Fprintf(cmd.OutOrStdout(), "# Dry-run: would modify %s\n", rcFile)
		if remove {
			fmt.Fprintln(cmd.OutOrStdout(), "# Would remove prst init block")
		} else {
			fmt.Fprintln(cmd.OutOrStdout(), "# Would append:")
			fmt.Fprintln(cmd.OutOrStdout(), generateBlock(sh, numbers))
		}
		return nil
	}

	if remove {
		return removeBlock(rcFile)
	}

	return appendBlock(rcFile, sh, numbers)
}

func generateBlock(sh shell.Shell, numbers []int) string {
	sort.Ints(numbers)
	var nums []string
	for _, n := range numbers {
		nums = append(nums, strconv.Itoa(n))
	}

	var b strings.Builder
	b.WriteString(prstBlockStart + "\n")
	b.WriteString(fmt.Sprintf(`eval "$(prst init %s %s)"`, sh, strings.Join(nums, " ")))
	b.WriteString("\n")
	b.WriteString(prstBlockEnd + "\n")
	return b.String()
}

func appendBlock(rcFile string, sh shell.Shell, numbers []int) error {
	block := generateBlock(sh, numbers)

	// Ensure directory exists.
	dir := filepath.Dir(rcFile)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("creating rc directory: %w", err)
	}

	// Read existing content.
	content, err := os.ReadFile(rcFile)
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("reading rc file: %w", err)
	}

	// Remove any existing block to ensure idempotency.
	content = stripBlock(content)

	// Append block with a leading newline for cleanliness.
	if len(content) > 0 && !bytes.HasSuffix(content, []byte("\n")) {
		content = append(content, '\n')
	}
	content = append(content, '\n')
	content = append(content, []byte(block)...)

	if err := os.WriteFile(rcFile, content, 0o600); err != nil {
		return fmt.Errorf("writing rc file: %w", err)
	}

	return nil
}

func removeBlock(rcFile string) error {
	content, err := os.ReadFile(rcFile)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return fmt.Errorf("reading rc file: %w", err)
	}

	content = stripBlock(content)

	// Clean up trailing newlines.
	content = bytes.TrimRight(content, "\n")
	content = append(content, '\n')

	if err := os.WriteFile(rcFile, content, 0o600); err != nil {
		return fmt.Errorf("writing rc file: %w", err)
	}

	return nil
}

func stripBlock(content []byte) []byte {
	scanner := bufio.NewScanner(bytes.NewReader(content))
	var lines []string
	inBlock := false
	for scanner.Scan() {
		line := scanner.Text()
		if strings.TrimSpace(line) == prstBlockStart {
			inBlock = true
			continue
		}
		if strings.TrimSpace(line) == prstBlockEnd {
			inBlock = false
			continue
		}
		if !inBlock {
			lines = append(lines, line)
		}
	}

	// Rejoin with newlines and ensure trailing newline.
	var b strings.Builder
	for i, line := range lines {
		if i > 0 {
			b.WriteByte('\n')
		}
		b.WriteString(line)
	}
	result := []byte(b.String())
	if len(result) > 0 && !bytes.HasSuffix(result, []byte("\n")) {
		result = append(result, '\n')
	}
	return result
}
