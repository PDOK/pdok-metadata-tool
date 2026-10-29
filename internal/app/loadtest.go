package app

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"time"

	"github.com/pdok/pdok-metadata-tool/v2/internal/common"
	"github.com/pdok/pdok-metadata-tool/v2/pkg/loadtest"
	"github.com/urfave/cli/v3"
)

const (
	defaultLoadTestTimeout   = 45 * time.Second
	defaultIDCacheTTLHrs     = 24
	defaultLoadTestWorkers   = 2
	defaultDatasetsPerSvc    = 2
	defaultWANIPLookupURL    = "https://api.ipify.org"
	defaultLoadTestOutDir    = "output/load-test"
	defaultLoadTestHostPlain = "https://nationaalgeoregister.nl"
	defaultLoadTestHostWWW   = "https://www.nationaalgeoregister.nl"
)

func init() {
	command := &cli.Command{
		Name:  "csw",
		Usage: "Tools to test CSW endpoints.",
		Commands: []*cli.Command{
			{
				Name: "load-test",
				Usage: "Load/stress test the NGR CSW by mimicking the INSPIRE validator, " +
					"with a self-contained HTML report.",
				Description: "First scrapes all PDOK service records (cached locally), then for every " +
					"service sends the requests the INSPIRE validator sends: GetRecordById for the service " +
					"and, for each dataset it operates on, the tags request and GetRecordById for INSPIRE " +
					"datasets. All hosts are tested alternately, never using a cache. Repeat with " +
					"--iterations and/or --duration. Writes results.jsonl, failures.csv, summary.json and report.html.",
				Flags:  loadTestFlags(),
				Action: runLoadTest,
			},
		},
	}
	PDOKMetadataToolCLI.Commands = append(PDOKMetadataToolCLI.Commands, command)
}

func loadTestFlags() []cli.Flag {
	return []cli.Flag{
		&cli.StringSliceFlag{
			Name:  "hosts",
			Value: []string{defaultLoadTestHostPlain, defaultLoadTestHostWWW},
			Usage: "Base URLs to test, requests alternate between them. The default compares with and without 'www.'.",
		},
		&cli.IntFlag{
			Name:  "workers",
			Value: defaultLoadTestWorkers,
			Usage: "Number of services processed in parallel.",
		},
		&cli.IntFlag{
			Name:  "requests-per-minute",
			Usage: "Maximum total requests per minute, spread over all hosts. 0 means unlimited (stress).",
		},
		&cli.DurationFlag{
			Name:  "timeout",
			Value: defaultLoadTestTimeout,
			Usage: "Timeout per request, including reading the response body.",
		},
		&cli.IntFlag{
			Name:  "iterations",
			Value: 1,
			Usage: "How many times all services are processed. 0 means until --duration is reached.",
		},
		&cli.DurationFlag{
			Name:  "duration",
			Usage: "Stop starting new work after this time (excluding the ID scrape), e.g. 30m. 0 means no limit.",
		},
		&cli.BoolFlag{
			Name:  "include-datasets",
			Value: true,
			Usage: "Also request tags and records of linked datasets, like the dataset validation does.",
		},
		&cli.IntFlag{
			Name:  "max-datasets-per-service",
			Value: defaultDatasetsPerSvc,
			Usage: "Maximum datasets requested per service. 0 means all.",
		},
		&cli.BoolFlag{
			Name:  "no-keep-alive",
			Usage: "Use a new connection for every request, to measure DNS/connect/TLS on each request.",
		},
		&cli.BoolFlag{Name: "shuffle", Usage: "Randomise the service order in every iteration."},
		&cli.IntFlag{
			Name:  "limit-services",
			Usage: "Only test the first N services (handy for a smoke test).",
		},
		&cli.StringFlag{
			Name:  "label",
			Value: "local",
			Usage: "Label for this run in the report, e.g. 'laptop' or 'cluster'.",
		},
		&cli.StringFlag{
			Name:  "output-dir",
			Value: defaultLoadTestOutDir,
			Usage: "Directory in which a new run directory with the results is created.",
		},
		&cli.StringFlag{
			Name:  "ids-cache-path",
			Value: filepath.Join(common.CachePath, "load-test", "service-ids.json"),
			Usage: "File in which the scraped service IDs are cached.",
		},
		&cli.IntFlag{
			Name:  "ids-cache-ttl",
			Value: defaultIDCacheTTLHrs,
			Usage: "TTL of the service ID cache in hours.",
		},
		&cli.BoolFlag{Name: "refresh-ids", Usage: "Ignore the service ID cache and scrape again."},
		&cli.StringFlag{
			Name:  "wan-ip",
			Usage: "Report this WAN IP instead of looking it up.",
		},
		&cli.StringFlag{
			Name:  "wan-ip-url",
			Value: defaultWANIPLookupURL,
			Usage: "Service that returns the WAN IP of this machine as plain text.",
		},
		&cli.BoolFlag{Name: "no-wan-ip-lookup", Usage: "Do not look up the WAN IP."},
	}
}

func runLoadTest(ctx context.Context, cmd *cli.Command) error {
	hosts := normalizeHosts(cmd.StringSlice("hosts"))
	if len(hosts) == 0 {
		return errors.New("--hosts must contain at least one URL")
	}

	if cmd.Int("iterations") == 0 && cmd.Duration("duration") == 0 {
		return errors.New("set --iterations and/or --duration")
	}

	lookupURL := cmd.String("wan-ip-url")
	if cmd.Bool("no-wan-ip-lookup") {
		lookupURL = ""
	}

	ctx, stop := signal.NotifyContext(ctx, os.Interrupt)
	defer stop()

	outcome, err := loadtest.Execute(ctx, loadtest.Options{
		Config: loadtest.Config{
			Hosts:                 hosts,
			Workers:               cmd.Int("workers"),
			RequestsPerMinute:     cmd.Int("requests-per-minute"),
			Timeout:               cmd.Duration("timeout"),
			Iterations:            cmd.Int("iterations"),
			Duration:              cmd.Duration("duration"),
			IncludeDatasets:       cmd.Bool("include-datasets"),
			MaxDatasetsPerService: cmd.Int("max-datasets-per-service"),
			KeepAlive:             !cmd.Bool("no-keep-alive"),
			Shuffle:               cmd.Bool("shuffle"),
		},
		IDCache: loadtest.IDCacheOptions{
			Path:    cmd.String("ids-cache-path"),
			TTL:     time.Duration(cmd.Int("ids-cache-ttl")) * time.Hour,
			Refresh: cmd.Bool("refresh-ids"),
		},
		LimitServices:  cmd.Int("limit-services"),
		Label:          cmd.String("label"),
		OutputDir:      cmd.String("output-dir"),
		WANIP:          cmd.String("wan-ip"),
		WANIPLookupURL: lookupURL,
		Progress:       os.Stderr,
	})
	if err != nil {
		return err
	}

	fmt.Print(loadtest.TextSummary(outcome))
	fmt.Printf("Report: %s\n", filepath.Join(outcome.Dir, "report.html"))

	return nil
}

func normalizeHosts(hosts []string) []string {
	out := make([]string, 0, len(hosts))

	for _, host := range hosts {
		host = strings.TrimRight(strings.TrimSpace(host), "/")
		if host != "" {
			out = append(out, host)
		}
	}

	return out
}
