package quick

import (
	"context"
	"net/http"

	https "github.com/deployport/airtls/http"
	"github.com/deployport/airtls/selfsigned"
	"github.com/deployport/airtls/store"
)

// ServeHTTPS starts an HTTPS server that uses an auto-signed certificate generator
// and the provided HTTP handler. It uses a memory store for certificate storage.
// The server will listen on the specified address and handle requests using the provided handler.
// It will continue to serve until the provided context is cancelled.
// It returns an error if the server fails to start or if there are issues with certificate generation.
// This function is a convenience wrapper for quick setup without needing to manage certificate generation and storage manually.
// It is suitable for development and testing purposes.
// Note: For production use, consider using a more robust certificate management solution, unless the purpose is to use self-signed certificates
// behind a global proxy like Cloudflare with Flex origins, private proxy or in a controlled environment.
func ServeHTTPS(
	ctx context.Context,
	laddr string,
	handler http.Handler,
) error {
	return https.ServeHTTPS(
		ctx,
		selfsigned.NewGenerator(),
		store.NewMemoryStore(),
		laddr,
		handler,
	)
}
