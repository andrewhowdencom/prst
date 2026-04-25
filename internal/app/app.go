// Package app provides the Cobra root command and global flag configuration.
package app

import (
	"log/slog"
	"os"

	"github.com/spf13/cobra"
)

// NewRootCommand creates and returns the root cobra command for prst.
func NewRootCommand(cmds []*cobra.Command) *cobra.Command {
	root := &cobra.Command{
		Use:   "prst",
		Short: "A tool for managing Bash prompt strings (PS0–PS4)",
		PersistentPreRunE: func(cmd *cobra.Command, _ []string) error {
			level, err := cmd.Flags().GetString("log-level")
			if err != nil {
				return err
			}

			var lvl slog.Level
			if err := lvl.UnmarshalText([]byte(level)); err != nil {
				return err
			}

			logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: lvl}))
			slog.SetDefault(logger)

			return nil
		},
	}

	root.PersistentFlags().String("log-level", "info", "Set the logging level (debug, info, warn, error)")
	root.PersistentFlags().Bool("no-color", false, "Disable colored output")
	root.AddCommand(cmds...)

	return root
}
