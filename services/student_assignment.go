package services

import (
	"errors"
	"zipcodereader/models"

	"gorm.io/gorm"
)

// StudentAssignmentService handles business logic for student assignments
type StudentAssignmentService struct {
	db *gorm.DB
}

// NewStudentAssignmentService creates a new student assignment service
func NewStudentAssignmentService(db *gorm.DB) *StudentAssignmentService {
	return &StudentAssignmentService{db: db}
}

// GetStudentAssignments retrieves all assignments for a student
func (s *StudentAssignmentService) GetStudentAssignments(studentID uint) ([]models.StudentAssignment, error) {
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
func (s *StudentAssignmentService) GetStudentAssignmentsByStatus(studentID uint, status string) ([]models.StudentAssignment, error) {
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

// GetStudentAssignment retrieves a specific student assignment
func (s *StudentAssignmentService) GetStudentAssignment(assignmentID uint, studentID uint) (*models.StudentAssignment, error) {
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

// UpdateAssignmentStatus updates the status of a student assignment
func (s *StudentAssignmentService) UpdateAssignmentStatus(assignmentID uint, studentID uint, status string) error {
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

// MarkAsCompleted marks an assignment as completed
func (s *StudentAssignmentService) MarkAsCompleted(assignmentID uint, studentID uint) error {
	// Get student assignment
	studentAssignment, err := s.GetStudentAssignment(assignmentID, studentID)
	if err != nil {
		return err
	}

	// Mark as completed
	return studentAssignment.MarkAsCompleted(s.db)
}

// MarkAsInProgress marks an assignment as in progress
func (s *StudentAssignmentService) MarkAsInProgress(assignmentID uint, studentID uint) error {
	// Get student assignment
	studentAssignment, err := s.GetStudentAssignment(assignmentID, studentID)
	if err != nil {
		return err
	}

	// Mark as in progress
	return studentAssignment.MarkAsInProgress(s.db)
}

// GetOverdueAssignments retrieves overdue assignments for a student
func (s *StudentAssignmentService) GetOverdueAssignments(studentID uint) ([]models.StudentAssignment, error) {
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
func (s *StudentAssignmentService) GetDashboardStats(studentID uint) (map[string]int, error) {
	// Validate student exists and has student role
	var student models.User
	if err := s.db.First(&student, studentID).Error; err != nil {
		return nil, errors.New("student not found")
	}

	if !student.IsStudent() {
		return nil, errors.New("user is not a student")
	}

	stats := map[string]int{
		"total":       0,
		"assigned":    0,
		"in_progress": 0,
		"completed":   0,
		"overdue":     0,
	}

	// Get all assignments
	assignments, err := models.GetStudentAssignmentsByStudent(s.db, studentID)
	if err != nil {
		return nil, err
	}

	stats["total"] = len(assignments)

	// Count by status
	for _, assignment := range assignments {
		stats[assignment.Status]++
	}

	// Count overdue assignments
	overdueAssignments, err := models.GetOverdueAssignments(s.db, studentID)
	if err != nil {
		return nil, err
	}

	stats["overdue"] = len(overdueAssignments)

	return stats, nil
}

// SearchStudentAssignments searches student assignments by query
func (s *StudentAssignmentService) SearchStudentAssignments(studentID uint, query string) ([]models.StudentAssignment, error) {
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
func (s *StudentAssignmentService) GetStudentAssignmentsByCategory(studentID uint, category string) ([]models.StudentAssignment, error) {
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
func (s *StudentAssignmentService) GetUpcomingAssignments(studentID uint, days int) ([]models.StudentAssignment, error) {
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
func (s *StudentAssignmentService) GetRecentlyCompleted(studentID uint, days int) ([]models.StudentAssignment, error) {
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
func (s *StudentAssignmentService) GetAssignmentCategories(studentID uint) ([]string, error) {
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
