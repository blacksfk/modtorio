package main

import (
	"bufio"
	"fmt"
	"golang.org/x/crypto/ssh/terminal"
	"modtorio/api"
	"modtorio/credentials"
	"modtorio/modlist"
	"os"
	"strings"
	"syscall"
)

const (
	MAX_LOGIN_ATTEMPTS = 5
)

func download(flags *ModtorioFlags, options []string) error {
	// get the mod results for each mod
	results, e := api.GetAll(options...)

	if e != nil {
		return e
	}

	var downloads []*api.Release
	var toBeEnabled []string

	for _, result := range results {
		found := false

		for i := len(result.Releases) - 1; i >= 0; i-- {
			if result.Releases[i].CmpFactorioVersion(flags.factorio) == 0 {
				found = true
				downloads = append(downloads, result.Releases[i])
				toBeEnabled = append(toBeEnabled, result.Name)
				break
			}
		}

		if !found {
			fmt.Printf("No matching factorio version (%v) found for mod %s\n", flags.factorio, result.Name)
		}
	}

	e = downloadReleases(flags.dir, downloads)

	if e != nil {
		return e
	}

	// enable (or add) all downloaded releases
	return modlist.Add(flags.dir, toBeEnabled...)
}

func attemptLogin() (*credentials.Credentials, error) {
	creds, e := credentials.FromCache()

	if e == nil {
		// credentials obtained from cache
		return creds, nil
	}

	// something went wrong with the cached credentials
	for attempts := 0; attempts < MAX_LOGIN_ATTEMPTS; attempts++ {
		creds, e = promptForCreds()

		if e != nil {
			return nil, e
		}

		fmt.Print("Retrieving token...")
		creds.Token, e = api.Login(creds.Username, creds.Password)

		if e != nil {
			fmt.Println(e)
		} else {
			// logged in successfully, cache creds
			fmt.Println("success")
			creds.ToCache()

			return creds, nil
		}
	}

	return nil, fmt.Errorf("Maximum login attempts reached")
}

// prompt the user for their login credentials
func promptForCreds() (*credentials.Credentials, error) {
	fmt.Println("Please enter your credentials to download from mods.factorio.com")

	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("Username: ")
	scanner.Scan()

	if e := scanner.Err(); e != nil {
		return nil, e
	}

	username := scanner.Text()

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

	return credentials.NewCredentials(username, string(bytes)), nil
}

// Download the releases. Authenticates the user prior to downloading.
func downloadReleases(dir string, releases []*api.Release) error {
	count := len(releases)

	if count == 0 {
		return fmt.Errorf("Nothing to download")
	}

	// print a summary of the releases to be downloaded
	fmt.Printf("\nDownloads (%d):", count)

	for i := 0; i < count; i++ {
		fmt.Printf(" %s", releases[i].File_name)
	}

	// prompt the user for confirmation of the releases to be downloaded
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Printf("\n\nContinue? (Y/n): ")
	scanner.Scan()

	if e := scanner.Err(); e != nil {
		return e
	}

	answer := scanner.Text()

	// cancel downloads if not "yes" or empty string (linefeed)
	if len(answer) > 0 && strings.ToLower(answer)[0] != 'y' {
		return fmt.Errorf("Downloads cancelled")
	}

	// log the user in
	creds, e := attemptLogin()

	if e != nil {
		return e
	}

	fmt.Println()

	// download all of the releases
	for i := 0; i < count; i++ {
		fmt.Printf("Downloading %s...", releases[i].File_name)
		e = releases[i].Download(dir, creds)
		fmt.Println("done")

		if e != nil {
			return e
		}
	}

	return nil
}
