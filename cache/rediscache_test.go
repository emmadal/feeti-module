package cache

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func setupTestRedis() {
	os.Setenv("REDIS_HOST", "localhost")
	os.Setenv("REDIS_PORT", "6379")
	_ = InitRedis()
}

type RedisTest struct {
	Name  string
	Value string
}

func TestInitRedis(t *testing.T) {
	setupTestRedis()

	assert.NotNil(t, rdb, "Redis client should be initialized")
	err := rdb.Ping(context.Background()).Err()
	assert.NoError(t, err, "Redis should be accessible")
}

func TestSetGetRedisData(t *testing.T) {
	setupTestRedis()
	ctx := context.Background()
	key := "test-key"
	value := RedisTest{Name: "test", Value: "test-value"}

	err := SetRedisData(ctx, key, value, 1)
	assert.NoError(t, err, "SetRedisData should not return an error")

	var result RedisTest
	result, err = GetRedisData[RedisTest](ctx, key)
	assert.NoError(t, err, "GetRedisData should not return an error")
	assert.Equal(t, value, result, "Stored and retrieved values should match")
}

func TestUpdateRedisData(t *testing.T) {
	setupTestRedis()
	ctx := context.Background()
	key := "test-key"
	initialValue := RedisTest{Name: "initial", Value: "initial"}
	updatedValue := RedisTest{Name: "updated", Value: "updated"}

	_ = SetRedisData(ctx, key, initialValue, 1)

	err := UpdateRedisData[RedisTest](ctx, key, updatedValue)
	assert.NoError(t, err, "UpdateRedisData should not return an error")

	var result RedisTest
	result, err = GetRedisData[RedisTest](ctx, key)
	assert.NoError(t, err, "GetRedisData should not return an error")
	assert.Equal(t, updatedValue, result, "Updated value should be retrievable")
}

func TestDeleteRedisData(t *testing.T) {
	setupTestRedis()
	ctx := context.Background()
	key := "test-key"

	_ = SetRedisData(ctx, key, RedisTest{Name: "some-value", Value: "some-value"}, 1)

	err := DeleteRedisData(ctx, key)
	assert.NoError(t, err, "DeleteRedisData should not return an error")

	_, err = GetRedisData[string](ctx, key)
	assert.Error(t, err, "Getting deleted key should return an error")
}

func TestFlushAllRedis(t *testing.T) {
	setupTestRedis()
	ctx := context.Background()

	_ = SetRedisData(ctx, "key1", "value1", 1)
	_ = SetRedisData(ctx, "key2", "value2", 1)

	err := FlushAllRedis(ctx)
	assert.NoError(t, err, "FlushAllRedis should not return an error")

	_, err1 := GetRedisData[string](ctx, "key1")
	_, err2 := GetRedisData[string](ctx, "key2")

	assert.Error(t, err1, "Flushed key1 should not be retrievable")
	assert.Error(t, err2, "Flushed key2 should not be retrievable")
}

func TestCloseRedis(t *testing.T) {
	setupTestRedis()
	err := CloseRedis()
	assert.NoError(t, err, "CloseRedis should not return an error")
}
