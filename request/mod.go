package request

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"modtorio/credentials"
	"net/http"
)

const (
	URL_DL  = "https://mods.factorio.com"
	URL_MOD = "https://mods.factorio.com/api/mods/"
)

type Mod struct {
	Downloads_count       uint
	Category, Name, Owner string
	Releases              []Release
}

type Release struct {
	Download_url, File_name, Released_at, Version, Sha1 string
	Info_json                                           Version
}

type Version struct {
	Factorio_version string
}

// get the "short" mod information
func GetMod(name string) (*Mod, error) {
	res, err := http.Get(URL_MOD + name)

	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(res.Body)

	if err != nil {
		return nil, err
	}

	if res.StatusCode >= 400 {
		// something went wrong
		ae := &apiError{}
		err = json.Unmarshal(body, ae)

		if err != nil {
			return nil, err
		}

		return nil, fmt.Errorf("%s: %s", res.Status, ae.message)
	}

	mod := &Mod{}
	err = json.Unmarshal(body, mod)

	if err != nil {
		return nil, err
	}

	return mod, nil
}

// download a specific mod (and version)
func DownloadMod(release *Release, creds *credentials.Credentials) error {
	// the url is of the form:
	// https://mods.factorio.com/download/{download_id}?username=<usr>&token=<tkn>
	mod := URL_DL + release.Download_url + "?username=" + creds.Username + "&token=" + creds.Token
	res, err := http.Get(mod)

	if err != nil {
		return err
	}

	if res.StatusCode >= 400 {
		// something went wrong with the request
		return fmt.Errorf("%s", res.Status)
	}

	// read the file into a []byte
	body, err := ioutil.ReadAll(res.Body)
	res.Body.Close()

	if err != nil {
		return err
	}

	// write the []byte to the specified file (release.File_name)
	err = ioutil.WriteFile(release.File_name, body, 0644)

	if err != nil {
		return err
	}

	// TODO: verify SHA1 hash? show progress?

	fmt.Printf("%s downloaded\n", release.File_name)

	return nil
}
