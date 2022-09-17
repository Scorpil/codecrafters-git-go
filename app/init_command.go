package main

import (
	"fmt"
	"os"
	"path/filepath"
)

type InitCommand struct{}

func (c InitCommand) Name() string { return "init" }

func (c InitCommand) Run(_ []string) error {
	for _, dir := range []string{
		GIT_BASE,
		filepath.Join(GIT_BASE, OBJECTS),
		filepath.Join(GIT_BASE, REFS),
	} {
		if err := os.Mkdir(dir, 0755); err != nil {
			fmt.Fprintf(os.Stderr, "Error creating directory: %s\n", err)
		}
	}

	headFileContents := []byte("ref: refs/heads/master\n")
	if err := os.WriteFile(filepath.Join(GIT_BASE, HEAD), headFileContents, 0644); err != nil {
		fmt.Fprintf(os.Stderr, "Error writing file: %s\n", err)
	}

	fmt.Println("Initialized git directory")
	return nil
}
