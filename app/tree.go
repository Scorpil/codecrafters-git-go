package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
)

const OBJECT_TYPE_TREE = "tree"

const (
	MODE_FILE = "100644" // normal file mode
	MODE_EXEC = "100755" // executable file mode
	MODE_LINK = "120000" // symbolic link mode
	MODE_DIR  = "040000" // directory mode
)

type TreeItem struct {
	Mode     string
	Filename string
	Hash     []byte
}

type Tree struct {
	Items []TreeItem
}

func (t Tree) Marshal() []byte {
	content := ""
	for _, item := range t.Items {
		content += fmt.Sprintf("%s %s\x00%s", item.Mode, item.Filename, item.Hash)
	}
	return []byte(content)
}

func FileInfoToGitMode(info os.FileInfo) string {
	if info.IsDir() {
		return MODE_DIR
	}
	// TODO: implement MODE_EXEC, MODE_LINK and submodule mode
	return MODE_FILE

}

func ReadTree(hash string) (Tree, error) {
	object, err := ReadObject(hash)
	if err != nil {
		return Tree{}, err
	}

	r := bufio.NewReader(bytes.NewReader(object.content))

	tree := Tree{
		make([]TreeItem, 0),
	}

	for {
		mode, err := r.ReadString(' ')
		if err == io.EOF {
			break
		} else if err != nil {
			return Tree{}, err
		}

		filename, err := r.ReadString(0)
		if err != nil {
			return Tree{}, err
		}

		targetHash := make([]byte, 20, 20)
		_, err = io.ReadFull(r, targetHash)
		if err != nil {
			return Tree{}, err
		}

		treeItem := TreeItem{
			mode[:len(mode)-1],         // Remove the trailing null byte
			filename[:len(filename)-1], // Remove the trailing null byte
			targetHash,
		}
		tree.Items = append(tree.Items, treeItem)
	}

	return tree, nil
}
