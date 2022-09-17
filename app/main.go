package main

import (
	"bufio"
	"bytes"
	"compress/zlib"
	"crypto/sha1"
	"encoding/hex"
	"flag"
	"fmt"
	"io"
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

func getObjectPath(objectName string) string {
	return filepath.Join(GIT_BASE, OBJECTS, objectName[0:2], objectName[2:])
}

func ReadObject(objectName string) (Object, error) {
	objectPath := getObjectPath(objectName)
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

func WriteObject(objectName string, content []byte) error {
	objectPath := getObjectPath(objectName)
	if err := os.MkdirAll(filepath.Dir(objectPath), 0755); err != nil {
		return err
	}
	return os.WriteFile(objectPath, content, 0644)
}

func (o Object) Marshal() (string, []byte, error) {
	var b bytes.Buffer
	hasher := sha1.New()
	zlibWriter := zlib.NewWriter(&b)
	w := io.MultiWriter(zlibWriter, hasher)

	fmt.Fprintf(w, "blob %d", len(o.content))
	w.Write([]byte{0})
	w.Write(o.content)
	zlibWriter.Close()

	hashStr := hex.EncodeToString(hasher.Sum(nil))
	return hashStr, b.Bytes(), nil
}

func main() {
	initCmd := flag.NewFlagSet("init", flag.ExitOnError)

	catFileCmd := flag.NewFlagSet("cat-file", flag.ExitOnError)
	prettyPrintPtr := catFileCmd.Bool("p", false, "Pretty-print the contents of an object based on its type.")

	hashObjectCmd := flag.NewFlagSet("hash-object", flag.ExitOnError)
	writePtr := hashObjectCmd.Bool("w", false, "Actually write the object into the object database.")

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
		if err := os.WriteFile(filepath.Join(GIT_BASE, HEAD), headFileContents, 0644); err != nil {
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
	case hashObjectCmd.Name():
		hashObjectCmd.Parse(os.Args[2:])

		tail := hashObjectCmd.Args()
		if len(tail) != 1 {
			hashObjectCmd.Usage()
			os.Exit(2)
		}
		filePath := tail[0]

		fileContent, err := os.ReadFile(filePath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to read file: %s", err)
			os.Exit(1)
		}

		object := Object{"blob", fileContent}
		objectName, objectBytes, err := object.Marshal()
		if err != nil {
			fmt.Fprintf(os.Stderr, "error while encoding an object %s", err)
			os.Exit(1)
		}

		if *writePtr {
			if err := WriteObject(objectName, objectBytes); err != nil {
				fmt.Fprintf(os.Stderr, "error while writing an object %s", err)
				os.Exit(1)
			}
		}

		fmt.Println(objectName)

	default:
		fmt.Fprintf(os.Stderr, "Unknown command %s\n", command)
		os.Exit(2)
	}
}
