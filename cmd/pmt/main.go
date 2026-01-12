// Package main provides a main function for running the PDOK Metadata tool CLI.
package main

import (
	"context"
	"fmt"
	"os"

	"github.com/pdok/pdok-metadata-tool/v2/internal/app"
)

func main() {
	ctx := context.Background()
	if err := app.PDOKMetadataToolCLI.Run(ctx, os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
