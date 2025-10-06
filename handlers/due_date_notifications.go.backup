package handlers

import (
	"net/http"
	"strconv"
	"zipcodereader/models"
	"zipcodereader/services"

	"github.com/gin-gonic/gin"
)

// DueDateNotificationHandlers handles due date notification operations
type DueDateNotificationHandlers struct {
	dueDateService *services.DueDateNotificationService
}

// NewDueDateNotificationHandlers creates new due date notification handlers
func NewDueDateNotificationHandlers(dueDateService *services.DueDateNotificationService) *DueDateNotificationHandlers {
	return &DueDateNotificationHandlers{
		dueDateService: dueDateService,
	}
}

// GetStudentDueDateAlerts handles GET /student/due-dates/alerts
func (h *DueDateNotificationHandlers) GetStudentDueDateAlerts(c *gin.Context) {
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

	// Get query parameters
	daysAhead := 7
	if days := c.Query("days"); days != "" {
		if d, err := strconv.Atoi(days); err == nil && d > 0 {
			daysAhead = d
		}
	}

	// Get upcoming alerts
	upcomingAlerts, err := h.dueDateService.GetUpcomingDueDateAlerts(userObj.ID, daysAhead)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Get overdue alerts
	overdueAlerts, err := h.dueDateService.GetOverdueDueDateAlerts(userObj.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"upcoming_alerts": upcomingAlerts,
		"overdue_alerts":  overdueAlerts,
		"total_upcoming":  len(upcomingAlerts),
		"total_overdue":   len(overdueAlerts),
	})
}

// GetStudentDueDateSummary handles GET /student/due-dates/summary
func (h *DueDateNotificationHandlers) GetStudentDueDateSummary(c *gin.Context) {
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

	// Get due date summary
	summary, err := h.dueDateService.GetDueDateSummary(userObj.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"summary": summary,
	})
}

// GetInstructorDueDateOverview handles GET /instructor/due-dates/overview
func (h *DueDateNotificationHandlers) GetInstructorDueDateOverview(c *gin.Context) {
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

	// Get due date overview
	overview, err := h.dueDateService.GetInstructorDueDateOverview(userObj.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"overview": overview,
	})
}

// GetDueDateNotifications handles GET /student/due-dates/notifications and GET /instructor/due-dates/notifications
func (h *DueDateNotificationHandlers) GetDueDateNotifications(c *gin.Context) {
	// Get user from context
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	userObj := user.(*models.User)
	if !userObj.IsStudent() && !userObj.IsInstructor() {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	// Get due date summary
	summary, err := h.dueDateService.GetDueDateSummary(userObj.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Generate notification messages
	var notifications []gin.H

	// Add overdue notifications
	for _, alert := range summary.OverdueAlerts {
		message := h.dueDateService.GenerateDueDateNotificationMessage(alert)
		notifications = append(notifications, gin.H{
			"type":     "overdue",
			"priority": alert.Priority,
			"message":  message,
			"alert":    alert,
		})
	}

	// Add due today notifications
	for _, alert := range summary.DueTodayAlerts {
		message := h.dueDateService.GenerateDueDateNotificationMessage(alert)
		notifications = append(notifications, gin.H{
			"type":     "due_today",
			"priority": alert.Priority,
			"message":  message,
			"alert":    alert,
		})
	}

	// Add upcoming notifications (high priority only)
	for _, alert := range summary.UpcomingAlerts {
		if alert.Priority == "high" || alert.Priority == "medium" {
			message := h.dueDateService.GenerateDueDateNotificationMessage(alert)
			notifications = append(notifications, gin.H{
				"type":     "upcoming",
				"priority": alert.Priority,
				"message":  message,
				"alert":    alert,
			})
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"notifications": notifications,
		"count":         len(notifications),
		"summary":       summary,
	})
}
