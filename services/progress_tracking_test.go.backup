package services

import (
	"testing"
	"time"
	"zipcodereader/models"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupProgressTrackingTestDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to test database")
	}

	// Migrate the schema
	db.AutoMigrate(&models.User{}, &models.Assignment{}, &models.StudentAssignment{})

	return db
}

func createProgressTrackingTestData(db *gorm.DB) (*models.User, *models.User, *models.Assignment) {
	// Create instructor
	instructor := &models.User{
		Username: "instructor",
		Email:    "instructor@test.com",
		Role:     "instructor",
	}
	db.Create(instructor)

	// Create student
	student := &models.User{
		Username: "student",
		Email:    "student@test.com",
		Role:     "student",
	}
	db.Create(student)

	// Create assignment
	dueDate := time.Now().AddDate(0, 0, 7) // Due in 7 days
	assignment := &models.Assignment{
		Title:       "Test Assignment",
		Description: "Test Description",
		URL:         "https://test.com",
		Category:    "test",
		DueDate:     &dueDate,
		CreatedByID: instructor.ID,
	}
	db.Create(assignment)

	return instructor, student, assignment
}

func TestProgressTrackingService_GetDetailedProgressReport(t *testing.T) {
	db := setupProgressTrackingTestDB()
	service := NewProgressTrackingService(db)

	instructor, student, assignment := createProgressTrackingTestData(db)

	// Create student assignment
	studentAssignment := &models.StudentAssignment{
		AssignmentID: assignment.ID,
		StudentID:    student.ID,
		Status:       models.StatusAssigned,
	}
	db.Create(studentAssignment)

	// Test getting detailed progress report
	report, err := service.GetDetailedProgressReport(assignment.ID, instructor.ID)
	if err != nil {
		t.Fatalf("Failed to get detailed progress report: %v", err)
	}

	// Verify report
	if report.AssignmentID != assignment.ID {
		t.Errorf("Expected assignment ID %d, got %d", assignment.ID, report.AssignmentID)
	}

	if report.Title != assignment.Title {
		t.Errorf("Expected title '%s', got '%s'", assignment.Title, report.Title)
	}

	if report.TotalStudents != 1 {
		t.Errorf("Expected 1 student, got %d", report.TotalStudents)
	}

	if report.CompletionRate != 0.0 {
		t.Errorf("Expected completion rate 0.0, got %f", report.CompletionRate)
	}

	if report.StatusBreakdown[models.StatusAssigned] != 1 {
		t.Errorf("Expected 1 assigned student, got %d", report.StatusBreakdown[models.StatusAssigned])
	}

	if len(report.StudentDetails) != 1 {
		t.Errorf("Expected 1 student detail, got %d", len(report.StudentDetails))
	}

	// Test access control
	otherInstructor := &models.User{
		Username: "other_instructor",
		Email:    "other@test.com",
		Role:     "instructor",
	}
	db.Create(otherInstructor)

	_, err = service.GetDetailedProgressReport(assignment.ID, otherInstructor.ID)
	if err == nil {
		t.Error("Expected access denied error for other instructor")
	}
}

func TestProgressTrackingService_GetInstructorProgressSummary(t *testing.T) {
	db := setupProgressTrackingTestDB()
	service := NewProgressTrackingService(db)

	instructor, student, assignment := createProgressTrackingTestData(db)

	// Create student assignment
	studentAssignment := &models.StudentAssignment{
		AssignmentID: assignment.ID,
		StudentID:    student.ID,
		Status:       models.StatusAssigned,
	}
	db.Create(studentAssignment)

	// Test getting instructor progress summary
	summary, err := service.GetInstructorProgressSummary(instructor.ID)
	if err != nil {
		t.Fatalf("Failed to get instructor progress summary: %v", err)
	}

	// Verify summary
	if summary.TotalAssignments != 1 {
		t.Errorf("Expected 1 assignment, got %d", summary.TotalAssignments)
	}

	if summary.TotalStudentAssignments != 1 {
		t.Errorf("Expected 1 student assignment, got %d", summary.TotalStudentAssignments)
	}

	if summary.OverallCompletionRate != 0.0 {
		t.Errorf("Expected completion rate 0.0, got %f", summary.OverallCompletionRate)
	}

	if summary.AssignmentsWithDueDates != 1 {
		t.Errorf("Expected 1 assignment with due date, got %d", summary.AssignmentsWithDueDates)
	}

	// Check category breakdown
	if len(summary.CategoryBreakdown) != 1 {
		t.Errorf("Expected 1 category, got %d", len(summary.CategoryBreakdown))
	}

	if stats, exists := summary.CategoryBreakdown["test"]; !exists {
		t.Error("Expected 'test' category to exist")
	} else {
		if stats.AssignmentCount != 1 {
			t.Errorf("Expected 1 assignment in 'test' category, got %d", stats.AssignmentCount)
		}
	}

	// Check student engagement
	if engagement, exists := summary.StudentEngagement["active_students"]; !exists {
		t.Error("Expected 'active_students' in engagement metrics")
	} else {
		if engagement.(int64) != 1 {
			t.Errorf("Expected 1 active student, got %v", engagement)
		}
	}
}

func TestProgressTrackingService_WithCompletedAssignment(t *testing.T) {
	db := setupProgressTrackingTestDB()
	service := NewProgressTrackingService(db)

	instructor, student, assignment := createProgressTrackingTestData(db)

	// Create and complete student assignment
	completedAt := time.Now()
	studentAssignment := &models.StudentAssignment{
		AssignmentID: assignment.ID,
		StudentID:    student.ID,
		Status:       models.StatusCompleted,
		CompletedAt:  &completedAt,
	}
	db.Create(studentAssignment)

	// Test getting detailed progress report
	report, err := service.GetDetailedProgressReport(assignment.ID, instructor.ID)
	if err != nil {
		t.Fatalf("Failed to get detailed progress report: %v", err)
	}

	// Verify completion metrics
	if report.CompletionRate != 100.0 {
		t.Errorf("Expected completion rate 100.0, got %f", report.CompletionRate)
	}

	if report.StatusBreakdown[models.StatusCompleted] != 1 {
		t.Errorf("Expected 1 completed student, got %d", report.StatusBreakdown[models.StatusCompleted])
	}

	// Check student details
	if len(report.StudentDetails) != 1 {
		t.Fatalf("Expected 1 student detail, got %d", len(report.StudentDetails))
	}

	studentDetail := report.StudentDetails[0]
	if studentDetail.Status != models.StatusCompleted {
		t.Errorf("Expected status '%s', got '%s'", models.StatusCompleted, studentDetail.Status)
	}

	if studentDetail.CompletedAt == nil {
		t.Error("Expected CompletedAt to be set")
	}

	if studentDetail.TimeToComplete == nil {
		t.Error("Expected TimeToComplete to be set")
	}

	// Test instructor summary
	summary, err := service.GetInstructorProgressSummary(instructor.ID)
	if err != nil {
		t.Fatalf("Failed to get instructor progress summary: %v", err)
	}

	if summary.OverallCompletionRate != 100.0 {
		t.Errorf("Expected completion rate 100.0, got %f", summary.OverallCompletionRate)
	}

	if len(summary.RecentCompletions) != 1 {
		t.Errorf("Expected 1 recent completion, got %d", len(summary.RecentCompletions))
	}
}

func TestProgressTrackingService_MultipleStudents(t *testing.T) {
	db := setupProgressTrackingTestDB()
	service := NewProgressTrackingService(db)

	instructor, student1, assignment := createProgressTrackingTestData(db)

	// Create second student
	student2 := &models.User{
		Username: "student2",
		Email:    "student2@test.com",
		Role:     "student",
	}
	db.Create(student2)

	// Create student assignments with different statuses
	studentAssignment1 := &models.StudentAssignment{
		AssignmentID: assignment.ID,
		StudentID:    student1.ID,
		Status:       models.StatusAssigned,
	}
	db.Create(studentAssignment1)

	completedAt := time.Now()
	studentAssignment2 := &models.StudentAssignment{
		AssignmentID: assignment.ID,
		StudentID:    student2.ID,
		Status:       models.StatusCompleted,
		CompletedAt:  &completedAt,
	}
	db.Create(studentAssignment2)

	// Test getting detailed progress report
	report, err := service.GetDetailedProgressReport(assignment.ID, instructor.ID)
	if err != nil {
		t.Fatalf("Failed to get detailed progress report: %v", err)
	}

	// Verify metrics
	if report.TotalStudents != 2 {
		t.Errorf("Expected 2 students, got %d", report.TotalStudents)
	}

	if report.CompletionRate != 50.0 {
		t.Errorf("Expected completion rate 50.0, got %f", report.CompletionRate)
	}

	if report.StatusBreakdown[models.StatusAssigned] != 1 {
		t.Errorf("Expected 1 assigned student, got %d", report.StatusBreakdown[models.StatusAssigned])
	}

	if report.StatusBreakdown[models.StatusCompleted] != 1 {
		t.Errorf("Expected 1 completed student, got %d", report.StatusBreakdown[models.StatusCompleted])
	}

	if len(report.StudentDetails) != 2 {
		t.Errorf("Expected 2 student details, got %d", len(report.StudentDetails))
	}
}
