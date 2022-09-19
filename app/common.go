package main

const GIT_BASE = ".git"
const OBJECTS = "objects"
const REFS = "refs"
const HEAD = "HEAD"

// Command is an interface that each CLI sub-command implements
type Command interface {

	// Name returns sub-command name, for example "clone"
	Name() string

	// Run exectues subcommand
	// args is a slice of CLI arguments excluding binary name and subcommand
	Run(args []string) error
}
