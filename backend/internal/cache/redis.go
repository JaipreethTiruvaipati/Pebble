// Package cache wraps Redis for Pebble: connection lifecycle, JSON get/set, and
// key helpers. Services share a single Client for market snapshots, portfolio
// summaries, OTP storage, and rate-limit counters accessed via RDB().
package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// Client is a thin wrapper around go-redis with JSON serialization helpers.
// Handlers and middleware receive *Client from service main; low-level commands
// (INCR, EXPIRE) use RDB() when helpers do not fit.
type Client struct {
	rdb *redis.Client
}

// Connect parses redisURL, opens a connection pool, and verifies reachability with PING.
//
// Inputs: ctx for dial timeout; redisURL from config.RedisURL (redis://host:port).
// Outputs: *Client on success; error if URL is invalid or Redis is unreachable.
//
// Called at startup in api-gateway, bill-service, and other services that cache or
// rate-limit; failure prevents boot so misconfigured REDIS_URL is caught early.
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

// Close drains and closes the underlying Redis connection pool.
//
// Inputs: none on receiver.
// Outputs: error from go-redis Close, if any.
//
// Registered with defer in service main for graceful shutdown alongside HTTP and DB.
func (c *Client) Close() error {
	return c.rdb.Close()
}

// SetJSON marshals v to JSON and stores it at key with optional TTL.
//
// Inputs: ctx; key from keys.go helpers; v any JSON-serializable value; ttl (0 = no expiry).
// Outputs: error from Marshal or Redis SET.
//
// Used for portfolio summaries, session blobs, and market poller writes; TTL aligns
// with constants like PortfolioSummaryTTL for cache freshness vs load on PostgreSQL.
func (c *Client) SetJSON(ctx context.Context, key string, v interface{}, ttl time.Duration) error {
	data, err := json.Marshal(v)
	if err != nil {
		return err
	}
	return c.rdb.Set(ctx, key, data, ttl).Err()
}

// GetJSON loads key and unmarshals JSON into dest.
//
// Inputs: ctx; key; dest pointer to populate.
// Outputs: (true, nil) on hit; (false, nil) if key missing (redis.Nil); error on Redis or JSON failure.
//
// Callers treat false as cache miss and fall back to database or upstream API, then
// optionally repopulate with SetJSON.
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

// RDB exposes the underlying go-redis client for operations without JSON helpers.
//
// Inputs: none.
// Outputs: *redis.Client for INCR, EXPIRE, pipelines, etc.
//
// Used by auth.RedisRateLimit and any service needing atomic counters or custom commands.
func (c *Client) RDB() *redis.Client {
	return c.rdb
}

// Delete removes a single key from Redis.
//
// Inputs: ctx; key (e.g. KeyPortfolioSummary after a trade invalidates cached digest).
// Outputs: error from DEL if Redis fails.
//
// Keeps portfolio and market caches coherent when underlying data changes in PostgreSQL.
func (c *Client) Delete(ctx context.Context, key string) error {
	return c.rdb.Del(ctx, key).Err()
}
