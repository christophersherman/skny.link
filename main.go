package main

import (
	"log"
	"net/http"
	"url_shortener/handlers"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()

	//move routes or keep simple? hmm
	r.HandleFunc("/shorten", handlers.ShortenURL).Methods("POST")
	r.HandleFunc("/{shortURL}", handlers.RedirectURL).Methods("GET")

	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatalf("Could not start server: %s\n", err)
	}
}
