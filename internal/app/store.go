package app

import (
	"context"
	"fmt"
	"net/url"
	"time"

	"github.com/pdok/pdok-metadata-tool/internal/common"
	"github.com/pdok/pdok-metadata-tool/pkg/client"
	"github.com/pdok/pdok-metadata-tool/pkg/model/csw"
	"github.com/pdok/pdok-metadata-tool/pkg/model/metadata"
	"github.com/pdok/pdok-metadata-tool/pkg/model/ngr"
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
						Usage: "Local path where the HVD Thesaurus is cached.",
					},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {

					// todo: refactor MetadataRepository in kangaroo API to use endpoint ipv host + path

					// todo: harvest metadata records from CSW
					// 	- Get master page and recursively follow links so we get a list of all available records
					//	- Iterate the list of records and download the metadata records
					//	- Store the metadata records in the cache directory

					ngrEndpoint, err := url.Parse(ngr.NgrEndpoint)
					if err != nil {
						return err
					}

					cswClient := client.NewCswClient(ngrEndpoint)
					cswClient.SetCache(common.MetadataCachePath, 7*24*time.Hour)
					service := csw.Service

					org := "Beheer PDOK"

					servicesByPDOKconstraint := csw.GetRecordsCQLConstraint{
						OrganisationName: &org,
						MetadataType:     &service,
					}

					mds, err := cswClient.HarvestByCQLConstraint(&servicesByPDOKconstraint)
					if err != nil {
						return err
					}

					for _, md := range mds {
						m := metadata.NewNLServiceMetadataFromMDMetadata(&md)

						fmt.Printf("%s\t%s\t%s\t%s\n", m.ServiceType, m.Title, m.OrganisationName, m.OperatesOn)
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

					// todo: Add filters to the harvest command to only harvest certain metadata types
					//	- Owner filter
					//	- Dataset type filter

					// todo: test marshaling of service data metadata
					// todo: create a flat service model like NLDatasetMetadata but then NLServiceMetadata

					// todo: create a tool to convert harvested records to NLServiceMetadata and load them into a database

					return nil
				},
			},
		},
	}
	PDOKMetadataToolCLI.Commands = append(PDOKMetadataToolCLI.Commands, command)
}
