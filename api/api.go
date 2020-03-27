package api

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

// auxiliary function to check for request errors
func handleResponse(res *http.Response) ([]byte, error) {
	defer res.Body.Close()
	body, e := ioutil.ReadAll(res.Body)

	if e != nil {
		return nil, e
	}

	if res.StatusCode >= http.StatusBadRequest && res.StatusCode < http.StatusInternalServerError {
		// only unmarshal the body if a 4xx error occurred
		reqError := &apiError{}
		e = json.Unmarshal(body, reqError)

		if e != nil {
			return nil, e
		}

		reqError.StatusText = res.Status

		return nil, reqError
	} else if res.StatusCode >= http.StatusInternalServerError {
		// the API crashed or some other bs
		return nil, &apiError{res.StatusCode, res.Status, string(body)}
	}

	return body, nil
}
