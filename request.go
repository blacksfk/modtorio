package main

import (
	"net/url"
	"net/http"
	"io/ioutil"
	"encoding/json"
)

const (
	loginUrl = "https://auth.factorio.com/api-login"
)

// factorio.com authentication data
type Auth struct {
	Token string
}

// log the user in with an HTTP POST request
func login(username, password string) (*Auth, error) {
	data := url.Values{username: {username}, password: {password}}
	res, err := http.PostForm(loginUrl, data)

	if err != nil {
		return nil, err
	}

	// read the response data
	body, err := ioutil.ReadAll(res.Body)
	res.Body.Close()

	if err != nil {
		return nil, err
	}

	// convert the response to JSON
	auth := &Auth{}
	err = json.Unmarshal(body, auth)

	if err != nil {
		return nil, err
	}

	// all is good, return no error
	return auth, nil
}
