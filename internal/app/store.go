package app

import (
	"context"
	"fmt"
	"github.com/urfave/cli/v3"
)

func init() {
	command := &cli.Command{
		Name:  "metadata",
		Usage: "The metadata toolchain is used to generate service metadata",
		Commands: []*cli.Command{
			{
				Name:  "bump-revision-date",
				Usage: "Bumps revision date to today or provided date in NGR metadata record.",
				Arguments: []cli.Argument{
					&cli.StringArg{
						Name:      "uuid",
						UsageText: "Metadata uuid of metadata record which needs to be bumped in NGR.",
					},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					fmt.Printf("Hello %q", cmd.Args().Get(0))
					return nil
				},
			},
		},
	}
	PDOKMetadataToolCLI.Commands = append(PDOKMetadataToolCLI.Commands, command)
}
