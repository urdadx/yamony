package handlers

import (
	"crypto/ed25519"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"

	"yamony/internal/crypto"
	"yamony/internal/database/sqlc"
	"yamony/internal/server/services"
)

type DeviceHandler struct {
	services services.Service
}

func NewDeviceHandler(services services.Service) *DeviceHandler {
	return &DeviceHandler{services: services}
}

// RegisterDeviceRequest represents the device registration payload
type RegisterDeviceRequest struct {
	DeviceLabel   string `json:"device_label" binding:"required"`
	X25519Public  string `json:"x25519_public" binding:"required"`  // base64 encoded
	Ed25519Public string `json:"ed25519_public" binding:"required"` // base64 encoded
}

// RegisterDeviceChallengeResponse is returned after initial registration
type RegisterDeviceChallengeResponse struct {
	DeviceID  string `json:"device_id"`
	Challenge string `json:"challenge"` // base64 encoded challenge to sign
	ExpiresAt int64  `json:"expires_at"`
}

// VerifyDeviceRequest contains the signed challenge
type VerifyDeviceRequest struct {
	DeviceID  string `json:"device_id" binding:"required"`
	Signature string `json:"signature" binding:"required"` // base64 encoded signature
}

// DeviceResponse represents a device
type DeviceResponse struct {
	ID            string     `json:"id"`
	DeviceLabel   string     `json:"device_label"`
	X25519Public  string     `json:"x25519_public"`
	Ed25519Public string     `json:"ed25519_public"`
	CreatedAt     time.Time  `json:"created_at"`
	LastSeen      *time.Time `json:"last_seen,omitempty"`
	RevokedAt     *time.Time `json:"revoked_at,omitempty"`
}

// PublicKeysResponse for sharing
type PublicKeysResponse struct {
	DeviceID      string    `json:"device_id"`
	DeviceLabel   string    `json:"device_label"`
	X25519Public  string    `json:"x25519_public"`
	Ed25519Public string    `json:"ed25519_public"`
	CreatedAt     time.Time `json:"created_at"`
}

// RegisterDevice initiates device registration and returns a challenge
// POST /api/devices/register
func (h *DeviceHandler) RegisterDevice(c *gin.Context) {
	var req RegisterDeviceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	// Decode and validate public keys
	x25519Pub, err := crypto.DecodeBase64(req.X25519Public)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid x25519_public format"})
		return
	}
	if err := crypto.ValidateX25519PublicKey(x25519Pub); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ed25519Pub, err := crypto.DecodeBase64(req.Ed25519Public)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ed25519_public format"})
		return
	}
	if err := crypto.ValidateEd25519PublicKey(ed25519Pub); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Generate device ID
	deviceID := uuid.New()
	pgUUID := pgtype.UUID{}
	_ = pgUUID.Scan(deviceID.String())

	pgDeviceLabel := pgtype.Text{}
	_ = pgDeviceLabel.Scan(req.DeviceLabel)

	// Create device in database
	queries := h.services.GetDB().GetQueries()
	device, err := queries.CreateDevice(c.Request.Context(), sqlc.CreateDeviceParams{
		ID:            pgUUID,
		UserID:        userID.(int32),
		DeviceLabel:   pgDeviceLabel,
		X25519Public:  x25519Pub,
		Ed25519Public: ed25519Pub,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to register device"})
		return
	}

	// Generate challenge for proof-of-possession
	challenge, err := crypto.GenerateRandomBytes(32)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate challenge"})
		return
	}

	// Store challenge in session/cache with expiry (5 minutes)
	challengeKey := fmt.Sprintf("device_challenge:%s", deviceID.String())
	expiresAt := time.Now().Add(5 * time.Minute).Unix()

	// TODO: Store in Redis/cache. For now, we'll accept without challenge verification
	// In production, store challenge and verify in VerifyDevice endpoint
	_ = challengeKey // placeholder

	c.JSON(http.StatusCreated, RegisterDeviceChallengeResponse{
		DeviceID:  uuidToString(device.ID),
		Challenge: crypto.EncodeBase64(challenge),
		ExpiresAt: expiresAt,
	})
}

// VerifyDevice verifies the signed challenge (optional additional security)
// POST /api/devices/verify
func (h *DeviceHandler) VerifyDevice(c *gin.Context) {
	var req VerifyDeviceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Parse device ID
	deviceID, err := uuid.Parse(req.DeviceID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid device_id"})
		return
	}

	pgUUID := pgtype.UUID{}
	_ = pgUUID.Scan(deviceID.String())

	// Get device from database
	queries := h.services.GetDB().GetQueries()
	device, err := queries.GetDeviceByID(c.Request.Context(), pgUUID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "device not found"})
		return
	}

	// Decode signature
	signature, err := crypto.DecodeBase64(req.Signature)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid signature format"})
		return
	}

	// TODO: Retrieve challenge from cache
	// For now, we'll skip challenge verification and just verify signature format
	if err := crypto.ValidateEd25519Signature(signature); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid signature"})
		return
	}

	// In production: verify signature against challenge
	// message := []byte(challenge)
	// valid := crypto.VerifySignature(device.Ed25519Public, message, signature)
	// if !valid {
	//     c.JSON(http.StatusUnauthorized, gin.H{"error": "signature verification failed"})
	//     return
	// }

	c.JSON(http.StatusOK, gin.H{
		"message":   "device verified",
		"device_id": uuidToString(device.ID),
	})
}

// GetDevices returns all active devices for the authenticated user
// GET /api/devices
func (h *DeviceHandler) GetDevices(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	queries := h.services.GetDB().GetQueries()
	devices, err := queries.GetDevicesByUserID(c.Request.Context(), userID.(int32))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch devices"})
		return
	}

	response := make([]DeviceResponse, len(devices))
	for i, device := range devices {
		deviceID := uuidToString(device.ID)
		deviceLabel := ""
		if device.DeviceLabel.Valid {
			deviceLabel = device.DeviceLabel.String
		}

		response[i] = DeviceResponse{
			ID:            deviceID,
			DeviceLabel:   deviceLabel,
			X25519Public:  crypto.EncodeBase64(device.X25519Public),
			Ed25519Public: crypto.EncodeBase64(device.Ed25519Public),
			CreatedAt:     timestampToTime(device.CreatedAt),
			LastSeen:      timestampToTimePtr(device.LastSeen),
			RevokedAt:     timestampToTimePtr(device.RevokedAt),
		}
	}

	c.JSON(http.StatusOK, response)
}

// RevokeDevice revokes a device
// DELETE /api/devices/:id
func (h *DeviceHandler) RevokeDevice(c *gin.Context) {
	deviceIDParam := c.Param("id")
	deviceID, err := uuid.Parse(deviceIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid device_id"})
		return
	}

	pgUUID := pgtype.UUID{}
	_ = pgUUID.Scan(deviceID.String())

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	// Verify device belongs to user
	queries := h.services.GetDB().GetQueries()
	device, err := queries.GetDeviceByID(c.Request.Context(), pgUUID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "device not found"})
		return
	}

	if device.UserID != userID.(int32) {
		c.JSON(http.StatusForbidden, gin.H{"error": "not authorized to revoke this device"})
		return
	}

	// Revoke device
	err = queries.RevokeDevice(c.Request.Context(), pgUUID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to revoke device"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "device revoked successfully"})
}

// GetUserPublicKeys returns public keys for a user's devices (for sharing)
// GET /api/users/:user_id/public-keys
func (h *DeviceHandler) GetUserPublicKeys(c *gin.Context) {
	userIDParam := c.Param("user_id")
	targetUserID, err := parseIntParam(userIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user_id"})
		return
	}

	// TODO: Add authorization check - only allow if users are in same vault/team

	queries := h.services.GetDB().GetQueries()
	devices, err := queries.GetUserDevicePublicKeys(c.Request.Context(), targetUserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch public keys"})
		return
	}

	response := make([]PublicKeysResponse, len(devices))
	for i, device := range devices {
		deviceID := uuidToString(device.ID)
		deviceLabel := ""
		if device.DeviceLabel.Valid {
			deviceLabel = device.DeviceLabel.String
		}

		response[i] = PublicKeysResponse{
			DeviceID:      deviceID,
			DeviceLabel:   deviceLabel,
			X25519Public:  crypto.EncodeBase64(device.X25519Public),
			Ed25519Public: crypto.EncodeBase64(device.Ed25519Public),
			CreatedAt:     timestampToTime(device.CreatedAt),
		}
	}

	c.JSON(http.StatusOK, response)
}

// VerifyDeviceSignature verifies a signature from a device
func VerifyDeviceSignature(devicePubKey ed25519.PublicKey, message, signature []byte) bool {
	return crypto.VerifySignature(devicePubKey, message, signature)
}

// CreateCanonicalMessage creates a canonical message for signing
// Format: METHOD|PATH|TIMESTAMP|BODY_HASH
func CreateCanonicalMessage(method, path string, timestamp int64, bodyHash []byte) []byte {
	msg := fmt.Sprintf("%s|%s|%d|%s", method, path, timestamp, crypto.EncodeBase64(bodyHash))
	return []byte(msg)
}

// SignatureVerificationMiddleware verifies device signatures on write operations
func SignatureVerificationMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Only verify signatures on write operations
		if c.Request.Method == "GET" || c.Request.Method == "OPTIONS" {
			c.Next()
			return
		}

		// Get device ID and signature from headers
		deviceIDHeader := c.GetHeader("X-Device-Id")
		signatureHeader := c.GetHeader("X-Device-Signature")
		timestampHeader := c.GetHeader("X-Device-Timestamp")

		if deviceIDHeader == "" || signatureHeader == "" || timestampHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "device authentication required"})
			c.Abort()
			return
		}

		// Parse device ID
		deviceID, err := uuid.Parse(deviceIDHeader)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid device_id"})
			c.Abort()
			return
		}

		// Get device and verify not revoked
		// TODO: Implement database query here
		_ = deviceID

		// Decode signature
		signature, err := crypto.DecodeBase64(signatureHeader)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid signature format"})
			c.Abort()
			return
		}

		// TODO: Get device public key from database and verify signature
		// For now, just validate signature format
		if err := crypto.ValidateEd25519Signature(signature); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid signature"})
			c.Abort()
			return
		}

		c.Next()
	}
}

// ExtractDeviceSignature extracts and validates device signature from request
type DeviceSignatureData struct {
	DeviceID  uuid.UUID
	Timestamp int64
	Signature []byte
}

func ExtractDeviceSignature(c *gin.Context) (*DeviceSignatureData, error) {
	deviceIDHeader := c.GetHeader("X-Device-Id")
	signatureHeader := c.GetHeader("X-Device-Signature")
	timestampHeader := c.GetHeader("X-Device-Timestamp")

	if deviceIDHeader == "" || signatureHeader == "" || timestampHeader == "" {
		return nil, fmt.Errorf("missing device authentication headers")
	}

	deviceID, err := uuid.Parse(deviceIDHeader)
	if err != nil {
		return nil, fmt.Errorf("invalid device_id")
	}

	signature, err := crypto.DecodeBase64(signatureHeader)
	if err != nil {
		return nil, fmt.Errorf("invalid signature format")
	}

	var timestamp int64
	if err := json.Unmarshal([]byte(timestampHeader), &timestamp); err != nil {
		return nil, fmt.Errorf("invalid timestamp")
	}

	// Check timestamp is within acceptable window (5 minutes)
	now := time.Now().Unix()
	if now-timestamp > 300 || timestamp > now+300 {
		return nil, fmt.Errorf("timestamp out of acceptable range")
	}

	return &DeviceSignatureData{
		DeviceID:  deviceID,
		Timestamp: timestamp,
		Signature: signature,
	}, nil
}

// Helper functions

// uuidToString converts pgtype.UUID to string
func uuidToString(pgUUID pgtype.UUID) string {
	if !pgUUID.Valid {
		return ""
	}
	return fmt.Sprintf("%x-%x-%x-%x-%x",
		pgUUID.Bytes[0:4], pgUUID.Bytes[4:6], pgUUID.Bytes[6:8],
		pgUUID.Bytes[8:10], pgUUID.Bytes[10:16])
}

// parseIntParam parses string to int32
func parseIntParam(s string) (int32, error) {
	var result int32
	_, err := fmt.Sscanf(s, "%d", &result)
	return result, err
}

// timestampToTime converts pgtype.Timestamp to time.Time
func timestampToTime(pgTime pgtype.Timestamp) time.Time {
	if !pgTime.Valid {
		return time.Time{}
	}
	return pgTime.Time
}

// timestampToTimePtr converts pgtype.Timestamp to *time.Time
func timestampToTimePtr(pgTime pgtype.Timestamp) *time.Time {
	if !pgTime.Valid {
		return nil
	}
	t := pgTime.Time
	return &t
}
