# Crypto Package

Zero-knowledge cryptographic primitives for the Yamony password manager.

## Overview

This package provides all cryptographic operations required for implementing zero-knowledge encryption, where the server stores only encrypted data and cannot decrypt vault contents.

## Security Model

The server **never** has access to:
- User master passwords
- Unencrypted vault data
- Encryption keys (except wrapped/encrypted keys)

All encryption and decryption happens client-side.

## Key Hierarchy

```
Master Password (user input)
    ↓ Argon2id
Master Key (MK)
    ↓ HKDF
Wrapping Key (WK)
    ↓ AES-GCM wraps
Vault Encryption Key (VEK) ← stored encrypted on server
    ↓ HKDF per-item
Item Encryption Key (IEK)
    ↓ AES-GCM
Encrypted Item Data
```

## Components

### Key Derivation (argon2.go)
- **Argon2id** for password-based key derivation
- Configurable parameters (time, memory, parallelism)
- Default params: 3 iterations, 64MB memory, 2 threads
- Mobile-optimized params available

```go
salt, _ := crypto.GenerateSalt(32)
params := crypto.DefaultKDFParams()
masterKey := crypto.DeriveMasterKey("user-password", salt, params)
```

### Symmetric Encryption (aes.go)
- **AES-256-GCM** for authenticated encryption
- 96-bit nonces, 128-bit authentication tags
- Support for Additional Authenticated Data (AAD)

```go
key, _ := crypto.GenerateRandomBytes(32)
encrypted, _ := crypto.EncryptAESGCM(key, plaintext, aad)
decrypted, _ := crypto.DecryptAESGCM(key, encrypted.Ciphertext, encrypted.IV, encrypted.Tag, aad)
```

### Key Derivation Function (hkdf.go)
- **HKDF-SHA256** for deriving multiple keys from a master key
- Per-vault and per-item key derivation
- Context separation for different use cases

```go
vek, _ := crypto.DeriveVaultEncryptionKey(masterKey, "vault-id")
iek, _ := crypto.DeriveItemEncryptionKey(vek, "item-id")
wrappingKey, _ := crypto.DeriveWrappingKey(masterKey)
```

### Digital Signatures (ed25519.go)
- **Ed25519** signatures for device authentication
- Used to sign all write operations
- Server verifies signatures before persisting data

```go
keyPair, _ := crypto.GenerateEd25519KeyPair()
signature := keyPair.Sign(message)
valid := crypto.VerifySignature(keyPair.PublicKey, message, signature)
```

### Key Exchange (x25519.go)
- **X25519** ECDH for secure key sharing
- Derives shared secrets between devices
- Used for vault sharing between users

```go
aliceKeyPair, _ := crypto.GenerateX25519KeyPair()
bobKeyPair, _ := crypto.GenerateX25519KeyPair()

// Both derive the same symmetric key
sharedKey, _ := crypto.DeriveSharedKey(
    aliceKeyPair.PrivateKey, 
    bobKeyPair.PublicKey, 
    "share-vek"
)
```

## High-Level Helpers (helpers.go)

### VaultKeyWrapper
Handles wrapping and unwrapping of Vault Encryption Keys:

```go
wrapper := crypto.NewVaultKeyWrapper(masterKey)
wrapped, _ := wrapper.WrapVEK(vek, aad)
unwrapped, _ := wrapper.UnwrapVEK(wrapped, aad)
```

### ItemEncryptor
Encrypts and decrypts vault items with per-item key derivation:

```go
encryptor := crypto.NewItemEncryptor(vek)
encrypted, _ := encryptor.EncryptItem(itemID, plaintext, aad)
decrypted, _ := encryptor.DecryptItem(itemID, encrypted, aad)

// JSON convenience methods
encrypted, _ := encryptor.EncryptItemJSON(itemID, dataStruct, aad)
err := encryptor.DecryptItemJSON(itemID, encrypted, aad, &targetStruct)
```

### ShareKeyWrapper
Wraps keys for sharing using ECDH:

```go
// Sender
senderWrapper := crypto.NewShareKeyWrapper(senderPrivateKey)
wrapped, _ := senderWrapper.WrapKeyForRecipient(recipientPubKey, vek, vaultID, aad)

// Recipient
recipientWrapper := crypto.NewShareKeyWrapper(recipientPrivateKey)
vek, _ := recipientWrapper.UnwrapKeyFromSender(senderPubKey, wrapped, vaultID, aad)
```

## Utilities

### Random Generation (random.go)
```go
bytes, _ := crypto.GenerateRandomBytes(32)
nonce, _ := crypto.GenerateNonce(12)
id, _ := crypto.GenerateID(16)
```

### Base64 Encoding (encoding.go)
```go
encoded := crypto.EncodeBase64(data)
decoded, _ := crypto.DecodeBase64(encoded)

// URL-safe variants
encodedURL := crypto.EncodeBase64URL(data)
decodedURL, _ := crypto.DecodeBase64URL(encodedURL)
```

## Usage Example: Complete Vault Flow

```go
// 1. User Registration
salt, _ := crypto.GenerateSalt(32)
params := crypto.DefaultKDFParams()
masterKey := crypto.DeriveMasterKey(password, salt, params)

// 2. Create Vault
vek, _ := crypto.GenerateRandomBytes(32)
wrapper := crypto.NewVaultKeyWrapper(masterKey)
wrappedVEK, _ := wrapper.WrapVEK(vek, []byte("vault-id"))
// Store wrappedVEK on server

// 3. Encrypt Item
encryptor := crypto.NewItemEncryptor(vek)
encrypted, _ := encryptor.EncryptItemJSON("item-id", loginItem, []byte("metadata"))
// Store encrypted on server

// 4. Share with Another User
senderKeyPair, _ := crypto.GenerateX25519KeyPair()
recipientPubKey := getRecipientPublicKey() // from server

shareWrapper := crypto.NewShareKeyWrapper(senderKeyPair.PrivateKey)
wrappedForRecipient, _ := shareWrapper.WrapKeyForRecipient(
    recipientPubKey, vek, "vault-id", []byte("sharing-context")
)
// Store sharing record on server
```

## Testing

Run the test suite:

```bash
go test ./internal/crypto/... -v
```

All cryptographic operations are thoroughly tested including:
- Key derivation determinism
- Encryption/decryption round-trips
- Signature verification
- ECDH key exchange
- High-level helper functions

## Security Considerations

1. **Never** log or expose master keys, VEKs, or IEKs
2. **Always** use constant-time comparison for secrets
3. **Always** verify signatures on write operations
4. **Never** reuse nonces (they're generated randomly)
5. Store private keys in OS secure storage (not implemented here)
6. Use TLS 1.3 for all network communication
7. Implement rate limiting on authentication endpoints

## Dependencies

- `golang.org/x/crypto/argon2` - Argon2id KDF
- `golang.org/x/crypto/curve25519` - X25519 ECDH
- `golang.org/x/crypto/hkdf` - HKDF key derivation
- `crypto/aes` - AES-GCM encryption
- `crypto/ed25519` - Ed25519 signatures
- `crypto/rand` - Cryptographically secure RNG

## References

- [Argon2 Specification](https://github.com/P-H-C/phc-winner-argon2)
- [RFC 5869 - HKDF](https://tools.ietf.org/html/rfc5869)
- [RFC 7748 - X25519 and Ed25519](https://tools.ietf.org/html/rfc7748)
- [NIST SP 800-38D - GCM](https://csrc.nist.gov/publications/detail/sp/800-38d/final)
