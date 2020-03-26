package main

import (
	"flag"
	"fmt"
	"modtorio/api"
	"regexp"
)

const (
	REGEXP_FLAGS    = "(?i)"
	MIN_OPTION_ARGS = 2
	S_FLAG_NAME     = "name-only"
	S_FLAG_OWNER    = "owner"
	S_FLAG_TAG      = "tag"
	FULL_RESULT_SEP = "--------------------"
)

type cmpFunc func(*SearchRE, *api.Result) bool

type SearchRE struct {
	*regexp.Regexp
	cmp cmpFunc
}

func newSearchRE(s string, cmp cmpFunc) (*SearchRE, error) {
	re, e := regexp.Compile(REGEXP_FLAGS + s)

	if e != nil {
		return nil, e
	}

	return &SearchRE{re, cmp}, nil
}

func (sre *SearchRE) match(result *api.Result) bool {
	return sre.cmp(sre, result)
}

// search for one or more mods.
// the search term is compiled as a regular expression
func search(flags *ModtorioFlags, options []string) error {
	var nameOnly bool
	var sres []*SearchRE
	var strOwner, strTag string

	searchFlags := flag.NewFlagSet("Search flags", flag.ContinueOnError)

	searchFlags.BoolVar(&nameOnly, S_FLAG_NAME, false, "Print a space-delimited list of mod names")
	searchFlags.StringVar(&strOwner, S_FLAG_OWNER, "", "Match a mod by owner")
	searchFlags.StringVar(&strTag, S_FLAG_TAG, "", "Match a mod by tag")
	searchFlags.Parse(options)

	sres = compileAndAppend(strOwner, matchOwner, sres)
	sres = compileAndAppend(strTag, matchTag, sres)
	sres = compileAndAppend(searchFlags.Arg(0), matchNameOrTitle, sres)

	if len(sres) == 0 {
		// nothing compiled properly
		return fmt.Errorf("search: regular expressions failed compilation")
	}

	results, e := api.GetAll()

	if e != nil {
		return e
	}

	var matches []*api.Result
	matchCount := 0

	// for each mod result:
	// if it returns true for all comparison functions add it to the matches
	for _, result := range results {
		match := true

		for _, sre := range sres {
			match = sre.match(result)

			if !match {
				// did not match so break early
				break
			}
		}

		if match {
			// matched all comparison functions
			matchCount++
			matches = append(matches, result)
		}
	}

	if nameOnly {
		// print all matched mods' names on a single line
		for i := 0; i < matchCount; i++ {
			fmt.Print(matches[i].Name)

			if i < matchCount-1 {
				// only print spaces if not the last element
				fmt.Print(" ")
			} else {
				fmt.Println()
			}
		}
	} else {
		// print the mod formatted, with a separator, and total match count
		for i := 0; i < matchCount; i++ {
			fmt.Println(matches[i])
			fmt.Println(FULL_RESULT_SEP)
		}

		fmt.Printf("%d matches\n", matchCount)
	}

	return nil
}

func matchOwner(sre *SearchRE, result *api.Result) bool {
	return sre.MatchString(result.Owner)
}

func matchTag(sre *SearchRE, result *api.Result) bool {
	return sre.MatchString(result.Category)
}

func matchNameOrTitle(sre *SearchRE, result *api.Result) bool {
	return sre.MatchString(result.Name) || sre.MatchString(result.Title)
}

func compileAndAppend(s string, cmp cmpFunc, sres []*SearchRE) []*SearchRE {
	if len(s) > 0 {
		sre, e := newSearchRE(s, cmp)

		if e == nil {
			// only append if no error was returned
			sres = append(sres, sre)
		}
	}

	return sres
}
