package main

import (
	"flag"
	"fmt"
	"os"
)

const VERSION = "0.0.1"

func validateValid(validate string) bool {
	var validValidateOptions = map[string]bool{
		"fail": true,
		"warn": true,
		"skip": true,
	}
	return validValidateOptions[validate]
}

func main() {

	flag.Usage = func() {
		fmt.Fprint(flag.CommandLine.Output(), "Usage: vcodeowners [options]\n\n")
		fmt.Fprint(flag.CommandLine.Output(), "Merges a VIRTUAL-CODEOWNERS.txt and a virtual-teams.yml into CODEOWNERS\n\n")
		flag.PrintDefaults()
	}

	/*
			Options from virtual-code-owners not yet implemented:
		    --emitLabeler                 Whether or not to emit a labeler.yml to be
		                                  used with actions/labeler
		                                  (default: false)
		    --labelerLocation [file-name] The location of the labeler.yml file
		                                  (default: ".github/labeler.yml")
		    --dryRun                      Just validate inputs, don't generate
		                                  outputs (default: false)
	*/
	versionPtr := flag.Bool("version", false, "output the version number")
	virtualCodeOwnersPtr := flag.String("virtualCodeOwners", ".github/VIRTUAL-CODEOWNERS.txt", "A CODEOWNERS file with team names in them that are defined in a virtual teams file")
	teamMapPtr := flag.String("virtualTeams", ".github/virtual-teams.json", "A JSON file listing teams and their members")
	codeOwnersPtr := flag.String("codeOwners", ".github/CODEOWNERS", "The CODEOWNERS file to merge the virtual teams into")
	validatePtr := flag.String("validate", "fail", "fail: exit on syntax errors, warn: print syntax errors & continue, skip: ignore syntax errors")

	flag.Parse()

	if *versionPtr {
		fmt.Println(VERSION)
		os.Exit(0)
	}
	if !validateValid(*validatePtr) {
		fmt.Fprintln(flag.CommandLine.Output(), "Invalid validate option. Valid options: fail, warn, skip")
		os.Exit(1)
	}

	bytes, readFileError := os.ReadFile(*virtualCodeOwnersPtr)

	if readFileError != nil {
		fmt.Fprintln(flag.CommandLine.Output(), readFileError)
		os.Exit(1)
	}

	codeOwnersLines, syntaxErrors := Parse(string(bytes))

	if len(syntaxErrors) > 0 && (*validatePtr != "skip") {
		fmt.Fprintln(flag.CommandLine.Output(), FormatAnomaliesAsText(syntaxErrors))
		if *validatePtr == "fail" {
			os.Exit(1)
		}
	}

	if *teamMapPtr != "" {
		teamMapBytes, teamMapReadError := os.ReadFile(*teamMapPtr)
		if teamMapReadError != nil {
			fmt.Fprintln(flag.CommandLine.Output(), teamMapReadError)
			os.Exit(1)
		}
		teamMap, teamMapParseError := ParseTeamMap(string(teamMapBytes))
		if teamMapParseError != nil {
			fmt.Fprintln(flag.CommandLine.Output(), teamMapParseError)
			os.Exit(1)
		}
		codeOwnersLines = ApplyTeamMap(codeOwnersLines, teamMap)
	}

	formatted, formatError := FormatCST(codeOwnersLines, "codeowners")
	if formatError != nil {
		fmt.Fprintln(flag.CommandLine.Output(), formatError)
		os.Exit(1)
	}
	writeError := os.WriteFile(*codeOwnersPtr, []byte(formatted), 0644)
	if writeError != nil {
		fmt.Fprintln(flag.CommandLine.Output(), writeError)
		os.Exit(1)
	}
	fmt.Fprintf(flag.CommandLine.Output(), "\nWrote '%s'\n\n", *codeOwnersPtr)
}
