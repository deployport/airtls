package caching_test

import (
	"crypto/tls"
	"sync"
	"testing"

	"github.com/deployport/airtls/caching"
	"github.com/deployport/airtls/selfsigned"
	"github.com/deployport/airtls/store"
)

func TestTieredStore(t *testing.T) {
	t.Run("GetCertificate", func(t *testing.T) {
		t.Run("GetCertificate first hit", func(t *testing.T) {
			store1 := NewMockStore()
			store2 := NewMockStore()
			tieredStore := caching.NewTieredStore(store1, store2)

			selfSignedGenerator := selfsigned.NewGenerator()

			cert, err := selfSignedGenerator.Generate("example.com")
			if err != nil {
				t.Fatalf("failed to generate self-signed certificate: %v", err)
			}

			// Test SetCertificate
			if err := store1.SetCertificate("example.com", *cert); err != nil {
				t.Fatalf("SetCertificate failed: %v", err)
			}

			// Test GetCertificate from first store
			retrievedCert, err := tieredStore.GetCertificate("example.com")
			if err != nil {
				t.Fatalf("GetCertificate failed: %v", err)
			}
			if retrievedCert == nil {
				t.Fatal("Expected certificate to be retrieved, got nil")
			}

			// Check calls to stores
			if len(store1.GetCalls) != 1 {
				t.Errorf("Expected GetCertificate to be called once on first store, got %d", len(store1.SetCalls))
			}
			if len(store2.GetCalls) != 0 {
				t.Errorf("Expected GetCertificate to not be called once on second store, got %d", len(store2.GetCalls))
			}
		})
		t.Run("GetCertificate miss", func(t *testing.T) {
			store1 := NewMockStore()
			store2 := NewMockStore()
			tieredStore := caching.NewTieredStore(store1, store2)

			// Test GetCertificate from first store
			retrievedCert, err := tieredStore.GetCertificate("example.com")

			if err == nil {
				t.Fatal("Expected error for missing certificate, got non-nil certificate")
			}
			if retrievedCert != nil {
				t.Fatal("Expected nil certificate for missing entry, got non-nil certificate")
			}
			// Check calls to stores
			if len(store1.GetCalls) != 1 {
				t.Errorf("Expected GetCertificate to be called once on first store, got %d", len(store1.SetCalls))
			}
			if len(store2.GetCalls) != 1 {
				t.Errorf("Expected GetCertificate to be called once on second store, got %d", len(store2.GetCalls))
			}

		})
	})
	t.Run("SetCertificate", func(t *testing.T) {
		t.Run("SetCertificate first hit", func(t *testing.T) {
			store1 := NewMockStore()
			store2 := NewMockStore()
			tieredStore := caching.NewTieredStore(store1, store2)

			selfSignedGenerator := selfsigned.NewGenerator()

			cert, err := selfSignedGenerator.Generate("example.com")
			if err != nil {
				t.Fatalf("failed to generate self-signed certificate: %v", err)
			}

			// Test SetCertificate
			if err := tieredStore.SetCertificate("example.com", *cert); err != nil {
				t.Fatalf("SetCertificate failed: %v", err)
			}

			// Check calls to stores
			if len(store1.SetCalls) != 1 {
				t.Errorf("Expected SetCertificate to be called once on first store, got %d", len(store1.SetCalls))
			}
			if len(store2.SetCalls) != 1 {
				t.Errorf("Expected SetCertificate to be called once on second store, got %d", len(store2.SetCalls))
			}
		})
	})
}

type MockStore struct {
	Certs    map[string]*tls.Certificate
	SetCalls []string
	GetCalls []string
	SetErr   error
	GetErr   error
	mu       sync.Mutex
}

func NewMockStore() *MockStore {
	return &MockStore{
		Certs: make(map[string]*tls.Certificate),
	}
}

func (m *MockStore) GetCertificate(serverName string) (*tls.Certificate, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.GetCalls = append(m.GetCalls, serverName)
	if m.GetErr != nil {
		return nil, m.GetErr
	}
	cert, ok := m.Certs[serverName]
	if !ok {
		return nil, store.NewCertificateNotFoundError()
	}
	return cert, nil
}

func (m *MockStore) SetCertificate(serverName string, cert tls.Certificate) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.SetCalls = append(m.SetCalls, serverName)
	if m.SetErr != nil {
		return m.SetErr
	}
	m.Certs[serverName] = &cert
	return nil
}
