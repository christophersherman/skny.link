package main

import (
	"log"
	"net/http"
	"url_shortener/handlers"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()

	// Define the routes
	r.HandleFunc("/shorten", handlers.ShortenURL).Methods("POST")
	r.HandleFunc("/{shortURL}", handlers.RedirectURL).Methods("GET")

	// Start the server
	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatalf("Could not start server: %s\n", err)
	}
}
