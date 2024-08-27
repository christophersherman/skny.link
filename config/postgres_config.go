package config

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq" // Import the PostgreSQL driver
)

func NewPostgresClient() *sql.DB {
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("POSTGRES_HOST"),
		os.Getenv("POSTGRES_PORT"),
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_DB"))

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Could not connect to database: %v", err)
	}

	if err := db.Ping(); err != nil {
		log.Fatalf("Could not ping database: %v", err)
	}

	log.Println("Connected to PostgreSQL database")

	initDB(db) // Initialize the database

	return db
}

func initDB(db *sql.DB) {
	createTableQuery := `
	CREATE TABLE IF NOT EXISTS urls (
		id SERIAL PRIMARY KEY,
		url TEXT NOT NULL,
		short_url TEXT NOT NULL UNIQUE,
		created_at TIMESTAMP NOT NULL,
		last_accessed TIMESTAMP NOT NULL,
		view_count INT NOT NULL
	);`
	_, err := db.Exec(createTableQuery)
	if err != nil {
		log.Fatalf("Could not create urls table: %v", err)
	}

	log.Println("Ensured urls table exists")
}
