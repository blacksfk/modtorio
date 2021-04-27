package main

import (
	"github.com/blacksfk/modtorio/modlist"
)

func disable(flags *ModtorioFlags, options []string) error {
	return modlist.SetStatus(flags.dir, false, options)
}
