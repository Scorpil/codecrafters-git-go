package main

import (
	"bufio"
	"bytes"
	"compress/zlib"
	"crypto/sha1"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
)

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
	type_, err := objectReader.ReadString(' ')
	if err != nil {
		return Object{}, err
	}
	type_ = type_[0 : len(type_)-1]

	sizeStr, err := objectReader.ReadString(0)
	if err != nil {
		return Object{}, err
	}

	size, err := strconv.Atoi(string(sizeStr[:len(sizeStr)-1]))
	if err != nil {
		return Object{}, err
	}

	var content = make([]byte, size, size)
	if _, err := objectReader.Read(content); err != nil {
		return Object{}, err
	}

	return Object{type_, content}, nil
}

func WriteObject(objectName string, content []byte) error {
	objectPath := getObjectPath(objectName)
	if err := os.MkdirAll(filepath.Dir(objectPath), 0755); err != nil {
		return err
	}
	return os.WriteFile(objectPath, content, 0644)
}

func (o Object) Marshal() ([]byte, []byte, error) {
	var b bytes.Buffer
	hasher := sha1.New()
	zlibWriter := zlib.NewWriter(&b)
	w := io.MultiWriter(zlibWriter, hasher)

	fmt.Fprintf(w, "%s %d", o.type_, len(o.content))
	w.Write([]byte{0})
	w.Write(o.content)
	zlibWriter.Close()

	return hasher.Sum(nil), b.Bytes(), nil
}
