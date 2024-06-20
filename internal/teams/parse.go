package teams

import (
	"cmp"
	"encoding/json"
	"fmt"
	"slices"
	"strings"

	"github.com/sverweij/vcodeowners/internal/codeowners"
)

type Map map[string][]string

func Parse(teamMapString string) (Map, error) {
	var teamMap Map
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

func cookOwner(ownerString string) codeowners.Owner {
	owner := codeowners.ParseOwner(ownerString)

	// if an owner in a team map doesn't start with an @, ParseOwner will
	// classify it as 'invalid'. However, that's how owners should appear in
	// team maps => chuck an '@' in front of it and call it a regular 'user-or-group'
	if owner.Type == "invalid" {
		owner.Type = "user-or-group"
		owner.Name = "@" + ownerString
	}

	return owner
}

func uniqOwners(owners []codeowners.Owner) []codeowners.Owner {
	var visited = map[string]bool{}
	var returnValue []codeowners.Owner

	for _, owner := range owners {
		if !visited[owner.Name] {
			returnValue = append(returnValue, owner)
			visited[owner.Name] = true
		}
	}
	return returnValue
}

func Apply(lines codeowners.CST, teamMap Map) codeowners.CST {
	transformedLines := codeowners.CST{}

	for _, line := range lines {
		if line.Type == "rule" || line.Type == "section-heading" {
			var newOwners []codeowners.Owner = []codeowners.Owner{}
			for _, owner := range line.Owners {
				if members := teamMap[strings.TrimPrefix(owner.Name, "@")]; members != nil && owner.Type == "user-or-group" {
					for _, member := range members {
						newOwners = append(newOwners, cookOwner(member))
					}
				} else {
					newOwners = append(newOwners, owner)
				}
			}
			slices.SortFunc(newOwners, func(a, b codeowners.Owner) int {
				return cmp.Compare(a.Name, b.Name)
			})
			line.Owners = uniqOwners(newOwners)
		}
		transformedLines = append(transformedLines, line)
	}

	return transformedLines
}
