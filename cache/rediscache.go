package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"

	redis "github.com/redis/go-redis/v9"
)

var (
	// Client is the redis client instance
	rdb            *redis.Client
	onceRedisCache sync.Once
)

// InitRedis initializes the redis client
func InitRedis() error {
	var initErr error
	onceRedisCache.Do(func() {
		host := os.Getenv("REDIS_HOST")
		port := os.Getenv("REDIS_PORT")
		password := os.Getenv("REDIS_PASSWORD")
		username := os.Getenv("REDIS_USERNAME")

		if host == "" || port == "" {
			initErr = fmt.Errorf("redis host or port not set")
			return
		}

		client := redis.NewClient(&redis.Options{
			Addr:     fmt.Sprintf("%s:%s", host, port),
			Password: password, // No need for fmt.Sprintf
			DB:       0,
			Username: username,
		})

		// Test Redis connection before assigning to rdb
		if err := client.Ping(context.Background()).Err(); err != nil {
			initErr = fmt.Errorf("failed to connect to Redis: %w", err)
			return
		}

		// Assign the successfully connected client
		rdb = client
	})

	return initErr
}

// GetRedisData gets data from cache with JSON encoding
func GetRedisData[T any](ctx context.Context, key string) (T, error) {
	var zero T

	// Get data from cache
	res, err := rdb.Get(ctx, key).Result()
	if err == redis.Nil {
		return zero, fmt.Errorf("data not found in cache for key %s", key)
	}
	if err != nil {
		return zero, fmt.Errorf("failed to get data from cache: %w", err)
	}

	// Allocate memory for a pointer type
	var result T
	if err := json.Unmarshal([]byte(res), &result); err != nil {
		return zero, fmt.Errorf("failed to unmarshal data from cache: %w", err)
	}
	return result, nil
}

// SetRedisData sets data in cache with JSON encoding. 0 means no expiration. ttl is in seconds
func SetRedisData(ctx context.Context, key string, value any, ttl int32) error {
	// Convert value to JSON
	jsonData, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal data: %w", err)
	}

	// Set data in cache
	err = rdb.Set(ctx, key, jsonData, time.Duration(ttl)*time.Minute).Err()
	if err != nil {
		return fmt.Errorf("failed to set data in cache: %w", err)
	}
	return nil
}

// UpdateRedisData updates data in cache. 0 means no expiration. ttl is in seconds
func UpdateRedisData(ctx context.Context, key string, newValue interface{}) error {
	// Check if the key exists
	ttl, err := rdb.TTL(ctx, key).Result()
	if err != nil {
		return fmt.Errorf("failed to retrieve TTL for key %s: %w", key, err)
	}

	// If TTL is -1, the key does not exist
	if ttl == -1 {
		return fmt.Errorf("data not found in cache for key %s", key)
	}

	// Serialize new value to JSON
	data, err := json.Marshal(newValue)
	if err != nil {
		return fmt.Errorf("failed to marshal new data: %w", err)
	}

	// Update the key, preserving its TTL
	err = rdb.Set(ctx, key, data, ttl).Err()
	if err != nil {
		return fmt.Errorf("failed to update data in cache: %w", err)
	}

	return nil
}

// DeleteRedisData deletes data from cache
func DeleteRedisData(ctx context.Context, key string) error {
	// Attempt to delete the key
	deleted, err := rdb.Del(ctx, key).Result()
	if err != nil {
		return fmt.Errorf("failed to delete data from cache: %w", err)
	}

	// Check if the key existed
	if deleted <= 0 {
		return fmt.Errorf("data not found in cache for key %s", key)
	}

	return nil
}

// FlushAll removes all items from cache
func FlushAllRedis(ctx context.Context) error {
	if err := rdb.FlushAll(ctx).Err(); err != nil {
		return fmt.Errorf("failed to flush cache: %w", err)
	}
	return nil
}

// CloseRedis closes the Redis connection
func CloseRedis() error {
	if rdb != nil {
		return rdb.Close()
	}
	return nil
}
