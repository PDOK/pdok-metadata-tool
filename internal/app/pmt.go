package app

import (
	"context"
	"fmt"
	"github.com/urfave/cli/v3"
	"pdok-metadata-tool/pkg/repository"
)

var PDOKMetadataToolCLI = &cli.Command{
	Name:  "pmt",
	Usage: "PDOK Metadata Tool - This tool is set up to handle various metadata related tasks.",
	Commands: []*cli.Command{
		{
			Name:  "metadata",
			Usage: "The metadata toolchain is used to generate service metadata",
			Commands: []*cli.Command{
				{
					Name:  "generate",
					Usage: "Generates service metadata in \"Nederlands profiel ISO 19119\" version 2.1.0.",
					Before: func(ctx context.Context, cmd *cli.Command) (context.Context, error) {
						hvdrepo := repository.NewHVDRepository("your-config")
						// Store hvdrepo in context
						return context.WithValue(ctx, "HVDRepoKey", hvdrepo), nil
					},
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
		},
	},
}
