// Package main is the entry point for the prst CLI.
package main

import (
	"context"
	"fmt"
	"os"

	"github.com/andrewhowdencom/prst/internal/di"
)

func main() {
	cmd, err := di.NewApplication()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	if err := cmd.ExecuteContext(context.Background()); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}
