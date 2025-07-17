package main

import (
	"flag"
	"log"

	"zipcodereader/config"
	"zipcodereader/database"
	"zipcodereader/handlers"
	"zipcodereader/middleware"
	"zipcodereader/models"
	"zipcodereader/services"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

func main() {
	// Parse command line flags
	useLocalAuth := flag.Bool("use_local_auth", false, "Use local authentication instead of GitHub OAuth2")
	flag.Parse()

	// Load configuration
	cfg := config.Load(*useLocalAuth)

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
	r.Use(sessions.Sessions("zipcodereader", store))

	// Add middleware
	r.Use(middleware.Logger())
	r.Use(middleware.CORS())

	// Load HTML templates
	r.LoadHTMLGlob("templates/*")

	// Serve static files
	r.Static("/static", "./static")

	// Initialize handlers
	h := handlers.New(db)

	// Setup authentication routes based on mode
	if cfg.UseLocalAuth {
		log.Println("Using local authentication mode")
		localAuthHandler := handlers.NewLocalAuthHandler(db)

		// Local authentication routes
		r.GET("/local/login", localAuthHandler.ShowLogin)
		r.POST("/local/login", localAuthHandler.Login)
		r.GET("/local/register", localAuthHandler.ShowRegister)
		r.POST("/local/register", localAuthHandler.Register)
		r.GET("/local/logout", localAuthHandler.Logout)

		// Dashboard route
		protected := r.Group("/")
		protected.Use(middleware.RequireAuth())
		{
			protected.GET("/dashboard", func(c *gin.Context) {
				session := sessions.Default(c)
				userID := session.Get("user_id")

				if userID == nil {
					c.Redirect(302, "/local/login")
					return
				}

				// Get user information
				user, err := models.GetUserByID(db, userID.(uint))
				if err != nil {
					c.JSON(500, gin.H{"error": "Failed to get user information"})
					return
				}

				c.HTML(200, "dashboard.html", gin.H{
					"title": "Dashboard",
					"user":  user,
				})
			})
		}

		// Update home page context
		r.GET("/", func(c *gin.Context) {
			c.HTML(200, "index.html", gin.H{
				"title":          "ZipCodeReader",
				"message":        "Welcome to ZipCodeReader - A reading list manager for students",
				"use_local_auth": true,
			})
		})
	} else {
		log.Println("Using GitHub OAuth2 authentication mode")
		authService := services.NewAuthService(db, cfg)
		authHandler := handlers.NewAuthHandler(authService)

		// GitHub OAuth2 routes
		r.GET("/auth/login", authHandler.Login)
		r.GET("/auth/callback", authHandler.Callback)
		r.GET("/auth/logout", authHandler.Logout)

		// Protected routes
		protected := r.Group("/")
		protected.Use(middleware.RequireAuth())
		{
			protected.GET("/dashboard", authHandler.Dashboard)
		}

		// Home page
		r.GET("/", h.Home)
	}

	// Common routes
	r.GET("/health", h.Health)

	// Start server
	log.Printf("Server starting on port %s", cfg.Port)
	log.Printf("Authentication mode: %s", func() string {
		if cfg.UseLocalAuth {
			return "Local"
		}
		return "GitHub OAuth2"
	}())

	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
