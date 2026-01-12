// Package app provides logic for the PMT CLI.
package app

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/pdok/pdok-metadata-tool/v2/internal/common"
	"github.com/pdok/pdok-metadata-tool/v2/pkg/generator/iso19110"
	"github.com/pdok/pdok-metadata-tool/v2/pkg/generator/iso19119"
	"github.com/urfave/cli/v3"
)

func init() {
	command := &cli.Command{
		Name:  "generate",
		Usage: "Used to generate metadata records.",
		Commands: []*cli.Command{
			getGenerateServiceCommand(),
			getServiceConfigExampleCommand(),
			getGenerateFeatureCatalogueCommand(),
			getFeatureCatalogueConfigExampleCommand(),
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
				Usage:    "Path to input file containing service specifics in json, yml or yaml format. See service-config-example for an example of the input file.",
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

			outputDir, err := getOutputDir(cmd)
			if err != nil {
				return err
			}

			var serviceSpecifics iso19119.ServiceSpecifics

			err = serviceSpecifics.LoadFromYamlOrJson(inputFile)
			if err != nil {
				return err
			}

			err = serviceSpecifics.Validate()
			if err != nil {
				return err
			}

			ISO19119generator, err := iso19119.NewGenerator(
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

func getServiceConfigExampleCommand() *cli.Command {
	return getExampleCommand(
		"service-config-example",
		"Shows example of <input_file_service_specifics> for users that are not familiar with the service specifics.",
		"examples/service_specifics/",
	)
}

func getGenerateFeatureCatalogueCommand() *cli.Command {
	return &cli.Command{
		Name:  "feature-catalogue",
		Usage: "Generates feature catalogue metadata in \"Nederlands profiel ISO 19110\".",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "input_file_feature_catalogue_specifics",
				Required: true,
				Usage:    "Path to input file containing feature catalogue specifics in json, yml or yaml format. See feature-catalogue-config-example for an example of the input file.",
			},
			&cli.StringFlag{
				Name:     "output_dir",
				Required: false,
				Usage:    "Location used to store feature catalogue metadata as xml. If omitted the current working directory is used.",
			},
		},
		Action: func(_ context.Context, cmd *cli.Command) error {
			inputFile := cmd.String("input_file_feature_catalogue_specifics")
			if inputFile == "" {
				return errors.New(
					"input file for feature catalogue (--input_file_feature_catalogue_specifics) is required",
				)
			}

			outputDir, err := getOutputDir(cmd)
			if err != nil {
				return err
			}

			var featureCatalogueSpecifics iso19110.FeatureCatalogueSpecifics

			err = featureCatalogueSpecifics.LoadFromYamlOrJson(inputFile)
			if err != nil {
				return err
			}

			err = featureCatalogueSpecifics.Validate()
			if err != nil {
				return err
			}

			ISO19110generator, err := iso19110.NewGenerator(
				featureCatalogueSpecifics,
				outputDir,
			)
			if err != nil {
				return err
			}

			err = ISO19110generator.Generate()
			if err != nil {
				return err
			}

			ISO19110generator.PrintSummary()

			return nil
		},
	}
}

func getFeatureCatalogueConfigExampleCommand() *cli.Command {
	return getExampleCommand(
		"feature-catalogue-example",
		"Shows example of <input_file_service_specifics> for users that are not familiar with the feature catalogue specifics.",
		"examples/feature_catalogue_specifics/",
	)
}

func getExampleCommand(name, usage, exampleDir string) *cli.Command {
	return &cli.Command{
		Name:  name,
		Usage: usage,
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

			path := filepath.Join(projRoot, exampleDir, filename)
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

func getOutputDir(cmd *cli.Command) (string, error) {
	outputDir := cmd.String("output_dir")
	if outputDir == "" {
		// Use current working directory if not provided
		cwd, err := os.Getwd()
		if err != nil {
			return "", fmt.Errorf("failed to get current working directory: %w", err)
		}

		return cwd, nil
	}

	return outputDir, nil
}
