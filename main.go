package main

import (
	"flag"
	"log"
	"net/http"

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
	progressTrackingService := services.NewProgressTrackingService(db)
	dueDateNotificationService := services.NewDueDateNotificationService(db)

	// Initialize assignment handlers
	instructorAssignmentHandlers := handlers.NewInstructorAssignmentHandlers(assignmentService)
	studentAssignmentHandlers := handlers.NewStudentAssignmentHandlers(studentAssignmentService)
	progressTrackingHandlers := handlers.NewProgressTrackingHandlers(progressTrackingService)
	dueDateNotificationHandlers := handlers.NewDueDateNotificationHandlers(dueDateNotificationService)
	dashboardHandlers := handlers.NewDashboardHandlers(assignmentService, studentAssignmentService, cfg.UseLocalAuth)

	// Setup authentication routes based on mode
	if cfg.UseLocalAuth {
		log.Println("Using local authentication mode (default)")
		localAuthHandler := handlers.NewLocalAuthHandler(db)

		// Local authentication routes
		r.GET("/local/login", localAuthHandler.ShowLogin)
		r.POST("/local/login", localAuthHandler.Login)
		r.GET("/local/register", localAuthHandler.ShowRegister)
		r.POST("/local/register", localAuthHandler.Register)
		r.GET("/local/logout", localAuthHandler.Logout)

		// Dashboard route - redirects to appropriate dashboard based on user role
		protected := r.Group("/")
		protected.Use(middleware.RequireAuthWithUser(db))
		{
			protected.GET("/dashboard", func(c *gin.Context) {
				user, exists := c.Get("user")
				if !exists {
					c.JSON(500, gin.H{"error": "User not found"})
					return
				}

				userObj := user.(*models.User)
				if userObj.IsInstructor() {
					c.Redirect(http.StatusSeeOther, "/instructor/dashboard")
				} else {
					c.Redirect(http.StatusSeeOther, "/student/dashboard")
				}
			})

			// Instructor assignment routes
			instructorGroup := protected.Group("/instructor")
			instructorGroup.Use(middleware.RequireRole("instructor"))
			{
				// Dashboard routes
				instructorGroup.GET("/dashboard", dashboardHandlers.ShowInstructorDashboard)
				instructorGroup.GET("/assignments/:id/detail", dashboardHandlers.ShowAssignmentDetail)
				instructorGroup.GET("/assignments/:id/progress-view", dashboardHandlers.ShowAssignmentProgress)

				// API routes
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

				// Advanced progress tracking routes
				instructorGroup.GET("/assignments/:id/detailed-progress", progressTrackingHandlers.GetDetailedProgressReport)
				instructorGroup.GET("/progress/summary", progressTrackingHandlers.GetInstructorProgressSummary)
				instructorGroup.GET("/progress/trends", progressTrackingHandlers.GetProgressTrends)
				instructorGroup.GET("/progress/completion-analytics", progressTrackingHandlers.GetCompletionAnalytics)

				// Due date notification routes for instructors
				instructorGroup.GET("/due-dates/overview", dueDateNotificationHandlers.GetInstructorDueDateOverview)
				instructorGroup.GET("/due-dates/notifications", dueDateNotificationHandlers.GetDueDateNotifications)
			}

			// Student assignment routes
			studentGroup := protected.Group("/student")
			studentGroup.Use(middleware.RequireRole("student"))
			{
				// Dashboard routes
				studentGroup.GET("/dashboard", dashboardHandlers.ShowStudentDashboard)
				studentGroup.GET("/assignments/:id/detail", dashboardHandlers.ShowAssignmentDetail)

				// API routes
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

				// Due date notification routes for students
				studentGroup.GET("/due-dates/alerts", dueDateNotificationHandlers.GetStudentDueDateAlerts)
				studentGroup.GET("/due-dates/summary", dueDateNotificationHandlers.GetStudentDueDateSummary)
				studentGroup.GET("/due-dates/notifications", dueDateNotificationHandlers.GetDueDateNotifications)
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
		log.Println("Using GitHub OAuth2 authentication mode (optional)")
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
				// Dashboard routes
				instructorGroup.GET("/dashboard", dashboardHandlers.ShowInstructorDashboard)
				instructorGroup.GET("/assignments/:id/detail", dashboardHandlers.ShowAssignmentDetail)
				instructorGroup.GET("/assignments/:id/progress-view", dashboardHandlers.ShowAssignmentProgress)
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

				// Advanced progress tracking routes
				instructorGroup.GET("/assignments/:id/detailed-progress", progressTrackingHandlers.GetDetailedProgressReport)
				instructorGroup.GET("/progress/summary", progressTrackingHandlers.GetInstructorProgressSummary)
				instructorGroup.GET("/progress/trends", progressTrackingHandlers.GetProgressTrends)
				instructorGroup.GET("/progress/completion-analytics", progressTrackingHandlers.GetCompletionAnalytics)

				// Due date notification routes for instructors
				instructorGroup.GET("/due-dates/overview", dueDateNotificationHandlers.GetInstructorDueDateOverview)
				instructorGroup.GET("/due-dates/notifications", dueDateNotificationHandlers.GetDueDateNotifications)
			}

			// Student assignment routes
			studentGroup := protected.Group("/student")
			studentGroup.Use(middleware.RequireRole("student"))
			{
				// Dashboard routes
				studentGroup.GET("/dashboard", dashboardHandlers.ShowStudentDashboard)
				studentGroup.GET("/assignments/:id/detail", dashboardHandlers.ShowAssignmentDetail)

				// API routes
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

				// Due date notification routes for students
				studentGroup.GET("/due-dates/alerts", dueDateNotificationHandlers.GetStudentDueDateAlerts)
				studentGroup.GET("/due-dates/summary", dueDateNotificationHandlers.GetStudentDueDateSummary)
				studentGroup.GET("/due-dates/notifications", dueDateNotificationHandlers.GetDueDateNotifications)
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
			return "Local (default)"
		}
		return "GitHub OAuth2 (optional)"
	}())

	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
