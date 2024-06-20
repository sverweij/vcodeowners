package codeowners

import "fmt"

// Anomaly represents a problem in a CODEOWNERS file
type Anomaly struct {
	LineNo int    `json:"lineNo"`
	Reason string `json:"reason"`
	Raw    string `json:"raw"`
}

func (anomaly Anomaly) String() string {
	return fmt.Sprintf("Line %4d, %s: \"%s\"", anomaly.LineNo, anomaly.Reason, anomaly.Raw)
}

// Anomalies is a collection of problems in a CODEOWNERS file
type Anomalies []Anomaly

func (anomalies Anomalies) String() string {
	var output string = "Syntax errors found in the input:\n"
	for _, anomaly := range anomalies {
		output += "  " + anomaly.String() + "\n"
	}
	return output
}
