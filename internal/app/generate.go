package app

import (
	"context"
	"fmt"
	"github.com/urfave/cli/v3"
)

func init() {
	command := &cli.Command{
		Name:  "generate",
		Usage: "Used to generate metadata records.",
		Commands: []*cli.Command{
			{
				Name:  "service",
				Usage: "Generates service metadata in \"Nederlands profiel ISO 19119\" version 2.1.0.",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "input_file_service_specifics",
						Required: true,
						Usage:    "Path to input file containing service specifics in json, yml or yaml format. See config-example for an example of the input file.",
					},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					fmt.Printf("Hello %q", cmd.Args().Get(0))
					return nil
				},
			},
			{
				Name:  "config-example",
				Usage: "Shows example of <input_file_service_specifics> for users that are not familiar with the service specifics.",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "o",
						Required: true,
						Usage:    "Output file in json, yml or yaml format.",
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
