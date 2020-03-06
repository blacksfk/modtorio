package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"modtorio/credentials"
	"net/http"
	"net/url"
)

const (
	URL_LOGIN = "https://auth.factorio.com/api-login"
)

func Login(creds *credentials.Credentials) error {
	data := url.Values{}

	// append the username and password
	data.Set("username", creds.Username)
	data.Set("password", creds.Password)

	// send the request
	res, e := http.PostForm(URL_LOGIN, data)

	if e != nil {
		return e
	}

	body, e := handleResponse(res)

	if e != nil {
		return e
	}

	// request/response was all good, convert to JSON
	var loginData []string
	e = json.Unmarshal(body, &loginData)

	if e != nil {
		return e
	}

	creds.Token = loginData[0]

	// all is good, return no error
	return nil
}