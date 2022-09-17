package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
)

type catFileParams struct {
	PrettyPrint bool
	ObjectName  string
}

type CatFileCommad struct{}

func (c CatFileCommad) Name() string { return "cat-file" }

func (c CatFileCommad) Run(args []string) error {
	return c.runWithParams(c.parseFlags(args))
}

func (c CatFileCommad) parseFlags(args []string) catFileParams {
	catFileCmd := flag.NewFlagSet(c.Name(), flag.ExitOnError)
	prettyPrintPtr := catFileCmd.Bool("p", false, "Pretty-print the contents of an object based on its type.")

	catFileCmd.Parse(args)

	tail := catFileCmd.Args()
	if len(tail) != 1 {
		catFileCmd.Usage()
		os.Exit(2)
	}
	objectName := tail[0]

	return catFileParams{*prettyPrintPtr, objectName}
}

func (c CatFileCommad) runWithParams(p catFileParams) error {
	if p.PrettyPrint {
		if object, err := ReadObject(p.ObjectName); err != nil {
			return errors.New(fmt.Sprintf("fatal: Not a valid object name %s\n", p.ObjectName))
		} else {
			fmt.Print(string(object.content))
		}
	}
	return nil
}
