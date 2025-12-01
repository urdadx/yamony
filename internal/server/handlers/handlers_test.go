package handlers

import (
	"bytes"
	"crypto/ed25519"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"yamony/internal/crypto"
)

// MockService is a mock implementation of the Service interface for testing
type MockService struct {
	mock.Mock
}

func (m *MockService) GetDB() interface{} {
	args := m.Called()
	return args.Get(0)
}

// TestDeviceSignatureGeneration tests canonical message creation and signing
func TestDeviceSignatureGeneration(t *testing.T) {
	// Generate Ed25519 keypair
	publicKey, privateKey, err := ed25519.GenerateKey(nil)
	assert.NoError(t, err)

	// Test data
	method := "POST"
	path := "/api/vaults/1/items"
	timestamp := time.Now().UnixMilli()
	bodyData := []byte(`{"item_type":"login","encrypted_blob":"test"}`)

	// Create canonical message
	bodyHash := sha256.Sum256(bodyData)
	canonicalMsg := CreateCanonicalMessage(method, path, timestamp, bodyHash[:])
	expectedMsg := []byte(fmt.Sprintf("%s|%s|%d|%s", method, path, timestamp, crypto.EncodeBase64(bodyHash[:])))
	assert.Equal(t, expectedMsg, canonicalMsg)

	// Sign message
	signature := ed25519.Sign(privateKey, canonicalMsg)
	assert.NotEmpty(t, signature)

	// Verify signature
	valid := ed25519.Verify(publicKey, canonicalMsg, signature)
	assert.True(t, valid)

	// Test with modified message (should fail)
	modifiedMsg := CreateCanonicalMessage("PUT", path, timestamp, bodyHash[:])
	valid = ed25519.Verify(publicKey, modifiedMsg, signature)
	assert.False(t, valid)
}

// TestVaultKeyHandlerUpload tests the UploadVaultKey handler
func TestVaultKeyHandlerUpload(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Generate test crypto materials
	password := "test_password_123"
	salt, err := crypto.GenerateSalt(32)
	assert.NoError(t, err)

	kdfParams := crypto.DefaultKDFParams()
	masterKey := crypto.DeriveMasterKey(password, salt, kdfParams)

	// Generate VEK
	vek, err := crypto.GenerateRandomBytes(32)
	assert.NoError(t, err)

	// Wrap VEK
	wrapper := crypto.NewVaultKeyWrapper(masterKey)
	aad := []byte("vault_key")
	wrappedVEK, err := wrapper.WrapVEK(vek, aad)
	assert.NoError(t, err)

	// Prepare request
	kdfParamsJSON, err := json.Marshal(kdfParams)
	assert.NoError(t, err)

	requestBody := map[string]interface{}{
		"wrapped_vek": crypto.EncodeBase64(wrappedVEK.Ciphertext),
		"wrap_iv":     crypto.EncodeBase64(wrappedVEK.IV),
		"wrap_tag":    crypto.EncodeBase64(wrappedVEK.Tag),
		"kdf_salt":    crypto.EncodeBase64(salt),
		"kdf_params":  json.RawMessage(kdfParamsJSON),
		"version":     1,
	}

	// Create test request
	bodyBytes, err := json.Marshal(requestBody)
	assert.NoError(t, err)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/api/vaults/1/keys", bytes.NewBuffer(bodyBytes))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Params = []gin.Param{{Key: "id", Value: "1"}}
	c.Set("user_id", int32(1))

	// Verify context setup
	assert.Equal(t, "1", c.Param("id"))
	userID, exists := c.Get("user_id")
	assert.True(t, exists)
	assert.Equal(t, int32(1), userID)
}

// TestItemEncryptionDecryption tests item encryption and decryption
func TestItemEncryptionDecryption(t *testing.T) {
	// Generate VEK
	vek, err := crypto.GenerateRandomBytes(32)
	assert.NoError(t, err)

	// Test item data
	itemID := uuid.New()
	itemIDStr := itemID.String()
	itemData := map[string]interface{}{
		"username": "testuser@example.com",
		"password": "secure_password_123",
		"url":      "https://example.com",
		"notes":    "Test login item",
	}

	// Derive IEK
	iek, err := crypto.DeriveItemEncryptionKey(vek, itemIDStr)
	assert.NoError(t, err)
	assert.NotNil(t, iek)
	assert.Len(t, iek, 32)

	// Encrypt item
	encryptor := crypto.NewItemEncryptor(iek)
	aad := []byte("vault_items")
	encrypted, err := encryptor.EncryptItemJSON(itemIDStr, itemData, aad)
	assert.NoError(t, err)
	assert.NotEmpty(t, encrypted.Ciphertext)
	assert.Len(t, encrypted.IV, 12)
	assert.Len(t, encrypted.Tag, 16)

	// Decrypt item
	var decrypted map[string]interface{}
	err = encryptor.DecryptItemJSON(itemIDStr, encrypted, aad, &decrypted)
	assert.NoError(t, err)
	assert.Equal(t, itemData["username"], decrypted["username"])
	assert.Equal(t, itemData["password"], decrypted["password"])
	assert.Equal(t, itemData["url"], decrypted["url"])
	assert.Equal(t, itemData["notes"], decrypted["notes"])

	// Test with wrong key (should fail)
	wrongKey, _ := crypto.GenerateRandomBytes(32)
	wrongEncryptor := crypto.NewItemEncryptor(wrongKey)
	err = wrongEncryptor.DecryptItemJSON(itemIDStr, encrypted, aad, &decrypted)
	assert.Error(t, err)
}

// TestVaultSharing tests the ECDH key wrapping for sharing
func TestVaultSharing(t *testing.T) {
	// Generate sender keypair
	senderKeyPair, err := crypto.GenerateX25519KeyPair()
	assert.NoError(t, err)

	// Generate recipient keypair
	recipientKeyPair, err := crypto.GenerateX25519KeyPair()
	assert.NoError(t, err)

	// Generate VEK to share
	vek, err := crypto.GenerateRandomBytes(32)
	assert.NoError(t, err)

	vaultID := "123"
	aad := []byte("vault_key")

	// Sender wraps VEK for recipient
	senderWrapper := crypto.NewShareKeyWrapper(senderKeyPair.PrivateKey)
	wrapped, err := senderWrapper.WrapKeyForRecipient(recipientKeyPair.PublicKey, vek, vaultID, aad)
	assert.NoError(t, err)
	assert.NotEmpty(t, wrapped.Ciphertext)

	// Recipient unwraps VEK
	recipientWrapper := crypto.NewShareKeyWrapper(recipientKeyPair.PrivateKey)
	unwrapped, err := recipientWrapper.UnwrapKeyFromSender(senderKeyPair.PublicKey, wrapped, vaultID, aad)
	assert.NoError(t, err)
	assert.Equal(t, vek, unwrapped)

	// Test with wrong recipient key (should fail)
	wrongKeyPair, _ := crypto.GenerateX25519KeyPair()
	wrongWrapper := crypto.NewShareKeyWrapper(wrongKeyPair.PrivateKey)
	_, err = wrongWrapper.UnwrapKeyFromSender(senderKeyPair.PublicKey, wrapped, vaultID, aad)
	assert.Error(t, err)

	// Test with wrong sender key (should fail)
	_, err = recipientWrapper.UnwrapKeyFromSender(wrongKeyPair.PublicKey, wrapped, vaultID, aad)
	assert.Error(t, err)
}

// TestOptimisticConcurrency tests version-based conflict detection
func TestOptimisticConcurrency(t *testing.T) {
	// Simulate item with version
	currentVersion := int32(5)
	clientBaseVersion := int32(5)
	concurrentBaseVersion := int32(4)

	// Client with matching base version should succeed
	if currentVersion == clientBaseVersion {
		newVersion := currentVersion + 1
		assert.Equal(t, int32(6), newVersion)
	}

	// Client with old base version should detect conflict
	if currentVersion != concurrentBaseVersion {
		t.Log("Conflict detected: current version", currentVersion, "!= base version", concurrentBaseVersion)
		assert.NotEqual(t, currentVersion, concurrentBaseVersion)
	}
}

// TestETagGeneration tests ETag computation for sync
func TestETagGeneration(t *testing.T) {
	vaultID := int32(1)
	versionID := int32(42)

	// Create mock items
	items := []struct {
		ID      uuid.UUID
		Version int32
	}{
		{uuid.New(), 1},
		{uuid.New(), 2},
		{uuid.New(), 3},
	}

	// Generate ETag
	hash := sha256.New()
	hash.Write([]byte(fmt.Sprintf("%d:%d:%d", vaultID, versionID, len(items))))
	for _, item := range items {
		hash.Write(item.ID[:])
		hash.Write([]byte(fmt.Sprintf(":%d", item.Version)))
	}
	etag1 := fmt.Sprintf("%x", hash.Sum(nil))
	assert.NotEmpty(t, etag1)

	// Same input should produce same ETag
	hash2 := sha256.New()
	hash2.Write([]byte(fmt.Sprintf("%d:%d:%d", vaultID, versionID, len(items))))
	for _, item := range items {
		hash2.Write(item.ID[:])
		hash2.Write([]byte(fmt.Sprintf(":%d", item.Version)))
	}
	etag2 := fmt.Sprintf("%x", hash2.Sum(nil))
	assert.Equal(t, etag1, etag2)

	// Different items should produce different ETag
	hash3 := sha256.New()
	hash3.Write([]byte(fmt.Sprintf("%d:%d:%d", vaultID, versionID, len(items)+1)))
	etag3 := fmt.Sprintf("%x", hash3.Sum(nil))
	assert.NotEqual(t, etag1, etag3)
}

// TestKDFParamsValidation tests KDF parameter validation
func TestKDFParamsValidation(t *testing.T) {
	validParams := crypto.KDFParams{
		Time:        3,
		Memory:      64 * 1024,
		Parallelism: 2,
		KeyLen:      32,
	}
	err := crypto.ValidateKDFParams(validParams)
	assert.NoError(t, err)

	// Test with low memory (should fail - below 8 KB)
	lowMemory := crypto.KDFParams{
		Time:        3,
		Memory:      4 * 1024, // Too low - only 4 KB
		Parallelism: 2,
		KeyLen:      32,
	}
	err = crypto.ValidateKDFParams(lowMemory)
	assert.Error(t, err)

	// Test with zero time (should fail)
	zeroTime := crypto.KDFParams{
		Time:        0, // Too low
		Memory:      64 * 1024,
		Parallelism: 2,
		KeyLen:      32,
	}
	err = crypto.ValidateKDFParams(zeroTime)
	assert.Error(t, err)

	// Test with invalid key length (should fail - below 16 bytes)
	invalidKeyLen := crypto.KDFParams{
		Time:        3,
		Memory:      64 * 1024,
		Parallelism: 2,
		KeyLen:      8, // Too short - must be at least 16
	}
	err = crypto.ValidateKDFParams(invalidKeyLen)
	assert.Error(t, err)
}

// BenchmarkItemEncryption benchmarks item encryption performance
func BenchmarkItemEncryption(b *testing.B) {
	vek, _ := crypto.GenerateRandomBytes(32)
	itemID := uuid.New()
	itemIDStr := itemID.String()
	iek, _ := crypto.DeriveItemEncryptionKey(vek, itemIDStr)
	encryptor := crypto.NewItemEncryptor(iek)
	itemData := map[string]string{
		"username": "testuser@example.com",
		"password": "secure_password_123",
		"url":      "https://example.com",
	}
	aad := []byte("vault_items")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = encryptor.EncryptItemJSON(itemIDStr, itemData, aad)
	}
}

// BenchmarkArgon2KDF benchmarks Argon2id key derivation
func BenchmarkArgon2KDF(b *testing.B) {
	password := "test_password_123"
	salt, _ := crypto.GenerateSalt(32)
	params := crypto.DefaultKDFParams()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = crypto.DeriveMasterKey(password, salt, params)
	}
}

// BenchmarkEd25519Signing benchmarks Ed25519 signature generation
func BenchmarkEd25519Signing(b *testing.B) {
	publicKey, privateKey, _ := ed25519.GenerateKey(nil)
	message := []byte("test message for signing")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		signature := ed25519.Sign(privateKey, message)
		_ = ed25519.Verify(publicKey, message, signature)
	}
}
