package main

import (
	"encoding/json"
	"fmt"
)

func formatOwners(owners []Owner) string {
	var output string
	for i, owner := range owners {
		output += owner.Name
		if i < len(owners)-1 {
			output += " "
		} else {
			output += ""
		}

	}
	return output
}

func formatInlineComment(inlineComment string) string {
	if inlineComment != "" {
		return " #" + inlineComment
	}
	return ""
}

func formatOptional(sectionOptional bool) string {
	if sectionOptional {
		return "^"
	}
	return ""
}

func formatMinApprovers(sectionMinApprovers int) string {
	if sectionMinApprovers > 0 {
		return fmt.Sprintf("[%d]", sectionMinApprovers)
	}
	return ""
}

// FormatCSTAsCodeOwners takes CST (a slice of CodeOwnersLines) and returns them
// in the CODEOWNERS format.
func FormatCSTAsCodeOwners(lines []CodeOwnersLine, headerComment string) (string, error) {
	var output string
	if headerComment != "" {
		output = headerComment
	}

	for _, line := range lines {
		switch line.Type {
		case "ignorable-comment":
		case "rule":
			output += line.RulePattern + line.Spaces + formatOwners(line.Owners) + formatInlineComment(line.InlineComment) + "\n"
		case "section-heading":
			if line.Owners != nil {
				output += formatOptional(line.SectionOptional) + "[" + line.SectionName + "]" + formatMinApprovers(line.SectionMinApprovers) + line.Spaces + formatOwners(line.Owners) + formatInlineComment(line.InlineComment) + "\n"
			} else {
				output += line.Raw + "\n"
			}
		default:
			output += line.Raw + "\n"
		}
	}
	return output, nil
}

// FormatCSTAsJSON takes CST (a slice of CodeOwnersLines) and returns them
// in JSON format.
func FormatCSTAsJSON(lines []CodeOwnersLine) (string, error) {
	jsonBytes, error := json.MarshalIndent(lines, "", "  ")
	return string(jsonBytes), error
}

func FormatCSTAsLabelerYML(lines []CodeOwnersLine, headerComment string) (string, error) {
	return "TODO", nil
}

func FormatCST(lines []CodeOwnersLine, format string) (string, error) {
	switch format {
	case "json":
		return FormatCSTAsJSON(lines)
	}
	return FormatCSTAsCodeOwners(lines, "")
}

// FormatAnomaliesAsText takes a slice of Anomalies and returns them in a human
// readable format.
func FormatAnomaliesAsText(anomalies []Anomaly) string {
	var output string = "Syntax errors found in the input:\n"
	for _, anomaly := range anomalies {
		output += "  " + anomaly.String() + "\n"
	}
	return output
}
