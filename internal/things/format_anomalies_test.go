package things

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

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
