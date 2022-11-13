package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
)

type Response struct {
	CurrentTime string
}

func testFunc(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello!")
}

func cloneFunc(w http.ResponseWriter, r *http.Request) {
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

	http.HandleFunc("/clone", cloneFunc)

	if err := http.ListenAndServe(":5000", nil); err != nil {
		log.Fatal(err)
	}
}
