package handlers

import (
	"fmt"
	"net/http"
	"yamony/internal/database/sqlc"
	"yamony/internal/server/middleware"
	"yamony/internal/server/services"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

type PageHandler struct {
	service services.Service
}

func NewPageHandler(service services.Service) *PageHandler {
	return &PageHandler{
		service: service,
	}
}

type CreatePageRequest struct {
	Handle   string `json:"handle"`
	IsActive bool   `json:"is_active"`
}

type UpdatePageRequest struct {
	IsActive    bool   `json:"is_active"`
	Name        string `json:"name"`
	Handle      string `json:"handle"`
	Image       string `json:"image"`
	BannerImage string `json:"banner_image"`
	Bio         string `json:"bio"`
}

type PageResponse struct {
	ID          int32  `json:"id"`
	Handle      string `json:"handle"`
	IsActive    bool   `json:"is_active"`
	Name        string `json:"name"`
	Image       string `json:"image"`
	BannerImage string `json:"banner_image"`
	Bio         string `json:"bio"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

func (h *AuthHandler) CreatePage(c *gin.Context) {
	var req CreatePageRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	userValue, exists := c.Get(middleware.UserKey)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "unauthorized",
		})
		return
	}

	user := userValue.(*sqlc.GetUserByIDRow)

	page, err := h.service.CreatePage(c.Request.Context(), user.ID, req.Handle, req.IsActive)

	if err != nil {
		if err == services.ErrPageAlreadyExists {
			c.JSON(http.StatusConflict, gin.H{
				"error": "page with handle already exists",
			})
			return
		}

		fmt.Println("Page creation error ", err)

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to create page",
		})
		return
	}

	session := sessions.Default(c)
	sessionToken := session.Get(middleware.SessionTokenKey)
	session.Set(middleware.ActivePageID, page.ID)

	if sessionToken != nil {
		err = h.service.SyncActivePageToSession(c.Request.Context(), sessionToken.(string), page.ID)
		if err != nil {
			fmt.Printf("Warning: failed to sync active page to database session: %v\n", err)
		}
	}

	if err := session.Save(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to save session",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "page created successfully",
		"user": PageResponse{
			ID:          user.ID,
			Handle:      page.Handle,
			IsActive:    page.IsActive,
			Name:        page.Name.String,
			Image:       page.Image.String,
			BannerImage: page.BannerImage.String,
			Bio:         page.Bio.String,
			CreatedAt:   page.CreatedAt.Time.String(),
			UpdatedAt:   page.UpdatedAt.Time.String(),
		},
	})

}

func (h *PageHandler) UpdatePage(c *gin.Context) {
	var req UpdatePageRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	session := sessions.Default(c)
	activePageID := session.Get(middleware.ActivePageID)
	if activePageID == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "no active page selected",
		})
		return
	}

	_, exists := c.Get(middleware.UserKey)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "unauthorized",
		})
		return
	}

	pageID := activePageID.(int32)

	page, err := h.service.UpdatePage(
		c.Request.Context(),
		pageID,
		req.IsActive,
		req.Name,
		req.Handle,
		req.Image,
		req.BannerImage,
		req.Bio,
	)

	if err != nil {
		if err == services.ErrPageNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "page not found",
			})
			return
		}
		if err == services.ErrHandleAlreadyExists {
			c.JSON(http.StatusConflict, gin.H{
				"error": "page handle already exists",
			})
			return
		}

		fmt.Println("Page update error:", err)

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to update page",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "page updated successfully",
		"page": PageResponse{
			ID:          page.ID,
			Handle:      page.Handle,
			IsActive:    page.IsActive,
			Name:        page.Name.String,
			Image:       page.Image.String,
			BannerImage: page.BannerImage.String,
			Bio:         page.Bio.String,
			CreatedAt:   page.CreatedAt.Time.String(),
			UpdatedAt:   page.UpdatedAt.Time.String(),
		},
	})
}

func (h *PageHandler) CheckHandleExists(c *gin.Context) {
	var req CreatePageRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	exists, err := h.service.CheckHandleExists(c.Request.Context(), req.Handle)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to check handle",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"exists": exists,
	})

}

func (h *PageHandler) GetActivePage(c *gin.Context) {
	session := sessions.Default(c)
	activePageID := session.Get(middleware.ActivePageID)
	if activePageID == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "no active page selected",
		})
		return
	}

	pageID := activePageID.(int32)
	page, err := h.service.GetPageByID(c.Request.Context(), pageID)
	if err != nil {
		if err == services.ErrPageNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "page not found",
			})
			return
		}

		fmt.Println("Get active page error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to get active page",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"page": PageResponse{
			ID:          page.ID,
			Handle:      page.Handle,
			IsActive:    page.IsActive,
			Name:        page.Name.String,
			Image:       page.Image.String,
			BannerImage: page.BannerImage.String,
			Bio:         page.Bio.String,
			CreatedAt:   page.CreatedAt.Time.String(),
			UpdatedAt:   page.UpdatedAt.Time.String(),
		},
	})
}

func (h *PageHandler) DeletePage(c *gin.Context) {
	session := sessions.Default(c)
	activePageID := session.Get(middleware.ActivePageID)
	if activePageID == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "no active page selected",
		})
		return
	}

	pageID := activePageID.(int32)
	userValue, exists := c.Get(middleware.UserKey)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "unauthorized",
		})
		return
	}

	user := userValue.(*sqlc.GetUserByIDRow)
	err := h.service.DeletePage(c.Request.Context(), pageID, user.ID)
	if err != nil {
		if err == services.ErrPageNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "page not found",
			})
			return
		}
		if err.Error() == "cannot delete the only page" {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "cannot delete the only page",
			})
			return
		}
		fmt.Println("Delete page error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to delete page",
		})
		return
	}

	nextPage, err := h.service.SetNextPageAsActive(c.Request.Context(), user.ID)
	if err != nil {
		fmt.Println("Set next active page error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to set next active page",
		})
		return
	}

	sessionToken := session.Get(middleware.SessionTokenKey)
	session.Set(middleware.ActivePageID, nextPage.ID)

	if sessionToken != nil {
		err = h.service.SyncActivePageToSession(c.Request.Context(), sessionToken.(string), nextPage.ID)
		if err != nil {
			fmt.Printf("Warning: failed to sync active page to database session: %v\n", err)
		}
	}

	if err := session.Save(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to save session",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "page deleted successfully",
	})

}

func (h *PageHandler) SetActivePage(c *gin.Context) {
	var req struct {
		PageID int32 `json:"page_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	session := sessions.Default(c)
	sessionToken := session.Get(middleware.SessionTokenKey)

	// Update the cookie session
	session.Set(middleware.ActivePageID, req.PageID)

	err := h.service.SetActivePage(c.Request.Context(), req.PageID)
	if err != nil {
		if err == services.ErrPageNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "page not found",
			})
			return
		}

		fmt.Println("Set active page error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to set active page",
		})
		return
	}

	if sessionToken != nil {
		err = h.service.SyncActivePageToSession(c.Request.Context(), sessionToken.(string), req.PageID)
		if err != nil {
			fmt.Printf("Warning: failed to sync active page to database session: %v\n", err)
		}
	}

	if err := session.Save(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to save session",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "active page set successfully",
	})
}

func (h *PageHandler) GetPageByHandle(c *gin.Context) {
	handle := c.Param("handle")
	page, err := h.service.GetPageByHandle(c.Request.Context(), handle)
	if err != nil {
		if err == services.ErrPageNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "page not found",
			})
			return
		}

		fmt.Println("Get page by handle error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to get page by handle",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"page": PageResponse{
			ID:          page.ID,
			Handle:      page.Handle,
			IsActive:    page.IsActive,
			Name:        page.Name.String,
			Image:       page.Image.String,
			BannerImage: page.BannerImage.String,
			Bio:         page.Bio.String,
			CreatedAt:   page.CreatedAt.Time.String(),
			UpdatedAt:   page.UpdatedAt.Time.String(),
		},
	})
}

func (h *PageHandler) GetAllUserPage(c *gin.Context) {
	userValue, exists := c.Get(middleware.UserKey)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "unauthorized",
		})
		return
	}

	user := userValue.(*sqlc.GetUserByIDRow)

	pages, err := h.service.GetAllUserPage(c.Request.Context(), user.ID)
	if err != nil {
		fmt.Println("Get all user pages error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to get all user pages",
		})
		return
	}

	var pageResponses []PageResponse
	for _, page := range pages {
		pageResponses = append(pageResponses, PageResponse{
			ID:          page.ID,
			Handle:      page.Handle,
			IsActive:    page.IsActive,
			Name:        page.Name.String,
			Image:       page.Image.String,
			BannerImage: page.BannerImage.String,
			Bio:         page.Bio.String,
			CreatedAt:   page.CreatedAt.Time.String(),
			UpdatedAt:   page.UpdatedAt.Time.String(),
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"pages": pageResponses,
	})
}
