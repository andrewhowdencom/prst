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
	cmd := &cobra.Command{
		Use:   "prompt [0|1|2|3|4]",
		Short: "Print a shell prompt string",
		Long: `Print a prompt string for the given PS level (0–4).

Only PS1 is currently implemented; PS0, PS2, PS3, and PS4 are reserved
for future expansion and currently print nothing.`,
		Example: `  # Raw ANSI (no shell wrapping)
  prst prompt --color=always 1

  # Wrapped for Bash PS1 (non-printing SOH/STX markers)
  prst prompt --color=always --shell=bash 1

  # Wrapped for Zsh PS1 (%{...%} markers)
  prst prompt --color=always --shell=zsh 1`,
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
				colorFlag, _ := cmd.Flags().GetString("color")
				cap := prompt.DefaultColorCapability(colorFlag, v)
				raw := g.Generate(cap)
				shell, _ := cmd.Flags().GetString("shell")
				if shell != "" {
					raw = prompt.FormatForShell(raw, shell)
				}
				_, err := fmt.Fprintln(cmd.OutOrStdout(), raw)
				return err
			}
			// 0, 2, 3, 4 are no-ops for now.
			return nil
		},
	}

	cmd.Flags().String("shell", "", "Target shell for non-printing sequence wrapping (bash, zsh)")

	return cmd
}
