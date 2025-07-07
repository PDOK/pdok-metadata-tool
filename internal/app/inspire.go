package app

import (
	"context"
	"fmt"
	"github.com/urfave/cli/v3"
	"pdok-metadata-tool/pkg/model/inspire"
	"pdok-metadata-tool/pkg/repository"
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
					if cmd.NArg() == 0 {
						return fmt.Errorf("please specify a register kind")
					}

					if cmd.Args().First() != "theme" && cmd.Args().First() != "layer" {
						return fmt.Errorf("please specify a register kind")
					}

					fmt.Println("list inspire: ", cmd.Args().First())

					repo := repository.NewInspireRepository(InspireLocalPath)

					if cmd.Args().First() == "theme" {
						themes, err := repo.GetThemes()
						if err != nil {
							return err
						}

						for _, theme := range themes {
							fmt.Println(theme.Id + "-" + theme.LabelDutch)
						}
					}

					if cmd.Args().First() == "layer" {
						layers, err := repo.GetLayers()
						if err != nil {
							return err
						}

						for _, layer := range layers {
							fmt.Println(layer.Id + "-" + layer.LabelDutch)
						}
					}

					return nil
				},
				ShellComplete: func(ctx context.Context, cmd *cli.Command) {
					// This will complete if no args are passed
					if cmd.NArg() > 0 {
						return
					}
					for _, t := range inspire.InspireRegisterKinds {
						fmt.Println(t)
					}
				},
			},
		},
	}
	PDOKMetadataToolCLI.Commands = append(PDOKMetadataToolCLI.Commands, command)
}
