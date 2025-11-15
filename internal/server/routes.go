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
		AllowHeaders:     []string{"Accept", "Authorization", "Content-Type"},
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
	pageHandler := handlers.NewPageHandler(s.services)

	r.GET("/", s.HelloWorldHandler)
	r.GET("/health", s.healthHandler)

	auth := r.Group("/api/auth")
	{
		auth.POST("/register", authHandler.Register)
		auth.POST("/login", authHandler.Login)
		auth.POST("/logout", authHandler.Logout)
	}

	r.GET("/api/:handle", pageHandler.GetPageByHandle)

	// Protected routes
	protected := r.Group("/api")
	protected.Use(middleware.AuthMiddleware(s.services))
	{
		protected.GET("/me", authHandler.Me)

		// Page routes
		protected.POST("/pages", authHandler.CreatePage)
		protected.PUT("/pages", pageHandler.UpdatePage)
		protected.DELETE("/pages/:id", pageHandler.DeletePage)
		protected.PUT("/pages/:id/activate", pageHandler.SetActivePage)
		protected.GET("/pages/:id", pageHandler.GetActivePage)
		protected.GET("/pages/:handle", pageHandler.CheckHandleExists)

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
