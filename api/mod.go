package api

import (
	"encoding/json"
	"fmt"
	"net/http"
)

const (
	URL_MOD = "https://mods.factorio.com/api/mods"
)

// get all mods
func GetAll() ([]*Result, error) {
	count, e := getModCount()

	if e != nil {
		return nil, e
	}

	// get all mods in one shot by requesting a "page" with all of the mods
	// i.e. page_size=<no_of_mods>
	res, e := http.Get(URL_MOD + fmt.Sprintf("?page_size=%d", count))

	if e != nil {
		return nil, e
	}

	body, e := handleResponse(res)

	if e != nil {
		return nil, e
	}

	modList := &ModListResponse{}
	e = json.Unmarshal(body, modList)

	if e != nil {
		return nil, e
	}

	// no need to return pagination data
	return modList.Results, nil
}

// get the count of mods
func getModCount() (int, error) {
	// get the count by sending a request with only one mod per page
	res, e := http.Get(URL_MOD + "?page_size=1")

	if e != nil {
		return 0, e
	}

	body, e := handleResponse(res)

	if e != nil {
		return 0, e
	}

	modList := &ModListResponse{}
	e = json.Unmarshal(body, modList)

	if e != nil {
		return 0, e
	}

	return modList.Pagination.Page_count, nil
}
