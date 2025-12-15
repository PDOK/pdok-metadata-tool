package app

import (
	"context"
	"fmt"

	"github.com/pdok/pdok-metadata-tool/internal/common"
	"github.com/pdok/pdok-metadata-tool/pkg/model/csw"
	"github.com/pdok/pdok-metadata-tool/pkg/model/inspire"
	"github.com/pdok/pdok-metadata-tool/pkg/model/metadata"
	"github.com/pdok/pdok-metadata-tool/pkg/model/ngr"
	"github.com/pdok/pdok-metadata-tool/pkg/repository"
	"github.com/urfave/cli/v3"
)

func init() {
	command := &cli.Command{
		Name:  "store",
		Usage: "The store is used to interact with metadata CSW store service.",
		Commands: []*cli.Command{
			{
				Name:  "harvest",
				Usage: "Harvest original metadata records from CSW source to cache directory.",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "csw-endpoint",
						Value: ngr.NgrEndpoint,
						Usage: "Endpoint of the CSW service to harvest metadata records from. Default is NGR.",
					},
					&cli.StringFlag{
						Name:  "cache-path",
						Value: common.MetadataCachePath,
						Usage: "Local path where CSW metadata records are cached.",
					},
					&cli.IntFlag{
						Name:  "cache-ttl",
						Value: 168, // hours (7 days)
						Usage: "Cache TTL in hours for CSW record cache (default: 168 hours = 7 days).",
					},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {

					// todo: refactor MetadataRepository in kangaroo API to use endpoint ipv host + path

					cswEndpoint := cmd.String("csw-endpoint")

					repo, err := repository.NewMetadataRepository(cswEndpoint)
					if err != nil {
						return err
					}

					cachePath := cmd.String("cache-path")
					cacheTTL := cmd.Int("cache-ttl")

					repo.CswClient.SetCache(cachePath, cacheTTL)

					service := csw.Service

					org := "Beheer PDOK"

					servicesByPDOKconstraint := csw.GetRecordsCQLConstraint{
						OrganisationName: &org,
						MetadataType:     &service,
					}

					mds, err := repo.CswClient.HarvestByCQLConstraint(&servicesByPDOKconstraint)
					if err != nil {
						return err
					}

					// todo: move this to a separate command and put this in repository
					for _, md := range mds {
						m := metadata.NewNLServiceMetadataFromMDMetadata(&md)

						determineDatasetInspireVariant(repo, m)

						fmt.Printf("%-15s %-15s %-100s %-30s\n",
							m.ServiceType, m.GetInspireVariant(), m.Title, m.OrganisationName)

					}

					//dataset := csw.Dataset
					//constraint2 := csw.GetRecordsCQLConstraint{
					//	MetadataType: &dataset,
					//}
					//
					//_, err = cswClient.HarvestByCQLConstraint(&constraint2)
					//if err != nil {
					//	return err
					//}

					// todo: find out how to discriminate between service metadata and dataset metadata

					// todo: create a tool to convert harvested records to NLServiceMetadata and load them into a database

					return nil
				},
			},
		},
	}
	PDOKMetadataToolCLI.Commands = append(PDOKMetadataToolCLI.Commands, command)
}

func determineDatasetInspireVariant(repo *repository.MetadataRepository, m *metadata.NLServiceMetadata) {
	for _, ref := range m.OperatesOn {

		// TODO This is slow, possibly implement a cache
		datasetMetadataId := ref.GetID()
		datasetMetadata, _ := repo.GetDatasetMetadataByID(datasetMetadataId)
		if datasetMetadata == nil {
			fmt.Println("No dataset metadata found for " + datasetMetadataId + ", for service record " + m.MetadataID)
			continue
		}

		if datasetMetadata.InspireVariant != nil {
			if len(m.OperatesOn) > 1 {
				// An INSPIRE service with multiple linked datasets can only be AsIs
				m.InspireVariant = common.Ptr(inspire.AsIs)
				return
			}
			m.InspireVariant = datasetMetadata.InspireVariant
			return
		}

	}
	m.InspireVariant = nil
	return
}
