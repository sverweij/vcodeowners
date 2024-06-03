package main

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

var ruleLineWithoutUsersPattern = regexp.MustCompile(`^(?<filesPattern>[^\s]+)(?<spaces>\s*)(?:#(?<comment>.*))?$`)
var ruleLinePattern = regexp.MustCompile(`^(?<filesPattern>[^\s]+)(?<spaces>\s+)(?<userNames>[^#]+)(?:#(?<comment>.*))?$`)
var sectionLineWithoutUsersPattern = regexp.MustCompile(`^(?<optionalIndicator>\^)?\[(?<name>[^\]]+)\](?:\[(?<minApprovers>[0-9]+)\])?(?:\s*)(?:#(?<comment>.*))?$`)
var sectionLinePattern = regexp.MustCompile(`^(?<optionalIndicator>\^)?\[(?<name>[^\]]+)\](?:\[(?<minApprovers>[0-9]+)\])?(?<spaces>\s+)(?<userNames>[^#]+)(?:#(?<comment>.*))?$`)
var userSeparatorPattern = regexp.MustCompile(`\s+`)
var emailPattern = regexp.MustCompile("^[^@]+@.+$")

// CodeOwnersLine represents a line in a CODEOWNERS file
type CodeOwnersLine struct {
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
	Spaces        string `json:"spaces"`
	Users         []User `json:"users"`
	InlineComment string `json:"inlineComment"`
}

// User represents a user in rule or section-heading
//
// Types:
//   - "user-or-group" (these start with an "@" symbol e.g. @john_doe or @the_a_team)
//   - "e-mail" (e-mail addresses. We're not checking against the entire RFC 5322,
//     just a simple check for presence of an "@" symbol)
//   - "invalid" (anything else)
type User struct {
	// user-or-group, e-mail, invalid
	Type string `json:"type"`
	Name string `json:"name"`
}

// Anomaly represents a problem in a CODEOWNERS file
type Anomaly struct {
	LineNo int    `json:"lineNo"`
	Reason string `json:"reason"`
	Raw    string `json:"raw"`
}

type parseState struct {
	currentSection              string
	currentSectionHasValidUsers bool
}

func (anomaly Anomaly) String() string {
	return fmt.Sprintf("Line %4d, %s: \"%s\"", anomaly.LineNo, anomaly.Reason, anomaly.Raw)
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

func parseUser(user string) User {
	if strings.HasPrefix(user, "@") {
		return User{
			Type: "user-or-group",
			Name: user,
		}
	}
	if emailPattern.MatchString(user) {
		return User{
			Type: "e-mail",
			Name: user,
		}
	}
	return User{
		Type: "invalid",
		Name: user,
	}
}

func parseUserString(usersString string) []User {
	users := userSeparatorPattern.Split(usersString, -1)
	var parsedUsers []User
	for _, user := range users {
		if user != "" {
			parsedUsers = append(parsedUsers, parseUser(user))
		}
	}
	return parsedUsers
}

func parseSectionHeadLine(line string, lineNo int, state parseState) (CodeOwnersLine, parseState) {
	noUserSectionPatternMatches := sectionLineWithoutUsersPattern.FindStringSubmatch(line)

	if noUserSectionPatternMatches != nil {
		return CodeOwnersLine{
			Type:   "section-heading",
			LineNo: lineNo,
			Raw:    line,

			SectionOptional:     noUserSectionPatternMatches[1] == "^",
			SectionName:         noUserSectionPatternMatches[2],
			SectionMinApprovers: getOptionalInt(noUserSectionPatternMatches[3]),
		}, parseState{currentSection: noUserSectionPatternMatches[2], currentSectionHasValidUsers: false}
	}

	sectionHeadPatternMatches := sectionLinePattern.FindStringSubmatch(line)
	if sectionHeadPatternMatches == nil {
		return CodeOwnersLine{
			Type:   "unknown",
			LineNo: lineNo,
			Raw:    line,
		}, state
	}

	users := parseUserString(sectionHeadPatternMatches[5])
	currentSectionHasValidUsers := false

	for _, user := range users {
		if user.Type != "invalid" {
			currentSectionHasValidUsers = true
			break
		}
	}
	return CodeOwnersLine{
			Type:   "section-heading",
			LineNo: lineNo,
			Raw:    line,

			SectionOptional:     sectionHeadPatternMatches[1] == "^",
			SectionName:         sectionHeadPatternMatches[2],
			SectionMinApprovers: getOptionalInt(sectionHeadPatternMatches[3]),
			Spaces:              sectionHeadPatternMatches[4],
			Users:               users,
			InlineComment:       sectionHeadPatternMatches[6],
		}, parseState{
			currentSection:              sectionHeadPatternMatches[2],
			currentSectionHasValidUsers: currentSectionHasValidUsers,
		}
}

func parseRuleLine(line string, lineNo int, state parseState) CodeOwnersLine {

	ruleLineWithoutUsersPatternMatches := ruleLineWithoutUsersPattern.FindStringSubmatch(line)

	if ruleLineWithoutUsersPatternMatches != nil && state.currentSectionHasValidUsers {
		return CodeOwnersLine{
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
		return CodeOwnersLine{
			Type:        "unknown",
			LineNo:      lineNo,
			Raw:         line,
			RuleSection: state.currentSection,
		}
	}
	return CodeOwnersLine{
		Type:          "rule",
		LineNo:        lineNo,
		Raw:           line,
		RulePattern:   ruleLinePatternMatches[1],
		RuleSection:   state.currentSection,
		Spaces:        ruleLinePatternMatches[2],
		Users:         parseUserString(ruleLinePatternMatches[3]),
		InlineComment: ruleLinePatternMatches[4],
	}

}

func parseLine(line string, lineNo int, state parseState) (CodeOwnersLine, parseState) {
	var trimmedLine = strings.TrimSpace(line)

	if trimmedLine == "" {
		return CodeOwnersLine{
			Type:   "empty",
			LineNo: lineNo,
			Raw:    line,
		}, state
	}
	if strings.HasPrefix(trimmedLine, "#!") {
		return CodeOwnersLine{
			Type:   "ignorable-comment",
			LineNo: lineNo,
			Raw:    line,
		}, state
	}
	if strings.HasPrefix(trimmedLine, "#") {
		return CodeOwnersLine{
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
func Parse(content string) ([]CodeOwnersLine, []Anomaly) {
	var codeOwnersLines []CodeOwnersLine
	var anomalies []Anomaly
	var state parseState
	var parsedLine CodeOwnersLine

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
			for _, user := range line.Users {
				if user.Type == "invalid" {
					anomalies = append(anomalies,
						Anomaly{
							LineNo: line.LineNo,
							Reason: fmt.Sprintf("Invalid user '%s'", user.Name),
							Raw:    line.Raw,
						},
					)
				}
			}
		}
	}

	return codeOwnersLines, anomalies
}
