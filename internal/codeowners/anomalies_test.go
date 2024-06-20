package codeowners

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestString(t *testing.T) {
	assert := assert.New(t)

	t.Run("invalid lines", func(t *testing.T) {
		var anomalies = Anomalies{
			{
				LineNo: 42,
				Reason: "Unknown line type",
				Raw:    "invalid line",
			},
		}

		expected := "Syntax errors found in the input:\n  Line   42, Unknown line type: \"invalid line\"\n"
		found := anomalies.String()

		assert.Equal(expected, found)
	})

}
