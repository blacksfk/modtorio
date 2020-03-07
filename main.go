/*
Package to search mods.factorio.com. Valid commands: search, download, help.
Authentication with factorio.com is requried to download mods, and the user
is prompted prior to downloading.
*/
package main

import (
	"fmt"
	"os"
)

const (
	MIN_ARGS = 3
)

// main function.
// command argument handling
func main() {
	// check if at least one command was provided
	if len(os.Args) < MIN_ARGS {
		fmt.Println("Not enough arguments")
		help()

		return
	}

	cmd := os.Args[1]

	// determine what the user wants to do
	switch cmd {
	case "search":
		search(os.Args[2:])
	case "download":
		download(os.Args[2:])
	case "-h":
		fallthrough
	case "--help":
		fallthrough
	case "help":
		help()
	default:
		fmt.Printf("Invalid command: %s\n", cmd)
	}
}

// display help for program usage
func help() {
	fmt.Printf("usage: modtorio [...flags] <command> [...options] <arguments>\n\n")

	fmt.Printf("Flags:\n")
	fmt.Printf("\t--dir\tSpecify the working directory for commands that interact with modlist.json. Leave blank if the current directory contains modlist.json\n")
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
