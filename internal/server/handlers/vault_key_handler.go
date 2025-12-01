package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"

	"yamony/internal/crypto"
	"yamony/internal/database/sqlc"
	"yamony/internal/server/services"
)

type VaultKeyHandler struct {
	services services.Service
}

func NewVaultKeyHandler(services services.Service) *VaultKeyHandler {
	return &VaultKeyHandler{services: services}
}

// UploadVaultKeyRequest represents the request to upload a wrapped VEK
type UploadVaultKeyRequest struct {
	WrappedVEK string          `json:"wrapped_vek" binding:"required"`
	WrapIV     string          `json:"wrap_iv" binding:"required"`
	WrapTag    string          `json:"wrap_tag" binding:"required"`
	KDFSalt    string          `json:"kdf_salt" binding:"required"`
	KDFParams  json.RawMessage `json:"kdf_params" binding:"required"`
	Version    int32           `json:"version"` // defaults to 1
}

// VaultKeyResponse represents a vault key
type VaultKeyResponse struct {
	VaultID    int32           `json:"vault_id"`
	WrappedVEK string          `json:"wrapped_vek"`
	WrapIV     string          `json:"wrap_iv"`
	WrapTag    string          `json:"wrap_tag"`
	KDFSalt    string          `json:"kdf_salt"`
	KDFParams  json.RawMessage `json:"kdf_params"`
	Version    int32           `json:"version"`
	CreatedAt  string          `json:"created_at"`
	UpdatedAt  string          `json:"updated_at"`
}

// UploadVaultKey uploads a wrapped VEK for a vault
// POST /api/vaults/:id/keys
func (h *VaultKeyHandler) UploadVaultKey(c *gin.Context) {
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

	var req UploadVaultKeyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Verify vault belongs to user
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

	// Decode base64 fields
	wrappedVEK, err := crypto.DecodeBase64(req.WrappedVEK)
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

	kdfSalt, err := crypto.DecodeBase64(req.KDFSalt)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid kdf_salt format"})
		return
	}

	// Validate KDF params
	var kdfParams crypto.KDFParams
	if err := json.Unmarshal(req.KDFParams, &kdfParams); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid kdf_params format"})
		return
	}

	if err := crypto.ValidateKDFParams(kdfParams); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Set default version if not provided
	version := req.Version
	if version == 0 {
		version = 1
	}

	// Create vault key
	vaultKey, err := queries.CreateVaultKey(c.Request.Context(), sqlc.CreateVaultKeyParams{
		VaultID:    vaultID,
		WrappedVek: wrappedVEK,
		WrapIv:     wrapIV,
		WrapTag:    wrapTag,
		KdfSalt:    kdfSalt,
		KdfParams:  req.KDFParams,
		Version:    version,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create vault key"})
		return
	}

	response := VaultKeyResponse{
		VaultID:    vaultKey.VaultID,
		WrappedVEK: crypto.EncodeBase64(vaultKey.WrappedVek),
		WrapIV:     crypto.EncodeBase64(vaultKey.WrapIv),
		WrapTag:    crypto.EncodeBase64(vaultKey.WrapTag),
		KDFSalt:    crypto.EncodeBase64(vaultKey.KdfSalt),
		KDFParams:  vaultKey.KdfParams,
		Version:    vaultKey.Version,
		CreatedAt:  timestampToTime(vaultKey.CreatedAt).Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:  timestampToTime(vaultKey.UpdatedAt).Format("2006-01-02T15:04:05Z07:00"),
	}

	c.JSON(http.StatusCreated, response)
}

// GetVaultKey retrieves the wrapped VEK for a vault
// GET /api/vaults/:id/keys
func (h *VaultKeyHandler) GetVaultKey(c *gin.Context) {
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

	// Verify vault belongs to user
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

	// Get vault key (latest version)
	vaultKey, err := queries.GetVaultKeyByVaultID(c.Request.Context(), vaultID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "vault key not found"})
		return
	}

	response := VaultKeyResponse{
		VaultID:    vaultKey.VaultID,
		WrappedVEK: crypto.EncodeBase64(vaultKey.WrappedVek),
		WrapIV:     crypto.EncodeBase64(vaultKey.WrapIv),
		WrapTag:    crypto.EncodeBase64(vaultKey.WrapTag),
		KDFSalt:    crypto.EncodeBase64(vaultKey.KdfSalt),
		KDFParams:  vaultKey.KdfParams,
		Version:    vaultKey.Version,
		CreatedAt:  timestampToTime(vaultKey.CreatedAt).Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:  timestampToTime(vaultKey.UpdatedAt).Format("2006-01-02T15:04:05Z07:00"),
	}

	c.JSON(http.StatusOK, response)
}

// GetVaultKeyVersion retrieves a specific version of the wrapped VEK
// GET /api/vaults/:id/keys/versions/:version
func (h *VaultKeyHandler) GetVaultKeyVersion(c *gin.Context) {
	vaultID, err := parseIntParam(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid vault_id"})
		return
	}

	version, err := parseIntParam(c.Param("version"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid version"})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	// Verify vault belongs to user
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

	// Get specific vault key version
	vaultKey, err := queries.GetVaultKeyByVaultIDAndVersion(c.Request.Context(), sqlc.GetVaultKeyByVaultIDAndVersionParams{
		VaultID: vaultID,
		Version: version,
	})
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "vault key version not found"})
		return
	}

	response := VaultKeyResponse{
		VaultID:    vaultKey.VaultID,
		WrappedVEK: crypto.EncodeBase64(vaultKey.WrappedVek),
		WrapIV:     crypto.EncodeBase64(vaultKey.WrapIv),
		WrapTag:    crypto.EncodeBase64(vaultKey.WrapTag),
		KDFSalt:    crypto.EncodeBase64(vaultKey.KdfSalt),
		KDFParams:  vaultKey.KdfParams,
		Version:    vaultKey.Version,
		CreatedAt:  timestampToTime(vaultKey.CreatedAt).Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:  timestampToTime(vaultKey.UpdatedAt).Format("2006-01-02T15:04:05Z07:00"),
	}

	c.JSON(http.StatusOK, response)
}

// GetVaultKeyVersions lists all versions of vault keys
// GET /api/vaults/:id/keys/versions
func (h *VaultKeyHandler) GetVaultKeyVersions(c *gin.Context) {
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

	// Verify vault belongs to user
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

	// Get all vault key versions
	vaultKeys, err := queries.GetAllVaultKeyVersions(c.Request.Context(), vaultID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch vault key versions"})
		return
	}

	response := make([]VaultKeyResponse, len(vaultKeys))
	for i, vaultKey := range vaultKeys {
		response[i] = VaultKeyResponse{
			VaultID:    vaultKey.VaultID,
			WrappedVEK: crypto.EncodeBase64(vaultKey.WrappedVek),
			WrapIV:     crypto.EncodeBase64(vaultKey.WrapIv),
			WrapTag:    crypto.EncodeBase64(vaultKey.WrapTag),
			KDFSalt:    crypto.EncodeBase64(vaultKey.KdfSalt),
			KDFParams:  vaultKey.KdfParams,
			Version:    vaultKey.Version,
			CreatedAt:  timestampToTime(vaultKey.CreatedAt).Format("2006-01-02T15:04:05Z07:00"),
			UpdatedAt:  timestampToTime(vaultKey.UpdatedAt).Format("2006-01-02T15:04:05Z07:00"),
		}
	}

	c.JSON(http.StatusOK, response)
}
