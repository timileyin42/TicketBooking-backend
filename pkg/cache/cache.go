package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"eventix-api/pkg/config"
	"eventix-api/pkg/logger"
	"go.uber.org/zap"
)

// Client holds the Redis client
var Client *redis.Client

// Connect establishes a connection to Redis
func Connect(cfg *config.RedisConfig) error {
	Client = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Test connection
	if err := Client.Ping(ctx).Err(); err != nil {
		return fmt.Errorf("failed to connect to Redis: %w", err)
	}

	logger.Info("Redis connected successfully",
		zap.String("host", cfg.Host),
		zap.Int("port", cfg.Port),
	)

	return nil
}

// Close closes the Redis connection
func Close() error {
	if Client == nil {
		return nil
	}
	return Client.Close()
}

// Get returns the Redis client
func Get() *redis.Client {
	return Client
}

// Set sets a key-value pair with expiration
func Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal value: %w", err)
	}

	return Client.Set(ctx, key, data, expiration).Err()
}

// GetValue gets a value by key and unmarshals it
func GetValue(ctx context.Context, key string, dest interface{}) error {
	data, err := Client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return fmt.Errorf("key not found")
		}
		return fmt.Errorf("failed to get value: %w", err)
	}

	if err := json.Unmarshal([]byte(data), dest); err != nil {
		return fmt.Errorf("failed to unmarshal value: %w", err)
	}

	return nil
}

// Delete deletes a key
func Delete(ctx context.Context, keys ...string) error {
	return Client.Del(ctx, keys...).Err()
}

// Exists checks if a key exists
func Exists(ctx context.Context, key string) (bool, error) {
	result, err := Client.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}
	return result > 0, nil
}

// SetNX sets a key-value pair only if it doesn't exist
func SetNX(ctx context.Context, key string, value interface{}, expiration time.Duration) (bool, error) {
	data, err := json.Marshal(value)
	if err != nil {
		return false, fmt.Errorf("failed to marshal value: %w", err)
	}

	return Client.SetNX(ctx, key, data, expiration).Result()
}

// Increment increments a key's value
func Increment(ctx context.Context, key string) (int64, error) {
	return Client.Incr(ctx, key).Result()
}

// Decrement decrements a key's value
func Decrement(ctx context.Context, key string) (int64, error) {
	return Client.Decr(ctx, key).Result()
}

// Expire sets expiration on a key
func Expire(ctx context.Context, key string, expiration time.Duration) error {
	return Client.Expire(ctx, key, expiration).Err()
}

// TTL gets the time to live for a key
func TTL(ctx context.Context, key string) (time.Duration, error) {
	return Client.TTL(ctx, key).Result()
}

// Health checks the Redis health
func Health() error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	return Client.Ping(ctx).Err()
}

// FlushDB flushes all keys in the current database (use with caution!)
func FlushDB(ctx context.Context) error {
	return Client.FlushDB(ctx).Err()
}

// Keys gets all keys matching a pattern
func Keys(ctx context.Context, pattern string) ([]string, error) {
	return Client.Keys(ctx, pattern).Result()
}

// ZAdd adds a member to a sorted set
func ZAdd(ctx context.Context, key string, members ...redis.Z) error {
	return Client.ZAdd(ctx, key, members...).Err()
}

// ZRange gets a range from a sorted set
func ZRange(ctx context.Context, key string, start, stop int64) ([]string, error) {
	return Client.ZRange(ctx, key, start, stop).Result()
}

// ZRem removes members from a sorted set
func ZRem(ctx context.Context, key string, members ...interface{}) error {
	return Client.ZRem(ctx, key, members...).Err()
}

// HSet sets a hash field
func HSet(ctx context.Context, key string, values ...interface{}) error {
	return Client.HSet(ctx, key, values...).Err()
}

// HGet gets a hash field
func HGet(ctx context.Context, key, field string) (string, error) {
	return Client.HGet(ctx, key, field).Result()
}

// HGetAll gets all fields from a hash
func HGetAll(ctx context.Context, key string) (map[string]string, error) {
	return Client.HGetAll(ctx, key).Result()
}

// HDel deletes hash fields
func HDel(ctx context.Context, key string, fields ...string) error {
	return Client.HDel(ctx, key, fields...).Err()
}
