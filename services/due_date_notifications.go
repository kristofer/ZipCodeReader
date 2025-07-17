package services

import (
	"fmt"
	"time"
	"zipcodereader/models"

	"gorm.io/gorm"
)

// DueDateNotificationService handles due date notifications and alerts
type DueDateNotificationService struct {
	db *gorm.DB
}

// NewDueDateNotificationService creates a new due date notification service
func NewDueDateNotificationService(db *gorm.DB) *DueDateNotificationService {
	return &DueDateNotificationService{db: db}
}

// DueDateAlert represents a due date alert
type DueDateAlert struct {
	StudentID       uint      `json:"student_id"`
	StudentName     string    `json:"student_name"`
	StudentEmail    string    `json:"student_email"`
	AssignmentID    uint      `json:"assignment_id"`
	AssignmentTitle string    `json:"assignment_title"`
	AssignmentURL   string    `json:"assignment_url"`
	DueDate         time.Time `json:"due_date"`
	DaysUntilDue    int       `json:"days_until_due"`
	Status          string    `json:"status"`
	AlertType       string    `json:"alert_type"` // "upcoming", "overdue", "due_today"
	Priority        string    `json:"priority"`   // "low", "medium", "high", "critical"
}

// DueDateSummary provides summary of due date information
type DueDateSummary struct {
	TotalUpcoming  int            `json:"total_upcoming"`
	DueToday       int            `json:"due_today"`
	DueTomorrow    int            `json:"due_tomorrow"`
	DueThisWeek    int            `json:"due_this_week"`
	Overdue        int            `json:"overdue"`
	UpcomingAlerts []DueDateAlert `json:"upcoming_alerts"`
	OverdueAlerts  []DueDateAlert `json:"overdue_alerts"`
	DueTodayAlerts []DueDateAlert `json:"due_today_alerts"`
}

// GetUpcomingDueDateAlerts retrieves upcoming due date alerts for a student
func (s *DueDateNotificationService) GetUpcomingDueDateAlerts(studentID uint, daysAhead int) ([]DueDateAlert, error) {
	if daysAhead <= 0 {
		daysAhead = 7 // Default to 7 days ahead
	}

	var alerts []DueDateAlert
	cutoffDate := time.Now().AddDate(0, 0, daysAhead)

	type AlertResult struct {
		StudentID       uint      `json:"student_id"`
		StudentName     string    `json:"student_name"`
		StudentEmail    string    `json:"student_email"`
		AssignmentID    uint      `json:"assignment_id"`
		AssignmentTitle string    `json:"assignment_title"`
		AssignmentURL   string    `json:"assignment_url"`
		DueDate         time.Time `json:"due_date"`
		Status          string    `json:"status"`
	}

	var results []AlertResult

	err := s.db.Table("student_assignments").
		Select("student_assignments.student_id, users.username as student_name, users.email as student_email, "+
			"assignments.id as assignment_id, assignments.title as assignment_title, assignments.url as assignment_url, "+
			"assignments.due_date, student_assignments.status").
		Joins("JOIN users ON users.id = student_assignments.student_id").
		Joins("JOIN assignments ON assignments.id = student_assignments.assignment_id").
		Where("student_assignments.student_id = ? AND assignments.due_date IS NOT NULL AND assignments.due_date >= ? AND assignments.due_date <= ? AND student_assignments.status != ?",
			studentID, time.Now(), cutoffDate, models.StatusCompleted).
		Order("assignments.due_date ASC").
		Find(&results).Error

	if err != nil {
		return nil, err
	}

	for _, result := range results {
		daysUntil := int(result.DueDate.Sub(time.Now()).Hours() / 24)

		alertType := "upcoming"
		priority := "low"

		if daysUntil == 0 {
			alertType = "due_today"
			priority = "high"
		} else if daysUntil == 1 {
			alertType = "due_tomorrow"
			priority = "medium"
		} else if daysUntil <= 3 {
			priority = "medium"
		}

		alerts = append(alerts, DueDateAlert{
			StudentID:       result.StudentID,
			StudentName:     result.StudentName,
			StudentEmail:    result.StudentEmail,
			AssignmentID:    result.AssignmentID,
			AssignmentTitle: result.AssignmentTitle,
			AssignmentURL:   result.AssignmentURL,
			DueDate:         result.DueDate,
			DaysUntilDue:    daysUntil,
			Status:          result.Status,
			AlertType:       alertType,
			Priority:        priority,
		})
	}

	return alerts, nil
}

// GetOverdueDueDateAlerts retrieves overdue assignments for a student
func (s *DueDateNotificationService) GetOverdueDueDateAlerts(studentID uint) ([]DueDateAlert, error) {
	var alerts []DueDateAlert

	type AlertResult struct {
		StudentID       uint      `json:"student_id"`
		StudentName     string    `json:"student_name"`
		StudentEmail    string    `json:"student_email"`
		AssignmentID    uint      `json:"assignment_id"`
		AssignmentTitle string    `json:"assignment_title"`
		AssignmentURL   string    `json:"assignment_url"`
		DueDate         time.Time `json:"due_date"`
		Status          string    `json:"status"`
	}

	var results []AlertResult

	err := s.db.Table("student_assignments").
		Select("student_assignments.student_id, users.username as student_name, users.email as student_email, "+
			"assignments.id as assignment_id, assignments.title as assignment_title, assignments.url as assignment_url, "+
			"assignments.due_date, student_assignments.status").
		Joins("JOIN users ON users.id = student_assignments.student_id").
		Joins("JOIN assignments ON assignments.id = student_assignments.assignment_id").
		Where("student_assignments.student_id = ? AND assignments.due_date IS NOT NULL AND assignments.due_date < ? AND student_assignments.status != ?",
			studentID, time.Now(), models.StatusCompleted).
		Order("assignments.due_date ASC").
		Find(&results).Error

	if err != nil {
		return nil, err
	}

	for _, result := range results {
		daysPastDue := int(time.Now().Sub(result.DueDate).Hours() / 24)

		priority := "high"
		if daysPastDue > 7 {
			priority = "critical"
		}

		alerts = append(alerts, DueDateAlert{
			StudentID:       result.StudentID,
			StudentName:     result.StudentName,
			StudentEmail:    result.StudentEmail,
			AssignmentID:    result.AssignmentID,
			AssignmentTitle: result.AssignmentTitle,
			AssignmentURL:   result.AssignmentURL,
			DueDate:         result.DueDate,
			DaysUntilDue:    -daysPastDue, // Negative for overdue
			Status:          result.Status,
			AlertType:       "overdue",
			Priority:        priority,
		})
	}

	return alerts, nil
}

// GetDueDateSummary provides a comprehensive summary of due date information for a student
func (s *DueDateNotificationService) GetDueDateSummary(studentID uint) (*DueDateSummary, error) {
	summary := &DueDateSummary{}

	// Get upcoming alerts
	upcomingAlerts, err := s.GetUpcomingDueDateAlerts(studentID, 7)
	if err != nil {
		return nil, err
	}

	// Get overdue alerts
	overdueAlerts, err := s.GetOverdueDueDateAlerts(studentID)
	if err != nil {
		return nil, err
	}

	// Process alerts for summary
	var dueTodayAlerts []DueDateAlert
	dueTomorrow := 0
	dueThisWeek := 0

	for _, alert := range upcomingAlerts {
		if alert.DaysUntilDue == 0 {
			dueTodayAlerts = append(dueTodayAlerts, alert)
		} else if alert.DaysUntilDue == 1 {
			dueTomorrow++
		}

		if alert.DaysUntilDue <= 7 {
			dueThisWeek++
		}
	}

	summary.TotalUpcoming = len(upcomingAlerts)
	summary.DueToday = len(dueTodayAlerts)
	summary.DueTomorrow = dueTomorrow
	summary.DueThisWeek = dueThisWeek
	summary.Overdue = len(overdueAlerts)
	summary.UpcomingAlerts = upcomingAlerts
	summary.OverdueAlerts = overdueAlerts
	summary.DueTodayAlerts = dueTodayAlerts

	return summary, nil
}

// GetInstructorDueDateOverview provides due date overview for all instructor's assignments
func (s *DueDateNotificationService) GetInstructorDueDateOverview(instructorID uint) (map[string]interface{}, error) {
	overview := make(map[string]interface{})

	// Get all assignments by instructor
	assignments, err := models.GetAssignmentsByInstructor(s.db, instructorID)
	if err != nil {
		return nil, err
	}

	totalAssignments := len(assignments)
	assignmentsWithDueDates := 0
	upcomingDueDates := 0
	overdueAssignments := 0

	var upcomingDeadlines []map[string]interface{}
	var overdueList []map[string]interface{}

	for _, assignment := range assignments {
		if assignment.DueDate != nil {
			assignmentsWithDueDates++

			// Check if upcoming (within 7 days)
			if assignment.DueDate.After(time.Now()) && assignment.DueDate.Before(time.Now().AddDate(0, 0, 7)) {
				upcomingDueDates++

				// Get student count for this assignment
				studentAssignments, _ := models.GetStudentAssignmentsByAssignment(s.db, assignment.ID)
				incompleteCount := 0
				for _, sa := range studentAssignments {
					if sa.Status != models.StatusCompleted {
						incompleteCount++
					}
				}

				upcomingDeadlines = append(upcomingDeadlines, map[string]interface{}{
					"assignment_id":    assignment.ID,
					"title":            assignment.Title,
					"due_date":         assignment.DueDate,
					"days_until_due":   int(assignment.DueDate.Sub(time.Now()).Hours() / 24),
					"incomplete_count": incompleteCount,
					"total_students":   len(studentAssignments),
				})
			}

			// Check if overdue
			if assignment.DueDate.Before(time.Now()) {
				// Get student count for this assignment
				studentAssignments, _ := models.GetStudentAssignmentsByAssignment(s.db, assignment.ID)
				incompleteCount := 0
				for _, sa := range studentAssignments {
					if sa.Status != models.StatusCompleted {
						incompleteCount++
					}
				}

				if incompleteCount > 0 {
					overdueAssignments++
					overdueList = append(overdueList, map[string]interface{}{
						"assignment_id":    assignment.ID,
						"title":            assignment.Title,
						"due_date":         assignment.DueDate,
						"days_overdue":     int(time.Now().Sub(*assignment.DueDate).Hours() / 24),
						"incomplete_count": incompleteCount,
						"total_students":   len(studentAssignments),
					})
				}
			}
		}
	}

	overview["total_assignments"] = totalAssignments
	overview["assignments_with_due_dates"] = assignmentsWithDueDates
	overview["upcoming_due_dates"] = upcomingDueDates
	overview["overdue_assignments"] = overdueAssignments
	overview["upcoming_deadlines"] = upcomingDeadlines
	overview["overdue_list"] = overdueList

	return overview, nil
}

// GenerateDueDateNotificationMessage generates a notification message for due date alerts
func (s *DueDateNotificationService) GenerateDueDateNotificationMessage(alert DueDateAlert) string {
	switch alert.AlertType {
	case "due_today":
		return fmt.Sprintf("📅 Assignment '%s' is due today! Complete it at: %s",
			alert.AssignmentTitle, alert.AssignmentURL)
	case "due_tomorrow":
		return fmt.Sprintf("⏰ Assignment '%s' is due tomorrow (%s). Complete it at: %s",
			alert.AssignmentTitle, alert.DueDate.Format("Jan 2"), alert.AssignmentURL)
	case "upcoming":
		return fmt.Sprintf("📚 Assignment '%s' is due in %d days (%s). Complete it at: %s",
			alert.AssignmentTitle, alert.DaysUntilDue, alert.DueDate.Format("Jan 2"), alert.AssignmentURL)
	case "overdue":
		return fmt.Sprintf("🚨 Assignment '%s' was due %d days ago (%s). Complete it now at: %s",
			alert.AssignmentTitle, -alert.DaysUntilDue, alert.DueDate.Format("Jan 2"), alert.AssignmentURL)
	default:
		return fmt.Sprintf("Assignment '%s' has an upcoming due date: %s",
			alert.AssignmentTitle, alert.DueDate.Format("Jan 2, 2006"))
	}
}
