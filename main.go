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
	FLAG_DIR = "dir"      // name of the working directory flag
	FLAG_VER = "factorio" // name of the factorio version flag
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
	{"list-enabled", 0, listEnabled},
	{"list-disabled", 0, listDisabled},
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

	// subtract 1 to compensate for program name
	argv := os.Args[1:]
	skip := 0

	// loop and subtract argc for each flag argument
	for _, arg := range argv {
		if arg == "--"+FLAG_DIR || arg == "--"+FLAG_VER {
			// skip 2; one for flag, one for the flag's value
			skip += 2
		}
	}

	cmd := argv[skip]
	// subtract the command arg and flag args
	argc := len(argv) - skip - 1
	found := false

	// loop through all the commands and check for a match
	for _, c := range commands {
		if c.cmp(cmd) {
			if argc >= c.min {
				found = true
				c.fn(argv)

				break
			} else {
				// not enough args given for the command
				fmt.Printf("Not enough arguments for command: %s. Minimum: %d, Found: %d\n", c.name, c.min, argc)
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

// display help for program usage
func help() {
	fmt.Printf("usage: modtorio [...flags] <command> [...options] <arguments>\n\n")

	fmt.Printf("Flags:\n")
	fmt.Printf("\t--dir\tSpecify the working directory for commands that interact with modlist.json. Leave blank if the current directory contains modlist.json or you want modlist.json to be created in the current directory.\n")
	fmt.Printf("\t--factorio\tSpecify the factorio version to compare releases against. Defaults to the latest version\n\n")

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

	// list-enabled command
	fmt.Printf("list-enabled\n")
	fmt.Printf("\tList all enabled mods.\n")
	fmt.Printf("\tExamples:\n")
	fmt.Printf("\t\tmodtorio list-enabled\n")
	fmt.Printf("\t\tmodtorio --dir ~/.config/factorio/mods list-enabled\n")

	// list-disabled command
	fmt.Printf("list-disabled\n")
	fmt.Printf("\tList all disabled mods.\n")
	fmt.Printf("\tExamples:\n")
	fmt.Printf("\t\tmodtorio list-disabled\n")
	fmt.Printf("\t\tmodtorio --dir ~/.config/factorio/mods list-disabled\n")
}
