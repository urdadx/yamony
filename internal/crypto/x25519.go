package crypto

import (
	"crypto/rand"
	"fmt"

	"golang.org/x/crypto/curve25519"
)

const (
	// X25519KeySize is the size of X25519 keys in bytes
	X25519KeySize = 32
)

// X25519KeyPair represents an X25519 key pair for ECDH
type X25519KeyPair struct {
	PublicKey  []byte
	PrivateKey []byte
}

// GenerateX25519KeyPair generates a new X25519 key pair for ECDH
func GenerateX25519KeyPair() (*X25519KeyPair, error) {
	privateKey := make([]byte, X25519KeySize)
	if _, err := rand.Read(privateKey); err != nil {
		return nil, fmt.Errorf("failed to generate X25519 private key: %w", err)
	}

	publicKey, err := curve25519.X25519(privateKey, curve25519.Basepoint)
	if err != nil {
		return nil, fmt.Errorf("failed to derive X25519 public key: %w", err)
	}

	return &X25519KeyPair{
		PublicKey:  publicKey,
		PrivateKey: privateKey,
	}, nil
}

// X25519SharedSecret computes a shared secret using ECDH
func X25519SharedSecret(privateKey, publicKey []byte) ([]byte, error) {
	if len(privateKey) != X25519KeySize {
		return nil, fmt.Errorf("invalid private key length: expected %d, got %d",
			X25519KeySize, len(privateKey))
	}
	if len(publicKey) != X25519KeySize {
		return nil, fmt.Errorf("invalid public key length: expected %d, got %d",
			X25519KeySize, len(publicKey))
	}

	sharedSecret, err := curve25519.X25519(privateKey, publicKey)
	if err != nil {
		return nil, fmt.Errorf("X25519 key exchange failed: %w", err)
	}

	return sharedSecret, nil
}

// DeriveSharedKey derives a symmetric key from X25519 ECDH shared secret using HKDF
// This is the recommended way to use X25519 shared secrets for encryption
func DeriveSharedKey(privateKey, publicKey []byte, info string) ([]byte, error) {
	sharedSecret, err := X25519SharedSecret(privateKey, publicKey)
	if err != nil {
		return nil, err
	}

	// Derive a proper symmetric key from the shared secret using HKDF
	symmetricKey, err := HKDFExtractAndExpand(sharedSecret, nil, []byte(info), 32)
	if err != nil {
		return nil, fmt.Errorf("failed to derive shared key: %w", err)
	}

	return symmetricKey, nil
}

// DeriveSharedKeyForSharing derives a key for sharing vault items
func DeriveSharedKeyForSharing(ourPrivateKey, theirPublicKey []byte, vaultID string) ([]byte, error) {
	info := fmt.Sprintf("share-vek:%s", vaultID)
	return DeriveSharedKey(ourPrivateKey, theirPublicKey, info)
}

// ValidateX25519PublicKey checks if a byte slice is a valid X25519 public key
func ValidateX25519PublicKey(pubKey []byte) error {
	if len(pubKey) != X25519KeySize {
		return fmt.Errorf("invalid X25519 public key length: expected %d, got %d",
			X25519KeySize, len(pubKey))
	}
	return nil
}

// ValidateX25519PrivateKey checks if a byte slice is a valid X25519 private key
func ValidateX25519PrivateKey(privKey []byte) error {
	if len(privKey) != X25519KeySize {
		return fmt.Errorf("invalid X25519 private key length: expected %d, got %d",
			X25519KeySize, len(privKey))
	}
	return nil
}
