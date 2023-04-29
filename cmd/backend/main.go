package main

import (
	"log"
	"net/http"
)

var commitHash string

func main() {
	log.Println("Starting backend")
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Hello, world!\nCommit hash: " + commitHash))
	})
	log.Println(http.ListenAndServe(":8080", nil))
}
