package app

import (
	"context"
	"fmt"
	"github.com/urfave/cli/v3"
)

func init() {
	command := &cli.Command{
		Name:  "generate",
		Usage: "The metadata toolchain is used to generate service metadata",
		Commands: []*cli.Command{
			{
				Name:  "service",
				Usage: "Generates service metadata in \"Nederlands profiel ISO 19119\" version 2.1.0.",
				Arguments: []cli.Argument{
					&cli.StringArg{
						Name:      "input_file_service_specifics",
						UsageText: "",
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
				Action: func(ctx context.Context, cmd *cli.Command) error {
					fmt.Printf("Hello %q", cmd.Args().Get(0))
					return nil
				},
			},
		},
	}
	PDOKMetadataToolCLI.Commands = append(PDOKMetadataToolCLI.Commands, command)
}
