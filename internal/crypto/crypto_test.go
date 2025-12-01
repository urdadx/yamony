package crypto

import (
	"bytes"
	"crypto/ed25519"
	"testing"
)

func TestArgon2KeyDerivation(t *testing.T) {
	password := "test-password-123"
	salt, err := GenerateSalt(32)
	if err != nil {
		t.Fatalf("Failed to generate salt: %v", err)
	}

	params := DefaultKDFParams()
	key1 := DeriveMasterKey(password, salt, params)
	key2 := DeriveMasterKey(password, salt, params)

	if !bytes.Equal(key1, key2) {
		t.Error("Key derivation is not deterministic")
	}

	if len(key1) != int(params.KeyLen) {
		t.Errorf("Expected key length %d, got %d", params.KeyLen, len(key1))
	}
}

func TestAESGCMEncryptionDecryption(t *testing.T) {
	key, _ := GenerateRandomBytes(32)
	plaintext := []byte("sensitive vault data")
	aad := []byte("vault-metadata")

	// Test separated format
	encrypted, err := EncryptAESGCM(key, plaintext, aad)
	if err != nil {
		t.Fatalf("Encryption failed: %v", err)
	}

	decrypted, err := DecryptAESGCM(key, encrypted.Ciphertext, encrypted.IV, encrypted.Tag, aad)
	if err != nil {
		t.Fatalf("Decryption failed: %v", err)
	}

	if !bytes.Equal(plaintext, decrypted) {
		t.Error("Decrypted data does not match original")
	}

	// Test combined format
	combined, iv, err := EncryptAESGCMCombined(key, plaintext, aad)
	if err != nil {
		t.Fatalf("Combined encryption failed: %v", err)
	}

	decrypted2, err := DecryptAESGCMCombined(key, combined, iv, aad)
	if err != nil {
		t.Fatalf("Combined decryption failed: %v", err)
	}

	if !bytes.Equal(plaintext, decrypted2) {
		t.Error("Combined decrypted data does not match original")
	}
}

func TestHKDFKeyDerivation(t *testing.T) {
	secret := []byte("master-secret")

	key1, err := DeriveKey(secret, "context-1", 32)
	if err != nil {
		t.Fatalf("Key derivation failed: %v", err)
	}

	key2, err := DeriveKey(secret, "context-2", 32)
	if err != nil {
		t.Fatalf("Key derivation failed: %v", err)
	}

	if bytes.Equal(key1, key2) {
		t.Error("Different contexts should produce different keys")
	}

	if len(key1) != 32 || len(key2) != 32 {
		t.Error("Keys should be 32 bytes")
	}
}

func TestEd25519SignatureVerification(t *testing.T) {
	keyPair, err := GenerateEd25519KeyPair()
	if err != nil {
		t.Fatalf("Failed to generate key pair: %v", err)
	}

	message := []byte("test message")
	signature := keyPair.Sign(message)

	if !VerifySignature(keyPair.PublicKey, message, signature) {
		t.Error("Valid signature failed verification")
	}

	tamperedMessage := []byte("tampered message")
	if VerifySignature(keyPair.PublicKey, tamperedMessage, signature) {
		t.Error("Tampered message incorrectly verified")
	}

	if len(keyPair.PublicKey) != ed25519.PublicKeySize {
		t.Errorf("Public key size incorrect: expected %d, got %d",
			ed25519.PublicKeySize, len(keyPair.PublicKey))
	}
}

func TestX25519KeyExchange(t *testing.T) {
	// Alice generates key pair
	aliceKeyPair, err := GenerateX25519KeyPair()
	if err != nil {
		t.Fatalf("Failed to generate Alice's key pair: %v", err)
	}

	// Bob generates key pair
	bobKeyPair, err := GenerateX25519KeyPair()
	if err != nil {
		t.Fatalf("Failed to generate Bob's key pair: %v", err)
	}

	// Both derive the same shared secret
	aliceShared, err := X25519SharedSecret(aliceKeyPair.PrivateKey, bobKeyPair.PublicKey)
	if err != nil {
		t.Fatalf("Alice failed to derive shared secret: %v", err)
	}

	bobShared, err := X25519SharedSecret(bobKeyPair.PrivateKey, aliceKeyPair.PublicKey)
	if err != nil {
		t.Fatalf("Bob failed to derive shared secret: %v", err)
	}

	if !bytes.Equal(aliceShared, bobShared) {
		t.Error("Shared secrets do not match")
	}

	// Derive symmetric keys
	aliceKey, err := DeriveSharedKey(aliceKeyPair.PrivateKey, bobKeyPair.PublicKey, "share-vek")
	if err != nil {
		t.Fatalf("Alice failed to derive key: %v", err)
	}

	bobKey, err := DeriveSharedKey(bobKeyPair.PrivateKey, aliceKeyPair.PublicKey, "share-vek")
	if err != nil {
		t.Fatalf("Bob failed to derive key: %v", err)
	}

	if !bytes.Equal(aliceKey, bobKey) {
		t.Error("Derived keys do not match")
	}
}

func TestVaultKeyWrapper(t *testing.T) {
	password := "user-master-password"
	salt, _ := GenerateSalt(32)
	params := DefaultKDFParams()
	masterKey := DeriveMasterKey(password, salt, params)

	wrapper := NewVaultKeyWrapper(masterKey)
	vek, _ := GenerateRandomBytes(32)
	aad := []byte("vault-123")

	// Wrap VEK
	wrapped, err := wrapper.WrapVEK(vek, aad)
	if err != nil {
		t.Fatalf("Failed to wrap VEK: %v", err)
	}

	// Unwrap VEK
	unwrapped, err := wrapper.UnwrapVEK(wrapped, aad)
	if err != nil {
		t.Fatalf("Failed to unwrap VEK: %v", err)
	}

	if !bytes.Equal(vek, unwrapped) {
		t.Error("Unwrapped VEK does not match original")
	}
}

func TestItemEncryptor(t *testing.T) {
	vek, _ := GenerateRandomBytes(32)
	encryptor := NewItemEncryptor(vek)

	itemID := "item-uuid-123"
	plaintext := []byte(`{"username":"test","password":"secret"}`)
	aad := []byte("login-item")

	// Encrypt item
	encrypted, err := encryptor.EncryptItem(itemID, plaintext, aad)
	if err != nil {
		t.Fatalf("Failed to encrypt item: %v", err)
	}

	// Decrypt item
	decrypted, err := encryptor.DecryptItem(itemID, encrypted, aad)
	if err != nil {
		t.Fatalf("Failed to decrypt item: %v", err)
	}

	if !bytes.Equal(plaintext, decrypted) {
		t.Error("Decrypted item does not match original")
	}

	// Test JSON convenience methods
	type LoginItem struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	original := LoginItem{Username: "test", Password: "secret"}
	encrypted2, err := encryptor.EncryptItemJSON(itemID, original, aad)
	if err != nil {
		t.Fatalf("Failed to encrypt item JSON: %v", err)
	}

	var decrypted2 LoginItem
	err = encryptor.DecryptItemJSON(itemID, encrypted2, aad, &decrypted2)
	if err != nil {
		t.Fatalf("Failed to decrypt item JSON: %v", err)
	}

	if original.Username != decrypted2.Username || original.Password != decrypted2.Password {
		t.Error("Decrypted JSON item does not match original")
	}
}

func TestShareKeyWrapper(t *testing.T) {
	// Sender key pair
	senderKeyPair, _ := GenerateX25519KeyPair()
	// Recipient key pair
	recipientKeyPair, _ := GenerateX25519KeyPair()

	vaultID := "vault-456"
	keyToShare, _ := GenerateRandomBytes(32)
	aad := []byte("sharing-context")

	// Sender wraps key for recipient
	senderWrapper := NewShareKeyWrapper(senderKeyPair.PrivateKey)
	wrapped, err := senderWrapper.WrapKeyForRecipient(recipientKeyPair.PublicKey, keyToShare, vaultID, aad)
	if err != nil {
		t.Fatalf("Failed to wrap key for recipient: %v", err)
	}

	// Recipient unwraps key from sender
	recipientWrapper := NewShareKeyWrapper(recipientKeyPair.PrivateKey)
	unwrapped, err := recipientWrapper.UnwrapKeyFromSender(senderKeyPair.PublicKey, wrapped, vaultID, aad)
	if err != nil {
		t.Fatalf("Failed to unwrap key from sender: %v", err)
	}

	if !bytes.Equal(keyToShare, unwrapped) {
		t.Error("Unwrapped key does not match original")
	}
}

func TestBase64Encoding(t *testing.T) {
	data := []byte("test data 123")

	encoded := EncodeBase64(data)
	decoded, err := DecodeBase64(encoded)
	if err != nil {
		t.Fatalf("Failed to decode: %v", err)
	}

	if !bytes.Equal(data, decoded) {
		t.Error("Decoded data does not match original")
	}

	// Test URL-safe encoding
	encodedURL := EncodeBase64URL(data)
	decodedURL, err := DecodeBase64URL(encodedURL)
	if err != nil {
		t.Fatalf("Failed to decode URL: %v", err)
	}

	if !bytes.Equal(data, decodedURL) {
		t.Error("Decoded URL data does not match original")
	}
}
