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
				Name:        "url",
				DefaultText: "eu-thesaurus-url",
				Value:       HvdCategoriesXmlRemote,
				Usage:       "HVD Thesaurus endpoint which should contain the HVD categories as RDF format.",
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
				Usage: "Downloads the RDF Thesaurus containing the HVD categories at local-path.",
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
				Usage: "Displays list of HVD categories.",
				Action: func(ctx context.Context, cmd *cli.Command) error {
					fmt.Println("pmt hvd list invoked")

					// Get the HVDRepository from context
					hvdrepo, ok := ctx.Value("HVDRepoKey").(*repository.HVDRepository)
					if !ok {
						return fmt.Errorf("failed to get HVDRepository from context")
					}

					// Call the GetAllHVDCategories method
					categories, err := hvdrepo.GetAllHVDCategories()
					if err != nil {
						return fmt.Errorf("failed to get HVD categories: %w", err)
					}

					// Print the categories
					fmt.Printf("Found %d HVD categories\n", len(categories))
					for _, category := range categories {
						fmt.Printf("ID: %s, Parent: %s, Order: %s\n", category.Id, category.Parent, category.Order)
						fmt.Printf("  Dutch: %s\n", category.LabelDutch)
						fmt.Printf("  English: %s\n", category.LabelEnglish)
						fmt.Println()
					}

					return nil
				},
			},
		},
	}
	PDOKMetadataToolCLI.Commands = append(PDOKMetadataToolCLI.Commands, command)
}
