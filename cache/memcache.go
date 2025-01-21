package cache

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/bradfitz/gomemcache/memcache"
)

var (
	// Client is the memcache client instance
	memCacheClient *memcache.Client
	onceMemcache   sync.Once
)

// validate cache options
func validateCacheOptions(ttl int32, key string) error {
	if ttl < 0 {
		return fmt.Errorf("ttl cannot be negative")
	}
	if key == "" || len(key) > 50 {
		return fmt.Errorf("key cannot be empty or longer than 50 characters")
	}
	if memCacheClient == nil {
		if err := InitMemcache(); err != nil {
			return fmt.Errorf("failed to initialize cache: %w", err)
		}
		return fmt.Errorf("cache client not initialized")
	}
	return nil
}

// InitMemcache initializes the memcache client
func InitMemcache() error {
	var initErr error
	// Initialize client once
	onceMemcache.Do(func() {
		host := os.Getenv("MEMCACHE_HOST")
		port := os.Getenv("MEMCACHE_PORT")
		if host == "" || port == "" {
			initErr = fmt.Errorf("memcache host or port not set")
			return
		}
		memCacheClient = memcache.New(fmt.Sprintf("%s:%s", host, port))
		memCacheClient.Timeout = 2 * time.Second // Set reasonable default timeout
		// Test connection
		if err := memCacheClient.Ping(); err != nil {
			initErr = fmt.Errorf("failed to connect to memcache: %w", err)
			memCacheClient = nil
			return
		}
	})
	return initErr
}

// SetDataInCache sets data in cache with JSON encoding. 0 means no expiration. ttl is in seconds
func SetDataInCache(key string, value interface{}, ttl int32) error {

	// Validate cache options
	if err := validateCacheOptions(ttl, key); err != nil {
		return err
	}

	// Convert value to JSON
	jsonData, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal data: %w", err)
	}

	// Set data in cache
	err = memCacheClient.Set(&memcache.Item{
		Key:        key,
		Value:      jsonData,
		Expiration: ttl,
	})
	if err != nil {
		return fmt.Errorf("failed to set data in cache: %w", err)
	}
	return nil
}

// GetDataFromCache gets data from cache and returns the result by type T
func GetDataFromCache[T any](key string) (T, error) {
	if memCacheClient == nil {
		if err := InitMemcache(); err != nil {
			return *new(T), fmt.Errorf("failed to initialize cache: %w", err)
		}
		return *new(T), fmt.Errorf("cache client not initialized")

	}
	if key == "" {
		return *new(T), fmt.Errorf("key cannot be empty")
	}

	item, err := memCacheClient.Get(key)
	if err != nil {
		if err == memcache.ErrCacheMiss {
			return *new(T), fmt.Errorf("data not found in cache for key %s", key)
		}
		return *new(T), fmt.Errorf("failed to get data from cache: %w", err)
	}

	var result T
	if err := json.Unmarshal(item.Value, &result); err != nil {
		return *new(T), fmt.Errorf("failed to unmarshal data from cache: %w", err)
	}
	return result, nil
}

// DeleteDataFromCache deletes data from cache
func DeleteDataFromCache(key string) error {
	if memCacheClient == nil {
		if err := InitMemcache(); err != nil {
			return fmt.Errorf("failed to initialize cache: %w", err)
		}
		return fmt.Errorf("cache client not initialized")
	}
	if key == "" {
		return fmt.Errorf("key cannot be empty")
	}

	err := memCacheClient.Delete(key)
	if err != nil && err != memcache.ErrCacheMiss {
		return fmt.Errorf("failed to delete data from cache: %w", err)
	}
	return nil
}

// UpdateDataInCache updates data in cache. 0 means no expiration. ttl is in seconds
func UpdateDataInCache(key string, value interface{}, ttl int32) error {
	// Validate cache options
	if err := validateCacheOptions(ttl, key); err != nil {
		return err
	}

	// Convert value to JSON
	jsonData, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal data: %w", err)
	}

	err = memCacheClient.Replace(&memcache.Item{
		Key:        key,
		Value:      jsonData,
		Expiration: ttl,
	})
	if err != nil {
		if err == memcache.ErrCacheMiss {
			return fmt.Errorf("no existing cache found for key %s", key)
		}
		return fmt.Errorf("failed to update data in cache: %w", err)
	}
	return nil
}

// FlushAll removes all items from cache
func FlushAll() error {
	if memCacheClient == nil {
		if err := InitMemcache(); err != nil {
			return fmt.Errorf("failed to initialize cache: %w", err)
		}
		return fmt.Errorf("cache client not initialized")
	}

	if err := memCacheClient.FlushAll(); err != nil {
		return fmt.Errorf("failed to flush cache: %w", err)
	}
	return nil
}
