/*
Package to search mods.factorio.com. Valid commands: search, download, help.
Authentication with factorio.com is requried to download mods, and the user
is prompted prior to downloading.
 */
package main

import (
	"os"
	"fmt"
	"bufio"
	"syscall"
	"strings"
	"golang.org/x/crypto/ssh/terminal"
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
	fmt.Println("Search functionality not implemented")
}

// download one or more mods.
// authenticates with the factorio.com web auth API via user
// password prompts.
func download() {
	stdin := bufio.NewReader(os.Stdin)

	fmt.Println("Please enter your credentials to download from mods.factorio.com")

	// prompt the user for their username
	fmt.Print("Username: ")
	s, err := stdin.ReadString('\n')

	if err != nil {
		fmt.Println(err)

		return
	}

	// drop the trailing new line char
	username := strings.TrimRight(s, "\n")

	// use terminal to read the password so it isn't
	// echoed back to the user in plaintext
	fmt.Print("Password: ")
	bytes, err := terminal.ReadPassword(int(syscall.Stdin))

	if err != nil {
		fmt.Println(err)

		return
	}

	// insert a newline so that the next prompt isn't printed on the
	// same line as the password prompt, then cast bytes to a string
	// and attempt to log the user in
	fmt.Printf("\n")
	password := string(bytes)
	token, err := login(username, string(password))

	if err != nil {
		fmt.Println(err)

		return
	}

	fmt.Printf("token: %s\n", token)
}

// display help for program usage
func help() {
	fmt.Println("Usage:")
	fmt.Printf("\tmodtorio <command> [...options] <arguments>\n")
	fmt.Printf("\tCommands: search, download, help\n")
}
