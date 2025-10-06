package handlers

import (
	"net/http"
	"strconv"
	"time"
	"zipcodereader/config"
	"zipcodereader/models"
	"zipcodereader/services"

	"github.com/gin-gonic/gin"
)

// InstructorHandlers consolidates all instructor-related handlers
// Replaces: InstructorAssignmentHandlers, ProgressTrackingHandlers,
// DueDateNotificationHandlers, and DashboardHandlers (instructor portion)
type InstructorHandlers struct {
	assignmentService *services.AssignmentService
	useLocalAuth      bool
}

// NewInstructorHandlers creates a new consolidated instructor handler
func NewInstructorHandlers(assignmentService *services.AssignmentService, cfg *config.Config) *InstructorHandlers {
	return &InstructorHandlers{
		assignmentService: assignmentService,
		useLocalAuth:      cfg.UseLocalAuth,
	}
}

// ============================================================================
// DASHBOARD & VIEWS
// ============================================================================

// Dashboard renders the instructor dashboard
func (h *InstructorHandlers) Dashboard(c *gin.Context) {
	user, _ := c.Get("user")
	userObj := user.(*models.User)

	c.HTML(http.StatusOK, "instructor_assignments.html", gin.H{
		"title":          "Instructor Dashboard",
		"user":           userObj,
		"use_local_auth": h.useLocalAuth,
		"template_type":  "instructor",
	})
}

// ManageAssignments renders the assignment management page
func (h *InstructorHandlers) ManageAssignments(c *gin.Context) {
	user, _ := c.Get("user")
	userObj := user.(*models.User)

	c.HTML(http.StatusOK, "assignment_management.html", gin.H{
		"title":          "Manage Assignments",
		"user":           userObj,
		"use_local_auth": h.useLocalAuth,
	})
}

// AssignmentDetail renders the assignment detail page
func (h *InstructorHandlers) AssignmentDetail(c *gin.Context) {
	user, _ := c.Get("user")
	userObj := user.(*models.User)

	assignmentID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.HTML(http.StatusBadRequest, "error.html", gin.H{"error": "Invalid assignment ID"})
		return
	}

	assignment, err := h.assignmentService.GetAssignmentByID(uint(assignmentID), userObj.ID)
	if err != nil {
		c.HTML(http.StatusNotFound, "error.html", gin.H{"error": "Assignment not found"})
		return
	}

	c.HTML(http.StatusOK, "assignment_detail.html", gin.H{
		"title":          assignment.Title,
		"user":           userObj,
		"assignment":     assignment,
		"use_local_auth": h.useLocalAuth,
	})
}

// ProgressView renders the assignment progress page
func (h *InstructorHandlers) ProgressView(c *gin.Context) {
	user, _ := c.Get("user")
	userObj := user.(*models.User)

	assignmentID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.HTML(http.StatusBadRequest, "error.html", gin.H{"error": "Invalid assignment ID"})
		return
	}

	c.HTML(http.StatusOK, "assignment_progress.html", gin.H{
		"title":          "Assignment Progress",
		"user":           userObj,
		"assignment_id":  assignmentID,
		"use_local_auth": h.useLocalAuth,
	})
}

// ShowStudentAssignments shows assignment management page for a specific student
func (h *InstructorHandlers) ShowStudentAssignments(c *gin.Context) {
	user, _ := c.Get("user")
	userObj := user.(*models.User)

	username := c.Param("username")
	student, err := models.GetUserByUsername(h.assignmentService.GetDB(), username)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Student not found"})
		return
	}

	if !student.IsStudent() {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User is not a student"})
		return
	}

	// Get all instructor's assignments
	assignments, err := h.assignmentService.GetAssignmentsByInstructor(userObj.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Get student's assigned readings
	studentAssignments, err := h.assignmentService.GetStudentAssignments(student.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Create a map for quick lookup
	assignedMap := make(map[uint]*models.StudentAssignment)
	for i := range studentAssignments {
		assignedMap[studentAssignments[i].AssignmentID] = &studentAssignments[i]
	}

	c.HTML(http.StatusOK, "student_assignment_management.html", gin.H{
		"title":              "Assign Readings to " + student.Username,
		"user":               userObj,
		"student":            student,
		"assignments":        assignments,
		"assigned_map":       assignedMap,
		"use_local_auth":     h.useLocalAuth,
		"template_type":      "student_assignment",
	})
}

// GetStudentProgress shows detailed progress for a specific student
func (h *InstructorHandlers) GetStudentProgress(c *gin.Context) {
	user, _ := c.Get("user")
	userObj := user.(*models.User)

	username := c.Param("username")
	student, err := models.GetUserByUsername(h.assignmentService.GetDB(), username)
	if err != nil {
		if c.GetHeader("Accept") == "application/json" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Student not found"})
		} else {
			c.HTML(http.StatusNotFound, "error.html", gin.H{"error": "Student not found"})
		}
		return
	}

	if !student.IsStudent() {
		if c.GetHeader("Accept") == "application/json" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "User is not a student"})
		} else {
			c.HTML(http.StatusBadRequest, "error.html", gin.H{"error": "User is not a student"})
		}
		return
	}

	// Get student assignments
	assignments, err := h.assignmentService.GetStudentAssignments(student.ID)
	if err != nil {
		if c.GetHeader("Accept") == "application/json" {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		} else {
			c.HTML(http.StatusInternalServerError, "error.html", gin.H{"error": err.Error()})
		}
		return
	}

	// Calculate statistics
	stats := h.assignmentService.CalculateStudentStats(assignments)

	if c.GetHeader("Accept") == "application/json" {
		c.JSON(http.StatusOK, gin.H{
			"student":     student,
			"assignments": assignments,
			"statistics":  stats,
		})
	} else {
		c.HTML(http.StatusOK, "student_progress.html", gin.H{
			"title":          "Progress for " + student.Username,
			"user":           userObj,
			"student":        student,
			"assignments":    assignments,
			"statistics":     stats,
			"use_local_auth": h.useLocalAuth,
		})
	}
}

// ============================================================================
// ASSIGNMENT CRUD
// ============================================================================

// ListAssignments handles GET /instructor/assignments
func (h *InstructorHandlers) ListAssignments(c *gin.Context) {
	user, _ := c.Get("user")
	userObj := user.(*models.User)

	if !userObj.IsInstructor() {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

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
	DueDate     string `json:"due_date"`
}

// CreateAssignment handles POST /instructor/assignments
func (h *InstructorHandlers) CreateAssignment(c *gin.Context) {
	user, _ := c.Get("user")
	userObj := user.(*models.User)

	if !userObj.IsInstructor() {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	var req CreateAssignmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var dueDate *time.Time
	if req.DueDate != "" {
		parsed, err := time.Parse(time.RFC3339, req.DueDate)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid due date format"})
			return
		}
		dueDate = &parsed
	}

	input := services.CreateAssignmentInput{
		Title:       req.Title,
		Description: req.Description,
		URL:         req.URL,
		Category:    req.Category,
		DueDate:     dueDate,
	}

	assignment, err := h.assignmentService.CreateAssignment(userObj.ID, input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"assignment": assignment,
		"message":    "Assignment created successfully",
	})
}

// GetAssignment handles GET /instructor/assignments/:id
func (h *InstructorHandlers) GetAssignment(c *gin.Context) {
	user, _ := c.Get("user")
	userObj := user.(*models.User)

	assignmentID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid assignment ID"})
		return
	}

	assignment, err := h.assignmentService.GetAssignmentByID(uint(assignmentID), userObj.ID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"assignment": assignment})
}

// UpdateAssignment handles PUT /instructor/assignments/:id
func (h *InstructorHandlers) UpdateAssignment(c *gin.Context) {
	user, _ := c.Get("user")
	userObj := user.(*models.User)

	assignmentID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid assignment ID"})
		return
	}

	var req CreateAssignmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var dueDate *time.Time
	if req.DueDate != "" {
		parsed, err := time.Parse(time.RFC3339, req.DueDate)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid due date format"})
			return
		}
		dueDate = &parsed
	}

	input := services.UpdateAssignmentInput{
		Title:       req.Title,
		Description: req.Description,
		URL:         req.URL,
		Category:    req.Category,
		DueDate:     dueDate,
	}

	assignment, err := h.assignmentService.UpdateAssignment(uint(assignmentID), userObj.ID, input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"assignment": assignment,
		"message":    "Assignment updated successfully",
	})
}

// DeleteAssignment handles DELETE /instructor/assignments/:id
func (h *InstructorHandlers) DeleteAssignment(c *gin.Context) {
	user, _ := c.Get("user")
	userObj := user.(*models.User)

	assignmentID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid assignment ID"})
		return
	}

	if err := h.assignmentService.DeleteAssignment(uint(assignmentID), userObj.ID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Assignment deleted successfully"})
}

// ============================================================================
// STUDENT MANAGEMENT
// ============================================================================

// AssignStudentsRequest represents the request to assign students
type AssignStudentsRequest struct {
	StudentIDs []uint `json:"student_ids" binding:"required"`
}

// AssignStudents handles POST /instructor/assignments/:id/assign
func (h *InstructorHandlers) AssignStudents(c *gin.Context) {
	user, _ := c.Get("user")
	userObj := user.(*models.User)

	assignmentID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid assignment ID"})
		return
	}

	var req AssignStudentsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.assignmentService.AssignToMultipleStudents(uint(assignmentID), req.StudentIDs, userObj.ID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Assignment assigned to students successfully",
		"count":   len(req.StudentIDs),
	})
}

// GetAssignmentStudents handles GET /instructor/assignments/:id/students
func (h *InstructorHandlers) GetAssignmentStudents(c *gin.Context) {
	user, _ := c.Get("user")
	userObj := user.(*models.User)

	assignmentID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid assignment ID"})
		return
	}

	students, err := h.assignmentService.GetAssignedStudents(uint(assignmentID), userObj.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"students": students,
		"total":    len(students),
	})
}

// RemoveStudent handles POST /instructor/assignments/:id/students/:student_id/remove
func (h *InstructorHandlers) RemoveStudent(c *gin.Context) {
	user, _ := c.Get("user")
	userObj := user.(*models.User)

	assignmentID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid assignment ID"})
		return
	}

	studentID, err := strconv.ParseUint(c.Param("student_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid student ID"})
		return
	}

	if err := h.assignmentService.RemoveStudentFromAssignment(uint(assignmentID), uint(studentID), userObj.ID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Student removed from assignment"})
}

// ListStudents handles GET /instructor/students
func (h *InstructorHandlers) ListStudents(c *gin.Context) {
	students, err := models.GetAllStudents(h.assignmentService.GetDB())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"students": students,
		"total":    len(students),
	})
}

// AssignToStudent handles POST /instructor/students/:username/assignments/:assignment_id/assign
func (h *InstructorHandlers) AssignToStudent(c *gin.Context) {
	user, _ := c.Get("user")
	userObj := user.(*models.User)

	username := c.Param("username")
	student, err := models.GetUserByUsername(h.assignmentService.GetDB(), username)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Student not found"})
		return
	}

	assignmentID, err := strconv.ParseUint(c.Param("assignment_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid assignment ID"})
		return
	}

	if err := h.assignmentService.AssignToStudent(uint(assignmentID), student.ID, userObj.ID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Assignment assigned to student successfully"})
}

// RemoveFromStudent handles DELETE /instructor/students/:username/assignments/:assignment_id/remove
func (h *InstructorHandlers) RemoveFromStudent(c *gin.Context) {
	user, _ := c.Get("user")
	userObj := user.(*models.User)

	username := c.Param("username")
	student, err := models.GetUserByUsername(h.assignmentService.GetDB(), username)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Student not found"})
		return
	}

	assignmentID, err := strconv.ParseUint(c.Param("assignment_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid assignment ID"})
		return
	}

	if err := h.assignmentService.RemoveStudentFromAssignment(uint(assignmentID), student.ID, userObj.ID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Student removed from assignment"})
}

// ============================================================================
// PROGRESS & ANALYTICS (Absorbed from ProgressTrackingHandlers)
// ============================================================================

// GetAssignmentProgress handles GET /instructor/assignments/:id/progress
func (h *InstructorHandlers) GetAssignmentProgress(c *gin.Context) {
	user, _ := c.Get("user")
	userObj := user.(*models.User)

	assignmentID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid assignment ID"})
		return
	}

	progress, err := h.assignmentService.GetAssignmentProgress(uint(assignmentID), userObj.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, progress)
}

// GetDetailedProgress handles GET /instructor/assignments/:id/detailed-progress
func (h *InstructorHandlers) GetDetailedProgress(c *gin.Context) {
	user, _ := c.Get("user")
	userObj := user.(*models.User)

	assignmentID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid assignment ID"})
		return
	}

	report, err := h.assignmentService.GetDetailedProgressReport(uint(assignmentID), userObj.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, report)
}

// GetProgressSummary handles GET /instructor/progress/summary
func (h *InstructorHandlers) GetProgressSummary(c *gin.Context) {
	user, _ := c.Get("user")
	userObj := user.(*models.User)

	summary, err := h.assignmentService.GetInstructorProgressSummary(userObj.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, summary)
}

// GetProgressTrends handles GET /instructor/progress/trends
func (h *InstructorHandlers) GetProgressTrends(c *gin.Context) {
	user, _ := c.Get("user")
	userObj := user.(*models.User)

	days := 30
	if daysStr := c.Query("days"); daysStr != "" {
		if d, err := strconv.Atoi(daysStr); err == nil {
			days = d
		}
	}

	trends, err := h.assignmentService.GetProgressTrends(userObj.ID, days)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, trends)
}

// GetCompletionAnalytics handles GET /instructor/progress/completion-analytics
func (h *InstructorHandlers) GetCompletionAnalytics(c *gin.Context) {
	user, _ := c.Get("user")
	userObj := user.(*models.User)

	analytics, err := h.assignmentService.GetCompletionAnalytics(userObj.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, analytics)
}

// ============================================================================
// DUE DATES & NOTIFICATIONS (Absorbed from DueDateNotificationHandlers)
// ============================================================================

// GetDueDateOverview handles GET /instructor/due-dates/overview
func (h *InstructorHandlers) GetDueDateOverview(c *gin.Context) {
	user, _ := c.Get("user")
	userObj := user.(*models.User)

	overview, err := h.assignmentService.GetInstructorDueDateOverview(userObj.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, overview)
}

// GetNotifications handles GET /instructor/due-dates/notifications
func (h *InstructorHandlers) GetNotifications(c *gin.Context) {
	user, _ := c.Get("user")
	userObj := user.(*models.User)

	notifications, err := h.assignmentService.GetInstructorDueDateNotifications(userObj.ID)
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

// GetDashboardStats handles GET /instructor/dashboard/stats
func (h *InstructorHandlers) GetDashboardStats(c *gin.Context) {
	user, _ := c.Get("user")
	userObj := user.(*models.User)

	if !userObj.IsInstructor() {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	assignments, err := h.assignmentService.GetAssignmentsByInstructor(userObj.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get assignments"})
		return
	}

	students, err := models.GetAllStudents(h.assignmentService.GetDB())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get students"})
		return
	}

	totalAssigned := 0
	for _, assignment := range assignments {
		count, _ := h.assignmentService.CountAssignedStudents(assignment.ID)
		totalAssigned += count
	}

	c.JSON(http.StatusOK, gin.H{
		"total_assignments": len(assignments),
		"active_students":   len(students),
		"total_assigned":    totalAssigned,
		"assignments":       assignments,
		"students":          students,
	})
}
