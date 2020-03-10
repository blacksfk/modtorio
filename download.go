package main

import (
	"bufio"
	"fmt"
	"golang.org/x/crypto/ssh/terminal"
	"modtorio/api"
	"modtorio/credentials"
	"os"
	"regexp"
	"strings"
	"syscall"
)

const (
	MAX_LOGIN_ATTEMPS = 5
)

func download(options []string) error {
	var creds *credentials.Credentials
	var e error
	attempts := 0

	// loop until the login is successful or until MAX_LOGIN_ATTEMPTS is reached
	for ; attempts < MAX_LOGIN_ATTEMPS; attempts++ {
		// prompt the user for their credentials
		creds, e = promptForCreds()

		if e != nil {
			// something went wrong with the input
			return e
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
		return fmt.Errorf("Maximum login attempts exceeded")
	}

	// get the mod results for each mod
	results, e := api.GetAll(options...)

	if e != nil {
		return e
	}

	if fVer == DEF_VER {
		// latest
		downloadLatest(results, creds)

		return nil
	}

	// a version was specified
	return downloadVersion(results, creds)
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
func downloadVersion(results []*api.Result, creds *credentials.Credentials) error {
	re := regexp.MustCompile(fVer)

	for _, result := range results {
		found := false

		// start from the end to find the latest release for
		// the specified factorio version
		for i := len(result.Releases) - 1; i >= 0; i-- {
			if result.Releases[i].CmpVersion(re) {
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
			return fmt.Errorf("Could not find version %s for mod %s", re, result.Name)
		}
	}

	return nil
}
