package crypto

import (
	"crypto/rand"
	"encoding/json"
	"fmt"

	"golang.org/x/crypto/argon2"
)

// KDFParams holds the parameters for Argon2id key derivation
type KDFParams struct {
	Time        uint32 `json:"time"`
	Memory      uint32 `json:"memory"`
	Parallelism uint8  `json:"parallelism"`
	KeyLen      uint32 `json:"keyLen"`
}

// DefaultKDFParams returns secure default parameters for Argon2id
func DefaultKDFParams() KDFParams {
	return KDFParams{
		Time:        3,
		Memory:      64 * 1024, // 64 MB
		Parallelism: 2,
		KeyLen:      32, // 256 bits
	}
}

// MobileKDFParams returns lighter parameters suitable for mobile devices
func MobileKDFParams() KDFParams {
	return KDFParams{
		Time:        2,
		Memory:      32 * 1024, // 32 MB
		Parallelism: 2,
		KeyLen:      32,
	}
}

// GenerateSalt creates a cryptographically secure random salt
func GenerateSalt(length int) ([]byte, error) {
	if length <= 0 {
		length = 32 // default 256 bits
	}
	salt := make([]byte, length)
	if _, err := rand.Read(salt); err != nil {
		return nil, fmt.Errorf("failed to generate salt: %w", err)
	}
	return salt, nil
}

// DeriveMasterKey derives a key from a password using Argon2id
func DeriveMasterKey(password string, salt []byte, params KDFParams) []byte {
	return argon2.IDKey(
		[]byte(password),
		salt,
		params.Time,
		params.Memory,
		params.Parallelism,
		params.KeyLen,
	)
}

// MarshalKDFParams converts KDFParams to JSON bytes
func MarshalKDFParams(params KDFParams) ([]byte, error) {
	return json.Marshal(params)
}

// UnmarshalKDFParams parses JSON bytes into KDFParams
func UnmarshalKDFParams(data []byte) (KDFParams, error) {
	var params KDFParams
	err := json.Unmarshal(data, &params)
	return params, err
}

// ValidateKDFParams checks if the parameters are within safe ranges
func ValidateKDFParams(params KDFParams) error {
	if params.Time < 1 {
		return fmt.Errorf("time parameter must be at least 1")
	}
	if params.Memory < 8*1024 {
		return fmt.Errorf("memory parameter must be at least 8 MB")
	}
	if params.Parallelism < 1 {
		return fmt.Errorf("parallelism must be at least 1")
	}
	if params.KeyLen < 16 {
		return fmt.Errorf("key length must be at least 16 bytes")
	}
	return nil
}
