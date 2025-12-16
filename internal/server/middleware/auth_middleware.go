package middleware

import (
	"net/http"
	"yamony/internal/server/services"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

const (
	SessionTokenKey = "session_token"
	UserKey         = "user"
	ActivePageID    = "active_page_id"
)

func AuthMiddleware(service services.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		sessionToken := session.Get(SessionTokenKey)

		if sessionToken == nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "unauthorized: no session found",
			})
			c.Abort()
			return
		}

		// Validate session
		user, err := service.ValidateSession(c.Request.Context(), sessionToken.(string))
		if err != nil {
			session.Delete(SessionTokenKey)
			session.Save()

			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "unauthorized: invalid or expired session",
			})
			c.Abort()
			return
		}

		c.Set("user_id", user.ID)
		c.Set(UserKey, user)
		c.Next()
	}
}
