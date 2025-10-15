# Metadata generation

The generator package contains functionality for generating service-metadata in XML based on a single input file.  
This can be used through the CLI or by usage of the code.  
In both cases the expected input is to be modelled according to the service_specifics.  
The resulting XML will be generated in the specified output directory.

## Usage in CLI

General usage of the CLI is described in the [CLI documentation](../../docs/README.md).

An example of <input_file_service_specifics> can be shown for users that are not familiar with the service specifics:
```
pmt generate config-example -o yml
```

Service metadata can be generated in "Nederlands profiel ISO 19119" version 2.1.0, based on the specified input file:
```
pmt generate service --input_file_service_specifics ./examples/service_specifics/example.yaml --output-dir ./output 
```

## Usage from code

Usage of the generator requires an instance of a `ServiceSpecifics` configuration.

This can be done by loading an existing config-file or by setting the fields on the structs.
To load the `ServiceSpecifics` from a config-file:
```
var serviceSpecifics generator.ServiceSpecifics
err := serviceSpecifics.LoadFromYAML(inputFile)
```
The alternative is to set the fields explicitly through code.

In both cases the resulting `ServiceSpecifics` may be validated, to prevent mistakes in the resulting metadata.  
This will check for required fields, duplicate metadata id's and specific configuration for INSPIRE metadata.  
```
err = serviceSpecifics.Validate()
```

After creating (and optionally validating) the `ServiceSpecifics`, the generator can be created and used:

```
ISO19119generator, err := generator.NewISO19119Generator(serviceSpecifics,outputDir)
err = ISO19119generator.Generate()
ISO19119generator.PrintSummary()
```

The generator will run for all services defined in the service_specifics.  
For each it will generate an XML-file based on the specified id, i.e.: `<id>.xml`.  
A summary of the generated metadata can be printed by using the PrintSummary() method.

## Metadata standards

The current implementation is solely aimed towards generating metadata according to
[Dutch ISO19119  standard](https://docs.geostandaarden.nl/md/mdprofiel-iso19119/).


## Service specifics

One of the main principles of using the metadata generator is that all configuration is done in one single input file.  
So no separate flags or settings which affect the output.  
Another principle is that multiple service metadata records which serve the same dataset, are likely to share certain metadata information.  
This is facilitated in the service_specifics input file by making a distinction between fields on the global and service level.  

See for example the following input file:
```yaml
globals:
  contactOrganisationName: "Example organisation name"
  contactOrganisationUri: "http://standaarden.overheid.nl/owms/terms/organisation"
  contactEmail: "contact@example.nl"
  contactUrl: "https://www.example.nl/contact"
  qosAvailability: 99.999
  qosPerformance: 1
  qosCapacity: 100
  title: "Example title"
  creationDate: "2019-09-26"
  revisionDate: "2025-09-26"
  abstract: "Example abstract"
  keywords:
    - "AA"
    - "BB"
  serviceLicense: "https://creativecommons.org/licenses/by/4.0/deed.nl"
  useLimitation: "Geen beperkingen"
  boundingBox:
    minX: "3.2062529"
    maxX: "7.2452583"
    minY: "50.733607"
    maxY: "53.582979"
  linkedDatasets:
    - "00000000-0000-0000-0000-000000000003"
  coordinateReferenceSystem: "EPSG:28992"
services:
  - type: wfs
    id: "00000000-0000-0000-0000-000000000001"
    accessPoint: "https://example.nl/example/wfs?request=GetCapabilities&service=WFS"
  - type: wms
    id: "00000000-0000-0000-0000-000000000002"
    accessPoint: "https://example.nl/example/wms?request=GetCapabilities&service=WMS"
    keywords:
      - "AA"
      - "BB"
      - "CC"
```
This example shows the configuration for:
- a WFS-service with metadata UUID `00000000-0000-0000-0000-000000000001`
- a WMS-service with metadata UUID `00000000-0000-0000-0000-000000000002`

All fields on global level will be used to generate the metadata for both services. 
On the service level there are some required (service specific) fields, such as type, id and accessPoint.
Apart from the required fields, each of the global fields can be overridden on the service level.  

See for example the keywords for the WMS-service.  
For the WFS-service, metadata will be created with global keywords AA and BB as `00000000-0000-0000-0000-000000000001.xml`.  
For the WMS-service, metadata will be created with service specific keywords AA, BB and CC as `00000000-0000-0000-0000-000000000002.xml`.  


## INSPIRE

For INSPIRE compliant services there are several additional INSPIRE requirements for the related metadata.  
These requirements depend on the kind of service and the corresponding Implementing Rule.

In the context of PDOK, the WMS and Atom are implementations for the INSPIRE Network Services.
Other services such as the WFS and OGC API Feature are implementations for the Interoperability rule. 

The following table shows the relation between INSPIRE Dataset type, the type of OGC Webservice and the corresponding implementing rule,  
according to the PDOK context:

| INSPIRE Dataset type | INSPIRE Theme | (OGC) Webservice  | Implementing Rule | INSPIRE Service type | PDOK metadata Service type |
|----------------------|---------------|-------------------|-------------------|----------------------|----------------------------|
| Harmonised           | 1             | WMS               | Network Services  | View                 | NetworkService             |
|                      |               | Atom              |                   | Download             | NetworkService             |
|                      |               | WFS               | Interoperability  | Other                | Interoperable              |
|                      |               | OGC API Feature   |                   | Other                | Interoperable              |
| as-is                | 1+            | WMS               | Network Services  | View                 | NetworkService             |
|                      |               | Atom              |                   | Download             | NetworkService             |
|                      |               | WFS               | Interoperability  | Other                | Invocable                  |
|                      |               | OGC API Feature   |                   | Other                | Invocable                  |

For usage in PDOK context, this means that setting the global level field `InspireDatasetType` should be sufficient!  
In the case of non-INSPIRE metadata, this field can be omitted.  
For INSPIRE Harmonised datasets, use `Harmonised`. For INSPIRE AsIs dataset, use `AsIs`.  
The corresponding INSPIRE implementation will be determined automatically.

Possible values for `InspireDatasetType` are:
- `nil`
- `Harmonised`
- `AsIs`

The requirements within the PDOK context may change in the future. Also other teams may use this functionality for their own needs.  
In these cases, the INSPIRE implementation can be overruled on the service level, by using the field `InspireServiceType`.

Possible values for `InspireServiceType` are:
- `nil`
- `NetworkService`
- `Interoperable`
- `Invocable`

Note that currently there is no support for SDS Harmonised service metadata for the Interoperability Implementing Rule.

## HVD

For HVD compliant services there are also some additional requirements for the related metadata.  
These relevant HVD categories can be set by using the `hvdCategories` field, both on the global and service level.  
Only existing id's of HVD cateogories should be used here. For example:    

```yaml
hvdCategories:
- "c_b79e35eb"
- "c_4b74ea13"
```

For each category, the parent categories will be collected if applicable.  
The order of the category keywords will be according to the HVD category hierarchy.  

An overview of the available categories van be viewed through the CLI using:  

```
pmt hvd list
```
