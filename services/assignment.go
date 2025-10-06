package services

import (
	"errors"
	"fmt"
	"time"
	"zipcodereader/models"

	"gorm.io/gorm"
)

// AssignmentService handles comprehensive assignment management functionality
// Consolidates: assignment CRUD, student assignments, progress tracking, and due date notifications
type AssignmentService struct {
	db *gorm.DB
}

// NewAssignmentService creates a new unified assignment service
func NewAssignmentService(db *gorm.DB) *AssignmentService {
	return &AssignmentService{db: db}
}

// GetDB returns the database instance
func (s *AssignmentService) GetDB() *gorm.DB {
	return s.db
}

// ============================================================================
// INPUT STRUCTURES
// ============================================================================

// CreateAssignmentInput represents input for creating an assignment
type CreateAssignmentInput struct {
	Title       string
	Description string
	URL         string
	Category    string
	DueDate     *time.Time
}

// UpdateAssignmentInput represents input for updating an assignment
type UpdateAssignmentInput struct {
	Title       string
	Description string
	URL         string
	Category    string
	DueDate     *time.Time
}

// ============================================================================
// PROGRESS TRACKING STRUCTURES
// ============================================================================

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

// ============================================================================
// DUE DATE NOTIFICATION STRUCTURES
// ============================================================================

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

// ============================================================================
// ASSIGNMENT CRUD OPERATIONS
// ============================================================================

// CreateAssignment creates a new assignment with validation
func (s *AssignmentService) CreateAssignment(instructorID uint, input CreateAssignmentInput) (*models.Assignment, error) {
	// Validate instructor exists and has instructor role
	var instructor models.User
	if err := s.db.First(&instructor, instructorID).Error; err != nil {
		return nil, errors.New("instructor not found")
	}

	if !instructor.IsInstructor() {
		return nil, errors.New("user is not an instructor")
	}

	// Validate input
	if input.Title == "" {
		return nil, errors.New("title is required")
	}

	if input.URL == "" {
		return nil, errors.New("URL is required")
	}

	// Create assignment
	assignment, err := models.CreateAssignment(s.db, input.Title, input.Description, input.URL, input.Category, input.DueDate, instructorID)
	if err != nil {
		return nil, err
	}

	return assignment, nil
}

// GetAssignmentByID retrieves an assignment by ID with authorization check
func (s *AssignmentService) GetAssignmentByID(assignmentID uint, userID uint) (*models.Assignment, error) {
	assignment, err := models.GetAssignmentByID(s.db, assignmentID)
	if err != nil {
		return nil, err
	}

	// Check if user is the creator or has access to this assignment
	if assignment.CreatedByID != userID {
		// Check if user is a student assigned to this assignment
		var user models.User
		if err := s.db.First(&user, userID).Error; err != nil {
			return nil, errors.New("user not found")
		}

		if user.IsStudent() {
			// Check if student is assigned to this assignment
			_, err := models.GetStudentAssignment(s.db, assignmentID, userID)
			if err != nil {
				return nil, errors.New("assignment not found or access denied")
			}
		} else {
			return nil, errors.New("access denied")
		}
	}

	return assignment, nil
}

// GetAssignmentsByInstructor retrieves all assignments for an instructor
func (s *AssignmentService) GetAssignmentsByInstructor(instructorID uint) ([]models.Assignment, error) {
	// Validate instructor exists and has instructor role
	var instructor models.User
	if err := s.db.First(&instructor, instructorID).Error; err != nil {
		return nil, errors.New("instructor not found")
	}

	if !instructor.IsInstructor() {
		return nil, errors.New("user is not an instructor")
	}

	assignments, err := models.GetAssignmentsByInstructor(s.db, instructorID)
	if err != nil {
		return nil, err
	}

	return assignments, nil
}

// UpdateAssignment updates an existing assignment
func (s *AssignmentService) UpdateAssignment(assignmentID uint, instructorID uint, input UpdateAssignmentInput) (*models.Assignment, error) {
	// Get assignment and validate ownership
	assignment, err := models.GetAssignmentByID(s.db, assignmentID)
	if err != nil {
		return nil, err
	}

	if assignment.CreatedByID != instructorID {
		return nil, errors.New("access denied")
	}

	// Validate input
	if input.Title == "" {
		return nil, errors.New("title is required")
	}

	if input.URL == "" {
		return nil, errors.New("URL is required")
	}

	// Update assignment
	err = assignment.UpdateAssignment(s.db, input.Title, input.Description, input.URL, input.Category, input.DueDate)
	if err != nil {
		return nil, err
	}

	return assignment, nil
}

// DeleteAssignment deletes an assignment
func (s *AssignmentService) DeleteAssignment(assignmentID uint, instructorID uint) error {
	// Get assignment and validate ownership
	assignment, err := models.GetAssignmentByID(s.db, assignmentID)
	if err != nil {
		return err
	}

	if assignment.CreatedByID != instructorID {
		return errors.New("access denied")
	}

	// Delete assignment (soft delete)
	return assignment.DeleteAssignment(s.db)
}

// SearchAssignments searches assignments by query
func (s *AssignmentService) SearchAssignments(query string, instructorID uint) ([]models.Assignment, error) {
	// Validate instructor exists and has instructor role
	var instructor models.User
	if err := s.db.First(&instructor, instructorID).Error; err != nil {
		return nil, errors.New("instructor not found")
	}

	if !instructor.IsInstructor() {
		return nil, errors.New("user is not an instructor")
	}

	return models.SearchAssignments(s.db, query, instructorID)
}

// GetAssignmentsByCategory gets assignments by category
func (s *AssignmentService) GetAssignmentsByCategory(category string, instructorID uint) ([]models.Assignment, error) {
	// Validate instructor exists and has instructor role
	var instructor models.User
	if err := s.db.First(&instructor, instructorID).Error; err != nil {
		return nil, errors.New("instructor not found")
	}

	if !instructor.IsInstructor() {
		return nil, errors.New("user is not an instructor")
	}

	return models.GetAssignmentsByCategory(s.db, category, instructorID)
}

// GetAllStudents gets all students for assignment purposes
func (s *AssignmentService) GetAllStudents(instructorID uint) ([]models.User, error) {
	// Validate instructor exists and has instructor role
	var instructor models.User
	if err := s.db.First(&instructor, instructorID).Error; err != nil {
		return nil, errors.New("instructor not found")
	}

	if !instructor.IsInstructor() {
		return nil, errors.New("user is not an instructor")
	}

	var students []models.User
	err := s.db.Where("role = ?", "student").Find(&students).Error
	if err != nil {
		return nil, err
	}

	return students, nil
}

// ============================================================================
// ASSIGNMENT-STUDENT RELATIONSHIP MANAGEMENT
// ============================================================================

// AssignToStudent assigns an assignment to a student
func (s *AssignmentService) AssignToStudent(assignmentID uint, studentID uint, instructorID uint) error {
	// Validate assignment exists and instructor owns it
	assignment, err := models.GetAssignmentByID(s.db, assignmentID)
	if err != nil {
		return err
	}

	if assignment.CreatedByID != instructorID {
		return errors.New("access denied")
	}

	// Validate student exists and has student role
	var student models.User
	if err := s.db.First(&student, studentID).Error; err != nil {
		return errors.New("student not found")
	}

	if !student.IsStudent() {
		return errors.New("user is not a student")
	}

	// Check if already assigned
	_, err = models.GetStudentAssignment(s.db, assignmentID, studentID)
	if err == nil {
		return errors.New("assignment already assigned to this student")
	}

	// Create student assignment
	_, err = models.CreateStudentAssignment(s.db, assignmentID, studentID)
	return err
}

// AssignToMultipleStudents assigns an assignment to multiple students
func (s *AssignmentService) AssignToMultipleStudents(assignmentID uint, studentIDs []uint, instructorID uint) error {
	// Validate assignment exists and instructor owns it
	assignment, err := models.GetAssignmentByID(s.db, assignmentID)
	if err != nil {
		return err
	}

	if assignment.CreatedByID != instructorID {
		return errors.New("access denied")
	}

	// Validate all students exist and have student role
	var students []models.User
	if err := s.db.Where("id IN ? AND role = ?", studentIDs, "student").Find(&students).Error; err != nil {
		return err
	}

	if len(students) != len(studentIDs) {
		return errors.New("some students not found or not valid students")
	}

	// Filter out already assigned students
	var validStudentIDs []uint
	for _, studentID := range studentIDs {
		_, err := models.GetStudentAssignment(s.db, assignmentID, studentID)
		if err != nil { // Student not assigned yet
			validStudentIDs = append(validStudentIDs, studentID)
		}
	}

	if len(validStudentIDs) == 0 {
		return errors.New("all students are already assigned to this assignment")
	}

	// Bulk create student assignments
	return models.BulkCreateStudentAssignments(s.db, assignmentID, validStudentIDs)
}

// RemoveStudentAssignment removes a student assignment (legacy method name)
func (s *AssignmentService) RemoveStudentAssignment(assignmentID uint, studentID uint, instructorID uint) error {
	return s.RemoveStudentFromAssignment(assignmentID, studentID, instructorID)
}

// RemoveStudentFromAssignment removes a student from an assignment
func (s *AssignmentService) RemoveStudentFromAssignment(assignmentID uint, studentID uint, instructorID uint) error {
	// Validate assignment exists and instructor owns it
	assignment, err := models.GetAssignmentByID(s.db, assignmentID)
	if err != nil {
		return err
	}

	if assignment.CreatedByID != instructorID {
		return errors.New("access denied")
	}

	// Remove student assignment
	return models.RemoveStudentAssignment(s.db, assignmentID, studentID)
}

// GetAssignmentProgress gets progress statistics for an assignment
func (s *AssignmentService) GetAssignmentProgress(assignmentID uint, instructorID uint) (map[string]int, error) {
	// Validate assignment exists and instructor owns it
	assignment, err := models.GetAssignmentByID(s.db, assignmentID)
	if err != nil {
		return nil, err
	}

	if assignment.CreatedByID != instructorID {
		return nil, errors.New("access denied")
	}

	// Get progress statistics
	return models.GetAssignmentProgress(s.db, assignmentID)
}

// GetAssignmentStudents gets all students assigned to an assignment
func (s *AssignmentService) GetAssignmentStudents(assignmentID uint, instructorID uint) ([]models.StudentAssignment, error) {
	// Validate assignment exists and instructor owns it
	assignment, err := models.GetAssignmentByID(s.db, assignmentID)
	if err != nil {
		return nil, err
	}

	if assignment.CreatedByID != instructorID {
		return nil, errors.New("access denied")
	}

	// Get assigned students
	return models.GetStudentAssignmentsByAssignment(s.db, assignmentID)
}

// GetAssignedStudents gets list of students assigned to an assignment
func (s *AssignmentService) GetAssignedStudents(assignmentID uint, instructorID uint) ([]models.User, error) {
	// Validate assignment exists and instructor owns it
	assignment, err := models.GetAssignmentByID(s.db, assignmentID)
	if err != nil {
		return nil, err
	}

	if assignment.CreatedByID != instructorID {
		return nil, errors.New("access denied")
	}

	// Get student assignments with student data
	var studentAssignments []models.StudentAssignment
	err = s.db.Preload("Student").Where("assignment_id = ?", assignmentID).Find(&studentAssignments).Error
	if err != nil {
		return nil, err
	}

	// Extract student objects
	var students []models.User
	for _, sa := range studentAssignments {
		students = append(students, sa.Student)
	}

	return students, nil
}

// CountAssignedStudents counts the number of students assigned to an assignment
func (s *AssignmentService) CountAssignedStudents(assignmentID uint) (int, error) {
	var count int64
	err := s.db.Model(&models.StudentAssignment{}).Where("assignment_id = ?", assignmentID).Count(&count).Error
	return int(count), err
}

// ============================================================================
// STUDENT ASSIGNMENT OPERATIONS
// ============================================================================

// GetStudentAssignments retrieves all assignments for a student
func (s *AssignmentService) GetStudentAssignments(studentID uint) ([]models.StudentAssignment, error) {
	// Validate student exists and has student role
	var student models.User
	if err := s.db.First(&student, studentID).Error; err != nil {
		return nil, errors.New("student not found")
	}

	if !student.IsStudent() {
		return nil, errors.New("user is not a student")
	}

	return models.GetStudentAssignmentsByStudent(s.db, studentID)
}

// GetStudentAssignmentsByStatus retrieves student assignments by status
func (s *AssignmentService) GetStudentAssignmentsByStatus(studentID uint, status string) ([]models.StudentAssignment, error) {
	// Validate student exists and has student role
	var student models.User
	if err := s.db.First(&student, studentID).Error; err != nil {
		return nil, errors.New("student not found")
	}

	if !student.IsStudent() {
		return nil, errors.New("user is not a student")
	}

	// Validate status
	if status != models.StatusAssigned && status != models.StatusInProgress && status != models.StatusCompleted {
		return nil, errors.New("invalid status")
	}

	return models.GetStudentAssignmentsByStatus(s.db, studentID, status)
}

// GetAssignmentsByStatus retrieves student assignments by status (alias for consistency)
func (s *AssignmentService) GetAssignmentsByStatus(studentID uint, status string) ([]models.StudentAssignment, error) {
	return s.GetStudentAssignmentsByStatus(studentID, status)
}

// GetStudentAssignment retrieves a specific student assignment by assignment ID and student ID
func (s *AssignmentService) GetStudentAssignment(assignmentID uint, studentID uint) (*models.StudentAssignment, error) {
	// Validate student exists and has student role
	var student models.User
	if err := s.db.First(&student, studentID).Error; err != nil {
		return nil, errors.New("student not found")
	}

	if !student.IsStudent() {
		return nil, errors.New("user is not a student")
	}

	return models.GetStudentAssignment(s.db, assignmentID, studentID)
}

// GetStudentAssignmentByID retrieves a specific student assignment by its ID
func (s *AssignmentService) GetStudentAssignmentByID(studentAssignmentID uint, studentID uint) (*models.StudentAssignment, error) {
	// Validate student exists and has student role
	var student models.User
	if err := s.db.First(&student, studentID).Error; err != nil {
		return nil, errors.New("student not found")
	}

	if !student.IsStudent() {
		return nil, errors.New("user is not a student")
	}

	return models.GetStudentAssignmentByID(s.db, studentAssignmentID, studentID)
}

// UpdateAssignmentStatus updates the status of a student assignment by assignment ID
func (s *AssignmentService) UpdateAssignmentStatus(assignmentID uint, studentID uint, status string) error {
	// Validate status
	if status != models.StatusAssigned && status != models.StatusInProgress && status != models.StatusCompleted {
		return errors.New("invalid status")
	}

	// Get student assignment
	studentAssignment, err := s.GetStudentAssignment(assignmentID, studentID)
	if err != nil {
		return err
	}

	// Update status
	return studentAssignment.UpdateStatus(s.db, status)
}

// UpdateStudentAssignmentStatus updates the status of a student assignment by student assignment ID
func (s *AssignmentService) UpdateStudentAssignmentStatus(studentAssignmentID uint, studentID uint, status string) error {
	// Validate status
	if status != models.StatusAssigned && status != models.StatusInProgress && status != models.StatusCompleted {
		return errors.New("invalid status")
	}

	// Get student assignment by ID
	studentAssignment, err := s.GetStudentAssignmentByID(studentAssignmentID, studentID)
	if err != nil {
		return err
	}

	// Update status
	return studentAssignment.UpdateStatus(s.db, status)
}

// MarkAsCompleted marks an assignment as completed by assignment ID
func (s *AssignmentService) MarkAsCompleted(assignmentID uint, studentID uint) error {
	// Get student assignment
	studentAssignment, err := s.GetStudentAssignment(assignmentID, studentID)
	if err != nil {
		return err
	}

	// Mark as completed
	return studentAssignment.MarkAsCompleted(s.db)
}

// MarkAsCompletedByID marks an assignment as completed using student assignment ID
func (s *AssignmentService) MarkAsCompletedByID(studentAssignmentID uint, studentID uint) error {
	// Get student assignment by ID
	studentAssignment, err := s.GetStudentAssignmentByID(studentAssignmentID, studentID)
	if err != nil {
		return err
	}

	// Mark as completed
	return studentAssignment.MarkAsCompleted(s.db)
}

// MarkAsInProgress marks an assignment as in progress by assignment ID
func (s *AssignmentService) MarkAsInProgress(assignmentID uint, studentID uint) error {
	// Get student assignment
	studentAssignment, err := s.GetStudentAssignment(assignmentID, studentID)
	if err != nil {
		return err
	}

	// Mark as in progress
	return studentAssignment.MarkAsInProgress(s.db)
}

// MarkAsInProgressByID marks an assignment as in progress using student assignment ID
func (s *AssignmentService) MarkAsInProgressByID(studentAssignmentID uint, studentID uint) error {
	// Get student assignment by ID
	studentAssignment, err := s.GetStudentAssignmentByID(studentAssignmentID, studentID)
	if err != nil {
		return err
	}

	// Mark as in progress
	return studentAssignment.MarkAsInProgress(s.db)
}

// GetOverdueAssignments retrieves overdue assignments for a student
func (s *AssignmentService) GetOverdueAssignments(studentID uint) ([]models.StudentAssignment, error) {
	// Validate student exists and has student role
	var student models.User
	if err := s.db.First(&student, studentID).Error; err != nil {
		return nil, errors.New("student not found")
	}

	if !student.IsStudent() {
		return nil, errors.New("user is not a student")
	}

	return models.GetOverdueAssignments(s.db, studentID)
}

// GetDashboardStats retrieves dashboard statistics for a student
func (s *AssignmentService) GetDashboardStats(studentID uint) (map[string]int, error) {
	// Validate student exists and has student role
	var student models.User
	if err := s.db.First(&student, studentID).Error; err != nil {
		return nil, errors.New("student not found")
	}

	if !student.IsStudent() {
		return nil, errors.New("user is not a student")
	}

	stats := map[string]int{
		"total_assignments":       0,
		"completed_assignments":   0,
		"in_progress_assignments": 0,
		"overdue_assignments":     0,
	}

	// Get all assignments
	assignments, err := models.GetStudentAssignmentsByStudent(s.db, studentID)
	if err != nil {
		return nil, err
	}

	stats["total_assignments"] = len(assignments)

	// Count by status
	for _, assignment := range assignments {
		switch assignment.Status {
		case models.StatusCompleted:
			stats["completed_assignments"]++
		case models.StatusInProgress:
			stats["in_progress_assignments"]++
		}
	}

	// Count overdue assignments
	overdueAssignments, err := models.GetOverdueAssignments(s.db, studentID)
	if err != nil {
		return nil, err
	}

	stats["overdue_assignments"] = len(overdueAssignments)

	return stats, nil
}

// SearchStudentAssignments searches student assignments by query
func (s *AssignmentService) SearchStudentAssignments(studentID uint, query string) ([]models.StudentAssignment, error) {
	// Validate student exists and has student role
	var student models.User
	if err := s.db.First(&student, studentID).Error; err != nil {
		return nil, errors.New("student not found")
	}

	if !student.IsStudent() {
		return nil, errors.New("user is not a student")
	}

	// Get all student assignments and filter by query
	var studentAssignments []models.StudentAssignment
	searchQuery := "%" + query + "%"

	err := s.db.Preload("Assignment").Preload("Assignment.CreatedBy").
		Joins("JOIN assignments ON assignments.id = student_assignments.assignment_id").
		Where("student_assignments.student_id = ? AND (assignments.title LIKE ? OR assignments.description LIKE ?)",
			studentID, searchQuery, searchQuery).
		Find(&studentAssignments).Error

	if err != nil {
		return nil, err
	}

	return studentAssignments, nil
}

// GetStudentAssignmentsByCategory retrieves student assignments by category
func (s *AssignmentService) GetStudentAssignmentsByCategory(studentID uint, category string) ([]models.StudentAssignment, error) {
	// Validate student exists and has student role
	var student models.User
	if err := s.db.First(&student, studentID).Error; err != nil {
		return nil, errors.New("student not found")
	}

	if !student.IsStudent() {
		return nil, errors.New("user is not a student")
	}

	// Get student assignments by category
	var studentAssignments []models.StudentAssignment

	err := s.db.Preload("Assignment").Preload("Assignment.CreatedBy").
		Joins("JOIN assignments ON assignments.id = student_assignments.assignment_id").
		Where("student_assignments.student_id = ? AND assignments.category = ?", studentID, category).
		Find(&studentAssignments).Error

	if err != nil {
		return nil, err
	}

	return studentAssignments, nil
}

// GetUpcomingAssignments retrieves assignments with upcoming due dates
func (s *AssignmentService) GetUpcomingAssignments(studentID uint, days int) ([]models.StudentAssignment, error) {
	// Validate student exists and has student role
	var student models.User
	if err := s.db.First(&student, studentID).Error; err != nil {
		return nil, errors.New("student not found")
	}

	if !student.IsStudent() {
		return nil, errors.New("user is not a student")
	}

	// Get upcoming assignments
	var studentAssignments []models.StudentAssignment

	err := s.db.Preload("Assignment").Preload("Assignment.CreatedBy").
		Joins("JOIN assignments ON assignments.id = student_assignments.assignment_id").
		Where("student_assignments.student_id = ? AND assignments.due_date IS NOT NULL AND assignments.due_date > NOW() AND assignments.due_date <= DATE_ADD(NOW(), INTERVAL ? DAY) AND student_assignments.status != ?",
			studentID, days, models.StatusCompleted).
		Order("assignments.due_date ASC").
		Find(&studentAssignments).Error

	if err != nil {
		return nil, err
	}

	return studentAssignments, nil
}

// GetRecentlyCompleted retrieves recently completed assignments
func (s *AssignmentService) GetRecentlyCompleted(studentID uint, days int) ([]models.StudentAssignment, error) {
	// Validate student exists and has student role
	var student models.User
	if err := s.db.First(&student, studentID).Error; err != nil {
		return nil, errors.New("student not found")
	}

	if !student.IsStudent() {
		return nil, errors.New("user is not a student")
	}

	// Get recently completed assignments
	var studentAssignments []models.StudentAssignment

	err := s.db.Preload("Assignment").Preload("Assignment.CreatedBy").
		Where("student_id = ? AND status = ? AND completed_at IS NOT NULL AND completed_at >= DATE_SUB(NOW(), INTERVAL ? DAY)",
			studentID, models.StatusCompleted, days).
		Order("completed_at DESC").
		Find(&studentAssignments).Error

	if err != nil {
		return nil, err
	}

	return studentAssignments, nil
}

// GetAssignmentCategories retrieves all categories for student's assignments
func (s *AssignmentService) GetAssignmentCategories(studentID uint) ([]string, error) {
	// Validate student exists and has student role
	var student models.User
	if err := s.db.First(&student, studentID).Error; err != nil {
		return nil, errors.New("student not found")
	}

	if !student.IsStudent() {
		return nil, errors.New("user is not a student")
	}

	// Get distinct categories
	var categories []string

	err := s.db.Model(&models.StudentAssignment{}).
		Select("DISTINCT assignments.category").
		Joins("JOIN assignments ON assignments.id = student_assignments.assignment_id").
		Where("student_assignments.student_id = ? AND assignments.category IS NOT NULL AND assignments.category != ''", studentID).
		Pluck("assignments.category", &categories).Error

	if err != nil {
		return nil, err
	}

	return categories, nil
}

// GetStudentCategories retrieves all categories for student's assignments (alias)
func (s *AssignmentService) GetStudentCategories(studentID uint) ([]string, error) {
	return s.GetAssignmentCategories(studentID)
}

// CalculateStudentStats calculates statistics from student assignments
func (s *AssignmentService) CalculateStudentStats(assignments []models.StudentAssignment) map[string]interface{} {
	stats := map[string]interface{}{
		"total":       len(assignments),
		"completed":   0,
		"in_progress": 0,
		"assigned":    0,
		"overdue":     0,
	}

	overdueCount := 0
	completedCount := 0
	inProgressCount := 0
	assignedCount := 0

	for _, sa := range assignments {
		switch sa.Status {
		case models.StatusCompleted:
			completedCount++
		case models.StatusInProgress:
			inProgressCount++
		case models.StatusAssigned:
			assignedCount++
		}

		// Check if overdue
		if sa.Assignment.DueDate != nil && sa.Status != models.StatusCompleted {
			if time.Now().After(*sa.Assignment.DueDate) {
				overdueCount++
			}
		}
	}

	stats["completed"] = completedCount
	stats["in_progress"] = inProgressCount
	stats["assigned"] = assignedCount
	stats["overdue"] = overdueCount

	// Calculate completion rate
	completionRate := 0.0
	if len(assignments) > 0 {
		completionRate = float64(completedCount) / float64(len(assignments)) * 100
	}
	stats["completion_rate"] = completionRate

	return stats
}

// ============================================================================
// PROGRESS TRACKING & ANALYTICS
// ============================================================================

// GetDetailedProgressReport generates a comprehensive progress report for an assignment
func (s *AssignmentService) GetDetailedProgressReport(assignmentID uint, instructorID uint) (*DetailedProgressReport, error) {
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
func (s *AssignmentService) GetInstructorProgressSummary(instructorID uint) (*InstructorProgressSummary, error) {
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

// GetProgressTrends retrieves progress trends over time (placeholder implementation)
func (s *AssignmentService) GetProgressTrends(instructorID uint, days int) (map[string]interface{}, error) {
	// Validate instructor exists
	var instructor models.User
	if err := s.db.First(&instructor, instructorID).Error; err != nil {
		return nil, errors.New("instructor not found")
	}

	if !instructor.IsInstructor() {
		return nil, errors.New("user is not an instructor")
	}

	trends := make(map[string]interface{})
	trends["days"] = days
	trends["instructor_id"] = instructorID

	// This is a placeholder - implement actual trend analysis as needed
	trends["message"] = "Progress trends feature coming soon"

	return trends, nil
}

// GetCompletionAnalytics retrieves completion analytics (placeholder implementation)
func (s *AssignmentService) GetCompletionAnalytics(instructorID uint) (map[string]interface{}, error) {
	// Validate instructor exists
	var instructor models.User
	if err := s.db.First(&instructor, instructorID).Error; err != nil {
		return nil, errors.New("instructor not found")
	}

	if !instructor.IsInstructor() {
		return nil, errors.New("user is not an instructor")
	}

	analytics := make(map[string]interface{})
	analytics["instructor_id"] = instructorID

	// This is a placeholder - implement actual analytics as needed
	analytics["message"] = "Completion analytics feature coming soon"

	return analytics, nil
}

// getRecentCompletions retrieves recent completion activities (internal helper)
func (s *AssignmentService) getRecentCompletions(instructorID uint, limit int) ([]RecentCompletionActivity, error) {
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

// calculateStudentEngagement calculates student engagement metrics (internal helper)
func (s *AssignmentService) calculateStudentEngagement(instructorID uint) map[string]interface{} {
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

// ============================================================================
// DUE DATE NOTIFICATIONS & ALERTS
// ============================================================================

// GetStudentDueDateAlerts retrieves upcoming due date alerts for a student
func (s *AssignmentService) GetStudentDueDateAlerts(studentID uint, daysAhead int) ([]DueDateAlert, error) {
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
func (s *AssignmentService) GetOverdueDueDateAlerts(studentID uint) ([]DueDateAlert, error) {
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

// GetStudentDueDateSummary provides a comprehensive summary of due date information for a student
func (s *AssignmentService) GetStudentDueDateSummary(studentID uint) (*DueDateSummary, error) {
	summary := &DueDateSummary{}

	// Get upcoming alerts
	upcomingAlerts, err := s.GetStudentDueDateAlerts(studentID, 7)
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

// GetStudentDueDateNotifications retrieves due date notifications for a student
func (s *AssignmentService) GetStudentDueDateNotifications(studentID uint) ([]string, error) {
	var notifications []string

	// Get upcoming alerts
	upcomingAlerts, err := s.GetStudentDueDateAlerts(studentID, 7)
	if err != nil {
		return nil, err
	}

	// Get overdue alerts
	overdueAlerts, err := s.GetOverdueDueDateAlerts(studentID)
	if err != nil {
		return nil, err
	}

	// Generate notification messages
	for _, alert := range upcomingAlerts {
		notifications = append(notifications, s.generateDueDateNotificationMessage(alert))
	}

	for _, alert := range overdueAlerts {
		notifications = append(notifications, s.generateDueDateNotificationMessage(alert))
	}

	return notifications, nil
}

// GetInstructorDueDateOverview provides due date overview for all instructor's assignments
func (s *AssignmentService) GetInstructorDueDateOverview(instructorID uint) (map[string]interface{}, error) {
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

// GetInstructorDueDateNotifications retrieves due date notifications for instructor
func (s *AssignmentService) GetInstructorDueDateNotifications(instructorID uint) ([]string, error) {
	var notifications []string

	// Get due date overview
	overview, err := s.GetInstructorDueDateOverview(instructorID)
	if err != nil {
		return nil, err
	}

	// Generate notifications from overview data
	if upcomingDeadlines, ok := overview["upcoming_deadlines"].([]map[string]interface{}); ok {
		for _, deadline := range upcomingDeadlines {
			title := deadline["title"].(string)
			daysUntilDue := deadline["days_until_due"].(int)
			incompleteCount := deadline["incomplete_count"].(int)

			msg := fmt.Sprintf("Assignment '%s' is due in %d day(s) with %d incomplete submission(s)",
				title, daysUntilDue, incompleteCount)
			notifications = append(notifications, msg)
		}
	}

	if overdueList, ok := overview["overdue_list"].([]map[string]interface{}); ok {
		for _, overdue := range overdueList {
			title := overdue["title"].(string)
			daysOverdue := overdue["days_overdue"].(int)
			incompleteCount := overdue["incomplete_count"].(int)

			msg := fmt.Sprintf("Assignment '%s' is %d day(s) overdue with %d incomplete submission(s)",
				title, daysOverdue, incompleteCount)
			notifications = append(notifications, msg)
		}
	}

	return notifications, nil
}

// generateDueDateNotificationMessage generates a notification message for due date alerts (internal helper)
func (s *AssignmentService) generateDueDateNotificationMessage(alert DueDateAlert) string {
	switch alert.AlertType {
	case "due_today":
		return fmt.Sprintf("Assignment '%s' is due today! Complete it at: %s",
			alert.AssignmentTitle, alert.AssignmentURL)
	case "due_tomorrow":
		return fmt.Sprintf("Assignment '%s' is due tomorrow (%s). Complete it at: %s",
			alert.AssignmentTitle, alert.DueDate.Format("Jan 2"), alert.AssignmentURL)
	case "upcoming":
		return fmt.Sprintf("Assignment '%s' is due in %d days (%s). Complete it at: %s",
			alert.AssignmentTitle, alert.DaysUntilDue, alert.DueDate.Format("Jan 2"), alert.AssignmentURL)
	case "overdue":
		return fmt.Sprintf("Assignment '%s' was due %d days ago (%s). Complete it now at: %s",
			alert.AssignmentTitle, -alert.DaysUntilDue, alert.DueDate.Format("Jan 2"), alert.AssignmentURL)
	default:
		return fmt.Sprintf("Assignment '%s' has an upcoming due date: %s",
			alert.AssignmentTitle, alert.DueDate.Format("Jan 2, 2006"))
	}
}
