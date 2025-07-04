package app

import (
	"context"
	"fmt"
	"github.com/urfave/cli/v3"
	"pdok-metadata-tool/pkg/repository"
)

func init() {
	command := &cli.Command{
		Name:  "hvd",
		Usage: "Used to retrieve and inspect high value dataset categories from the HVD Thesaurus.",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "url",
				Value: HvdCategoriesXmlRemote,
				Usage: "HVD Thesaurus endpoint which should contain the HVD categories as RDF format.",
			},
			&cli.StringFlag{
				Name:  "local-path",
				Value: HvdLocalPath,
				Usage: "HVD Thesaurus endpoint which should contain the HVD categories as RDF format.",
			},
		},
		Before: func(ctx context.Context, cmd *cli.Command) (context.Context, error) {
			url := cmd.String("url")
			localPath := cmd.String("local-path")
			hvdrepo := repository.NewHVDRepository(url, localPath)
			// Store hvdrepo in context
			return context.WithValue(ctx, "HVDRepoKey", hvdrepo), nil
		},
		Commands: []*cli.Command{
			{
				Name:  "download",
				Usage: "Download HVD",
				Action: func(ctx context.Context, cmd *cli.Command) error {
					fmt.Println("pmt hvd download invoked")

					// Get the HVDRepository from context
					hvdrepo, ok := ctx.Value("HVDRepoKey").(*repository.HVDRepository)
					if !ok {
						return fmt.Errorf("failed to get HVDRepository from context")
					}

					// Call the Download method
					err := hvdrepo.Download()
					if err != nil {
						return fmt.Errorf("failed to download HVD: %w", err)
					}

					return nil
				},
			},
			{
				Name:  "list",
				Usage: "List HVDs",
				Action: func(ctx context.Context, cmd *cli.Command) error {
					fmt.Println("pmt hvd list invoked")

					// Get the HVDRepository from context
					hvdrepo, ok := ctx.Value("HVDRepoKey").(*repository.HVDRepository)
					if !ok {
						return fmt.Errorf("failed to get HVDRepository from context")
					}

					// Call the GetAllHVDCategories method
					err := hvdrepo.GetAllHVDCategories()
					if err != nil {
						return fmt.Errorf("failed to get HVD categories: %w", err)
					}

					return nil
				},
			},
		},
	}
	PDOKMetadataToolCLI.Commands = append(PDOKMetadataToolCLI.Commands, command)
}
