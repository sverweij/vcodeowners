package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

const JSON_TEST_DIR = "testdata/json"
const CODEOWNERS_TEST_DIR = "testdata/codeowners"

func TestFormatCST(t *testing.T) {
	assert := assert.New(t)

	t.Run("codeowners", func(t *testing.T) {
		lines := []CodeOwnersLine{
			{Type: "comment", LineNo: 1, Raw: "# this is a comment"},
			{Type: "empty", LineNo: 2, Raw: ""},
			{
				Type:        "rule",
				LineNo:      3,
				Raw:         "*   @owner",
				RulePattern: "*",
				Users: []User{
					{Type: "user-or-team", Name: "@owner"},
				},
				InlineComment: "",
				RuleSection:   "",
				Spaces:        "   ",
			},
		}

		expected := "# this is a comment\n\n*   @owner\n"
		found, _ := FormatCST(lines, "codeowners")

		assert.Equal(expected, found)
	})

	t.Run("json", func(t *testing.T) {
		lines := []CodeOwnersLine{
			{Type: "empty", LineNo: 1, Raw: ""},
			{
				Type:        "rule",
				LineNo:      2,
				Raw:         "*   @owner",
				RulePattern: "*",
				Users: []User{
					{Type: "user-or-team", Name: "@owner"},
				},
				InlineComment: "",
				RuleSection:   "",
				Spaces:        "   ",
			},
		}

		expected := `[
  {
    "type": "empty",
    "lineNo": 1,
    "raw": "",
    "rulePattern": "",
    "ruleSection": "",
    "sectionOptional": false,
    "sectionName": "",
    "sectionMinApprovers": 0,
    "spaces": "",
    "users": null,
    "inlineComment": ""
  },
  {
    "type": "rule",
    "lineNo": 2,
    "raw": "*   @owner",
    "rulePattern": "*",
    "ruleSection": "",
    "sectionOptional": false,
    "sectionName": "",
    "sectionMinApprovers": 0,
    "spaces": "   ",
    "users": [
      {
        "type": "user-or-team",
        "name": "@owner"
      }
    ],
    "inlineComment": ""
  }
]`
		found, _ := FormatCST(lines, "json")

		assert.Equal(expected, found)
	})
}

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

				parsed, _ := Parse(string(content))
				found, _ := FormatCSTAsJSON(parsed)

				assert.Equal(string(expected), found)
			})
		}
	}
}

func TestFormatCodeOwners(t *testing.T) {
	assert := assert.New(t)

	testInputs, _ := os.ReadDir(CODEOWNERS_TEST_DIR)

	for _, thing := range testInputs {
		if strings.HasSuffix(thing.Name(), "-mock.txt") {
			root := strings.Replace(thing.Name(), "-mock.txt", "", -1)
			t.Run(root, func(t *testing.T) {
				content, errorContent := os.ReadFile(filepath.Join(CODEOWNERS_TEST_DIR, thing.Name()))
				assert.Nil(errorContent)

				expected, errorExpected := os.ReadFile(filepath.Join(CODEOWNERS_TEST_DIR, root+"-fixture.txt"))
				assert.Nil(errorExpected)

				parsed, _ := Parse(string(content))
				found, _ := FormatCSTAsCodeOwners(parsed)

				assert.Equal(string(expected), found)
			})
		}
	}
}

func TestFormatAnomalies(t *testing.T) {
	assert := assert.New(t)

	t.Run("invalid lines", func(t *testing.T) {
		var anomalies = []Anomaly{
			{
				LineNo: 42,
				Reason: "Unknown line type",
				Raw:    "invalid line",
			},
		}

		expected := "Syntax errors found in the input:\n  Line   42, Unknown line type: \"invalid line\"\n"
		found := FormatAnomaliesAsText(anomalies)

		assert.Equal(expected, found)
	})

}
