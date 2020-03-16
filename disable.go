package main

import (
	"modtorio/modlist"
)

func disable(options []string) error {
	return modlist.SetStatus(FLAGS.dir, false, options)
}
