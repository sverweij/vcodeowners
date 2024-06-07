package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

const LABELER_TEST_DIR = "testdata/labeler"

func TestFormatLabeler(t *testing.T) {
	assert := assert.New(t)

	testInputs, _ := os.ReadDir(LABELER_TEST_DIR)
	for _, thing := range testInputs {
		if strings.HasSuffix(thing.Name(), "-cst.json") {
			fmt.Println("Running test for", thing.Name())
			root := strings.Replace(thing.Name(), "-cst.json", "", -1)
			t.Run(root, func(t *testing.T) {
				// codeowners CST
				var codeOwnerCST []CodeOwnersLine = []CodeOwnersLine{}
				codeOwnerCSTAsJSON, errorCodeOwnerCST := os.ReadFile(filepath.Join(LABELER_TEST_DIR, thing.Name()))
				assert.Nil(errorCodeOwnerCST)
				codeOwnerCSTJSONError := json.Unmarshal([]byte(codeOwnerCSTAsJSON), &codeOwnerCST)
				assert.Nil(codeOwnerCSTJSONError)

				// team map
				teamJSON, errorTeamJSON := os.ReadFile(filepath.Join(LABELER_TEST_DIR, root+"-teams.json"))
				assert.Nil(errorTeamJSON)
				teamMap, errorTeamMap := ParseTeamMap(string(teamJSON))
				assert.Nil(errorTeamMap)

				// expected
				expected, errorExpected := os.ReadFile(filepath.Join(LABELER_TEST_DIR, root+"-labeler.yml"))
				assert.Nil(errorExpected)

				found, _ := FormatCSTAsLabelerYML(codeOwnerCST, teamMap, "")

				assert.Equal(string(expected), found)
			})
		}
	}
	t.Run("adds a comment header when passed a non-empty string", func(t *testing.T) {
		content := `* @owner`
		parsed, _ := Parse(content)
		var teamMap = make(map[string][]string)
		expected := `# The man in black fled across the desert, and the gunslinger followed.`
		found, _ := FormatCSTAsLabelerYML(parsed, teamMap, "# The man in black fled across the desert, and the gunslinger followed.\n\n")

		assert.Contains(found, expected)
	})
}
