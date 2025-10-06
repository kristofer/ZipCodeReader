package handlers

import (
	"net/http"
	"strconv"
	"zipcodereader/models"
	"zipcodereader/services"

	"github.com/gin-gonic/gin"
)

// ProgressTrackingHandlers handles enhanced progress tracking operations
type ProgressTrackingHandlers struct {
	progressService *services.ProgressTrackingService
}

// NewProgressTrackingHandlers creates new progress tracking handlers
func NewProgressTrackingHandlers(progressService *services.ProgressTrackingService) *ProgressTrackingHandlers {
	return &ProgressTrackingHandlers{
		progressService: progressService,
	}
}

// GetDetailedProgressReport handles GET /instructor/assignments/:id/detailed-progress
func (h *ProgressTrackingHandlers) GetDetailedProgressReport(c *gin.Context) {
	// Get user from context
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	userObj := user.(*models.User)
	if !userObj.IsInstructor() {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	// Get assignment ID from URL
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid assignment ID"})
		return
	}

	// Get detailed progress report
	report, err := h.progressService.GetDetailedProgressReport(uint(id), userObj.ID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"report": report,
	})
}

// GetInstructorProgressSummary handles GET /instructor/progress/summary
func (h *ProgressTrackingHandlers) GetInstructorProgressSummary(c *gin.Context) {
	// Get user from context
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	userObj := user.(*models.User)
	if !userObj.IsInstructor() {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	// Get progress summary
	summary, err := h.progressService.GetInstructorProgressSummary(userObj.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"summary": summary,
	})
}

// GetProgressTrends handles GET /instructor/progress/trends
func (h *ProgressTrackingHandlers) GetProgressTrends(c *gin.Context) {
	// Get user from context
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	userObj := user.(*models.User)
	if !userObj.IsInstructor() {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	// Get query parameters
	period := c.DefaultQuery("period", "30")              // days
	granularity := c.DefaultQuery("granularity", "daily") // daily, weekly, monthly

	// For now, return a placeholder response
	// This can be expanded to include actual trend calculations
	c.JSON(http.StatusOK, gin.H{
		"trends": gin.H{
			"period":      period,
			"granularity": granularity,
			"message":     "Progress trends feature coming soon",
		},
	})
}

// GetCompletionAnalytics handles GET /instructor/progress/completion-analytics
func (h *ProgressTrackingHandlers) GetCompletionAnalytics(c *gin.Context) {
	// Get user from context
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	userObj := user.(*models.User)
	if !userObj.IsInstructor() {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	// Get progress summary (contains completion analytics)
	summary, err := h.progressService.GetInstructorProgressSummary(userObj.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Extract completion analytics from summary
	analytics := gin.H{
		"overall_completion_rate":   summary.OverallCompletionRate,
		"average_completion_time":   summary.AverageCompletionTime,
		"total_assignments":         summary.TotalAssignments,
		"total_student_assignments": summary.TotalStudentAssignments,
		"overdue_assignments":       summary.OverdueAssignments,
		"category_breakdown":        summary.CategoryBreakdown,
		"recent_completions":        summary.RecentCompletions,
		"student_engagement":        summary.StudentEngagement,
	}

	c.JSON(http.StatusOK, gin.H{
		"analytics": analytics,
	})
}
