package handlers

import (
	"net/http"
	"strconv"
	"zipcodereader/models"
	"zipcodereader/services"

	"github.com/gin-gonic/gin"
)

// DashboardHandlers handles dashboard-related operations
type DashboardHandlers struct {
	assignmentService        *services.AssignmentService
	studentAssignmentService *services.StudentAssignmentService
	useLocalAuth             bool
}

// NewDashboardHandlers creates new dashboard handlers
func NewDashboardHandlers(assignmentService *services.AssignmentService, studentAssignmentService *services.StudentAssignmentService, useLocalAuth bool) *DashboardHandlers {
	return &DashboardHandlers{
		assignmentService:        assignmentService,
		studentAssignmentService: studentAssignmentService,
		useLocalAuth:             useLocalAuth,
	}
}

// ShowInstructorDashboard renders the instructor assignment dashboard
func (h *DashboardHandlers) ShowInstructorDashboard(c *gin.Context) {
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

	c.HTML(http.StatusOK, "instructor_assignments.html", gin.H{
		"title":          "Assignment Management",
		"user":           userObj,
		"use_local_auth": h.useLocalAuth,
		"template_type":  "instructor",
	})
}

// ShowStudentDashboard renders the student assignment dashboard
func (h *DashboardHandlers) ShowStudentDashboard(c *gin.Context) {
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

	c.HTML(http.StatusOK, "student_assignments.html", gin.H{
		"title":          "My Assignments",
		"user":           userObj,
		"use_local_auth": h.useLocalAuth,
		"template_type":  "student",
	})
}

// ShowAssignmentDetail renders the assignment detail page
func (h *DashboardHandlers) ShowAssignmentDetail(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	userObj := user.(*models.User)

	// Get assignment ID from URL
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid assignment ID"})
		return
	}

	// Check if user can access this assignment
	if userObj.IsInstructor() {
		// Instructor can access any assignment they created
		assignment, err := h.assignmentService.GetAssignmentByID(uint(id), userObj.ID)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Assignment not found"})
			return
		}

		c.HTML(http.StatusOK, "assignment_detail.html", gin.H{
			"title":          assignment.Title,
			"user":           userObj,
			"assignment":     assignment,
			"use_local_auth": h.useLocalAuth,
		})
	} else {
		// Student can only access assigned assignments
		studentAssignment, err := h.studentAssignmentService.GetStudentAssignmentByID(uint(id), userObj.ID)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Assignment not found or not assigned to you"})
			return
		}

		c.HTML(http.StatusOK, "assignment_detail.html", gin.H{
			"title":             studentAssignment.Assignment.Title,
			"user":              userObj,
			"assignment":        studentAssignment.Assignment,
			"studentAssignment": studentAssignment,
			"use_local_auth":    h.useLocalAuth,
			"template_type":     "student",
		})
	}
}

// ShowAssignmentProgress renders the assignment progress page for instructors
func (h *DashboardHandlers) ShowAssignmentProgress(c *gin.Context) {
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

	// Get assignment details
	assignment, err := h.assignmentService.GetAssignmentByID(uint(id), userObj.ID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Assignment not found"})
		return
	}

	// Get progress data
	progress, err := h.assignmentService.GetAssignmentProgress(uint(id), userObj.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error loading progress data"})
		return
	}

	c.HTML(http.StatusOK, "assignment_progress.html", gin.H{
		"title":          "Assignment Progress - " + assignment.Title,
		"user":           userObj,
		"assignment":     assignment,
		"progress":       progress,
		"use_local_auth": h.useLocalAuth,
	})
}

// ShowAssignmentManagement renders the assignment management page with full CRUD functionality
func (h *DashboardHandlers) ShowAssignmentManagement(c *gin.Context) {
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

	c.HTML(http.StatusOK, "assignment_management.html", gin.H{
		"title":          "Assignment Management",
		"user":           userObj,
		"use_local_auth": h.useLocalAuth,
		"template_type":  "instructor",
	})
}
