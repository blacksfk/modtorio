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
	MAX_LOGIN_ATTEMPTS = 5
)

func download(options []string) error {
	// attempt login of the user
	creds, e := attemptLogin()

	// get the mod results for each mod
	results, e := api.GetAll(options...)

	if e != nil {
		return e
	}

	for _, result := range results {
		found := false

		for i := len(result.Releases) - 1; i >= 0; i-- {
			if result.Releases[i].CmpVersion(FLAGS.factorio) == 0 {
				found = true
				e = result.Releases[i].Download(FLAGS.dir, creds)

				if e != nil {
					// don't return on download error
					// continue processing rest of the downloads
					fmt.Printf("Error downloading %s: %v\n", result.Name, e)
				}
			}
		}

		if !found {
			fmt.Printf("No matching factorio version (%v) found for mod %s\n", FLAGS.factorio, result.Name)
		}
	}

	return nil
}

func attemptLogin() (*credentials.Credentials, error) {
	attempts := 0

	for ; attempts < MAX_LOGIN_ATTEMPTS; attempts++ {
		creds, e := promptForCreds()

		if e != nil {
			return nil, e
		}

		e = api.Login(creds)

		if e != nil {
			fmt.Println(e)
		} else {
			// logged in successfully
			return creds, nil
		}
	}

	return nil, fmt.Errorf("Maximum login attempts reached")

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

/*// download the latest release of a mod
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
	re := regexp.MustCompile(FLAGS.fVer)

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
}*/
