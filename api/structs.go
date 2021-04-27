package api

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/blacksfk/modtorio/common"
	"github.com/blacksfk/modtorio/credentials"
)

const (
	MODE   = 0644
	URL_DL = "https://mods.factorio.com"
)

// returned from mods.factorio.com/api/mods
type ModListResponse struct {
	Pagination Pagination
	Results    []*Result
}

// pagination data (sub-struct) of a ModListResponse
type Pagination struct {
	Count, Page           int
	Page_count, Page_size int
	Links                 struct {
		First, Prev string
		Next, Last  string
	}
}

// mod data
type Result struct {
	Downloads_count              uint
	Name, Owner, Summary         string
	Title, Changelog, Created_at string
	Description, Github_path     string
	Category, Homepage           string
	Latest_release               *Release
	Releases                     []*Release
	Tag                          []*Tag
}

// pretty print a mod's information
func (r *Result) String() string {
	return fmt.Sprintf("%s\n\tName: %s\n\tOwner: %s\n\tCategory: %s\n\tSummary: %s\n", r.Title, r.Name, r.Owner, r.Category, r.Summary)
}

// specific release information of a mod
type Release struct {
	Download_url, File_name    string
	Released_at, Version, sha1 string
	Semver                     *common.Semver
	Info_json                  struct {
		Factorio_version string
		Semver           *common.Semver
	}
}

// compare release version
func (r *Release) CmpVersion(semver *common.Semver) int {
	return r.Semver.Cmp(semver)
}

// compare release factorio version
func (r *Release) CmpFactorioVersion(semver *common.Semver) int {
	return r.Info_json.Semver.Cmp(semver)
}

// download a release
func (r *Release) Download(dir string, creds *credentials.Credentials) error {
	b := strings.Builder{}
	b.WriteString(URL_DL)
	b.WriteString(r.Download_url)
	b.WriteString("?username=")
	b.WriteString(creds.Username)
	b.WriteString("&token=")
	b.WriteString(creds.Token)

	res, e := http.Get(b.String())

	if e != nil {
		return e
	}

	body, e := handleResponse(res)

	if e != nil {
		return e
	}

	if dir[len(dir)-1] != '/' {
		// append a slash
		dir += "/"
	}

	path := dir + r.File_name
	e = ioutil.WriteFile(path, body, MODE)

	if e != nil {
		return e
	}

	return nil
}

// parse the version and factorio_version as semantic versions
func (r *Release) ParseVersions() error {
	var e error

	r.Semver, e = common.NewSemver(r.Version)

	if e != nil {
		return e
	}

	r.Info_json.Semver, e = common.NewSemver(r.Info_json.Factorio_version)

	return e
}

// mod tags (refer to array above)
type Tag struct {
	Id                             int
	Name, Title, Description, Type string
}

// JSON API errors
type apiError struct {
	Status              int
	StatusText, Message string
}

func (ae apiError) Error() string {
	return fmt.Sprintf("%s: %s", ae.StatusText, ae.Message)
}
