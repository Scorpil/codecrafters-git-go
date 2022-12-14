package main

import (
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

type WriteTreeCommand struct{}

type writeTreeParams struct {
	dir string
}

func (c WriteTreeCommand) Name() string { return "write-tree" }

func (c WriteTreeCommand) Run(args []string) error {
	wd, err := os.Getwd()
	if err != nil {
		return errors.New(fmt.Sprintf("Faield to find current working directory: %s", err.Error()))
	}

	hash, err := c.runWithParams(writeTreeParams{wd})
	if err == nil {
		fmt.Println(hex.EncodeToString(hash))
	}
	return err
}

func (c WriteTreeCommand) runWithParams(p writeTreeParams) ([]byte, error) {
	fileInfos, err := readDir(p.dir)
	if err != nil {
		return nil, err
	}

	treeItems := make([]TreeItem, 0, len(fileInfos))
	for _, fileInfo := range fileInfos {
		gitMode := FileInfoToGitMode(fileInfo)
		var hash []byte
		if gitMode == MODE_FILE {
			hashObjectCommand := HashObjectCommand{}
			hash, err = hashObjectCommand.runWithParams(hashObjectParams{
				Write:    true,
				FilePath: filepath.Join(p.dir, fileInfo.Name()),
			})
			if err != nil {
				return nil, err
			}
		}
		if gitMode == MODE_DIR {
			writeTreeCommand := WriteTreeCommand{}
			hash, err = writeTreeCommand.runWithParams(writeTreeParams{
				dir: filepath.Join(p.dir, fileInfo.Name()),
			})
			if err != nil {
				return nil, err
			}
		}
		treeItems = append(treeItems, TreeItem{
			Mode:     gitMode,
			Filename: fileInfo.Name(),
			Hash:     hash,
		})
	}

	tree := Tree{treeItems}
	treeBytes := tree.Marshal()

	object := Object{OBJECT_TYPE_TREE, treeBytes}

	hash, objectBytes, err := object.Marshal()
	if err != nil {
		return nil, nil
	}
	hashStr := hex.EncodeToString(hash)
	WriteObject(hashStr, objectBytes)

	return hash, nil
}

// readDir reads a directory and returns and array of FileInfo structs
func readDir(dir string) ([]os.FileInfo, error) {
	dirContents, err := os.ReadDir(dir)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Failed to read directory contents: %s", err.Error()))
	}

	fileInfos := make([]os.FileInfo, 0, len(dirContents))
	for _, dirEntry := range dirContents {
		fileInfo, err := dirEntry.Info()
		if err != nil {
			return nil, err
		}

		// do not include .git directory
		if fileInfo.Name() != GIT_BASE {
			fileInfos = append(fileInfos, fileInfo)
		}
	}
	return fileInfos, nil
}
