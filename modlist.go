package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

const (
	MODE = 0644
	FILE_NAME = "mod-list.json"
)

type ModList struct {
	Mods []Mod `json:"mods"`
}

type Mod struct {
	Name    string `json:"name"`
	Enabled bool   `json:"enabled"`
}

// base mod should always be present in the file,
// but does not have an archive. so it is removed before
// a readModList is returned, and is added during a writeModList
var base Mod = Mod{"base", true}

func readModList(dir string) (*ModList, error) {
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
			} else if i == length - 1 {
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

func writeModList(list *ModList, dir string) error {
	path := genPath(dir)
	bytes, e := json.Marshal(list)

	if e != nil {
		return e
	}

	// base mod should not be present, so append it
	list.Mods = append(list.Mods, base)

	return ioutil.WriteFile(path, bytes, MODE)
}

func genPath(dir string) string {
	if dir[len(dir)-1] != '/' {
		// append a slash
		dir += "/"
	}

	return dir + FILE_NAME
}
