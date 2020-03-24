package main

import (
	"modtorio/modlist"
)

func enable(flags *ModtorioFlags, options []string) error {
	return modlist.SetStatus(flags.dir, true, options)
}
