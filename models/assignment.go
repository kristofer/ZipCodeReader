package models

import (
	"time"

	"gorm.io/gorm"
)

// Assignment represents a reading assignment in the system
type Assignment struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	Title       string         `json:"title" gorm:"not null"`
	Description string         `json:"description"`
	URL         string         `json:"url" gorm:"not null"`
	Category    string         `json:"category"`
	DueDate     *time.Time     `json:"due_date"`
	CreatedByID uint           `json:"created_by_id"`
	CreatedBy   User           `json:"created_by" gorm:"foreignKey:CreatedByID"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}

// CreateAssignment creates a new assignment with validation
func CreateAssignment(db *gorm.DB, title, description, url, category string, dueDate *time.Time, createdByID uint) (*Assignment, error) {
	assignment := &Assignment{
		Title:       title,
		Description: description,
		URL:         url,
		Category:    category,
		DueDate:     dueDate,
		CreatedByID: createdByID,
	}

	result := db.Create(assignment)
	if result.Error != nil {
		return nil, result.Error
	}

	return assignment, nil
}

// GetAssignmentByID retrieves an assignment by ID
func GetAssignmentByID(db *gorm.DB, id uint) (*Assignment, error) {
	var assignment Assignment
	result := db.Preload("CreatedBy").First(&assignment, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &assignment, nil
}

// GetAssignmentsByInstructor retrieves all assignments created by a specific instructor
func GetAssignmentsByInstructor(db *gorm.DB, instructorID uint) ([]Assignment, error) {
	var assignments []Assignment
	result := db.Where("created_by_id = ?", instructorID).Find(&assignments)
	if result.Error != nil {
		return nil, result.Error
	}
	return assignments, nil
}

// UpdateAssignment updates an existing assignment
func (a *Assignment) UpdateAssignment(db *gorm.DB, title, description, url, category string, dueDate *time.Time) error {
	updates := map[string]interface{}{
		"title":       title,
		"description": description,
		"url":         url,
		"category":    category,
		"due_date":    dueDate,
	}

	result := db.Model(a).Updates(updates)
	return result.Error
}

// DeleteAssignment soft deletes an assignment
func (a *Assignment) DeleteAssignment(db *gorm.DB) error {
	result := db.Delete(a)
	return result.Error
}

// IsOverdue checks if the assignment is overdue
func (a *Assignment) IsOverdue() bool {
	if a.DueDate == nil {
		return false
	}
	return time.Now().After(*a.DueDate)
}

// GetAssignmentsByCategory retrieves assignments by category
func GetAssignmentsByCategory(db *gorm.DB, category string, instructorID uint) ([]Assignment, error) {
	var assignments []Assignment
	result := db.Where("category = ? AND created_by_id = ?", category, instructorID).Find(&assignments)
	if result.Error != nil {
		return nil, result.Error
	}
	return assignments, nil
}

// SearchAssignments searches assignments by title or description
func SearchAssignments(db *gorm.DB, query string, instructorID uint) ([]Assignment, error) {
	var assignments []Assignment
	searchQuery := "%" + query + "%"
	result := db.Where("(title LIKE ? OR description LIKE ?) AND created_by_id = ?", searchQuery, searchQuery, instructorID).Find(&assignments)
	if result.Error != nil {
		return nil, result.Error
	}
	return assignments, nil
}
