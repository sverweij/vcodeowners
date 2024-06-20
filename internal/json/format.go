package json

import (
	"encoding/json"

	"github.com/sverweij/vcodeowners/internal/codeowners"
)

// FormatCST takes a CST (a slice of CodeOwnersLines) and returns them
// in JSON format.
func FormatCST(cst codeowners.CST) (string, error) {
	jsonBytes, error := json.MarshalIndent(cst, "", "  ")
	return string(jsonBytes), error
}
