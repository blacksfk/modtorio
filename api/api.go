package api

import (
	"encoding/json"
	"fmt"
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

	if res.StatusCode >= http.StatusBadRequest {
		reqError := &apiError{}
		e = json.Unmarshal(body, reqError)

		if e != nil {
			return nil, e
		}

		// something went wrong with the request
		return nil, fmt.Errorf("%d %s: %s", res.StatusCode, res.Status, reqError.message)
	}

	return body, nil
}
