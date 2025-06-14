package https

import (
	"crypto/tls"
	"fmt"

	certstore "github.com/deployport/airtls/store"
)

// GetCertificateFunc is a function type that retrieves or generates a TLS certificate for a given host
type GetCertificateFunc func(chi *tls.ClientHelloInfo) (*tls.Certificate, error)

// NewGetCertificate returns a function that retrieves or generates a TLS certificate for a given host
// using the provided generator and store. If the certificate is not found in the store, it generates a new one.
// you can use this function as the GetCertificate callback in a tls.Config.
func NewGetCertificate(
	generator certstore.Generator,
	store certstore.Store,
) (GetCertificateFunc, error) {
	if generator == nil {
		return nil, fmt.Errorf("generator is nil")
	}
	if store == nil {
		return nil, fmt.Errorf("store is nil")
	}
	return GetCertificateFunc(func(chi *tls.ClientHelloInfo) (*tls.Certificate, error) {
		host := chi.ServerName
		if host == "" {
			host = "localhost"
		}
		cert, err := store.GetCertificate(host)
		if certstore.IsCertificateNotFound(err) {
			cert, err = generator.Generate(host)
			if err != nil {
				return nil, fmt.Errorf("failed to generate certificate for %s: %w", host, err)
			}
		} else if err != nil {
			return nil, fmt.Errorf("failed to get certificate for %s: %w", host, err)
		}
		if err := store.SetCertificate(host, *cert); err != nil {
			return nil, fmt.Errorf("failed to store certificate for %s: %w", host, err)
		}
		return cert, nil
	}), nil
}
