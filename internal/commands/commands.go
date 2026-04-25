// Package commands provides the CLI subcommands for prst.
package commands

import (
	"github.com/andrewhowdencom/prst/internal/prompt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

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
		NewPromptCommand(v, g),
		NewInitCommand(),
		NewInstallCommand(),
		NewVersionCommand(),
	}
}
