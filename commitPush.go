package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
)

func main() {

	directory := os.Args[1]
	fmt.Println(directory)

	// Opens an already existing repository.
	r, err := git.PlainOpen(directory)
	fmt.Println(err)

	w, err := r.Worktree()

	//Info("echo \"hello world!\" > example-git-file")
	filename := filepath.Join(directory, "example-git-file")
	err = ioutil.WriteFile(filename, []byte("hello world!"), 0644)

	// Adds the new file to the staging area.
	//Info("git add example-git-file")
	_, err = w.Add("example-git-file")

	// We can verify the current status of the worktree using the method Status.
	//Info("git status --porcelain")
	status, err := w.Status()

	fmt.Println(status)

	// Commits the current staging area to the repository, with the new file
	// just created. We should provide the object.Signature of Author of the
	// commit Since version 5.0.1, we can omit the Author signature, being read
	// from the git config files.
	//Info("git commit -m \"example go-git commit\"")
	commit, err := w.Commit("example go-git commit", &git.CommitOptions{
		Author: &object.Signature{
			Name:  "Avik Sengupta",
			Email: "avik.sengupta27@gmail.com",
			When:  time.Now(),
		},
	})

	// Prints the current HEAD to verify that all worked well.
	//Info("git show -s")
	obj, err := r.CommitObject(commit)

	fmt.Println(obj)

	// push using default options
	err = r.Push(&git.PushOptions{})
}
