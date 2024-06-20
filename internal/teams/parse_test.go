package teams

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/sverweij/vcodeowners/internal/codeowners"
)

func TestParseTeamMap(t *testing.T) {
	assert := assert.New(t)

	t.Run("empty team map", func(t *testing.T) {
		teamMapString := `{}`
		teamMap, error := Parse(teamMapString)

		assert.Equal(0, len(teamMap))
		assert.Nil(error)
	})
	t.Run("valid team map", func(t *testing.T) {
		teamMapString := `{"team1": ["user1", "user2"], "team2": ["user3", "user4", "user5@somehwere.else.com"]}`
		teamMap, error := Parse(teamMapString)

		assert.NotNil(teamMap)
		assert.Equal(2, len(teamMap))
		assert.Equal(2, len(teamMap["team1"]))
		assert.Equal(3, len(teamMap["team2"]))
		assert.Nil(error)
	})

	t.Run("valid team map, but one team member starts with an @", func(t *testing.T) {
		teamMapString := `{"team1": ["user1", "@user-starts-with-at"]}`
		teamMap, error := Parse(teamMapString)
		expected := "don't start team member names with an '@'; '@user-starts-with-at' (team 'team1', member 1)"

		assert.Nil(teamMap)
		assert.NotNil(error)
		assert.Equal(
			expected,
			error.Error(),
		)
	})

	t.Run("error: team map is an array", func(t *testing.T) {
		teamMapString := `["it's", "a", "trap"]`
		teamMap, error := Parse(teamMapString)

		assert.Nil(teamMap)
		assert.NotNil(error)
	})

	// other validations might follow
}

func TestApplyTeamMap(t *testing.T) {
	assert := assert.New(t)
	t.Run("replaces teams in CodeOwners CSTs - rule", func(t *testing.T) {
		codeOwners := codeowners.CST{
			{
				Type:   "rule",
				LineNo: 1,
				Raw:    "* @team2 @team1 not@replaced.nl # comment with @team1",

				RulePattern: "*",
				Spaces:      " ",
				Owners: []codeowners.Owner{
					{
						Type: "user-or-group",
						Name: "@team2",
					},
					{
						Type: "user-or-group",
						Name: "@team1",
					},
					{
						Type: "e-mail",
						Name: "not@replaced.nl",
					},
				},
				InlineComment: " comment with @team1",
			},
		}
		teamMap := Map{
			"team1": {"user1", "user2"},
			"team2": {"user3", "user4", "user5@somewhere.else.com"},
		}
		expectedCodeOwners := codeowners.CST{
			{
				Type:   "rule",
				LineNo: 1,
				Raw:    "* @team2 @team1 not@replaced.nl # comment with @team1",

				RulePattern: "*",
				Spaces:      " ",
				Owners: []codeowners.Owner{
					{
						Type: "user-or-group",
						Name: "@user1",
					},
					{
						Type: "user-or-group",
						Name: "@user2",
					},
					{
						Type: "user-or-group",
						Name: "@user3",
					},
					{
						Type: "user-or-group",
						Name: "@user4",
					},
					{
						Type: "e-mail",
						Name: "not@replaced.nl",
					},
					{
						Type: "e-mail",
						Name: "user5@somewhere.else.com",
					},
				},
				InlineComment: " comment with @team1",
			},
		}
		transformedCodeOwners := Apply(codeOwners, teamMap)

		assert.Equal(expectedCodeOwners, transformedCodeOwners)
	})

	t.Run("replaces teams in CodeOwners CSTs & uniqs them - rule", func(t *testing.T) {
		codeOwners := codeowners.CST{
			{
				Type:   "rule",
				LineNo: 1,
				Raw:    "* @team2 @team1 not@replaced.nl # comment with @team1",

				RulePattern: "*",
				Spaces:      " ",
				Owners: []codeowners.Owner{
					{
						Type: "user-or-group",
						Name: "@team2",
					},
					{
						Type: "user-or-group",
						Name: "@team1",
					},
					{
						Type: "e-mail",
						Name: "not@replaced.nl",
					},
				},
				InlineComment: " comment with @team1",
			},
		}
		teamMap := Map{
			"team1": {"user1", "user2", "not@replaced.nl"},
			"team2": {"user3", "user4", "user2", "user1", "user5@somewhere.else.com"},
		}
		expectedCodeOwners := codeowners.CST{
			{
				Type:   "rule",
				LineNo: 1,
				Raw:    "* @team2 @team1 not@replaced.nl # comment with @team1",

				RulePattern: "*",
				Spaces:      " ",
				Owners: []codeowners.Owner{
					{
						Type: "user-or-group",
						Name: "@user1",
					},
					{
						Type: "user-or-group",
						Name: "@user2",
					},
					{
						Type: "user-or-group",
						Name: "@user3",
					},
					{
						Type: "user-or-group",
						Name: "@user4",
					},
					{
						Type: "e-mail",
						Name: "not@replaced.nl",
					},
					{
						Type: "e-mail",
						Name: "user5@somewhere.else.com",
					},
				},
				InlineComment: " comment with @team1",
			},
		}
		transformedCodeOwners := Apply(codeOwners, teamMap)

		assert.Equal(expectedCodeOwners, transformedCodeOwners)
	})

	t.Run("replaces teams in CodeOwners CSTs - section-heading", func(t *testing.T) {
		codeOwners := codeowners.CST{
			{
				Type:   "section-heading",
				LineNo: 1,
				Raw:    "[some_section] @team2 @team1 not@replaced.nl # comment with @team1",

				SectionOptional:     false,
				SectionName:         "some_section",
				SectionMinApprovers: 0,
				Spaces:              " ",
				Owners: []codeowners.Owner{
					{
						Type: "user-or-group",
						Name: "@team2",
					},
					{
						Type: "user-or-group",
						Name: "@team1",
					},
					{
						Type: "user-or-group",
						// same as before - not a typo
						Name: "@team1",
					},
					{
						Type: "e-mail",
						Name: "not@replaced.nl",
					},
				},
				InlineComment: " comment with @team1",
			},
		}
		teamMap := Map{
			"team1": {"user1", "user2"},
			"team2": {"user3", "user4", "user5@somewhere.else.com"},
		}
		expectedCodeOwners := codeowners.CST{
			{
				Type:   "section-heading",
				LineNo: 1,
				Raw:    "[some_section] @team2 @team1 not@replaced.nl # comment with @team1",

				SectionOptional:     false,
				SectionName:         "some_section",
				SectionMinApprovers: 0,
				Spaces:              " ",
				Owners: []codeowners.Owner{
					{
						Type: "user-or-group",
						Name: "@user1",
					},
					{
						Type: "user-or-group",
						Name: "@user2",
					},
					{
						Type: "user-or-group",
						Name: "@user3",
					},
					{
						Type: "user-or-group",
						Name: "@user4",
					},
					{
						Type: "e-mail",
						Name: "not@replaced.nl",
					},
					{
						Type: "e-mail",
						Name: "user5@somewhere.else.com",
					},
				},
				InlineComment: " comment with @team1",
			},
		}
		transformedCodeOwners := Apply(codeOwners, teamMap)

		assert.Equal(expectedCodeOwners, transformedCodeOwners)
	})

	t.Run("only sorts owners when the team map is empty", func(t *testing.T) {
		codeOwners := codeowners.CST{
			{
				Type:   "rule",
				LineNo: 1,
				Raw:    "* not@replaced.nl @team2 @team3 @team1 # comment with @team1",

				RulePattern: "*",
				Spaces:      " ",
				Owners: []codeowners.Owner{
					{
						Type: "e-mail",
						Name: "not@replaced.nl",
					},
					{
						Type: "user-or-group",
						Name: "@team2",
					},
					{
						Type: "user-or-group",
						Name: "@team3",
					},
					{
						Type: "user-or-group",
						Name: "@team1",
					},
				},
				InlineComment: " comment with @team1",
			},
		}
		teamMap := Map{}
		expectedCodeOwners := codeowners.CST{
			{
				Type:   "rule",
				LineNo: 1,
				Raw:    "* not@replaced.nl @team2 @team3 @team1 # comment with @team1",

				RulePattern: "*",
				Spaces:      " ",
				Owners: []codeowners.Owner{
					{
						Type: "user-or-group",
						Name: "@team1",
					},
					{
						Type: "user-or-group",
						Name: "@team2",
					},
					{
						Type: "user-or-group",
						Name: "@team3",
					},
					{
						Type: "e-mail",
						Name: "not@replaced.nl",
					},
				},
				InlineComment: " comment with @team1",
			},
		}
		transformedCodeOwners := Apply(codeOwners, teamMap)

		assert.Equal(expectedCodeOwners, transformedCodeOwners)
	})
}
