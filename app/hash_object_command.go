package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
)

type hashObjectParams struct {
	Write    bool
	FilePath string
}

type HashObjectCommand struct{}

func (c HashObjectCommand) Name() string { return "hash-object" }

func (c HashObjectCommand) Run(args []string) error {
	objectName, err := c.runWithParams(c.parseFlags(args))
	if err == nil {
		fmt.Println(objectName)
	}
	return err
}

func (c HashObjectCommand) parseFlags(args []string) hashObjectParams {
	hashObjectCmd := flag.NewFlagSet(c.Name(), flag.ExitOnError)
	writePtr := hashObjectCmd.Bool("w", false, "Actually write the object into the object database.")

	hashObjectCmd.Parse(args)

	tail := hashObjectCmd.Args()
	if len(tail) != 1 {
		hashObjectCmd.Usage()
		os.Exit(2)
	}
	filePath := tail[0]

	return hashObjectParams{*writePtr, filePath}
}

func (c HashObjectCommand) runWithParams(p hashObjectParams) (string, error) {
	fileContent, err := os.ReadFile(p.FilePath)
	if err != nil {
		return "", errors.New(fmt.Sprintf("Failed to read file: %s", err.Error()))
	}

	object := Object{"blob", fileContent}
	objectName, objectBytes, err := object.Marshal()
	if err != nil {
		return "", errors.New(fmt.Sprintf("Error while encoding an object %s", err.Error()))
	}

	if p.Write {
		if err := WriteObject(objectName, objectBytes); err != nil {
			return "", errors.New(fmt.Sprintf("Error while writing an object %s", err.Error()))
		}
	}

	return objectName, nil
}
