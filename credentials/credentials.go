/*
Package to hold the credentials for a user logged in via the
auth.factorio.com web API.
*/
package credentials

import (
	"encoding/json"
	"os"
)

const (
	MODE  = 0600
	CACHE = "./modtorio_user.json"
)

type Credentials struct {
	Username, Token string
	Password        string `json:"-"` // don't store the user's password
}

// Create a new set of credentials (minus the token).
func NewCredentials(username, password string) *Credentials {
	creds := Credentials{}

	creds.Username = username
	creds.Password = password

	return &creds
}

// Get a set of credentials from the cache.
func FromCache() (*Credentials, error) {
	bytes, e := os.ReadFile(CACHE)

	if e != nil {
		return nil, e
	}

	creds := &Credentials{}
	e = json.Unmarshal(bytes, creds)

	if e != nil {
		return nil, e
	}

	return creds, nil
}

// Write a set of credentials to the cache (minus the password).
func (c *Credentials) ToCache() error {
	bytes, e := json.Marshal(c)

	if e != nil {
		return e
	}

	return os.WriteFile(CACHE, bytes, MODE)
}
