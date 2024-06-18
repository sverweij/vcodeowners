package main

import (
	"fmt"
)

func formatOwners(owners []Owner) string {
	var output string
	for i, owner := range owners {
		output += owner.Name
		if i < len(owners)-1 {
			output += " "
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
	returnValue := headerComment

	for _, line := range lines {
		switch line.Type {
		case "ignorable-comment":
		case "rule":
			returnValue += line.RulePattern + line.Spaces + formatOwners(line.Owners) + formatInlineComment(line.InlineComment) + "\n"
		case "section-heading":
			if line.Owners != nil {
				returnValue += formatOptional(line.SectionOptional) + "[" + line.SectionName + "]" + formatMinApprovers(line.SectionMinApprovers) + line.Spaces + formatOwners(line.Owners) + formatInlineComment(line.InlineComment) + "\n"
			} else {
				returnValue += line.Raw + "\n"
			}
		default:
			returnValue += line.Raw + "\n"
		}
	}
	return returnValue, nil
}
