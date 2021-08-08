package main

import (
	"fmt"
)

func help(flags *ModtorioFlags, options []string) error {
	if len(options) > 0 {
		switch options[0] {
		case CMD_SEARCH:
			helpSearch()
		case CMD_DOWNLOAD:
			helpDownload()
		case CMD_UPDATE:
			helpUpdate()
		case CMD_ENABLE:
			helpEnable()
		case CMD_DISABLE:
			helpDisable()
		case CMD_LIST:
			helpList()
		case CMD_HELP:
			helpHelp()
		default:
			return fmt.Errorf("Unknown command: %s", options[0])
		}

		return nil
	}

	// default to all if no command provided
	helpAll()

	return nil
}

func helpAll() {
	fmt.Printf("usage: modtorio [...flags] <command> [...options] <arguments>\n\n")
	fmt.Printf("Flags:\n")
	fmt.Printf("\t--dir\tSpecify the working directory for commands that interact with modlist.json. Leave blank if the current directory contains modlist.json or you want modlist.json to be created in the current directory.\n")
	fmt.Printf("\t--factorio\tSpecify the factorio version to compare releases against. Defaults to the latest version.\n\n")
	fmt.Printf("Commands:\n")
	helpHelp()
	helpSearch()
	helpDownload()
	helpUpdate()
	helpEnable()
	helpDisable()
	helpList()
}

func helpSearch() {
	// search command
	fmt.Printf("search\n")
	fmt.Printf("\tSearch for a mod. The argument is compiled as a regular expression.\n")
	fmt.Printf("\tOptions:\n")
	fmt.Printf("\t\t--tag\t\tSearch for mods based on a tag\n")
	fmt.Printf("\t\t--owner\t\tSearch for mods created by a user\n")
	fmt.Printf("\t\t--name-only\tOnly print out the mod name for matching mods\n")
	fmt.Printf("\tExamples:\n")
	fmt.Printf("\t\tmodtorio search ^bob\n")
	fmt.Printf("\t\tmodtorio search --tag general\n")
	fmt.Printf("\t\tmodtorio search --owner py.*\n")
	fmt.Printf("\t\tmodtorio search --name-only angel\n")
	fmt.Printf("\t\tmodtorio search --owner bobingabout --tag general --name-only\n")
}

func helpDownload() {
	// download command
	fmt.Printf("download\n")
	fmt.Printf("\tDownload any number of mods. Must be listed by the mod name.\n")
	fmt.Printf("\tExamples:\n")
	fmt.Printf("\t\tmodtorio download bobinserters miniloader pyhightech\n")
	fmt.Printf("\t\tmodtorio --factorio 0.17 --dir ~/.config/factorio/mods download bobinserters helicopters\n")
}

func helpUpdate() {
	// update command
	fmt.Printf("update\n")
	fmt.Printf("\tUpdate all mods to their latest release for the factorio version (if specified).\n")
	fmt.Printf("\tExamples:\n")
	fmt.Printf("\t\tmodtorio update\n")
	fmt.Printf("\t\tmodtorio --factorio 0.18 update\n")
	fmt.Printf("\t\tmodtorio --factorio 0.18 --dir ~/.config/factorio/mods update\n")
}

func helpEnable() {
	// enable command
	fmt.Printf("enable\n")
	fmt.Printf("\tEnable mods. Arguments are compiled as regular expressions.\n")
	fmt.Printf("\tExamples:\n")
	fmt.Printf("\t\tmodtorio enable bob.* pyhightech ^angel\n")
	fmt.Printf("\t\tmodtorio --dir ~/.config/factorio/mods enable bob.*\n")
}

func helpDisable() {
	// disable command
	fmt.Printf("disable\n")
	fmt.Printf("\tDisable mods. Arguments are compiled as regular expressions.\n")
	fmt.Printf("\tExamples:\n")
	fmt.Printf("\t\tmodtorio disable bob.* pyhightech ^angel\n")
	fmt.Printf("\t\tmodtorio --dir ~/.config/factorio/mods disable bob.*\n")
}

func helpList() {
	// list command
	fmt.Printf("list\n")
	fmt.Printf("\tList mods. Base mod is intentionally left out as it should not be manipulated.\n")
	fmt.Printf("\tOptions:\n")
	fmt.Printf("\t\t--all\t\tList all installed mods (default)\n")
	fmt.Printf("\t\t--enabled\tList all enabled mods\n")
	fmt.Printf("\t\t--disabled\tList all disabled mods\n")
	fmt.Printf("\tExamples:\n")
	fmt.Printf("\t\tmodtorio list\n")
	fmt.Printf("\t\tmodtorio list --all\n")
	fmt.Printf("\t\tmodtorio list --enabled\n")
	fmt.Printf("\t\tmodtorio list --disabled\n")
	fmt.Printf("\t\tmodtorio --dir ~/.config/factorio/mods list\n")
}

func helpHelp() {
	// help command
	fmt.Printf("help\n")
	fmt.Printf("\tShow help text for a command or print this help text.\n")
	fmt.Printf("\tExamples:\n")
	fmt.Printf("\t\tmodtorio help\n")
	fmt.Printf("\t\tmodtorio help search\n")
}
