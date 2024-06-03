package main

import (
	"cmp"
	"encoding/json"
	"slices"
	"strings"
)

func ParseTeamMap(teamMapString string) (map[string][]string, error) {
	var teamMap map[string][]string
	error := json.Unmarshal([]byte(teamMapString), &teamMap)
	return teamMap, error
}

func ApplyTeamMap(codeOwnersLines []CodeOwnersLine, teamMap map[string][]string) []CodeOwnersLine {
	transformedCodeOwnersLines := []CodeOwnersLine{}

	for _, line := range codeOwnersLines {
		if line.Type == "rule" || line.Type == "section-heading" {
			var newUsers []User = []User{}
			for _, user := range line.Users {
				if members := teamMap[strings.TrimPrefix(user.Name, "@")]; members != nil && user.Type == "user-or-group" {
					for _, member := range members {
						newUsers = append(newUsers, User{Type: "user-or-group", Name: "@" + member})
					}
				} else {
					newUsers = append(newUsers, user)
				}
			}
			slices.SortFunc(newUsers, func(a, b User) int {
				return cmp.Compare(a.Name, b.Name)
			})
			line.Users = newUsers
		}
		transformedCodeOwnersLines = append(transformedCodeOwnersLines, line)
	}

	return transformedCodeOwnersLines
}
