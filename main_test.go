package main

import (
	"os"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

var versionPattern = regexp.MustCompile(`^[0-9]+\.[0-9]+\.[0-9]+$`)

func initCliOptions() cliOptionsType {

	version := false
	virtualCodeOwners := ".github/VIRTUAL-CODEOWNERS.txt"
	teamMap := ".github/virtual-teams.json"
	codeOwners := ".github/CODEOWNERS"
	validate := "fail"
	dryRun := false
	emitLabeler := false
	labelerLocation := ".github/labeler.yml"
	json := false

	return cliOptionsType{
		version:           &version,
		virtualCodeOwners: &virtualCodeOwners,
		teamMap:           &teamMap,
		codeOwners:        &codeOwners,
		validate:          &validate,
		dryRun:            &dryRun,
		emitLabeler:       &emitLabeler,
		labelerLocation:   &labelerLocation,
		json:              &json,
	}
}

func TestCli(t *testing.T) {
	assert := assert.New(t)
	t.Run("--version returns the version", func(t *testing.T) {
		options := initCliOptions()
		showVersion := true
		options.version = &showVersion
		message, error := cli(options)

		assert.Regexp(versionPattern, message)
		assert.Nil(error)
	})

	t.Run("invalid validate option returns an error", func(t *testing.T) {
		options := initCliOptions()
		validateValue := "invalid"
		options.validate = &validateValue
		_, error := cli(options)

		assert.NotNil(error)
		assert.Equal("invalid validate option 'invalid'; valid options: fail, warn, skip", error.Error())
	})

	t.Run("invalid virtualCodeOwners file returns an error", func(t *testing.T) {
		options := initCliOptions()
		nonExistentFile := "non_existent_file.txt"
		options.virtualCodeOwners = &nonExistentFile
		_, error := cli(options)

		assert.NotNil(error)
		assert.Equal("open non_existent_file.txt: no such file or directory", error.Error())
	})

	t.Run("happy day everything", func(t *testing.T) {
		vcoFileName := "delete_me_VIRTUAL_CODEOWNERS.txt"
		teamsFileName := "delete_me_virtual-teams.json"
		coFileName := "delete_me_CODEOWNERS"
		doEmitLabeler := true
		labelerFileName := "delete_me_labeler.yml"
		defer func() {
			os.Remove(vcoFileName)
			os.Remove(coFileName)
			os.Remove(teamsFileName)
			os.Remove(labelerFileName)
		}()

		_, vcoFileCreateError := os.Create(vcoFileName)
		teamsFile, teamsFileCreateError := os.Create(teamsFileName)

		assert.Nil(vcoFileCreateError)
		assert.Nil(teamsFileCreateError)

		teamsFile.WriteString("{}")

		options := initCliOptions()
		options.virtualCodeOwners = &vcoFileName
		options.teamMap = &teamsFileName
		options.codeOwners = &coFileName
		options.emitLabeler = &doEmitLabeler
		options.labelerLocation = &labelerFileName
		foundMessage, error := cli(options)

		assert.Nil(error)
		assert.Equal("\nWrote 'delete_me_CODEOWNERS'\n", foundMessage)
	})

	t.Run("happy day everything --dryRun", func(t *testing.T) {
		vcoFileName := "delete_me_VIRTUAL_CODEOWNERS.txt"
		teamsFileName := "delete_me_virtual-teams.json"
		coFileName := "delete_me_CODEOWNERS_should_not_be_created"
		doEmitLabeler := true
		dryRun := true
		labelerFileName := "delete_me_labeler_should_not_be_created.yml"
		defer func() {
			os.Remove(vcoFileName)
			os.Remove(coFileName)
			os.Remove(teamsFileName)
			os.Remove(labelerFileName)
		}()

		_, vcoFileCreateError := os.Create(vcoFileName)
		teamsFile, teamsFileCreateError := os.Create(teamsFileName)

		assert.Nil(vcoFileCreateError)
		assert.Nil(teamsFileCreateError)

		teamsFile.WriteString("{}")

		options := initCliOptions()
		options.virtualCodeOwners = &vcoFileName
		options.teamMap = &teamsFileName
		options.codeOwners = &coFileName
		options.emitLabeler = &doEmitLabeler
		options.labelerLocation = &labelerFileName
		options.dryRun = &dryRun
		foundMessage, error := cli(options)

		assert.Nil(error)
		_, coFileOpenError := os.Open(coFileName)
		assert.NotNil(coFileOpenError)
		_, labelerFileOpenError := os.Open(labelerFileName)
		assert.NotNil(labelerFileOpenError)
		assert.Equal("\nWrote 'delete_me_CODEOWNERS_should_not_be_created' (dry run)\n\n", foundMessage)
	})
}
