package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

const (
	FILE_NAME = "mod-list.json"
)

type ModList struct {
	Mods []Mod `json:"mods"`
}

type Mod struct {
	Name    string `json:"name"`
	Enabled bool   `json:"enabled"`
}

func readModList(dir string) (*ModList, error) {
	if dir[len(dir)-1] != '/' {
		// append a trailing slash
		dir += "/"
	}

	path := dir + FILE_NAME
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

	return list, nil
}
