package main

const GIT_BASE = ".git"
const OBJECTS = "objects"
const REFS = "refs"
const HEAD = "HEAD"

type Command interface {
	Name() string
	Run([]string) error
}
