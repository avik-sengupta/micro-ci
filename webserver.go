package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
)

type Response struct {
	CurrentTime string
}

func testFunc(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello!")
}

func clone(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		fmt.Fprintf(w, "ParseForm() err: %v", err)
		return
	}
	fmt.Fprintf(w, "POST request successful")
	url := r.FormValue("cloneurl")
	directory := r.FormValue("directory")
	pk := r.FormValue("pk")

	fmt.Fprintf(w, "URL = %s\n", url)
	fmt.Fprintf(w, "Directory = %s\n", directory)
	fmt.Fprintf(w, "PK = %s\n", pk)

	publicKeys, err := ssh.NewPublicKeysFromFile("git", pk, "")
	if err != nil {
		fmt.Println("generate publickeys failed: %s\n", err.Error())
		return
	}

	req, err := git.PlainClone(directory, false, &git.CloneOptions{
		Auth:     publicKeys,
		URL:      url,
		Progress: os.Stdout,
	})

	fmt.Println(req, err)

}

func push(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		fmt.Fprintf(w, "ParseForm() err: %v", err)
		return
	}
	fmt.Fprintf(w, "POST request successful")
	directory := r.FormValue("directory")

	fmt.Fprintf(w, "Directory = %s\n", directory)

	// Opens an already existing repository.
	req, err := git.PlainOpen(directory)

	fmt.Println(req, err)

	workTree, err := req.Worktree()

	//Info("echo \"hello world!\" > example-git-file")
	filename := filepath.Join(directory, "git-push-example")
	err = ioutil.WriteFile(filename, []byte("first git push!"), 0644)

	// Adds the new file to the staging area.
	//Info("git add example-git-file")
	_, err = workTree.Add("example-git-file")

	// We can verify the current status of the worktree using the method Status.
	//Info("git status --porcelain")
	status, err := workTree.Status()

	fmt.Println(status)

	// Commits the current staging area to the repository, with the new file
	// just created. We should provide the object.Signature of Author of the
	// commit Since version 5.0.1, we can omit the Author signature, being read
	// from the git config files.
	//Info("git commit -m \"example go-git commit\"")
	commit, err := workTree.Commit("test go-git commit", &git.CommitOptions{
		Author: &object.Signature{
			Name:  "Avik Sengupta",
			Email: "avik.sengupta27@gmail.com",
			When:  time.Now(),
		},
	})

	// Prints the current HEAD to verify that all worked well.
	//Info("git show -s")
	obj, err := req.CommitObject(commit)

	fmt.Println(obj)

	// push using default options
	err = req.Push(&git.PushOptions{})

	fmt.Println(err)

}

func main() {
	http.Handle("/", http.FileServer(http.Dir("web")))

	http.HandleFunc("/get-time", func(rw http.ResponseWriter, r *http.Request) {
		ctime := Response{
			CurrentTime: time.Now().Format(time.RFC3339),
		}
		byteArray, err := json.Marshal(ctime)
		if err != nil {
			fmt.Println(err)
		}
		rw.Write(byteArray)
	})

	http.HandleFunc("/hello", testFunc)

	http.HandleFunc("/clone", clone)
	http.HandleFunc("/push", push)

	if err := http.ListenAndServe(":5000", nil); err != nil {
		log.Fatal(err)
	}
}
