# NAME

pmt - PDOK Metadata Tool - This tool is set up to handle various metadata related tasks.

# SYNOPSIS

pmt

**Usage**:

```
pmt [GLOBAL OPTIONS] [command [COMMAND OPTIONS]] [ARGUMENTS...]
```

# COMMANDS

## generate

Used to generate metadata records.

### service

Generates service metadata in "Nederlands profiel ISO 19119" version 2.1.0.

**--input_file_service_specifics**="": Path to input file containing service specifics in json, yml or yaml format. See config-example for an example of the input file.

### config-example

Shows example of <input_file_service_specifics> for users that are not familiar with the service specifics.

## hvd

Used to retrieve and inspect high value dataset categories from the HVD Thesaurus.

**--local-path**="": Local path where the HVD Thesaurus is cached. (default: ./cache/high-value-dataset-category.rdf)

**--url**="": HVD Thesaurus endpoint which should contain the HVD categories as RDF format. (default: https://op.europa.eu/o/opportal-service/euvoc-download-handler?cellarURI=http%3A%2F%2Fpublications.europa.eu%2Fresource%2Fdistribution%2Fhigh-value-dataset-category%2F20241002-0%2Frdf%2Fskos_core%2Fhigh-value-dataset-category.rdf&fileName=high-value-dataset-category.rdf)

### download

Downloads the RDF Thesaurus containing the HVD categories at local-path.

### list

Displays list of HVD categories.

### csv

Exports HVD categories to a CSV file.

**-o**="": Output file path for the CSV file. (default: ./cache/high-value-dataset-category.csv)

## inspire

The metadata toolchain is used to generate service metadata.

### list

List inspire themes or layers. Usage: pmt inspire list <themes|layers>

### csv

Exports inspire themes or layers to a CSV file. Usage: pmt inspire csv <themes|layers>

**-o**="": Output file path for the CSV file.

## store

The store is used to interact with metadata CSW store service.

### bump-revision-date

Bumps revision date to today or provided date in NGR metadata record.
