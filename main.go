package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

func main() {
	fmt.Println("Hello world")

	r := mux.NewRouter()

	r.HandleFunc("/hash", createHash).Methods(http.MethodPost)
	r.HandleFunc("/hash/{id}", getHash).Methods(http.MethodGet)
	r.HandleFunc("/stats", getStats).Methods(http.MethodGet)

	http.ListenAndServe(":8080", r)
}

func createHash(w http.ResponseWriter, r *http.Request) {
	fmt.Println("create hash and return request #")
}

func getHash(w http.ResponseWriter, r *http.Request) {
	fmt.Println("get hash")
}

func getStats(w http.ResponseWriter, r *http.Request) {
	fmt.Println(`get stats {"total": "number of requests made", "average": "average of time taken to process the requests"}`)
}

// must track id, hashedpw, average time to process the requests
