package main

import (
	"flag"
	"log"

	"zipcodereader/config"
	"zipcodereader/database"
	"zipcodereader/handlers"
	"zipcodereader/middleware"
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

	// Initialize assignment services
	assignmentService := services.NewAssignmentService(db)
	studentAssignmentService := services.NewStudentAssignmentService(db)

	// Initialize assignment handlers
	instructorAssignmentHandlers := handlers.NewInstructorAssignmentHandlers(assignmentService)
	studentAssignmentHandlers := handlers.NewStudentAssignmentHandlers(studentAssignmentService)

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
		protected.Use(middleware.RequireAuthWithUser(db))
		{
			protected.GET("/dashboard", func(c *gin.Context) {
				user, exists := c.Get("user")
				if !exists {
					c.JSON(500, gin.H{"error": "User not found"})
					return
				}

				c.HTML(200, "dashboard.html", gin.H{
					"title": "Dashboard",
					"user":  user,
				})
			})

			// Instructor assignment routes
			instructorGroup := protected.Group("/instructor")
			instructorGroup.Use(middleware.RequireRole("instructor"))
			{
				instructorGroup.GET("/assignments", instructorAssignmentHandlers.GetAssignments)
				instructorGroup.POST("/assignments", instructorAssignmentHandlers.CreateAssignment)
				instructorGroup.GET("/assignments/:id", instructorAssignmentHandlers.GetAssignment)
				instructorGroup.PUT("/assignments/:id", instructorAssignmentHandlers.UpdateAssignment)
				instructorGroup.DELETE("/assignments/:id", instructorAssignmentHandlers.DeleteAssignment)
				instructorGroup.POST("/assignments/:id/assign", instructorAssignmentHandlers.AssignStudents)
				instructorGroup.GET("/assignments/:id/progress", instructorAssignmentHandlers.GetAssignmentProgress)
				instructorGroup.GET("/assignments/:id/students", instructorAssignmentHandlers.GetAssignmentStudents)
				instructorGroup.POST("/assignments/:id/students/:student_id/remove", instructorAssignmentHandlers.RemoveStudent)
				instructorGroup.GET("/students", instructorAssignmentHandlers.GetAllStudents)
				instructorGroup.GET("/dashboard/stats", instructorAssignmentHandlers.GetDashboardStats)
			}

			// Student assignment routes
			studentGroup := protected.Group("/student")
			studentGroup.Use(middleware.RequireRole("student"))
			{
				studentGroup.GET("/assignments", studentAssignmentHandlers.GetAssignments)
				studentGroup.GET("/assignments/:id", studentAssignmentHandlers.GetAssignment)
				studentGroup.POST("/assignments/:id/status", studentAssignmentHandlers.UpdateStatus)
				studentGroup.POST("/assignments/:id/complete", studentAssignmentHandlers.MarkAsCompleted)
				studentGroup.POST("/assignments/:id/progress", studentAssignmentHandlers.MarkAsInProgress)
				studentGroup.GET("/dashboard/stats", studentAssignmentHandlers.GetDashboardStats)
				studentGroup.GET("/assignments/overdue", studentAssignmentHandlers.GetOverdueAssignments)
				studentGroup.GET("/assignments/upcoming", studentAssignmentHandlers.GetUpcomingAssignments)
				studentGroup.GET("/assignments/recent", studentAssignmentHandlers.GetRecentlyCompleted)
				studentGroup.GET("/categories", studentAssignmentHandlers.GetCategories)
				studentGroup.GET("/assignments/by-status", studentAssignmentHandlers.GetAssignmentsByStatus)
				studentGroup.GET("/assignments/by-category", studentAssignmentHandlers.GetAssignmentsByCategory)
				studentGroup.GET("/assignments/search", studentAssignmentHandlers.SearchAssignments)
			}
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
		protected.Use(middleware.RequireAuthWithUser(db))
		{
			protected.GET("/dashboard", authHandler.Dashboard)

			// Instructor assignment routes
			instructorGroup := protected.Group("/instructor")
			instructorGroup.Use(middleware.RequireRole("instructor"))
			{
				instructorGroup.GET("/assignments", instructorAssignmentHandlers.GetAssignments)
				instructorGroup.POST("/assignments", instructorAssignmentHandlers.CreateAssignment)
				instructorGroup.GET("/assignments/:id", instructorAssignmentHandlers.GetAssignment)
				instructorGroup.PUT("/assignments/:id", instructorAssignmentHandlers.UpdateAssignment)
				instructorGroup.DELETE("/assignments/:id", instructorAssignmentHandlers.DeleteAssignment)
				instructorGroup.POST("/assignments/:id/assign", instructorAssignmentHandlers.AssignStudents)
				instructorGroup.GET("/assignments/:id/progress", instructorAssignmentHandlers.GetAssignmentProgress)
				instructorGroup.GET("/assignments/:id/students", instructorAssignmentHandlers.GetAssignmentStudents)
				instructorGroup.POST("/assignments/:id/students/:student_id/remove", instructorAssignmentHandlers.RemoveStudent)
				instructorGroup.GET("/students", instructorAssignmentHandlers.GetAllStudents)
				instructorGroup.GET("/dashboard/stats", instructorAssignmentHandlers.GetDashboardStats)
			}

			// Student assignment routes
			studentGroup := protected.Group("/student")
			studentGroup.Use(middleware.RequireRole("student"))
			{
				studentGroup.GET("/assignments", studentAssignmentHandlers.GetAssignments)
				studentGroup.GET("/assignments/:id", studentAssignmentHandlers.GetAssignment)
				studentGroup.POST("/assignments/:id/status", studentAssignmentHandlers.UpdateStatus)
				studentGroup.POST("/assignments/:id/complete", studentAssignmentHandlers.MarkAsCompleted)
				studentGroup.POST("/assignments/:id/progress", studentAssignmentHandlers.MarkAsInProgress)
				studentGroup.GET("/dashboard/stats", studentAssignmentHandlers.GetDashboardStats)
				studentGroup.GET("/assignments/overdue", studentAssignmentHandlers.GetOverdueAssignments)
				studentGroup.GET("/assignments/upcoming", studentAssignmentHandlers.GetUpcomingAssignments)
				studentGroup.GET("/assignments/recent", studentAssignmentHandlers.GetRecentlyCompleted)
				studentGroup.GET("/categories", studentAssignmentHandlers.GetCategories)
				studentGroup.GET("/assignments/by-status", studentAssignmentHandlers.GetAssignmentsByStatus)
				studentGroup.GET("/assignments/by-category", studentAssignmentHandlers.GetAssignmentsByCategory)
				studentGroup.GET("/assignments/search", studentAssignmentHandlers.SearchAssignments)
			}
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
