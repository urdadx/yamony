# Yamony API Documentation

## Overview

Yamony is a zero-knowledge password manager with end-to-end encryption. The API implements a secure architecture where:

- All sensitive data is encrypted client-side before transmission
- The server only stores encrypted blobs and cannot decrypt user data
- Device-based authentication with Ed25519 signatures
- Per-vault encryption keys (VEK) wrapped with user's master key
- Per-item encryption keys (IEK) derived from VEK using HKDF
- Secure sharing using X25519 ECDH key exchange

## Base URL

```
http://localhost:8080/api
```

## Authentication

All protected endpoints require session authentication. After login, a session cookie is set automatically.

### Protected Headers

For write operations (POST, PUT, DELETE) on vault data:

- `X-Device-Id`: UUID of the registered device
- `X-Device-Signature`: Ed25519 signature of canonical message
- `X-Device-Timestamp`: Unix timestamp in milliseconds

## Endpoints

### Authentication

#### Register User
```http
POST /register
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "secure_password",
  "name": "John Doe"
}
```

#### Login
```http
POST /login
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "secure_password"
}
```

#### Get Current User
```http
GET /me
```

Response:
```json
{
  "id": 1,
  "email": "user@example.com",
  "name": "John Doe"
}
```

### Device Management

#### Register Device
```http
POST /devices/register
Content-Type: application/json

{
  "device_label": "MacBook Pro",
  "x25519_public": "base64_encoded_x25519_public_key",
  "ed25519_public": "base64_encoded_ed25519_public_key"
}
```

Response includes a challenge for verification.

#### Verify Device
```http
POST /devices/verify
Content-Type: application/json

{
  "device_id": "uuid",
  "challenge_signature": "base64_encoded_signature"
}
```

#### List Devices
```http
GET /devices
```

#### Revoke Device
```http
DELETE /devices/{device_id}
X-Device-Id: {device_id}
X-Device-Signature: {signature}
X-Device-Timestamp: {timestamp}
```

#### Get User Public Keys
```http
GET /users/{user_id}/public-keys
```

Returns all active devices' public keys for sharing.

### Vault Keys

#### Upload Vault Key
```http
POST /vaults/{vault_id}/keys
Content-Type: application/json

{
  "wrapped_vek": "base64_encoded_wrapped_vault_encryption_key",
  "wrap_iv": "base64_encoded_iv",
  "wrap_tag": "base64_encoded_auth_tag",
  "kdf_salt": "base64_encoded_salt",
  "kdf_params": {
    "time": 3,
    "memory": 65536,
    "parallelism": 2,
    "key_len": 32
  },
  "version": 1
}
```

The VEK is wrapped with a key derived from the master password using Argon2id.

#### Get Vault Key
```http
GET /vaults/{vault_id}/keys
```

Returns the latest version of the wrapped VEK.

#### Get Vault Key Version
```http
GET /vaults/{vault_id}/keys/versions/{version}
```

#### List Vault Key Versions
```http
GET /vaults/{vault_id}/keys/versions
```

### Vault Items

#### Create Vault Item
```http
POST /vaults/{vault_id}/items
Content-Type: application/json
X-Device-Id: {device_id}
X-Device-Signature: {signature}
X-Device-Timestamp: {timestamp}

{
  "item_type": "login",
  "encrypted_blob": "base64_encoded_encrypted_item_data",
  "iv": "base64_encoded_iv",
  "tag": "base64_encoded_auth_tag",
  "meta": {
    "title": "Example Login",
    "url": "https://example.com"
  },
  "version": 1
}
```

Item types: `login`, `note`, `card`, `alias`

#### List Vault Items
```http
GET /vaults/{vault_id}/items?type=login
```

Query parameters:
- `type` (optional): Filter by item type

Returns simplified list without encrypted blobs for efficiency.

#### Get Vault Item
```http
GET /vaults/{vault_id}/items/{item_id}
```

Returns full item including encrypted blob.

#### Update Vault Item
```http
PUT /vaults/{vault_id}/items/{item_id}
Content-Type: application/json
X-Device-Id: {device_id}
X-Device-Signature: {signature}
X-Device-Timestamp: {timestamp}

{
  "encrypted_blob": "base64_encoded_updated_data",
  "iv": "base64_encoded_iv",
  "tag": "base64_encoded_auth_tag",
  "meta": {...},
  "base_version": 1
}
```

Uses optimistic concurrency control. Returns `409 Conflict` if version mismatch.

#### Delete Vault Item
```http
DELETE /vaults/{vault_id}/items/{item_id}
X-Device-Id: {device_id}
X-Device-Signature: {signature}
X-Device-Timestamp: {timestamp}
```

### Sharing

#### Share Vault
```http
POST /vaults/{vault_id}/share
Content-Type: application/json
X-Device-Id: {device_id}
X-Device-Signature: {signature}
X-Device-Timestamp: {timestamp}

{
  "recipient_user_id": 123,
  "wrapped_vek": "base64_encoded_vek_wrapped_with_ecdh_key",
  "wrap_iv": "base64_encoded_iv",
  "wrap_tag": "base64_encoded_auth_tag"
}
```

The VEK is wrapped using a symmetric key derived from X25519 ECDH shared secret.

#### Get Pending Shares
```http
GET /shares/pending
```

Returns sharing invitations awaiting acceptance.

#### Accept Share
```http
POST /shares/{share_id}/accept
```

#### Reject Share
```http
POST /shares/{share_id}/reject
```

#### Get Shared Vaults
```http
GET /vaults/shared
```

Returns all vaults shared with the current user.

#### Revoke Share
```http
DELETE /shares/{share_id}
X-Device-Id: {device_id}
X-Device-Signature: {signature}
X-Device-Timestamp: {timestamp}
```

Only vault owner can revoke shares.

### Synchronization

#### Pull Vault Changes
```http
POST /vaults/{vault_id}/sync/pull
Content-Type: application/json
If-None-Match: {etag}

{
  "last_synced_version_id": 42
}
```

Returns:
- 200 OK with vault state
- 304 Not Modified if ETag matches

Response includes `ETag` header for optimistic concurrency.

#### Commit Vault Changes
```http
POST /vaults/{vault_id}/sync/commit
Content-Type: application/json
X-Device-Id: {device_id}
X-Device-Signature: {signature}
X-Device-Timestamp: {timestamp}
If-Match: {etag}

{
  "base_version_id": 42,
  "items": [
    {
      "id": null,
      "item_type": "login",
      "encrypted_blob": "...",
      "iv": "...",
      "tag": "...",
      "meta": {...}
    },
    {
      "id": "existing-item-uuid",
      "item_type": "note",
      "encrypted_blob": "...",
      "iv": "...",
      "tag": "...",
      "meta": {...},
      "base_version": 2
    }
  ],
  "deleted_items": ["uuid1", "uuid2"]
}
```

Returns:
- 200 OK if successful
- 409 Conflict if version conflicts detected
- 412 Precondition Failed if If-Match header doesn't match current ETag

#### Get Vault Versions
```http
GET /vaults/{vault_id}/versions?limit=10
```

Returns version history with metadata.

## Cryptographic Operations

### Key Hierarchy

```
Master Password
    ↓ (Argon2id)
Master Key (MK)
    ↓ (HKDF with "wrap-vek")
Wrapping Key
    ↓ (AES-256-GCM)
Vault Encryption Key (VEK) [stored wrapped]
    ↓ (HKDF with "item-key:{item_id}")
Item Encryption Key (IEK)
    ↓ (AES-256-GCM)
Item Data [stored encrypted]
```

### Device Signature Verification

For all write operations, the client must sign a canonical message:

```
{HTTP_METHOD}|{URL_PATH}|{TIMESTAMP}|{SHA256(BODY)}
```

Example:
```
POST|/api/vaults/1/items|1701432000000|e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855
```

### Sharing with ECDH

When sharing a vault:

1. Sender fetches recipient's X25519 public key
2. Sender computes ECDH shared secret: `shared = X25519(sender_private, recipient_public)`
3. Derive symmetric key: `wrap_key = HKDF(shared, salt=vault_id, info="vault-share")`
4. Wrap VEK: `wrapped_vek = AES-256-GCM(wrap_key, vek)`
5. Send wrapped VEK to recipient

Recipient unwraps using their private key and sender's public key.

## Error Responses

All errors follow this format:

```json
{
  "error": "error message"
}
```

Common status codes:
- 400 Bad Request - Invalid input
- 401 Unauthorized - Not authenticated or invalid signature
- 403 Forbidden - No access to resource
- 404 Not Found - Resource doesn't exist
- 409 Conflict - Version conflict (optimistic concurrency)
- 412 Precondition Failed - ETag mismatch
- 500 Internal Server Error - Server error

## Rate Limiting

Not currently implemented. Consider adding rate limiting for production deployments.

## Security Considerations

1. **Zero Knowledge**: Server never has access to plaintext data
2. **Device Authentication**: All writes require Ed25519 signatures
3. **Perfect Forward Secrecy**: Each item uses a unique encryption key
4. **Secure Sharing**: ECDH ensures only intended recipients can decrypt
5. **Version Control**: Optimistic concurrency prevents data loss
6. **Audit Trail**: Version history tracks all changes

## Client Implementation Guide

### Initial Setup Flow

1. Register user account
2. Derive master key from password: `MK = Argon2id(password, salt)`
3. Generate device keypairs (X25519 + Ed25519)
4. Register device with public keys
5. Verify device by signing challenge
6. Create vault and generate VEK
7. Derive wrapping key and wrap VEK: `wrapped_vek = encrypt(derive_wrapping_key(MK), VEK)`
8. Upload wrapped VEK to server

### Adding an Item

1. Fetch wrapped VEK from server
2. Unwrap VEK: `VEK = decrypt(derive_wrapping_key(MK), wrapped_vek)`
3. Derive item encryption key: `IEK = HKDF(VEK, "item-key:{item_id}")`
4. Encrypt item data: `encrypted_blob = AES-256-GCM(IEK, item_data)`
5. Sign request with device Ed25519 key
6. Upload encrypted item to server

### Syncing

1. Pull current vault state with ETag
2. Make local changes
3. Commit changes with If-Match header containing ETag
4. Handle conflicts if returned
5. Repeat until successful

## Example Code Snippets

### JavaScript Client - Device Registration

```javascript
// Generate keypairs
const x25519KeyPair = await crypto.subtle.generateKey(
  { name: "X25519" },
  true,
  ["deriveKey"]
);

const ed25519KeyPair = await crypto.subtle.generateKey(
  { name: "Ed25519" },
  true,
  ["sign", "verify"]
);

// Export public keys
const x25519Public = await crypto.subtle.exportKey("raw", x25519KeyPair.publicKey);
const ed25519Public = await crypto.subtle.exportKey("raw", ed25519KeyPair.publicKey);

// Register device
const response = await fetch('/api/devices/register', {
  method: 'POST',
  headers: { 'Content-Type': 'application/json' },
  body: JSON.stringify({
    device_label: 'My Device',
    x25519_public: btoa(String.fromCharCode(...new Uint8Array(x25519Public))),
    ed25519_public: btoa(String.fromCharCode(...new Uint8Array(ed25519Public)))
  })
});

const { device_id, challenge } = await response.json();

// Sign challenge
const signature = await crypto.subtle.sign(
  "Ed25519",
  ed25519KeyPair.privateKey,
  new TextEncoder().encode(challenge)
);

// Verify device
await fetch('/api/devices/verify', {
  method: 'POST',
  headers: { 'Content-Type': 'application/json' },
  body: JSON.stringify({
    device_id,
    challenge_signature: btoa(String.fromCharCode(...new Uint8Array(signature)))
  })
});
```

### Go Client - Creating Encrypted Item

```go
import "yamony/internal/crypto"

// Derive item encryption key
vek := []byte{...} // 32-byte VEK
itemID := uuid.New()
iek := crypto.DeriveItemEncryptionKey(vek, itemID)

// Encrypt item data
itemData := map[string]string{
    "username": "user@example.com",
    "password": "secret123",
}
encryptor := crypto.NewItemEncryptor(iek)
encrypted, err := encryptor.EncryptItemJSON(itemData)
if err != nil {
    log.Fatal(err)
}

// Upload to server
client.CreateVaultItem(vaultID, VaultItemRequest{
    ItemType:      "login",
    EncryptedBlob: crypto.EncodeBase64(encrypted.Ciphertext),
    IV:            crypto.EncodeBase64(encrypted.IV),
    Tag:           crypto.EncodeBase64(encrypted.Tag),
    Meta:          json.RawMessage(`{"title": "Example Login"}`),
})
```
