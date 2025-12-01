package crypto

import (
	"encoding/base64"
	"fmt"
)

// EncodeBase64 encodes bytes to base64 string (standard encoding)
func EncodeBase64(data []byte) string {
	return base64.StdEncoding.EncodeToString(data)
}

// DecodeBase64 decodes base64 string to bytes (standard encoding)
func DecodeBase64(encoded string) ([]byte, error) {
	decoded, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return nil, fmt.Errorf("failed to decode base64: %w", err)
	}
	return decoded, nil
}

// EncodeBase64URL encodes bytes to base64 URL-safe string
func EncodeBase64URL(data []byte) string {
	return base64.URLEncoding.EncodeToString(data)
}

// DecodeBase64URL decodes base64 URL-safe string to bytes
func DecodeBase64URL(encoded string) ([]byte, error) {
	decoded, err := base64.URLEncoding.DecodeString(encoded)
	if err != nil {
		return nil, fmt.Errorf("failed to decode base64 URL: %w", err)
	}
	return decoded, nil
}

// EncodeBase64RawURL encodes bytes to base64 raw URL-safe string (no padding)
func EncodeBase64RawURL(data []byte) string {
	return base64.RawURLEncoding.EncodeToString(data)
}

// DecodeBase64RawURL decodes base64 raw URL-safe string to bytes (no padding)
func DecodeBase64RawURL(encoded string) ([]byte, error) {
	decoded, err := base64.RawURLEncoding.DecodeString(encoded)
	if err != nil {
		return nil, fmt.Errorf("failed to decode base64 raw URL: %w", err)
	}
	return decoded, nil
}
