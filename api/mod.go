package api

import (
	"encoding/json"
	"net/http"
	"strings"
)

const (
	PAGE_SIZE = "?page_size=max"
	URL_MOD   = "https://mods.factorio.com/api/mods"
)

// get all mods or mods exactly matching the names provided
func GetAll(names ...string) ([]*Result, error) {
	count := len(names)
	url := strings.Builder{}
	namelist := strings.Builder{}

	// build a string of mod names separated by a comma
	for i := 0; i < count; i++ {
		namelist.WriteString(names[i])

		if i < count-1 {
			// only append a comma if not the last element
			namelist.WriteString(",")
		}
	}

	// build the initial URL string
	url.WriteString(URL_MOD)
	url.WriteString(PAGE_SIZE)

	// only append &namelist=<list> if mod names were provided
	if count > 0 {
		url.WriteString("&namelist=")
		url.WriteString(namelist.String())
	}

	// get all mods in one shot by requesting a "page" with all of the mods
	// i.e. page_size=max
	res, e := http.Get(url.String())

	if e != nil {
		return nil, e
	}

	body, e := handleResponse(res)

	if e != nil {
		return nil, e
	}

	mlr := &ModListResponse{}
	e = json.Unmarshal(body, mlr)

	if e != nil {
		return nil, e
	}

	// parse the factorio and release versions as semvers
	for _, result := range mlr.Results {
		for _, release := range result.Releases {
			e := release.ParseVersions()

			if e != nil {
				return nil, e
			}
		}
	}

	return mlr.Results, nil
}
