package cachingredis_test

import (
	"context"
	"net/http"

	"github.com/deployport/airtls/caching/cachingredis"
	"github.com/deployport/airtls/https"
	"github.com/deployport/airtls/selfsigned"
	redis "github.com/redis/go-redis/v9"
)

// ExampleRedisCache demonstrates how to use RedisCache with JSON encoding for certificate storage.
func ExampleRedisCache() {
	ctx := context.Background()
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		DB:   0,
	})

	// Create the RedisCache with a custom prefix (optional)
	cache := cachingredis.New(client, cachingredis.WithPrefix("mycerts:"))

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Handle the request, e.g., serve a simple response
		w.Write([]byte("Hello, HTTPS with Redis Cache!"))
	})

	err := https.ServeHTTPS(
		ctx,
		selfsigned.NewGenerator(),
		cache,
		":443",
		handler,
	)
	if err != nil {
		panic(err) // Handle error appropriately in production code
	}
}
