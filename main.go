package main

import (
	"flag"
	"fmt"
	"os"
)

const VERSION = "0.0.1"

func formatValid(format string) bool {
	return format == "codeowners" || format == "json"
}

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
		fmt.Fprintf(flag.CommandLine.Output(), "Usage of %s:\n", os.Args[0])
		flag.PrintDefaults()
	}

	formatPtr := flag.String("format", "codeowners", "output format")
	versionPtr := flag.Bool("version", false, "print version")
	validatePtr := flag.String("validate", "fail", "fail: exit on syntax errors, warn: print syntax errors & continue, skip: ignore syntax errors")
	teamMapPtr := flag.String("team-map", "", "file path to team map json file")

	flag.Parse()

	if *versionPtr {
		fmt.Println(VERSION)
		os.Exit(0)
	}

	if len(flag.Args()) <= 0 {
		fmt.Fprintln(flag.CommandLine.Output(), "Please provide a CODEOWNERS file")
		os.Exit(1)
	}
	if !formatValid(*formatPtr) {
		fmt.Fprintln(flag.CommandLine.Output(), "Invalid format option. Valid options: codeowners, json")
		os.Exit(1)
	}
	if !validateValid(*validatePtr) {
		fmt.Fprintln(flag.CommandLine.Output(), "Invalid validate option. Valid options: fail, warn, skip")
		os.Exit(1)
	}

	bytes, readFileError := os.ReadFile(flag.Arg(0))

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

	formatted, formatError := FormatCST(codeOwnersLines, *formatPtr)
	if formatError != nil {
		fmt.Fprintln(flag.CommandLine.Output(), formatError)
		os.Exit(1)
	}
	fmt.Print(formatted)
}
