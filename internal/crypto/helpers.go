package crypto

import (
	"encoding/json"
	"fmt"
)

// VaultKeyWrapper handles the encryption and decryption of Vault Encryption Keys
type VaultKeyWrapper struct {
	masterKey []byte
}

// NewVaultKeyWrapper creates a new wrapper with the given master key
func NewVaultKeyWrapper(masterKey []byte) *VaultKeyWrapper {
	return &VaultKeyWrapper{masterKey: masterKey}
}

// WrapVEK encrypts a Vault Encryption Key for storage
func (w *VaultKeyWrapper) WrapVEK(vek []byte, aad []byte) (*EncryptedData, error) {
	wrappingKey, err := DeriveWrappingKey(w.masterKey)
	if err != nil {
		return nil, fmt.Errorf("failed to derive wrapping key: %w", err)
	}

	encrypted, err := EncryptAESGCM(wrappingKey, vek, aad)
	if err != nil {
		return nil, fmt.Errorf("failed to wrap VEK: %w", err)
	}

	return encrypted, nil
}

// UnwrapVEK decrypts a Vault Encryption Key
func (w *VaultKeyWrapper) UnwrapVEK(encrypted *EncryptedData, aad []byte) ([]byte, error) {
	wrappingKey, err := DeriveWrappingKey(w.masterKey)
	if err != nil {
		return nil, fmt.Errorf("failed to derive wrapping key: %w", err)
	}

	vek, err := DecryptAESGCM(wrappingKey, encrypted.Ciphertext, encrypted.IV, encrypted.Tag, aad)
	if err != nil {
		return nil, fmt.Errorf("failed to unwrap VEK: %w", err)
	}

	return vek, nil
}

// ItemEncryptor handles encryption and decryption of vault items
type ItemEncryptor struct {
	vek []byte
}

// NewItemEncryptor creates a new encryptor with the given VEK
func NewItemEncryptor(vek []byte) *ItemEncryptor {
	return &ItemEncryptor{vek: vek}
}

// EncryptItem encrypts a vault item using per-item key derivation
func (e *ItemEncryptor) EncryptItem(itemID string, plaintext []byte, aad []byte) (*EncryptedData, error) {
	itemKey, err := DeriveItemEncryptionKey(e.vek, itemID)
	if err != nil {
		return nil, fmt.Errorf("failed to derive item key: %w", err)
	}

	encrypted, err := EncryptAESGCM(itemKey, plaintext, aad)
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt item: %w", err)
	}

	return encrypted, nil
}

// DecryptItem decrypts a vault item using per-item key derivation
func (e *ItemEncryptor) DecryptItem(itemID string, encrypted *EncryptedData, aad []byte) ([]byte, error) {
	itemKey, err := DeriveItemEncryptionKey(e.vek, itemID)
	if err != nil {
		return nil, fmt.Errorf("failed to derive item key: %w", err)
	}

	plaintext, err := DecryptAESGCM(itemKey, encrypted.Ciphertext, encrypted.IV, encrypted.Tag, aad)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt item: %w", err)
	}

	return plaintext, nil
}

// EncryptItemJSON is a convenience method for encrypting JSON-serializable data
func (e *ItemEncryptor) EncryptItemJSON(itemID string, data interface{}, aad []byte) (*EncryptedData, error) {
	plaintext, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal item: %w", err)
	}

	return e.EncryptItem(itemID, plaintext, aad)
}

// DecryptItemJSON is a convenience method for decrypting into a struct
func (e *ItemEncryptor) DecryptItemJSON(itemID string, encrypted *EncryptedData, aad []byte, target interface{}) error {
	plaintext, err := e.DecryptItem(itemID, encrypted, aad)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(plaintext, target); err != nil {
		return fmt.Errorf("failed to unmarshal item: %w", err)
	}

	return nil
}

// ShareKeyWrapper handles wrapping keys for sharing using ECDH
type ShareKeyWrapper struct {
	privateKey []byte
}

// NewShareKeyWrapper creates a new wrapper with the device's X25519 private key
func NewShareKeyWrapper(privateKey []byte) *ShareKeyWrapper {
	return &ShareKeyWrapper{privateKey: privateKey}
}

// WrapKeyForRecipient wraps a key for a specific recipient using ECDH
func (w *ShareKeyWrapper) WrapKeyForRecipient(recipientPublicKey, keyToWrap []byte, vaultID string, aad []byte) (*EncryptedData, error) {
	// Derive shared key using ECDH
	sharedKey, err := DeriveSharedKeyForSharing(w.privateKey, recipientPublicKey, vaultID)
	if err != nil {
		return nil, fmt.Errorf("failed to derive shared key: %w", err)
	}

	// Encrypt the key with the shared key
	encrypted, err := EncryptAESGCM(sharedKey, keyToWrap, aad)
	if err != nil {
		return nil, fmt.Errorf("failed to wrap key: %w", err)
	}

	return encrypted, nil
}

// UnwrapKeyFromSender unwraps a key received from a sender using ECDH
func (w *ShareKeyWrapper) UnwrapKeyFromSender(senderPublicKey []byte, encrypted *EncryptedData, vaultID string, aad []byte) ([]byte, error) {
	// Derive the same shared key using ECDH
	sharedKey, err := DeriveSharedKeyForSharing(w.privateKey, senderPublicKey, vaultID)
	if err != nil {
		return nil, fmt.Errorf("failed to derive shared key: %w", err)
	}

	// Decrypt the key
	key, err := DecryptAESGCM(sharedKey, encrypted.Ciphertext, encrypted.IV, encrypted.Tag, aad)
	if err != nil {
		return nil, fmt.Errorf("failed to unwrap key: %w", err)
	}

	return key, nil
}
