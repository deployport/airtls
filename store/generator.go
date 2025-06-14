package store

import "crypto/tls"

// Generator defines an interface for generating TLS certificates on the fly
type Generator interface {
	// Generate generates a new certificate for the given server name.
	Generate(serverName string) (*tls.Certificate, error)
}
