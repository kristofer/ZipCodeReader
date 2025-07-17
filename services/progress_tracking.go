package services

import (
	"errors"
	"time"
	"zipcodereader/models"

	"gorm.io/gorm"
)

// ProgressTrackingService handles advanced progress tracking functionality
type ProgressTrackingService struct {
	db *gorm.DB
}

// NewProgressTrackingService creates a new progress tracking service
func NewProgressTrackingService(db *gorm.DB) *ProgressTrackingService {
	return &ProgressTrackingService{db: db}
}

// DetailedProgressReport contains comprehensive progress information
type DetailedProgressReport struct {
	AssignmentID          uint                    `json:"assignment_id"`
	Title                 string                  `json:"title"`
	TotalStudents         int                     `json:"total_students"`
	CompletionRate        float64                 `json:"completion_rate"`
	AverageTimeToComplete int                     `json:"average_time_to_complete_hours"`
	StatusBreakdown       map[string]int          `json:"status_breakdown"`
	OverdueCount          int                     `json:"overdue_count"`
	StudentDetails        []StudentProgressDetail `json:"student_details"`
	CreatedAt             time.Time               `json:"created_at"`
	DueDate               *time.Time              `json:"due_date"`
}

// StudentProgressDetail contains individual student progress information
type StudentProgressDetail struct {
	StudentID      uint       `json:"student_id"`
	StudentName    string     `json:"student_name"`
	StudentEmail   string     `json:"student_email"`
	Status         string     `json:"status"`
	AssignedAt     time.Time  `json:"assigned_at"`
	CompletedAt    *time.Time `json:"completed_at"`
	TimeToComplete *int       `json:"time_to_complete_hours"`
	IsOverdue      bool       `json:"is_overdue"`
}

// InstructorProgressSummary contains overall instructor progress statistics
type InstructorProgressSummary struct {
	TotalAssignments        int                        `json:"total_assignments"`
	TotalStudentAssignments int                        `json:"total_student_assignments"`
	OverallCompletionRate   float64                    `json:"overall_completion_rate"`
	AssignmentsWithDueDates int                        `json:"assignments_with_due_dates"`
	OverdueAssignments      int                        `json:"overdue_assignments"`
	AverageCompletionTime   int                        `json:"average_completion_time_hours"`
	CategoryBreakdown       map[string]CategoryStats   `json:"category_breakdown"`
	RecentCompletions       []RecentCompletionActivity `json:"recent_completions"`
	StudentEngagement       map[string]interface{}     `json:"student_engagement"`
}

// CategoryStats contains statistics for a specific category
type CategoryStats struct {
	AssignmentCount       int     `json:"assignment_count"`
	CompletionRate        float64 `json:"completion_rate"`
	AverageTimeToComplete int     `json:"average_time_to_complete_hours"`
}

// RecentCompletionActivity represents recent completion activity
type RecentCompletionActivity struct {
	StudentName     string    `json:"student_name"`
	AssignmentTitle string    `json:"assignment_title"`
	CompletedAt     time.Time `json:"completed_at"`
	TimeTaken       int       `json:"time_taken_hours"`
}

// GetDetailedProgressReport generates a comprehensive progress report for an assignment
func (s *ProgressTrackingService) GetDetailedProgressReport(assignmentID uint, instructorID uint) (*DetailedProgressReport, error) {
	// Validate assignment exists and instructor owns it
	assignment, err := models.GetAssignmentByID(s.db, assignmentID)
	if err != nil {
		return nil, errors.New("assignment not found")
	}

	if assignment.CreatedByID != instructorID {
		return nil, errors.New("access denied")
	}

	// Get all student assignments for this assignment
	studentAssignments, err := models.GetStudentAssignmentsByAssignment(s.db, assignmentID)
	if err != nil {
		return nil, err
	}

	// Calculate basic statistics
	totalStudents := len(studentAssignments)
	statusBreakdown := make(map[string]int)
	var completedCount int
	var totalCompletionTime int
	var overdueCount int
	var studentDetails []StudentProgressDetail

	// Initialize status breakdown
	statusBreakdown[models.StatusAssigned] = 0
	statusBreakdown[models.StatusInProgress] = 0
	statusBreakdown[models.StatusCompleted] = 0

	for _, sa := range studentAssignments {
		// Update status breakdown
		statusBreakdown[sa.Status]++

		// Check if overdue
		isOverdue := false
		if assignment.DueDate != nil && sa.Status != models.StatusCompleted {
			isOverdue = time.Now().After(*assignment.DueDate)
			if isOverdue {
				overdueCount++
			}
		}

		// Calculate time to complete
		var timeToComplete *int
		if sa.CompletedAt != nil {
			hours := int(sa.CompletedAt.Sub(sa.CreatedAt).Hours())
			timeToComplete = &hours
			totalCompletionTime += hours
		}

		if sa.Status == models.StatusCompleted {
			completedCount++
		}

		// Add student detail
		studentDetails = append(studentDetails, StudentProgressDetail{
			StudentID:      sa.StudentID,
			StudentName:    sa.Student.Username,
			StudentEmail:   sa.Student.Email,
			Status:         sa.Status,
			AssignedAt:     sa.CreatedAt,
			CompletedAt:    sa.CompletedAt,
			TimeToComplete: timeToComplete,
			IsOverdue:      isOverdue,
		})
	}

	// Calculate completion rate
	completionRate := 0.0
	if totalStudents > 0 {
		completionRate = float64(completedCount) / float64(totalStudents) * 100
	}

	// Calculate average time to complete
	averageTimeToComplete := 0
	if completedCount > 0 {
		averageTimeToComplete = totalCompletionTime / completedCount
	}

	return &DetailedProgressReport{
		AssignmentID:          assignmentID,
		Title:                 assignment.Title,
		TotalStudents:         totalStudents,
		CompletionRate:        completionRate,
		AverageTimeToComplete: averageTimeToComplete,
		StatusBreakdown:       statusBreakdown,
		OverdueCount:          overdueCount,
		StudentDetails:        studentDetails,
		CreatedAt:             assignment.CreatedAt,
		DueDate:               assignment.DueDate,
	}, nil
}

// GetInstructorProgressSummary generates comprehensive instructor progress summary
func (s *ProgressTrackingService) GetInstructorProgressSummary(instructorID uint) (*InstructorProgressSummary, error) {
	// Get all assignments by instructor
	assignments, err := models.GetAssignmentsByInstructor(s.db, instructorID)
	if err != nil {
		return nil, err
	}

	totalAssignments := len(assignments)
	assignmentsWithDueDates := 0
	categoryBreakdown := make(map[string]CategoryStats)

	var totalStudentAssignments int
	var totalCompleted int
	var totalCompletionTime int
	var completedAssignments int
	var overdueAssignments int

	// Process each assignment
	for _, assignment := range assignments {
		// Count assignments with due dates
		if assignment.DueDate != nil {
			assignmentsWithDueDates++
		}

		// Get student assignments for this assignment
		studentAssignments, err := models.GetStudentAssignmentsByAssignment(s.db, assignment.ID)
		if err != nil {
			continue
		}

		assignmentCompleted := 0
		assignmentTotalTime := 0
		assignmentOverdue := 0

		for _, sa := range studentAssignments {
			totalStudentAssignments++

			if sa.Status == models.StatusCompleted {
				totalCompleted++
				assignmentCompleted++
				completedAssignments++

				if sa.CompletedAt != nil {
					hours := int(sa.CompletedAt.Sub(sa.CreatedAt).Hours())
					totalCompletionTime += hours
					assignmentTotalTime += hours
				}
			}

			// Check if overdue
			if assignment.DueDate != nil && sa.Status != models.StatusCompleted {
				if time.Now().After(*assignment.DueDate) {
					overdueAssignments++
					assignmentOverdue++
				}
			}
		}

		// Update category breakdown
		category := assignment.Category
		if category == "" {
			category = "uncategorized"
		}

		if stats, exists := categoryBreakdown[category]; exists {
			stats.AssignmentCount++
			// Update completion rate and average time
			if len(studentAssignments) > 0 {
				stats.CompletionRate = (stats.CompletionRate*(float64(stats.AssignmentCount-1)) +
					float64(assignmentCompleted)/float64(len(studentAssignments))*100) / float64(stats.AssignmentCount)
			}
			if assignmentCompleted > 0 {
				newAvgTime := assignmentTotalTime / assignmentCompleted
				stats.AverageTimeToComplete = (stats.AverageTimeToComplete*(stats.AssignmentCount-1) + newAvgTime) / stats.AssignmentCount
			}
			categoryBreakdown[category] = stats
		} else {
			completionRate := 0.0
			if len(studentAssignments) > 0 {
				completionRate = float64(assignmentCompleted) / float64(len(studentAssignments)) * 100
			}
			avgTime := 0
			if assignmentCompleted > 0 {
				avgTime = assignmentTotalTime / assignmentCompleted
			}
			categoryBreakdown[category] = CategoryStats{
				AssignmentCount:       1,
				CompletionRate:        completionRate,
				AverageTimeToComplete: avgTime,
			}
		}
	}

	// Calculate overall completion rate
	overallCompletionRate := 0.0
	if totalStudentAssignments > 0 {
		overallCompletionRate = float64(totalCompleted) / float64(totalStudentAssignments) * 100
	}

	// Calculate average completion time
	averageCompletionTime := 0
	if completedAssignments > 0 {
		averageCompletionTime = totalCompletionTime / completedAssignments
	}

	// Get recent completions
	recentCompletions, err := s.getRecentCompletions(instructorID, 10)
	if err != nil {
		recentCompletions = []RecentCompletionActivity{}
	}

	// Calculate student engagement metrics
	studentEngagement := s.calculateStudentEngagement(instructorID)

	return &InstructorProgressSummary{
		TotalAssignments:        totalAssignments,
		TotalStudentAssignments: totalStudentAssignments,
		OverallCompletionRate:   overallCompletionRate,
		AssignmentsWithDueDates: assignmentsWithDueDates,
		OverdueAssignments:      overdueAssignments,
		AverageCompletionTime:   averageCompletionTime,
		CategoryBreakdown:       categoryBreakdown,
		RecentCompletions:       recentCompletions,
		StudentEngagement:       studentEngagement,
	}, nil
}

// getRecentCompletions retrieves recent completion activities
func (s *ProgressTrackingService) getRecentCompletions(instructorID uint, limit int) ([]RecentCompletionActivity, error) {
	var results []RecentCompletionActivity

	type CompletionResult struct {
		StudentName     string    `json:"student_name"`
		AssignmentTitle string    `json:"assignment_title"`
		CompletedAt     time.Time `json:"completed_at"`
		AssignedAt      time.Time `json:"assigned_at"`
	}

	var completionResults []CompletionResult

	err := s.db.Table("student_assignments").
		Select("users.username as student_name, assignments.title as assignment_title, student_assignments.completed_at, student_assignments.created_at as assigned_at").
		Joins("JOIN users ON users.id = student_assignments.student_id").
		Joins("JOIN assignments ON assignments.id = student_assignments.assignment_id").
		Where("assignments.created_by_id = ? AND student_assignments.completed_at IS NOT NULL", instructorID).
		Order("student_assignments.completed_at DESC").
		Limit(limit).
		Find(&completionResults).Error

	if err != nil {
		return nil, err
	}

	for _, result := range completionResults {
		timeTaken := int(result.CompletedAt.Sub(result.AssignedAt).Hours())
		results = append(results, RecentCompletionActivity{
			StudentName:     result.StudentName,
			AssignmentTitle: result.AssignmentTitle,
			CompletedAt:     result.CompletedAt,
			TimeTaken:       timeTaken,
		})
	}

	return results, nil
}

// calculateStudentEngagement calculates student engagement metrics
func (s *ProgressTrackingService) calculateStudentEngagement(instructorID uint) map[string]interface{} {
	engagement := make(map[string]interface{})

	// Count active students (students with at least one assignment)
	var activeStudents int64
	s.db.Table("student_assignments").
		Joins("JOIN assignments ON assignments.id = student_assignments.assignment_id").
		Where("assignments.created_by_id = ?", instructorID).
		Distinct("student_assignments.student_id").
		Count(&activeStudents)

	engagement["active_students"] = activeStudents

	// Calculate average assignments per student
	var totalAssignments int64
	s.db.Table("student_assignments").
		Joins("JOIN assignments ON assignments.id = student_assignments.assignment_id").
		Where("assignments.created_by_id = ?", instructorID).
		Count(&totalAssignments)

	avgAssignmentsPerStudent := 0.0
	if activeStudents > 0 {
		avgAssignmentsPerStudent = float64(totalAssignments) / float64(activeStudents)
	}
	engagement["average_assignments_per_student"] = avgAssignmentsPerStudent

	// Calculate completion rate by time period (last 7 days, last 30 days)
	sevenDaysAgo := time.Now().AddDate(0, 0, -7)
	thirtyDaysAgo := time.Now().AddDate(0, 0, -30)

	var completionsLast7Days int64
	var completionsLast30Days int64

	s.db.Table("student_assignments").
		Joins("JOIN assignments ON assignments.id = student_assignments.assignment_id").
		Where("assignments.created_by_id = ? AND student_assignments.completed_at >= ?", instructorID, sevenDaysAgo).
		Count(&completionsLast7Days)

	s.db.Table("student_assignments").
		Joins("JOIN assignments ON assignments.id = student_assignments.assignment_id").
		Where("assignments.created_by_id = ? AND student_assignments.completed_at >= ?", instructorID, thirtyDaysAgo).
		Count(&completionsLast30Days)

	engagement["completions_last_7_days"] = completionsLast7Days
	engagement["completions_last_30_days"] = completionsLast30Days

	return engagement
}
