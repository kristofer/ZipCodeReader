package handlers

import (
	"net/http"
	"strconv"
	"strings"
	"zipcodereader/models"
	"zipcodereader/services"

	"github.com/gin-gonic/gin"
)

// StudentAssignmentHandlers handles student assignment operations
type StudentAssignmentHandlers struct {
	studentService *services.StudentAssignmentService
}

// NewStudentAssignmentHandlers creates a new student assignment handlers
func NewStudentAssignmentHandlers(studentService *services.StudentAssignmentService) *StudentAssignmentHandlers {
	return &StudentAssignmentHandlers{
		studentService: studentService,
	}
}

// GetAssignments handles GET /student/assignments
func (h *StudentAssignmentHandlers) GetAssignments(c *gin.Context) {
	// Get user from context (set by auth middleware)
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	userObj := user.(*models.User)
	if !userObj.IsStudent() {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	// Get query parameters for filtering
	status := c.Query("status")
	category := c.Query("category")
	search := c.Query("search")

	var assignments []models.StudentAssignment
	var err error

	if status != "" {
		assignments, err = h.studentService.GetStudentAssignmentsByStatus(userObj.ID, status)
	} else if category != "" {
		assignments, err = h.studentService.GetStudentAssignmentsByCategory(userObj.ID, category)
	} else if search != "" {
		assignments, err = h.studentService.SearchStudentAssignments(userObj.ID, search)
	} else {
		assignments, err = h.studentService.GetStudentAssignments(userObj.ID)
	}

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
func (h *StudentAssignmentHandlers) GetAssignment(c *gin.Context) {
	// Get user from context
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	userObj := user.(*models.User)
	if !userObj.IsStudent() {
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

	// Get student assignment
	assignment, err := h.studentService.GetStudentAssignment(uint(id), userObj.ID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			c.JSON(http.StatusNotFound, gin.H{"error": "Assignment not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"assignment": assignment,
	})
}

// UpdateStatusRequest represents the request body for updating assignment status
type UpdateStatusRequest struct {
	Status string `json:"status" binding:"required"`
}

// UpdateStatus handles POST /student/assignments/:id/status
func (h *StudentAssignmentHandlers) UpdateStatus(c *gin.Context) {
	// Get user from context
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	userObj := user.(*models.User)
	if !userObj.IsStudent() {
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

	// Parse request body
	var req UpdateStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update assignment status
	err = h.studentService.UpdateAssignmentStatus(uint(id), userObj.ID, req.Status)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			c.JSON(http.StatusNotFound, gin.H{"error": "Assignment not found"})
			return
		}
		if strings.Contains(err.Error(), "invalid status") {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid status"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Assignment status updated successfully",
	})
}

// MarkAsCompleted handles POST /student/assignments/:id/complete
func (h *StudentAssignmentHandlers) MarkAsCompleted(c *gin.Context) {
	// Get user from context
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	userObj := user.(*models.User)
	if !userObj.IsStudent() {
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

	// Mark assignment as completed
	err = h.studentService.MarkAsCompletedByID(uint(id), userObj.ID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			c.JSON(http.StatusNotFound, gin.H{"error": "Assignment not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Assignment marked as completed successfully",
	})
}

// MarkAsInProgress handles POST /student/assignments/:id/progress
func (h *StudentAssignmentHandlers) MarkAsInProgress(c *gin.Context) {
	// Get user from context
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	userObj := user.(*models.User)
	if !userObj.IsStudent() {
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

	// Mark assignment as in progress
	err = h.studentService.MarkAsInProgressByID(uint(id), userObj.ID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			c.JSON(http.StatusNotFound, gin.H{"error": "Assignment not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Assignment marked as in progress successfully",
	})
}

// GetDashboardStats handles GET /student/dashboard/stats
func (h *StudentAssignmentHandlers) GetDashboardStats(c *gin.Context) {
	// Get user from context
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	userObj := user.(*models.User)
	if !userObj.IsStudent() {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	// Get dashboard statistics
	stats, err := h.studentService.GetDashboardStats(userObj.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, stats)
}

// GetOverdueAssignments handles GET /student/assignments/overdue
func (h *StudentAssignmentHandlers) GetOverdueAssignments(c *gin.Context) {
	// Get user from context
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	userObj := user.(*models.User)
	if !userObj.IsStudent() {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	// Get overdue assignments
	assignments, err := h.studentService.GetOverdueAssignments(userObj.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"assignments": assignments,
		"total":       len(assignments),
	})
}

// GetUpcomingAssignments handles GET /student/assignments/upcoming
func (h *StudentAssignmentHandlers) GetUpcomingAssignments(c *gin.Context) {
	// Get user from context
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	userObj := user.(*models.User)
	if !userObj.IsStudent() {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	// Get days parameter (default to 7 days)
	daysStr := c.DefaultQuery("days", "7")
	days, err := strconv.Atoi(daysStr)
	if err != nil {
		days = 7
	}

	// Get upcoming assignments
	assignments, err := h.studentService.GetUpcomingAssignments(userObj.ID, days)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"assignments": assignments,
		"total":       len(assignments),
		"days":        days,
	})
}

// GetRecentlyCompleted handles GET /student/assignments/recent
func (h *StudentAssignmentHandlers) GetRecentlyCompleted(c *gin.Context) {
	// Get user from context
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	userObj := user.(*models.User)
	if !userObj.IsStudent() {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	// Get days parameter (default to 7 days)
	daysStr := c.DefaultQuery("days", "7")
	days, err := strconv.Atoi(daysStr)
	if err != nil {
		days = 7
	}

	// Get recently completed assignments
	assignments, err := h.studentService.GetRecentlyCompleted(userObj.ID, days)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"assignments": assignments,
		"total":       len(assignments),
		"days":        days,
	})
}

// GetCategories handles GET /student/assignments/categories
func (h *StudentAssignmentHandlers) GetCategories(c *gin.Context) {
	// Get user from context
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	userObj := user.(*models.User)
	if !userObj.IsStudent() {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	// Get assignment categories
	categories, err := h.studentService.GetAssignmentCategories(userObj.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"categories": categories,
		"total":      len(categories),
	})
}

// GetAssignmentsByStatus handles GET /student/assignments/status/:status
func (h *StudentAssignmentHandlers) GetAssignmentsByStatus(c *gin.Context) {
	// Get user from context
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	userObj := user.(*models.User)
	if !userObj.IsStudent() {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	// Get status from URL
	status := c.Param("status")

	// Get assignments by status
	assignments, err := h.studentService.GetStudentAssignmentsByStatus(userObj.ID, status)
	if err != nil {
		if strings.Contains(err.Error(), "invalid status") {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid status"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"assignments": assignments,
		"total":       len(assignments),
		"status":      status,
	})
}

// GetAssignmentsByCategory handles GET /student/assignments/category/:category
func (h *StudentAssignmentHandlers) GetAssignmentsByCategory(c *gin.Context) {
	// Get user from context
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	userObj := user.(*models.User)
	if !userObj.IsStudent() {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	// Get category from URL
	category := c.Param("category")

	// Get assignments by category
	assignments, err := h.studentService.GetStudentAssignmentsByCategory(userObj.ID, category)
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

// SearchAssignments handles GET /student/assignments/search
func (h *StudentAssignmentHandlers) SearchAssignments(c *gin.Context) {
	// Get user from context
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	userObj := user.(*models.User)
	if !userObj.IsStudent() {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	// Get search query
	query := c.Query("q")
	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Search query is required"})
		return
	}

	// Search assignments
	assignments, err := h.studentService.SearchStudentAssignments(userObj.ID, query)
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
