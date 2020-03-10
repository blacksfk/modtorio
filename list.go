package main

import (
	"fmt"
)

const (
	H_SEP = "-"
)

func list(args []string) {
	if len(args) > 0 {
		switch args[0] {
		case "--all":
			listAll()
		case "--enabled":
			listMods(true)
		case "--disabled":
			listMods(false)
		default:
			fmt.Printf("Unknown option %s for command list\n", args[0])
			help()
		}

		return
	}

	// if no args default to all()
	listAll()
}

func listMods(enabled bool) {
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

func listAll() {
	list, e := readModList(dir)

	if e != nil {
		fmt.Println(e)

		return
	}

	// default to 4 for "Name" header
	longest := 4
	modCount := len(list.Mods)

	for i := 0; i < modCount; i++ {
		if l := len(list.Mods[i].Name); l > longest {
			longest = l
		}
	}

	// how many dashes to print
	// +1: 3 pipes - 2 surrounding
	// +4: 4 spaces between pipes and strings
	// +7: len("Enabled")
	hSepCount := longest + 1 + 4 + 7

	fmt.Print("|")
	printStringTimes(H_SEP, hSepCount)
	fmt.Print("|")
	fmt.Println()
	fmt.Printf("| % -*s | Enabled |\n", longest, "Name")
	fmt.Print("|")
	printStringTimes(H_SEP, hSepCount)
	fmt.Print("|")
	fmt.Println()

	for _, mod := range list.Mods {
		fmt.Printf("| % -*s | %-7t |\n", longest, mod.Name, mod.Enabled)
	}

	fmt.Print("|")
	printStringTimes(H_SEP, hSepCount)
	fmt.Print("|")
	fmt.Println()
}

func printStringTimes(s string, times int) {
	for i := 0; i < times; i++ {
		fmt.Print(s)
	}
}
