package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// Client wraps go-redis for typed JSON cache operations.
type Client struct {
	rdb *redis.Client
}

// Connect parses REDIS_URL and verifies connectivity with PING.
func Connect(ctx context.Context, redisURL string) (*Client, error) {
	opt, err := redis.ParseURL(redisURL)
	if err != nil {
		return nil, fmt.Errorf("parse redis url: %w", err)
	}
	rdb := redis.NewClient(opt)
	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("redis ping: %w", err)
	}
	return &Client{rdb: rdb}, nil
}

// Close shuts down the Redis connection pool.
func (c *Client) Close() error {
	return c.rdb.Close()
}

// SetJSON marshals v and stores it with an optional TTL.
func (c *Client) SetJSON(ctx context.Context, key string, v interface{}, ttl time.Duration) error {
	data, err := json.Marshal(v)
	if err != nil {
		return err
	}
	return c.rdb.Set(ctx, key, data, ttl).Err()
}

// GetJSON loads and unmarshals a cached JSON value. Returns false if the key is missing.
func (c *Client) GetJSON(ctx context.Context, key string, dest interface{}) (bool, error) {
	data, err := c.rdb.Get(ctx, key).Bytes()
	if err == redis.Nil {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	if err := json.Unmarshal(data, dest); err != nil {
		return false, err
	}
	return true, nil
}

// RDB exposes the underlying client for advanced use (e.g. INCR rate limits).
func (c *Client) RDB() *redis.Client {
	return c.rdb
}
