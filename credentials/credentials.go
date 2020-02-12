/*
Package to hold the credentials for a user logged in via the
auth.factorio.com web API.
*/
package credentials

type Credentials struct {
	Username, Password, Token string
}

func NewCredentials(username, password string) *Credentials {
	creds := Credentials{}

	creds.Username = username
	creds.Password = password

	return &creds
}
