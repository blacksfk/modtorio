package main

import (
	"modtorio/api"
	"modtorio/common"
	"modtorio/modlist"
)

type ModResult struct {
	*modlist.Mod
	*api.Result
}

// returns a release that matches the factorio version and only if
// there is no archive or a release is newer than what we have
func (mr *ModResult) FindRelease(factorio *common.Semver) *api.Release {
	for i := len(mr.Releases) - 1; i >= 0; i-- {
		release := mr.Releases[i]

		if release.CmpFactorioVersion(factorio) == 0 && (mr.Archive == nil || release.CmpVersion(mr.Archive.Semver) == 1) {
			return release
		}
	}

	// no match found
	return nil
}

func update(flags *ModtorioFlags, options []string) error {
	// first, get a list of mods
	list, e := modlist.Read(flags.dir)

	if e != nil {
		return e
	}

	// then scour the directory for the files
	e = list.FindArchives(flags.dir)

	if e != nil {
		return e
	}

	// get all of the same mods from the API
	results, e := api.GetAll(list.GetAllModNames()...)

	if e != nil {
		return e
	}

	// the following functionality has been broken up into smaller,
	// more easily maintainable, loops. Otherwise, the condensed functionality
	// is an n^3 pyramid of doom
	//
	// Premise:
	// for all mods in the mod list have an archive:
	// if the api contains a newer version and it matches the factorio version, download it
	//
	// for all mods in the mod list that don't have an archive:
	// check that the factorio version matches the flag

	// first, create an array of hybrid mod results, where each element is matched
	// to both mod (from the list) and result (from the api)
	var modResults []*ModResult

	for _, mod := range list.Mods {
		length := len(results)

		for i := 0; i < length; i++ {
			if mod.Name == results[i].Name {
				modResults = append(modResults, &ModResult{mod, results[i]})

				if length > 1 {
					// remove the result from the slice
					if i == 0 {
						// result is at the beginning
						results = results[i+1:]
					} else if i == length-1 {
						// result is at the end
						results = results[:i]
					} else {
						// result is in the middle
						results = append(results[:i], results[i+1:]...)
					}

					// break the inner loop
					break
				}
			}
		}
	}

	// second, loop through all of the mod results and generate an array of
	// releases to download
	var downloads []*api.Release

	for _, mr := range modResults {
		release := mr.FindRelease(flags.factorio)

		if release != nil {
			downloads = append(downloads, release)
		}
	}

	// last, attempt to login and download the releases
	return downloadReleases(flags.dir, downloads)
}
