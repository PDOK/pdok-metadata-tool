// Package app provides logic for the PMT CLI.
package app

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/pdok/pdok-metadata-tool/internal/common"
	"github.com/pdok/pdok-metadata-tool/pkg/generator"
	"github.com/urfave/cli/v3"
)

func init() {
	command := &cli.Command{
		Name:  "generate",
		Usage: "Used to generate metadata records.",
		Commands: []*cli.Command{
			getGenerateServiceCommand(),
			getGenerateConfigExampleCommand(),
		},
	}
	PDOKMetadataToolCLI.Commands = append(PDOKMetadataToolCLI.Commands, command)
}

func getGenerateServiceCommand() *cli.Command {
	return &cli.Command{
		Name:  "service",
		Usage: "Generates service metadata in \"Nederlands profiel ISO 19119\" version 2.1.0.",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "input_file_service_specifics",
				Required: true,
				Usage:    "Path to input file containing service specifics in json, yml or yaml format. See config-example for an example of the input file.",
			},
			&cli.StringFlag{
				Name:     "output_dir",
				Required: false,
				Usage:    "Location used to store service metadata as xml. If omitted the current working directory is used.",
			},
		},
		Action: func(_ context.Context, cmd *cli.Command) error {
			inputFile := cmd.String("input_file_service_specifics")
			if inputFile == "" {
				return errors.New(
					"input file for service specifics (--input_file_service_specifics) is required",
				)
			}

			outputDir := cmd.String("output_dir")
			if outputDir == "" {
				// Use current working directory if not provided
				cwd, err := os.Getwd()
				if err != nil {
					return fmt.Errorf("failed to get current working directory: %w", err)
				}
				outputDir = cwd
			}

			var serviceSpecifics generator.ServiceSpecifics
			err := serviceSpecifics.LoadFromYAML(inputFile)
			if err != nil {
				return err
			}

			err = serviceSpecifics.Validate()
			if err != nil {
				return err
			}

			ISO19119generator, err := generator.NewISO19119Generator(
				serviceSpecifics,
				outputDir,
				nil,
				nil,
			)
			if err != nil {
				return err
			}

			err = ISO19119generator.Generate()
			if err != nil {
				return err
			}

			ISO19119generator.PrintSummary()

			return nil
		},
	}
}

func getGenerateConfigExampleCommand() *cli.Command {
	return &cli.Command{
		Name:  "config-example",
		Usage: "Shows example of <input_file_service_specifics> for users that are not familiar with the service specifics.",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "o",
				Required: true,
				Usage:    "Output file in json, yml or yaml format.",
			},
		},
		Action: func(_ context.Context, cmd *cli.Command) error {
			outputFormat := cmd.String("o")
			if outputFormat == "" {
				return errors.New("outputFormat (--o) is required")
			}

			projRoot := common.GetProjectRoot()
			var filename string

			switch strings.ToLower(outputFormat) {
			case "yaml", "yml":
				filename = "example.yaml"
			case "json":
				filename = "example.json"
			default:
				return fmt.Errorf("unsupported output format: %s", outputFormat)
			}

			path := filepath.Join(projRoot, "examples/service_specifics/", filename)
			//nolint:gosec
			data, err := os.ReadFile(path)
			if err != nil {
				return fmt.Errorf("failed to read input file: %w", err)
			}
			fmt.Print(string(data))

			return nil
		},
	}
}
