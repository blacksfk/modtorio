package main

import (
	"fmt"
	"strings"
	"net/url"
	"net/http"
	"io/ioutil"
	"encoding/json"
)

const (
	URL_LOGIN = "https://auth.factorio.com/api-login"
)

type ApiError struct {
	Message string
}

// log the user in with an HTTP POST request
func login(username, password string) (string, error) {
	data := url.Values{}

	data.Set("username", username)
	data.Set("password", password)

	// create the request and client
	req, err := http.NewRequest("POST", URL_LOGIN, strings.NewReader(data.Encode()))
	client := http.Client{}

	// add the content type header and send the request
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	res, err := client.Do(req)

	if err != nil {
		return "", err
	}

	// read the response data
	body, err := ioutil.ReadAll(res.Body)
	res.Body.Close()

	if err != nil {
		return "", err
	}

	if res.StatusCode >= 400 {
		// something went wrong with the request
		ae := &ApiError{}
		err = json.Unmarshal(body, ae)

		if err != nil {
			// conversion to JSON failed
			return "", err
		}

		// return the API error
		return "", fmt.Errorf("%s: %s", res.Status, ae.Message)
	}

	// request/response was all good, convert to JSON
	var resData []string
	err = json.Unmarshal(body, &resData)

	if err != nil {
		// conversion to JSON failed
		return "", err
	}

	// all is good, return no error
	return resData[0], nil
}
