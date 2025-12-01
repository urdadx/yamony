# Yamony - Zero-Knowledge Password Manager

Yamony is a modern password manager built with end-to-end encryption and zero-knowledge architecture. The server stores only encrypted data and cannot decrypt user information, ensuring maximum privacy and security.

## Features

- 🔐 **Zero-Knowledge Encryption**: All data encrypted client-side before transmission
- 🔑 **Advanced Key Management**: Per-vault encryption keys (VEK) with per-item keys (IEK)
- 🚀 **Device-Based Authentication**: Ed25519 signatures for all write operations
- 🤝 **Secure Sharing**: X25519 ECDH key exchange for sharing vaults with other users
- 🔄 **Sync & Versioning**: Optimistic concurrency control with ETags and conflict detection
- 🛡️ **Strong Cryptography**: Argon2id, AES-256-GCM, HKDF-SHA256, Ed25519, X25519

## Architecture

### Key Hierarchy

```
Master Password → Argon2id → Master Key (MK)
                                  ↓
                          HKDF("wrap-vek")
                                  ↓
                           Wrapping Key
                                  ↓
                          AES-256-GCM encrypt
                                  ↓
            Vault Encryption Key (VEK) [stored wrapped]
                                  ↓
                     HKDF("item-key:{item_id}")
                                  ↓
                     Item Encryption Key (IEK)
                                  ↓
                          AES-256-GCM encrypt
                                  ↓
                    Item Data [stored encrypted]
```

### Security Features

- **Client-Side Encryption**: All sensitive data encrypted before leaving the device
- **Device Authentication**: Ed25519 signatures prevent unauthorized access
- **Perfect Forward Secrecy**: Each item uses a unique encryption key
- **Secure Sharing**: ECDH ensures only intended recipients can decrypt shared data
- **Optimistic Concurrency**: Version control prevents data loss during conflicts
- **Audit Trail**: Complete version history for all changes

## Getting Started

### Prerequisites

- Go 1.21 or higher
- PostgreSQL 14+
- Make (for using Makefile commands)
- [sqlc](https://sqlc.dev/) (for generating type-safe SQL code)
- [goose](https://github.com/pressly/goose) (for database migrations)

### Installation

1. Clone the repository:
```bash
git clone https://github.com/urdadx/yamony.git
cd yamony
```

2. Install dependencies:
```bash
go mod download
```

3. Set up environment variables (create `.env` file):
```env
DATABASE_URL=postgres://user:password@localhost:5432/yamony?sslmode=disable
PORT=8080
SESSION_SECRET=your-secret-key-change-in-production
```

4. Start PostgreSQL with Docker:
```bash
make docker-run
```

5. Run database migrations:
```bash
goose -dir internal/database/schema postgres "${DATABASE_URL}" up
```

6. Generate SQL code:
```bash
sqlc generate
```

7. Build and run the application:
```bash
make build
make run
```

The API will be available at `http://localhost:8080`

### Quick Start with Make

The project includes a Makefile for common tasks:

## MakeFile

Run build make command with tests
```bash
make all
```

Build the application
```bash
make build
```

Run the application
```bash
make run
```
Create DB container
```bash
make docker-run
```

Shutdown DB Container
```bash
make docker-down
```

DB Integrations Test:
```bash
make itest
```

Live reload the application:
```bash
make watch
```

Run the test suite:
```bash
make test
```

Clean up binary from the last build:
```bash
make clean
```

## API Documentation

See [API_DOCS.md](./API_DOCS.md) for complete API documentation including:

- Authentication endpoints
- Device management
- Vault key operations
- Vault item CRUD
- Sharing workflows
- Synchronization API
- Cryptographic details
- Example code snippets

## Project Structure

```
yamony/
├── cmd/
│   └── api/
│       └── main.go           # Application entry point
├── internal/
│   ├── crypto/               # Cryptographic primitives
│   │   ├── argon2.go         # Argon2id KDF
│   │   ├── aes.go            # AES-256-GCM encryption
│   │   ├── hkdf.go           # HKDF key derivation
│   │   ├── ed25519.go        # Ed25519 signatures
│   │   ├── x25519.go         # X25519 ECDH
│   │   ├── helpers.go        # High-level crypto helpers
│   │   └── crypto_test.go    # Comprehensive tests
│   ├── database/
│   │   ├── schema/           # SQL migration files
│   │   ├── queries/          # SQL queries for sqlc
│   │   └── sqlc/             # Generated type-safe code
│   └── server/
│       ├── handlers/         # HTTP request handlers
│       │   ├── device_handler.go
│       │   ├── vault_key_handler.go
│       │   ├── vault_item_handler.go
│       │   ├── share_handler.go
│       │   └── sync_handler.go
│       ├── middleware/       # Authentication middleware
│       ├── services/         # Business logic layer
│       └── routes.go         # Route registration
├── frontend/                 # React frontend (TanStack Router)
├── API_AUTH.md              # Authentication documentation
├── API_DOCS.md              # Complete API documentation
├── docs/
│   └── implementation.md    # Zero-knowledge implementation details
└── README.md
```

## Development

### Database Migrations

Create a new migration:
```bash
goose -dir internal/database/schema create migration_name sql
```

Run migrations:
```bash
goose -dir internal/database/schema postgres "${DATABASE_URL}" up
```

Rollback migration:
```bash
goose -dir internal/database/schema postgres "${DATABASE_URL}" down
```

### Generating SQL Code

After modifying SQL queries in `internal/database/queries/`, regenerate Go code:
```bash
sqlc generate
```

### Running Tests

Run all tests:
```bash
make test
```

Run crypto tests:
```bash
go test ./internal/crypto -v
```

Run integration tests:
```bash
make itest
```

### Live Reload

For development with automatic reload on file changes:
```bash
make watch
```

## Testing

### Crypto Package Tests

The crypto package includes comprehensive tests covering:
- Argon2id key derivation
- AES-256-GCM encryption/decryption
- HKDF key derivation
- Ed25519 signature verification
- X25519 key exchange
- VaultKeyWrapper operations
- ItemEncryptor operations
- ShareKeyWrapper operations

Run crypto tests:
```bash
go test ./internal/crypto -v -cover
```

### Integration Tests

Integration tests cover database operations and handler logic. Ensure PostgreSQL is running before executing:
```bash
make itest
```

## Deployment

### Environment Variables

Required environment variables for production:

```env
# Database
DATABASE_URL=postgres://user:password@host:5432/yamony?sslmode=require

# Server
PORT=8080
GIN_MODE=release

# Security
SESSION_SECRET=cryptographically-secure-random-key

# CORS (adjust for your domain)
ALLOWED_ORIGINS=https://yourdomain.com

# Optional
LOG_LEVEL=info
```

### Docker Deployment

Build Docker image:
```bash
docker build -t yamony:latest .
```

Run with Docker Compose:
```bash
docker-compose up -d
```

### Security Checklist

Before deploying to production:

- [ ] Generate strong SESSION_SECRET (32+ random bytes)
- [ ] Enable HTTPS/TLS for all connections
- [ ] Configure CORS with specific allowed origins
- [ ] Set `Secure` flag on session cookies
- [ ] Enable PostgreSQL SSL mode
- [ ] Set up database backups
- [ ] Configure rate limiting
- [ ] Enable audit logging
- [ ] Review and update ALLOWED_ORIGINS
- [ ] Set GIN_MODE=release

## Client Implementation

For building a client application, see the [API_DOCS.md](./API_DOCS.md) for:

- Complete endpoint documentation
- Cryptographic operation details
- Example code in JavaScript and Go
- Authentication flow
- Device registration process
- Key derivation procedures
- Encryption/decryption workflows

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## Security

### Reporting Vulnerabilities

If you discover a security vulnerability, please email security@yamony.com instead of using the issue tracker.

### Security Considerations

- All encryption happens client-side
- Server never has access to plaintext data or encryption keys
- Device private keys must be stored securely on client devices
- Master password is never transmitted to the server
- Consider implementing additional security measures like:
  - Two-factor authentication (2FA)
  - Biometric authentication
  - Hardware security keys (WebAuthn)
  - Rate limiting on API endpoints
  - Account recovery mechanisms

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Acknowledgments

- Built with [Gin](https://gin-gonic.com/) web framework
- Database migrations with [goose](https://github.com/pressly/goose)
- Type-safe SQL with [sqlc](https://sqlc.dev/)
- Cryptography powered by Go's crypto packages
- Frontend built with React and TanStack Router

## Tech Stack

### Backend
- **Language**: Go 1.21+
- **Web Framework**: Gin
- **Database**: PostgreSQL 14+ with pgx driver
- **SQL Code Generation**: sqlc
- **Migrations**: goose
- **Cryptography**: golang.org/x/crypto

### Frontend
- **Framework**: React
- **Router**: TanStack Router
- **Build Tool**: Vite
- **Styling**: Tailwind CSS

### DevOps
- **Containerization**: Docker
- **Orchestration**: Docker Compose
- **CI/CD**: GitHub Actions (recommended)

## Roadmap

- [ ] Mobile app (iOS/Android)
- [ ] Browser extensions (Chrome, Firefox, Safari)
- [ ] Two-factor authentication (TOTP)
- [ ] Hardware security key support (WebAuthn)
- [ ] Emergency access/account recovery
- [ ] Password breach monitoring
- [ ] Password strength analyzer
- [ ] Secure notes and file attachments
- [ ] Audit logs and activity monitoring
- [ ] Team/organization features
- [ ] Admin dashboard

## Support

For questions, issues, or suggestions:

- Open an issue on GitHub
- Check existing documentation in `docs/`
- Review API documentation in `API_DOCS.md`
- Contact: support@yamony.com

---

