package main

import (
	"bufio"
	"bytes"
	"encoding/hex"
	"io"
)

type Tree struct {
	Items []TreeItem
	Hash  string
}

type TreeItem struct {
	Mode     string
	Filename string
	Hash     string
}

func ReadTree(hash string) (Tree, error) {
	object, err := ReadObject(hash)
	if err != nil {
		return Tree{}, err
	}

	r := bufio.NewReader(bytes.NewReader(object.content))

	tree := Tree{
		make([]TreeItem, 0),
		hash,
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
			mode[:len(mode)-1],
			filename[:len(filename)-1],
			hex.EncodeToString([]byte(targetHash))}
		tree.Items = append(tree.Items, treeItem)
	}

	return tree, nil
}
