package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseTeamMap(t *testing.T) {
	assert := assert.New(t)

	t.Run("empty team map", func(t *testing.T) {
		teamMapString := `{}`
		teamMap, _ := ParseTeamMap(teamMapString)

		assert.Equal(0, len(teamMap))
	})
	t.Run("valid team map", func(t *testing.T) {
		teamMapString := `{"team1": ["user1", "user2"], "team2": ["user3", "user4", "user5@somehwere.else.com"]}`
		teamMap, _ := ParseTeamMap(teamMapString)

		assert.NotNil(teamMap)
		assert.Equal(2, len(teamMap))
		assert.Equal(2, len(teamMap["team1"]))
		assert.Equal(3, len(teamMap["team2"]))
	})

	t.Run("error: team map is an array", func(t *testing.T) {
		teamMapString := `["it's", "a", "trap"]`
		teamMap, error := ParseTeamMap(teamMapString)

		assert.Nil(teamMap)
		assert.NotNil(error)
	})

	// other validations might follow
}
