package app

import (
	"context"
	"fmt"
	"github.com/urfave/cli/v3"
)

func init() {
	command := &cli.Command{
		Name:  "inspire",
		Usage: "The metadata toolchain is used to generate service metadata.",
		Commands: []*cli.Command{
			{
				Name:  "list",
				Usage: "List inspire themes or layers. Usage: pmt inspire list <themes|layers>",
				Action: func(ctx context.Context, cmd *cli.Command) error {
					fmt.Println("list inspire: ", cmd.Args().First())
					return nil
				},
				ShellComplete: func(ctx context.Context, cmd *cli.Command) {
					// This will complete if no args are passed
					if cmd.NArg() > 0 {
						return
					}
					for _, t := range GetInspireSourceNames() {
						fmt.Println(t)
					}
				},
			},
		},
	}
	PDOKMetadataToolCLI.Commands = append(PDOKMetadataToolCLI.Commands, command)
}
