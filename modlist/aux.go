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
		found := false

		for _, mod := range list.Mods {
			if mod.Name == name {
				// mod already exists, enable it
				found = true
				mod.Enabled = true
				break
			}
		}

		if !found {
			// mod does not exist, add and enable it
			newMod := &Mod{name, true, nil}
			list.Mods = append(list.Mods, newMod)
		}
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
