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

**--input_file_service_specifics**="": Path to input file containing service specifics in json, yml or yaml format. See config-example for an example of the input file.

**--output_dir**="": Location used to store service metadata as xml. If omitted the current working directory is used.

### config-example

Shows example of <input_file_service_specifics> for users that are not familiar with the service specifics.

**-o**="": Output file in json, yml or yaml format.

## hvd

Used to retrieve and inspect high value dataset categories from the HVD Thesaurus.

**--local-path**="": Local path where the HVD Thesaurus is cached. (default: /var/git/pdok-metadata-tool/cache/high-value-dataset-category.rdf)

**--url**="": HVD Thesaurus endpoint which should contain the HVD categories as RDF format. (default: https://op.europa.eu/o/opportal-service/euvoc-download-handler?cellarURI=http%3A%2F%2Fpublications.europa.eu%2Fresource%2Fdistribution%2Fhigh-value-dataset-category%2F20241002-0%2Frdf%2Fskos_core%2Fhigh-value-dataset-category.rdf&fileName=high-value-dataset-category.rdf)

### download

Downloads the RDF Thesaurus containing the HVD categories at local-path.

### list

Displays list of HVD categories.

### csv

Exports HVD categories to a CSV file.

**-o**="": Output file path for the CSV file. (default: /var/git/pdok-metadata-tool/cache/high-value-dataset-category.csv)

## inspire

The metadata toolchain is used to generate service metadata.

### list

List inspire themes or layers. Usage: pmt inspire list <theme|layer>

### csv

Exports inspire themes or layers to a CSV file. Usage: pmt inspire csv <theme|layer>

**-o**="": Output file path for the CSV file.

## store

The store is used to interact with metadata CSW store service.

### harvest

Harvest original XML metadata records from a CSW source using optional CQL filters. Records are cached on disk for inspection and reuse.

**--cache-path**="": Local path where raw CSW metadata records (XML) are cached. (default: /var/git/pdok-metadata-tool/cache/records)

**--cache-ttl**="": Cache TTL in hours for CSW record cache (default: 168 hours = 7 days). (default: 168)

**--csw-endpoint**="": Endpoint of the CSW service to harvest metadata records from. Default is NGR. (default: https://nationaalgeoregister.nl/geonetwork/srv/dut/csw)

**--filter-org**="": Optional filter by organisation name (CQL field 'OrganisationName'). Matches exact value.

**--filter-type**="": Optional filter by metadata type: 'service' or 'dataset'. If omitted, all types are harvested.

### harvest-service

Harvest service metadata (flat model) as JSON. Supports optional organisation filter and caching options.

**--cache-path**="": Local path where raw CSW metadata records (XML) are cached. (default: /var/git/pdok-metadata-tool/cache/records)

**--cache-ttl**="": Cache TTL in hours for CSW record cache (default: 168 hours = 7 days). (default: 168)

**--csw-endpoint**="": Endpoint of the CSW service to harvest metadata records from. Default is NGR. (default: https://nationaalgeoregister.nl/geonetwork/srv/dut/csw)

**--filter-org**="": Optional filter by organisation name (CQL field 'OrganisationName'). Matches exact value.

**--hvd-local-path**="": Local cache path for the HVD Thesaurus RDF. (default: /var/git/pdok-metadata-tool/cache/high-value-dataset-category.rdf)

**--hvd-url**="": HVD Thesaurus endpoint (RDF). Used to enrich HVD categories. (default: https://op.europa.eu/o/opportal-service/euvoc-download-handler?cellarURI=http%3A%2F%2Fpublications.europa.eu%2Fresource%2Fdistribution%2Fhigh-value-dataset-category%2F20241002-0%2Frdf%2Fskos_core%2Fhigh-value-dataset-category.rdf&fileName=high-value-dataset-category.rdf)

### harvest-dataset

Harvest dataset metadata (flat model) as JSON. Supports optional organisation filter and caching options.

**--cache-path**="": Local path where raw CSW metadata records (XML) are cached. (default: /var/git/pdok-metadata-tool/cache/records)

**--cache-ttl**="": Cache TTL in hours for CSW record cache (default: 168 hours = 7 days). (default: 168)

**--csw-endpoint**="": Endpoint of the CSW service to harvest metadata records from. Default is NGR. (default: https://nationaalgeoregister.nl/geonetwork/srv/dut/csw)

**--filter-org**="": Optional filter by organisation name (CQL field 'OrganisationName'). Matches exact value.

**--hvd-local-path**="": Local cache path for the HVD Thesaurus RDF. (default: /var/git/pdok-metadata-tool/cache/high-value-dataset-category.rdf)

**--hvd-url**="": HVD Thesaurus endpoint (RDF). Used to enrich HVD categories. (default: https://op.europa.eu/o/opportal-service/euvoc-download-handler?cellarURI=http%3A%2F%2Fpublications.europa.eu%2Fresource%2Fdistribution%2Fhigh-value-dataset-category%2F20241002-0%2Frdf%2Fskos_core%2Fhigh-value-dataset-category.rdf&fileName=high-value-dataset-category.rdf)
