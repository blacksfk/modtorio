/*
Package to search mods.factorio.com.
Authentication with factorio.com is requried to download mods, and the user
is prompted prior to downloading.
*/
package main

import (
	"flag"
	"fmt"
	"os"
)

const (
	MIN_ARGS = 2          // including program name
	FLAG_DIR = "dir"      // name of the working directory flag
	FLAG_VER = "factorio" // name of the factorio version flag
	RAW_FLAG_DIR = "--" + FLAG_DIR
	RAW_FLAG_VER = "--" + FLAG_VER
	DEF_DIR  = "./"       // default to the current directory
	DEF_VER  = "*"        // default to match any version
)

type command struct {
	name string         // command string
	min  int            // minimum args for the command
	fn   func([]string) // function to handle the command
}

// compare a commandline string to the command's name
func (c command) cmp(name string) bool {
	return c.name == name
}

// command definitions
var commands []command = []command{
	{"search", 1, search},
	{"download", 1, download},
	{"update", 0, update},
	{"enable", 1, enable},
	{"disable", 1, disable},
	{"list", 0, list},
}

// package wide flag values
var dir, fVer string

// main function.
// command argument handling
func main() {
	// define flags
	flag.StringVar(&dir, FLAG_DIR, DEF_DIR, "Working directory")
	flag.StringVar(&fVer, FLAG_VER, DEF_VER, "Factorio version")

	// parse the flags
	flag.Parse()

	// validate all arguments and extract the command its options
	cmd, options, e := validate(os.Args)

	if e != nil {
		fmt.Println(e)

		return
	}

	optionCount := len(options)
	found := false

	// loop through all the commands and check for a match
	for _, c := range commands {
		if c.cmp(cmd) {
			if optionCount >= c.min {
				found = true
				c.fn(options)

				break
			} else {
				// not enough args given for the command
				fmt.Printf("Not enough arguments for command: %s. Minimum: %d, Found: %d\n", c.name, c.min, optionCount)
				help()

				return
			}
		}
	}

	if !found {
		fmt.Printf("Invalid command: %s\n", cmd)
		help()
	}
}

// validate the arguments, extract the command and any options
func validate(args []string) (string, []string, error) {
	var options []string

	// check if at least MIN_ARGS were provided
	// eg. modtorio list
	if len(args) < MIN_ARGS {
		return "", options, fmt.Errorf("No command specified\n")
	}

	skip := 0
	argv := args[1:] // skip program name

	for _, arg := range argv {
		if arg == RAW_FLAG_DIR || arg == RAW_FLAG_VER {
			// skip 2; one for flag, one for the flag's value
			skip += 2
		}
	}

	argc := len(argv)

	// check if something was provided after the flags
	if argc <= skip {
		return "", options, fmt.Errorf("Invalid arguments: %v", argv)
	}

	// check if options were provided after the command string
	// so skip an extra
	if skip + 1 < argc {
		options = argv[skip+1:]
	}

	return argv[skip], options, nil
}

// display help for program usage
func help() {
	fmt.Printf("usage: modtorio [...flags] <command> [...options] <arguments>\n\n")

	fmt.Printf("Flags:\n")
	fmt.Printf("\t--dir\tSpecify the working directory for commands that interact with modlist.json. Leave blank if the current directory contains modlist.json or you want modlist.json to be created in the current directory.\n")
	fmt.Printf("\t--factorio\tSpecify the factorio version to compare releases against. Defaults to the latest version.\n\n")

	fmt.Printf("Commands:\n")

	// search command
	fmt.Printf("search\n")
	fmt.Printf("\tSearch for a mod. The argument is compiled as a regular expression.\n")
	fmt.Printf("\tOptions:\n")
	fmt.Printf("\t\t--tag\tSearch for mods based on a tag\n")
	fmt.Printf("\t\t--owner\tSearch for mods created by a user\n")
	fmt.Printf("\tExamples:\n")
	fmt.Printf("\t\tmodtorio search ^bob\n")
	fmt.Printf("\t\tmodtorio search --tag general\n")
	fmt.Printf("\t\tmodtorio search --owner py*\n")

	// download command
	fmt.Printf("download\n")
	fmt.Printf("\tDownload any number of mods. Must be listed by the mod name.\n")
	fmt.Printf("\tExamples:\n")
	fmt.Printf("\t\tmodtorio download bobinserters miniloader pyhightech\n")
	fmt.Printf("\t\tmodtorio --factorio 0.17 --dir ~/.config/factorio/mods download bobinserters helicopters\n")

	// update command
	fmt.Printf("update\n")
	fmt.Printf("\tUpdate all mods to their latest release for the factorio version (if specified).\n")
	fmt.Printf("\tExamples:\n")
	fmt.Printf("\t\tmodtorio update\n")
	fmt.Printf("\t\tmodtorio --factorio 0.18 update\n")
	fmt.Printf("\t\tmodtorio --factorio 0.18 --dir ~/.config/factorio/mods update\n")

	// enable command
	fmt.Printf("enable\n")
	fmt.Printf("\tEnable mods. Arguments are compiled as regular expressions.\n")
	fmt.Printf("\tExamples:\n")
	fmt.Printf("\t\tmodtorio enable bob* pyhightech ^angel\n")
	fmt.Printf("\t\tmodtorio --dir ~/.config/factorio/mods enable bob*\n")

	// disable command
	fmt.Printf("disable\n")
	fmt.Printf("\tDisable mods. Arguments are compiled as regular expressions.\n")
	fmt.Printf("\tExamples:\n")
	fmt.Printf("\t\tmodtorio disable bob* pyhightech ^angel\n")
	fmt.Printf("\t\tmodtorio --dir ~/.config/factorio/mods disable bob*\n")

	// list command
	fmt.Printf("list\n")
	fmt.Printf("\tList mods.\n")
	fmt.Printf("\tOptions:\n")
	fmt.Printf("\t\t--all\tList all installed mods (default)\n")
	fmt.Printf("\t\t--enabled\tList all enabled mods\n")
	fmt.Printf("\t\t--disabled\tList all disabled mods\n")
	fmt.Printf("\tExamples:\n")
	fmt.Printf("\t\tmodtorio list\n")
	fmt.Printf("\t\tmodtorio list --all\n")
	fmt.Printf("\t\tmodtorio list --enabled\n")
	fmt.Printf("\t\tmodtorio list --disabled\n")
	fmt.Printf("\t\tmodtorio --dir ~/.config/factorio/mods list\n")
}
