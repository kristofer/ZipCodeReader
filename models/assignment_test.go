package models

import (
	"testing"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// setupTestDB creates a test database for testing
func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	// Auto-migrate models
	err = db.AutoMigrate(&User{}, &Assignment{}, &StudentAssignment{})
	if err != nil {
		t.Fatalf("Failed to migrate test database: %v", err)
	}

	return db
}

// createTestUser creates a test user
func createTestUser(t *testing.T, db *gorm.DB, username, role string) *User {
	user := &User{
		Username: username,
		Email:    username + "@example.com",
		Role:     role,
	}

	result := db.Create(user)
	if result.Error != nil {
		t.Fatalf("Failed to create test user: %v", result.Error)
	}

	return user
}

func TestCreateAssignment(t *testing.T) {
	db := setupTestDB(t)
	instructor := createTestUser(t, db, "instructor1", "instructor")

	// Test creating a basic assignment
	assignment, err := CreateAssignment(db, "Test Assignment", "Test Description", "https://example.com", "reading", nil, instructor.ID)
	if err != nil {
		t.Fatalf("Failed to create assignment: %v", err)
	}

	if assignment.Title != "Test Assignment" {
		t.Errorf("Expected title 'Test Assignment', got '%s'", assignment.Title)
	}

	if assignment.CreatedByID != instructor.ID {
		t.Errorf("Expected created_by_id %d, got %d", instructor.ID, assignment.CreatedByID)
	}

	// Test creating assignment with due date
	dueDate := time.Now().Add(24 * time.Hour)
	assignment2, err := CreateAssignment(db, "Test Assignment 2", "Test Description 2", "https://example.com/2", "homework", &dueDate, instructor.ID)
	if err != nil {
		t.Fatalf("Failed to create assignment with due date: %v", err)
	}

	if assignment2.DueDate == nil {
		t.Error("Expected due date to be set")
	}
}

func TestGetAssignmentByID(t *testing.T) {
	db := setupTestDB(t)
	instructor := createTestUser(t, db, "instructor1", "instructor")

	// Create test assignment
	assignment, err := CreateAssignment(db, "Test Assignment", "Test Description", "https://example.com", "reading", nil, instructor.ID)
	if err != nil {
		t.Fatalf("Failed to create assignment: %v", err)
	}

	// Test getting assignment by ID
	retrievedAssignment, err := GetAssignmentByID(db, assignment.ID)
	if err != nil {
		t.Fatalf("Failed to get assignment by ID: %v", err)
	}

	if retrievedAssignment.ID != assignment.ID {
		t.Errorf("Expected assignment ID %d, got %d", assignment.ID, retrievedAssignment.ID)
	}

	if retrievedAssignment.Title != assignment.Title {
		t.Errorf("Expected title '%s', got '%s'", assignment.Title, retrievedAssignment.Title)
	}

	// Test getting non-existent assignment
	_, err = GetAssignmentByID(db, 999)
	if err == nil {
		t.Error("Expected error when getting non-existent assignment")
	}
}

func TestGetAssignmentsByInstructor(t *testing.T) {
	db := setupTestDB(t)
	instructor1 := createTestUser(t, db, "instructor1", "instructor")
	instructor2 := createTestUser(t, db, "instructor2", "instructor")

	// Create assignments for instructor1
	CreateAssignment(db, "Assignment 1", "Description 1", "https://example.com/1", "reading", nil, instructor1.ID)
	CreateAssignment(db, "Assignment 2", "Description 2", "https://example.com/2", "homework", nil, instructor1.ID)

	// Create assignment for instructor2
	CreateAssignment(db, "Assignment 3", "Description 3", "https://example.com/3", "reading", nil, instructor2.ID)

	// Test getting assignments by instructor1
	assignments, err := GetAssignmentsByInstructor(db, instructor1.ID)
	if err != nil {
		t.Fatalf("Failed to get assignments by instructor: %v", err)
	}

	if len(assignments) != 2 {
		t.Errorf("Expected 2 assignments for instructor1, got %d", len(assignments))
	}

	// Test getting assignments by instructor2
	assignments2, err := GetAssignmentsByInstructor(db, instructor2.ID)
	if err != nil {
		t.Fatalf("Failed to get assignments by instructor: %v", err)
	}

	if len(assignments2) != 1 {
		t.Errorf("Expected 1 assignment for instructor2, got %d", len(assignments2))
	}
}

func TestAssignmentIsOverdue(t *testing.T) {
	db := setupTestDB(t)
	instructor := createTestUser(t, db, "instructor1", "instructor")

	// Test assignment without due date
	assignment1, _ := CreateAssignment(db, "Assignment 1", "Description 1", "https://example.com/1", "reading", nil, instructor.ID)
	if assignment1.IsOverdue() {
		t.Error("Assignment without due date should not be overdue")
	}

	// Test assignment with future due date
	futureDate := time.Now().Add(24 * time.Hour)
	assignment2, _ := CreateAssignment(db, "Assignment 2", "Description 2", "https://example.com/2", "reading", &futureDate, instructor.ID)
	if assignment2.IsOverdue() {
		t.Error("Assignment with future due date should not be overdue")
	}

	// Test assignment with past due date
	pastDate := time.Now().Add(-24 * time.Hour)
	assignment3, _ := CreateAssignment(db, "Assignment 3", "Description 3", "https://example.com/3", "reading", &pastDate, instructor.ID)
	if !assignment3.IsOverdue() {
		t.Error("Assignment with past due date should be overdue")
	}
}

func TestUpdateAssignment(t *testing.T) {
	db := setupTestDB(t)
	instructor := createTestUser(t, db, "instructor1", "instructor")

	// Create test assignment
	assignment, err := CreateAssignment(db, "Original Title", "Original Description", "https://example.com", "reading", nil, instructor.ID)
	if err != nil {
		t.Fatalf("Failed to create assignment: %v", err)
	}

	// Update assignment
	newDueDate := time.Now().Add(48 * time.Hour)
	err = assignment.UpdateAssignment(db, "Updated Title", "Updated Description", "https://updated.com", "homework", &newDueDate)
	if err != nil {
		t.Fatalf("Failed to update assignment: %v", err)
	}

	// Verify updates
	updatedAssignment, err := GetAssignmentByID(db, assignment.ID)
	if err != nil {
		t.Fatalf("Failed to get updated assignment: %v", err)
	}

	if updatedAssignment.Title != "Updated Title" {
		t.Errorf("Expected title 'Updated Title', got '%s'", updatedAssignment.Title)
	}

	if updatedAssignment.Description != "Updated Description" {
		t.Errorf("Expected description 'Updated Description', got '%s'", updatedAssignment.Description)
	}

	if updatedAssignment.Category != "homework" {
		t.Errorf("Expected category 'homework', got '%s'", updatedAssignment.Category)
	}
}

func TestDeleteAssignment(t *testing.T) {
	db := setupTestDB(t)
	instructor := createTestUser(t, db, "instructor1", "instructor")

	// Create test assignment
	assignment, err := CreateAssignment(db, "Test Assignment", "Test Description", "https://example.com", "reading", nil, instructor.ID)
	if err != nil {
		t.Fatalf("Failed to create assignment: %v", err)
	}

	// Delete assignment
	err = assignment.DeleteAssignment(db)
	if err != nil {
		t.Fatalf("Failed to delete assignment: %v", err)
	}

	// Verify assignment is deleted (soft delete)
	_, err = GetAssignmentByID(db, assignment.ID)
	if err == nil {
		t.Error("Expected error when getting deleted assignment")
	}
}

func TestSearchAssignments(t *testing.T) {
	db := setupTestDB(t)
	instructor := createTestUser(t, db, "instructor1", "instructor")

	// Create test assignments
	CreateAssignment(db, "JavaScript Basics", "Learn JavaScript fundamentals", "https://example.com/js", "programming", nil, instructor.ID)
	CreateAssignment(db, "Python Tutorial", "Introduction to Python", "https://example.com/python", "programming", nil, instructor.ID)
	CreateAssignment(db, "Reading Assignment", "Read chapter 1", "https://example.com/reading", "reading", nil, instructor.ID)

	// Test search by title
	assignments, err := SearchAssignments(db, "JavaScript", instructor.ID)
	if err != nil {
		t.Fatalf("Failed to search assignments: %v", err)
	}

	if len(assignments) != 1 {
		t.Errorf("Expected 1 assignment for 'JavaScript' search, got %d", len(assignments))
	}

	// Test search by description
	assignments2, err := SearchAssignments(db, "Python", instructor.ID)
	if err != nil {
		t.Fatalf("Failed to search assignments: %v", err)
	}

	if len(assignments2) != 1 {
		t.Errorf("Expected 1 assignment for 'Python' search, got %d", len(assignments2))
	}

	// Test search with no results
	assignments3, err := SearchAssignments(db, "NonexistentTerm", instructor.ID)
	if err != nil {
		t.Fatalf("Failed to search assignments: %v", err)
	}

	if len(assignments3) != 0 {
		t.Errorf("Expected 0 assignments for 'NonexistentTerm' search, got %d", len(assignments3))
	}
}

func TestGetAssignmentsByCategory(t *testing.T) {
	db := setupTestDB(t)
	instructor := createTestUser(t, db, "instructor1", "instructor")

	// Create test assignments
	CreateAssignment(db, "Assignment 1", "Description 1", "https://example.com/1", "reading", nil, instructor.ID)
	CreateAssignment(db, "Assignment 2", "Description 2", "https://example.com/2", "reading", nil, instructor.ID)
	CreateAssignment(db, "Assignment 3", "Description 3", "https://example.com/3", "homework", nil, instructor.ID)

	// Test getting assignments by category
	readingAssignments, err := GetAssignmentsByCategory(db, "reading", instructor.ID)
	if err != nil {
		t.Fatalf("Failed to get assignments by category: %v", err)
	}

	if len(readingAssignments) != 2 {
		t.Errorf("Expected 2 reading assignments, got %d", len(readingAssignments))
	}

	homeworkAssignments, err := GetAssignmentsByCategory(db, "homework", instructor.ID)
	if err != nil {
		t.Fatalf("Failed to get assignments by category: %v", err)
	}

	if len(homeworkAssignments) != 1 {
		t.Errorf("Expected 1 homework assignment, got %d", len(homeworkAssignments))
	}
}
