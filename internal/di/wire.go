//go:build wireinject

// Package di provides compile-time dependency injection via Google Wire.
package di

import (
	"github.com/andrewhowdencom/prst/internal/app"
	"github.com/andrewhowdencom/prst/internal/commands"
	"github.com/andrewhowdencom/prst/internal/configuration"
	"github.com/andrewhowdencom/prst/internal/prompt"
	"github.com/google/wire"
	"github.com/spf13/cobra"
)

// NewApplication constructs the cobra.Command application graph.
func NewApplication() (*cobra.Command, error) {
	wire.Build(
		configuration.NewViper,
		prompt.NewPS1Config,
		prompt.NewPS1Generator,
		wire.Bind(new(prompt.Generator), new(*prompt.PS1Generator)),
		commands.NewCommands,
		app.NewRootCommand,
	)
	return nil, nil
}
