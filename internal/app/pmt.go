package app

import (
	"github.com/urfave/cli/v3"
)

// PDOKMetadataToolCLI contains the logic of the PDOK Metadata Tool CLI.
var PDOKMetadataToolCLI = &cli.Command{
	Name:                  "pmt",
	Usage:                 "PDOK Metadata Tool - This tool is set up to handle various metadata related tasks.",
	EnableShellCompletion: true,
}
