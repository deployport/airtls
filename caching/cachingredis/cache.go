package cachingredis

import (
	"bytes"
	"context"
	"crypto/tls"
	"fmt"

	"github.com/deployport/airtls/certencoding/json"
	"github.com/deployport/airtls/store"
	redis "github.com/redis/go-redis/v9"
)

// RedisCache implements a certificate cache using Redis as the backend.
type RedisCache struct {
	client *redis.Client
	prefix string
}

// RedisCacheOption configures a RedisCache.
type RedisCacheOption func(*RedisCache)

// WithPrefix sets the prefix for cache entries in Redis.
func WithPrefix(prefix string) RedisCacheOption {
	return func(c *RedisCache) {
		c.prefix = prefix
	}
}

// NewRedisCache creates a new RedisCache with the given Redis client and options.
func NewRedisCache(client *redis.Client, opts ...RedisCacheOption) *RedisCache {
	cache := &RedisCache{
		client: client,
		prefix: "airtls:",
	}
	for _, opt := range opts {
		opt(cache)
	}
	return cache
}

func (c *RedisCache) key(serverName string) string {
	return c.prefix + serverName
}

// GetCertificate retrieves a certificate by server name from Redis using JSON marshaling.
func (c *RedisCache) GetCertificate(serverName string) (*tls.Certificate, error) {
	ctx := context.Background()
	val, err := c.client.Get(ctx, c.key(serverName)).Bytes()
	if err == redis.Nil {
		return nil, store.NewCertificateNotFoundError()
	}
	if err != nil {
		return nil, fmt.Errorf("redis get error: %w", err)
	}
	marshaler := json.Marshaler{}
	cert, err := marshaler.Unmarshal(bytes.NewReader(val))
	if err != nil {
		return nil, fmt.Errorf("json unmarshal error: %w", err)
	}
	return cert, nil
}

// SetCertificate stores a certificate by server name in Redis using JSON marshaling.
func (c *RedisCache) SetCertificate(serverName string, cert tls.Certificate) error {
	ctx := context.Background()
	var buf bytes.Buffer
	marshaler := json.Marshaler{}
	if err := marshaler.Marshal(cert, &buf); err != nil {
		return fmt.Errorf("json marshal error: %w", err)
	}
	if err := c.client.Set(ctx, c.key(serverName), buf.Bytes(), 0).Err(); err != nil {
		return fmt.Errorf("redis set error: %w", err)
	}
	return nil
}
