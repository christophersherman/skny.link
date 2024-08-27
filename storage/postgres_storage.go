package storage

import (
	"context"
	"database/sql"
	"time"
	"url_shortener/models"
)

type PostgresStorage struct {
	client *sql.DB
	ctx    context.Context
}

func NewPostgresStorage(client *sql.DB) *PostgresStorage {
	return &PostgresStorage{
		client: client,
		ctx:    context.Background(),
	}
}

// InsertURL inserts the URL if it does not exist in the database already
func (p *PostgresStorage) InsertURL(urlData models.URLDataEntry) error {
	_, err := p.client.ExecContext(p.ctx, `
        INSERT INTO urls (url, short_url, created_at, last_accessed, view_count) 
        VALUES ($1, $2, $3, $4, $5)
        ON CONFLICT (short_url) DO NOTHING`, // Prevent duplicate entries
		urlData.URL, urlData.ShortURL, urlData.CreatedAt, urlData.LastAccessed, urlData.ViewCount)
	if err != nil {
		return err
	}
	return nil
}

func (p *PostgresStorage) UpdateViewCountByShortURL(shortenedURL string) error {
	_, err := p.client.ExecContext(p.ctx, `
        UPDATE urls 
        SET view_count = view_count + 1, last_accessed = $1 
        WHERE short_url = $2`,
		time.Now(), shortenedURL)
	if err != nil {
		return err
	}
	return nil
}

// GetURLDataIfExist gets the URL data if it exists in the database
func (p *PostgresStorage) GetURLDataIfExist(shortenedURL string) (models.URLDataEntry, error) {
	var urlData models.URLDataEntry
	err := p.client.QueryRowContext(p.ctx, `
        SELECT url, short_url, created_at, last_accessed, view_count 
        FROM urls 
        WHERE short_url = $1`, shortenedURL).Scan(
		&urlData.URL, &urlData.ShortURL, &urlData.CreatedAt, &urlData.LastAccessed, &urlData.ViewCount)
	if err != nil {
		return models.URLDataEntry{}, err
	}
	return urlData, nil
}
