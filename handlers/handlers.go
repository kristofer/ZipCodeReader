package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Handler holds dependencies for HTTP handlers
type Handler struct {
	db *gorm.DB
}

// New creates a new Handler instance
func New(db *gorm.DB) *Handler {
	return &Handler{
		db: db,
	}
}

// Home handles the home page
func (h *Handler) Home(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", gin.H{
		"title":   "ZipCodeReader",
		"message": "Welcome to ZipCodeReader - A reading list manager for students",
	})
}

// Health handles health check endpoint
func (h *Handler) Health(c *gin.Context) {
	// Check database connection
	sqlDB, err := h.db.DB()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":    "error",
			"message":   "Database connection failed",
			"timestamp": time.Now().Format(time.RFC3339),
		})
		return
	}

	if err := sqlDB.Ping(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":    "error",
			"message":   "Database ping failed",
			"timestamp": time.Now().Format(time.RFC3339),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":    "healthy",
		"message":   "ZipCodeReader is running",
		"database":  "connected",
		"timestamp": time.Now().Format(time.RFC3339),
	})
}
