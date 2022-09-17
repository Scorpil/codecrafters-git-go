package main

import (
	"bufio"
	"compress/zlib"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
)

const GIT_BASE = ".git"
const OBJECTS = "objects"
const REFS = "refs"
const HEAD = "HEAD"

type Object struct {
	type_   string
	content []byte
}

func ReadObject(objectName string) (Object, error) {
	objectPath := filepath.Join(GIT_BASE, OBJECTS, objectName[0:2], objectName[2:])
	if _, err := os.Stat(objectPath); err != nil {
		// object file does not exist
		return Object{}, err
	}

	objectFile, err := os.Open(objectPath)
	if err != nil {
		return Object{}, err
	}

	decodedReader, err := zlib.NewReader(objectFile)
	defer decodedReader.Close()
	if err != nil {
		return Object{}, err
	}

	objectReader := bufio.NewReader(decodedReader)
	typeBytes, err := objectReader.ReadBytes(' ')
	if err != nil {
		return Object{}, err
	}

	sizeBytes, err := objectReader.ReadBytes(0)
	if err != nil {
		return Object{}, err
	}

	size, err := strconv.Atoi(string(sizeBytes[:len(sizeBytes)-1]))
	if err != nil {
		return Object{}, err
	}

	var content = make([]byte, size, size)
	if _, err := objectReader.Read(content); err != nil {
		return Object{}, err
	}

	return Object{string(typeBytes), content}, nil
}

func main() {
	initCmd := flag.NewFlagSet("init", flag.ExitOnError)

	catFileCmd := flag.NewFlagSet("cat-file", flag.ExitOnError)
	prettyPrintPtr := catFileCmd.Bool("p", false, "pretty-print object's content")

	flag.Parse()

	if len(os.Args) < 2 {
		fmt.Println("Usage: git <command> [<args>]")
		os.Exit(2)
	}

	switch command := os.Args[1]; command {
	case initCmd.Name():
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
		if err := ioutil.WriteFile(filepath.Join(GIT_BASE, HEAD), headFileContents, 0644); err != nil {
			fmt.Fprintf(os.Stderr, "Error writing file: %s\n", err)
		}

		fmt.Println("Initialized git directory")
	case catFileCmd.Name():
		catFileCmd.Parse(os.Args[2:])

		tail := catFileCmd.Args()
		if len(tail) != 1 {
			catFileCmd.Usage()
			os.Exit(2)
		}
		objectName := tail[0]

		if *prettyPrintPtr {
			if object, err := ReadObject(objectName); err != nil {
				fmt.Fprintf(os.Stderr, "fatal: Not a valid object name %s\n", objectName)
				os.Exit(1)
			} else {
				fmt.Print(string(object.content))
			}
		}
	default:
		fmt.Fprintf(os.Stderr, "Unknown command %s\n", command)
		os.Exit(2)
	}
}
