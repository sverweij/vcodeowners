package main

// FormatAnomaliesAsText takes a slice of Anomalies and returns them in a human
// readable format.
func FormatAnomaliesAsText(anomalies []Anomaly) string {
	var output string = "Syntax errors found in the input:\n"
	for _, anomaly := range anomalies {
		output += "  " + anomaly.String() + "\n"
	}
	return output
}
