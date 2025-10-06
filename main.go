package main

import (
	"flag"
	"log"
	"net/http"
	"zipcodereader/config"
	"zipcodereader/database"
	"zipcodereader/middleware"
	"zipcodereader/routes"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

func main() {
	// Parse command line flags
	useOAuth2 := flag.Bool("use_oauth2", false, "Use GitHub OAuth2 authentication instead of local authentication")
	flag.Parse()

	// Load configuration (local auth is default, OAuth2 is optional)
	cfg := config.Load(!*useOAuth2)

	// Initialize database
	db, err := database.Initialize(cfg.DatabaseURL)
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}

	// Set Gin mode
	if cfg.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Create Gin router
	r := gin.Default()

	// Session middleware
	store := cookie.NewStore([]byte(cfg.SessionSecret))
	store.Options(sessions.Options{
		Path:     "/",
		MaxAge:   86400 * 30, // 30 days
		HttpOnly: true,
		Secure:   false, // Set to true in production with HTTPS
		SameSite: http.SameSiteLaxMode,
	})
	r.Use(sessions.Sessions("zipcodereader", store))

	// Add middleware
	r.Use(middleware.Logger())
	r.Use(middleware.CORS())

	// Load HTML templates
	r.LoadHTMLGlob("templates/*")

	// Serve static files
	r.Static("/static", "./static")

	// Register all routes (authentication mode handled inside)
	routes.Register(r, db, cfg)

	// Start server
	log.Printf("Server starting on port %s", cfg.Port)
	log.Printf("Authentication mode: %s", func() string {
		if cfg.UseLocalAuth {
			return "Local (default)"
		}
		return "GitHub OAuth2 (optional)"
	}())

	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
