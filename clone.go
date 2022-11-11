package main

import (
	"fmt"
	"os"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
)

func main() {

	url, directory, privateKey := os.Args[1], os.Args[2], os.Args[3]
	fmt.Println(url, directory, privateKey)

	publicKeys, err := ssh.NewPublicKeysFromFile("git", privateKey, "")
	if err != nil {
		fmt.Println("generate publickeys failed: %s\n", err.Error())
		return
	}

	r, err := git.PlainClone(directory, false, &git.CloneOptions{
		Auth:     publicKeys,
		URL:      url,
		Progress: os.Stdout,
	})

	fmt.Println(r, err)
	// ... retrieving the branch being pointed by HEAD
	ref, err := r.Head()

	fmt.Println(ref, err)
	// ... retrieving the commit object
	commit, err := r.CommitObject(ref.Hash())

	fmt.Println(commit, err)
}
