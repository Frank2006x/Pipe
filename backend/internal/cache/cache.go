
package cache

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"
)

// ErrCacheMiss is returned when a key is not found in the cache.
var ErrCacheMiss = errors.New("cache miss")

// Cache wraps a Redis client and exposes typed get/set/del helpers.
type Cache struct {
	client *redis.Client
}

// New creates a Cache from the given Redis URL (e.g. "redis://localhost:6379").
func New(redisURL string) (*Cache, error) {
	opts, err := redis.ParseURL(redisURL)
	if err != nil {
		return nil, err
	}
	client := redis.NewClient(opts)

	// Ping to verify connectivity on startup.
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, err
	}

	return &Cache{client: client}, nil
}

// Get retrieves and JSON-unmarshals a value from the cache into dst.
// Returns ErrCacheMiss if the key does not exist.
func (c *Cache) Get(ctx context.Context, key string, dst any) error {
	data, err := c.client.Get(ctx, key).Bytes()
	if errors.Is(err, redis.Nil) {
		return ErrCacheMiss
	}
	if err != nil {
		return err
	}
	return json.Unmarshal(data, dst)
}

// Set JSON-marshals value and stores it under key with the given TTL.
func (c *Cache) Set(ctx context.Context, key string, value any, ttl time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return c.client.Set(ctx, key, data, ttl).Err()
}

// Del removes one or more keys from the cache.
func (c *Cache) Del(ctx context.Context, keys ...string) error {
	return c.client.Del(ctx, keys...).Err()
}

// Close closes the underlying Redis connection.
func (c *Cache) Close() error {
	return c.client.Close()
}
