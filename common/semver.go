/*
Common re-usable functions, types, and constants.
 */
package common

import (
	"fmt"
	"regexp"
	"strconv"
)

const (
	MIN_MATCHES = 4 // full, major, minor, patch
)

// extract a semantic version from a string
var re *regexp.Regexp = regexp.MustCompile(`\d+\.\d+\.\d+`)

type Semver struct {
	Major, Minor, Patch int
}

// Compare two Semantic Versions. Returns the result of cmp(a, b)
// for the first non-matching version (major, minor, patch) or 0
// if all versions match.
func (a *Semver) Cmp(b *Semver) int {
	if v := cmp(a.Major, b.Major); v != 0 {
		// major versions differ
		return v
	}

	if v := cmp(a.Minor, b.Minor); v != 0 {
		// minor versions differ
		return v
	}

	// major and minor match, compare patch versions
	return cmp(a.Patch, b.Patch)
}

// Compare a semver with a string equivalent. Converts the
// string to semver for comparison.
func (a *Semver) CmpString(version string) (int, error) {
	b, e := NewSemver(version)

	if e != nil {
		return 0, e
	}

	return a.Cmp(b), nil
}

// Compare two ints. Returns:
// -1 if a < b
// 0 if a == b
// 1 if a > b
func cmp(a, b int) int {
	if a < b {
		return -1
	} else if a > b {
		return 1
	}

	return 0
}

// Create a semver struct from a string.
// `version` must be of the form: a.b.c, where a, b, and c are integers
func NewSemver(version string) (*Semver, error) {
	s := Semver{}
	var e error
	matches := re.FindStringSubmatch(version)

	if len(matches) < MIN_MATCHES {
		return nil, fmt.Errorf("Invalid semantic version: %v", matches)
	}

	// match found:
	// [0]: full match (eg. 1.2.3)
	// [1]: major version
	// [2]: minor version
	// [3]: patch version
	s.Major, e = strconv.Atoi(matches[1])

	if e != nil {
		return nil, e
	}

	s.Minor, e = strconv.Atoi(matches[2])

	if e != nil {
		return nil, e
	}

	s.Patch, e = strconv.Atoi(matches[3])

	if e != nil {
		return nil, e
	}

	return &s, nil
}
