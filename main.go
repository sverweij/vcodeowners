package main

import (
	_ "embed"
	"flag"
	"fmt"
	"io"
	"os"
)

const VERSION = "0.0.1"

//go:embed headers/codeowners.txt
var codeOwnersHeaderComment []byte

//go:embed headers/labeler.txt
var labelerHeaderComment []byte

func validateValid(validate string) bool {
	var validValidateOptions = map[string]bool{
		"fail": true,
		"warn": true,
		"skip": true,
	}
	return validValidateOptions[validate]
}

type cliOptionsType struct {
	version           *bool
	virtualCodeOwners *string
	teamMap           *string
	codeOwners        *string
	validate          *string
	dryRun            *bool
	emitLabeler       *bool
	labelerLocation   *string
	json              *bool
}

const EXIT_CODE_ERROR = 1

func getOptions(stderr io.Writer) cliOptionsType {
	flag.Usage = func() {
		fmt.Fprint(stderr, "Usage: vcodeowners [options]\n\n")
		fmt.Fprint(stderr, "Merges a VIRTUAL-CODEOWNERS.txt and a virtual-teams.json into CODEOWNERS\n\n")
		flag.PrintDefaults()
	}

	cliOptions := cliOptionsType{
		version:           flag.Bool("version", false, "output the version number"),
		virtualCodeOwners: flag.String("virtualCodeOwners", ".github/VIRTUAL-CODEOWNERS.txt", "A CODEOWNERS file with team names in them that are defined in a virtual teams file"),
		teamMap:           flag.String("virtualTeams", ".github/virtual-teams.json", "A JSON file listing teams and their members"),
		codeOwners:        flag.String("codeOwners", ".github/CODEOWNERS", "The CODEOWNERS file to merge the virtual teams into"),
		validate:          flag.String("validate", "fail", "fail: exit on syntax errors, warn: print syntax errors & continue, skip: ignore syntax errors"),
		dryRun:            flag.Bool("dryRun", false, "Just validate inputs, don't generate outputs"),
		emitLabeler:       flag.Bool("emitLabeler", false, "Whether or not to emit a labeler.yml to be used with actions/labeler"),
		labelerLocation:   flag.String("labelerLocation", ".github/labeler.yml", "The location of the labeler.yml file"),
		json:              flag.Bool("json", false, "Output JSON to stdout (in addition to writing CODEOWNERS)"),
	}

	flag.Parse()
	return cliOptions
}

func cli(options cliOptionsType) (string, error) {
	returnMessage := ""

	if *options.version {
		return VERSION, nil
	}
	if !validateValid(*options.validate) {
		return "",
			fmt.Errorf("invalid validate option '%s'; valid options: fail, warn, skip", *options.validate)
	}

	bytes, readFileError := os.ReadFile(*options.virtualCodeOwners)

	if readFileError != nil {
		return "", readFileError
	}

	codeOwnersLines, syntaxErrors := Parse(string(bytes))

	if len(syntaxErrors) > 0 && (*options.validate != "skip") {
		if *options.validate == "fail" {
			return "", fmt.Errorf("%s", FormatAnomaliesAsText(syntaxErrors))
		}
		returnMessage = returnMessage + FormatAnomaliesAsText(syntaxErrors)
	}

	teamMap := map[string][]string{}

	if *options.teamMap != "" {
		teamMapBytes, teamMapReadError := os.ReadFile(*options.teamMap)
		if teamMapReadError != nil {
			return "", teamMapReadError
		}
		var teamMapParseError error
		teamMap, teamMapParseError = ParseTeamMap(string(teamMapBytes))
		if teamMapParseError != nil {
			return "", teamMapParseError
		}
	}
	transformedCodeOwnersLines := ApplyTeamMap(codeOwnersLines, teamMap)

	formatted, formatError := FormatCSTAsCodeOwners(transformedCodeOwnersLines, string(codeOwnersHeaderComment))
	if formatError != nil {
		return "", formatError
	}

	if !*options.dryRun {
		writeError := os.WriteFile(*options.codeOwners, []byte(formatted), 0644)
		if writeError != nil {
			return "", writeError
		}
		returnMessage = returnMessage + fmt.Sprintf("\nWrote '%s'\n", *options.codeOwners)
		if *options.emitLabeler {
			labelerFormatted, labelerFormatError := FormatCSTAsLabelerYML(codeOwnersLines, teamMap, string(labelerHeaderComment))
			if labelerFormatError != nil {
				return "", labelerFormatError
			}
			labelerWriteError := os.WriteFile(*options.labelerLocation, []byte(labelerFormatted), 0644)
			if labelerWriteError != nil {
				return "", labelerWriteError
			}
		}
	} else {
		returnMessage = returnMessage + fmt.Sprintf("\nWrote '%s' (dry run)\n\n", *options.codeOwners)
	}

	if *options.json {
		jsonFormatted, jsonFormatError := FormatCSTAsJSON(transformedCodeOwnersLines)
		if jsonFormatError != nil {
			return "", jsonFormatError
		}
		fmt.Println(jsonFormatted)
	}
	return returnMessage, nil
}

func main() {
	cliOptions := getOptions(flag.CommandLine.Output())
	message, error := cli(cliOptions)

	if error != nil {
		fmt.Fprintln(flag.CommandLine.Output(), error.Error())
		os.Exit(EXIT_CODE_ERROR)
	} else {
		fmt.Fprintln(flag.CommandLine.Output(), message)
	}
}
