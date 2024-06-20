package json

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/sverweij/vcodeowners/internal/codeowners"
)

const JSON_TEST_DIR = "testdata/json"

func TestFormatJSON(t *testing.T) {
	assert := assert.New(t)

	testInputs, _ := os.ReadDir(JSON_TEST_DIR)

	for _, thing := range testInputs {
		if filepath.Ext(thing.Name()) == ".txt" {
			root := strings.Replace(thing.Name(), ".txt", "", -1)
			t.Run(root, func(t *testing.T) {
				content, errorContent := os.ReadFile(filepath.Join(JSON_TEST_DIR, thing.Name()))
				assert.Nil(errorContent)

				expected, errorExpected := os.ReadFile(filepath.Join(JSON_TEST_DIR, root+".json"))
				assert.Nil(errorExpected)

				parsed, _ := codeowners.Parse(string(content))
				found, _ := FormatCST(parsed)

				assert.Equal(string(expected), found)
			})
		}
	}
}
