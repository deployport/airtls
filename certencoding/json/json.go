package json

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"io"
)

// Certificate represents a TLS certificate in JSON format
type Certificate struct {
	CertPEM string `json:"c,omitempty"` // PEM-encoded certificate chain
	// KeyPEM is the PEM-encoded private key.
	KeyPEM string `json:"k,omitempty"` // PEM-encoded private key
}

// MarshalTLSCert converts a tls.Certificate to a Certificate
func MarshalTLSCert(cert tls.Certificate) (Certificate, error) {
	// Combine the cert chain
	var certPEM bytes.Buffer
	for _, der := range cert.Certificate {
		err := pem.Encode(&certPEM, &pem.Block{
			Type:  "CERTIFICATE",
			Bytes: der,
		})
		if err != nil {
			return Certificate{}, err
		}
	}

	// Encode the private key
	privKeyBytes, err := x509.MarshalPKCS8PrivateKey(cert.PrivateKey)
	if err != nil {
		return Certificate{}, err
	}
	keyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: privKeyBytes,
	})

	return Certificate{
		CertPEM: certPEM.String(),
		KeyPEM:  string(keyPEM),
	}, nil
}

// UnmarshalTLSCert converts a Certificate to a tls.Certificate
func UnmarshalTLSCert(jsonCert Certificate) (*tls.Certificate, error) {
	cert, err := tls.X509KeyPair([]byte(jsonCert.CertPEM), []byte(jsonCert.KeyPEM))
	if err != nil {
		return nil, err
	}
	return &cert, nil
}

// Marshaler implements Marshaler and Unmarshaler for JSON encoding.
type Marshaler struct{}

// Marshal serializes a tls.Certificate to compact JSON and writes to w.
func (j *Marshaler) Marshal(cert tls.Certificate, w io.Writer) error {
	jsonCert, err := MarshalTLSCert(cert)
	if err != nil {
		return err
	}
	enc := json.NewEncoder(w)
	enc.SetIndent("", "") // ensure compact (no indent)
	return enc.Encode(jsonCert)
}

// Unmarshal deserializes a tls.Certificate from compact JSON read from r.
func (j *Marshaler) Unmarshal(r io.Reader) (*tls.Certificate, error) {
	var jsonCert Certificate
	dec := json.NewDecoder(r)
	if err := dec.Decode(&jsonCert); err != nil {
		return nil, err
	}
	return UnmarshalTLSCert(jsonCert)
}
