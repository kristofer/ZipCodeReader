package handlers

import (
	"net/http"
	"strconv"
	"strings"
	"time"
	"zipcodereader/models"
	"zipcodereader/services"

	"github.com/gin-gonic/gin"
)

// InstructorAssignmentHandlers handles instructor assignment operations
type InstructorAssignmentHandlers struct {
	assignmentService *services.AssignmentService
}

// NewInstructorAssignmentHandlers creates a new instructor assignment handlers
func NewInstructorAssignmentHandlers(assignmentService *services.AssignmentService) *InstructorAssignmentHandlers {
	return &InstructorAssignmentHandlers{
		assignmentService: assignmentService,
	}
}

// GetAssignments handles GET /instructor/assignments
func (h *InstructorAssignmentHandlers) GetAssignments(c *gin.Context) {
	// Get user from context (set by auth middleware)
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

	// Get query parameters for filtering
	category := c.Query("category")
	search := c.Query("search")

	var assignments []models.Assignment
	var err error

	if category != "" {
		assignments, err = h.assignmentService.GetAssignmentsByCategory(category, userObj.ID)
	} else if search != "" {
		assignments, err = h.assignmentService.SearchAssignments(search, userObj.ID)
	} else {
		assignments, err = h.assignmentService.GetAssignmentsByInstructor(userObj.ID)
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

// CreateAssignmentRequest represents the request body for creating an assignment
type CreateAssignmentRequest struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description"`
	URL         string `json:"url" binding:"required"`
	Category    string `json:"category"`
	DueDate     string `json:"due_date"` // ISO 8601 format
}

// CreateAssignment handles POST /instructor/assignments
func (h *InstructorAssignmentHandlers) CreateAssignment(c *gin.Context) {
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

	// Parse request body
	var req CreateAssignmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		println("Parse error:", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Debug: Log the parsed request
	println("Parsed request - Title:", req.Title, "URL:", req.URL, "Category:", req.Category)

	// Parse due date if provided
	var dueDate *time.Time
	if req.DueDate != "" {
		// Try different date formats
		var parsedDate time.Time
		var err error

		// First try RFC3339 format (ISO 8601)
		parsedDate, err = time.Parse(time.RFC3339, req.DueDate)
		if err != nil {
			// Try datetime-local format (YYYY-MM-DDTHH:MM) - assume local timezone
			parsedDate, err = time.ParseInLocation("2006-01-02T15:04", req.DueDate, time.Local)
			if err != nil {
				// Try date only format - assume local timezone, end of day
				parsedDate, err = time.ParseInLocation("2006-01-02", req.DueDate, time.Local)
				if err != nil {
					println("Due date parsing failed for:", req.DueDate)
					c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid due date format. Expected YYYY-MM-DD or YYYY-MM-DDTHH:MM"})
					return
				}
				// If it's just a date, set it to end of day (23:59:59)
				parsedDate = time.Date(parsedDate.Year(), parsedDate.Month(), parsedDate.Day(), 23, 59, 59, 0, parsedDate.Location())
			}
		}
		dueDate = &parsedDate
	}

	// Create assignment
	input := services.CreateAssignmentInput{
		Title:       req.Title,
		Description: req.Description,
		URL:         req.URL,
		Category:    req.Category,
		DueDate:     dueDate,
	}

	assignment, err := h.assignmentService.CreateAssignment(userObj.ID, input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":    "Assignment created successfully",
		"assignment": assignment,
	})
}

// GetAssignment handles GET /instructor/assignments/:id
func (h *InstructorAssignmentHandlers) GetAssignment(c *gin.Context) {
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

	// Get assignment
	assignment, err := h.assignmentService.GetAssignmentByID(uint(id), userObj.ID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Assignment not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"assignment": assignment,
	})
}

// UpdateAssignmentRequest represents the request body for updating an assignment
type UpdateAssignmentRequest struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description"`
	URL         string `json:"url" binding:"required"`
	Category    string `json:"category"`
	DueDate     string `json:"due_date"` // ISO 8601 format
}

// UpdateAssignment handles PUT /instructor/assignments/:id
func (h *InstructorAssignmentHandlers) UpdateAssignment(c *gin.Context) {
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

	// Parse request body
	var req UpdateAssignmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Parse due date if provided
	var dueDate *time.Time
	if req.DueDate != "" {
		// Try different date formats
		var parsedDate time.Time
		var err error

		// First try RFC3339 format (ISO 8601)
		parsedDate, err = time.Parse(time.RFC3339, req.DueDate)
		if err != nil {
			// Try datetime-local format (YYYY-MM-DDTHH:MM) - assume local timezone
			parsedDate, err = time.ParseInLocation("2006-01-02T15:04", req.DueDate, time.Local)
			if err != nil {
				// Try date only format - assume local timezone, end of day
				parsedDate, err = time.ParseInLocation("2006-01-02", req.DueDate, time.Local)
				if err != nil {
					println("Due date parsing failed for:", req.DueDate)
					c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid due date format. Expected YYYY-MM-DD or YYYY-MM-DDTHH:MM"})
					return
				}
				// If it's just a date, set it to end of day (23:59:59)
				parsedDate = time.Date(parsedDate.Year(), parsedDate.Month(), parsedDate.Day(), 23, 59, 59, 0, parsedDate.Location())
			}
		}
		dueDate = &parsedDate
	}

	// Update assignment
	input := services.UpdateAssignmentInput{
		Title:       req.Title,
		Description: req.Description,
		URL:         req.URL,
		Category:    req.Category,
		DueDate:     dueDate,
	}

	err = h.assignmentService.UpdateAssignment(uint(id), userObj.ID, input)
	if err != nil {
		if strings.Contains(err.Error(), "access denied") {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Assignment updated successfully",
	})
}

// DeleteAssignment handles DELETE /instructor/assignments/:id
func (h *InstructorAssignmentHandlers) DeleteAssignment(c *gin.Context) {
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

	// Delete assignment
	err = h.assignmentService.DeleteAssignment(uint(id), userObj.ID)
	if err != nil {
		if strings.Contains(err.Error(), "access denied") {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Assignment deleted successfully",
	})
}

// AssignStudentsRequest represents the request body for assigning students
type AssignStudentsRequest struct {
	StudentIDs []uint `json:"student_ids" binding:"required"`
}

// AssignStudents handles POST /instructor/assignments/:id/assign
func (h *InstructorAssignmentHandlers) AssignStudents(c *gin.Context) {
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

	// Parse request body
	var req AssignStudentsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if len(req.StudentIDs) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "At least one student ID is required"})
		return
	}

	// Assign students to assignment
	err = h.assignmentService.AssignToMultipleStudents(uint(id), req.StudentIDs, userObj.ID)
	if err != nil {
		if strings.Contains(err.Error(), "access denied") {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Students assigned successfully",
	})
}

// GetAssignmentProgress handles GET /instructor/assignments/:id/progress
func (h *InstructorAssignmentHandlers) GetAssignmentProgress(c *gin.Context) {
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

	// Get progress
	progress, err := h.assignmentService.GetAssignmentProgress(uint(id), userObj.ID)
	if err != nil {
		if strings.Contains(err.Error(), "access denied") {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Calculate percentages
	total := progress[models.StatusAssigned] + progress[models.StatusInProgress] + progress[models.StatusCompleted]
	var percentages map[string]float64
	if total > 0 {
		percentages = map[string]float64{
			models.StatusAssigned:   float64(progress[models.StatusAssigned]) / float64(total) * 100,
			models.StatusInProgress: float64(progress[models.StatusInProgress]) / float64(total) * 100,
			models.StatusCompleted:  float64(progress[models.StatusCompleted]) / float64(total) * 100,
		}
	} else {
		percentages = map[string]float64{
			models.StatusAssigned:   0,
			models.StatusInProgress: 0,
			models.StatusCompleted:  0,
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"progress":    progress,
		"percentages": percentages,
		"total":       total,
	})
}

// GetAssignmentStudents handles GET /instructor/assignments/:id/students
func (h *InstructorAssignmentHandlers) GetAssignmentStudents(c *gin.Context) {
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

	// Get assigned students
	students, err := h.assignmentService.GetAssignmentStudents(uint(id), userObj.ID)
	if err != nil {
		if strings.Contains(err.Error(), "access denied") {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"students": students,
		"total":    len(students),
	})
}

// RemoveStudentRequest represents the request body for removing a student
type RemoveStudentRequest struct {
	StudentID uint `json:"student_id" binding:"required"`
}

// RemoveStudent handles DELETE /instructor/assignments/:id/students
func (h *InstructorAssignmentHandlers) RemoveStudent(c *gin.Context) {
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

	// Parse request body
	var req RemoveStudentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Remove student from assignment
	err = h.assignmentService.RemoveStudentAssignment(uint(id), req.StudentID, userObj.ID)
	if err != nil {
		if strings.Contains(err.Error(), "access denied") {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Student removed from assignment successfully",
	})
}

// GetAllStudents handles GET /instructor/students
func (h *InstructorAssignmentHandlers) GetAllStudents(c *gin.Context) {
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

	// Get all students
	students, err := h.assignmentService.GetAllStudents(userObj.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"students": students,
		"total":    len(students),
	})
}

// GetDashboardStats handles GET /instructor/dashboard/stats
func (h *InstructorAssignmentHandlers) GetDashboardStats(c *gin.Context) {
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

	// Get all assignments for instructor
	assignments, err := h.assignmentService.GetAssignmentsByInstructor(userObj.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Get all students
	students, err := models.GetAllStudents(h.assignmentService.GetDB())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Calculate overall statistics
	totalAssignments := len(assignments)
	totalStudentAssignments := 0
	totalCompleted := 0
	totalInProgress := 0
	totalAssigned := 0
	overdueCount := 0

	for _, assignment := range assignments {
		progress, err := h.assignmentService.GetAssignmentProgress(assignment.ID, userObj.ID)
		if err != nil {
			continue
		}

		assignmentTotal := progress[models.StatusAssigned] + progress[models.StatusInProgress] + progress[models.StatusCompleted]
		totalStudentAssignments += assignmentTotal
		totalAssigned += progress[models.StatusAssigned]
		totalInProgress += progress[models.StatusInProgress]
		totalCompleted += progress[models.StatusCompleted]

		// Check for overdue assignments
		if assignment.DueDate != nil && assignment.DueDate.Before(time.Now()) {
			overdueCount++
		}
	}

	// Calculate completion rate
	var completionRate float64
	if totalStudentAssignments > 0 {
		completionRate = float64(totalCompleted) / float64(totalStudentAssignments) * 100
	}

	c.JSON(http.StatusOK, gin.H{
		"total_assignments":         totalAssignments,
		"active_students":           len(students),
		"total_student_assignments": totalStudentAssignments,
		"total_assigned":            totalAssigned,
		"total_in_progress":         totalInProgress,
		"total_completed":           totalCompleted,
		"completion_rate":           completionRate,
		"overdue_count":             overdueCount,
		"students":                  students,
	})
}
