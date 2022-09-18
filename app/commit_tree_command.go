package main

import (
	"encoding/hex"
	"errors"
	"flag"
	"fmt"
	"time"
)

const OBJECT_TYPE_COMMIT = "commit"

type CommitTreeCommand struct{}

type commitTreeParams struct {
	TreeSha   []byte
	ParentSha []byte
	Message   string
}

func (c CommitTreeCommand) Name() string { return "commit-tree" }

func (c CommitTreeCommand) Run(args []string) error {
	p, err := c.parseFlags(args)
	if err != nil {
		return err
	}

	commitHash, err := c.runWithParams(p)
	if err != nil {
		return err
	}

	fmt.Println(hex.EncodeToString(commitHash))
	return nil
}

func (c CommitTreeCommand) parseFlags(args []string) (commitTreeParams, error) {
	if len(args) < 1 {
		return commitTreeParams{}, errors.New("Tree sha not found")
	}
	treeShaStr := args[0]
	if len(treeShaStr) != 40 {
		return commitTreeParams{}, errors.New(
			fmt.Sprintf("Expected a 40-char SHA as tree SHA. Got: %s\n", treeShaStr),
		)
	}

	commitTreeCmd := flag.NewFlagSet(c.Name(), flag.ExitOnError)
	parentShaStrPtr := commitTreeCmd.String("p", "", "SHA of a parent commit object")
	messagePtr := commitTreeCmd.String("m", "", "commit message")

	commitTreeCmd.Parse(args[1:])

	if len(*parentShaStrPtr) != 40 {
		return commitTreeParams{}, errors.New(
			fmt.Sprintf("Expected a 40-char SHA as parent SHA (-p value). Got: %s\n", *parentShaStrPtr),
		)
	}

	treeSha, err := hex.DecodeString(treeShaStr)
	if err != nil {
		return commitTreeParams{}, errors.New(fmt.Sprintf("Failed to parse tree SHA: %s", err.Error()))
	}

	parentSha, err := hex.DecodeString(*parentShaStrPtr)
	if err != nil {
		return commitTreeParams{}, errors.New(fmt.Sprintf("Failed to parse parent commit SHA: %s", err.Error()))
	}

	return commitTreeParams{treeSha, parentSha, *messagePtr}, nil
}

func (c CommitTreeCommand) runWithParams(p commitTreeParams) ([]byte, error) {

	commit := Commit{
		Author:    "Andrew Savchyn <dev@scorpil.com>",
		Timestamp: time.Now().Unix(),
		Timezone:  "+0200",
		TreeSha:   p.TreeSha,
		ParentSha: p.ParentSha,
		Message:   p.Message,
	}

	commitBytes := commit.Marshal()

	object := Object{OBJECT_TYPE_COMMIT, commitBytes}

	hash, objectBytes, err := object.Marshal()
	if err != nil {
		return nil, nil
	}
	hashStr := hex.EncodeToString(hash)
	WriteObject(hashStr, objectBytes)

	return hash, nil
}
