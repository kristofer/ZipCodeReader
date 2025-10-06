package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"
	"zipcodereader/models"
	"zipcodereader/services"

	"github.com/gin-gonic/gin"
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

// setupTestRouter creates a test router with authentication middleware
func setupTestRouter(handlers *InstructorAssignmentHandlers, user *models.User) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	// Mock auth middleware
	router.Use(func(c *gin.Context) {
		if user != nil {
			c.Set("user", user)
		}
		c.Next()
	})

	// Setup routes
	instructorGroup := router.Group("/instructor")
	{
		instructorGroup.GET("/assignments", handlers.GetAssignments)
		instructorGroup.POST("/assignments", handlers.CreateAssignment)
		instructorGroup.GET("/assignments/:id", handlers.GetAssignment)
		instructorGroup.PUT("/assignments/:id", handlers.UpdateAssignment)
		instructorGroup.DELETE("/assignments/:id", handlers.DeleteAssignment)
		instructorGroup.POST("/assignments/:id/assign", handlers.AssignStudents)
		instructorGroup.GET("/assignments/:id/progress", handlers.GetAssignmentProgress)
		instructorGroup.GET("/assignments/:id/students", handlers.GetAssignmentStudents)
		instructorGroup.DELETE("/assignments/:id/students", handlers.RemoveStudent)
		instructorGroup.GET("/students", handlers.GetAllStudents)
		instructorGroup.GET("/dashboard/stats", handlers.GetDashboardStats)
	}

	return router
}

func TestGetAssignments(t *testing.T) {
	db := setupTestDB(t)
	assignmentService := services.NewAssignmentService(db)
	handlers := NewInstructorAssignmentHandlers(assignmentService)

	instructor := createTestUser(t, db, "instructor1", "instructor")
	student := createTestUser(t, db, "student1", "student")

	// Create test assignments
	input1 := services.CreateAssignmentInput{
		Title:       "Assignment 1",
		Description: "Description 1",
		URL:         "https://example.com/1",
		Category:    "reading",
	}
	assignmentService.CreateAssignment(instructor.ID, input1)

	input2 := services.CreateAssignmentInput{
		Title:       "Assignment 2",
		Description: "Description 2",
		URL:         "https://example.com/2",
		Category:    "homework",
	}
	assignmentService.CreateAssignment(instructor.ID, input2)

	// Test as instructor
	router := setupTestRouter(handlers, instructor)
	req, _ := http.NewRequest("GET", "/instructor/assignments", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	assignments := response["assignments"].([]interface{})
	if len(assignments) != 2 {
		t.Errorf("Expected 2 assignments, got %d", len(assignments))
	}

	// Test as student (should fail)
	router = setupTestRouter(handlers, student)
	req, _ = http.NewRequest("GET", "/instructor/assignments", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusForbidden {
		t.Errorf("Expected status 403, got %d", w.Code)
	}
}

func TestCreateAssignment(t *testing.T) {
	db := setupTestDB(t)
	assignmentService := services.NewAssignmentService(db)
	handlers := NewInstructorAssignmentHandlers(assignmentService)

	instructor := createTestUser(t, db, "instructor1", "instructor")
	student := createTestUser(t, db, "student1", "student")

	// Test creating assignment as instructor
	router := setupTestRouter(handlers, instructor)

	requestBody := CreateAssignmentRequest{
		Title:       "Test Assignment",
		Description: "Test Description",
		URL:         "https://example.com",
		Category:    "reading",
	}
	jsonBody, _ := json.Marshal(requestBody)

	req, _ := http.NewRequest("POST", "/instructor/assignments", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected status 201, got %d", w.Code)
	}

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	if response["message"] != "Assignment created successfully" {
		t.Errorf("Expected success message, got %v", response["message"])
	}

	// Test creating assignment as student (should fail)
	router = setupTestRouter(handlers, student)
	req, _ = http.NewRequest("POST", "/instructor/assignments", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusForbidden {
		t.Errorf("Expected status 403, got %d", w.Code)
	}

	// Test creating assignment with invalid data
	router = setupTestRouter(handlers, instructor)
	invalidRequestBody := CreateAssignmentRequest{
		Title:       "", // Empty title should fail
		Description: "Test Description",
		URL:         "https://example.com",
		Category:    "reading",
	}
	jsonBody, _ = json.Marshal(invalidRequestBody)

	req, _ = http.NewRequest("POST", "/instructor/assignments", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}

func TestGetAssignment(t *testing.T) {
	db := setupTestDB(t)
	assignmentService := services.NewAssignmentService(db)
	handlers := NewInstructorAssignmentHandlers(assignmentService)

	instructor := createTestUser(t, db, "instructor1", "instructor")
	instructor2 := createTestUser(t, db, "instructor2", "instructor")

	// Create test assignment
	input := services.CreateAssignmentInput{
		Title:       "Test Assignment",
		Description: "Test Description",
		URL:         "https://example.com",
		Category:    "reading",
	}
	assignment, _ := assignmentService.CreateAssignment(instructor.ID, input)

	// Test getting assignment as creator
	router := setupTestRouter(handlers, instructor)
	req, _ := http.NewRequest("GET", "/instructor/assignments/"+strconv.Itoa(int(assignment.ID)), nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	assignmentData := response["assignment"].(map[string]interface{})
	if assignmentData["title"] != "Test Assignment" {
		t.Errorf("Expected title 'Test Assignment', got %v", assignmentData["title"])
	}

	// Test getting assignment as different instructor (should fail)
	router = setupTestRouter(handlers, instructor2)
	req, _ = http.NewRequest("GET", "/instructor/assignments/"+strconv.Itoa(int(assignment.ID)), nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status 404, got %d", w.Code)
	}

	// Test getting non-existent assignment
	router = setupTestRouter(handlers, instructor)
	req, _ = http.NewRequest("GET", "/instructor/assignments/999", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status 404, got %d", w.Code)
	}
}

func TestUpdateAssignment(t *testing.T) {
	db := setupTestDB(t)
	assignmentService := services.NewAssignmentService(db)
	handlers := NewInstructorAssignmentHandlers(assignmentService)

	instructor := createTestUser(t, db, "instructor1", "instructor")
	instructor2 := createTestUser(t, db, "instructor2", "instructor")

	// Create test assignment
	input := services.CreateAssignmentInput{
		Title:       "Original Title",
		Description: "Original Description",
		URL:         "https://example.com",
		Category:    "reading",
	}
	assignment, _ := assignmentService.CreateAssignment(instructor.ID, input)

	// Test updating assignment as creator
	router := setupTestRouter(handlers, instructor)

	updateRequest := UpdateAssignmentRequest{
		Title:       "Updated Title",
		Description: "Updated Description",
		URL:         "https://updated.com",
		Category:    "homework",
	}
	jsonBody, _ := json.Marshal(updateRequest)

	req, _ := http.NewRequest("PUT", "/instructor/assignments/"+strconv.Itoa(int(assignment.ID)), bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	if response["message"] != "Assignment updated successfully" {
		t.Errorf("Expected success message, got %v", response["message"])
	}

	// Test updating assignment as different instructor (should fail)
	router = setupTestRouter(handlers, instructor2)
	req, _ = http.NewRequest("PUT", "/instructor/assignments/"+strconv.Itoa(int(assignment.ID)), bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusForbidden {
		t.Errorf("Expected status 403, got %d", w.Code)
	}
}

func TestDeleteAssignment(t *testing.T) {
	db := setupTestDB(t)
	assignmentService := services.NewAssignmentService(db)
	handlers := NewInstructorAssignmentHandlers(assignmentService)

	instructor := createTestUser(t, db, "instructor1", "instructor")
	instructor2 := createTestUser(t, db, "instructor2", "instructor")

	// Create test assignment
	input := services.CreateAssignmentInput{
		Title:       "Test Assignment",
		Description: "Test Description",
		URL:         "https://example.com",
		Category:    "reading",
	}
	assignment, _ := assignmentService.CreateAssignment(instructor.ID, input)

	// Test deleting assignment as different instructor (should fail)
	router := setupTestRouter(handlers, instructor2)
	req, _ := http.NewRequest("DELETE", "/instructor/assignments/"+strconv.Itoa(int(assignment.ID)), nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusForbidden {
		t.Errorf("Expected status 403, got %d", w.Code)
	}

	// Test deleting assignment as creator
	router = setupTestRouter(handlers, instructor)
	req, _ = http.NewRequest("DELETE", "/instructor/assignments/"+strconv.Itoa(int(assignment.ID)), nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	if response["message"] != "Assignment deleted successfully" {
		t.Errorf("Expected success message, got %v", response["message"])
	}
}

func TestAssignStudents(t *testing.T) {
	db := setupTestDB(t)
	assignmentService := services.NewAssignmentService(db)
	handlers := NewInstructorAssignmentHandlers(assignmentService)

	instructor := createTestUser(t, db, "instructor1", "instructor")
	student1 := createTestUser(t, db, "student1", "student")
	student2 := createTestUser(t, db, "student2", "student")

	// Create test assignment
	input := services.CreateAssignmentInput{
		Title:       "Test Assignment",
		Description: "Test Description",
		URL:         "https://example.com",
		Category:    "reading",
	}
	assignment, _ := assignmentService.CreateAssignment(instructor.ID, input)

	// Test assigning students
	router := setupTestRouter(handlers, instructor)

	assignRequest := AssignStudentsRequest{
		StudentIDs: []uint{student1.ID, student2.ID},
	}
	jsonBody, _ := json.Marshal(assignRequest)

	req, _ := http.NewRequest("POST", "/instructor/assignments/"+strconv.Itoa(int(assignment.ID))+"/assign", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	if response["message"] != "Students assigned successfully" {
		t.Errorf("Expected success message, got %v", response["message"])
	}

	// Test assigning with empty student list
	emptyRequest := AssignStudentsRequest{
		StudentIDs: []uint{},
	}
	jsonBody, _ = json.Marshal(emptyRequest)

	req, _ = http.NewRequest("POST", "/instructor/assignments/"+strconv.Itoa(int(assignment.ID))+"/assign", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}

func TestGetAssignmentProgress(t *testing.T) {
	db := setupTestDB(t)
	assignmentService := services.NewAssignmentService(db)
	studentService := services.NewStudentAssignmentService(db)
	handlers := NewInstructorAssignmentHandlers(assignmentService)

	instructor := createTestUser(t, db, "instructor1", "instructor")
	student1 := createTestUser(t, db, "student1", "student")
	student2 := createTestUser(t, db, "student2", "student")
	student3 := createTestUser(t, db, "student3", "student")

	// Create test assignment
	input := services.CreateAssignmentInput{
		Title:       "Test Assignment",
		Description: "Test Description",
		URL:         "https://example.com",
		Category:    "reading",
	}
	assignment, _ := assignmentService.CreateAssignment(instructor.ID, input)

	// Assign students
	assignmentService.AssignToMultipleStudents(assignment.ID, []uint{student1.ID, student2.ID, student3.ID}, instructor.ID)

	// Update some statuses
	studentService.MarkAsInProgress(assignment.ID, student2.ID)
	studentService.MarkAsCompleted(assignment.ID, student3.ID)

	// Test getting progress
	router := setupTestRouter(handlers, instructor)
	req, _ := http.NewRequest("GET", "/instructor/assignments/"+strconv.Itoa(int(assignment.ID))+"/progress", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	progress := response["progress"].(map[string]interface{})
	if int(progress["assigned"].(float64)) != 1 {
		t.Errorf("Expected 1 assigned student, got %v", progress["assigned"])
	}

	if int(progress["in_progress"].(float64)) != 1 {
		t.Errorf("Expected 1 in-progress student, got %v", progress["in_progress"])
	}

	if int(progress["completed"].(float64)) != 1 {
		t.Errorf("Expected 1 completed student, got %v", progress["completed"])
	}

	if int(response["total"].(float64)) != 3 {
		t.Errorf("Expected 3 total students, got %v", response["total"])
	}
}

func TestGetAssignmentStudents(t *testing.T) {
	db := setupTestDB(t)
	assignmentService := services.NewAssignmentService(db)
	handlers := NewInstructorAssignmentHandlers(assignmentService)

	instructor := createTestUser(t, db, "instructor1", "instructor")
	student1 := createTestUser(t, db, "student1", "student")
	student2 := createTestUser(t, db, "student2", "student")

	// Create test assignment
	input := services.CreateAssignmentInput{
		Title:       "Test Assignment",
		Description: "Test Description",
		URL:         "https://example.com",
		Category:    "reading",
	}
	assignment, _ := assignmentService.CreateAssignment(instructor.ID, input)

	// Assign students
	assignmentService.AssignToMultipleStudents(assignment.ID, []uint{student1.ID, student2.ID}, instructor.ID)

	// Test getting assigned students
	router := setupTestRouter(handlers, instructor)
	req, _ := http.NewRequest("GET", "/instructor/assignments/"+strconv.Itoa(int(assignment.ID))+"/students", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	students := response["students"].([]interface{})
	if len(students) != 2 {
		t.Errorf("Expected 2 students, got %d", len(students))
	}

	if int(response["total"].(float64)) != 2 {
		t.Errorf("Expected 2 total students, got %v", response["total"])
	}
}

func TestGetAllStudents(t *testing.T) {
	db := setupTestDB(t)
	assignmentService := services.NewAssignmentService(db)
	handlers := NewInstructorAssignmentHandlers(assignmentService)

	instructor := createTestUser(t, db, "instructor1", "instructor")
	createTestUser(t, db, "student1", "student")
	createTestUser(t, db, "student2", "student")
	createTestUser(t, db, "instructor2", "instructor")

	// Test getting all students
	router := setupTestRouter(handlers, instructor)
	req, _ := http.NewRequest("GET", "/instructor/students", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	students := response["students"].([]interface{})
	if len(students) != 2 {
		t.Errorf("Expected 2 students, got %d", len(students))
	}

	if int(response["total"].(float64)) != 2 {
		t.Errorf("Expected 2 total students, got %v", response["total"])
	}
}

func TestCreateAssignmentWithDueDate(t *testing.T) {
	db := setupTestDB(t)
	assignmentService := services.NewAssignmentService(db)
	handlers := NewInstructorAssignmentHandlers(assignmentService)

	instructor := createTestUser(t, db, "instructor1", "instructor")

	// Test creating assignment with due date
	router := setupTestRouter(handlers, instructor)

	dueDate := time.Now().Add(24 * time.Hour)
	requestBody := CreateAssignmentRequest{
		Title:       "Test Assignment",
		Description: "Test Description",
		URL:         "https://example.com",
		Category:    "reading",
		DueDate:     dueDate.Format(time.RFC3339),
	}
	jsonBody, _ := json.Marshal(requestBody)

	req, _ := http.NewRequest("POST", "/instructor/assignments", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected status 201, got %d", w.Code)
	}

	// Test creating assignment with invalid due date format
	requestBody.DueDate = "invalid-date"
	jsonBody, _ = json.Marshal(requestBody)

	req, _ = http.NewRequest("POST", "/instructor/assignments", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}
