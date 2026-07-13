package repository

import (
	"context"
	"testing"

	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
)

func TestNewRedisLinkRepository(t *testing.T) {
	s, err := miniredis.Run()
	assert.NoError(t, err)
	defer s.Close()

	client := redis.NewClient(&redis.Options{
		Addr: s.Addr(),
	})
	defer client.Close()

	repo := NewRedisLinkRepository(client)
	assert.NotNil(t, repo)
	_, ok := repo.(*RedisLinkRepository)
	assert.True(t, ok)
}

func TestRedisLinkRepository_Save(t *testing.T) {
	s, err := miniredis.Run()
	assert.NoError(t, err)
	defer s.Close()

	client := redis.NewClient(&redis.Options{
		Addr: s.Addr(),
	})
	defer client.Close()

	repo := NewRedisLinkRepository(client)
	ctx := context.Background()

	err = repo.Save(ctx, "test", "http://google.com", 0)
	assert.NoError(t, err)

	val, err := s.Get("test")
	assert.NoError(t, err)
	assert.Equal(t, "http://google.com", val)
}

func TestRedisLinkRepository_GetByCode(t *testing.T) {
	s, err := miniredis.Run()
	assert.NoError(t, err)
	defer s.Close()

	client := redis.NewClient(&redis.Options{
		Addr: s.Addr(),
	})
	defer client.Close()

	repo := NewRedisLinkRepository(client)
	ctx := context.Background()

	s.Set("test", "http://google.com")

	val, err := repo.GetByCode(ctx, "test")
	assert.NoError(t, err)
	assert.Equal(t, "http://google.com", val)

	val, err = repo.GetByCode(ctx, "non-existent")
	assert.ErrorIs(t, err, redis.Nil)
	assert.Empty(t, val)
}

func TestRedisLinkRepository_Exists(t *testing.T) {
	s, err := miniredis.Run()
	assert.NoError(t, err)
	defer s.Close()

	client := redis.NewClient(&redis.Options{
		Addr: s.Addr(),
	})
	defer client.Close()

	repo := NewRedisLinkRepository(client)
	ctx := context.Background()

	exists, err := repo.Exists(ctx, "test")
	assert.NoError(t, err)
	assert.False(t, exists)

	s.Set("test", "http://google.com")

	exists, err = repo.Exists(ctx, "test")
	assert.NoError(t, err)
	assert.True(t, exists)
}

func TestRedisLinkRepository_Ping(t *testing.T) {
	s, err := miniredis.Run()
	assert.NoError(t, err)
	defer s.Close()

	client := redis.NewClient(&redis.Options{
		Addr: s.Addr(),
	})
	defer client.Close()

	repo := NewRedisLinkRepository(client)
	ctx := context.Background()

	err = repo.Ping(ctx)
	assert.NoError(t, err)
}
