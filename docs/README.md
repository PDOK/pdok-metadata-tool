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

The metadata toolchain is used to generate service metadata

### service

Generates service metadata in "Nederlands profiel ISO 19119" version 2.1.0.

### config-example

Shows example of <input_file_service_specifics> for users that are not familiar with the service specifics.

## hvd

Used to retrieve and inspect high value dataset categories from the HVD Thesaurus.

**--local-path**="": HVD Thesaurus endpoint which should contain the HVD categories as RDF format. (default: high-value-dataset-category.rdf)

**--url**="": HVD Thesaurus endpoint which should contain the HVD categories as RDF format. (default: https://op.europa.eu/o/opportal-service/euvoc-download-handler?cellarURI=http%3A%2F%2Fpublications.europa.eu%2Fresource%2Fdistribution%2Fhigh-value-dataset-category%2F20241002-0%2Frdf%2Fskos_core%2Fhigh-value-dataset-category.rdf&fileName=high-value-dataset-category.rdf)

### download

Downloads the RDF Thesaurus containing the HVD categories at local-path.

### list

Displays list of HVD categories.

## inspire

The metadata toolchain is used to generate service metadata

### list

Bumps revision date to today or provided date in NGR metadata record.

**-k**="": Inspire resource kind, choose between (theme, layer) (default: theme)

## metadata

The metadata toolchain is used to generate service metadata

### bump-revision-date

Bumps revision date to today or provided date in NGR metadata record.
