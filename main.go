package main

import (
	"github.com/gorilla/mux"
	"net/http"
	"password-encoder/handlers"
)

func main() {
	r := mux.NewRouter()
	h := handlers.InitializeHandler()

	r.HandleFunc("/hash", h.CreateHash).Methods(http.MethodPost)
	r.HandleFunc("/hash/{id}", h.GetHash).Methods(http.MethodGet)
	r.HandleFunc("/stats", h.GetStats).Methods(http.MethodGet)

	http.ListenAndServe(":8080", r)
}