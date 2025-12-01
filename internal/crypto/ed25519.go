package crypto

import (
	"crypto/ed25519"
	"crypto/rand"
	"fmt"
)

// Ed25519KeyPair represents an Ed25519 signing key pair
type Ed25519KeyPair struct {
	PublicKey  ed25519.PublicKey
	PrivateKey ed25519.PrivateKey
}

// GenerateEd25519KeyPair generates a new Ed25519 key pair for signatures
func GenerateEd25519KeyPair() (*Ed25519KeyPair, error) {
	pub, priv, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return nil, fmt.Errorf("failed to generate Ed25519 key pair: %w", err)
	}
	return &Ed25519KeyPair{
		PublicKey:  pub,
		PrivateKey: priv,
	}, nil
}

// Sign creates a signature for the given message
func (kp *Ed25519KeyPair) Sign(message []byte) []byte {
	return ed25519.Sign(kp.PrivateKey, message)
}

// SignMessage is a convenience function to sign a message with a private key
func SignMessage(privateKey ed25519.PrivateKey, message []byte) []byte {
	return ed25519.Sign(privateKey, message)
}

// VerifySignature verifies an Ed25519 signature
func VerifySignature(publicKey ed25519.PublicKey, message, signature []byte) bool {
	return ed25519.Verify(publicKey, message, signature)
}

// ValidateEd25519PublicKey checks if a byte slice is a valid Ed25519 public key
func ValidateEd25519PublicKey(pubKey []byte) error {
	if len(pubKey) != ed25519.PublicKeySize {
		return fmt.Errorf("invalid Ed25519 public key length: expected %d, got %d",
			ed25519.PublicKeySize, len(pubKey))
	}
	return nil
}

// ValidateEd25519PrivateKey checks if a byte slice is a valid Ed25519 private key
func ValidateEd25519PrivateKey(privKey []byte) error {
	if len(privKey) != ed25519.PrivateKeySize {
		return fmt.Errorf("invalid Ed25519 private key length: expected %d, got %d",
			ed25519.PrivateKeySize, len(privKey))
	}
	return nil
}

// ValidateEd25519Signature checks if a byte slice is a valid Ed25519 signature
func ValidateEd25519Signature(sig []byte) error {
	if len(sig) != ed25519.SignatureSize {
		return fmt.Errorf("invalid Ed25519 signature length: expected %d, got %d",
			ed25519.SignatureSize, len(sig))
	}
	return nil
}
