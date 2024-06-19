package things

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {
	assert := assert.New(t)
	t.Run("empty content", func(t *testing.T) {
		content := ``
		codeOwnersLines, anomalies := Parse(content)

		assert.Equal(1, len(codeOwnersLines))
		assert.Equal(CodeOwnersLine{
			Type:   "empty",
			LineNo: 1,
			Raw:    content,
		}, codeOwnersLines[0])

		assert.Equal(0, len(anomalies))
	})
	t.Run("comment", func(t *testing.T) {
		content := `# This is a comment`
		codeOwnersLines, anomalies := Parse(content)

		assert.Equal(1, len(codeOwnersLines))
		assert.Equal(CodeOwnersLine{
			Type:   "comment",
			LineNo: 1,
			Raw:    content,
		}, codeOwnersLines[0])

		assert.Equal(0, len(anomalies))
	})

	t.Run("ignorable comment", func(t *testing.T) {
		content := `#! This is an ignorable comment`
		codeOwnersLines, anomalies := Parse(content)

		assert.Equal(1, len(codeOwnersLines))
		assert.Equal(CodeOwnersLine{
			Type:   "ignorable-comment",
			LineNo: 1,
			Raw:    "#! This is an ignorable comment",
		}, codeOwnersLines[0])

		assert.Equal(0, len(anomalies))
	})
	t.Run("section without owners", func(t *testing.T) {
		content := `[section]`
		codeOwnersLines, anomalies := Parse(content)

		assert.Equal(1, len(codeOwnersLines))
		assert.Equal(CodeOwnersLine{
			Type:                "section-heading",
			LineNo:              1,
			Raw:                 content,
			SectionOptional:     false,
			SectionName:         "section",
			SectionMinApprovers: 0,
		}, codeOwnersLines[0])

		assert.Equal(0, len(anomalies))
	})
	t.Run("section without owners", func(t *testing.T) {
		content := `[section]`
		codeOwnersLines, anomalies := Parse(content)

		assert.Equal(1, len(codeOwnersLines))
		assert.Equal(CodeOwnersLine{
			Type:                "section-heading",
			LineNo:              1,
			Raw:                 content,
			SectionOptional:     false,
			SectionName:         "section",
			SectionMinApprovers: 0,
		}, codeOwnersLines[0])

		assert.Equal(0, len(anomalies))
	})
	t.Run("section kitchensink", func(t *testing.T) {
		content := `^[section][42]      @user1 @user2 # inline comment`
		codeOwnersLines, anomalies := Parse(content)

		assert.Equal(1, len(codeOwnersLines))
		assert.Equal(CodeOwnersLine{
			Type:                "section-heading",
			LineNo:              1,
			Raw:                 content,
			SectionOptional:     true,
			SectionName:         "section",
			SectionMinApprovers: 42,
			Spaces:              "      ",
			Owners: []Owner{
				{Name: "@user1", Type: "user-or-group"},
				{Name: "@user2", Type: "user-or-group"},
			},
			InlineComment: " inline comment",
		}, codeOwnersLines[0])

		assert.Equal(0, len(anomalies))
	})
	t.Run("starts as a section, but isn't", func(t *testing.T) {
		content := `[section aap noot mies`
		codeOwnersLines, anomalies := Parse(content)

		assert.Equal(1, len(codeOwnersLines))
		assert.Equal(CodeOwnersLine{
			Type:   "unknown",
			LineNo: 1,
			Raw:    content,
		}, codeOwnersLines[0])

		assert.Equal(1, len(anomalies))
		assert.Equal(Anomaly{LineNo: 1, Reason: "Unknown line type", Raw: "[section aap noot mies"}, anomalies[0])
	})

	t.Run("section with an invalid number of approvers", func(t *testing.T) {
		content := `[section][999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999]`
		codeOwnersLines, anomalies := Parse(content)

		assert.Equal(1, len(codeOwnersLines))
		assert.Equal(CodeOwnersLine{
			Type:                "section-heading",
			LineNo:              1,
			Raw:                 content,
			SectionOptional:     false,
			SectionName:         "section",
			SectionMinApprovers: 0,
		}, codeOwnersLines[0])

		assert.Equal(0, len(anomalies))
	})

	t.Run("rule", func(t *testing.T) {
		content := `*    @user1 @user2`
		codeOwnersLines, anomalies := Parse(content)

		assert.Equal(1, len(codeOwnersLines))
		assert.Equal(CodeOwnersLine{
			Type:            "rule",
			LineNo:          1,
			Raw:             content,
			SectionOptional: false,
			RulePattern:     "*",
			Spaces:          "    ",
			Owners: []Owner{
				{Name: "@user1", Type: "user-or-group"},
				{Name: "@user2", Type: "user-or-group"},
			},
			InlineComment: "",
		}, codeOwnersLines[0])

		assert.Equal(0, len(anomalies))
	})
	t.Run("rule with an inline comment", func(t *testing.T) {
		content := `*    @user1 @user2 # This is a comment`
		codeOwnersLines, anomalies := Parse(content)

		assert.Equal(CodeOwnersLine{
			Type:            "rule",
			LineNo:          1,
			Raw:             content,
			SectionOptional: false,
			RulePattern:     "*",
			Spaces:          "    ",
			Owners: []Owner{
				{Name: "@user1", Type: "user-or-group"},
				{Name: "@user2", Type: "user-or-group"},
			},
			InlineComment: " This is a comment",
		}, codeOwnersLines[0])

		assert.Equal(0, len(anomalies))
	})

	t.Run("rule without owners", func(t *testing.T) {
		content := `*`
		codeOwnersLines, anomalies := Parse(content)

		assert.Equal(1, len(codeOwnersLines))
		assert.Equal(CodeOwnersLine{
			Type:   "unknown",
			LineNo: 1,
			Raw:    content,
		}, codeOwnersLines[0])

		assert.Equal(1, len(anomalies))
		assert.Equal(Anomaly{LineNo: 1, Reason: "Unknown line type", Raw: "*"}, anomalies[0])
	})

	t.Run("rule with classified owners", func(t *testing.T) {
		content := `*    @user1 invalid invalid-too@ email@address.org`
		codeOwnersLines, anomalies := Parse(content)

		assert.Equal(1, len(codeOwnersLines))
		assert.Equal(CodeOwnersLine{
			Type:            "rule",
			LineNo:          1,
			Raw:             content,
			SectionOptional: false,
			RulePattern:     "*",
			Spaces:          "    ",
			Owners: []Owner{
				{Name: "@user1", Type: "user-or-group"},
				{Name: "invalid", Type: "invalid"},
				{Name: "invalid-too@", Type: "invalid"},
				{Name: "email@address.org", Type: "e-mail"},
			},
			InlineComment: "",
		}, codeOwnersLines[0])

		assert.Equal([]Anomaly{
			{LineNo: 1, Reason: "Invalid user 'invalid'", Raw: "*    @user1 invalid invalid-too@ email@address.org"},
			{LineNo: 1, Reason: "Invalid user 'invalid-too@'", Raw: "*    @user1 invalid invalid-too@ email@address.org"},
		}, anomalies)
	})

	t.Run("rule without owners in the context of a section with", func(t *testing.T) {
		content := "^[section] @some_group\n*"
		codeOwnersLines, anomalies := Parse(content)

		assert.Equal(2, len(codeOwnersLines))
		assert.Equal(CodeOwnersLine{
			Type:                "section-heading",
			LineNo:              1,
			Raw:                 "^[section] @some_group",
			SectionOptional:     true,
			SectionName:         "section",
			SectionMinApprovers: 0,
			Spaces:              " ",
			Owners:              []Owner{{Name: "@some_group", Type: "user-or-group"}},
		}, codeOwnersLines[0])
		assert.Equal(CodeOwnersLine{
			Type:                "rule",
			LineNo:              2,
			Raw:                 "*",
			SectionOptional:     false,
			SectionName:         "",
			SectionMinApprovers: 0,
			RulePattern:         "*",
			RuleSection:         "section",
			Spaces:              "",
			Owners:              nil,
		}, codeOwnersLines[1])

		assert.Equal(0, len(anomalies))
	})

	t.Run("rule without owners in the context of a section with", func(t *testing.T) {
		content := "^[section] @some_group\n*"
		codeOwnersLines, anomalies := Parse(content)

		assert.Equal(2, len(codeOwnersLines))
		assert.Equal(CodeOwnersLine{
			Type:                "section-heading",
			LineNo:              1,
			Raw:                 "^[section] @some_group",
			SectionOptional:     true,
			SectionName:         "section",
			SectionMinApprovers: 0,
			Spaces:              " ",
			Owners:              []Owner{{Name: "@some_group", Type: "user-or-group"}},
		}, codeOwnersLines[0])
		assert.Equal(CodeOwnersLine{
			Type:                "rule",
			LineNo:              2,
			Raw:                 "*",
			SectionOptional:     false,
			SectionName:         "",
			SectionMinApprovers: 0,
			RulePattern:         "*",
			RuleSection:         "section",
			Spaces:              "",
			Owners:              nil,
		}, codeOwnersLines[1])

		assert.Equal(0, len(anomalies))
	})

	t.Run("rule without owners in the context of a section with valid and invalid owners", func(t *testing.T) {
		content := "^[section] @some_group invalid_group @valid_group\n*"
		codeOwnersLines, anomalies := Parse(content)

		assert.Equal(2, len(codeOwnersLines))
		assert.Equal(CodeOwnersLine{
			Type:                "section-heading",
			LineNo:              1,
			Raw:                 "^[section] @some_group invalid_group @valid_group",
			SectionOptional:     true,
			SectionName:         "section",
			SectionMinApprovers: 0,
			Spaces:              " ",
			Owners: []Owner{
				{Name: "@some_group", Type: "user-or-group"},
				{Name: "invalid_group", Type: "invalid"},
				{Name: "@valid_group", Type: "user-or-group"},
			},
		}, codeOwnersLines[0])
		assert.Equal(CodeOwnersLine{
			Type:                "rule",
			LineNo:              2,
			Raw:                 "*",
			SectionOptional:     false,
			SectionName:         "",
			SectionMinApprovers: 0,
			RulePattern:         "*",
			RuleSection:         "section",
			Spaces:              "",
			Owners:              nil,
		}, codeOwnersLines[1])

		assert.Equal(1, len(anomalies))
		assert.Equal([]Anomaly{
			{LineNo: 1, Reason: "Invalid user 'invalid_group'", Raw: "^[section] @some_group invalid_group @valid_group"},
		}, anomalies)
	})

	t.Run("rule without owners in the context of a section also without owners", func(t *testing.T) {
		content := "^[section]\n*"
		codeOwnersLines, anomalies := Parse(content)

		assert.Equal(2, len(codeOwnersLines))
		assert.Equal(CodeOwnersLine{
			Type:                "section-heading",
			LineNo:              1,
			Raw:                 "^[section]",
			SectionOptional:     true,
			SectionName:         "section",
			SectionMinApprovers: 0,
			Spaces:              "",
			Owners:              nil,
		}, codeOwnersLines[0])
		assert.Equal(CodeOwnersLine{
			Type:        "unknown",
			LineNo:      2,
			Raw:         "*",
			RuleSection: "section",
		}, codeOwnersLines[1])

		assert.Equal(1, len(anomalies))
		assert.Equal([]Anomaly{
			{LineNo: 2, Reason: "Unknown line type", Raw: "*"},
		}, anomalies)
	})

	t.Run("rule without owners in the context of a section without valid owners", func(t *testing.T) {
		content := "^[section] invalid_user\n*"
		codeOwnersLines, anomalies := Parse(content)

		assert.Equal(2, len(codeOwnersLines))
		assert.Equal(CodeOwnersLine{
			Type:                "section-heading",
			LineNo:              1,
			Raw:                 "^[section] invalid_user",
			SectionOptional:     true,
			SectionName:         "section",
			SectionMinApprovers: 0,
			Spaces:              " ",
			Owners:              []Owner{{Name: "invalid_user", Type: "invalid"}},
		}, codeOwnersLines[0])
		assert.Equal(CodeOwnersLine{
			Type:        "unknown",
			LineNo:      2,
			Raw:         "*",
			RuleSection: "section",
		}, codeOwnersLines[1])

		assert.Equal(2, len(anomalies))
		assert.Equal([]Anomaly{
			{LineNo: 1, Reason: "Invalid user 'invalid_user'", Raw: "^[section] invalid_user"},
			{LineNo: 2, Reason: "Unknown line type", Raw: "*"},
		}, anomalies)
	})

}
