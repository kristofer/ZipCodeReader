package routes

import (
	"net/http"
	"zipcodereader/config"
	"zipcodereader/handlers"
	"zipcodereader/middleware"
	"zipcodereader/models"
	"zipcodereader/services"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Register sets up all routes for the application
func Register(r *gin.Engine, db *gorm.DB, cfg *config.Config) {
	// Common routes
	h := handlers.New(db)
	r.GET("/health", h.Health)
	r.GET("/", homeHandler(db, cfg))

	// Auth routes (mode-specific)
	registerAuthRoutes(r, db, cfg)

	// Protected routes (shared by both auth modes)
	registerProtectedRoutes(r, db, cfg)
}

// registerAuthRoutes sets up authentication routes based on auth mode
func registerAuthRoutes(r *gin.Engine, db *gorm.DB, cfg *config.Config) {
	if cfg.UseLocalAuth {
		// Local authentication routes
		local := handlers.NewLocalAuthHandler(db)
		r.GET("/local/login", local.ShowLogin)
		r.POST("/local/login", local.Login)
		r.GET("/local/register", local.ShowRegister)
		r.POST("/local/register", local.Register)
		r.GET("/local/logout", local.Logout)
	} else {
		// GitHub OAuth2 routes
		authService := services.NewAuthService(db, cfg)
		oauth := handlers.NewAuthHandler(authService)
		r.GET("/auth/login", oauth.Login)
		r.GET("/auth/callback", oauth.Callback)
		r.GET("/auth/logout", oauth.Logout)
	}
}

// registerProtectedRoutes sets up protected routes (requires authentication)
// These routes work for BOTH auth modes - no duplication!
func registerProtectedRoutes(r *gin.Engine, db *gorm.DB, cfg *config.Config) {
	protected := r.Group("/")
	protected.Use(middleware.RequireAuthWithUser(db))

	// Dashboard redirect
	protected.GET("/dashboard", dashboardRedirect)

	// Initialize services once
	assignmentService := services.NewAssignmentService(db)

	// Instructor routes
	registerInstructorRoutes(protected, assignmentService, cfg)

	// Student routes
	registerStudentRoutes(protected, assignmentService, cfg)
}

// registerInstructorRoutes sets up all instructor endpoints
func registerInstructorRoutes(rg *gin.RouterGroup, svc *services.AssignmentService, cfg *config.Config) {
	ig := rg.Group("/instructor")
	ig.Use(middleware.RequireRole("instructor"))

	// Initialize handlers
	h := handlers.NewInstructorHandlers(svc, cfg)

	// Dashboard & Views
	ig.GET("/dashboard", h.Dashboard)
	ig.GET("/assignments/manage", h.ManageAssignments)
	ig.GET("/assignments/:id/detail", h.AssignmentDetail)
	ig.GET("/assignments/:id/progress-view", h.ProgressView)
	ig.GET("/students/:username/assignments", h.ShowStudentAssignments)
	ig.GET("/students/:username/progress", h.GetStudentProgress)

	// Assignment CRUD API
	ig.GET("/assignments", h.ListAssignments)
	ig.POST("/assignments", h.CreateAssignment)
	ig.GET("/assignments/:id", h.GetAssignment)
	ig.PUT("/assignments/:id", h.UpdateAssignment)
	ig.DELETE("/assignments/:id", h.DeleteAssignment)

	// Student Management
	ig.POST("/assignments/:id/assign", h.AssignStudents)
	ig.GET("/assignments/:id/students", h.GetAssignmentStudents)
	ig.POST("/assignments/:id/students/:student_id/remove", h.RemoveStudent)
	ig.GET("/students", h.ListStudents)
	ig.POST("/students/:username/assignments/:assignment_id/assign", h.AssignToStudent)
	ig.DELETE("/students/:username/assignments/:assignment_id/remove", h.RemoveFromStudent)

	// Progress & Analytics
	ig.GET("/assignments/:id/progress", h.GetAssignmentProgress)
	ig.GET("/assignments/:id/detailed-progress", h.GetDetailedProgress)
	ig.GET("/progress/summary", h.GetProgressSummary)
	ig.GET("/progress/trends", h.GetProgressTrends)
	ig.GET("/progress/completion-analytics", h.GetCompletionAnalytics)

	// Due Dates & Notifications
	ig.GET("/due-dates/overview", h.GetDueDateOverview)
	ig.GET("/due-dates/notifications", h.GetNotifications)

	// Dashboard Stats
	ig.GET("/dashboard/stats", h.GetDashboardStats)
}

// registerStudentRoutes sets up all student endpoints
func registerStudentRoutes(rg *gin.RouterGroup, svc *services.AssignmentService, cfg *config.Config) {
	sg := rg.Group("/student")
	sg.Use(middleware.RequireRole("student"))

	// Initialize handlers
	h := handlers.NewStudentHandlers(svc, cfg)

	// Dashboard & Views
	sg.GET("/dashboard", h.Dashboard)
	sg.GET("/assignments/:id/detail", h.AssignmentDetail)

	// Assignment API
	sg.GET("/assignments", h.ListAssignments)
	sg.GET("/assignments/:id", h.GetAssignment)
	sg.POST("/assignments/:id/status", h.UpdateStatus)
	sg.POST("/assignments/:id/complete", h.MarkCompleted)
	sg.POST("/assignments/:id/progress", h.MarkInProgress)

	// Filtering & Search
	sg.GET("/assignments/overdue", h.GetOverdue)
	sg.GET("/assignments/upcoming", h.GetUpcoming)
	sg.GET("/assignments/recent", h.GetRecentlyCompleted)
	sg.GET("/assignments/status/:status", h.GetByStatus)
	sg.GET("/assignments/category/:category", h.GetByCategory)
	sg.GET("/assignments/search", h.Search)
	sg.GET("/categories", h.GetCategories)

	// Due Dates
	sg.GET("/due-dates/alerts", h.GetDueDateAlerts)
	sg.GET("/due-dates/summary", h.GetDueDateSummary)
	sg.GET("/due-dates/notifications", h.GetNotifications)

	// Dashboard Stats
	sg.GET("/dashboard/stats", h.GetDashboardStats)
}

// dashboardRedirect redirects users to their role-specific dashboard
func dashboardRedirect(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}

	userObj := user.(*models.User)
	if userObj.IsInstructor() {
		c.Redirect(http.StatusSeeOther, "/instructor/dashboard")
	} else {
		c.Redirect(http.StatusSeeOther, "/student/dashboard")
	}
}

// homeHandler renders the home page
func homeHandler(db *gorm.DB, cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		userID := session.Get("user_id")

		data := gin.H{
			"title":          "ZipCodeReader",
			"use_local_auth": cfg.UseLocalAuth,
		}

		if userID != nil {
			if user, err := models.GetUserByID(db, userID.(uint)); err == nil {
				data["user"] = user
			}
		}

		c.HTML(http.StatusOK, "index.html", data)
	}
}
