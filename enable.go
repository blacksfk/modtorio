package main

import (
	"github.com/blacksfk/modtorio/modlist"
)

func enable(flags *ModtorioFlags, options []string) error {
	return modlist.SetStatus(flags.dir, true, options)
}
