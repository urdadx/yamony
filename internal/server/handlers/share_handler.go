package handlers

import (
	"crypto/sha256"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"

	"yamony/internal/crypto"
	"yamony/internal/database/sqlc"
	"yamony/internal/server/services"
)

type ShareHandler struct {
	services services.Service
}

func NewShareHandler(services services.Service) *ShareHandler {
	return &ShareHandler{services: services}
}

// ShareVaultRequest represents a request to share a vault with another user
type ShareVaultRequest struct {
	RecipientUserID int32  `json:"recipient_user_id" binding:"required"`
	WrappedVEK      string `json:"wrapped_vek" binding:"required"` // base64 encoded VEK wrapped with ECDH
	WrapIV          string `json:"wrap_iv" binding:"required"`     // base64 encoded
	WrapTag         string `json:"wrap_tag" binding:"required"`    // base64 encoded
}

// ShareItemRequest represents a request to share a specific vault item
type ShareItemRequest struct {
	RecipientUserID int32  `json:"recipient_user_id" binding:"required"`
	ItemID          string `json:"item_id" binding:"required"`
	WrappedIEK      string `json:"wrapped_iek" binding:"required"` // base64 encoded IEK wrapped with ECDH
	WrapIV          string `json:"wrap_iv" binding:"required"`     // base64 encoded
	WrapTag         string `json:"wrap_tag" binding:"required"`    // base64 encoded
}

// SharingRecordResponse represents a sharing record
type SharingRecordResponse struct {
	ID              string  `json:"id"`
	VaultID         int32   `json:"vault_id"`
	ItemID          *string `json:"item_id,omitempty"`
	SenderUserID    int32   `json:"sender_user_id"`
	RecipientUserID int32   `json:"recipient_user_id"`
	WrappedKey      string  `json:"wrapped_key"`
	WrapIV          string  `json:"wrap_iv"`
	WrapTag         string  `json:"wrap_tag"`
	Status          string  `json:"status"`
	CreatedAt       string  `json:"created_at"`
	AcceptedAt      *string `json:"accepted_at,omitempty"`
}

// SharedVaultResponse represents a vault shared with the user
type SharedVaultResponse struct {
	ID        int32  `json:"id"`
	Name      string `json:"name"`
	UserID    int32  `json:"user_id"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

// ShareVault shares a vault with another user
// POST /api/vaults/:id/share
func (h *ShareHandler) ShareVault(c *gin.Context) {
	vaultID, err := parseIntParam(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid vault_id"})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	var req ShareVaultRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Verify device signature
	sigData, err := ExtractDeviceSignature(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "device authentication required: " + err.Error()})
		return
	}

	queries := h.services.GetDB().GetQueries()

	// Verify vault belongs to user
	vault, err := queries.GetVaultByID(c.Request.Context(), vaultID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "vault not found"})
		return
	}

	if vault.UserID != userID.(int32) {
		c.JSON(http.StatusForbidden, gin.H{"error": "only vault owner can share"})
		return
	}

	// Verify device
	pgDeviceID := pgtype.UUID{}
	_ = pgDeviceID.Scan(sigData.DeviceID.String())

	device, err := queries.GetDeviceByID(c.Request.Context(), pgDeviceID)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "device not found"})
		return
	}

	if device.UserID != userID.(int32) {
		c.JSON(http.StatusForbidden, gin.H{"error": "device does not belong to user"})
		return
	}

	// Verify signature
	bodyBytes, _ := c.Get("raw_body")
	bodyHash := sha256.Sum256(bodyBytes.([]byte))
	canonicalMsg := CreateCanonicalMessage(c.Request.Method, c.Request.URL.Path, sigData.Timestamp, bodyHash[:])

	if !VerifyDeviceSignature(device.Ed25519Public, canonicalMsg, sigData.Signature) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid device signature"})
		return
	}

	// Verify recipient user exists (get their public keys)
	recipientKeys, err := queries.GetUserDevicePublicKeys(c.Request.Context(), req.RecipientUserID)
	if err != nil || len(recipientKeys) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "recipient user not found or has no devices"})
		return
	}

	// Decode wrapped key
	wrappedKey, err := crypto.DecodeBase64(req.WrappedVEK)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid wrapped_vek format"})
		return
	}

	wrapIV, err := crypto.DecodeBase64(req.WrapIV)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid wrap_iv format"})
		return
	}

	wrapTag, err := crypto.DecodeBase64(req.WrapTag)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid wrap_tag format"})
		return
	}

	// Create sharing record
	sharingRecord, err := queries.CreateSharingRecord(c.Request.Context(), sqlc.CreateSharingRecordParams{
		VaultID:         vaultID,
		ItemID:          pgtype.UUID{Valid: false}, // NULL for vault sharing
		SenderUserID:    userID.(int32),
		RecipientUserID: req.RecipientUserID,
		WrappedKey:      wrappedKey,
		WrapIv:          wrapIV,
		WrapTag:         wrapTag,
		Status:          "pending",
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create sharing record"})
		return
	}

	response := SharingRecordResponse{
		ID:              uuidToString(sharingRecord.ID),
		VaultID:         sharingRecord.VaultID,
		ItemID:          uuidToStringPtr(sharingRecord.ItemID),
		SenderUserID:    sharingRecord.SenderUserID,
		RecipientUserID: sharingRecord.RecipientUserID,
		WrappedKey:      crypto.EncodeBase64(sharingRecord.WrappedKey),
		WrapIV:          crypto.EncodeBase64(sharingRecord.WrapIv),
		WrapTag:         crypto.EncodeBase64(sharingRecord.WrapTag),
		Status:          sharingRecord.Status,
		CreatedAt:       timestampToTime(sharingRecord.CreatedAt).Format("2006-01-02T15:04:05Z07:00"),
		AcceptedAt:      timestampToStringPtr(sharingRecord.AcceptedAt),
	}

	c.JSON(http.StatusCreated, response)
}

// GetPendingShares retrieves pending sharing invitations for the current user
func (h *ShareHandler) GetPendingShares(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	queries := h.services.GetDB().GetQueries()

	records, err := queries.GetPendingSharingRecordsByRecipientID(c.Request.Context(), userID.(int32))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch pending shares"})
		return
	}

	response := make([]SharingRecordResponse, len(records))
	for i, record := range records {
		response[i] = SharingRecordResponse{
			ID:              uuidToString(record.ID),
			VaultID:         record.VaultID,
			ItemID:          uuidToStringPtr(record.ItemID),
			SenderUserID:    record.SenderUserID,
			RecipientUserID: record.RecipientUserID,
			WrappedKey:      crypto.EncodeBase64(record.WrappedKey),
			WrapIV:          crypto.EncodeBase64(record.WrapIv),
			WrapTag:         crypto.EncodeBase64(record.WrapTag),
			Status:          record.Status,
			CreatedAt:       timestampToTime(record.CreatedAt).Format("2006-01-02T15:04:05Z07:00"),
			AcceptedAt:      timestampToStringPtr(record.AcceptedAt),
		}
	}

	c.JSON(http.StatusOK, response)
}

// AcceptShare accepts a sharing invitation
func (h *ShareHandler) AcceptShare(c *gin.Context) {
	shareID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid share_id"})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	queries := h.services.GetDB().GetQueries()

	pgShareID := pgtype.UUID{}
	_ = pgShareID.Scan(shareID.String())

	// Get sharing record to verify recipient
	sharingRecord, err := queries.GetSharingRecordByID(c.Request.Context(), pgShareID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "sharing record not found"})
		return
	}

	// Verify user is the recipient
	if sharingRecord.RecipientUserID != userID.(int32) {
		c.JSON(http.StatusForbidden, gin.H{"error": "not authorized to accept this share"})
		return
	}

	// Accept the share
	updatedRecord, err := queries.AcceptSharingRecord(c.Request.Context(), pgShareID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to accept share"})
		return
	}

	response := SharingRecordResponse{
		ID:              uuidToString(updatedRecord.ID),
		VaultID:         updatedRecord.VaultID,
		ItemID:          uuidToStringPtr(updatedRecord.ItemID),
		SenderUserID:    updatedRecord.SenderUserID,
		RecipientUserID: updatedRecord.RecipientUserID,
		WrappedKey:      crypto.EncodeBase64(updatedRecord.WrappedKey),
		WrapIV:          crypto.EncodeBase64(updatedRecord.WrapIv),
		WrapTag:         crypto.EncodeBase64(updatedRecord.WrapTag),
		Status:          updatedRecord.Status,
		CreatedAt:       timestampToTime(updatedRecord.CreatedAt).Format("2006-01-02T15:04:05Z07:00"),
		AcceptedAt:      timestampToStringPtr(updatedRecord.AcceptedAt),
	}

	c.JSON(http.StatusOK, response)
}

// RejectShare rejects a sharing invitation
// POST /api/shares/:id/reject
func (h *ShareHandler) RejectShare(c *gin.Context) {
	shareID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid share_id"})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	queries := h.services.GetDB().GetQueries()

	pgShareID := pgtype.UUID{}
	_ = pgShareID.Scan(shareID.String())

	// Get sharing record to verify recipient
	sharingRecord, err := queries.GetSharingRecordByID(c.Request.Context(), pgShareID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "sharing record not found"})
		return
	}

	// Verify user is the recipient
	if sharingRecord.RecipientUserID != userID.(int32) {
		c.JSON(http.StatusForbidden, gin.H{"error": "not authorized to reject this share"})
		return
	}

	// Reject the share
	updatedRecord, err := queries.RejectSharingRecord(c.Request.Context(), pgShareID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to reject share"})
		return
	}

	response := SharingRecordResponse{
		ID:              uuidToString(updatedRecord.ID),
		VaultID:         updatedRecord.VaultID,
		ItemID:          uuidToStringPtr(updatedRecord.ItemID),
		SenderUserID:    updatedRecord.SenderUserID,
		RecipientUserID: updatedRecord.RecipientUserID,
		WrappedKey:      crypto.EncodeBase64(updatedRecord.WrappedKey),
		WrapIV:          crypto.EncodeBase64(updatedRecord.WrapIv),
		WrapTag:         crypto.EncodeBase64(updatedRecord.WrapTag),
		Status:          updatedRecord.Status,
		CreatedAt:       timestampToTime(updatedRecord.CreatedAt).Format("2006-01-02T15:04:05Z07:00"),
		AcceptedAt:      timestampToStringPtr(updatedRecord.AcceptedAt),
	}

	c.JSON(http.StatusOK, response)
}

// GetSharedVaults retrieves all vaults shared with the current user
func (h *ShareHandler) GetSharedVaults(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	queries := h.services.GetDB().GetQueries()

	vaults, err := queries.GetSharedVaultsForUser(c.Request.Context(), userID.(int32))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch shared vaults"})
		return
	}

	response := make([]SharedVaultResponse, len(vaults))
	for i, vault := range vaults {
		response[i] = SharedVaultResponse{
			ID:        vault.ID,
			Name:      vault.Name,
			UserID:    vault.UserID,
			CreatedAt: timestampToTime(vault.CreatedAt).Format("2006-01-02T15:04:05Z07:00"),
			UpdatedAt: timestampToTime(vault.UpdatedAt).Format("2006-01-02T15:04:05Z07:00"),
		}
	}

	c.JSON(http.StatusOK, response)
}

// RevokeShare revokes a sharing record (vault owner only)
func (h *ShareHandler) RevokeShare(c *gin.Context) {
	shareID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid share_id"})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	// Verify device signature
	sigData, err := ExtractDeviceSignature(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "device authentication required: " + err.Error()})
		return
	}

	queries := h.services.GetDB().GetQueries()

	pgShareID := pgtype.UUID{}
	_ = pgShareID.Scan(shareID.String())

	// Get sharing record
	sharingRecord, err := queries.GetSharingRecordByID(c.Request.Context(), pgShareID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "sharing record not found"})
		return
	}

	// Get vault to verify ownership
	vault, err := queries.GetVaultByID(c.Request.Context(), sharingRecord.VaultID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "vault not found"})
		return
	}

	// Only vault owner can revoke shares
	if vault.UserID != userID.(int32) {
		c.JSON(http.StatusForbidden, gin.H{"error": "only vault owner can revoke shares"})
		return
	}

	// Verify device
	pgDeviceID := pgtype.UUID{}
	_ = pgDeviceID.Scan(sigData.DeviceID.String())

	device, err := queries.GetDeviceByID(c.Request.Context(), pgDeviceID)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "device not found"})
		return
	}

	if device.UserID != userID.(int32) {
		c.JSON(http.StatusForbidden, gin.H{"error": "device does not belong to user"})
		return
	}

	// Verify signature
	bodyBytes, _ := c.Get("raw_body")
	bodyHash := sha256.Sum256(bodyBytes.([]byte))
	canonicalMsg := CreateCanonicalMessage(c.Request.Method, c.Request.URL.Path, sigData.Timestamp, bodyHash[:])

	if !VerifyDeviceSignature(device.Ed25519Public, canonicalMsg, sigData.Signature) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid device signature"})
		return
	}

	// Revoke the share
	err = queries.RevokeSharingRecord(c.Request.Context(), pgShareID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to revoke share"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "share revoked successfully"})
}

// Helper function to convert pgtype.UUID to string pointer (for nullable UUID)
func uuidToStringPtr(u pgtype.UUID) *string {
	if !u.Valid {
		return nil
	}
	str := uuidToString(u)
	return &str
}

// Helper function to convert pgtype.Timestamp to string pointer (for nullable timestamps)
func timestampToStringPtr(t pgtype.Timestamp) *string {
	if !t.Valid {
		return nil
	}
	str := timestampToTime(t).Format("2006-01-02T15:04:05Z07:00")
	return &str
}
