package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
	"url_shortener/config"
	"url_shortener/models"
	"url_shortener/storage"

	"github.com/gorilla/mux"
)

var redisClient = config.NewRedisClient()
var store = storage.NewRedisStorage(redisClient)
var databaseClient = config.NewPostgresClient()

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

	// Asynchronously insert the shortened URL into PostgreSQL
	go func() {
		urlData := models.URLDataEntry{
			URL:          request.URL,
			ShortURL:     shortURL,
			CreatedAt:    time.Now(),
			LastAccessed: time.Now(),
			ViewCount:    0,
		}
		_, err := databaseClient.Exec("INSERT INTO urls (url, short_url, created_at, last_accessed, view_count) VALUES ($1, $2, $3, $4, $5)",
			urlData.URL, urlData.ShortURL, urlData.CreatedAt, urlData.LastAccessed, urlData.ViewCount)
		if err != nil {
			log.Printf("Error inserting URL into PostgreSQL: %v\n", err)
		} else {
			log.Printf("Successfully inserted URL into PostgreSQL: %s -> %s\n", urlData.URL, urlData.ShortURL)
		}
	}()

	response := models.URLResponse{ShortURL: shortURL}
	json.NewEncoder(w).Encode(response)
	log.Printf("Successfully shortened URL: %s -> %s\n", request.URL, shortURL)
}

func RedirectURL(w http.ResponseWriter, r *http.Request) {
	log.Println("Received request to redirect URL")

	vars := mux.Vars(r)
	shortURL := vars["shortURL"]

	originalURL, err := store.GetURL(shortURL)
	if err != nil {
		http.Error(w, "URL not found", http.StatusNotFound)
		log.Printf("Error retrieving URL: %v\n", err)
		return
	}

	// Asynchronously update the view count and last accessed time in PostgreSQL
	go func() {
		_, err := databaseClient.Exec("UPDATE urls SET view_count = view_count + 1, last_accessed = $1 WHERE short_url = $2", time.Now(), shortURL)
		if err != nil {
			log.Printf("Error updating view count in PostgreSQL: %v\n", err)
		} else {
			log.Printf("Successfully updated view count for short URL: %s\n", shortURL)
		}
	}()

	http.Redirect(w, r, originalURL, http.StatusFound)
	log.Printf("Redirected short URL: %s -> %s\n", shortURL, originalURL)
}
