package main

import (
	"fmt"
	"modtorio/api"
	"regexp"
)

const (
	REGEXP_FLAGS    = "(?i)"
	MIN_OPTION_ARGS = 2
)

// search for one or more mods.
// the search term is compiled as a regular expression
func search(flags *ModtorioFlags, options []string) error {
	var e error
	var re *regexp.Regexp
	var cmp func(*regexp.Regexp, *api.Result) bool
	count := len(options)

	switch options[0] {
	// determine search type
	case "--tag":
		cmp = matchTag
		fallthrough
	case "--owner":
		if count < MIN_OPTION_ARGS {
			return fmt.Errorf("Not enough arguments to search %s", options[0])
		}

		if cmp == nil {
			// cmp is still nil so tag must be "--owner"
			// (no fall through from "--tag")
			cmp = matchOwner
		}

		re, e = regexp.Compile(REGEXP_FLAGS + options[1])
	default:
		cmp = matchName
		re, e = regexp.Compile(REGEXP_FLAGS + options[0])
	}

	if e != nil {
		return e
	}

	results, e := api.GetAll()

	if e != nil {
		return e
	}

	matches := 0

	for _, result := range results {
		if cmp(re, result) {
			// increment the match count and print a divider and mod information
			matches++
			fmt.Println("-------------")
			fmt.Println(result)
		}
	}

	// print the ending divider and the match count
	fmt.Println("-------------")
	fmt.Printf("Found %d mods matching %v\n", matches, re)

	return nil
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
