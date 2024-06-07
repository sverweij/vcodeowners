package main

import (
	_ "embed"
	"flag"
	"fmt"
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

func main() {

	flag.Usage = func() {
		fmt.Fprint(flag.CommandLine.Output(), "Usage: vcodeowners [options]\n\n")
		fmt.Fprint(flag.CommandLine.Output(), "Merges a VIRTUAL-CODEOWNERS.txt and a virtual-teams.json into CODEOWNERS\n\n")
		flag.PrintDefaults()
	}

	versionPtr := flag.Bool("version", false, "output the version number")
	virtualCodeOwnersPtr := flag.String("virtualCodeOwners", ".github/VIRTUAL-CODEOWNERS.txt", "A CODEOWNERS file with team names in them that are defined in a virtual teams file")
	teamMapPtr := flag.String("virtualTeams", ".github/virtual-teams.json", "A JSON file listing teams and their members")
	codeOwnersPtr := flag.String("codeOwners", ".github/CODEOWNERS", "The CODEOWNERS file to merge the virtual teams into")
	validatePtr := flag.String("validate", "fail", "fail: exit on syntax errors, warn: print syntax errors & continue, skip: ignore syntax errors")
	dryRunPtr := flag.Bool("dryRun", false, "Just validate inputs, don't generate outputs")
	emitLabelerPtr := flag.Bool("emitLabeler", false, "Whether or not to emit a labeler.yml to be used with actions/labeler")
	labelerLocationPtr := flag.String("labelerLocation", ".github/labeler.yml", "The location of the labeler.yml file")
	jsonPtr := flag.Bool("json", false, "Output JSON to stdout (in addition to writing CODEOWNERS)")

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

	teamMap := map[string][]string{}

	if *teamMapPtr != "" {
		teamMapBytes, teamMapReadError := os.ReadFile(*teamMapPtr)
		if teamMapReadError != nil {
			fmt.Fprintln(flag.CommandLine.Output(), teamMapReadError)
			os.Exit(1)
		}
		var teamMapParseError error
		teamMap, teamMapParseError = ParseTeamMap(string(teamMapBytes))
		if teamMapParseError != nil {
			fmt.Fprintln(flag.CommandLine.Output(), teamMapParseError)
			os.Exit(1)
		}
	}
	transformedCodeOwnersLines := ApplyTeamMap(codeOwnersLines, teamMap)

	formatted, formatError := FormatCSTAsCodeOwners(transformedCodeOwnersLines, string(codeOwnersHeaderComment))
	if formatError != nil {
		fmt.Fprintln(flag.CommandLine.Output(), formatError)
		os.Exit(1)
	}

	if !*dryRunPtr {
		writeError := os.WriteFile(*codeOwnersPtr, []byte(formatted), 0644)
		if writeError != nil {
			fmt.Fprintln(flag.CommandLine.Output(), writeError)
			os.Exit(1)
		}
		fmt.Fprintf(flag.CommandLine.Output(), "\nWrote '%s'\n\n", *codeOwnersPtr)
		if *emitLabelerPtr {
			labelerFormatted, labelerFormatError := FormatCSTAsLabelerYML(codeOwnersLines, teamMap, string(labelerHeaderComment))
			if labelerFormatError != nil {
				fmt.Fprintln(flag.CommandLine.Output(), labelerFormatError)
			}
			labelerWriteError := os.WriteFile(*labelerLocationPtr, []byte(labelerFormatted), 0644)
			if labelerWriteError != nil {
				fmt.Fprintln(flag.CommandLine.Output(), labelerWriteError)
				os.Exit(1)
			}

		}
	} else {
		fmt.Fprintf(flag.CommandLine.Output(), "\nWrote '%s' (dry run)\n\n", *codeOwnersPtr)
	}

	if *jsonPtr {
		jsonFormatted, jsonFormatError := FormatCSTAsJSON(transformedCodeOwnersLines)
		if jsonFormatError != nil {
			fmt.Fprintln(flag.CommandLine.Output(), jsonFormatError)
			os.Exit(1)
		}
		fmt.Println(jsonFormatted)
	}

}
