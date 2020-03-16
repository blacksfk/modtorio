/*
Sub-package containing all operations related to interactions
with the mod list file.
*/
package modlist

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"modtorio/common"
	"os"
	"regexp"
	"strings"
)

const (
	MODE       = 0644
	VERSION_RE = `(\d(?:\.\d+)+)`
	FILE_NAME  = "mod-list.json"
)

type ModList struct {
	Mods []*Mod `json:"mods"`
}

// write the mod list in the specified directory
func (list *ModList) Write(dir string) error {
	path := genPath(dir)
	bytes, e := json.Marshal(list)

	if e != nil {
		return e
	}

	// base mod should not be in the list, so append it
	list.Mods = append(list.Mods, base)

	return ioutil.WriteFile(path, bytes, MODE)
}

// populate the archive file data for all mods in this list.
// does not return an error if no match for a mod is found.
func (list *ModList) FindArchives(dir string) error {
	files, e := ioutil.ReadDir(dir)

	if e != nil {
		return e
	}

	for _, mod := range list.Mods {
		b := strings.Builder{}
		b.WriteString(mod.Name)   // match the mod name exactly
		b.WriteString("_")        // underscore between name and version
		b.WriteString(VERSION_RE) // match any version

		re, e := regexp.Compile(b.String())

		if e != nil {
			// regexp compilation failure
			return e
		}

		length := len(files)

		for i := 0; i < length; i++ {
			matches := re.FindStringSubmatch(files[i].Name())

			if matches != nil {
				// match found:
				// [0]: full match (<mod_name>_<mod_version>)
				// [1]: version sub-group (<mod_version> eg. 0.17.3333)
				mod.Archive, e = NewArchive(matches[0], matches[1])

				if e != nil {
					// something went wrong with semantic version extraction,
					// no reason to stop processing
					fmt.Println(e)
				}

				// remove the file from the list (no longer needed for comparison)
				if i == 0 {
					// file is at the beginning
					files = files[1:]
				} else if i == length-1 {
					// file is at the end
					files = files[:i]
				} else {
					// file is in the middle somewhere
					files = append(files[:i], files[i+1:]...)
				}

				break
			}
		}
	}

	return nil
}

// get an array of all mods' names in this list
func (list *ModList) GetAllModNames() []string {
	var names []string

	for _, mod := range list.Mods {
		names = append(names, mod.Name)
	}

	return names
}

type Mod struct {
	Name    string `json:"name"`
	Enabled bool   `json:"enabled"`
	// the archive data should not be written to mod-list.json,
	// so keep it hidden with tag: "-"
	Archive *Archive `json:"-"`
}

type Archive struct {
	Name, Version string
	Semver        *common.Semver
}

// Extract the semantic version from `version` and create
// a new archive. Returns an error if semantic version extraction
// failed.
func NewArchive(name, version string) (*Archive, error) {
	semver, e := common.NewSemver(version)

	if e != nil {
		return nil, e
	}

	return &Archive{name, version, semver}, nil
}

// base mod should always be present in the file,
// but does not have an archive. so it is removed before
// a read is returned, and is added during a write
var base *Mod = &Mod{"base", true, nil}

func Read(dir string) (*ModList, error) {
	path := genPath(dir)
	bytes, e := ioutil.ReadFile(path)

	if e != nil {
		if os.IsNotExist(e) {
			// return an empty mod list
			return &ModList{}, nil
		} else {
			// permission error, or something else
			return nil, e
		}
	}

	// file exists and is readable
	list := &ModList{}
	e = json.Unmarshal(bytes, list)

	if e != nil {
		return nil, e
	}

	length := len(list.Mods)

	// loop through the mod list and remove the base mod
	// so as not to interfere with downloading, updating etc
	for i := 0; i < length; i++ {
		if list.Mods[i].Name == base.Name {
			// base mod found
			if i == 0 {
				// base mod is at the beginning
				list.Mods = list.Mods[1:]
			} else if i == length-1 {
				// base mod is at the end
				list.Mods = list.Mods[:i]
			} else {
				// base mod is in the middle
				list.Mods = append(list.Mods[:i], list.Mods[i+1:]...)
			}

			break
		}
	}

	// base mod removed (or new, empty list)
	return list, nil
}

func genPath(dir string) string {
	if dir[len(dir)-1] != '/' {
		// append a slash
		dir += "/"
	}

	return dir + FILE_NAME
}
