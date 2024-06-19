package things

import (
	"fmt"
	"strings"
)

// transform takes a string and returns a string that is equivalent in the
// yaml / minimatch sense.
func transform(pOriginalString string) string {
	var lReturnValue = pOriginalString

	// as documented in CODEOWNERS "*" means all files
	// in minimatch "*" means all files _in the root folder only_; all files over
	// there is "**" so ...
	if pOriginalString == "*" {
		lReturnValue = "**"
	}

	// naked, unquoted "*" apparently mean something different in yaml then just
	// the string "*" (yarn parsers & validators will loudly howl if you enter
	// them naked).
	// Something similar seems to go for values _starting_ with "*"
	// Quoted they're, OK, though so that's what we'll do:
	if strings.HasPrefix(lReturnValue, "*") {
		lReturnValue = "\"" + lReturnValue + "\""
	}

	// in CODEOWNERS a file pattern like src/bla/ means 'everything
	// starting with "src/bla/"'. In minimatch it means 'everything
	// equal to "src/bla/"'. If you want to convey the original meaning minimatch-
	// wise you'd use "src/bla/**". So that's what we'll do.
	//
	// TODO: This does leave things like "src/bla" in the dark, though. Maybe
	// just treat these the same?
	if strings.HasSuffix(pOriginalString, "/") {
		lReturnValue = lReturnValue + "**"
	}

	// TODO: These transformations cover my current _known_ use cases,
	// but are these _all_ anomalies? Does a conversion lib for this exist maybe?

	return lReturnValue
}

// getPatternsForTeam returns the patterns CODEOWNERS has for a given team.
func getPatternsForTeam(team string, lines []CodeOwnersLine) []string {
	returnValue := []string{}

	for _, line := range lines {
		if line.Type == "rule" {
			for _, owner := range line.Owners {
				if owner.Name == "@"+team {
					returnValue = append(returnValue, transform(line.RulePattern))
				}
			}
		}
	}

	return returnValue
}

// FormatCSTAsLabelerYML takes CST (a slice of CodeOwnersLines) and returns them
// in the format of a labeler.yml file.
func FormatCSTAsLabelerYML(lines []CodeOwnersLine, teamMap map[string][]string, headerComment string) (string, error) {
	returnValue := headerComment

	for team := range teamMap {
		patterns := getPatternsForTeam(team, lines)
		if len(patterns) > 0 {
			returnValue += fmt.Sprintf("%s:\n", team)
			returnValue += "  - changed-files:\n"
			for _, pattern := range patterns {
				returnValue += fmt.Sprintf("    - any-glob-to-any-file: %s\n", pattern)
			}
			returnValue += "\n"
		}
	}
	return returnValue, nil
}
