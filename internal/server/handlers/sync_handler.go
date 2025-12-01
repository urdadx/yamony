package handlers

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"

	"yamony/internal/crypto"
	"yamony/internal/database/sqlc"
	"yamony/internal/server/services"
)

type SyncHandler struct {
	services services.Service
}

func NewSyncHandler(services services.Service) *SyncHandler {
	return &SyncHandler{services: services}
}

// SyncPullRequest represents a request to pull vault changes
type SyncPullRequest struct {
	LastSyncedVersionID *int32 `json:"last_synced_version_id,omitempty"`
}

// SyncPullResponse represents the response with vault changes
type SyncPullResponse struct {
	VaultID         int32               `json:"vault_id"`
	CurrentVersion  int32               `json:"current_version"`
	Items           []VaultItemResponse `json:"items"`
	DeletedItemIDs  []string            `json:"deleted_item_ids,omitempty"`
	ETag            string              `json:"etag"`
	HasMoreVersions bool                `json:"has_more_versions"`
}

// SyncCommitRequest represents a request to commit vault changes
type SyncCommitRequest struct {
	BaseVersionID *int32           `json:"base_version_id,omitempty"` // for optimistic concurrency
	Items         []SyncItemCommit `json:"items" binding:"required"`
	DeletedItems  []string         `json:"deleted_items,omitempty"`
}

// SyncItemCommit represents an item being committed
type SyncItemCommit struct {
	ID            *string `json:"id,omitempty"` // nil for new items
	ItemType      string  `json:"item_type" binding:"required"`
	EncryptedBlob string  `json:"encrypted_blob" binding:"required"`
	IV            string  `json:"iv" binding:"required"`
	Tag           string  `json:"tag" binding:"required"`
	Meta          []byte  `json:"meta,omitempty"`
	BaseVersion   *int32  `json:"base_version,omitempty"` // for optimistic concurrency on updates
}

// SyncCommitResponse represents the response after committing changes
type SyncCommitResponse struct {
	VaultID        int32               `json:"vault_id"`
	NewVersionID   int32               `json:"new_version_id"`
	CommittedItems []VaultItemResponse `json:"committed_items"`
	Conflicts      []SyncConflict      `json:"conflicts,omitempty"`
	ETag           string              `json:"etag"`
}

// SyncConflict represents a conflict detected during commit
type SyncConflict struct {
	ItemID          string `json:"item_id"`
	ConflictType    string `json:"conflict_type"` // "version_mismatch", "concurrent_edit"
	CurrentVersion  int32  `json:"current_version"`
	AttemptedAction string `json:"attempted_action"` // "create", "update", "delete"
}

// VaultVersionResponse represents a vault version/snapshot
type VaultVersionResponse struct {
	ID              int32   `json:"id"`
	VaultID         int32   `json:"vault_id"`
	ObjectKey       string  `json:"object_key"`
	MAC             *string `json:"mac,omitempty"`
	CreatedByDevice *string `json:"created_by_device,omitempty"`
	CreatedAt       string  `json:"created_at"`
}

// PullVaultChanges retrieves changes since the last sync
// POST /api/vaults/:id/sync/pull
func (h *SyncHandler) PullVaultChanges(c *gin.Context) {
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

	var req SyncPullRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		// No body is also acceptable - sync from beginning
		req.LastSyncedVersionID = nil
	}

	queries := h.services.GetDB().GetQueries()

	// Check if user has access to vault (owner or shared)
	accessLevel, err := queries.CheckUserVaultAccess(c.Request.Context(), sqlc.CheckUserVaultAccessParams{
		ID:     vaultID,
		UserID: userID.(int32),
	})
	if err != nil || accessLevel == nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "no access to this vault"})
		return
	}

	// Get all vault items (in a real system with large vaults, implement pagination)
	items, err := queries.GetVaultItemsByVaultID(c.Request.Context(), vaultID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch vault items"})
		return
	}

	// Get latest version for ETag generation
	latestVersion, err := queries.GetLatestVaultVersion(c.Request.Context(), vaultID)
	var versionID int32 = 0
	if err == nil {
		versionID = latestVersion.ID
	}

	// Convert items to response format
	itemResponses := make([]VaultItemResponse, len(items))
	for i, item := range items {
		itemResponses[i] = VaultItemResponse{
			ID:            uuidToString(item.ID),
			VaultID:       item.VaultID,
			ItemType:      item.ItemType,
			EncryptedBlob: crypto.EncodeBase64(item.EncryptedBlob),
			IV:            crypto.EncodeBase64(item.Iv),
			Tag:           crypto.EncodeBase64(item.Tag),
			Meta:          item.Meta,
			Version:       item.Version,
			CreatedAt:     timestampToTime(item.CreatedAt).Format("2006-01-02T15:04:05Z07:00"),
			UpdatedAt:     timestampToTime(item.UpdatedAt).Format("2006-01-02T15:04:05Z07:00"),
		}
	}

	// Generate ETag based on vault state
	etag := generateVaultETag(vaultID, versionID, items)

	// Check If-None-Match header
	if match := c.GetHeader("If-None-Match"); match != "" && match == etag {
		c.Status(http.StatusNotModified)
		return
	}

	response := SyncPullResponse{
		VaultID:         vaultID,
		CurrentVersion:  versionID,
		Items:           itemResponses,
		DeletedItemIDs:  []string{}, // Would track deletions in production
		ETag:            etag,
		HasMoreVersions: false,
	}

	c.Header("ETag", etag)
	c.JSON(http.StatusOK, response)
}

// CommitVaultChanges commits local changes to the vault
// POST /api/vaults/:id/sync/commit
func (h *SyncHandler) CommitVaultChanges(c *gin.Context) {
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

	var req SyncCommitRequest
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

	// Check vault access
	accessLevel, err := queries.CheckUserVaultAccess(c.Request.Context(), sqlc.CheckUserVaultAccessParams{
		ID:     vaultID,
		UserID: userID.(int32),
	})
	if err != nil || accessLevel == nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "no access to this vault"})
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

	// Check If-Match header for optimistic concurrency
	ifMatch := c.GetHeader("If-Match")
	if ifMatch != "" {
		// Get current vault state
		items, _ := queries.GetVaultItemsByVaultID(c.Request.Context(), vaultID)
		latestVersion, err := queries.GetLatestVaultVersion(c.Request.Context(), vaultID)
		var versionID int32 = 0
		if err == nil {
			versionID = latestVersion.ID
		}
		currentETag := generateVaultETag(vaultID, versionID, items)

		if ifMatch != currentETag {
			c.JSON(http.StatusPreconditionFailed, gin.H{
				"error":        "vault state has changed, pull latest changes first",
				"current_etag": currentETag,
			})
			return
		}
	}

	// Process commits
	var committedItems []VaultItemResponse
	var conflicts []SyncConflict

	// Process new items and updates
	for _, itemCommit := range req.Items {
		if itemCommit.ID == nil || *itemCommit.ID == "" {
			// Create new item
			encryptedBlob, err := crypto.DecodeBase64(itemCommit.EncryptedBlob)
			if err != nil {
				continue
			}
			iv, err := crypto.DecodeBase64(itemCommit.IV)
			if err != nil {
				continue
			}
			tag, err := crypto.DecodeBase64(itemCommit.Tag)
			if err != nil {
				continue
			}

			itemID := uuid.New()
			pgItemID := pgtype.UUID{}
			_ = pgItemID.Scan(itemID.String())

			createdItem, err := queries.CreateVaultItem(c.Request.Context(), sqlc.CreateVaultItemParams{
				ID:            pgItemID,
				VaultID:       vaultID,
				ItemType:      itemCommit.ItemType,
				EncryptedBlob: encryptedBlob,
				Iv:            iv,
				Tag:           tag,
				Meta:          itemCommit.Meta,
				Version:       1,
			})
			if err != nil {
				continue
			}

			committedItems = append(committedItems, VaultItemResponse{
				ID:            uuidToString(createdItem.ID),
				VaultID:       createdItem.VaultID,
				ItemType:      createdItem.ItemType,
				EncryptedBlob: crypto.EncodeBase64(createdItem.EncryptedBlob),
				IV:            crypto.EncodeBase64(createdItem.Iv),
				Tag:           crypto.EncodeBase64(createdItem.Tag),
				Meta:          createdItem.Meta,
				Version:       createdItem.Version,
				CreatedAt:     timestampToTime(createdItem.CreatedAt).Format("2006-01-02T15:04:05Z07:00"),
				UpdatedAt:     timestampToTime(createdItem.UpdatedAt).Format("2006-01-02T15:04:05Z07:00"),
			})
		} else {
			// Update existing item
			itemID, err := uuid.Parse(*itemCommit.ID)
			if err != nil {
				continue
			}

			pgItemID := pgtype.UUID{}
			_ = pgItemID.Scan(itemID.String())

			// Get current item for version check
			currentItem, err := queries.GetVaultItemByID(c.Request.Context(), pgItemID)
			if err != nil {
				conflicts = append(conflicts, SyncConflict{
					ItemID:          *itemCommit.ID,
					ConflictType:    "not_found",
					AttemptedAction: "update",
				})
				continue
			}

			// Check optimistic concurrency
			if itemCommit.BaseVersion != nil && currentItem.Version != *itemCommit.BaseVersion {
				conflicts = append(conflicts, SyncConflict{
					ItemID:          *itemCommit.ID,
					ConflictType:    "version_mismatch",
					CurrentVersion:  currentItem.Version,
					AttemptedAction: "update",
				})
				continue
			}

			encryptedBlob, err := crypto.DecodeBase64(itemCommit.EncryptedBlob)
			if err != nil {
				continue
			}
			iv, err := crypto.DecodeBase64(itemCommit.IV)
			if err != nil {
				continue
			}
			tag, err := crypto.DecodeBase64(itemCommit.Tag)
			if err != nil {
				continue
			}

			updatedItem, err := queries.UpdateVaultItem(c.Request.Context(), sqlc.UpdateVaultItemParams{
				ID:            pgItemID,
				EncryptedBlob: encryptedBlob,
				Iv:            iv,
				Tag:           tag,
				Meta:          itemCommit.Meta,
				Version:       currentItem.Version + 1,
			})
			if err != nil {
				continue
			}

			committedItems = append(committedItems, VaultItemResponse{
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
			})
		}
	}

	// Process deletions
	for _, itemIDStr := range req.DeletedItems {
		itemID, err := uuid.Parse(itemIDStr)
		if err != nil {
			continue
		}

		pgItemID := pgtype.UUID{}
		_ = pgItemID.Scan(itemID.String())

		err = queries.DeleteVaultItem(c.Request.Context(), pgItemID)
		if err != nil {
			conflicts = append(conflicts, SyncConflict{
				ItemID:          itemIDStr,
				ConflictType:    "not_found",
				AttemptedAction: "delete",
			})
		}
	}

	// Create vault version snapshot (simplified - just storing metadata)
	objectKey := fmt.Sprintf("vaults/%d/versions/%d.snapshot", vaultID, time.Now().Unix())
	vaultVersion, err := queries.CreateVaultVersion(c.Request.Context(), sqlc.CreateVaultVersionParams{
		VaultID:   vaultID,
		ObjectKey: objectKey,
		Mac:       nil, // Would compute MAC in production
		CreatedByDevice: pgtype.UUID{
			Bytes: sigData.DeviceID,
			Valid: true,
		},
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create version"})
		return
	}

	// Generate new ETag
	items, _ := queries.GetVaultItemsByVaultID(c.Request.Context(), vaultID)
	newETag := generateVaultETag(vaultID, vaultVersion.ID, items)

	response := SyncCommitResponse{
		VaultID:        vaultID,
		NewVersionID:   vaultVersion.ID,
		CommittedItems: committedItems,
		Conflicts:      conflicts,
		ETag:           newETag,
	}

	c.Header("ETag", newETag)
	if len(conflicts) > 0 {
		c.JSON(http.StatusConflict, response)
	} else {
		c.JSON(http.StatusOK, response)
	}
}

// GetVaultVersions retrieves version history for a vault
// GET /api/vaults/:id/versions
func (h *SyncHandler) GetVaultVersions(c *gin.Context) {
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

	// Parse limit parameter
	limitStr := c.DefaultQuery("limit", "10")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 || limit > 100 {
		limit = 10
	}

	queries := h.services.GetDB().GetQueries()

	// Check vault access
	accessLevel, err := queries.CheckUserVaultAccess(c.Request.Context(), sqlc.CheckUserVaultAccessParams{
		ID:     vaultID,
		UserID: userID.(int32),
	})
	if err != nil || accessLevel == nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "no access to this vault"})
		return
	}

	// Get versions
	versions, err := queries.GetVaultVersionsByVaultID(c.Request.Context(), sqlc.GetVaultVersionsByVaultIDParams{
		VaultID: vaultID,
		Limit:   int32(limit),
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch versions"})
		return
	}

	response := make([]VaultVersionResponse, len(versions))
	for i, version := range versions {
		response[i] = VaultVersionResponse{
			ID:              version.ID,
			VaultID:         version.VaultID,
			ObjectKey:       version.ObjectKey,
			MAC:             bytesToStringPtr(version.Mac),
			CreatedByDevice: uuidToStringPtr(version.CreatedByDevice),
			CreatedAt:       timestampToTime(version.CreatedAt).Format("2006-01-02T15:04:05Z07:00"),
		}
	}

	c.JSON(http.StatusOK, response)
}

// generateVaultETag generates an ETag for vault state
func generateVaultETag(vaultID int32, versionID int32, items []sqlc.VaultItem) string {
	hash := sha256.New()
	hash.Write([]byte(fmt.Sprintf("%d:%d:%d", vaultID, versionID, len(items))))
	for _, item := range items {
		hash.Write(item.ID.Bytes[:])
		hash.Write([]byte(fmt.Sprintf(":%d", item.Version)))
	}
	return hex.EncodeToString(hash.Sum(nil))
}

// Helper function to convert []byte to string pointer
func bytesToStringPtr(b []byte) *string {
	if b == nil {
		return nil
	}
	str := crypto.EncodeBase64(b)
	return &str
}
