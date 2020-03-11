package main

import (
	"fmt"
	"modtorio/modlist"
)

const (
	H_SEP = "-"
)

func list(options []string) error {
	if len(options) > 0 {
		switch options[0] {
		case "--all":
			return listAll()
		case "--enabled":
			return listMods(true)
		case "--disabled":
			return listMods(false)
		default:
			return fmt.Errorf("Unknown option %s for command list", options[0])
		}
	}

	// if no options default to all()
	return listAll()
}

func listMods(enabled bool) error {
	list, e := modlist.Read(FLAGS.dir)

	if e != nil {
		return e
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

	return nil
}

// display all mods by name (column 1) and their status (column 2)
func listAll() error {
	list, e := modlist.Read(FLAGS.dir)

	if e != nil {
		return e
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

	return nil
}

func printStringTimes(s string, times int) {
	for i := 0; i < times; i++ {
		fmt.Print(s)
	}
}
