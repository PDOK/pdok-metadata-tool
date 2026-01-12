package app

import (
	"context"
	"encoding/csv"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/pdok/pdok-metadata-tool/v2/internal/common"
	"github.com/pdok/pdok-metadata-tool/v2/pkg/model/hvd"
	"github.com/pdok/pdok-metadata-tool/v2/pkg/repository"

	"github.com/urfave/cli/v3"
)

type contextKey string

const hvdRepoKey contextKey = "HVDRepoKey"

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
				Value: common.HvdLocalRDFPath,
				Usage: "Local path where the HVD Thesaurus is cached.",
			},
		},
		Before: func(ctx context.Context, cmd *cli.Command) (context.Context, error) {
			url := cmd.String("url")
			localPath := cmd.String("local-path")
			repo := repository.NewHVDRepository(url, localPath)

			return context.WithValue(ctx, hvdRepoKey, repo), nil
		},
		Commands: []*cli.Command{
			getHvdDownloadCommand(),
			getHvdListCommand(),
			getHvdCSVCommand(),
		},
	}
	PDOKMetadataToolCLI.Commands = append(PDOKMetadataToolCLI.Commands, command)
}

func getHvdDownloadCommand() *cli.Command {
	return &cli.Command{

		Name:  "download",
		Usage: "Downloads the RDF Thesaurus containing the HVD categories at local-path.",
		Action: func(ctx context.Context, cmd *cli.Command) (err error) {
			fmt.Println(
				"Downloading HVD Thesaurus from " + cmd.String(
					"url",
				) + " to " + cmd.String(
					"local-path",
				),
			)

			repo, ok := ctx.Value(hvdRepoKey).(*repository.HVDRepository)
			if !ok {
				return errors.New("failed to get HVDRepository from context")
			}

			_, err = repo.Download()

			return err
		},
	}
}

func getHvdListCommand() *cli.Command {
	return &cli.Command{
		Name:  "list",
		Usage: "Displays list of HVD categories.",
		Action: func(ctx context.Context, _ *cli.Command) error {
			fmt.Println("pmt hvd list invoked")

			hvdRepo, ok := ctx.Value(hvdRepoKey).(*repository.HVDRepository)
			if !ok {
				return errors.New("failed to get HVDRepository from context")
			}

			categories, err := hvdRepo.GetAllHVDCategories()
			if err != nil {
				return fmt.Errorf("failed to get HVD categories: %w", err)
			}

			fmt.Printf("Found %d HVD categories\n", len(categories))

			// Print table header
			fmt.Printf(
				"%-10s %-10s %-6s %-30s %-30s\n",
				"ID",
				"PARENT",
				"ORDER",
				"DUTCH",
				"ENGLISH",
			)

			repeatCount := 90
			fmt.Println(strings.Repeat("-", repeatCount))

			// Print each category as a row

			maxLength := 30
			for _, category := range categories {
				fmt.Printf("%-10s %-10s %-6s %-30s %-30s\n",
					category.ID,
					category.Parent,
					category.Order,
					common.TruncateString(category.LabelDutch, maxLength),
					common.TruncateString(category.LabelEnglish, maxLength))
			}

			return nil
		},
	}
}

func getHvdCSVCommand() *cli.Command {
	return &cli.Command{
		Name:  "csv",
		Usage: "Exports HVD categories to a CSV file.",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "o",
				Value: common.HvdLocalCSVPath,
				Usage: "Output file path for the CSV file.",
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			fmt.Println("pmt hvd csv invoked")

			hvdRepo, ok := ctx.Value(hvdRepoKey).(*repository.HVDRepository)
			if !ok {
				return errors.New("failed to get HVDRepository from context")
			}

			categories, err := hvdRepo.GetAllHVDCategories()
			if err != nil {
				return fmt.Errorf("failed to get HVD categories: %w", err)
			}

			outputPath := cmd.String("o")
			fmt.Printf("Exporting %d HVD categories to %s\n", len(categories), outputPath)

			// Create the CSV file
			//nolint:gosec
			file, err := os.Create(outputPath)
			if err != nil {
				return fmt.Errorf("failed to create CSV file: %w", err)
			}
			defer common.SafeClose(file)

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
					category.ID,
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
	}
}
