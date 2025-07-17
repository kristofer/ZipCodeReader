package models

import (
	"testing"
	"time"
)

func TestCreateStudentAssignment(t *testing.T) {
	db := setupTestDB(t)
	instructor := createTestUser(t, db, "instructor1", "instructor")
	student := createTestUser(t, db, "student1", "student")

	// Create test assignment
	assignment, err := CreateAssignment(db, "Test Assignment", "Test Description", "https://example.com", "reading", nil, instructor.ID)
	if err != nil {
		t.Fatalf("Failed to create assignment: %v", err)
	}

	// Create student assignment
	studentAssignment, err := CreateStudentAssignment(db, assignment.ID, student.ID)
	if err != nil {
		t.Fatalf("Failed to create student assignment: %v", err)
	}

	if studentAssignment.AssignmentID != assignment.ID {
		t.Errorf("Expected assignment ID %d, got %d", assignment.ID, studentAssignment.AssignmentID)
	}

	if studentAssignment.StudentID != student.ID {
		t.Errorf("Expected student ID %d, got %d", student.ID, studentAssignment.StudentID)
	}

	if studentAssignment.Status != StatusAssigned {
		t.Errorf("Expected status '%s', got '%s'", StatusAssigned, studentAssignment.Status)
	}
}

func TestGetStudentAssignment(t *testing.T) {
	db := setupTestDB(t)
	instructor := createTestUser(t, db, "instructor1", "instructor")
	student := createTestUser(t, db, "student1", "student")

	// Create test assignment and student assignment
	assignment, _ := CreateAssignment(db, "Test Assignment", "Test Description", "https://example.com", "reading", nil, instructor.ID)
	studentAssignment, _ := CreateStudentAssignment(db, assignment.ID, student.ID)

	// Test getting student assignment
	retrievedSA, err := GetStudentAssignment(db, assignment.ID, student.ID)
	if err != nil {
		t.Fatalf("Failed to get student assignment: %v", err)
	}

	if retrievedSA.ID != studentAssignment.ID {
		t.Errorf("Expected student assignment ID %d, got %d", studentAssignment.ID, retrievedSA.ID)
	}

	// Test getting non-existent student assignment
	_, err = GetStudentAssignment(db, 999, student.ID)
	if err == nil {
		t.Error("Expected error when getting non-existent student assignment")
	}
}

func TestGetStudentAssignmentsByStudent(t *testing.T) {
	db := setupTestDB(t)
	instructor := createTestUser(t, db, "instructor1", "instructor")
	student1 := createTestUser(t, db, "student1", "student")
	student2 := createTestUser(t, db, "student2", "student")

	// Create test assignments
	assignment1, _ := CreateAssignment(db, "Assignment 1", "Description 1", "https://example.com/1", "reading", nil, instructor.ID)
	assignment2, _ := CreateAssignment(db, "Assignment 2", "Description 2", "https://example.com/2", "homework", nil, instructor.ID)

	// Create student assignments
	CreateStudentAssignment(db, assignment1.ID, student1.ID)
	CreateStudentAssignment(db, assignment2.ID, student1.ID)
	CreateStudentAssignment(db, assignment1.ID, student2.ID)

	// Test getting assignments for student1
	assignments, err := GetStudentAssignmentsByStudent(db, student1.ID)
	if err != nil {
		t.Fatalf("Failed to get student assignments: %v", err)
	}

	if len(assignments) != 2 {
		t.Errorf("Expected 2 assignments for student1, got %d", len(assignments))
	}

	// Test getting assignments for student2
	assignments2, err := GetStudentAssignmentsByStudent(db, student2.ID)
	if err != nil {
		t.Fatalf("Failed to get student assignments: %v", err)
	}

	if len(assignments2) != 1 {
		t.Errorf("Expected 1 assignment for student2, got %d", len(assignments2))
	}
}

func TestGetStudentAssignmentsByAssignment(t *testing.T) {
	db := setupTestDB(t)
	instructor := createTestUser(t, db, "instructor1", "instructor")
	student1 := createTestUser(t, db, "student1", "student")
	student2 := createTestUser(t, db, "student2", "student")

	// Create test assignments
	assignment1, _ := CreateAssignment(db, "Assignment 1", "Description 1", "https://example.com/1", "reading", nil, instructor.ID)
	assignment2, _ := CreateAssignment(db, "Assignment 2", "Description 2", "https://example.com/2", "homework", nil, instructor.ID)

	// Create student assignments
	CreateStudentAssignment(db, assignment1.ID, student1.ID)
	CreateStudentAssignment(db, assignment1.ID, student2.ID)
	CreateStudentAssignment(db, assignment2.ID, student1.ID)

	// Test getting students for assignment1
	students, err := GetStudentAssignmentsByAssignment(db, assignment1.ID)
	if err != nil {
		t.Fatalf("Failed to get student assignments: %v", err)
	}

	if len(students) != 2 {
		t.Errorf("Expected 2 students for assignment1, got %d", len(students))
	}

	// Test getting students for assignment2
	students2, err := GetStudentAssignmentsByAssignment(db, assignment2.ID)
	if err != nil {
		t.Fatalf("Failed to get student assignments: %v", err)
	}

	if len(students2) != 1 {
		t.Errorf("Expected 1 student for assignment2, got %d", len(students2))
	}
}

func TestUpdateStatus(t *testing.T) {
	db := setupTestDB(t)
	instructor := createTestUser(t, db, "instructor1", "instructor")
	student := createTestUser(t, db, "student1", "student")

	// Create test assignment and student assignment
	assignment, _ := CreateAssignment(db, "Test Assignment", "Test Description", "https://example.com", "reading", nil, instructor.ID)
	studentAssignment, _ := CreateStudentAssignment(db, assignment.ID, student.ID)

	// Test updating status to in_progress
	err := studentAssignment.UpdateStatus(db, StatusInProgress)
	if err != nil {
		t.Fatalf("Failed to update status: %v", err)
	}

	// Verify status update
	updatedSA, err := GetStudentAssignment(db, assignment.ID, student.ID)
	if err != nil {
		t.Fatalf("Failed to get updated student assignment: %v", err)
	}

	if updatedSA.Status != StatusInProgress {
		t.Errorf("Expected status '%s', got '%s'", StatusInProgress, updatedSA.Status)
	}

	// Test updating status to completed
	err = studentAssignment.UpdateStatus(db, StatusCompleted)
	if err != nil {
		t.Fatalf("Failed to update status to completed: %v", err)
	}

	// Verify completed status and completed_at timestamp
	completedSA, err := GetStudentAssignment(db, assignment.ID, student.ID)
	if err != nil {
		t.Fatalf("Failed to get completed student assignment: %v", err)
	}

	if completedSA.Status != StatusCompleted {
		t.Errorf("Expected status '%s', got '%s'", StatusCompleted, completedSA.Status)
	}

	if completedSA.CompletedAt == nil {
		t.Error("Expected completed_at timestamp to be set")
	}
}

func TestMarkAsCompleted(t *testing.T) {
	db := setupTestDB(t)
	instructor := createTestUser(t, db, "instructor1", "instructor")
	student := createTestUser(t, db, "student1", "student")

	// Create test assignment and student assignment
	assignment, _ := CreateAssignment(db, "Test Assignment", "Test Description", "https://example.com", "reading", nil, instructor.ID)
	studentAssignment, _ := CreateStudentAssignment(db, assignment.ID, student.ID)

	// Mark as completed
	err := studentAssignment.MarkAsCompleted(db)
	if err != nil {
		t.Fatalf("Failed to mark as completed: %v", err)
	}

	// Verify completion
	completedSA, err := GetStudentAssignment(db, assignment.ID, student.ID)
	if err != nil {
		t.Fatalf("Failed to get completed student assignment: %v", err)
	}

	if !completedSA.IsCompleted() {
		t.Error("Expected assignment to be completed")
	}
}

func TestIsOverdue(t *testing.T) {
	db := setupTestDB(t)
	instructor := createTestUser(t, db, "instructor1", "instructor")
	student := createTestUser(t, db, "student1", "student")

	// Test assignment without due date
	assignment1, _ := CreateAssignment(db, "Assignment 1", "Description 1", "https://example.com/1", "reading", nil, instructor.ID)
	CreateStudentAssignment(db, assignment1.ID, student.ID)

	// Load the assignment relation
	sa1, _ := GetStudentAssignment(db, assignment1.ID, student.ID)
	if sa1.IsOverdue() {
		t.Error("Assignment without due date should not be overdue")
	}

	// Test assignment with future due date
	futureDate := time.Now().Add(24 * time.Hour)
	assignment2, _ := CreateAssignment(db, "Assignment 2", "Description 2", "https://example.com/2", "reading", &futureDate, instructor.ID)
	CreateStudentAssignment(db, assignment2.ID, student.ID)

	sa2, _ := GetStudentAssignment(db, assignment2.ID, student.ID)
	if sa2.IsOverdue() {
		t.Error("Assignment with future due date should not be overdue")
	}

	// Test assignment with past due date (not completed)
	pastDate := time.Now().Add(-24 * time.Hour)
	assignment3, _ := CreateAssignment(db, "Assignment 3", "Description 3", "https://example.com/3", "reading", &pastDate, instructor.ID)
	CreateStudentAssignment(db, assignment3.ID, student.ID)

	sa3, _ := GetStudentAssignment(db, assignment3.ID, student.ID)
	if !sa3.IsOverdue() {
		t.Error("Uncompleted assignment with past due date should be overdue")
	}

	// Test completed assignment with past due date
	assignment4, _ := CreateAssignment(db, "Assignment 4", "Description 4", "https://example.com/4", "reading", &pastDate, instructor.ID)
	studentAssignment4, _ := CreateStudentAssignment(db, assignment4.ID, student.ID)
	studentAssignment4.MarkAsCompleted(db)

	sa4, _ := GetStudentAssignment(db, assignment4.ID, student.ID)
	if sa4.IsOverdue() {
		t.Error("Completed assignment should not be overdue")
	}
}

func TestGetStudentAssignmentsByStatus(t *testing.T) {
	db := setupTestDB(t)
	instructor := createTestUser(t, db, "instructor1", "instructor")
	student := createTestUser(t, db, "student1", "student")

	// Create test assignments
	assignment1, _ := CreateAssignment(db, "Assignment 1", "Description 1", "https://example.com/1", "reading", nil, instructor.ID)
	assignment2, _ := CreateAssignment(db, "Assignment 2", "Description 2", "https://example.com/2", "reading", nil, instructor.ID)
	assignment3, _ := CreateAssignment(db, "Assignment 3", "Description 3", "https://example.com/3", "reading", nil, instructor.ID)

	// Create student assignments with different statuses
	CreateStudentAssignment(db, assignment1.ID, student.ID)

	sa2, _ := CreateStudentAssignment(db, assignment2.ID, student.ID)
	sa2.UpdateStatus(db, StatusInProgress)

	sa3, _ := CreateStudentAssignment(db, assignment3.ID, student.ID)
	sa3.MarkAsCompleted(db)

	// Test getting assigned assignments
	assignedAssignments, err := GetStudentAssignmentsByStatus(db, student.ID, StatusAssigned)
	if err != nil {
		t.Fatalf("Failed to get assigned assignments: %v", err)
	}

	if len(assignedAssignments) != 1 {
		t.Errorf("Expected 1 assigned assignment, got %d", len(assignedAssignments))
	}

	// Test getting in-progress assignments
	inProgressAssignments, err := GetStudentAssignmentsByStatus(db, student.ID, StatusInProgress)
	if err != nil {
		t.Fatalf("Failed to get in-progress assignments: %v", err)
	}

	if len(inProgressAssignments) != 1 {
		t.Errorf("Expected 1 in-progress assignment, got %d", len(inProgressAssignments))
	}

	// Test getting completed assignments
	completedAssignments, err := GetStudentAssignmentsByStatus(db, student.ID, StatusCompleted)
	if err != nil {
		t.Fatalf("Failed to get completed assignments: %v", err)
	}

	if len(completedAssignments) != 1 {
		t.Errorf("Expected 1 completed assignment, got %d", len(completedAssignments))
	}
}

func TestGetAssignmentProgress(t *testing.T) {
	db := setupTestDB(t)
	instructor := createTestUser(t, db, "instructor1", "instructor")
	student1 := createTestUser(t, db, "student1", "student")
	student2 := createTestUser(t, db, "student2", "student")
	student3 := createTestUser(t, db, "student3", "student")

	// Create test assignment
	assignment, _ := CreateAssignment(db, "Test Assignment", "Test Description", "https://example.com", "reading", nil, instructor.ID)

	// Create student assignments with different statuses
	CreateStudentAssignment(db, assignment.ID, student1.ID)

	sa2, _ := CreateStudentAssignment(db, assignment.ID, student2.ID)
	sa2.UpdateStatus(db, StatusInProgress)

	sa3, _ := CreateStudentAssignment(db, assignment.ID, student3.ID)
	sa3.MarkAsCompleted(db)

	// Test getting assignment progress
	progress, err := GetAssignmentProgress(db, assignment.ID)
	if err != nil {
		t.Fatalf("Failed to get assignment progress: %v", err)
	}

	if progress[StatusAssigned] != 1 {
		t.Errorf("Expected 1 assigned student, got %d", progress[StatusAssigned])
	}

	if progress[StatusInProgress] != 1 {
		t.Errorf("Expected 1 in-progress student, got %d", progress[StatusInProgress])
	}

	if progress[StatusCompleted] != 1 {
		t.Errorf("Expected 1 completed student, got %d", progress[StatusCompleted])
	}
}

func TestBulkCreateStudentAssignments(t *testing.T) {
	db := setupTestDB(t)
	instructor := createTestUser(t, db, "instructor1", "instructor")
	student1 := createTestUser(t, db, "student1", "student")
	student2 := createTestUser(t, db, "student2", "student")
	student3 := createTestUser(t, db, "student3", "student")

	// Create test assignment
	assignment, _ := CreateAssignment(db, "Test Assignment", "Test Description", "https://example.com", "reading", nil, instructor.ID)

	// Test bulk creating student assignments
	studentIDs := []uint{student1.ID, student2.ID, student3.ID}
	err := BulkCreateStudentAssignments(db, assignment.ID, studentIDs)
	if err != nil {
		t.Fatalf("Failed to bulk create student assignments: %v", err)
	}

	// Verify all student assignments were created
	students, err := GetStudentAssignmentsByAssignment(db, assignment.ID)
	if err != nil {
		t.Fatalf("Failed to get student assignments: %v", err)
	}

	if len(students) != 3 {
		t.Errorf("Expected 3 student assignments, got %d", len(students))
	}

	// Verify all have assigned status
	for _, student := range students {
		if student.Status != StatusAssigned {
			t.Errorf("Expected status '%s', got '%s'", StatusAssigned, student.Status)
		}
	}
}

func TestRemoveStudentAssignment(t *testing.T) {
	db := setupTestDB(t)
	instructor := createTestUser(t, db, "instructor1", "instructor")
	student := createTestUser(t, db, "student1", "student")

	// Create test assignment and student assignment
	assignment, _ := CreateAssignment(db, "Test Assignment", "Test Description", "https://example.com", "reading", nil, instructor.ID)
	CreateStudentAssignment(db, assignment.ID, student.ID)

	// Verify student assignment exists
	_, err := GetStudentAssignment(db, assignment.ID, student.ID)
	if err != nil {
		t.Fatalf("Student assignment should exist: %v", err)
	}

	// Remove student assignment
	err = RemoveStudentAssignment(db, assignment.ID, student.ID)
	if err != nil {
		t.Fatalf("Failed to remove student assignment: %v", err)
	}

	// Verify student assignment is removed
	_, err = GetStudentAssignment(db, assignment.ID, student.ID)
	if err == nil {
		t.Error("Expected error when getting removed student assignment")
	}
}
