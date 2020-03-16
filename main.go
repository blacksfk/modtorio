/*
Package to search mods.factorio.com.
Authentication with factorio.com is requried to download mods, and the user
is prompted prior to downloading.
*/
package main

import (
	"flag"
	"fmt"
	"modtorio/common"
	"os"
)

const (
	MIN_ARGS     = 2          // including program name
	FLAG_DIR     = "dir"      // name of the working directory flag
	FLAG_VER     = "factorio" // name of the factorio version flag
	RAW_FLAG_DIR = "--" + FLAG_DIR
	RAW_FLAG_VER = "--" + FLAG_VER
	DEF_DIR      = "./"             // default to the current directory
	DEF_VER      = common.MATCH_ANY // default to match any version
	CMD_SEARCH   = "search"
	CMD_DOWNLOAD = "download"
	CMD_UPDATE   = "update"
	CMD_ENABLE   = "enable"
	CMD_DISABLE  = "disable"
	CMD_LIST     = "list"
)

type Command struct {
	name string               // command string
	min  int                  // minimum args for the command
	fn   func([]string) error // function to handle the command
}

// compare a commandline string to the command's name
func (c Command) cmp(name string) bool {
	return c.name == name
}

// command definitions
var commands []Command = []Command{
	{CMD_SEARCH, 1, search},
	{CMD_DOWNLOAD, 1, download},
	{CMD_UPDATE, 0, update},
	{CMD_ENABLE, 1, enable},
	{CMD_DISABLE, 1, disable},
	{CMD_LIST, 0, list},
	{"help", 0, help},
}

// package wide flag values
type ModtorioFlags struct {
	dir      string
	factorio *common.Semver
}

var stringVer string
var FLAGS ModtorioFlags = ModtorioFlags{}

// main function.
// command argument handling
func main() {
	// define flags
	flag.StringVar(&FLAGS.dir, FLAG_DIR, DEF_DIR, "Working directory")
	flag.StringVar(&stringVer, FLAG_VER, DEF_VER, "Factorio version")

	// parse the flags
	flag.Parse()

	semver, e := common.NewSemver(stringVer)

	if e != nil {
		fmt.Println("Factorio version flag:", e)

		return
	}

	FLAGS.factorio = semver

	// validate all arguments and extract the command and its options
	cmd, options, e := validate(os.Args)

	if e != nil {
		fmt.Println("Validation failed:", e)

		return
	}

	e = matchCommand(cmd, options)

	if e != nil {
		fmt.Println(e)

		return
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
	if skip+1 < argc {
		options = argv[skip+1:]
	}

	return argv[skip], options, nil
}

func matchCommand(name string, options []string) error {
	optionCount := len(options)

	// loop through all defined commands
	for _, cmd := range commands {
		if cmd.cmp(name) {
			// match found, compare minimum arguments
			if optionCount >= cmd.min {
				// minimum arguments found, call the handler
				return cmd.fn(options)
			} else {
				// argument count does not meet the minimum
				return fmt.Errorf("Not enough arguments for command: %s. Minimum: %d, Found: %d", cmd.name, cmd.min, optionCount)
			}
		}
	}

	// no match found
	return fmt.Errorf("Invalid command: %s", name)
}
