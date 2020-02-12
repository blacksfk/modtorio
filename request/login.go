/*
Package to send login, search, and download requests to *.factorio.com
*/
package request

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

type apiError struct {
	message string
}

// log the user in with an HTTP POST request
func Login(creds *credentials.Credentials) error {
	data := url.Values{}

	// append username and password
	data.Set("username", creds.Username)
	data.Set("password", creds.Password)

	// send the request
	res, err := http.PostForm(URL_LOGIN, data)

	if err != nil {
		return err
	}

	// read the response data
	body, err := ioutil.ReadAll(res.Body)
	res.Body.Close()

	if err != nil {
		return err
	}

	if res.StatusCode >= 400 {
		// something went wrong with the request
		ae := &apiError{}
		err = json.Unmarshal(body, ae)

		if err != nil {
			// conversion to JSON failed
			return err
		}

		// return the API error
		return fmt.Errorf("%s: %s", res.Status, ae.message)
	}

	// request/response was all good, convert to JSON
	var resData []string
	err = json.Unmarshal(body, &resData)

	if err != nil {
		// conversion to JSON failed
		return err
	}

	creds.Token = resData[0]

	// all is good, return no error
	return nil
}
