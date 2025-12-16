package handlers

import (
	"crypto/sha256"
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"

	"yamony/internal/crypto"
	"yamony/internal/database/sqlc"
	"yamony/internal/server/services"
)

type VaultItemHandler struct {
	services services.Service
}

func NewVaultItemHandler(services services.Service) *VaultItemHandler {
	return &VaultItemHandler{services: services}
}

// CreateVaultItemRequest represents the request to create a vault item
type CreateVaultItemRequest struct {
	ItemType      string          `json:"item_type" binding:"required"` // login, note, card, alias
	EncryptedBlob string          `json:"encrypted_blob" binding:"required"`
	IV            string          `json:"iv" binding:"required"`
	Tag           string          `json:"tag" binding:"required"`
	Meta          json.RawMessage `json:"meta,omitempty"`
	Version       int32           `json:"version"`
}

// UpdateVaultItemRequest represents the request to update a vault item
type UpdateVaultItemRequest struct {
	EncryptedBlob string          `json:"encrypted_blob" binding:"required"`
	IV            string          `json:"iv" binding:"required"`
	Tag           string          `json:"tag" binding:"required"`
	Meta          json.RawMessage `json:"meta,omitempty"`
	BaseVersion   int32           `json:"base_version" binding:"required"`
}

// VaultItemResponse represents a vault item
type VaultItemResponse struct {
	ID            string          `json:"id"`
	VaultID       int32           `json:"vault_id"`
	ItemType      string          `json:"item_type"`
	EncryptedBlob string          `json:"encrypted_blob"`
	IV            string          `json:"iv"`
	Tag           string          `json:"tag"`
	Meta          json.RawMessage `json:"meta,omitempty"`
	Version       int32           `json:"version"`
	CreatedAt     string          `json:"created_at"`
	UpdatedAt     string          `json:"updated_at"`
}

// VaultItemListResponse is a simplified response for listing items
type VaultItemListResponse struct {
	ID        string          `json:"id"`
	ItemType  string          `json:"item_type"`
	Meta      json.RawMessage `json:"meta,omitempty"`
	Version   int32           `json:"version"`
	CreatedAt string          `json:"created_at"`
	UpdatedAt string          `json:"updated_at"`
}

// CreateVaultItem creates a new encrypted vault item
// POST /api/vaults/:id/items
func (h *VaultItemHandler) CreateVaultItem(c *gin.Context) {
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

	var req CreateVaultItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Verify device signature on write operation
	sigData, err := ExtractDeviceSignature(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "device authentication required: " + err.Error()})
		return
	}

	queries := h.services.GetDB().GetQueries()

	vault, err := queries.GetVaultByID(c.Request.Context(), vaultID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "vault not found"})
		return
	}

	if vault.UserID != userID.(int32) {
		c.JSON(http.StatusForbidden, gin.H{"error": "not authorized to access this vault"})
		return
	}

	// Verify device belongs to user and is not revoked
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

	// Decode encrypted data
	encryptedBlob, err := crypto.DecodeBase64(req.EncryptedBlob)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid encrypted_blob format"})
		return
	}

	iv, err := crypto.DecodeBase64(req.IV)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid iv format"})
		return
	}

	tag, err := crypto.DecodeBase64(req.Tag)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid tag format"})
		return
	}

	// Set default version
	version := req.Version
	if version == 0 {
		version = 1
	}

	// Generate item ID
	itemID := uuid.New()
	pgItemID := pgtype.UUID{}
	_ = pgItemID.Scan(itemID.String())

	// Create vault item
	vaultItem, err := queries.CreateVaultItem(c.Request.Context(), sqlc.CreateVaultItemParams{
		ID:            pgItemID,
		VaultID:       vaultID,
		ItemType:      req.ItemType,
		EncryptedBlob: encryptedBlob,
		Iv:            iv,
		Tag:           tag,
		Meta:          req.Meta,
		Version:       version,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create vault item"})
		return
	}

	response := VaultItemResponse{
		ID:            uuidToString(vaultItem.ID),
		VaultID:       vaultItem.VaultID,
		ItemType:      vaultItem.ItemType,
		EncryptedBlob: crypto.EncodeBase64(vaultItem.EncryptedBlob),
		IV:            crypto.EncodeBase64(vaultItem.Iv),
		Tag:           crypto.EncodeBase64(vaultItem.Tag),
		Meta:          vaultItem.Meta,
		Version:       vaultItem.Version,
		CreatedAt:     timestampToTime(vaultItem.CreatedAt).Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:     timestampToTime(vaultItem.UpdatedAt).Format("2006-01-02T15:04:05Z07:00"),
	}

	c.JSON(http.StatusCreated, response)
}

// GetVaultItems lists all items in a vault (without encrypted blobs for efficiency)
// GET /api/vaults/:id/items
func (h *VaultItemHandler) GetVaultItems(c *gin.Context) {
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

	queries := h.services.GetDB().GetQueries()

	// Verify vault belongs to user
	vault, err := queries.GetVaultByID(c.Request.Context(), vaultID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "vault not found"})
		return
	}

	if vault.UserID != userID.(int32) {
		c.JSON(http.StatusForbidden, gin.H{"error": "not authorized to access this vault"})
		return
	}

	// Optional filter by item type
	itemType := c.Query("type")

	var items []sqlc.VaultItem
	if itemType != "" {
		items, err = queries.GetVaultItemsByVaultIDAndType(c.Request.Context(), sqlc.GetVaultItemsByVaultIDAndTypeParams{
			VaultID:  vaultID,
			ItemType: itemType,
		})
	} else {
		items, err = queries.GetVaultItemsByVaultID(c.Request.Context(), vaultID)
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch vault items"})
		return
	}

	// Return simplified list without encrypted blobs
	response := make([]VaultItemListResponse, len(items))
	for i, item := range items {
		response[i] = VaultItemListResponse{
			ID:        uuidToString(item.ID),
			ItemType:  item.ItemType,
			Meta:      item.Meta,
			Version:   item.Version,
			CreatedAt: timestampToTime(item.CreatedAt).Format("2006-01-02T15:04:05Z07:00"),
			UpdatedAt: timestampToTime(item.UpdatedAt).Format("2006-01-02T15:04:05Z07:00"),
		}
	}

	c.JSON(http.StatusOK, response)
}

// GetVaultItem retrieves a specific vault item with encrypted blob
// GET /api/vaults/:id/items/:item_id
func (h *VaultItemHandler) GetVaultItem(c *gin.Context) {
	vaultID, err := parseIntParam(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid vault_id"})
		return
	}

	itemID, err := uuid.Parse(c.Param("item_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid item_id"})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
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
		c.JSON(http.StatusForbidden, gin.H{"error": "not authorized to access this vault"})
		return
	}

	pgItemID := pgtype.UUID{}
	_ = pgItemID.Scan(itemID.String())

	// Get vault item
	vaultItem, err := queries.GetVaultItemByID(c.Request.Context(), pgItemID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "vault item not found"})
		return
	}

	// Verify item belongs to the vault
	if vaultItem.VaultID != vaultID {
		c.JSON(http.StatusForbidden, gin.H{"error": "item does not belong to this vault"})
		return
	}

	response := VaultItemResponse{
		ID:            uuidToString(vaultItem.ID),
		VaultID:       vaultItem.VaultID,
		ItemType:      vaultItem.ItemType,
		EncryptedBlob: crypto.EncodeBase64(vaultItem.EncryptedBlob),
		IV:            crypto.EncodeBase64(vaultItem.Iv),
		Tag:           crypto.EncodeBase64(vaultItem.Tag),
		Meta:          vaultItem.Meta,
		Version:       vaultItem.Version,
		CreatedAt:     timestampToTime(vaultItem.CreatedAt).Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:     timestampToTime(vaultItem.UpdatedAt).Format("2006-01-02T15:04:05Z07:00"),
	}

	c.JSON(http.StatusOK, response)
}

// UpdateVaultItem updates an encrypted vault item with optimistic concurrency
// PUT /api/vaults/:id/items/:item_id
func (h *VaultItemHandler) UpdateVaultItem(c *gin.Context) {
	vaultID, err := parseIntParam(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid vault_id"})
		return
	}

	itemID, err := uuid.Parse(c.Param("item_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid item_id"})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	var req UpdateVaultItemRequest
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
		c.JSON(http.StatusForbidden, gin.H{"error": "not authorized to access this vault"})
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

	pgItemID := pgtype.UUID{}
	_ = pgItemID.Scan(itemID.String())

	// Get current item to check version (optimistic concurrency)
	currentItem, err := queries.GetVaultItemByID(c.Request.Context(), pgItemID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "vault item not found"})
		return
	}

	// Verify item belongs to vault
	if currentItem.VaultID != vaultID {
		c.JSON(http.StatusForbidden, gin.H{"error": "item does not belong to this vault"})
		return
	}

	// Check version for optimistic concurrency control
	if currentItem.Version != req.BaseVersion {
		c.JSON(http.StatusConflict, gin.H{
			"error":            "version conflict",
			"current_version":  currentItem.Version,
			"provided_version": req.BaseVersion,
		})
		return
	}

	// Decode encrypted data
	encryptedBlob, err := crypto.DecodeBase64(req.EncryptedBlob)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid encrypted_blob format"})
		return
	}

	iv, err := crypto.DecodeBase64(req.IV)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid iv format"})
		return
	}

	tag, err := crypto.DecodeBase64(req.Tag)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid tag format"})
		return
	}

	// Update with incremented version
	newVersion := currentItem.Version + 1

	updatedItem, err := queries.UpdateVaultItem(c.Request.Context(), sqlc.UpdateVaultItemParams{
		ID:            pgItemID,
		EncryptedBlob: encryptedBlob,
		Iv:            iv,
		Tag:           tag,
		Meta:          req.Meta,
		Version:       newVersion,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update vault item"})
		return
	}

	response := VaultItemResponse{
		ID:            uuidToString(updatedItem.ID),
		VaultID:       updatedItem.VaultID,
		ItemType:      updatedItem.ItemType,
		EncryptedBlob: crypto.EncodeBase64(updatedItem.EncryptedBlob),
		IV:            crypto.EncodeBase64(updatedItem.Iv),
		Tag:           crypto.EncodeBase64(updatedItem.Tag),
		Meta:          updatedItem.Meta,
		Version:       updatedItem.Version,
		CreatedAt:     timestampToTime(updatedItem.CreatedAt).Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:     timestampToTime(updatedItem.UpdatedAt).Format("2006-01-02T15:04:05Z07:00"),
	}

	c.JSON(http.StatusOK, response)
}

// DeleteVaultItem deletes a vault item
// DELETE /api/vaults/:id/items/:item_id
func (h *VaultItemHandler) DeleteVaultItem(c *gin.Context) {
	vaultID, err := parseIntParam(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid vault_id"})
		return
	}

	itemID, err := uuid.Parse(c.Param("item_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid item_id"})
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

	// Verify vault belongs to user
	vault, err := queries.GetVaultByID(c.Request.Context(), vaultID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "vault not found"})
		return
	}

	if vault.UserID != userID.(int32) {
		c.JSON(http.StatusForbidden, gin.H{"error": "not authorized to access this vault"})
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

	pgItemID := pgtype.UUID{}
	_ = pgItemID.Scan(itemID.String())

	// Verify item exists and belongs to vault
	item, err := queries.GetVaultItemByID(c.Request.Context(), pgItemID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "vault item not found"})
		return
	}

	if item.VaultID != vaultID {
		c.JSON(http.StatusForbidden, gin.H{"error": "item does not belong to this vault"})
		return
	}

	// Delete item
	err = queries.DeleteVaultItem(c.Request.Context(), pgItemID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete vault item"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "vault item deleted successfully"})
}
