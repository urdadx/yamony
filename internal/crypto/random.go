package crypto

import (
	"crypto/rand"
	"fmt"
	"io"
)

// GenerateRandomBytes generates cryptographically secure random bytes
func GenerateRandomBytes(n int) ([]byte, error) {
	if n <= 0 {
		return nil, fmt.Errorf("length must be positive")
	}
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		return nil, fmt.Errorf("failed to generate random bytes: %w", err)
	}
	return b, nil
}

// GenerateNonce generates a nonce of specified size
func GenerateNonce(size int) ([]byte, error) {
	return GenerateRandomBytes(size)
}

// SecureRandomReader returns a reader for cryptographically secure random data
func SecureRandomReader() io.Reader {
	return rand.Reader
}

// GenerateID generates a cryptographically secure random ID as hex string
func GenerateID(byteLength int) (string, error) {
	bytes, err := GenerateRandomBytes(byteLength)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", bytes), nil
}
