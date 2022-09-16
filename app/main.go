package main

import (
	"fmt"
	"io/ioutil"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: git <command> [<args>]")
		os.Exit(0)
	}

	switch command := os.Args[1]; command {
	case "init":
		for _, dir := range []string{".git", ".git/objects", ".git/refs"} {
			if err := os.Mkdir(dir, 0755); err != nil {
				fmt.Printf("Error creating directory: %s\n", err)
			}
		}

		headFileContents := []byte("ref: refs/heads/master\n")
		if err := ioutil.WriteFile(".git/HEAD", headFileContents, 0644); err != nil {
			fmt.Printf("Error writing file: %s\n", err)
		}

		fmt.Println("Initialized git directory")

	default:
		fmt.Printf("Unknown command %s\n", command)
		os.Exit(1)
	}
}
