package main

import (
	"errors"
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: git <command> [<args>]")
		os.Exit(2)
	}

	command := os.Args[1]

	commandHandlers := []Command{
		InitCommand{},
		CatFileCommad{},
		HashObjectCommand{},
		LsTreeCommand{},
		WriteTreeCommand{},
		CommitTreeCommand{},
	}

	err := errors.New(fmt.Sprintf("Unknown command %s\n", command))

	for _, commandHandler := range commandHandlers {
		if commandHandler.Name() == command {
			err = commandHandler.Run(os.Args[2:])
		}
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
		os.Exit(1)
	}
}
