package cache

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"sync"

	"github.com/bradfitz/gomemcache/memcache"
)

var (
	// Client is the memcache client instance
	Client                 *memcache.Client
	once                   sync.Once
	ErrCacheNotInitialized = errors.New("cache client not initialized")
)

// InitMemcache initializes the memcache client
func InitMemcache() error {
	var initErr error
	once.Do(func() {
		host := os.Getenv("MEMCACHE_HOST")
		port := os.Getenv("MEMCACHE_PORT")
		if host == "" || port == "" {
			initErr = fmt.Errorf("memcache host or port not set")
			return
		}
		Client = memcache.New(fmt.Sprintf("%s:%s", host, port))
		// Test connection
		if err := Client.Ping(); err != nil {
			initErr = fmt.Errorf("failed to connect to memcache: %w", err)
			Client = nil
			return
		}
	})
	return initErr
}

// SetDataInCache sets data in cache with JSON encoding. 0 means no expiration. ttl is in seconds
func SetDataInCache(key string, value interface{}, ttl int32) error {
	if Client == nil {
		if err := InitMemcache(); err != nil {
			return fmt.Errorf("failed to initialize cache: %w", err)
		}
	}
	if key == "" {
		return fmt.Errorf("key cannot be empty")
	}

	// Convert value to JSON
	jsonData, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal data: %w", err)
	}

	err = Client.Set(&memcache.Item{
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
	var empty T
	if Client == nil {
		if err := InitMemcache(); err != nil {
			return empty, fmt.Errorf("failed to initialize cache: %w", err)
		}
	}
	if key == "" {
		return empty, fmt.Errorf("key cannot be empty")
	}

	item, err := Client.Get(key)
	if err != nil {
		if err == memcache.ErrCacheMiss {
			return empty, fmt.Errorf("data not found in cache for key %s", key)
		}
		return empty, fmt.Errorf("failed to get data from cache: %w", err)
	}

	var result T
	if err := json.Unmarshal(item.Value, &result); err != nil {
		return empty, fmt.Errorf("failed to unmarshal data from cache: %w", err)
	}
	return result, nil
}

// DeleteDataFromCache deletes data from cache
func DeleteDataFromCache(key string) error {
	if Client == nil {
		if err := InitMemcache(); err != nil {
			return fmt.Errorf("failed to initialize cache: %w", err)
		}
	}
	if key == "" {
		return fmt.Errorf("key cannot be empty")
	}

	err := Client.Delete(key)
	if err != nil && err != memcache.ErrCacheMiss {
		return fmt.Errorf("failed to delete data from cache: %w", err)
	}
	return nil
}

// UpdateDataInCache updates data in cache. 0 means no expiration. ttl is in seconds
func UpdateDataInCache(key string, value interface{}, ttl int32) error {
	if Client == nil {
		if err := InitMemcache(); err != nil {
			return fmt.Errorf("failed to initialize cache: %w", err)
		}
	}
	if key == "" {
		return fmt.Errorf("key cannot be empty")
	}

	// Convert value to JSON
	jsonData, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal data: %w", err)
	}

	err = Client.Replace(&memcache.Item{
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
	if Client == nil {
		if err := InitMemcache(); err != nil {
			return fmt.Errorf("failed to initialize cache: %w", err)
		}
	}

	if err := Client.FlushAll(); err != nil {
		return fmt.Errorf("failed to flush cache: %w", err)
	}
	return nil
}
