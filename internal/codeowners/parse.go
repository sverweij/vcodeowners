package codeowners

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

var ruleLineWithoutOwnersPattern = regexp.MustCompile(`^(?<filesPattern>[^\s]+)(?<spaces>\s*)(?:#(?<comment>.*))?$`)
var ruleLinePattern = regexp.MustCompile(`^(?<filesPattern>[^\s]+)(?<spaces>\s+)(?<ownerNames>[^#]+)(?:#(?<comment>.*))?$`)
var sectionLineWithoutOwnersPattern = regexp.MustCompile(`^(?<optionalIndicator>\^)?\[(?<name>[^\]]+)\](?:\[(?<minApprovers>[0-9]+)\])?(?:\s*)(?:#(?<comment>.*))?$`)
var sectionLinePattern = regexp.MustCompile(`^(?<optionalIndicator>\^)?\[(?<name>[^\]]+)\](?:\[(?<minApprovers>[0-9]+)\])?(?<spaces>\s+)(?<userNames>[^#]+)(?:#(?<comment>.*))?$`)
var ownerSeparatorPattern = regexp.MustCompile(`\s+`)
var userOrGroupPattern = regexp.MustCompile("^@.+$")
var emailPattern = regexp.MustCompile("^[^@]+@.+$")

// Line represents a line in a CODEOWNERS file
type Line struct {
	// rule, section-heading, comment, ignorable-comment, empty, unknown,
	Type   string `json:"type"`
	LineNo int    `json:"lineNo"`
	Raw    string `json:"raw"`

	// rule only
	RulePattern string `json:"rulePattern"`
	RuleSection string `json:"ruleSection"`

	// section heading
	SectionOptional     bool   `json:"sectionOptional"`
	SectionName         string `json:"sectionName"`
	SectionMinApprovers int    `json:"sectionMinApprovers"`

	// rule and section heading
	Spaces        string  `json:"spaces"`
	Owners        []Owner `json:"owners"`
	InlineComment string  `json:"inlineComment"`
}

// CST represents the Concrete Syntax Tree of a CODEOWNERS file
type CST []Line

// Owner represents an owner in rule or section-heading
//
// Types:
//   - "user-or-group" (these start with an "@" symbol e.g. @john_doe or @the_a_team)
//   - "e-mail" (e-mail addresses. We're not checking against the entire RFC 5322,
//     just a simple check for presence of an "@" symbol)
//   - "invalid" (anything else)
type Owner struct {
	// user-or-group, e-mail, invalid
	Type string `json:"type"`
	Name string `json:"name"`
}

type parseState struct {
	currentSection              string
	currentSectionHasValidUsers bool
}

func getOptionalInt(fragment string) int {
	if fragment == "" {
		return 0
	}
	optionalInt, error := strconv.Atoi(fragment)
	if error != nil {
		// probably should bork with an error here
		return 0
	}
	return optionalInt

}

func ParseOwner(owner string) Owner {
	if userOrGroupPattern.MatchString(owner) {
		return Owner{
			Type: "user-or-group",
			Name: owner,
		}
	}
	if emailPattern.MatchString(owner) {
		return Owner{
			Type: "e-mail",
			Name: owner,
		}
	}
	return Owner{
		Type: "invalid",
		Name: owner,
	}
}

func parseOwnersString(ownersString string) []Owner {
	owners := ownerSeparatorPattern.Split(ownersString, -1)
	var parsedOwners []Owner
	for _, owner := range owners {
		if owner != "" {
			parsedOwners = append(parsedOwners, ParseOwner(owner))
		}
	}
	return parsedOwners
}

func parseSectionHeadLine(line string, lineNo int, state parseState) (Line, parseState) {
	ownerlessSectionPatternMatches := sectionLineWithoutOwnersPattern.FindStringSubmatch(line)

	if ownerlessSectionPatternMatches != nil {
		return Line{
			Type:   "section-heading",
			LineNo: lineNo,
			Raw:    line,

			SectionOptional:     ownerlessSectionPatternMatches[1] == "^",
			SectionName:         ownerlessSectionPatternMatches[2],
			SectionMinApprovers: getOptionalInt(ownerlessSectionPatternMatches[3]),
		}, parseState{currentSection: ownerlessSectionPatternMatches[2], currentSectionHasValidUsers: false}
	}

	sectionHeadPatternMatches := sectionLinePattern.FindStringSubmatch(line)
	if sectionHeadPatternMatches == nil {
		return Line{
			Type:   "unknown",
			LineNo: lineNo,
			Raw:    line,
		}, state
	}

	owners := parseOwnersString(sectionHeadPatternMatches[5])
	currentSectionHasValidUsers := false

	for _, owner := range owners {
		if owner.Type != "invalid" {
			currentSectionHasValidUsers = true
			break
		}
	}
	return Line{
			Type:   "section-heading",
			LineNo: lineNo,
			Raw:    line,

			SectionOptional:     sectionHeadPatternMatches[1] == "^",
			SectionName:         sectionHeadPatternMatches[2],
			SectionMinApprovers: getOptionalInt(sectionHeadPatternMatches[3]),
			Spaces:              sectionHeadPatternMatches[4],
			Owners:              owners,
			InlineComment:       sectionHeadPatternMatches[6],
		}, parseState{
			currentSection:              sectionHeadPatternMatches[2],
			currentSectionHasValidUsers: currentSectionHasValidUsers,
		}
}

func parseRuleLine(line string, lineNo int, state parseState) Line {

	ruleLineWithoutUsersPatternMatches := ruleLineWithoutOwnersPattern.FindStringSubmatch(line)

	if ruleLineWithoutUsersPatternMatches != nil && state.currentSectionHasValidUsers {
		return Line{
			Type:          "rule",
			LineNo:        lineNo,
			Raw:           line,
			RulePattern:   ruleLineWithoutUsersPatternMatches[1],
			RuleSection:   state.currentSection,
			Spaces:        ruleLineWithoutUsersPatternMatches[2],
			InlineComment: ruleLineWithoutUsersPatternMatches[3],
		}
	}

	ruleLinePatternMatches := ruleLinePattern.FindStringSubmatch(line)

	if ruleLinePatternMatches == nil {
		return Line{
			Type:        "unknown",
			LineNo:      lineNo,
			Raw:         line,
			RuleSection: state.currentSection,
		}
	}
	return Line{
		Type:          "rule",
		LineNo:        lineNo,
		Raw:           line,
		RulePattern:   ruleLinePatternMatches[1],
		RuleSection:   state.currentSection,
		Spaces:        ruleLinePatternMatches[2],
		Owners:        parseOwnersString(ruleLinePatternMatches[3]),
		InlineComment: ruleLinePatternMatches[4],
	}

}

func parseLine(line string, lineNo int, state parseState) (Line, parseState) {
	var trimmedLine = strings.TrimSpace(line)

	if trimmedLine == "" {
		return Line{
			Type:   "empty",
			LineNo: lineNo,
			Raw:    line,
		}, state
	}
	if strings.HasPrefix(trimmedLine, "#!") {
		return Line{
			Type:   "ignorable-comment",
			LineNo: lineNo,
			Raw:    line,
		}, state
	}
	if strings.HasPrefix(trimmedLine, "#") {
		return Line{
			Type:   "comment",
			LineNo: lineNo,
			Raw:    line,
		}, state
	}
	if strings.HasPrefix(trimmedLine, "[") || strings.HasPrefix(trimmedLine, "^[") {
		return parseSectionHeadLine(line, lineNo, state)
	}

	return parseRuleLine(trimmedLine, lineNo, state), state
}

// Parse parses the content of a CODEOWNERS file and returns a CST (and a list
// of syntax errors and other anomalies)
func Parse(content string) (CST, Anomalies) {
	var codeOwnersLines CST
	var anomalies Anomalies
	var state parseState
	var parsedLine Line

	lines := strings.Split(content, "\n")
	for lineNo, line := range lines {
		parsedLine, state = parseLine(line, lineNo+1, state)
		codeOwnersLines = append(codeOwnersLines, parsedLine)
	}

	for _, line := range codeOwnersLines {
		if line.Type == "unknown" {
			anomalies = append(anomalies,
				Anomaly{
					LineNo: line.LineNo,
					Reason: "Unknown line type",
					Raw:    line.Raw,
				},
			)
		}
		if line.Type == "rule" || line.Type == "section-heading" {
			for _, owner := range line.Owners {
				if owner.Type == "invalid" {
					anomalies = append(anomalies,
						Anomaly{
							LineNo: line.LineNo,
							Reason: fmt.Sprintf("Invalid user '%s'", owner.Name),
							Raw:    line.Raw,
						},
					)
				}
			}
		}
	}

	return codeOwnersLines, anomalies
}
