package services

import (
	"testing"
	"time"
	"zipcodereader/models"

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
	err = db.AutoMigrate(&models.User{}, &models.Assignment{}, &models.StudentAssignment{})
	if err != nil {
		t.Fatalf("Failed to migrate test database: %v", err)
	}

	return db
}

// createTestUser creates a test user
func createTestUser(t *testing.T, db *gorm.DB, username, role string) *models.User {
	user := &models.User{
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
	service := NewAssignmentService(db)

	instructor := createTestUser(t, db, "instructor1", "instructor")
	student := createTestUser(t, db, "student1", "student")

	// Test creating assignment as instructor
	input := CreateAssignmentInput{
		Title:       "Test Assignment",
		Description: "Test Description",
		URL:         "https://example.com",
		Category:    "reading",
		DueDate:     nil,
	}

	assignment, err := service.CreateAssignment(instructor.ID, input)
	if err != nil {
		t.Fatalf("Failed to create assignment: %v", err)
	}

	if assignment.Title != input.Title {
		t.Errorf("Expected title '%s', got '%s'", input.Title, assignment.Title)
	}

	if assignment.CreatedByID != instructor.ID {
		t.Errorf("Expected created_by_id %d, got %d", instructor.ID, assignment.CreatedByID)
	}

	// Test creating assignment as student (should fail)
	_, err = service.CreateAssignment(student.ID, input)
	if err == nil {
		t.Error("Expected error when student tries to create assignment")
	}

	// Test creating assignment with invalid input
	invalidInput := CreateAssignmentInput{
		Title:       "",
		Description: "Test Description",
		URL:         "https://example.com",
		Category:    "reading",
	}

	_, err = service.CreateAssignment(instructor.ID, invalidInput)
	if err == nil {
		t.Error("Expected error when creating assignment with empty title")
	}

	// Test creating assignment with invalid URL
	invalidInput2 := CreateAssignmentInput{
		Title:       "Test Assignment",
		Description: "Test Description",
		URL:         "",
		Category:    "reading",
	}

	_, err = service.CreateAssignment(instructor.ID, invalidInput2)
	if err == nil {
		t.Error("Expected error when creating assignment with empty URL")
	}
}

func TestGetAssignmentByID(t *testing.T) {
	db := setupTestDB(t)
	service := NewAssignmentService(db)

	instructor := createTestUser(t, db, "instructor1", "instructor")
	student := createTestUser(t, db, "student1", "student")

	// Create test assignment
	input := CreateAssignmentInput{
		Title:       "Test Assignment",
		Description: "Test Description",
		URL:         "https://example.com",
		Category:    "reading",
	}

	assignment, _ := service.CreateAssignment(instructor.ID, input)

	// Test getting assignment as instructor (creator)
	retrievedAssignment, err := service.GetAssignmentByID(assignment.ID, instructor.ID)
	if err != nil {
		t.Fatalf("Failed to get assignment by ID: %v", err)
	}

	if retrievedAssignment.ID != assignment.ID {
		t.Errorf("Expected assignment ID %d, got %d", assignment.ID, retrievedAssignment.ID)
	}

	// Test getting assignment as unassigned student (should fail)
	_, err = service.GetAssignmentByID(assignment.ID, student.ID)
	if err == nil {
		t.Error("Expected error when unassigned student tries to get assignment")
	}

	// Assign student to assignment
	service.AssignToStudent(assignment.ID, student.ID, instructor.ID)

	// Test getting assignment as assigned student (should succeed)
	_, err = service.GetAssignmentByID(assignment.ID, student.ID)
	if err != nil {
		t.Errorf("Expected success when assigned student gets assignment: %v", err)
	}
}

func TestAssignToStudent(t *testing.T) {
	db := setupTestDB(t)
	service := NewAssignmentService(db)

	instructor := createTestUser(t, db, "instructor1", "instructor")
	student := createTestUser(t, db, "student1", "student")
	instructor2 := createTestUser(t, db, "instructor2", "instructor")

	// Create test assignment
	input := CreateAssignmentInput{
		Title:       "Test Assignment",
		Description: "Test Description",
		URL:         "https://example.com",
		Category:    "reading",
	}

	assignment, _ := service.CreateAssignment(instructor.ID, input)

	// Test assigning to student as creator
	err := service.AssignToStudent(assignment.ID, student.ID, instructor.ID)
	if err != nil {
		t.Fatalf("Failed to assign student: %v", err)
	}

	// Test assigning same student again (should fail)
	err = service.AssignToStudent(assignment.ID, student.ID, instructor.ID)
	if err == nil {
		t.Error("Expected error when assigning student twice")
	}

	// Test assigning as different instructor (should fail)
	err = service.AssignToStudent(assignment.ID, student.ID, instructor2.ID)
	if err == nil {
		t.Error("Expected error when different instructor tries to assign")
	}

	// Test assigning non-existent student
	err = service.AssignToStudent(assignment.ID, 999, instructor.ID)
	if err == nil {
		t.Error("Expected error when assigning non-existent student")
	}
}

func TestAssignToMultipleStudents(t *testing.T) {
	db := setupTestDB(t)
	service := NewAssignmentService(db)

	instructor := createTestUser(t, db, "instructor1", "instructor")
	student1 := createTestUser(t, db, "student1", "student")
	student2 := createTestUser(t, db, "student2", "student")
	student3 := createTestUser(t, db, "student3", "student")

	// Create test assignment
	input := CreateAssignmentInput{
		Title:       "Test Assignment",
		Description: "Test Description",
		URL:         "https://example.com",
		Category:    "reading",
	}

	assignment, _ := service.CreateAssignment(instructor.ID, input)

	// Test bulk assigning students
	studentIDs := []uint{student1.ID, student2.ID, student3.ID}
	err := service.AssignToMultipleStudents(assignment.ID, studentIDs, instructor.ID)
	if err != nil {
		t.Fatalf("Failed to assign multiple students: %v", err)
	}

	// Verify all students were assigned
	students, err := service.GetAssignmentStudents(assignment.ID, instructor.ID)
	if err != nil {
		t.Fatalf("Failed to get assignment students: %v", err)
	}

	if len(students) != 3 {
		t.Errorf("Expected 3 assigned students, got %d", len(students))
	}

	// Test assigning already assigned students (should filter out)
	err = service.AssignToMultipleStudents(assignment.ID, studentIDs, instructor.ID)
	if err == nil {
		t.Error("Expected error when all students are already assigned")
	}
}

func TestGetAssignmentsByInstructor(t *testing.T) {
	db := setupTestDB(t)
	service := NewAssignmentService(db)

	instructor1 := createTestUser(t, db, "instructor1", "instructor")
	instructor2 := createTestUser(t, db, "instructor2", "instructor")
	student := createTestUser(t, db, "student1", "student")

	// Create assignments for instructor1
	input1 := CreateAssignmentInput{
		Title:       "Assignment 1",
		Description: "Description 1",
		URL:         "https://example.com/1",
		Category:    "reading",
	}
	service.CreateAssignment(instructor1.ID, input1)

	input2 := CreateAssignmentInput{
		Title:       "Assignment 2",
		Description: "Description 2",
		URL:         "https://example.com/2",
		Category:    "homework",
	}
	service.CreateAssignment(instructor1.ID, input2)

	// Create assignment for instructor2
	input3 := CreateAssignmentInput{
		Title:       "Assignment 3",
		Description: "Description 3",
		URL:         "https://example.com/3",
		Category:    "reading",
	}
	service.CreateAssignment(instructor2.ID, input3)

	// Test getting assignments for instructor1
	assignments, err := service.GetAssignmentsByInstructor(instructor1.ID)
	if err != nil {
		t.Fatalf("Failed to get assignments by instructor: %v", err)
	}

	if len(assignments) != 2 {
		t.Errorf("Expected 2 assignments for instructor1, got %d", len(assignments))
	}

	// Test getting assignments for instructor2
	assignments2, err := service.GetAssignmentsByInstructor(instructor2.ID)
	if err != nil {
		t.Fatalf("Failed to get assignments by instructor: %v", err)
	}

	if len(assignments2) != 1 {
		t.Errorf("Expected 1 assignment for instructor2, got %d", len(assignments2))
	}

	// Test getting assignments as student (should fail)
	_, err = service.GetAssignmentsByInstructor(student.ID)
	if err == nil {
		t.Error("Expected error when student tries to get assignments")
	}
}

func TestUpdateAssignment(t *testing.T) {
	db := setupTestDB(t)
	service := NewAssignmentService(db)

	instructor := createTestUser(t, db, "instructor1", "instructor")
	instructor2 := createTestUser(t, db, "instructor2", "instructor")

	// Create test assignment
	input := CreateAssignmentInput{
		Title:       "Original Title",
		Description: "Original Description",
		URL:         "https://example.com",
		Category:    "reading",
	}

	assignment, _ := service.CreateAssignment(instructor.ID, input)

	// Test updating assignment as creator
	updateInput := UpdateAssignmentInput{
		Title:       "Updated Title",
		Description: "Updated Description",
		URL:         "https://updated.com",
		Category:    "homework",
		DueDate:     &time.Time{},
	}

	err := service.UpdateAssignment(assignment.ID, instructor.ID, updateInput)
	if err != nil {
		t.Fatalf("Failed to update assignment: %v", err)
	}

	// Test updating assignment as different instructor (should fail)
	err = service.UpdateAssignment(assignment.ID, instructor2.ID, updateInput)
	if err == nil {
		t.Error("Expected error when different instructor tries to update")
	}

	// Test updating with invalid input
	invalidInput := UpdateAssignmentInput{
		Title:       "",
		Description: "Updated Description",
		URL:         "https://updated.com",
		Category:    "homework",
	}

	err = service.UpdateAssignment(assignment.ID, instructor.ID, invalidInput)
	if err == nil {
		t.Error("Expected error when updating with empty title")
	}
}

func TestDeleteAssignment(t *testing.T) {
	db := setupTestDB(t)
	service := NewAssignmentService(db)

	instructor := createTestUser(t, db, "instructor1", "instructor")
	instructor2 := createTestUser(t, db, "instructor2", "instructor")

	// Create test assignment
	input := CreateAssignmentInput{
		Title:       "Test Assignment",
		Description: "Test Description",
		URL:         "https://example.com",
		Category:    "reading",
	}

	assignment, _ := service.CreateAssignment(instructor.ID, input)

	// Test deleting assignment as different instructor (should fail)
	err := service.DeleteAssignment(assignment.ID, instructor2.ID)
	if err == nil {
		t.Error("Expected error when different instructor tries to delete")
	}

	// Test deleting assignment as creator
	err = service.DeleteAssignment(assignment.ID, instructor.ID)
	if err != nil {
		t.Fatalf("Failed to delete assignment: %v", err)
	}

	// Test getting deleted assignment (should fail)
	_, err = service.GetAssignmentByID(assignment.ID, instructor.ID)
	if err == nil {
		t.Error("Expected error when getting deleted assignment")
	}
}

func TestGetAssignmentProgress(t *testing.T) {
	db := setupTestDB(t)
	service := NewAssignmentService(db)

	instructor := createTestUser(t, db, "instructor1", "instructor")
	student1 := createTestUser(t, db, "student1", "student")
	student2 := createTestUser(t, db, "student2", "student")
	student3 := createTestUser(t, db, "student3", "student")

	// Create test assignment
	input := CreateAssignmentInput{
		Title:       "Test Assignment",
		Description: "Test Description",
		URL:         "https://example.com",
		Category:    "reading",
	}

	assignment, _ := service.CreateAssignment(instructor.ID, input)

	// Assign to students
	studentIDs := []uint{student1.ID, student2.ID, student3.ID}
	service.AssignToMultipleStudents(assignment.ID, studentIDs, instructor.ID)

	// Update some statuses
	studentService := NewStudentAssignmentService(db)
	studentService.MarkAsInProgress(assignment.ID, student2.ID)
	studentService.MarkAsCompleted(assignment.ID, student3.ID)

	// Test getting progress
	progress, err := service.GetAssignmentProgress(assignment.ID, instructor.ID)
	if err != nil {
		t.Fatalf("Failed to get assignment progress: %v", err)
	}

	if progress[models.StatusAssigned] != 1 {
		t.Errorf("Expected 1 assigned student, got %d", progress[models.StatusAssigned])
	}

	if progress[models.StatusInProgress] != 1 {
		t.Errorf("Expected 1 in-progress student, got %d", progress[models.StatusInProgress])
	}

	if progress[models.StatusCompleted] != 1 {
		t.Errorf("Expected 1 completed student, got %d", progress[models.StatusCompleted])
	}
}

func TestSearchAssignments(t *testing.T) {
	db := setupTestDB(t)
	service := NewAssignmentService(db)

	instructor := createTestUser(t, db, "instructor1", "instructor")

	// Create test assignments
	input1 := CreateAssignmentInput{
		Title:       "JavaScript Basics",
		Description: "Learn JavaScript fundamentals",
		URL:         "https://example.com/js",
		Category:    "programming",
	}
	service.CreateAssignment(instructor.ID, input1)

	input2 := CreateAssignmentInput{
		Title:       "Python Tutorial",
		Description: "Introduction to Python",
		URL:         "https://example.com/python",
		Category:    "programming",
	}
	service.CreateAssignment(instructor.ID, input2)

	// Test search
	assignments, err := service.SearchAssignments("JavaScript", instructor.ID)
	if err != nil {
		t.Fatalf("Failed to search assignments: %v", err)
	}

	if len(assignments) != 1 {
		t.Errorf("Expected 1 assignment for 'JavaScript' search, got %d", len(assignments))
	}

	if assignments[0].Title != "JavaScript Basics" {
		t.Errorf("Expected 'JavaScript Basics', got '%s'", assignments[0].Title)
	}
}

func TestGetAssignmentsByCategory(t *testing.T) {
	db := setupTestDB(t)
	service := NewAssignmentService(db)

	instructor := createTestUser(t, db, "instructor1", "instructor")

	// Create test assignments
	input1 := CreateAssignmentInput{
		Title:       "Assignment 1",
		Description: "Description 1",
		URL:         "https://example.com/1",
		Category:    "reading",
	}
	service.CreateAssignment(instructor.ID, input1)

	input2 := CreateAssignmentInput{
		Title:       "Assignment 2",
		Description: "Description 2",
		URL:         "https://example.com/2",
		Category:    "reading",
	}
	service.CreateAssignment(instructor.ID, input2)

	input3 := CreateAssignmentInput{
		Title:       "Assignment 3",
		Description: "Description 3",
		URL:         "https://example.com/3",
		Category:    "homework",
	}
	service.CreateAssignment(instructor.ID, input3)

	// Test getting assignments by category
	readingAssignments, err := service.GetAssignmentsByCategory("reading", instructor.ID)
	if err != nil {
		t.Fatalf("Failed to get assignments by category: %v", err)
	}

	if len(readingAssignments) != 2 {
		t.Errorf("Expected 2 reading assignments, got %d", len(readingAssignments))
	}

	homeworkAssignments, err := service.GetAssignmentsByCategory("homework", instructor.ID)
	if err != nil {
		t.Fatalf("Failed to get assignments by category: %v", err)
	}

	if len(homeworkAssignments) != 1 {
		t.Errorf("Expected 1 homework assignment, got %d", len(homeworkAssignments))
	}
}

func TestGetAllStudents(t *testing.T) {
	db := setupTestDB(t)
	service := NewAssignmentService(db)

	instructor := createTestUser(t, db, "instructor1", "instructor")
	student1 := createTestUser(t, db, "student1", "student")
	student2 := createTestUser(t, db, "student2", "student")
	createTestUser(t, db, "instructor2", "instructor")

	// Test getting all students
	students, err := service.GetAllStudents(instructor.ID)
	if err != nil {
		t.Fatalf("Failed to get all students: %v", err)
	}

	if len(students) != 2 {
		t.Errorf("Expected 2 students, got %d", len(students))
	}

	// Verify students are returned
	studentNames := []string{students[0].Username, students[1].Username}
	expectedNames := []string{student1.Username, student2.Username}

	for _, expected := range expectedNames {
		found := false
		for _, actual := range studentNames {
			if actual == expected {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected student '%s' not found in results", expected)
		}
	}
}
