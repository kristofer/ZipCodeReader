package handlers

import (
	"net/http"
	"strconv"
	"zipcodereader/config"
	"zipcodereader/models"
	"zipcodereader/services"

	"github.com/gin-gonic/gin"
)

// StudentHandlers consolidates all student-related handlers
// Replaces: StudentAssignmentHandlers and DueDateNotificationHandlers (student portion)
type StudentHandlers struct {
	assignmentService *services.AssignmentService
	useLocalAuth      bool
}

// NewStudentHandlers creates a new consolidated student handler
func NewStudentHandlers(assignmentService *services.AssignmentService, cfg *config.Config) *StudentHandlers {
	return &StudentHandlers{
		assignmentService: assignmentService,
		useLocalAuth:      cfg.UseLocalAuth,
	}
}

// ============================================================================
// DASHBOARD & VIEWS
// ============================================================================

// Dashboard renders the student dashboard
func (h *StudentHandlers) Dashboard(c *gin.Context) {
	user, _ := c.Get("user")
	userObj := user.(*models.User)

	c.HTML(http.StatusOK, "student_assignments.html", gin.H{
		"title":          "My Assignments",
		"user":           userObj,
		"use_local_auth": h.useLocalAuth,
	})
}

// AssignmentDetail renders the assignment detail page for students
func (h *StudentHandlers) AssignmentDetail(c *gin.Context) {
	user, _ := c.Get("user")
	userObj := user.(*models.User)

	assignmentID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.HTML(http.StatusBadRequest, "error.html", gin.H{"error": "Invalid assignment ID"})
		return
	}

	studentAssignment, err := h.assignmentService.GetStudentAssignmentByID(uint(assignmentID), userObj.ID)
	if err != nil {
		c.HTML(http.StatusNotFound, "error.html", gin.H{"error": "Assignment not found"})
		return
	}

	c.HTML(http.StatusOK, "assignment_detail.html", gin.H{
		"title":               studentAssignment.Assignment.Title,
		"user":                userObj,
		"assignment":          studentAssignment.Assignment,
		"student_assignment":  studentAssignment,
		"use_local_auth":      h.useLocalAuth,
	})
}

// ============================================================================
// ASSIGNMENT API
// ============================================================================

// ListAssignments handles GET /student/assignments
func (h *StudentHandlers) ListAssignments(c *gin.Context) {
	user, _ := c.Get("user")
	userObj := user.(*models.User)

	if !userObj.IsStudent() {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	assignments, err := h.assignmentService.GetStudentAssignments(userObj.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"assignments": assignments,
		"total":       len(assignments),
	})
}

// GetAssignment handles GET /student/assignments/:id
func (h *StudentHandlers) GetAssignment(c *gin.Context) {
	user, _ := c.Get("user")
	userObj := user.(*models.User)

	assignmentID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid assignment ID"})
		return
	}

	studentAssignment, err := h.assignmentService.GetStudentAssignmentByID(uint(assignmentID), userObj.ID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"assignment": studentAssignment})
}

// UpdateStatusRequest represents the request to update assignment status
type UpdateStatusRequest struct {
	Status string `json:"status" binding:"required"`
}

// UpdateStatus handles POST /student/assignments/:id/status
func (h *StudentHandlers) UpdateStatus(c *gin.Context) {
	user, _ := c.Get("user")
	userObj := user.(*models.User)

	assignmentID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid assignment ID"})
		return
	}

	var req UpdateStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.assignmentService.UpdateStudentAssignmentStatus(uint(assignmentID), userObj.ID, req.Status); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Status updated successfully"})
}

// MarkCompleted handles POST /student/assignments/:id/complete
func (h *StudentHandlers) MarkCompleted(c *gin.Context) {
	user, _ := c.Get("user")
	userObj := user.(*models.User)

	assignmentID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid assignment ID"})
		return
	}

	if err := h.assignmentService.MarkAsCompleted(uint(assignmentID), userObj.ID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Assignment marked as completed"})
}

// MarkInProgress handles POST /student/assignments/:id/progress
func (h *StudentHandlers) MarkInProgress(c *gin.Context) {
	user, _ := c.Get("user")
	userObj := user.(*models.User)

	assignmentID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid assignment ID"})
		return
	}

	if err := h.assignmentService.MarkAsInProgress(uint(assignmentID), userObj.ID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Assignment marked as in progress"})
}

// ============================================================================
// FILTERING & SEARCH
// ============================================================================

// GetOverdue handles GET /student/assignments/overdue
func (h *StudentHandlers) GetOverdue(c *gin.Context) {
	user, _ := c.Get("user")
	userObj := user.(*models.User)

	assignments, err := h.assignmentService.GetOverdueAssignments(userObj.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"assignments": assignments,
		"total":       len(assignments),
	})
}

// GetUpcoming handles GET /student/assignments/upcoming
func (h *StudentHandlers) GetUpcoming(c *gin.Context) {
	user, _ := c.Get("user")
	userObj := user.(*models.User)

	days := 7
	if daysStr := c.Query("days"); daysStr != "" {
		if d, err := strconv.Atoi(daysStr); err == nil {
			days = d
		}
	}

	assignments, err := h.assignmentService.GetUpcomingAssignments(userObj.ID, days)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"assignments": assignments,
		"total":       len(assignments),
	})
}

// GetRecentlyCompleted handles GET /student/assignments/recent
func (h *StudentHandlers) GetRecentlyCompleted(c *gin.Context) {
	user, _ := c.Get("user")
	userObj := user.(*models.User)

	limit := 10
	if limitStr := c.Query("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil {
			limit = l
		}
	}

	assignments, err := h.assignmentService.GetRecentlyCompleted(userObj.ID, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"assignments": assignments,
		"total":       len(assignments),
	})
}

// GetByStatus handles GET /student/assignments/status/:status
func (h *StudentHandlers) GetByStatus(c *gin.Context) {
	user, _ := c.Get("user")
	userObj := user.(*models.User)

	status := c.Param("status")

	assignments, err := h.assignmentService.GetAssignmentsByStatus(userObj.ID, status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"assignments": assignments,
		"total":       len(assignments),
		"status":      status,
	})
}

// GetByCategory handles GET /student/assignments/category/:category
func (h *StudentHandlers) GetByCategory(c *gin.Context) {
	user, _ := c.Get("user")
	userObj := user.(*models.User)

	category := c.Param("category")

	assignments, err := h.assignmentService.GetStudentAssignmentsByCategory(userObj.ID, category)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"assignments": assignments,
		"total":       len(assignments),
		"category":    category,
	})
}

// Search handles GET /student/assignments/search
func (h *StudentHandlers) Search(c *gin.Context) {
	user, _ := c.Get("user")
	userObj := user.(*models.User)

	query := c.Query("q")
	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Search query required"})
		return
	}

	assignments, err := h.assignmentService.SearchStudentAssignments(userObj.ID, query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"assignments": assignments,
		"total":       len(assignments),
		"query":       query,
	})
}

// GetCategories handles GET /student/categories
func (h *StudentHandlers) GetCategories(c *gin.Context) {
	user, _ := c.Get("user")
	userObj := user.(*models.User)

	categories, err := h.assignmentService.GetStudentCategories(userObj.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"categories": categories,
		"total":      len(categories),
	})
}

// ============================================================================
// DUE DATES (Absorbed from DueDateNotificationHandlers)
// ============================================================================

// GetDueDateAlerts handles GET /student/due-dates/alerts
func (h *StudentHandlers) GetDueDateAlerts(c *gin.Context) {
	user, _ := c.Get("user")
	userObj := user.(*models.User)

	days := 7
	if daysStr := c.Query("days"); daysStr != "" {
		if d, err := strconv.Atoi(daysStr); err == nil {
			days = d
		}
	}

	alerts, err := h.assignmentService.GetStudentDueDateAlerts(userObj.ID, days)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"alerts": alerts,
		"total":  len(alerts),
	})
}

// GetDueDateSummary handles GET /student/due-dates/summary
func (h *StudentHandlers) GetDueDateSummary(c *gin.Context) {
	user, _ := c.Get("user")
	userObj := user.(*models.User)

	summary, err := h.assignmentService.GetStudentDueDateSummary(userObj.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, summary)
}

// GetNotifications handles GET /student/due-dates/notifications
func (h *StudentHandlers) GetNotifications(c *gin.Context) {
	user, _ := c.Get("user")
	userObj := user.(*models.User)

	notifications, err := h.assignmentService.GetStudentDueDateNotifications(userObj.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"notifications": notifications,
		"total":         len(notifications),
	})
}

// ============================================================================
// DASHBOARD STATS
// ============================================================================

// GetDashboardStats handles GET /student/dashboard/stats
func (h *StudentHandlers) GetDashboardStats(c *gin.Context) {
	user, _ := c.Get("user")
	userObj := user.(*models.User)

	if !userObj.IsStudent() {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	assignments, err := h.assignmentService.GetStudentAssignments(userObj.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get assignments"})
		return
	}

	stats := h.assignmentService.CalculateStudentStats(assignments)

	c.JSON(http.StatusOK, stats)
}
