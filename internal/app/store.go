package app

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"path/filepath"

	"github.com/pdok/pdok-metadata-tool/internal/common"
	"github.com/pdok/pdok-metadata-tool/pkg/client"
	"github.com/pdok/pdok-metadata-tool/pkg/model/csw"
	"github.com/pdok/pdok-metadata-tool/pkg/model/hvd"
	"github.com/pdok/pdok-metadata-tool/pkg/model/metadata"
	"github.com/pdok/pdok-metadata-tool/pkg/model/ngr"
	"github.com/pdok/pdok-metadata-tool/pkg/repository"
	"github.com/urfave/cli/v3"
)

// Shared CLI flags reused across subcommands
var (
	flagCswEndpoint = &cli.StringFlag{
		Name:  "csw-endpoint",
		Value: ngr.NgrEndpoint,
		Usage: "Endpoint of the CSW service to harvest metadata records from. Default is NGR.",
	}
	flagCachePath = &cli.StringFlag{
		Name:  "cache-path",
		Value: common.MetadataCachePath,
		Usage: "Local path where raw CSW metadata records (XML) are cached.",
	}
	flagCacheTTL = &cli.IntFlag{
		Name:  "cache-ttl",
		Value: 168, // hours (7 days)
		Usage: "Cache TTL in hours for CSW record cache (default: 168 hours = 7 days).",
	}
	flagFilterOrg = &cli.StringFlag{
		Name:  "filter-org",
		Usage: "Optional filter by organisation name (CQL field 'OrganisationName'). Matches exact value.",
	}
	flagFilterType = &cli.StringFlag{
		Name:  "filter-type",
		Usage: "Optional filter by metadata type: 'service' or 'dataset'. If omitted, all types are harvested.",
	}
	// HVD repository flags used to enrich HVD categories in flat outputs
	flagHvdURL = &cli.StringFlag{
		Name:        "hvd-url",
		DefaultText: "eu-thesaurus-url",
		Value:       hvd.HvdEndpoint,
		Usage:       "HVD Thesaurus endpoint (RDF). Used to enrich HVD categories.",
	}
	flagHvdLocalPath = &cli.StringFlag{
		Name:  "hvd-local-path",
		Value: common.HvdLocalRDFPath,
		Usage: "Local cache path for the HVD Thesaurus RDF.",
	}
)

func init() {
	command := &cli.Command{
		Name:  "store",
		Usage: "The store is used to interact with metadata CSW store service.",
		Commands: []*cli.Command{
			{
				Name:  "harvest",
				Usage: "Harvest original XML metadata records from a CSW source using optional CQL filters. Records are cached on disk for inspection and reuse.",
				Flags: []cli.Flag{
					flagCswEndpoint,
					flagCachePath,
					flagCacheTTL,
					flagFilterType,
					flagFilterOrg,
				},
				Action: func(_ context.Context, cmd *cli.Command) error {

					cswEndpoint := cmd.String("csw-endpoint")
					u, err := url.Parse(cswEndpoint)
					if err != nil {
						return err
					}
					cswClient := client.NewCswClient(u)

					cachePath := cmd.String("cache-path")
					cacheTTL := cmd.Int("cache-ttl")
					cswClient.SetCache(cachePath, cacheTTL)

					// Build CQL constraint from flags
					var constraint csw.GetRecordsCQLConstraint
					if t := cmd.String("filter-type"); t != "" {
						switch t {
						case csw.Service.String():
							mt := csw.Service
							constraint.MetadataType = &mt
						case csw.Dataset.String():
							mt := csw.Dataset
							constraint.MetadataType = &mt
						default:
							return fmt.Errorf("invalid --filter-type: %s (allowed: service, dataset)", t)
						}
					}
					if org := cmd.String("filter-org"); org != "" {
						constraint.OrganisationName = &org
					}

					mds, err := cswClient.HarvestByCQLConstraint(&constraint)
					if err != nil {
						return err
					}

					fmt.Printf("Harvested %d records. Cached XML in %s (TTL %d hours).\n", len(mds), cachePath, cacheTTL)

					return nil
				},
			},
			{
				Name:  "harvest-service",
				Usage: "Harvest service metadata (flat model) as JSON. Supports optional organisation filter and caching options.",
				Flags: []cli.Flag{
					flagCswEndpoint,
					flagCachePath,
					flagCacheTTL,
					flagFilterOrg,
					flagHvdURL,
					flagHvdLocalPath,
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					return harvestFlatToFile[metadata.NLServiceMetadata](cmd, csw.Service, "service-metadata", "service metadata items")
				},
			},
			{
				Name:  "harvest-dataset",
				Usage: "Harvest dataset metadata (flat model) as JSON. Supports optional organisation filter and caching options.",
				Flags: []cli.Flag{
					flagCswEndpoint,
					flagCachePath,
					flagCacheTTL,
					flagFilterOrg,
					flagHvdURL,
					flagHvdLocalPath,
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					return harvestFlatToFile[metadata.NLDatasetMetadata](cmd, csw.Dataset, "dataset-metadata", "dataset metadata items")
				},
			},
		},
	}
	PDOKMetadataToolCLI.Commands = append(PDOKMetadataToolCLI.Commands, command)
}

// harvestFlatToFile centralizes the shared logic for harvesting flat models (service/dataset),
// marshalling them to JSON, and writing the output to a file under the parent of cache-path.
func harvestFlatToFile[T any](cmd *cli.Command, mt csw.MetadataType, outBase string, summaryLabel string) error {
	// Init repository and propagate cache
	cswEndpoint := cmd.String("csw-endpoint")
	repo, err := repository.NewMetadataRepository(cswEndpoint)
	if err != nil {
		return err
	}

	cachePath := cmd.String("cache-path")
	cacheTTL := cmd.Int("cache-ttl")
	repo.SetCache(cachePath, cacheTTL)

	// Configure HVD Repository for enrichment
	hvdRepo := repository.NewHVDRepository(cmd.String("hvd-url"), cmd.String("hvd-local-path"))
	repo.SetHVDRepo(hvdRepo)

	// Build constraint with static MetadataType and optional org filter
	var constraint csw.GetRecordsCQLConstraint
	constraint.MetadataType = &mt
	if org := cmd.String("filter-org"); org != "" {
		constraint.OrganisationName = &org
	}

	// Harvest using generic repo method
	res, err := repository.HarvestByCQLConstraint[T](repo, &constraint)
	if err != nil {
		return err
	}

	// Marshal to pretty JSON
	b, err := json.MarshalIndent(res, "", "  ")
	if err != nil {
		return err
	}

	// Determine output file under parent dir of cachePath
	parentDir := filepath.Dir(cachePath)
	org := cmd.String("filter-org")
	norm := common.NormalizeForFilename(org)
	outPath := filepath.Join(parentDir, fmt.Sprintf("%s-%s.json", outBase, norm))
	if err := os.MkdirAll(parentDir, 0o755); err != nil {
		return err
	}
	if err := os.WriteFile(outPath, b, 0o644); err != nil {
		return err
	}

	fmt.Printf("Wrote %d %s to %s\n", len(res), summaryLabel, outPath)
	return nil
}
