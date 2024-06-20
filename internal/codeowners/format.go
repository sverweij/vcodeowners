package codeowners

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

func formatRule(line Line) string {
	return line.RulePattern + line.Spaces + formatOwners(line.Owners) + formatInlineComment(line.InlineComment)
}

func formatSectionHeading(line Line) string {
	return formatOptional(line.SectionOptional) + "[" + line.SectionName + "]" + formatMinApprovers(line.SectionMinApprovers) + line.Spaces + formatOwners(line.Owners) + formatInlineComment(line.InlineComment)
}

func (line Line) String() string {
	returnValue := ""
	switch line.Type {
	case "ignorable-comment":
	case "rule":
		returnValue += formatRule(line) + "\n"
	case "section-heading":
		if line.Owners != nil {
			returnValue += formatSectionHeading(line) + "\n"
		} else {
			returnValue += line.Raw + "\n"
		}
	default:
		returnValue += line.Raw + "\n"
	}
	return returnValue
}

// Format returns the cst as a string in CODEOWNERS format.
func (cst CST) Format(headerComment string) (string, error) {
	returnValue := headerComment

	for _, line := range cst {
		returnValue += line.String()
	}
	return returnValue, nil
}
