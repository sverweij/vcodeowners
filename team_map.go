package main

import (
	"cmp"
	"encoding/json"
	"fmt"
	"slices"
	"strings"
)

func ParseTeamMap(teamMapString string) (map[string][]string, error) {
	var teamMap map[string][]string
	error := json.Unmarshal([]byte(teamMapString), &teamMap)
	// unmarshal & strong typing already ensure the map is in the right
	// shape. Only thing left is to check if the usernames don't accidentally
	// contain the "@" prefix.
	for team, members := range teamMap {
		for i, member := range members {
			if strings.HasPrefix(member, "@") {
				return nil, fmt.Errorf("don't start team member names with an '@'; '%s' (team '%s', member %d)", member, team, i)
			}
		}
	}
	return teamMap, error
}

func cookOwner(ownerString string) Owner {
	owner := ParseOwner(ownerString)

	// if an owner in a team map doesn't start with an @, ParseOwner will
	// classify it as 'invalid'. However, that's how owners should appear in
	// team maps => chuck an '@' in front of it and call it a regular 'user-or-group'
	if owner.Type == "invalid" {
		owner.Type = "user-or-group"
		owner.Name = "@" + ownerString
	}

	return owner
}

func uniqOwners(owners []Owner) []Owner {
	var visited = map[string]bool{}
	var returnValue []Owner

	for _, owner := range owners {
		if !visited[owner.Name] {
			returnValue = append(returnValue, owner)
			visited[owner.Name] = true
		}
	}
	return returnValue
}

func ApplyTeamMap(codeOwnersLines []CodeOwnersLine, teamMap map[string][]string) []CodeOwnersLine {
	transformedLines := []CodeOwnersLine{}

	for _, line := range codeOwnersLines {
		if line.Type == "rule" || line.Type == "section-heading" {
			var newOwners []Owner = []Owner{}
			for _, owner := range line.Owners {
				if members := teamMap[strings.TrimPrefix(owner.Name, "@")]; members != nil && owner.Type == "user-or-group" {
					for _, member := range members {
						newOwners = append(newOwners, cookOwner(member))
					}
				} else {
					newOwners = append(newOwners, owner)
				}
			}
			slices.SortFunc(newOwners, func(a, b Owner) int {
				return cmp.Compare(a.Name, b.Name)
			})
			line.Owners = uniqOwners(newOwners)
		}
		transformedLines = append(transformedLines, line)
	}

	return transformedLines
}
