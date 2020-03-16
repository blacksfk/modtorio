package main

import (
	"modtorio/modlist"
)

func enable(options []string) error {
	return modlist.SetStatus(FLAGS.dir, true, options)
}
