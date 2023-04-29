package main

import (
	"fmt"
	"log"
	"net/http"
)

var commitHash string

func main() {
	log.Println("Starting backend")
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)

		link := fmt.Sprintf("<a href=\"https://github.com/mdma-backend/mdma-backend/tree/%s\">%s</a>", commitHash, commitHash)

		w.Write([]byte("<html><body>Hello, world!<br />Commit hash: " + link + "</body></html>"))
	})
	log.Println(http.ListenAndServe(":8080", nil))
}
