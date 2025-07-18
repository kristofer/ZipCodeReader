package models

import (
	"time"

	"gorm.io/gorm"
)

// StudentAssignment represents the relationship between a student and an assignment
type StudentAssignment struct {
	ID           uint           `json:"id" gorm:"primaryKey"`
	AssignmentID uint           `json:"assignment_id" gorm:"not null"`
	Assignment   Assignment     `json:"assignment" gorm:"foreignKey:AssignmentID"`
	StudentID    uint           `json:"student_id" gorm:"not null"`
	Student      User           `json:"student" gorm:"foreignKey:StudentID"`
	Status       string         `json:"status" gorm:"default:assigned"` // assigned, in_progress, completed
	CompletedAt  *time.Time     `json:"completed_at"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}

// Assignment status constants
const (
	StatusAssigned   = "assigned"
	StatusInProgress = "in_progress"
	StatusCompleted  = "completed"
)

// CreateStudentAssignment creates a new student assignment
func CreateStudentAssignment(db *gorm.DB, assignmentID, studentID uint) (*StudentAssignment, error) {
	studentAssignment := &StudentAssignment{
		AssignmentID: assignmentID,
		StudentID:    studentID,
		Status:       StatusAssigned,
	}

	result := db.Create(studentAssignment)
	if result.Error != nil {
		return nil, result.Error
	}

	return studentAssignment, nil
}

// GetStudentAssignment retrieves a student assignment by assignment ID and student ID
func GetStudentAssignment(db *gorm.DB, assignmentID, studentID uint) (*StudentAssignment, error) {
	var studentAssignment StudentAssignment
	result := db.Preload("Assignment").Preload("Student").Where("assignment_id = ? AND student_id = ?", assignmentID, studentID).First(&studentAssignment)
	if result.Error != nil {
		return nil, result.Error
	}
	return &studentAssignment, nil
}

// GetStudentAssignmentByID retrieves a student assignment by its ID and student ID
func GetStudentAssignmentByID(db *gorm.DB, studentAssignmentID, studentID uint) (*StudentAssignment, error) {
	var studentAssignment StudentAssignment
	result := db.Preload("Assignment").Preload("Assignment.CreatedBy").Preload("Student").Where("id = ? AND student_id = ?", studentAssignmentID, studentID).First(&studentAssignment)
	if result.Error != nil {
		return nil, result.Error
	}
	return &studentAssignment, nil
}

// GetStudentAssignmentsByStudent retrieves all assignments for a specific student
func GetStudentAssignmentsByStudent(db *gorm.DB, studentID uint) ([]StudentAssignment, error) {
	var studentAssignments []StudentAssignment
	result := db.Preload("Assignment").Preload("Assignment.CreatedBy").Where("student_id = ? AND deleted_at IS NULL", studentID).Find(&studentAssignments)
	if result.Error != nil {
		return nil, result.Error
	}
	return studentAssignments, nil
}

// GetStudentAssignmentsByAssignment retrieves all student assignments for a specific assignment
func GetStudentAssignmentsByAssignment(db *gorm.DB, assignmentID uint) ([]StudentAssignment, error) {
	var studentAssignments []StudentAssignment
	result := db.Preload("Student").Where("assignment_id = ?", assignmentID).Find(&studentAssignments)
	if result.Error != nil {
		return nil, result.Error
	}
	return studentAssignments, nil
}

// UpdateStatus updates the status of a student assignment
func (sa *StudentAssignment) UpdateStatus(db *gorm.DB, status string) error {
	updates := map[string]interface{}{
		"status": status,
	}

	// If marking as completed, set completed_at timestamp
	if status == StatusCompleted {
		now := time.Now()
		updates["completed_at"] = &now
	}

	result := db.Model(sa).Updates(updates)
	return result.Error
}

// MarkAsCompleted marks the assignment as completed
func (sa *StudentAssignment) MarkAsCompleted(db *gorm.DB) error {
	return sa.UpdateStatus(db, StatusCompleted)
}

// MarkAsInProgress marks the assignment as in progress
func (sa *StudentAssignment) MarkAsInProgress(db *gorm.DB) error {
	return sa.UpdateStatus(db, StatusInProgress)
}

// IsCompleted checks if the assignment is completed
func (sa *StudentAssignment) IsCompleted() bool {
	return sa.Status == StatusCompleted
}

// IsOverdue checks if the assignment is overdue
func (sa *StudentAssignment) IsOverdue() bool {
	if sa.Assignment.DueDate == nil {
		return false
	}
	return time.Now().After(*sa.Assignment.DueDate) && !sa.IsCompleted()
}

// GetStudentAssignmentsByStatus retrieves student assignments by status
func GetStudentAssignmentsByStatus(db *gorm.DB, studentID uint, status string) ([]StudentAssignment, error) {
	var studentAssignments []StudentAssignment
	result := db.Preload("Assignment").Preload("Assignment.CreatedBy").Where("student_id = ? AND status = ?", studentID, status).Find(&studentAssignments)
	if result.Error != nil {
		return nil, result.Error
	}
	return studentAssignments, nil
}

// GetOverdueAssignments retrieves overdue assignments for a student
func GetOverdueAssignments(db *gorm.DB, studentID uint) ([]StudentAssignment, error) {
	var studentAssignments []StudentAssignment
	result := db.Preload("Assignment").Preload("Assignment.CreatedBy").
		Joins("JOIN assignments ON assignments.id = student_assignments.assignment_id").
		Where("student_assignments.student_id = ? AND assignments.due_date < ? AND student_assignments.status != ?",
			studentID, time.Now(), StatusCompleted).
		Find(&studentAssignments)
	if result.Error != nil {
		return nil, result.Error
	}
	return studentAssignments, nil
}

// GetAssignmentProgress calculates the completion progress for an assignment
func GetAssignmentProgress(db *gorm.DB, assignmentID uint) (map[string]int, error) {
	var results []struct {
		Status string
		Count  int
	}

	err := db.Model(&StudentAssignment{}).
		Select("status, COUNT(*) as count").
		Where("assignment_id = ?", assignmentID).
		Group("status").
		Scan(&results).Error

	if err != nil {
		return nil, err
	}

	progress := map[string]int{
		StatusAssigned:   0,
		StatusInProgress: 0,
		StatusCompleted:  0,
	}

	for _, result := range results {
		progress[result.Status] = result.Count
	}

	return progress, nil
}

// BulkCreateStudentAssignments creates multiple student assignments at once
func BulkCreateStudentAssignments(db *gorm.DB, assignmentID uint, studentIDs []uint) error {
	var studentAssignments []StudentAssignment

	for _, studentID := range studentIDs {
		studentAssignments = append(studentAssignments, StudentAssignment{
			AssignmentID: assignmentID,
			StudentID:    studentID,
			Status:       StatusAssigned,
		})
	}

	result := db.Create(&studentAssignments)
	return result.Error
}

// RemoveStudentAssignment removes a student assignment (soft delete)
func RemoveStudentAssignment(db *gorm.DB, assignmentID, studentID uint) error {
	result := db.Where("assignment_id = ? AND student_id = ?", assignmentID, studentID).Delete(&StudentAssignment{})
	return result.Error
}
