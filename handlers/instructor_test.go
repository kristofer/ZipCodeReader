package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
	"zipcodereader/config"
	"zipcodereader/models"
	"zipcodereader/services"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// setupTestDB creates an in-memory SQLite database for testing
func setupTestDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		panic("failed to connect to test database")
	}

	// Auto-migrate test schema
	db.AutoMigrate(&models.User{}, &models.Assignment{}, &models.StudentAssignment{})

	return db
}

// createTestInstructor creates a test instructor user
func createTestInstructor(db *gorm.DB) *models.User {
	user := &models.User{
		Username: "testinstructor",
		Email:    "instructor@test.com",
		Role:     "instructor",
	}
	db.Create(user)
	return user
}

// createTestStudent creates a test student user
func createTestStudent(db *gorm.DB, username string) *models.User {
	user := &models.User{
		Username: username,
		Email:    username + "@test.com",
		Role:     "student",
	}
	db.Create(user)
	return user
}

// createTestAssignment creates a test assignment
func createTestAssignment(db *gorm.DB, instructorID uint, title string) *models.Assignment {
	assignment := &models.Assignment{
		Title:       title,
		Description: "Test assignment",
		URL:         "https://test.com",
		Category:    "Programming",
		CreatedByID: instructorID,
		Type:        "reading",
		EstimatedMinutes: 30,
	}
	db.Create(assignment)
	return assignment
}

// Test InstructorHandlers.ListAssignments
func TestInstructorHandlers_ListAssignments(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := setupTestDB()

	instructor := createTestInstructor(db)
	assignment1 := createTestAssignment(db, instructor.ID, "Assignment 1")
	assignment2 := createTestAssignment(db, instructor.ID, "Assignment 2")

	svc := services.NewAssignmentService(db)
	cfg := &config.Config{UseLocalAuth: true}
	handlers := NewInstructorHandlers(svc, cfg)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Set("user", instructor)
	c.Request = httptest.NewRequest("GET", "/instructor/assignments", nil)

	handlers.ListAssignments(c)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	assert.Equal(t, float64(2), response["total"])
	assignments := response["assignments"].([]interface{})
	assert.Len(t, assignments, 2)

	// Verify assignment IDs are present
	assignmentMap := make(map[float64]bool)
	for _, a := range assignments {
		assignmentMap[a.(map[string]interface{})["id"].(float64)] = true
	}
	assert.True(t, assignmentMap[float64(assignment1.ID)])
	assert.True(t, assignmentMap[float64(assignment2.ID)])
}

// Test InstructorHandlers.CreateAssignment
func TestInstructorHandlers_CreateAssignment(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := setupTestDB()

	instructor := createTestInstructor(db)

	svc := services.NewAssignmentService(db)
	cfg := &config.Config{UseLocalAuth: true}
	handlers := NewInstructorHandlers(svc, cfg)

	requestBody := CreateAssignmentRequest{
		Title:       "New Assignment",
		Description: "Test description",
		URL:         "https://example.com",
		Category:    "Reading",
		DueDate:     time.Now().Add(7 * 24 * time.Hour).Format(time.RFC3339),
	}

	bodyBytes, _ := json.Marshal(requestBody)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Set("user", instructor)
	c.Request = httptest.NewRequest("POST", "/instructor/assignments", bytes.NewBuffer(bodyBytes))
	c.Request.Header.Set("Content-Type", "application/json")

	handlers.CreateAssignment(c)

	assert.Equal(t, http.StatusCreated, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	assert.Equal(t, "Assignment created successfully", response["message"])
	assignment := response["assignment"].(map[string]interface{})
	assert.Equal(t, "New Assignment", assignment["title"])
	assert.Equal(t, "Reading", assignment["category"])
}

// Test InstructorHandlers.GetAssignment
func TestInstructorHandlers_GetAssignment(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := setupTestDB()

	instructor := createTestInstructor(db)
	assignment := createTestAssignment(db, instructor.ID, "Test Assignment")

	svc := services.NewAssignmentService(db)
	cfg := &config.Config{UseLocalAuth: true}
	handlers := NewInstructorHandlers(svc, cfg)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Set("user", instructor)
	c.Params = []gin.Param{{Key: "id", Value: "1"}}
	c.Request = httptest.NewRequest("GET", "/instructor/assignments/1", nil)

	handlers.GetAssignment(c)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	assignmentData := response["assignment"].(map[string]interface{})
	assert.Equal(t, assignment.Title, assignmentData["title"])
	assert.Equal(t, assignment.URL, assignmentData["url"])
}

// Test InstructorHandlers.AssignStudents
func TestInstructorHandlers_AssignStudents(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := setupTestDB()

	instructor := createTestInstructor(db)
	student1 := createTestStudent(db, "student1")
	student2 := createTestStudent(db, "student2")
	assignment := createTestAssignment(db, instructor.ID, "Test Assignment")

	svc := services.NewAssignmentService(db)
	cfg := &config.Config{UseLocalAuth: true}
	handlers := NewInstructorHandlers(svc, cfg)

	requestBody := AssignStudentsRequest{
		StudentIDs: []uint{student1.ID, student2.ID},
	}

	bodyBytes, _ := json.Marshal(requestBody)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Set("user", instructor)
	c.Params = []gin.Param{{Key: "id", Value: "1"}}
	c.Request = httptest.NewRequest("POST", "/instructor/assignments/1/assign", bytes.NewBuffer(bodyBytes))
	c.Request.Header.Set("Content-Type", "application/json")

	handlers.AssignStudents(c)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	assert.Equal(t, float64(2), response["count"])

	// Verify assignments were created in database
	var count int64
	db.Model(&models.StudentAssignment{}).Where("assignment_id = ?", assignment.ID).Count(&count)
	assert.Equal(t, int64(2), count)
}

// Test InstructorHandlers.ListStudents
func TestInstructorHandlers_ListStudents(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := setupTestDB()

	instructor := createTestInstructor(db)
	createTestStudent(db, "student1")
	createTestStudent(db, "student2")
	createTestStudent(db, "student3")

	svc := services.NewAssignmentService(db)
	cfg := &config.Config{UseLocalAuth: true}
	handlers := NewInstructorHandlers(svc, cfg)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Set("user", instructor)
	c.Request = httptest.NewRequest("GET", "/instructor/students", nil)

	handlers.ListStudents(c)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	assert.Equal(t, float64(3), response["total"])
	students := response["students"].([]interface{})
	assert.Len(t, students, 3)
}

// Test InstructorHandlers.GetDashboardStats
func TestInstructorHandlers_GetDashboardStats(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := setupTestDB()

	instructor := createTestInstructor(db)
	createTestAssignment(db, instructor.ID, "Assignment 1")
	createTestAssignment(db, instructor.ID, "Assignment 2")
	createTestStudent(db, "student1")
	createTestStudent(db, "student2")

	svc := services.NewAssignmentService(db)
	cfg := &config.Config{UseLocalAuth: true}
	handlers := NewInstructorHandlers(svc, cfg)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Set("user", instructor)
	c.Request = httptest.NewRequest("GET", "/instructor/dashboard/stats", nil)

	handlers.GetDashboardStats(c)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	assert.Equal(t, float64(2), response["total_assignments"])
	assert.Equal(t, float64(2), response["active_students"])
}

// Test InstructorHandlers.DeleteAssignment
func TestInstructorHandlers_DeleteAssignment(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := setupTestDB()

	instructor := createTestInstructor(db)
	assignment := createTestAssignment(db, instructor.ID, "To Delete")

	svc := services.NewAssignmentService(db)
	cfg := &config.Config{UseLocalAuth: true}
	handlers := NewInstructorHandlers(svc, cfg)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Set("user", instructor)
	c.Params = []gin.Param{{Key: "id", Value: "1"}}
	c.Request = httptest.NewRequest("DELETE", "/instructor/assignments/1", nil)

	handlers.DeleteAssignment(c)

	assert.Equal(t, http.StatusOK, w.Code)

	// Verify assignment is soft deleted
	var deletedAssignment models.Assignment
	result := db.Unscoped().First(&deletedAssignment, assignment.ID)
	assert.NoError(t, result.Error)
	assert.NotNil(t, deletedAssignment.DeletedAt)
}

// Test access control - non-instructor cannot access
func TestInstructorHandlers_AccessControl(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := setupTestDB()

	student := createTestStudent(db, "student1")

	svc := services.NewAssignmentService(db)
	cfg := &config.Config{UseLocalAuth: true}
	handlers := NewInstructorHandlers(svc, cfg)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Set("user", student)
	c.Request = httptest.NewRequest("GET", "/instructor/assignments", nil)

	handlers.ListAssignments(c)

	assert.Equal(t, http.StatusForbidden, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	assert.Equal(t, "Access denied", response["error"])
}
