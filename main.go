/*
Package to search mods.factorio.com.
Authentication with factorio.com is requried to download mods, and the user
is prompted prior to downloading.
*/
package main

import (
	"flag"
	"fmt"

	"github.com/blacksfk/modtorio/common"
)

const (
	FLAG_DIR     = "dir"            // name of the working directory flag
	FLAG_VER     = "factorio"       // name of the factorio version flag
	DEF_DIR      = "./"             // default to the current directory
	DEF_VER      = common.MATCH_ANY // default to match any version
	CMD_SEARCH   = "search"
	CMD_DOWNLOAD = "download"
	CMD_UPDATE   = "update"
	CMD_ENABLE   = "enable"
	CMD_DISABLE  = "disable"
	CMD_LIST     = "list"
	CMD_HELP     = "help"
)

type Command struct {
	name string                               // command string
	min  int                                  // minimum args for the command
	fn   func(*ModtorioFlags, []string) error // function to handle the command
}

// compare a commandline string to the command's name
func (c Command) cmp(name string) bool {
	return c.name == name
}

type ModtorioFlags struct {
	dir      string
	factorio *common.Semver
}

// main function.
// command argument handling
func main() {
	// define flags
	var strVer string
	flags := &ModtorioFlags{}

	flag.StringVar(&flags.dir, FLAG_DIR, DEF_DIR, "Working directory")
	flag.StringVar(&strVer, FLAG_VER, DEF_VER, "Factorio version")

	// parse the flags
	flag.Parse()

	semver, e := common.NewSemver(strVer)

	if e != nil {
		fmt.Println("Factorio version flag:", e)

		return
	}

	flags.factorio = semver

	// validate remaining arguments
	argv := flag.Args()
	argc := len(argv)

	if argc == 0 {
		fmt.Println("No command specified")

		return
	}

	var options []string
	cmd := argv[0]

	if argc > 1 {
		// options or flags for command
		options = argv[1:]
	}

	e = matchAndRun(cmd, flags, options)

	if e != nil {
		fmt.Println(e)

		return
	}
}

func matchAndRun(name string, flags *ModtorioFlags, options []string) error {
	optionCount := len(options)
	commands := []Command{
		{CMD_SEARCH, 1, search},
		{CMD_DOWNLOAD, 1, download},
		{CMD_UPDATE, 0, update},
		{CMD_ENABLE, 1, enable},
		{CMD_DISABLE, 1, disable},
		{CMD_LIST, 0, list},
		{CMD_HELP, 0, help},
	}

	// loop through all defined commands
	for _, cmd := range commands {
		if cmd.cmp(name) {
			// match found, compare minimum arguments
			if optionCount >= cmd.min {
				// minimum arguments found, call the handler
				return cmd.fn(flags, options)
			} else {
				// argument count does not meet the minimum
				return fmt.Errorf("Not enough arguments for command: %s. Minimum: %d, Found: %d", cmd.name, cmd.min, optionCount)
			}
		}
	}

	// no match found
	return fmt.Errorf("Invalid command: %s", name)
}
