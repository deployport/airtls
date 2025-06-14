package store

import (
	"crypto/tls"
	"errors"
)

// CertificateGetter defines an interface for retrieving TLS certificates by server name.
//
// The GetCertificate method attempts to find and return a tls.Certificate for the given serverName.
// If a certificate cannot be found or an error occurs during retrieval, it returns a non-nil error.
//
// If the certificate is not found, it returns a *CertificateNotFoundError.
type CertificateGetter interface {
	GetCertificate(serverName string) (*tls.Certificate, error)
}

// CertificateSetter implements storing TLS certificates
// by host name
type CertificateSetter interface {
	SetCertificate(serverName string, cert tls.Certificate) error
}

// Store implements retrieving and storing certificates
type Store interface {
	CertificateGetter
	CertificateSetter
}

// CertificateNotFoundError is returned when a certificate is not found in the store.
type CertificateNotFoundError struct{}

// NewCertificateNotFoundError creates a new instance of CertificateNotFoundError
func NewCertificateNotFoundError() *CertificateNotFoundError {
	return &CertificateNotFoundError{}
}

func (e *CertificateNotFoundError) Error() string {
	return "certificate not found"
}

// IsCertificateNotFound checks if the error is a CertificateNotFoundError.
func IsCertificateNotFound(err error) bool {
	var target *CertificateNotFoundError
	return errors.As(err, &target)
}

// Is implements errors.Is for CertificateNotFoundError.
func (e *CertificateNotFoundError) Is(target error) bool {
	_, ok := target.(*CertificateNotFoundError)
	return ok
}
