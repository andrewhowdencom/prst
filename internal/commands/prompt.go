package commands

import (
	"fmt"
	"strconv"

	"github.com/andrewhowdencom/prst/internal/prompt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// NewPromptCommand returns the prst prompt command.
func NewPromptCommand(v *viper.Viper, g prompt.Generator) *cobra.Command {
	return &cobra.Command{
		Use:   "prompt [0|1|2|3|4]",
		Short: "Print a shell prompt string",
		Long: `Print a prompt string for the given PS level (0–4).

Only PS1 is currently implemented; PS0, PS2, PS3, and PS4 are reserved
for future expansion and currently print nothing.`,
		Args: func(_ *cobra.Command, args []string) error {
			if len(args) != 1 {
				return fmt.Errorf("requires exactly one prompt number (0-4)")
			}
			n, err := strconv.Atoi(args[0])
			if err != nil || n < 0 || n > 4 {
				return fmt.Errorf("invalid prompt number %q: must be 0-4", args[0])
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			// Args were validated above, so strconv.Atoi cannot fail here.
			n, _ := strconv.Atoi(args[0])
			if n == 1 {
				noColor, _ := cmd.Flags().GetBool("no-color")
				cap := prompt.DefaultColorCapability(noColor, v)
				_, err := fmt.Fprintln(cmd.OutOrStdout(), g.Generate(cap))
				return err
			}
			// 0, 2, 3, 4 are no-ops for now.
			return nil
		},
	}
}
