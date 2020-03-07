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
	fmt.Println("Usage:")
	fmt.Printf("\tmodtorio <command> [...options] <arguments>\n")
	fmt.Printf("\tCommands: search, download, help\n")
}
