package app

import (
	"context"
	"fmt"
	"github.com/urfave/cli/v3"
	"pdok-metadata-tool/internal/common"
	"pdok-metadata-tool/pkg/repository"
	"strings"
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
			repo := repository.NewHVDRepository(url, localPath)
			return context.WithValue(ctx, "HVDRepoKey", repo), nil
		},
		Commands: []*cli.Command{
			{
				Name:  "download",
				Usage: "Downloads the RDF Thesaurus containing the HVD categories at local-path.",
				Action: func(ctx context.Context, cmd *cli.Command) (err error) {
					fmt.Println("Downloading HVD Thesaurus from " + cmd.String("url") + " to " + cmd.String("local-path"))

					repo, ok := ctx.Value("HVDRepoKey").(*repository.HVDRepository)
					if !ok {
						return fmt.Errorf("failed to get HVDRepository from context")
					}

					err = repo.Download()
					return err
				},
			},
			{
				Name:  "list",
				Usage: "Displays list of HVD categories.",
				Action: func(ctx context.Context, cmd *cli.Command) error {
					fmt.Println("pmt hvd list invoked")

					hvdrepo, ok := ctx.Value("HVDRepoKey").(*repository.HVDRepository)
					if !ok {
						return fmt.Errorf("failed to get HVDRepository from context")
					}

					categories, err := hvdrepo.GetAllHVDCategories()
					if err != nil {
						return fmt.Errorf("failed to get HVD categories: %w", err)
					}

					fmt.Printf("Found %d HVD categories\n", len(categories))

					// Print table header
					fmt.Printf("%-10s %-10s %-6s %-30s %-30s\n", "ID", "PARENT", "ORDER", "DUTCH", "ENGLISH")
					fmt.Println(strings.Repeat("-", 90))

					// Print each category as a row
					for _, category := range categories {
						fmt.Printf("%-10s %-10s %-6s %-30s %-30s\n",
							category.Id,
							category.Parent,
							category.Order,
							common.TruncateString(category.LabelDutch, 30),
							common.TruncateString(category.LabelEnglish, 30))
					}

					return nil
				},
			},
		},
	}
	PDOKMetadataToolCLI.Commands = append(PDOKMetadataToolCLI.Commands, command)
}
