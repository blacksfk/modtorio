package main

import (
	"fmt"
)

func listEnabled(args []string) {
	printMods(true)
}

func listDisabled(args []string) {
	printMods(false)
}

func printMods(enabled bool) {
	list, e := readModList(dir)

	if e != nil {
		fmt.Println(e)

		return
	}

	for _, mod := range list.Mods {
		// print if:
		// the mod is enabled and list enabled mods (enabled = true)
		// or:
		// the mod is disabled and list disabled mods (enabled = false)
		if enabled && mod.Enabled || !enabled && !mod.Enabled {
			fmt.Println(mod.Name)
		}
	}
}
