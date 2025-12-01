// Package crypto provides cryptographic primitives for zero-knowledge encryption.
//
// This package implements the cryptographic operations required for a zero-knowledge
// password manager, including:
//   - Argon2id key derivation for master password hashing
//   - AES-256-GCM symmetric encryption/decryption
//   - HKDF for key derivation and key hierarchies
//   - Ed25519 for digital signatures (device authentication)
//   - X25519 for ECDH key exchange (secure sharing)
//
// Security Model:
//
// The server stores only encrypted data and cannot decrypt vault contents.
// Clients derive encryption keys from the master password using Argon2id,
// then use HKDF to derive per-vault and per-item keys.
//
// Key Hierarchy:
//  1. Master Password (user input) → Argon2id → Master Key (MK)
//  2. MK → HKDF → Wrapping Key (WK)
//  3. WK wraps VEK (Vault Encryption Key) → stored on server
//  4. VEK → HKDF → Item Encryption Key (IEK) per item
//
// Device Authentication:
//   - Each device has Ed25519 keypair for signatures
//   - Device signs all write operations
//   - Server verifies signatures before persisting
//
// Secure Sharing:
//   - Each device has X25519 keypair for ECDH
//   - Sender and recipient derive shared secret
//   - Shared secret wraps the VEK or IEK
//
// Example usage:
//
//	// Key derivation
//	salt, _ := crypto.GenerateSalt(32)
//	params := crypto.DefaultKDFParams()
//	masterKey := crypto.DeriveMasterKey(password, salt, params)
//	wrappingKey, _ := crypto.DeriveWrappingKey(masterKey)
//
//	// Encrypt vault data
//	vek, _ := crypto.GenerateRandomBytes(32)
//	encrypted, _ := crypto.EncryptAESGCM(wrappingKey, vek, nil)
//
//	// Device signing
//	keyPair, _ := crypto.GenerateEd25519KeyPair()
//	signature := keyPair.Sign(message)
//	valid := crypto.VerifySignature(keyPair.PublicKey, message, signature)
//
//	// Secure sharing
//	sharedKey, _ := crypto.DeriveSharedKey(myPrivKey, theirPubKey, "share-vek")
//	wrappedKey, _ := crypto.EncryptAESGCM(sharedKey, vek, nil)
package crypto
