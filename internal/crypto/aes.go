package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
)

const (
	// GCMNonceSize is the standard nonce size for GCM (96 bits / 12 bytes)
	GCMNonceSize = 12
	// GCMTagSize is the authentication tag size for GCM (128 bits / 16 bytes)
	GCMTagSize = 16
)

// EncryptedData holds the components of an AES-GCM encrypted message
type EncryptedData struct {
	Ciphertext []byte `json:"ciphertext"`
	IV         []byte `json:"iv"`
	Tag        []byte `json:"tag"`
}

// EncryptAESGCM encrypts plaintext using AES-256-GCM
// The key must be 32 bytes (256 bits)
// AAD (Additional Authenticated Data) is optional and can be nil
func EncryptAESGCM(key, plaintext, aad []byte) (*EncryptedData, error) {
	if len(key) != 32 {
		return nil, fmt.Errorf("key must be 32 bytes, got %d", len(key))
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("failed to create cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("failed to create GCM: %w", err)
	}

	// Generate a random nonce
	nonce := make([]byte, gcm.NonceSize())
	if _, err := rand.Read(nonce); err != nil {
		return nil, fmt.Errorf("failed to generate nonce: %w", err)
	}

	// Seal encrypts and authenticates plaintext
	// Format: ciphertext || tag (tag is appended by Seal)
	sealed := gcm.Seal(nil, nonce, plaintext, aad)

	// Split ciphertext and tag
	if len(sealed) < GCMTagSize {
		return nil, fmt.Errorf("sealed data too short")
	}

	ciphertextLen := len(sealed) - GCMTagSize
	ciphertext := sealed[:ciphertextLen]
	tag := sealed[ciphertextLen:]

	return &EncryptedData{
		Ciphertext: ciphertext,
		IV:         nonce,
		Tag:        tag,
	}, nil
}

// DecryptAESGCM decrypts ciphertext using AES-256-GCM
// The key must be 32 bytes (256 bits)
// AAD must match what was used during encryption
func DecryptAESGCM(key, ciphertext, iv, tag, aad []byte) ([]byte, error) {
	if len(key) != 32 {
		return nil, fmt.Errorf("key must be 32 bytes, got %d", len(key))
	}

	if len(iv) != GCMNonceSize {
		return nil, fmt.Errorf("IV must be %d bytes, got %d", GCMNonceSize, len(iv))
	}

	if len(tag) != GCMTagSize {
		return nil, fmt.Errorf("tag must be %d bytes, got %d", GCMTagSize, len(tag))
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("failed to create cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("failed to create GCM: %w", err)
	}

	// Combine ciphertext and tag for Open
	sealed := append(ciphertext, tag...)

	// Open verifies and decrypts
	plaintext, err := gcm.Open(nil, iv, sealed, aad)
	if err != nil {
		return nil, fmt.Errorf("decryption failed: %w", err)
	}

	return plaintext, nil
}

// EncryptAESGCMCombined is a convenience function that returns the ciphertext with tag appended
// This matches the format that Go's Seal produces naturally
func EncryptAESGCMCombined(key, plaintext, aad []byte) (ciphertext, iv []byte, err error) {
	if len(key) != 32 {
		return nil, nil, fmt.Errorf("key must be 32 bytes, got %d", len(key))
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create GCM: %w", err)
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := rand.Read(nonce); err != nil {
		return nil, nil, fmt.Errorf("failed to generate nonce: %w", err)
	}

	sealed := gcm.Seal(nil, nonce, plaintext, aad)
	return sealed, nonce, nil
}

// DecryptAESGCMCombined decrypts data where ciphertext and tag are combined
func DecryptAESGCMCombined(key, combined, iv, aad []byte) ([]byte, error) {
	if len(key) != 32 {
		return nil, fmt.Errorf("key must be 32 bytes, got %d", len(key))
	}

	if len(iv) != GCMNonceSize {
		return nil, fmt.Errorf("IV must be %d bytes, got %d", GCMNonceSize, len(iv))
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("failed to create cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("failed to create GCM: %w", err)
	}

	plaintext, err := gcm.Open(nil, iv, combined, aad)
	if err != nil {
		return nil, fmt.Errorf("decryption failed: %w", err)
	}

	return plaintext, nil
}
