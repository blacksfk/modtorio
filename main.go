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
	"os"
	"strings"
	"syscall"
)

type credentials struct {
	username, password, token string
}

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
// authenticates with the factorio.com web auth API
func download() {
	creds, err := promptForCreds()

	if err != nil {
		fmt.Println(err)

		return
	}

	creds.token, err = login(creds.username, creds.password)

	if err != nil {
		fmt.Println(err)

		return
	}

	fmt.Printf("token: %s\n", creds.token)
}

// prompt the user for their login credentials
func promptForCreds() (*credentials, error) {
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

	return &credentials{username: strings.TrimRight(username, "\n"), password: string(bytes)}, nil
}

// display help for program usage
func help() {
	fmt.Println("Usage:")
	fmt.Printf("\tmodtorio <command> [...options] <arguments>\n")
	fmt.Printf("\tCommands: search, download, help\n")
}
