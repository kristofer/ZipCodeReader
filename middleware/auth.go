package middleware

import (
	"net/http"
	"zipcodereader/models"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Logger middleware for request logging
func Logger() gin.HandlerFunc {
	return gin.LoggerWithConfig(gin.LoggerConfig{
		SkipPaths: []string{"/health"},
	})
}

// CORS middleware for cross-origin requests
func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

// RequireAuth middleware ensures user is authenticated
func RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		userID := session.Get("user_id")

		if userID == nil {
			c.Redirect(http.StatusTemporaryRedirect, "/")
			c.Abort()
			return
		}

		// Set user info in context for handlers to use
		c.Set("user_id", userID)
		c.Set("user_role", session.Get("user_role"))
		c.Next()
	}
}

// RequireAuthWithUser middleware ensures user is authenticated and loads user object
func RequireAuthWithUser(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		userID := session.Get("user_id")

		if userID == nil {
			c.Redirect(http.StatusTemporaryRedirect, "/")
			c.Abort()
			return
		}

		// Get user from database
		user, err := models.GetUserByID(db, userID.(uint))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "User not found"})
			c.Abort()
			return
		}

		// Set user info in context for handlers to use
		c.Set("user", user)
		c.Set("user_id", userID)
		c.Set("user_role", user.Role)
		c.Next()
	}
}

// RequireRole middleware ensures user has specific role
func RequireRole(role string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, exists := c.Get("user_role")
		if !exists || userRole != role {
			c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
			c.Abort()
			return
		}
		c.Next()
	}
}

// RequireInstructor middleware ensures user is an instructor
func RequireInstructor() gin.HandlerFunc {
	return RequireRole("instructor")
}

// RequireStudent middleware ensures user is a student
func RequireStudent() gin.HandlerFunc {
	return RequireRole("student")
}
