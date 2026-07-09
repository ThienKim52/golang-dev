package repository

import (
	"context"
	"reflect"
	"time"

	"github.com/redis/go-redis/v9"
)

// LinkRepository defines the interface for link storage operations
type LinkRepository interface {
	Save(ctx context.Context, code, url string, exp time.Duration) error
	GetByCode(ctx context.Context, code string) (string, error)
	Exists(ctx context.Context, code string) (bool, error)
	Ping(ctx context.Context) error
}

type redisClient interface {
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd
	Get(ctx context.Context, key string) *redis.StringCmd
	Exists(ctx context.Context, keys ...string) *redis.IntCmd
	Ping(ctx context.Context) *redis.StatusCmd
}

// RedisLinkRepository implements LinkRepository using Redis
type RedisLinkRepository struct {
	c redisClient
}

func isNil(c redisClient) bool {
	if c == nil {
		return true
	}
	val := reflect.ValueOf(c)
	if val.Kind() == reflect.Ptr && val.IsNil() {
		return true
	}
	return false
}

// NewRedisLinkRepository creates a new Redis link repository
func NewRedisLinkRepository(c redisClient) LinkRepository {
	if isNil(c) {
		return &RedisLinkRepository{c: nil}
	}
	return &RedisLinkRepository{c: c}
}

// Save saves a link to Redis with expiration
func (r *RedisLinkRepository) Save(ctx context.Context, code, url string, exp time.Duration) error {
	if r.c == nil {
		return nil
	}
	return r.c.Set(ctx, code, url, exp).Err()
}

// GetByCode retrieves a link by its code
func (r *RedisLinkRepository) GetByCode(ctx context.Context, code string) (string, error) {
	if r.c == nil {
		return "", nil
	}
	return r.c.Get(ctx, code).Result()
}

// Exists checks if a code already exists in Redis
func (r *RedisLinkRepository) Exists(ctx context.Context, code string) (bool, error) {
	if r.c == nil {
		return false, nil
	}
	// TODO: Implement actual Redis existence check
	// Example implementation:
	exists, err := r.c.Exists(ctx, code).Result()
	return exists > 0, err
}

// Ping checks Redis connection
func (r *RedisLinkRepository) Ping(ctx context.Context) error {
	if r.c == nil {
		return nil
	}
	// TODO: Implement actual Redis ping
	// Example implementation:
	return r.c.Ping(ctx).Err()
}
