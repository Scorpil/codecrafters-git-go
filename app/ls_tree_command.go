package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
)

type lsTreeParams struct {
	nameOnly bool
	hash     string
}

type LsTreeCommand struct{}

func (c LsTreeCommand) Name() string { return "ls-tree" }

func (c LsTreeCommand) Run(args []string) error {
	p := c.parseFlags(args)
	tree, err := ReadTree(p.hash)

	if err != nil {
		return errors.New(fmt.Sprintf("Failed to read a tree %s\n", p.hash))
	}

	for _, item := range tree.Items {
		if p.nameOnly {
			fmt.Println(item.Filename)
		}
		// TODO: implement default output (without --names-only flag)
	}

	return nil
}

func (c LsTreeCommand) parseFlags(args []string) lsTreeParams {
	lsTreeCmd := flag.NewFlagSet(c.Name(), flag.ExitOnError)
	nameOnlyPtr := lsTreeCmd.Bool("name-only", false, "List only filenames (instead of the \"long\" output), one per line.")

	lsTreeCmd.Parse(args)

	tail := lsTreeCmd.Args()
	if len(tail) != 1 {
		lsTreeCmd.Usage()
		os.Exit(2)
	}
	hash := tail[0]

	return lsTreeParams{*nameOnlyPtr, hash}
}
