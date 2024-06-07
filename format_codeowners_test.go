package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

const CODEOWNER_TEST_DIR = "testdata/codeowners"

func TestFormatCodeOwners(t *testing.T) {
	assert := assert.New(t)

	testInputs, _ := os.ReadDir(CODEOWNER_TEST_DIR)

	for _, thing := range testInputs {
		if strings.HasSuffix(thing.Name(), "-mock.txt") {
			root := strings.Replace(thing.Name(), "-mock.txt", "", -1)
			t.Run(root, func(t *testing.T) {
				content, errorContent := os.ReadFile(filepath.Join(CODEOWNER_TEST_DIR, thing.Name()))
				assert.Nil(errorContent)

				expected, errorExpected := os.ReadFile(filepath.Join(CODEOWNER_TEST_DIR, root+"-fixture.txt"))
				assert.Nil(errorExpected)

				parsed, _ := Parse(string(content))
				found, _ := FormatCSTAsCodeOwners(parsed, "")

				assert.Equal(string(expected), found)
			})
		}
	}
	t.Run("adds a comment header when passed a non-empty string", func(t *testing.T) {
		content := `* @owner`
		parsed, _ := Parse(content)
		expected := `# The man in black fled across the desert, and the gunslinger followed.`
		found, _ := FormatCSTAsCodeOwners(parsed, "# The man in black fled across the desert, and the gunslinger followed.")

		assert.Contains(found, expected)
	})
}
