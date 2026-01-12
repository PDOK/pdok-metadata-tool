package utils

import (
	"encoding/xml"
	"os"
	"strings"

	"github.com/ucarion/c14n"
)

func CanonicalizeXML(path string) (string, error) {
	//nolint:gosec
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}

	decoder := xml.NewDecoder(strings.NewReader(string(data)))

	canonical, err := c14n.Canonicalize(decoder)
	if err != nil {
		return "", err
	}

	return string(canonical), nil
}
