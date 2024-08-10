package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"url_shortener/config"
	"url_shortener/models"
	"url_shortener/storage"

	"github.com/go-redis/redis/v8"
	"github.com/gorilla/mux"
)

var redisClient = config.NewRedisClient()
var store = storage.NewRedisStorage(redisClient)

func ShortenURL(w http.ResponseWriter, r *http.Request) {
	log.Println("Received request to shorten URL")

	var request models.URLRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Printf("Failed to decode request: %v\n", err)
		return
	}

	log.Printf("Processing URL: %s\n", request.URL)
	shortURL, err := store.SaveURL(request.URL)
	if err != nil {
		http.Error(w, "Could not shorten URL", http.StatusInternalServerError)
		log.Printf("Error saving URL: %v\n", err)
		return
	}

	response := models.URLResponse{ShortURL: shortURL}
	json.NewEncoder(w).Encode(response)
	log.Printf("Successfully shortened URL: %s -> %s\n", request.URL, shortURL)
}

func RedirectURL(w http.ResponseWriter, r *http.Request) {
	log.Println("Received request to redirect URL")

	vars := mux.Vars(r)
	shortURL := vars["shortURL"]

	originalURL, err := store.GetURL(shortURL)
	if err == redis.Nil {
		http.NotFound(w, r)
		log.Printf("Short URL not found: %s\n", shortURL)
		return
	} else if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		log.Printf("Error retrieving URL: %v\n", err)
		return
	}

	http.Redirect(w, r, originalURL, http.StatusFound)
	log.Printf("Redirecting short URL: %s -> %s\n", shortURL, originalURL)
}
