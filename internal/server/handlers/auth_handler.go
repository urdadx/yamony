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

type AuthHandler struct {
	service services.Service
}

func NewAuthHandler(service services.Service) *AuthHandler {
	return &AuthHandler{
		service: service,
	}
}

type RegisterRequest struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type UserResponse struct {
	ID            int32  `json:"id"`
	Username      string `json:"username"`
	Email         string `json:"email"`
	EmailVerified bool   `json:"emailVerified"`
	Image         string `json:"image"`
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	user, sessionToken, activePageID, err := h.service.RegisterUser(
		c.Request.Context(),
		req.Username,
		req.Email,
		req.Password,
	)
	if err != nil {
		if err == services.ErrEmailAlreadyExists {
			c.JSON(http.StatusConflict, gin.H{
				"error": "user with this email already exists",
			})
			return
		}

		fmt.Println("Registration error ", err)

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to register user",
		})
		return
	}

	session := sessions.Default(c)
	session.Set(middleware.SessionTokenKey, sessionToken)

	// Set active page in session if user has pages (activePageID > 0)
	if activePageID > 0 {
		session.Set(middleware.ActivePageID, activePageID)
	}

	if err := session.Save(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to save session",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "user registered successfully",
		"user": UserResponse{
			ID:            user.ID,
			Username:      user.Username,
			Email:         user.Email,
			EmailVerified: user.EmailVerified,
			Image:         user.Image,
		},
	})
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	user, sessionToken, activePageID, err := h.service.LoginUser(
		c.Request.Context(),
		req.Email,
		req.Password,
	)
	if err != nil {
		if err == services.ErrInvalidCredentials {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "invalid email or password",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to login",
		})
		fmt.Println("Login error ", err)
		return
	}

	session := sessions.Default(c)
	session.Set(middleware.SessionTokenKey, sessionToken)

	// Restore the user's last active page in the session
	if activePageID > 0 {
		session.Set(middleware.ActivePageID, activePageID)
	}

	if err := session.Save(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to save session",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "login successful",
		"user": UserResponse{
			ID:            user.ID,
			Username:      user.Username,
			Email:         user.Email,
			EmailVerified: user.EmailVerified,
			Image:         user.Image,
		},
	})
}

func (h *AuthHandler) Logout(c *gin.Context) {
	session := sessions.Default(c)
	sessionToken := session.Get(middleware.SessionTokenKey)

	if sessionToken != nil {
		_ = h.service.LogoutUser(c.Request.Context(), sessionToken.(string))
	}

	session.Clear()
	if err := session.Save(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to clear session",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "logout successful",
	})
}

func (h *AuthHandler) Me(c *gin.Context) {
	userValue, exists := c.Get(middleware.UserKey)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "unauthorized",
		})
		return
	}

	user := userValue.(*sqlc.GetUserByIDRow)
	c.JSON(http.StatusOK, gin.H{
		"user": UserResponse{
			ID:            user.ID,
			Username:      user.Username,
			Email:         user.Email,
			EmailVerified: user.EmailVerified,
			Image:         user.Image,
		},
	})
}
