package main

import (
	"fmt"
)

const (
	H_SEP = "-"
)

func listInstalled(args []string) {
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

	fmt.Println("longest:", longest)

	// how many dashes to print
	// +1: 3 pipes - 2 surrounding
	// +4: 4 spaces between pipes and strings
	// +7: len("Enabled")
	hSepCount := longest + 1 + 4 + 7

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
