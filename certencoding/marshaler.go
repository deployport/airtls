package certencoding

import (
	"crypto/tls"
	"io"
)

// Marshaler defines an interface for serializing a tls.Certificate to an io.Writer.
type Marshaler interface {
	Marshal(cert tls.Certificate, w io.Writer) error
}

// Unmarshaler defines an interface for deserializing a tls.Certificate from an io.Reader.
type Unmarshaler interface {
	Unmarshal(r io.Reader) (*tls.Certificate, error)
}
