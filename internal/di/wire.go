//go:build wireinject

// Package di provides compile-time dependency injection via Google Wire.
package di

import (
	"github.com/andrewhowdencom/prst/internal/app"
	"github.com/andrewhowdencom/prst/internal/commands"
	"github.com/google/wire"
	"github.com/spf13/cobra"
)

// NewApplication constructs the cobra.Command application graph.
func NewApplication() *cobra.Command {
	wire.Build(
		commands.NewCommands,
		app.NewRootCommand,
	)
	return nil
}
