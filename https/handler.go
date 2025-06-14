package https

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"

	"github.com/deployport/airtls/store"
)

// ServeHTTPS starts an HTTPS server that uses the provided generator to create certificates
// and using the http package shared mux handler
func ServeHTTPS(
	ctx context.Context,
	generator store.Generator,
	store store.Store,
	laddr string,
	handler http.Handler,
) error {
	getter, err := NewGetCertificate(generator, store)
	if err != nil {
		return fmt.Errorf("failed to create get certificate: %w", err)
	}

	tlsConfig := &tls.Config{
		GetCertificate: getter,
	}

	ln, err := tls.Listen("tcp", laddr, tlsConfig)
	if err != nil {
		return fmt.Errorf("failed to listen for TLS: %w", err)
	}
	defer ln.Close()
	go func() {
		<-ctx.Done()
		ln.Close()
	}()
	srv := &http.Server{Handler: handler}
	return srv.Serve(ln)
}
