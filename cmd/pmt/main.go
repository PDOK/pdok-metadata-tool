package main

import (
	"context"
	"fmt"
	"github.com/pdok/pdok-metadata-tool/internal/app"
	"os"
)

func main() {

	ctx := context.Background()
	if err := app.PDOKMetadataToolCLI.Run(ctx, os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
