package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"

	"yamony/internal/database/sqlc"
	"yamony/internal/server/services"
)

type VaultHandler struct {
	services services.Service
}

func NewVaultHandler(services services.Service) *VaultHandler {
	return &VaultHandler{services: services}
}

// CreateVaultRequest represents the request to create a vault
type CreateVaultRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
	Theme       string `json:"theme"`
	IsFavorite  bool   `json:"is_favorite"`
}

// VaultResponse represents a vault
type VaultResponse struct {
	ID          int32  `json:"id"`
	UserID      int32  `json:"user_id"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Icon        string `json:"icon,omitempty"`
	Theme       string `json:"theme,omitempty"`
	IsFavorite  bool   `json:"is_favorite"`
	ItemCount   int32  `json:"item_count"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

// CreateVault creates a new vault
func (h *VaultHandler) CreateVault(c *gin.Context) {
	var req CreateVaultRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body", "details": err.Error()})
		return
	}

	// Get authenticated user
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	queries := h.services.GetDB().GetQueries()

	// Create description pgtype.Text
	var description pgtype.Text
	if req.Description != "" {
		description = pgtype.Text{
			String: req.Description,
			Valid:  true,
		}
	}

	// Create icon pgtype.Text
	var icon pgtype.Text
	if req.Icon != "" {
		icon = pgtype.Text{
			String: req.Icon,
			Valid:  true,
		}
	}

	// Create theme pgtype.Text
	var theme pgtype.Text
	if req.Theme != "" {
		theme = pgtype.Text{
			String: req.Theme,
			Valid:  true,
		}
	}

	vault, err := queries.CreateVault(c.Request.Context(), sqlc.CreateVaultParams{
		UserID:      userID.(int32),
		Name:        req.Name,
		Description: description,
		Icon:        icon,
		Theme:       theme,
		IsFavorite:  req.IsFavorite,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create vault", "details": err.Error()})
		return
	}

	response := VaultResponse{
		ID:         vault.ID,
		UserID:     vault.UserID,
		Name:       vault.Name,
		IsFavorite: vault.IsFavorite,
		ItemCount:  0, // New vaults have no items
		CreatedAt:  vault.CreatedAt.Time.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:  vault.UpdatedAt.Time.Format("2006-01-02T15:04:05Z07:00"),
	}

	if vault.Description.Valid {
		response.Description = vault.Description.String
	}

	if vault.Icon.Valid {
		response.Icon = vault.Icon.String
	}

	if vault.Theme.Valid {
		response.Theme = vault.Theme.String
	}

	c.JSON(http.StatusCreated, response)
}

// GetVaults retrieves all vaults for the authenticated user
func (h *VaultHandler) GetVaults(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	queries := h.services.GetDB().GetQueries()

	vaults, err := queries.GetUserVaults(c.Request.Context(), userID.(int32))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve vaults", "details": err.Error()})
		return
	}

	var response []VaultResponse
	for _, vault := range vaults {
		vaultResp := VaultResponse{
			ID:         vault.ID,
			UserID:     vault.UserID,
			Name:       vault.Name,
			IsFavorite: vault.IsFavorite,
			ItemCount:  vault.ItemCount,
			CreatedAt:  vault.CreatedAt.Time.Format("2006-01-02T15:04:05Z07:00"),
			UpdatedAt:  vault.UpdatedAt.Time.Format("2006-01-02T15:04:05Z07:00"),
		}

		if vault.Description.Valid {
			vaultResp.Description = vault.Description.String
		}

		if vault.Icon.Valid {
			vaultResp.Icon = vault.Icon.String
		}

		if vault.Theme.Valid {
			vaultResp.Theme = vault.Theme.String
		}

		response = append(response, vaultResp)
	}

	c.JSON(http.StatusOK, response)
}

// GetVault retrieves a specific vault
func (h *VaultHandler) GetVault(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var vaultID int32
	if _, err := fmt.Sscanf(c.Param("id"), "%d", &vaultID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid vault ID"})
		return
	}

	queries := h.services.GetDB().GetQueries()

	vault, err := queries.GetVaultByIDAndUserID(c.Request.Context(), sqlc.GetVaultByIDAndUserIDParams{
		ID:     vaultID,
		UserID: userID.(int32),
	})

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Vault not found"})
		return
	}

	response := VaultResponse{
		ID:         vault.ID,
		UserID:     vault.UserID,
		Name:       vault.Name,
		IsFavorite: vault.IsFavorite,
		CreatedAt:  vault.CreatedAt.Time.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:  vault.UpdatedAt.Time.Format("2006-01-02T15:04:05Z07:00"),
	}

	if vault.Description.Valid {
		response.Description = vault.Description.String
	}

	if vault.Icon.Valid {
		response.Icon = vault.Icon.String
	}

	if vault.Theme.Valid {
		response.Theme = vault.Theme.String
	}

	c.JSON(http.StatusOK, response)
}

// UpdateVault updates a vault
func (h *VaultHandler) UpdateVault(c *gin.Context) {
	var req CreateVaultRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body", "details": err.Error()})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var vaultID int32
	if _, err := fmt.Sscanf(c.Param("id"), "%d", &vaultID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid vault ID"})
		return
	}

	queries := h.services.GetDB().GetQueries()

	// Create description pgtype.Text
	var description pgtype.Text
	if req.Description != "" {
		description = pgtype.Text{
			String: req.Description,
			Valid:  true,
		}
	}

	// Create icon pgtype.Text
	var icon pgtype.Text
	if req.Icon != "" {
		icon = pgtype.Text{
			String: req.Icon,
			Valid:  true,
		}
	}

	// Create theme pgtype.Text
	var theme pgtype.Text
	if req.Theme != "" {
		theme = pgtype.Text{
			String: req.Theme,
			Valid:  true,
		}
	}

	vault, err := queries.UpdateVault(c.Request.Context(), sqlc.UpdateVaultParams{
		ID:          vaultID,
		Name:        req.Name,
		Description: description,
		Icon:        icon,
		Theme:       theme,
		IsFavorite:  req.IsFavorite,
		UserID:      userID.(int32),
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update vault", "details": err.Error()})
		return
	}

	response := VaultResponse{
		ID:         vault.ID,
		UserID:     vault.UserID,
		Name:       vault.Name,
		IsFavorite: vault.IsFavorite,
		CreatedAt:  vault.CreatedAt.Time.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:  vault.UpdatedAt.Time.Format("2006-01-02T15:04:05Z07:00"),
	}

	if vault.Description.Valid {
		response.Description = vault.Description.String
	}

	if vault.Icon.Valid {
		response.Icon = vault.Icon.String
	}

	if vault.Theme.Valid {
		response.Theme = vault.Theme.String
	}

	c.JSON(http.StatusOK, response)
}

// DeleteVault deletes a vault
func (h *VaultHandler) DeleteVault(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var vaultID int32
	if _, err := fmt.Sscanf(c.Param("id"), "%d", &vaultID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid vault ID"})
		return
	}

	queries := h.services.GetDB().GetQueries()

	err := queries.DeleteVault(c.Request.Context(), sqlc.DeleteVaultParams{
		ID:     vaultID,
		UserID: userID.(int32),
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete vault", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Vault deleted successfully"})
}
