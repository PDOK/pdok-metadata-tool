package app

import (
	"context"
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/pdok/pdok-metadata-tool/v2/internal/common"
	"github.com/pdok/pdok-metadata-tool/v2/pkg/model/inspire"
	"github.com/pdok/pdok-metadata-tool/v2/pkg/repository"
	"github.com/urfave/cli/v3"
)

const theme = "theme"
const layer = "layer"

// errSpecifyRegisterKind returns a unified error message including available kinds.
func errSpecifyRegisterKind() error {
	kinds := make([]string, 0, len(inspire.InspireRegisterKinds))
	for _, k := range inspire.InspireRegisterKinds {
		kinds = append(kinds, string(k))
	}

	return fmt.Errorf(
		"please specify a register kind. Available kinds: %s",
		strings.Join(kinds, ", "),
	)
}

func init() {
	command := &cli.Command{
		Name:  "inspire",
		Usage: "The metadata toolchain is used to generate service metadata.",
		Commands: []*cli.Command{
			getInspireListCommand(),
			getInspireCSVCommand(),
		},
	}
	PDOKMetadataToolCLI.Commands = append(PDOKMetadataToolCLI.Commands, command)
}

func getInspireListCommand() *cli.Command {
	return &cli.Command{
		Name:  "list",
		Usage: "List inspire themes or layers. Usage: pmt inspire list <theme|layer>",
		Action: func(_ context.Context, cmd *cli.Command) error {
			if cmd.NArg() == 0 {
				return errSpecifyRegisterKind()
			}

			if cmd.Args().First() != theme && cmd.Args().First() != layer {
				return errSpecifyRegisterKind()
			}

			fmt.Println("list inspire: ", cmd.Args().First())

			repo := repository.NewInspireRepository(common.InspireLocalPath)

			if cmd.Args().First() == theme {
				themes, err := repo.GetThemes()
				if err != nil {
					return err
				}

				for _, theme := range themes {
					fmt.Println(theme.ID + "-" + theme.LabelDutch)
				}
			}

			if cmd.Args().First() == layer {
				layers, err := repo.GetLayers()
				if err != nil {
					return err
				}

				for _, layer := range layers {
					fmt.Println(layer.ID + "-" + layer.LabelDutch)
				}
			}

			return nil
		},
		ShellComplete: func(_ context.Context, cmd *cli.Command) {
			// This will complete if no args are passed
			if cmd.NArg() > 0 {
				return
			}

			for _, t := range inspire.InspireRegisterKinds {
				fmt.Println(t)
			}
		},
	}
}

//nolint:gocognit
func getInspireCSVCommand() *cli.Command {
	return &cli.Command{
		Name:  "csv",
		Usage: "Exports inspire themes or layers to a CSV file. Usage: pmt inspire csv <theme|layer>",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "o",
				Value: "",
				Usage: "Output file path for the CSV file.",
			},
		},
		Action: func(_ context.Context, cmd *cli.Command) error {
			if cmd.NArg() == 0 {
				return errSpecifyRegisterKind()
			}

			choice := cmd.Args().First()

			if choice != theme && choice != layer {
				return errSpecifyRegisterKind()
			}

			fmt.Println("csv inspire: ", choice)

			repo := repository.NewInspireRepository(common.InspireLocalPath)
			outputPath := cmd.String("o")

			if outputPath == "" {
				outputPath = filepath.Join(common.CachePath, choice+".csv")
			}

			// Create the CSV file
			//nolint:gosec
			file, err := os.Create(outputPath)
			if err != nil {
				return fmt.Errorf("failed to create CSV file: %w", err)
			}
			defer common.SafeClose(file)

			writer := csv.NewWriter(file)
			defer writer.Flush()

			if choice == theme {
				themes, err := repo.GetThemes()
				if err != nil {
					return err
				}

				fmt.Printf("Exporting %d INSPIRE themes to %s\n", len(themes), outputPath)

				// Write header
				header := []string{"ID", "Order", "LabelDutch", "LabelEnglish", "URL"}
				if err := writer.Write(header); err != nil {
					return fmt.Errorf("failed to write CSV header: %w", err)
				}

				// Write data
				for _, theme := range themes {
					row := []string{
						theme.ID,
						strconv.Itoa(theme.Order),
						theme.LabelDutch,
						theme.LabelEnglish,
						theme.URL,
					}
					if err := writer.Write(row); err != nil {
						return fmt.Errorf("failed to write CSV row: %w", err)
					}
				}
			}

			if choice == layer {
				layers, err := repo.GetLayers()
				if err != nil {
					return err
				}

				fmt.Printf("Exporting %d INSPIRE layers to %s\n", len(layers), outputPath)

				// Write header
				header := []string{"ID", "LabelDutch", "LabelEnglish"}
				if err := writer.Write(header); err != nil {
					return fmt.Errorf("failed to write CSV header: %w", err)
				}

				// Write data
				for _, layer := range layers {
					row := []string{
						layer.ID,
						layer.LabelDutch,
						layer.LabelEnglish,
					}
					if err := writer.Write(row); err != nil {
						return fmt.Errorf("failed to write CSV row: %w", err)
					}
				}
			}

			fmt.Printf("Successfully exported INSPIRE %s to %s\n", choice, outputPath)

			return nil
		},
		ShellComplete: func(_ context.Context, cmd *cli.Command) {
			// This will complete if no args are passed
			if cmd.NArg() > 0 {
				return
			}

			for _, t := range inspire.InspireRegisterKinds {
				fmt.Println(t)
			}
		},
	}
}
