package main

import (
	"context"
	"fmt"
	"os"
	"pdok-metadata-tool/internal/app"
)

func main() {

	ctx := context.Background()
	if err := app.PDOKMetadataToolCLI.Run(ctx, os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
