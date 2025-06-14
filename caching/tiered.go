package caching

import (
	"crypto/tls"

	"github.com/deployport/airtls/store"
)

// TieredStore implements a Store that uses multiple Store implementations in order for tiered caching
type TieredStore struct {
	stores []store.Store
}

// NewTieredStore creates a new TieredStore with the given stores in order of priority, first to last where first is the highest priority
func NewTieredStore(stores ...store.Store) *TieredStore {
	return &TieredStore{stores: stores}
}

// GetCertificate tries to retrieve a certificate from each store in order, returning the first found.
// If none are found, returns a *store.CertificateNotFoundError.
func (t *TieredStore) GetCertificate(serverName string) (*tls.Certificate, error) {
	var lastErr error
	for _, s := range t.stores {
		cert, err := s.GetCertificate(serverName)
		if err == nil {
			return cert, nil
		}
		if !store.IsCertificateNotFound(err) {
			return nil, err
		}
		lastErr = err
	}
	if lastErr == nil {
		lastErr = store.NewCertificateNotFoundError()
	}
	return nil, lastErr
}

// SetCertificate sets the certificate in all stores in order. Returns the first error encountered, if any.
func (t *TieredStore) SetCertificate(serverName string, cert tls.Certificate) error {
	var firstErr error
	for _, s := range t.stores {
		err := s.SetCertificate(serverName, cert)
		if err != nil && firstErr == nil {
			firstErr = err
		}
	}
	return firstErr
}
