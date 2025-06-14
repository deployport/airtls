package cachingredis

import (
	"crypto/tls"
	"fmt"
	"log"

	redis "github.com/redis/go-redis/v9"
)

// ExampleRedisCacheUsage demonstrates how to use RedisCache with JSON encoding for certificate storage.
func ExampleRedisCacheUsage() {
	// Create a Redis client (adjust options as needed)
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		DB:   0,
	})

	// Create the RedisCache with a custom prefix (optional)
	cache := NewRedisCache(client, WithPrefix("mycerts:"))

	// Example certificate (replace with a real one in production)
	cert := tls.Certificate{
		Certificate: [][]byte{[]byte("dummy-cert")},
		PrivateKey:  nil, // For demonstration only
	}

	serverName := "example.com"

	// Store the certificate
	err := cache.SetCertificate(serverName, cert)
	if err != nil {
		log.Fatalf("failed to set certificate: %v", err)
	}
	fmt.Println("Certificate stored in Redis.")

	// Retrieve the certificate
	retrievedCert, err := cache.GetCertificate(serverName)
	if err != nil {
		log.Fatalf("failed to get certificate: %v", err)
	}
	fmt.Printf("Retrieved certificate for %s: %+v\n", serverName, retrievedCert)
}
