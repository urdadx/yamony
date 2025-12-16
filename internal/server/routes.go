package server

import (
	"net/http"

	"yamony/internal/server/handlers"
	"yamony/internal/server/middleware"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

func (s *Server) RegisterRoutes() http.Handler {
	r := gin.Default()

	// CORS middleware
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3001", "https://yamony.com"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowHeaders:     []string{"Accept", "Authorization", "Content-Type", "X-Device-Id", "X-Device-Signature", "X-Device-Timestamp"},
		AllowCredentials: true,
	}))

	sessionSecret := []byte("secret-key-change-this-in-production")
	store := cookie.NewStore(sessionSecret)
	store.Options(sessions.Options{
		Path:     "/",
		MaxAge:   86400 * 7,
		HttpOnly: true,
		Secure:   false, // Set to true in production with HTTPS
		SameSite: http.SameSiteLaxMode,
	})
	r.Use(sessions.Sessions("yamony_session", store))

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(s.services)
	deviceHandler := handlers.NewDeviceHandler(s.services)
	vaultHandler := handlers.NewVaultHandler(s.services)
	vaultKeyHandler := handlers.NewVaultKeyHandler(s.services)
	vaultItemHandler := handlers.NewVaultItemHandler(s.services)
	shareHandler := handlers.NewShareHandler(s.services)
	syncHandler := handlers.NewSyncHandler(s.services)
	uploadHandler := handlers.NewUploadHandler()

	r.GET("/", s.HelloWorldHandler)
	r.GET("/health", s.healthHandler)

	auth := r.Group("/api")
	{
		auth.POST("/register", authHandler.Register)
		auth.POST("/login", authHandler.Login)
		auth.POST("/logout", authHandler.Logout)

		// Google OAuth routes
		auth.GET("/auth/google", authHandler.GoogleLogin)
		auth.GET("/auth/google/callback", authHandler.GoogleCallback)
	}

	// Protected routes
	protected := r.Group("/api")
	protected.Use(middleware.AuthMiddleware(s.services))
	{
		protected.GET("/me", authHandler.Me)
		protected.GET("/sessions", authHandler.GetSessionByUserID)

		// Device routes
		protected.POST("/devices/register", deviceHandler.RegisterDevice)
		protected.POST("/devices/verify", deviceHandler.VerifyDevice)
		protected.GET("/devices", deviceHandler.GetDevices)
		protected.DELETE("/devices/:id", deviceHandler.RevokeDevice)
		protected.GET("/users/:user_id/public-keys", deviceHandler.GetUserPublicKeys)

		// Vault routes
		protected.POST("/vaults", vaultHandler.CreateVault)
		protected.GET("/vaults", vaultHandler.GetVaults)
		protected.GET("/vaults/:id", vaultHandler.GetVault)
		protected.PUT("/vaults/:id", vaultHandler.UpdateVault)
		protected.DELETE("/vaults/:id", vaultHandler.DeleteVault)

		// Vault key routes
		protected.POST("/vaults/:id/keys", vaultKeyHandler.UploadVaultKey)
		protected.GET("/vaults/:id/keys", vaultKeyHandler.GetVaultKey)
		protected.GET("/vaults/:id/keys/versions", vaultKeyHandler.GetVaultKeyVersions)
		protected.GET("/vaults/:id/keys/versions/:version", vaultKeyHandler.GetVaultKeyVersion)

		// Vault item routes
		protected.POST("/vaults/:id/items", vaultItemHandler.CreateVaultItem)
		protected.GET("/vaults/:id/items", vaultItemHandler.GetVaultItems)
		protected.GET("/vaults/:id/items/:item_id", vaultItemHandler.GetVaultItem)
		protected.PUT("/vaults/:id/items/:item_id", vaultItemHandler.UpdateVaultItem)
		protected.DELETE("/vaults/:id/items/:item_id", vaultItemHandler.DeleteVaultItem)

		// Sharing routes
		protected.POST("/vaults/:id/share", shareHandler.ShareVault)
		protected.GET("/vaults/shared", shareHandler.GetSharedVaults)
		protected.GET("/shares/pending", shareHandler.GetPendingShares)
		protected.POST("/shares/:id/accept", shareHandler.AcceptShare)
		protected.POST("/shares/:id/reject", shareHandler.RejectShare)
		protected.DELETE("/shares/:id", shareHandler.RevokeShare)

		// Sync and versioning routes
		protected.POST("/vaults/:id/sync/pull", syncHandler.PullVaultChanges)
		protected.POST("/vaults/:id/sync/commit", syncHandler.CommitVaultChanges)
		protected.GET("/vaults/:id/versions", syncHandler.GetVaultVersions)

		// Storage upload route
		protected.POST("/storage/upload", uploadHandler.UploadToGCS)

	}

	return r
}

func (s *Server) HelloWorldHandler(c *gin.Context) {
	resp := make(map[string]string)
	resp["message"] = "Hello World"

	c.JSON(http.StatusOK, resp)
}

func (s *Server) healthHandler(c *gin.Context) {
	c.JSON(http.StatusOK, s.db.Health())
}
