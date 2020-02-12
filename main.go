/*
Package to search mods.factorio.com. Valid commands: search, download, help.
Authentication with factorio.com is requried to download mods, and the user
is prompted prior to downloading.
*/
package main

import (
	"bufio"
	"fmt"
	"golang.org/x/crypto/ssh/terminal"
	"modtorio/credentials"
	"modtorio/request"
	"os"
	"strings"
	"syscall"
)

// main function.
// command argument handling
func main() {
	// check if at least one command was provided
	if len(os.Args) < 3 {
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
// authenticates with the factorio.com web auth API
func download() {
	// first, log the user in
	creds, err := promptForCreds()

	if err != nil {
		fmt.Println(err)

		return
	}

	err = request.Login(creds)

	if err != nil {
		fmt.Println(err)

		return
	}

	// then get the "short" information for each mod and
	// download the archives
	for _, term := range os.Args[2:] {
		mod, err := request.GetMod(term)

		if err != nil {
			// if an error occurred with this search term,
			// print the error and continue
			fmt.Printf("%s: %s\n", term, err)
		} else {
			err = request.DownloadMod(&mod.Releases[len(mod.Releases)-1], creds)

			if err != nil {
				fmt.Printf("%s: %s\n", term, err)
			}
		}
	}
}

// prompt the user for their login credentials
func promptForCreds() (*credentials.Credentials, error) {
	fmt.Println("Please enter your credentials to download from mods.factorio.com")

	stdin := bufio.NewReader(os.Stdin)
	fmt.Print("Username: ")
	username, err := stdin.ReadString('\n')

	if err != nil {
		return nil, err
	}

	// use the terminal package to read the password so it isn't
	// echoed back to the user in plain sight
	fmt.Print("Password: ")
	bytes, err := terminal.ReadPassword(int(syscall.Stdin))

	if err != nil {
		return nil, err
	}

	// insert a linefeed ('\n') so that the next prompt isn't printed on the
	// same line as the password prompt
	fmt.Printf("\n")

	return credentials.NewCredentials(strings.TrimRight(username, "\n"), string(bytes)), nil
}

// display help for program usage
func help() {
	fmt.Println("Usage:")
	fmt.Printf("\tmodtorio <command> [...options] <arguments>\n")
	fmt.Printf("\tCommands: search, download, help\n")
}
