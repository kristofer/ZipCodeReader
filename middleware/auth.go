package middleware

import (
	"net/http"
	"strings"
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

// isAPIRequest checks if the request is likely an API/AJAX request
func isAPIRequest(c *gin.Context) bool {
	// Check for common API request indicators
	accept := c.GetHeader("Accept")
	xRequestedWith := c.GetHeader("X-Requested-With")
	contentType := c.GetHeader("Content-Type")

	// If request explicitly asks for JSON or is XMLHttpRequest
	if strings.Contains(accept, "application/json") || xRequestedWith == "XMLHttpRequest" {
		return true
	}

	// If it's a JSON content type (for POST/PUT requests)
	if strings.Contains(contentType, "application/json") {
		return true
	}

	// Check if the URL path suggests it's an API endpoint
	path := c.Request.URL.Path
	return strings.Contains(path, "/api/") ||
		strings.Contains(path, "/stats") ||
		strings.HasSuffix(path, "/assignments") ||
		strings.Contains(path, "/assignments/")
}

// RequireAuth middleware ensures user is authenticated
func RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		userID := session.Get("user_id")

		if userID == nil {
			// Check if this is an API request (AJAX/fetch)
			if isAPIRequest(c) {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
			} else {
				c.Redirect(http.StatusTemporaryRedirect, "/")
			}
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
			// Check if this is an API request (AJAX/fetch)
			if isAPIRequest(c) {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
			} else {
				c.Redirect(http.StatusTemporaryRedirect, "/")
			}
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
