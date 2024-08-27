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
var databaseStore = storage.NewPostgresStorage(databaseClient)

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
		log.Printf("Looking for URL: %s\n", shortURL)
		urlData, err := databaseStore.GetURLDataIfExist(shortURL)
		if err != nil {
			log.Printf("URL not found in PostgreSQL, inserting new entry: %v\n", err)
			urlData = models.URLDataEntry{
				URL:          request.URL,
				ShortURL:     shortURL,
				CreatedAt:    time.Now(),
				LastAccessed: time.Now(),
				ViewCount:    0,
			}
			if err := databaseStore.InsertURL(urlData); err != nil {
				log.Printf("Error inserting URL into PostgreSQL: %v\n", err)
			}
		} else {
			log.Printf("Found short URL in PostgreSQL: %s\n", urlData.URL)
			if err := databaseStore.UpdateViewCountByShortURL(urlData.ShortURL); err != nil {
				log.Printf("Error updating view count in PostgreSQL: %v\n", err)
			}
		}
	}()

	response := models.URLResponse{ShortURL: shortURL}
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Could not encode response", http.StatusInternalServerError)
		log.Printf("Error encoding response: %v\n", err)
		return
	}
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
		if err := databaseStore.UpdateViewCountByShortURL(shortURL); err != nil {
			log.Printf("Error updating URL in PostgreSQL: %v\n", err)
		}
	}()

	http.Redirect(w, r, originalURL, http.StatusFound)
	log.Printf("Redirected short URL: %s -> %s\n", shortURL, originalURL)
}
