// Package commands provides the CLI subcommands for prst.
package commands

import (
	"fmt"

	"github.com/andrewhowdencom/prst/internal/prompt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// NewCommand0 returns the no-op prst 0 command (PS0).
func NewCommand0() *cobra.Command {
	return &cobra.Command{
		Use:   "0",
		Short: "Bash prompt string 0 (PS0)",
		RunE: func(_ *cobra.Command, _ []string) error {
			return nil
		},
	}
}

// NewCommand1 returns the prst 1 command that prints a PS1 string.
func NewCommand1(v *viper.Viper, g prompt.Generator) *cobra.Command {
	return &cobra.Command{
		Use:   "1",
		Short: "Bash prompt string 1 (PS1)",
		RunE: func(cmd *cobra.Command, _ []string) error {
			noColor, _ := cmd.Flags().GetBool("no-color")
			cap := prompt.DefaultColorCapability(noColor, v)
			_, err := fmt.Fprintln(cmd.OutOrStdout(), g.Generate(cap))
			return err
		},
	}
}

// NewCommand2 returns the no-op prst 2 command (PS2).
func NewCommand2() *cobra.Command {
	return &cobra.Command{
		Use:   "2",
		Short: "Bash prompt string 2 (PS2)",
		RunE: func(_ *cobra.Command, _ []string) error {
			return nil
		},
	}
}

// NewCommand3 returns the no-op prst 3 command (PS3).
func NewCommand3() *cobra.Command {
	return &cobra.Command{
		Use:   "3",
		Short: "Bash prompt string 3 (PS3)",
		RunE: func(_ *cobra.Command, _ []string) error {
			return nil
		},
	}
}

// NewCommand4 returns the no-op prst 4 command (PS4).
func NewCommand4() *cobra.Command {
	return &cobra.Command{
		Use:   "4",
		Short: "Bash prompt string 4 (PS4)",
		RunE: func(_ *cobra.Command, _ []string) error {
			return nil
		},
	}
}

// NewVersionCommand returns the version command showing Go build info.
func NewVersionCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print the version of prst",
		RunE: func(cmd *cobra.Command, _ []string) error {
			cmd.Println("v0.0.0-dev")
			return nil
		},
	}
}

// NewCommands returns all subcommands as a slice.
func NewCommands(v *viper.Viper, g prompt.Generator) []*cobra.Command {
	return []*cobra.Command{
		NewCommand0(),
		NewCommand1(v, g),
		NewCommand2(),
		NewCommand3(),
		NewCommand4(),
		NewVersionCommand(),
	}
}
