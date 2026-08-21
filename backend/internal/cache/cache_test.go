package cache

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type TestUser struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

func getTestCache(t *testing.T) *Cache {
	t.Helper()
	c, err := New("redis://localhost:6379")
	if err != nil {
		t.Skipf("Skipping Redis test (Redis not reachable on localhost:6379): %v", err)
	}
	return c
}

func TestCache_SetAndGet(t *testing.T) {
	c := getTestCache(t)
	defer c.Close()

	ctx := context.Background()
	key := "test:user:1"
	user := TestUser{
		ID:       1,
		Username: "frank",
		Email:    "frank@example.com",
	}

	// Cleanup
	_ = c.Del(ctx, key)
	defer func() { _ = c.Del(ctx, key) }()

	// 1. Set in cache
	err := c.Set(ctx, key, user, 10*time.Second)
	require.NoError(t, err)

	// 2. Get from cache
	var retrieved TestUser
	err = c.Get(ctx, key, &retrieved)
	require.NoError(t, err)
	assert.Equal(t, user, retrieved)
}

func TestCache_CacheMiss(t *testing.T) {
	c := getTestCache(t)
	defer c.Close()

	ctx := context.Background()
	key := "test:nonexistent:key"

	var retrieved TestUser
	err := c.Get(ctx, key, &retrieved)
	require.Error(t, err)
	assert.ErrorIs(t, err, ErrCacheMiss)
}

func TestCache_Del(t *testing.T) {
	c := getTestCache(t)
	defer c.Close()

	ctx := context.Background()
	key := "test:user:delete"
	user := TestUser{ID: 2, Username: "bob"}

	err := c.Set(ctx, key, user, 10*time.Second)
	require.NoError(t, err)

	// Delete key
	err = c.Del(ctx, key)
	require.NoError(t, err)

	// Verify key is gone
	var retrieved TestUser
	err = c.Get(ctx, key, &retrieved)
	require.ErrorIs(t, err, ErrCacheMiss)
}

func TestCache_TTL_Expiration(t *testing.T) {
	c := getTestCache(t)
	defer c.Close()

	ctx := context.Background()
	key := "test:user:ttl"
	user := TestUser{ID: 3, Username: "charlie"}

	// Set with short TTL
	err := c.Set(ctx, key, user, 200*time.Millisecond)
	require.NoError(t, err)

	// Sleep past expiration
	time.Sleep(300 * time.Millisecond)

	var retrieved TestUser
	err = c.Get(ctx, key, &retrieved)
	require.ErrorIs(t, err, ErrCacheMiss)
}
