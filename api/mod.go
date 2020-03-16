package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
)

const (
	URL_MOD = "https://mods.factorio.com/api/mods"
)

// get all mods
func GetAll(names ...string) ([]*Result, error) {
	var e error
	var size int
	namelist := strings.Builder{}
	count := len(names)

	if count == 0 {
		// no mod names passed so get all mods
		size, e = getModCount()
	} else {
		size = count

		// build a string of mod names separated by a comma
		for i := 0; i < count; i++ {
			namelist.WriteString(names[i])

			if i < count-1 {
				// only append a comma if not the last element
				namelist.WriteString(",")
			}
		}
	}

	if e != nil {
		return nil, e
	}

	url := strings.Builder{}
	url.WriteString(URL_MOD)
	url.WriteString("?page_size=")
	url.WriteString(strconv.Itoa(size))

	// only append &namelist=<list> if count > 0
	if count > 0 {
		url.WriteString("&namelist=")
		url.WriteString(namelist.String())
	}

	// get all mods in one shot by requesting a "page" with all of the mods
	// i.e. page_size=<no_of_mods>
	res, e := http.Get(url.String())

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

	// parse the factorio version of each release as a semver
	for _, result := range modList.Results {
		for _, release := range result.Releases {
			e = release.ParseVersions()

			if e != nil {
				return nil, e
			}
		}
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
