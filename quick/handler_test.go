package quick_test

import (
	"context"
	"net/http"

	"github.com/deployport/airtls/quick"
)

// ExampleServeHTTPS demonstrates how to set up an HTTPS server using a self-signed certificate generator and in memory caching
func ExampleServeHTTPS() {
	ctx := context.Background()
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, HTTPS world with in-memory self-signed certificates!"))
	})

	err := quick.ServeHTTPS(
		ctx,
		":443",
		handler,
	)
	if err != nil {
		panic(err) // Handle error appropriately in production code
	}
}
