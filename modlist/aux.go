package modlist

import (
	"fmt"
	"regexp"
)

// add mods to the list and enable them
func Add(dir string, names ...string) error {
	list, e := Read(dir)

	if e != nil {
		return e
	}

	for _, name := range names {
		mod := &Mod{name, true, nil}
		list.Mods = append(list.Mods, mod)
	}

	return list.Write(dir)
}

// set the "enabled" status of mods
func SetStatus(dir string, enabled bool, names []string) error {
	list, e := Read(dir)

	if e != nil {
		return e
	}

	for _, name := range names {
		re, e := regexp.Compile(name)

		if e != nil {
			return e
		}

		found := false

		for _, mod := range list.Mods {
			if re.MatchString(mod.Name) {
				mod.Enabled = enabled
				found = true
			}
		}

		if !found {
			return fmt.Errorf("%s not found in the mod list", name)
		}
	}

	return list.Write(dir)
}
