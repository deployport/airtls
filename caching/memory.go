package caching

import (
	"crypto/tls"
	"sync"

	"github.com/deployport/airtls/store"
)

// MemoryStore is a concurrent in-memory implementation of Store
// that stores certificates by server name.
type MemoryStore struct {
	mu    sync.RWMutex
	certs map[string]*tls.Certificate
}

// DefaultMemoryStoreCapacity is the default capacity for the MemoryStore.
var DefaultMemoryStoreCapacity = 100

// MemoryStoreOption configures a MemoryStore.
type MemoryStoreOption func(*MemoryStoreConfig)

// MemoryStoreConfig holds configuration for MemoryStore.
type MemoryStoreConfig struct {
	Capacity int
}

// WithCapacity sets the capacity for the MemoryStore.
func WithCapacity(capacity int) MemoryStoreOption {
	return func(cfg *MemoryStoreConfig) {
		cfg.Capacity = capacity
	}
}

// NewMemoryStore creates a new MemoryStore instance with options.
func NewMemoryStore(opts ...MemoryStoreOption) *MemoryStore {
	cfg := &MemoryStoreConfig{
		Capacity: DefaultMemoryStoreCapacity,
	}
	for _, opt := range opts {
		opt(cfg)
	}
	if cfg.Capacity <= 0 {
		cfg.Capacity = DefaultMemoryStoreCapacity
	}
	return &MemoryStore{
		certs: make(map[string]*tls.Certificate, cfg.Capacity),
	}
}

// GetCertificate retrieves a certificate by server name.
func (m *MemoryStore) GetCertificate(serverName string) (*tls.Certificate, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	cert, ok := m.certs[serverName]
	if !ok {
		return nil, store.NewCertificateNotFoundError()
	}
	return cert, nil
}

// SetCertificate stores a certificate by server name.
func (m *MemoryStore) SetCertificate(serverName string, cert tls.Certificate) error {
	m.mu.Lock()
	m.certs[serverName] = &cert
	m.mu.Unlock()
	return nil
}
