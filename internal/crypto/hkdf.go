package crypto

import (
	"crypto/sha256"
	"fmt"
	"io"

	"golang.org/x/crypto/hkdf"
)

// HKDFExtract performs the extract step of HKDF
func HKDFExtract(secret, salt []byte) []byte {
	if salt == nil {
		salt = make([]byte, sha256.Size)
	}
	mac := sha256.New
	h := mac()
	h.Write(salt)
	h.Write(secret)
	return h.Sum(nil)
}

// HKDFExpand performs the expand step of HKDF
func HKDFExpand(prk, info []byte, length int) ([]byte, error) {
	h := hkdf.Expand(sha256.New, prk, info)
	out := make([]byte, length)
	if _, err := io.ReadFull(h, out); err != nil {
		return nil, fmt.Errorf("HKDF expand failed: %w", err)
	}
	return out, nil
}

// HKDFExtractAndExpand performs both extract and expand steps in one call
// This is the most common usage pattern for HKDF
func HKDFExtractAndExpand(secret, salt, info []byte, outLen int) ([]byte, error) {
	if outLen <= 0 {
		return nil, fmt.Errorf("output length must be positive")
	}

	if outLen > 255*sha256.Size {
		return nil, fmt.Errorf("output length too large for HKDF-SHA256")
	}

	hk := hkdf.New(sha256.New, secret, salt, info)
	out := make([]byte, outLen)
	if _, err := io.ReadFull(hk, out); err != nil {
		return nil, fmt.Errorf("HKDF failed: %w", err)
	}
	return out, nil
}

// DeriveKey is a convenience function for common key derivation scenarios
// It derives a key of specified length from a master key and context info
func DeriveKey(masterKey []byte, context string, keyLen int) ([]byte, error) {
	if keyLen <= 0 {
		keyLen = 32 // default to 256 bits
	}
	return HKDFExtractAndExpand(masterKey, nil, []byte(context), keyLen)
}

// DeriveVaultEncryptionKey derives a Vault Encryption Key (VEK) from a master key
func DeriveVaultEncryptionKey(masterKey []byte, vaultID string) ([]byte, error) {
	info := []byte(fmt.Sprintf("vault-key:%s", vaultID))
	return HKDFExtractAndExpand(masterKey, nil, info, 32)
}

// DeriveItemEncryptionKey derives an Item Encryption Key (IEK) from a VEK and item ID
func DeriveItemEncryptionKey(vek []byte, itemID string) ([]byte, error) {
	info := []byte(fmt.Sprintf("item-key:%s", itemID))
	return HKDFExtractAndExpand(vek, nil, info, 32)
}

// DeriveWrappingKey derives a key specifically for wrapping other keys
func DeriveWrappingKey(masterKey []byte) ([]byte, error) {
	return HKDFExtractAndExpand(masterKey, nil, []byte("wrap-vek"), 32)
}
