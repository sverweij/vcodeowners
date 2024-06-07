package main

import (
	"encoding/json"
)

// FormatCSTAsJSON takes a CST (a slice of CodeOwnersLines) and returns them
// in JSON format.
func FormatCSTAsJSON(lines []CodeOwnersLine) (string, error) {
	jsonBytes, error := json.MarshalIndent(lines, "", "  ")
	return string(jsonBytes), error
}
