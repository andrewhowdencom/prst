// Package main is the entry point for the prst CLI.
package main

import (
	"context"
	"fmt"
	"os"

	"github.com/andrewhowdencom/prst/internal/di"
)

func main() {
	if err := di.NewApplication().ExecuteContext(context.Background()); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}
