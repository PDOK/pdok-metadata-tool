package app

import (
	"context"
	"fmt"
	"github.com/urfave/cli/v3"
)

func init() {
	command := &cli.Command{
		Name:  "inspire",
		Usage: "The metadata toolchain is used to generate service metadata",
		Commands: []*cli.Command{
			{
				Name:  "list",
				Usage: "Bumps revision date to today or provided date in NGR metadata record.",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "k",
						Value: "theme",
						Usage: "Inspire resource kind, choose between (theme, layer)",
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
