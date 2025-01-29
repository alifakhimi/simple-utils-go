/*
Package cache provides a Redis-based caching layer with utility functions for
connecting to Redis, and storing, retrieving, and deleting JSON-encoded objects.

Key Features:
1. **Connect**:
   - Establishes a connection to a Redis server using a URL string.

2. **Set**:
   - Stores a JSON-encoded object in Redis with an optional expiration time.

3. **Get**:
   - Retrieves and decodes a JSON-encoded object from Redis.

4. **Del**:
   - Deletes an object from Redis by its key.

Usage:
- Use `Connect` to initialize the Redis client.
- Use `Set`, `Get`, and `Del` for cache operations.

*/

package cache

import (
	"context"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"

	simutils "github.com/alifakhimi/simple-utils-go"
)

var (
	client *redis.Client // Redis client instance
)

func IsConnect(ctx context.Context) bool {
	return client != nil // && client.Ping(ctx).Err() == nil
}

// Connect initializes the Redis client using the given connection string.
func Connect(str string) (*redis.Client, error) {
	if opt, err := redis.ParseURL(str); err != nil {
		return nil, err
	} else {
		client = redis.NewClient(opt)
	}
	return client, nil
}

// Set stores a JSON-encoded object in Redis with an expiration time.
func Set(ctx context.Context, val any, exp time.Duration) error {
	key := simutils.GetTKey(val) // Generate the key using TKey logic.
	if b, err := json.Marshal(val); err != nil {
		return err
	} else if _, err := client.Set(ctx, string(key), b, exp).Result(); err != nil {
		return err
	}
	return nil
}

// Get retrieves a JSON-encoded object from Redis and decodes it into `dst`.
func Get(ctx context.Context, key simutils.TKey, dst any) error {
	if res, err := client.Get(ctx, string(key)).Result(); err != nil {
		return err
	} else if err := json.Unmarshal([]byte(res), dst); err != nil {
		return err
	}
	return nil
}

// Del deletes a cached object from Redis using its key.
func Del(ctx context.Context, key simutils.TKey) error {
	if _, err := client.Del(ctx, string(key)).Result(); err != nil {
		return err
	}
	return nil
}
