package app

import (
	"context"
	"encoding/csv"
	"fmt"
	"os"
	"pdok-metadata-tool/internal/common"
	"pdok-metadata-tool/pkg/model/hvd"
	"pdok-metadata-tool/pkg/repository"
	"strings"

	"github.com/urfave/cli/v3"
)

func init() {
	command := &cli.Command{
		Name:  "hvd",
		Usage: "Used to retrieve and inspect high value dataset categories from the HVD Thesaurus.",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "url",
				DefaultText: "eu-thesaurus-url",
				Value:       hvd.HvdEndpoint,
				Usage:       "HVD Thesaurus endpoint which should contain the HVD categories as RDF format.",
			},
			&cli.StringFlag{
				Name:  "local-path",
				Value: HvdLocalRDFPath,
				Usage: "Local path where the HVD Thesaurus is cached.",
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

					_, err = repo.Download()
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
			{
				Name:  "csv",
				Usage: "Exports HVD categories to a CSV file.",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "o",
						Value: HvdLocalCSVPath,
						Usage: "Output file path for the CSV file.",
					},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					fmt.Println("pmt hvd csv invoked")

					hvdrepo, ok := ctx.Value("HVDRepoKey").(*repository.HVDRepository)
					if !ok {
						return fmt.Errorf("failed to get HVDRepository from context")
					}

					categories, err := hvdrepo.GetAllHVDCategories()
					if err != nil {
						return fmt.Errorf("failed to get HVD categories: %w", err)
					}

					outputPath := cmd.String("o")
					fmt.Printf("Exporting %d HVD categories to %s\n", len(categories), outputPath)

					// Create the CSV file
					file, err := os.Create(outputPath)
					if err != nil {
						return fmt.Errorf("failed to create CSV file: %w", err)
					}
					defer file.Close()

					writer := csv.NewWriter(file)
					defer writer.Flush()

					// Write header
					header := []string{"ID", "Parent", "Order", "LabelDutch", "LabelEnglish"}
					if err := writer.Write(header); err != nil {
						return fmt.Errorf("failed to write CSV header: %w", err)
					}

					// Write data
					for _, category := range categories {
						row := []string{
							category.Id,
							category.Parent,
							category.Order,
							category.LabelDutch,
							category.LabelEnglish,
						}
						if err := writer.Write(row); err != nil {
							return fmt.Errorf("failed to write CSV row: %w", err)
						}
					}

					fmt.Printf("Successfully exported HVD categories to %s\n", outputPath)
					return nil
				},
			},
		},
	}
	PDOKMetadataToolCLI.Commands = append(PDOKMetadataToolCLI.Commands, command)
}
