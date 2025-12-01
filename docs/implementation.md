# IMPLEMENTATION.md

## 1. Goals & constraints
- **Zero-knowledge** server: The server must store **only** ciphertext and public keys. The server must not be able to decrypt vault contents.
- Use standard, well-reviewed cryptographic primitives and libraries.
- Implement wrapped per-vault key (VEK) model: the server stores a wrapped VEK; clients derive a key from the master password and unwrap the VEK locally.
- Support per-item encryption for efficient rotation and sharing.
- Support device registration (public key) and device revocation.
- Implement secure sharing via X25519 ECDH (per-recipient wrapped keys).

---

## 2. High-level architecture and responsibilities
**Client responsibility (must be done by client, not server):**
- Derive MK from master password using Argon2id (client-side).
- Unwrap VEK using MK.
- Derive per-item encryption keys (IEK) from VEK and item id.
- Encrypt/decrypt items and attachments locally.
- Perform PAKE/OPAQUE/SRP for authentication or use an SRP-backed login flow.
- Generate device keypairs (X25519 for ECDH, Ed25519 for signatures) and keep the private key in OS keystore.

**Server responsibility (Gin backend):**
- Store encrypted blobs, metadata, device public keys, and sharing records.
- Verify device signatures on write/commit operations.
- Provide endpoints to register devices and exchange public keys.
- Provide sync endpoints and versioned storage.
- Do not attempt to decrypt vault contents at any time.

---

## 3. Crypto primitives & parameters
**Primitives (recommended)**
- KDF: `Argon2id` (from golang.org/x/crypto/argon2)
- Symmetric encryption: `AES-256-GCM` (crypto/aes + crypto/cipher)
- Key derivation: `HKDF-SHA256` (crypto/hkdf)
- Asymmetric: `X25519` (golang.org/x/crypto/curve25519) for shared secret; `Ed25519` (crypto/ed25519) for signatures.
- Randomness: `crypto/rand`.

**Suggested Argon2id parameters (configurable per environment):**
- Time: 3 (desktop), 2 (mobile) — make configurable.
- Memory: 64*1024 (64 MB) — make configurable.
- Parallelism: 2

**Notes:**
- Keep KDF params in DB so clients can fetch them on login (e.g., `users.kdf_params`).
- Use 96-bit IV (12 bytes) for AES-GCM and include an `aad` (vault metadata) where appropriate.

---

## 4. Database schema (SQL)
You already have `vaults` table. We need additional tables for wrapped keys, vault items (encrypted), attachments, devices, sharing records, and vault versions.

### New/modified tables
```sql
-- Vaults (existing, preserved)
CREATE TABLE IF NOT EXISTS vaults (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    icon VARCHAR(100),
    theme VARCHAR(50),
    is_favorite BOOLEAN NOT NULL DEFAULT false,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Per-vault wrapped key. Each vault has one wrapped VEK (Vault Encryption Key) per vault per user.
CREATE TABLE IF NOT EXISTS vault_keys (
    id SERIAL PRIMARY KEY,
    vault_id INTEGER NOT NULL REFERENCES vaults(id) ON DELETE CASCADE,
    wrapped_vek BYTEA NOT NULL,        -- ciphertext (wrapped by key derived from master password)
    wrap_iv BYTEA NOT NULL,            -- IV used for wrapping
    wrap_tag BYTEA NOT NULL,           -- tag for AES-GCM wrapping
    kdf_salt BYTEA NOT NULL,           -- salt used for Argon2 when creating the MK (store so client can derive MK)
    kdf_params JSONB NOT NULL,         -- stores time, memory, parallelism
    version INTEGER NOT NULL DEFAULT 1,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_vault_keys_vault_id ON vault_keys(vault_id);

-- Encrypted items per vault (per-item encryption)
CREATE TABLE IF NOT EXISTS vault_items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    vault_id INTEGER NOT NULL REFERENCES vaults(id) ON DELETE CASCADE,
    item_type VARCHAR(50) NOT NULL, -- login, note, card, etc
    encrypted_blob BYTEA NOT NULL,  -- AES-GCM ciphertext of plaintext JSON
    iv BYTEA NOT NULL,
    tag BYTEA NOT NULL,
    meta JSONB NULL,                -- non-sensitive searchable metadata (optional: keep minimal)
    version INTEGER NOT NULL DEFAULT 1,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_vault_items_vault_id ON vault_items(vault_id);

-- Attachments stored as encrypted objects; metadata tracked in DB.
CREATE TABLE IF NOT EXISTS vault_attachments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    vault_id INTEGER NOT NULL REFERENCES vaults(id) ON DELETE CASCADE,
    item_id UUID NULL REFERENCES vault_items(id) ON DELETE SET NULL,
    object_key TEXT NOT NULL,  -- path to object in object store
    size BIGINT NOT NULL,
    iv BYTEA NOT NULL,
    tag BYTEA NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Devices: store public keys and metadata
CREATE TABLE IF NOT EXISTS devices (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    device_label VARCHAR(255),
    x25519_public BYTEA NOT NULL,
    ed25519_public BYTEA NOT NULL,
    revoked_at TIMESTAMP NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    last_seen TIMESTAMP NULL
);

CREATE INDEX idx_devices_user_id ON devices(user_id);

-- Sharing records (encrypted wrapped keys per recipient)
CREATE TABLE IF NOT EXISTS sharing_records (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    vault_id INTEGER NOT NULL REFERENCES vaults(id) ON DELETE CASCADE,
    item_id UUID NULL REFERENCES vault_items(id) ON DELETE CASCADE,
    sender_user_id INTEGER NOT NULL REFERENCES users(id),
    recipient_user_id INTEGER NOT NULL REFERENCES users(id),
    wrapped_key BYTEA NOT NULL, -- key encrypted with ECDH-derived symmetric key
    wrap_iv BYTEA NOT NULL,
    wrap_tag BYTEA NOT NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'pending',
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Optional: vault_versions table to store versions/ETags of vault snapshots
CREATE TABLE IF NOT EXISTS vault_versions (
    id SERIAL PRIMARY KEY,
    vault_id INTEGER NOT NULL REFERENCES vaults(id) ON DELETE CASCADE,
    object_key TEXT NOT NULL, -- path to encrypted snapshot in object store
    mac BYTEA NULL,           -- optional authenticated mac over snapshot metadata
    created_by_device UUID NULL, -- device id
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);
```

**Notes**
- All ciphertext fields are `BYTEA` and should store raw bytes (not base64); clients can send base64 and server decodes.
- Keep `meta` minimal — never include secrets there. Prefer search or tags only if they are not sensitive.

---

## 5. API surface (Gin routes + handler descriptions)
Assume JWT-based session tokens for authenticated endpoints. All endpoints must require TLS and device signature for write endpoints.

### Auth & device
- `POST /auth/register` — create user account (store email, kdf params placeholder). *Client does NOT send master password.*
- `GET  /auth/kdf-params?email={email}` — return KDF params and salt for email/user so client can derive MK. (Alternatively fetch during login flow.)
- `POST /auth/login` — PAKE/SRP/OPAQUE flow; returns JWT and device provisioning challenge.
- `POST /devices/register` — register device public keys. Body includes device_id, x25519_public (base64), ed25519_public (base64), device_label; server returns challenge to sign for verification.
- `DELETE /devices/{id}` — revoke device (mark revoked_at).

### Vaults
- `POST /vaults` — create new vault (server stores metadata). Client must also `POST /vault_keys` to upload wrapped VEK for this vault.
- `GET /vaults` — list vaults (metadata only).
- `GET /vaults/{id}` — get vault metadata and latest version info (not decrypted content).
- `PUT /vaults/{id}` — update vault metadata (name, description). Must be device-signed.
- `DELETE /vaults/{id}` — delete vault (cascade removes items & keys).

### Vault keys
- `POST /vaults/{id}/keys` — upload wrapped VEK, kdf_salt, kdf_params. The server stores wrapped_vek as provided.
- `GET  /vaults/{id}/keys` — retrieve wrapped_vek, kdf_salt, kdf_params (for client to derive MK and unwrap). *Note:* ensure only owner can fetch.

### Items & attachments
- `POST /vaults/{id}/items` — push encrypted_item (encrypted_blob, iv, tag, meta). Must be signed by device ed25519 key.
- `GET  /vaults/{id}/items` — list item metadata (meta, id, version, timestamps).
- `GET  /vaults/{id}/items/{itemId}` — fetch encrypted_blob (ciphertext, iv, tag).
- `PUT  /vaults/{id}/items/{itemId}` — update encrypted item (new ciphertext, increment version). Device-signed.
- `DELETE /vaults/{id}/items/{itemId}` — delete item.

- `POST /vaults/{id}/attachments` — upload attachment metadata; actual binary upload goes to object store. Server returns pre-signed URL for object upload.
- `GET  /vaults/{id}/attachments/{id}` — download attachment (server redirects to presigned URL). Server must verify requester has access.

### Sync & versions
- `GET /sync/pull?vault_id={id}&since_version={n}` — return changed items and versions since `n`.
- `POST /sync/commit` — client commits a version snapshot; body contains `object_key` (where encrypted snapshot is stored) and device signature.

### Sharing
- `GET  /users/{userId}/public-keys` — fetch device public keys for ECDH.
- `POST /vaults/{id}/share` — create a sharing record for a recipient (sender provides wrapped_key and metadata). Server stores it and marks pending.
- `GET  /vaults/{id}/shares` — list sharing records for a user.
- `POST /vaults/{id}/share/accept` — recipient accepts share; server marks accepted.

---

## 6. Go implementation notes & code snippets
Use clear, well-tested helper functions encapsulating crypto. Always treat input as potentially tainted.

### Recommended packages
- `github.com/gin-gonic/gin` — HTTP
- `github.com/jackc/pgx/v5` or `database/sql` + `lib/pq` — Postgres
- `golang.org/x/crypto/argon2` — Argon2id
- `crypto/aes`, `crypto/cipher` — AES-GCM
- `crypto/hmac`, `crypto/sha256` — HMAC and HKDF
- `crypto/ed25519` — signatures
- `golang.org/x/crypto/curve25519` — X25519
- `crypto/rand` — CSPRNG
- `encoding/base64` — base64 encoding from clients

### Utility: derive master key (Argon2id)
```go
func DeriveMasterKey(password string, salt []byte, time, memory uint32, threads uint8, keyLen uint32) []byte {
    return argon2.IDKey([]byte(password), salt, time, memory, threads, keyLen)
}
```

### Utility: HKDF
```go
func HKDFExtractAndExpand(secret, salt, info []byte, outLen int) ([]byte, error) {
    hk := hkdf.New(sha256.New, secret, salt, info)
    out := make([]byte, outLen)
    if _, err := io.ReadFull(hk, out); err != nil {
        return nil, err
    }
    return out, nil
}
```

### AES-GCM decrypt & encrypt helpers
```go
func EncryptAESGCM(key, plaintext, aad []byte) (ciphertext, iv, tag []byte, err error) {
    block, err := aes.NewCipher(key)
    if err != nil { return nil, nil, nil, err }
    gcm, err := cipher.NewGCM(block)
    if err != nil { return nil, nil, nil, err }
    iv = make([]byte, gcm.NonceSize())
    if _, err := rand.Read(iv); err != nil { return nil, nil, nil, err }
    ct := gcm.Seal(nil, iv, plaintext, aad)
    // ct contains ciphertext||tag, but Go's Seal returns both; tag length = gcm.Overhead()
    return ct, iv, nil, nil
}

func DecryptAESGCM(key, iv, ciphertext, aad []byte) (plaintext []byte, err error) {
    block, err := aes.NewCipher(key)
    if err != nil { return nil, err }
    gcm, err := cipher.NewGCM(block)
    if err != nil { return nil, err }
    pt, err := gcm.Open(nil, iv, ciphertext, aad)
    if err != nil { return nil, err }
    return pt, nil
}
```
**Note:** Go's AES-GCM appends tag into ciphertext. You don't need a separate `tag` field if you store the full ciphertext. However the DB schema includes tag for clarity; store as combined or separate consistently.

### Unwrap VEK (client-side responsibility)
Server provides `wrapped_vek`, `wrap_iv`, `kdf_salt`, `kdf_params`.
Client does:
1. MK := Argon2id(masterPassword, kdf_salt, params)
2. wrappingKey := HKDF(MK, nil, []byte("wrap-vek"), 32)
3. VEK := AES-GCM-Decrypt(wrappingKey, wrap_iv, wrapped_vek, aad)

**Server** only stores the bytes; it must not attempt to decrypt.

### Verify device signatures for mutations
- All push/commit endpoints must include `device_id` and `signature` headers (signature is over the canonical request body + vault id + timestamp). Server fetches device's `ed25519_public` and verifies signature before persisting.

Snippet to verify Ed25519 signature:
```go
func VerifySignature(pubKey, message, sig []byte) bool {
    return ed25519.Verify(pubKey, message, sig)
}
```

### X25519 derive shared key for sharing
```go
func X25519SharedSecret(ourPriv, theirPub []byte) ([]byte, error) {
    var shared [32]byte
    var priv, pub [32]byte
    copy(priv[:], ourPriv)
    copy(pub[:], theirPub)
    curve25519.ScalarMult(&shared, &priv, &pub)
    // feed shared into HKDF to derive symmetric key
    symKey, _ := HKDFExtractAndExpand(shared[:], nil, []byte("share-vek"), 32)
    return symKey, nil
}
```

**Important**: Private keys must be generated and stored only on client devices (in keystore). Server stores only public keys.

---

## 7. Device registration and authentication
**Flow**
1. Client generates `x25519` keypair and `ed25519` keypair.
2. Client `POST /devices/register` with its public keys and device_label.
3. Server issues a per-device challenge (nonce) that client signs with Ed25519 and returns as proof-of-possession.
4. Server verifies signature and marks device as registered.

**Authentication**
- Use SRP or OPAQUE to avoid sending password-equivalent secrets to server. If SRP/OPAQUE is not implemented, at minimum use a secure login flow with Argon2id-salted verifier and short-lived tokens.
- After login, issue JWT access token (short-lived) and refresh token (rotate on use).

**Device revocation**
- `DELETE /devices/{id}` marks a device as revoked. Server refuses commits from revoked devices. Client should fetch device list and show revoked devices.

---

## 8. Sharing and team vaults
**One-off share (per-item)**
- Sender gets recipient's `x25519_public` from server.
- Sender computes ECDH shared key, derives symmetric key via HKDF.
- Sender wraps IEK or VEK with derived symmetric key (AES-GCM) and uploads a sharing record to `sharing_records` with `wrapped_key`.
- Recipient downloads the record, derives the same symmetric key using their private X25519 and sender's public X25519, un-wraps the IEK/VEK and can decrypt the item.

**Team vaults / shared vaults**
- Maintain a `team_vaults` ACL: store per-member wrapped VEKs (VEK encrypted under each member's public key) so each member can unwrap locally.
- Implement roles (owner/editor/viewer) and support rotating team VEK to revoke access.

**Revocation**
- To revoke, rotate VEK and re-wrap for remaining members (lazy re-encryption of items is acceptable but note complexity).

---

## 9. Sync, versioning, and conflict resolution
**Versioning**
- Each `vault_items` row has a `version` integer.
- Clients should provide `base_version` when updating; server must reject updates if `base_version` != current `version` (optimistic concurrency). On conflict, server returns 409 with latest item data.

**Conflict resolution**
- Simple approach: last-writer-wins with server timestamp. Preferred: return conflict to client and let client merge.

**Snapshotting**
- Clients may upload a full encrypted snapshot to object storage and POST a `sync/commit` with `object_key` and device signature. Server will record snapshot in `vault_versions`.

---

## 10. Migration and backwards compatibility
- Add `version` fields to vault/vault_keys schema to support changes in wrapping scheme.
- Provide migration scripts in `/migrations` for goose (up/down).
- When changing KDF params or wrapping scheme, increment `vault_keys.version` and include migration path in clients.

---

## 11. Testing plan
**Unit tests**
- KDF parameter tests (consistency across runs).
- AES-GCM encryption/decryption for various plaintext sizes.
- HKDF derivation tests.
- ECDH handshake derivation correctness.

**Integration tests**
- Device registration flow (register→challenge→verify)
- Push/pull sync workflows across two simulated clients.
- Sharing flow with two simulated users/devices.

**Security tests**
- Fuzz vault item parsing and DB insertion points.
- Test server rejects malformed base64/ciphertext.
- Penetration tests focusing on server endpoints.

---

## 12. Security checklist & operational notes
- All endpoints require TLS 1.3 with HSTS.
- Enforce rate-limiting on auth endpoints and device registration.
- Validate input sizes and types; do not trust client-side validation.
- Log audit events for device registration, revocation, large data downloads, and repeated failed decrypt attempts (counts only). Avoid logging ciphertext or keys.
- Use signed client binaries or produce release signing to mitigate supply-chain risk.
- Run regular third-party crypto and security audits.

---

## 13. Environment / configuration
Environment variables to configure the server:
```
DATABASE_URL=postgres://...
S3_BUCKET=...
S3_REGION=...
KDF_TIME=3
KDF_MEMORY=65536
KDF_PARALLELISM=2
JWT_SECRET=...
JWT_EXPIRY=15m
REFRESH_TOKEN_SECRET=...
RATE_LIMIT=100
```

All KDF params should be stored with user/vault metadata so clients can derive key with the exact params used for creation.

---

## 14. Deliverables & acceptance criteria
- DB migrations for `vault_keys`, `vault_items`, `vault_attachments`, `devices`, `sharing_records`, `vault_versions` using goose.
- Gin handlers for all endpoints listed in **API surface** with unit/integration tests.
- Helper crypto library in Go (`/internal/crypto`) that exposes: `DeriveMasterKey`, `HKDF`, `EncryptAESGCM`, `DecryptAESGCM`, `VerifySignature`, `X25519SharedSecret`.
- Device registration flow (challenge/sign/verify) implemented and tested.
- Sync endpoints and optimistic concurrency controls.
- Example Postman/Insomnia collection showing usage flows.
- README describing how to run locally including env variables and migrations.

---

### Appendix A — example request/response payloads
**Upload a wrapped VEK (client -> server)**
```json
POST /vaults/123/keys
Authorization: Bearer <token>
Content-Type: application/json
{
  "wrapped_vek": "BASE64(...)",
  "wrap_iv": "BASE64(...)",
  "kdf_salt": "BASE64(...)",
  "kdf_params": {"time":3,"memory":65536,"parallelism":2},
  "version": 1
}
```

**Push encrypted item (client -> server)**
```json
POST /vaults/123/items
Authorization: Bearer <token>
X-Device-Id: <device-id>
X-Device-Sig: BASE64(signature)
Content-Type: application/json
{
  "id": "uuid",
  "item_type": "login",
  "encrypted_blob": "BASE64(...)",
  "iv": "BASE64(...)",
  "meta": {"title":"example"},
  "version": 1
}
```

**Fetch wrapped VEK (server -> client)**
```json
GET /vaults/123/keys
Authorization: Bearer <token>

200 OK
{
  "wrapped_vek": "BASE64(...)",
  "wrap_iv": "BASE64(...)",
  "kdf_salt": "BASE64(...)",
  "kdf_params": {"time":3,"memory":65536,"parallelism":2},
  "version": 1
}
```


