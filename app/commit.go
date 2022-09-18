package main

import (
	"encoding/hex"
	"fmt"
)

type Commit struct {
	Author    string
	Timestamp int64
	Timezone  string
	TreeSha   []byte
	ParentSha []byte
	Message   string
}

func (c Commit) Marshal() []byte {
	content := fmt.Sprintf("tree %s\n", hex.EncodeToString(c.TreeSha))
	content += fmt.Sprintf("parent %s\n", hex.EncodeToString(c.ParentSha))
	content += fmt.Sprintf("author %s %d %s\n", c.Author, c.Timestamp, c.Timezone)
	content += fmt.Sprintf("committer %s %d %s\n", c.Author, c.Timestamp, c.Timezone)
	content += "\n"
	content += (c.Message + "\n")
	return []byte(content)
}
