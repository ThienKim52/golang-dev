package repository

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type URLStorage interface {
	StoreURL(ctx context.Context, code, url string, exp time.Duration) error
	GetURL(ctx context.Context, code string) (string, error)
}

type urlStorage struct {
	c *redis.Client
}

func NewURLStorage(c *redis.Client) URLStorage {
	return &urlStorage{c: c}
}
func (s *urlStorage) StoreURL(ctx context.Context, code, url string, exp time.Duration) error {
	return s.c.Set(ctx, code, url, exp).Err()
}

func (s *urlStorage) GetURL(ctx context.Context, code string) (string, error) {
	return s.c.Get(ctx, code).Result()
}
