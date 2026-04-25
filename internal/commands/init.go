package commands

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/andrewhowdencom/prst/internal/shell"
	"github.com/spf13/cobra"
)

// NewInitCommand returns the prst init command.
func NewInitCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "init <shell> [0] [1] [2] [3] [4]",
		Short: "Print shell-specific initialization code",
		Long: `Outputs a shell script that sets up prst for the given shell.
Pass one or more prompt numbers (0-4) to configure which PS variables
are managed. If no numbers are given, defaults to 1.`,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return fmt.Errorf("requires a shell argument: %s", strings.Join(shell.ValidShells(), ", "))
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			sh, err := shell.ParseShell(args[0])
			if err != nil {
				return err
			}

			numbers, err := parsePromptNumbers(args[1:])
			if err != nil {
				return err
			}
			if len(numbers) == 0 {
				numbers = []int{1}
			}

			script := sh.InitScript(numbers)
			_, err = fmt.Fprint(cmd.OutOrStdout(), script)
			return err
		},
	}
}

func parsePromptNumbers(args []string) ([]int, error) {
	var numbers []int
	for _, arg := range args {
		n, err := strconv.Atoi(arg)
		if err != nil || n < 0 || n > 4 {
			return nil, fmt.Errorf("invalid prompt number %q: must be 0-4", arg)
		}
		numbers = append(numbers, n)
	}
	return numbers, nil
}
