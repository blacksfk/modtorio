package main

import (
	"bufio"
	"fmt"
	"golang.org/x/crypto/ssh/terminal"
	"modtorio/api"
	"modtorio/credentials"
	"os"
	"strings"
	"syscall"
)

const (
	FLAG_VERSION = "--factorio"
	MAX_LOGIN_ATTEMPS = 5
)

func download(args []string) {
	latest := true
	length := len(args)
	var version string

	for i := 0; i < length; i++ {
		if args[i] == FLAG_VERSION {
			if i == length-1 {
				// user forgot the version number
				fmt.Println("Version number must follow the version option")
				help()

				return
			}

			latest = false
			version = args[i+1]

			// remove the --version <version> from the args slice
			if i == 0 {
				// version was specified at the beginning
				args = args[i+2:]
			} else if i == length-2 {
				// version was specified at the end
				args = args[:i-1]
			} else {
				// the bastard put it in the middle
				args = append(args[:i-1], args[i+2:]...)
			}

			break
		}
	}

	var creds *credentials.Credentials
	var e error
	attempts := 0

	for ; attempts < MAX_LOGIN_ATTEMPS; attempts++ {
		// prompt the user for their credentials
		creds, e = promptForCreds()

		if e != nil {
			// something went wrong with the input
			fmt.Println(e)

			return
		}

		e = api.Login(creds)

		if e != nil {
			fmt.Println(e)
		} else {
			// logged in successfully
			break
		}
	}

	if attempts >= MAX_LOGIN_ATTEMPS {
		fmt.Println("Maximum login attempts exceeded")

		return
	}

	// get the mod results for each mod
	results, e := api.GetAll(args...)

	if e != nil {
		fmt.Println(e)

		return
	}

	if latest {
		downloadLatest(results, creds)
	} else {
		downloadVersion(results, creds, version)
	}
}

// prompt the user for their login credentials
func promptForCreds() (*credentials.Credentials, error) {
	fmt.Println("Please enter your credentials to download from mods.factorio.com")

	stdin := bufio.NewReader(os.Stdin)
	fmt.Print("Username: ")
	username, e := stdin.ReadString('\n')

	if e != nil {
		return nil, e
	}

	// use the terminal package to read the password so it isn't
	// echoed back to the user in plain sight
	fmt.Print("Password: ")
	bytes, e := terminal.ReadPassword(int(syscall.Stdin))

	if e != nil {
		return nil, e
	}

	// insert a linefeed ('\n') so that the next prompt isn't printed on the
	// same line as the password prompt
	fmt.Printf("\n")

	return credentials.NewCredentials(strings.TrimRight(username, "\n"), string(bytes)), nil
}

// download the latest release of a mod
func downloadLatest(results []*api.Result, creds *credentials.Credentials) {
	for _, result := range results {
		// latest release should be the last in the array
		e := result.Releases[len(result.Releases)-1].Download(creds)

		if e != nil {
			// don't return on any error
			// continue processing rest of the downloads
			fmt.Println(e)
		}
	}
}

// download a the latest release of the specified FACTORIO version
func downloadVersion(results []*api.Result, creds *credentials.Credentials, version string) {
	for _, result := range results {
		found := false

		// start from the end to find the latest release for
		// the specified factorio version
		for i := len(result.Releases); i >= 0; i-- {
			if result.Releases[i].CmpVersion(version) {
				found = true
				e := result.Releases[i].Download(creds)

				if e != nil {
					// don't return on any error
					// continue processing rest of the downloads
					fmt.Println(e)
				}

				// found the latest release so bust out of the inner loop
				break
			}
		}

		if !found {
			fmt.Printf("Could not find version %s for mod %s\n", version, result.Name)
		}
	}
}
