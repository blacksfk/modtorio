package main

import (
	"fmt"
	"modtorio/api"
	"regexp"
)

const (
	REGEXP_FLAGS = "(?i)"
	MIN_OPTION_ARGS = 2
)

// search for one or more mods.
// the search term is compiled as a regular expression
func search(args []string) {
	var e error
	var re *regexp.Regexp
	var cmp func(*regexp.Regexp, *api.Result) bool
	count := len(args)

	switch args[0] {
	// determine search type
	case "--tag":
		cmp = matchTag
		fallthrough
	case "--owner":
		if count < MIN_OPTION_ARGS {
			fmt.Printf("Not enough arguments to search %s\n", args[0])
			help()

			return
		}

		if cmp == nil {
			// cmp is still nil so tag must be "--owner"
			// (no fall through from "--tag")
			cmp = matchOwner
		}

		re, e = regexp.Compile(REGEXP_FLAGS + args[1])
	default:
		cmp = matchName
		re, e = regexp.Compile(REGEXP_FLAGS + args[0])
	}

	if e != nil {
		fmt.Println(e)

		return
	}

	results, e := api.GetAll()

	if e != nil {
		fmt.Println(e)

		return
	}

	matches := 0

	for _, result := range results {
		if cmp(re, result) {
			// increment the match count and print a divider and mod information
			matches++
			fmt.Println("-------------")
			fmt.Println(result.String())
		}
	}

	// print the ending divider and the match count
	fmt.Println("-------------")
	fmt.Printf("Found %d mods matching %v\n", matches, re)
}

// compare the search term with the name and title properties
func matchName(re *regexp.Regexp, r *api.Result) bool {
	return re.MatchString(r.Name) || re.MatchString(r.Title)
}

// compare the search term with the category (tag) property
func matchTag(re *regexp.Regexp, r *api.Result) bool {
	return re.MatchString(r.Category)
}

// compare the search term with the owner property
func matchOwner(re *regexp.Regexp, r *api.Result) bool {
	return re.MatchString(r.Owner)
}
