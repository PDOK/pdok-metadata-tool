package core

import (
	"encoding/xml"
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

type MetadataEntry[M any, C interface{ Config }] struct {
	// Config type, either service- or feature catalogue- specifics.
	Config C
	// Metadata type, either ISO19119 or ISO19110
	Metadata M
	Output   []byte
	Filename string
}

func (e *MetadataEntry[M, C]) SetFilename() {
	e.Filename = e.Config.GetID() + ".xml"
}

func (e *MetadataEntry[M, C]) GetID() string {
	return e.Config.GetID()
}

// Generator is generic, used by both ISO19110 and ISO19119.
type Generator[M any, C interface{ Config }] struct {
	MetadataHolder map[string]*MetadataEntry[M, C]
	CurrentID      *string
	OutputDir      string
}

func (g *Generator[M, C]) CurrentEntry() (*MetadataEntry[M, C], error) {
	if g.CurrentID == nil {
		return nil, errors.New("CurrentID is nil")
	}

	entry, ok := g.MetadataHolder[*g.CurrentID]
	if !ok {
		return nil, fmt.Errorf("no entry found for id: %s", *g.CurrentID)
	}

	return entry, nil
}

// CreateXML creates XML based on the available metadata.
func (g *Generator[M, C]) CreateXML() error {
	entry, err := g.CurrentEntry()
	if err != nil {
		return err
	}

	output, err := xml.MarshalIndent(entry.Metadata, "", "  ")
	if err != nil {
		return err
	}

	entry.Output = output

	return nil
}

// WriteToFile writes the available metadata to a file.
func (g *Generator[M, C]) WriteToFile() error {
	entry, err := g.CurrentEntry()
	if err != nil {
		return err
	}

	perm := 750
	if err1 := os.MkdirAll(g.OutputDir, os.FileMode(perm)); err1 != nil {
		return err1
	}

	entry.SetFilename()

	path := filepath.Join(g.OutputDir, entry.Filename)

	perm = 0600
	if err = os.WriteFile(path, entry.Output, os.FileMode(perm)); err != nil {
		return err
	}

	return nil
}

// PrintSummary prints a summary of the generated metadata files.
func (g *Generator[M, C]) PrintSummary() {
	fmt.Printf("The following metadata has been created in %s: \n", g.OutputDir)

	for _, entry := range g.MetadataHolder {
		fmt.Printf("  - %s\n", entry.Filename)
	}
}
