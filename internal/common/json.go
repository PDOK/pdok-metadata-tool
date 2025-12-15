package common

import "encoding/json"

// MarshalJSON is a helper to marshal any slice to pretty JSON (for CLI output).
func MarshalJSON(v any) (string, error) {
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return "", err
	}
	return string(b), nil
}
