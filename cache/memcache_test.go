package cache

import (
	"os"
	"sync"
	"testing"
	"time"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

type TestStruct struct {
	Name  string `json:"name"`
	Value int    `json:"value"`
}

func setupTest(t *testing.T) func() {
	// Reset client and once for clean state
	memCacheClient = nil
	onceMemcache = sync.Once{}

	// Load env file
	if err := godotenv.Load("../.env"); err != nil {
		// If .env file doesn't exist, set default values
		os.Setenv("MEMCACHE_HOST", "localhost")
		os.Setenv("MEMCACHE_PORT", "11211")
	}

	// Initialize memcache
	if err := InitMemcache(); err != nil {
		t.Fatalf("Failed to initialize memcache: %v", err)
	}

	// Return cleanup function
	return func() {
		if memCacheClient != nil {
			if err := FlushAll(); err != nil {
				t.Logf("Warning: Failed to flush cache during cleanup: %v", err)
			}
		}
		memCacheClient = nil
		onceMemcache = sync.Once{}
	}
}

func TestInitMemcache(t *testing.T) {
	tests := []struct {
		name    string
		setup   func()
		wantErr bool
	}{
		{
			name: "successful initialization",
			setup: func() {
				os.Setenv("MEMCACHE_HOST", "localhost")
				os.Setenv("MEMCACHE_PORT", "11211")
			},
			wantErr: false,
		},
		{
			name: "missing host",
			setup: func() {
				os.Unsetenv("MEMCACHE_HOST")
				os.Setenv("MEMCACHE_PORT", "11211")
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset client for each test
			memCacheClient = nil
			onceMemcache = sync.Once{}

			tt.setup()
			err := InitMemcache()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, memCacheClient)
			}
		})
	}
}

func TestSetAndGetDataInCache(t *testing.T) {
	cleanup := setupTest(t)
	defer cleanup()

	assert.NotNil(t, memCacheClient, "memCacheClient should be initialized")

	tests := []struct {
		name    string
		key     string
		value   interface{}
		ttl     int32
		wantErr bool
	}{
		{
			name:    "string value",
			key:     "test_string",
			value:   "test value",
			ttl:     60,
			wantErr: false,
		},
		{
			name: "struct value",
			key:  "test_struct",
			value: TestStruct{
				Name:  "test",
				Value: 123,
			},
			ttl:     60,
			wantErr: false,
		},
		{
			name:    "empty key",
			key:     "",
			value:   "test",
			ttl:     60,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test Set
			err := SetDataInCache(tt.key, tt.value, tt.ttl)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)

			// Test Get
			switch v := tt.value.(type) {
			case string:
				result, err := GetDataFromCache[string](tt.key)
				assert.NoError(t, err)
				assert.Equal(t, v, result)
			case TestStruct:
				result, err := GetDataFromCache[TestStruct](tt.key)
				assert.NoError(t, err)
				assert.Equal(t, v, result)
			}
		})
	}
}

func TestUpdateDataInCache(t *testing.T) {
	cleanup := setupTest(t)
	defer cleanup()

	assert.NotNil(t, memCacheClient, "memCacheClient should be initialized")

	key := "test_update"
	initialValue := TestStruct{Name: "initial", Value: 1}
	updatedValue := TestStruct{Name: "updated", Value: 2}

	// Set initial value
	err := SetDataInCache(key, initialValue, 60)
	assert.NoError(t, err)

	// Update value
	err = UpdateDataInCache(key, updatedValue, 60)
	assert.NoError(t, err)

	// Verify update
	result, err := GetDataFromCache[TestStruct](key)
	assert.NoError(t, err)
	assert.Equal(t, updatedValue, result)

	// Try to update non-existent key
	err = UpdateDataInCache("non_existent", updatedValue, 60)
	assert.Error(t, err)
}

func TestDeleteDataFromCache(t *testing.T) {
	cleanup := setupTest(t)
	defer cleanup()

	assert.NotNil(t, memCacheClient, "memCacheClient should be initialized")

	key := "test_delete"
	value := "test value"

	// Set data
	err := SetDataInCache(key, value, 60)
	assert.NoError(t, err)

	// Delete data
	err = DeleteDataFromCache(key)
	assert.NoError(t, err)

	// Verify deletion
	_, err = GetDataFromCache[string](key)
	assert.Error(t, err)

	// Try to delete non-existent key
	err = DeleteDataFromCache(key)
	assert.NoError(t, err) // Should not return error for non-existent key
}

func TestCacheExpiration(t *testing.T) {
	cleanup := setupTest(t)
	defer cleanup()

	assert.NotNil(t, memCacheClient, "memCacheClient should be initialized")

	key := "test_expiration"
	value := "test value"
	ttl := int32(1) // 1 second TTL

	err := SetDataInCache(key, value, ttl)
	assert.NoError(t, err)

	// Wait for expiration
	time.Sleep(2 * time.Second)

	// Try to get expired data
	_, err = GetDataFromCache[string](key)
	assert.Error(t, err)
}

func TestFlushAll(t *testing.T) {
	cleanup := setupTest(t)
	defer cleanup()

	assert.NotNil(t, memCacheClient, "memCacheClient should be initialized")

	// Set multiple test data
	testData := map[string]string{
		"key1": "value1",
		"key2": "value2",
		"key3": "value3",
	}

	for k, v := range testData {
		err := SetDataInCache(k, v, 60)
		assert.NoError(t, err)
	}

	// Flush all data
	err := FlushAll()
	assert.NoError(t, err)

	// Verify all data is removed
	for k := range testData {
		_, err := GetDataFromCache[string](k)
		assert.Error(t, err)
	}
}
