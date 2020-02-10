/*
Package to search mods.factorio.com. Valid commands: search, download, help.
Authentication with factorio.com is requried to download mods, and the user
is prompted prior to downloading.
 */
package main

import (
	"os"
	"fmt"
)

// main function.
// command argument handling
func main() {
	// check if at least one command was provided
	if len(os.Args) < 2 {
		fmt.Println("Invalid command")
		help()

		return
	}

	cmd := os.Args[1]

	// determine what the user wants to do
	switch cmd {
	case "search":
		search()
	case "download":
		download()
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

// search for one or mods.
func search() {
	//
}

// download one or more mods.
// authenticates with the factorio.com web auth API via user
// password prompts.
func download() {
	fmt.Println("Download funcitonality not implemented")
}

// display help for program usage
func help() {
	fmt.Println("Usage:")
	fmt.Printf("\tmodtorio <command> [...options] <arguments>\n")
	fmt.Printf("\tCommands: search, download, help\n")
}
