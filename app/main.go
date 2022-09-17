package main

import (
	"errors"
	"fmt"
	"os"
)

const GIT_BASE = ".git"
const OBJECTS = "objects"
const REFS = "refs"
const HEAD = "HEAD"

type Command interface {
	Name() string
	Run([]string) error
}

func main() {

	initCommand := InitCommand{}
	catFileCommad := CatFileCommad{}
	HashObjectCommand := HashObjectCommand{}

	if len(os.Args) < 2 {
		fmt.Println("Usage: git <command> [<args>]")
		os.Exit(2)
	}

	var err error

	switch command := os.Args[1]; command {
	case initCommand.Name():
		err = initCommand.Run(os.Args[2:])
	case catFileCommad.Name():
		err = catFileCommad.Run(os.Args[2:])
	case HashObjectCommand.Name():
		err = HashObjectCommand.Run(os.Args[2:])
	default:
		err = errors.New(fmt.Sprintf("Unknown command %s\n", command))
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
		os.Exit(1)
	}
}
