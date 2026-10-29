# NAME

pmt - PDOK Metadata Tool - This tool is set up to handle various metadata related tasks.

# SYNOPSIS

pmt

```
[--log-level]=[value]
```

**Usage**:

```
pmt [GLOBAL OPTIONS] [command [COMMAND OPTIONS]] [ARGUMENTS...]
```

# GLOBAL OPTIONS

**--log-level**="": Set log level: debug, info, warn, error (env: PMT_LOG_LEVEL) (default: info)


# COMMANDS

## generate

Used to generate metadata records.

### service

Generates service metadata in "Nederlands profiel ISO 19119" version 2.1.0.

**--input_file_service_specifics**="": Path to input file containing service specifics in json, yml or yaml format. See service-config-example for an example of the input file.

**--output_dir**="": Location used to store service metadata as xml. If omitted the current working directory is used.

### service-config-example

Shows example of <input_file_service_specifics> for users that are not familiar with the service specifics.

**-o**="": Output file in json, yml or yaml format.

### feature-catalogue

Generates feature catalogue metadata in "Nederlands profiel ISO 19110".

**--input_file_feature_catalogue_specifics**="": Path to input file containing feature catalogue specifics in json, yml or yaml format. See feature-catalogue-config-example for an example of the input file.

**--output_dir**="": Location used to store feature catalogue metadata as xml. If omitted the current working directory is used.

### feature-catalogue-example

Shows example of <input_file_service_specifics> for users that are not familiar with the feature catalogue specifics.

**-o**="": Output file in json, yml or yaml format.

## hvd

Used to retrieve and inspect high value dataset categories from the HVD Thesaurus.

**--local-path**="": Local path where the HVD Thesaurus is cached. (default: cache/high-value-dataset-category.rdf)

**--url**="": HVD Thesaurus endpoint which should contain the HVD categories as RDF format. (default: https://publications.europa.eu/resource/distribution/high-value-dataset-category/20241002-0/rdf/skos_core/high-value-dataset-category.rdf)

### download

Downloads the RDF Thesaurus containing the HVD categories at local-path.

### list

Displays list of HVD categories.

### csv

Exports HVD categories to a CSV file.

**-o**="": Output file path for the CSV file. (default: cache/high-value-dataset-category.csv)

## inspire

The metadata toolchain is used to generate service metadata.

### list

List inspire themes or layers. Usage: pmt inspire list <theme|layer>

### csv

Exports inspire themes or layers to a CSV file. Usage: pmt inspire csv <theme|layer>

**-o**="": Output file path for the CSV file.

## csw

Tools to test CSW endpoints.

### load-test

Load/stress test the NGR CSW by mimicking the INSPIRE validator, with a self-contained HTML report.

**--duration**="": Stop starting new work after this time (excluding the ID scrape), e.g. 30m. 0 means no limit. (default: 0s)

**--hosts**="": Base URLs to test, requests alternate between them. The default compares with and without 'www.'. (default: [https://nationaalgeoregister.nl https://www.nationaalgeoregister.nl])

**--ids-cache-path**="": File in which the scraped service IDs are cached. (default: cache/load-test/service-ids.json)

**--ids-cache-ttl**="": TTL of the service ID cache in hours. (default: 24)

**--include-datasets**: Also request tags and records of linked datasets, like the dataset validation does.

**--iterations**="": How many times all services are processed. 0 means until --duration is reached. (default: 1)

**--label**="": Label for this run in the report, e.g. 'laptop' or 'cluster'. (default: local)

**--limit-services**="": Only test the first N services (handy for a smoke test). (default: 0)

**--max-datasets-per-service**="": Maximum datasets requested per service. 0 means all. (default: 2)

**--no-keep-alive**: Use a new connection for every request, to measure DNS/connect/TLS on each request.

**--no-wan-ip-lookup**: Do not look up the WAN IP.

**--output-dir**="": Directory in which a new run directory with the results is created. (default: output/load-test)

**--refresh-ids**: Ignore the service ID cache and scrape again.

**--requests-per-minute**="": Maximum total requests per minute, spread over all hosts. 0 means unlimited (stress). (default: 0)

**--shuffle**: Randomise the service order in every iteration.

**--timeout**="": Timeout per request, including reading the response body. (default: 45s)

**--wan-ip**="": Report this WAN IP instead of looking it up.

**--wan-ip-url**="": Service that returns the WAN IP of this machine as plain text. (default: https://api.ipify.org)

**--workers**="": Number of services processed in parallel. (default: 2)

## store

The store is used to interact with metadata CSW store service.

### harvest

Harvest original XML metadata records from a CSW source using optional CQL filters. Records are cached on disk for inspection and reuse.

**--cache-path**="": Local path where raw CSW metadata records (XML) are cached. (default: cache/records)

**--cache-ttl**="": Cache TTL in hours for CSW record cache (default: 168 hours = 7 days). (default: 168)

**--csw-endpoint**="": Endpoint of the CSW service to harvest metadata records from. Default is NGR. (default: https://nationaalgeoregister.nl/geonetwork/srv/dut/csw)

**--filter-org**="": Optional filter by organisation name (CQL field 'OrganisationName'). Matches exact value.

**--filter-type**="": Optional filter by metadata type: 'service' or 'dataset'. If omitted, all types are harvested.

### harvest-service

Harvest service metadata (flat model) as JSON. Supports optional organisation filter and caching options.

**--cache-path**="": Local path where raw CSW metadata records (XML) are cached. (default: cache/records)

**--cache-ttl**="": Cache TTL in hours for CSW record cache (default: 168 hours = 7 days). (default: 168)

**--csw-endpoint**="": Endpoint of the CSW service to harvest metadata records from. Default is NGR. (default: https://nationaalgeoregister.nl/geonetwork/srv/dut/csw)

**--filter-org**="": Optional filter by organisation name (CQL field 'OrganisationName'). Matches exact value.

**--hvd-local-path**="": Local cache path for the HVD Thesaurus RDF. (default: cache/high-value-dataset-category.rdf)

**--hvd-url**="": HVD Thesaurus endpoint (RDF). Used to enrich HVD categories. (default: https://publications.europa.eu/resource/distribution/high-value-dataset-category/20241002-0/rdf/skos_core/high-value-dataset-category.rdf)

### harvest-dataset

Harvest dataset metadata (flat model) as JSON. Supports optional organisation filter and caching options.

**--cache-path**="": Local path where raw CSW metadata records (XML) are cached. (default: cache/records)

**--cache-ttl**="": Cache TTL in hours for CSW record cache (default: 168 hours = 7 days). (default: 168)

**--csw-endpoint**="": Endpoint of the CSW service to harvest metadata records from. Default is NGR. (default: https://nationaalgeoregister.nl/geonetwork/srv/dut/csw)

**--filter-org**="": Optional filter by organisation name (CQL field 'OrganisationName'). Matches exact value.

**--hvd-local-path**="": Local cache path for the HVD Thesaurus RDF. (default: cache/high-value-dataset-category.rdf)

**--hvd-url**="": HVD Thesaurus endpoint (RDF). Used to enrich HVD categories. (default: https://publications.europa.eu/resource/distribution/high-value-dataset-category/20241002-0/rdf/skos_core/high-value-dataset-category.rdf)
