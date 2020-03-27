package api

import (
	"encoding/json"
	"net/http"
	"net/url"
)

const (
	URL_LOGIN = "https://auth.factorio.com/api-login"
)

func Login(username, password string) (string, error) {
	data := url.Values{}

	// append the username and password
	data.Set("username", username)
	data.Set("password", password)

	// send the request
	res, e := http.PostForm(URL_LOGIN, data)

	if e != nil {
		return "", e
	}

	body, e := handleResponse(res)

	if e != nil {
		return "", e
	}

	// request/response was all good, convert to JSON
	var loginData []string
	e = json.Unmarshal(body, &loginData)

	if e != nil {
		return "", e
	}

	return loginData[0], nil
}
