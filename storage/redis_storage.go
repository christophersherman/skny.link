package storage

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"time"

	"github.com/go-redis/redis/v8"
)

type RedisStorage struct {
	client *redis.Client
	ctx    context.Context
}

func NewRedisStorage(client *redis.Client) *RedisStorage {
	return &RedisStorage{
		client: client,
		ctx:    context.Background(),
	}
}

func (r *RedisStorage) SaveURL(originalURL string) (string, error) {
	urlHash := r.hashURL(originalURL)

	// Check if the hash already exists in Redis
	shortURL, err := r.client.Get(r.ctx, "hash:"+urlHash).Result()
	if err == redis.Nil { // Not found, create a new short URL
		shortURL = generateShortURL()

		// Store the hash to short URL mapping
		err = r.client.Set(r.ctx, "hash:"+urlHash, shortURL, 0).Err()
		if err != nil {
			return "", err
		}

		// Store the short URL to original URL mapping
		err = r.client.Set(r.ctx, "short:"+shortURL, originalURL, 0).Err()
		if err != nil {
			return "", err
		}
	} else if err != nil {
		return "", err
	}

	return shortURL, nil
}

func (r *RedisStorage) GetURL(shortURL string) (string, error) {
	return r.client.Get(r.ctx, "short:"+shortURL).Result()
}

func (r *RedisStorage) hashURL(url string) string {
	hash := sha256.Sum256([]byte(url))
	return hex.EncodeToString(hash[:])
}

// Generate a unique short URL
func generateShortURL() string {
	b := make([]byte, 4)
	_, err := rand.Read(b)
	if err != nil {
		panic(err) // handle error in production code
	}
	return hex.EncodeToString(b) + time.Now().Format("060102150405")
}
