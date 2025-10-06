package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"zipcodereader/config"
	"zipcodereader/models"
	"zipcodereader/services"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// Test StudentHandlers.ListAssignments
func TestStudentHandlers_ListAssignments(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := setupTestDB()

	instructor := createTestInstructor(db)
	student := createTestStudent(db, "student1")
	assignment1 := createTestAssignment(db, instructor.ID, "Assignment 1")
	assignment2 := createTestAssignment(db, instructor.ID, "Assignment 2")

	// Assign to student
	models.CreateStudentAssignment(db, assignment1.ID, student.ID)
	models.CreateStudentAssignment(db, assignment2.ID, student.ID)

	svc := services.NewAssignmentService(db)
	cfg := &config.Config{UseLocalAuth: true}
	handlers := NewStudentHandlers(svc, cfg)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Set("user", student)
	c.Request = httptest.NewRequest("GET", "/student/assignments", nil)

	handlers.ListAssignments(c)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	assert.Equal(t, float64(2), response["total"])
	assignments := response["assignments"].([]interface{})
	assert.Len(t, assignments, 2)
}

// Test StudentHandlers.GetAssignment
func TestStudentHandlers_GetAssignment(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := setupTestDB()

	instructor := createTestInstructor(db)
	student := createTestStudent(db, "student1")
	assignment := createTestAssignment(db, instructor.ID, "Test Assignment")
	studentAssignment, _ := models.CreateStudentAssignment(db, assignment.ID, student.ID)

	svc := services.NewAssignmentService(db)
	cfg := &config.Config{UseLocalAuth: true}
	handlers := NewStudentHandlers(svc, cfg)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Set("user", student)
	c.Params = []gin.Param{{Key: "id", Value: "1"}}
	c.Request = httptest.NewRequest("GET", "/student/assignments/1", nil)

	handlers.GetAssignment(c)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	assignmentData := response["assignment"].(map[string]interface{})
	assert.Equal(t, float64(studentAssignment.ID), assignmentData["id"])
	assert.Equal(t, models.StatusAssigned, assignmentData["status"])
}

// Test StudentHandlers.MarkCompleted
func TestStudentHandlers_MarkCompleted(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := setupTestDB()

	instructor := createTestInstructor(db)
	student := createTestStudent(db, "student1")
	assignment := createTestAssignment(db, instructor.ID, "Test Assignment")
	studentAssignment, _ := models.CreateStudentAssignment(db, assignment.ID, student.ID)

	svc := services.NewAssignmentService(db)
	cfg := &config.Config{UseLocalAuth: true}
	handlers := NewStudentHandlers(svc, cfg)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Set("user", student)
	c.Params = []gin.Param{{Key: "id", Value: "1"}}
	c.Request = httptest.NewRequest("POST", "/student/assignments/1/complete", nil)

	handlers.MarkCompleted(c)

	assert.Equal(t, http.StatusOK, w.Code)

	// Verify status was updated in database
	var updatedSA models.StudentAssignment
	db.First(&updatedSA, studentAssignment.ID)
	assert.Equal(t, models.StatusCompleted, updatedSA.Status)
	assert.NotNil(t, updatedSA.CompletedAt)
}

// Test StudentHandlers.MarkInProgress
func TestStudentHandlers_MarkInProgress(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := setupTestDB()

	instructor := createTestInstructor(db)
	student := createTestStudent(db, "student1")
	assignment := createTestAssignment(db, instructor.ID, "Test Assignment")
	studentAssignment, _ := models.CreateStudentAssignment(db, assignment.ID, student.ID)

	svc := services.NewAssignmentService(db)
	cfg := &config.Config{UseLocalAuth: true}
	handlers := NewStudentHandlers(svc, cfg)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Set("user", student)
	c.Params = []gin.Param{{Key: "id", Value: "1"}}
	c.Request = httptest.NewRequest("POST", "/student/assignments/1/progress", nil)

	handlers.MarkInProgress(c)

	assert.Equal(t, http.StatusOK, w.Code)

	// Verify status was updated in database
	var updatedSA models.StudentAssignment
	db.First(&updatedSA, studentAssignment.ID)
	assert.Equal(t, models.StatusInProgress, updatedSA.Status)
}

// Test StudentHandlers.UpdateStatus
func TestStudentHandlers_UpdateStatus(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := setupTestDB()

	instructor := createTestInstructor(db)
	student := createTestStudent(db, "student1")
	assignment := createTestAssignment(db, instructor.ID, "Test Assignment")
	studentAssignment, _ := models.CreateStudentAssignment(db, assignment.ID, student.ID)

	svc := services.NewAssignmentService(db)
	cfg := &config.Config{UseLocalAuth: true}
	handlers := NewStudentHandlers(svc, cfg)

	requestBody := UpdateStatusRequest{
		Status: models.StatusInProgress,
	}

	bodyBytes, _ := json.Marshal(requestBody)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Set("user", student)
	c.Params = []gin.Param{{Key: "id", Value: "1"}}
	c.Request = httptest.NewRequest("POST", "/student/assignments/1/status", bytes.NewBuffer(bodyBytes))
	c.Request.Header.Set("Content-Type", "application/json")

	handlers.UpdateStatus(c)

	assert.Equal(t, http.StatusOK, w.Code)

	// Verify status was updated in database
	var updatedSA models.StudentAssignment
	db.First(&updatedSA, studentAssignment.ID)
	assert.Equal(t, models.StatusInProgress, updatedSA.Status)
}

// Test StudentHandlers.GetDashboardStats
func TestStudentHandlers_GetDashboardStats(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := setupTestDB()

	instructor := createTestInstructor(db)
	student := createTestStudent(db, "student1")

	// Create 3 assignments with different statuses
	assignment1 := createTestAssignment(db, instructor.ID, "Assignment 1")
	assignment2 := createTestAssignment(db, instructor.ID, "Assignment 2")
	assignment3 := createTestAssignment(db, instructor.ID, "Assignment 3")

	sa1, _ := models.CreateStudentAssignment(db, assignment1.ID, student.ID)
	sa2, _ := models.CreateStudentAssignment(db, assignment2.ID, student.ID)
	models.CreateStudentAssignment(db, assignment3.ID, student.ID)

	// Mark one as in progress, one as completed
	sa1.MarkAsInProgress(db)
	sa2.MarkAsCompleted(db)

	svc := services.NewAssignmentService(db)
	cfg := &config.Config{UseLocalAuth: true}
	handlers := NewStudentHandlers(svc, cfg)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Set("user", student)
	c.Request = httptest.NewRequest("GET", "/student/dashboard/stats", nil)

	handlers.GetDashboardStats(c)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	assert.Equal(t, float64(3), response["total"])
	assert.Equal(t, float64(1), response["completed"])
	assert.Equal(t, float64(1), response["in_progress"])
	assert.Equal(t, float64(1), response["assigned"])
}

// Test StudentHandlers.GetCategories
func TestStudentHandlers_GetCategories(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := setupTestDB()

	instructor := createTestInstructor(db)
	student := createTestStudent(db, "student1")

	// Create assignments with different categories
	assignment1 := createTestAssignment(db, instructor.ID, "Assignment 1")
	assignment1.Category = "Programming"
	db.Save(assignment1)

	assignment2 := createTestAssignment(db, instructor.ID, "Assignment 2")
	assignment2.Category = "Reading"
	db.Save(assignment2)

	assignment3 := createTestAssignment(db, instructor.ID, "Assignment 3")
	assignment3.Category = "Programming"
	db.Save(assignment3)

	models.CreateStudentAssignment(db, assignment1.ID, student.ID)
	models.CreateStudentAssignment(db, assignment2.ID, student.ID)
	models.CreateStudentAssignment(db, assignment3.ID, student.ID)

	svc := services.NewAssignmentService(db)
	cfg := &config.Config{UseLocalAuth: true}
	handlers := NewStudentHandlers(svc, cfg)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Set("user", student)
	c.Request = httptest.NewRequest("GET", "/student/categories", nil)

	handlers.GetCategories(c)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	categories := response["categories"].([]interface{})
	assert.Contains(t, categories, "Programming")
	assert.Contains(t, categories, "Reading")
}

// Test access control - non-student cannot access
func TestStudentHandlers_AccessControl(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := setupTestDB()

	instructor := createTestInstructor(db)

	svc := services.NewAssignmentService(db)
	cfg := &config.Config{UseLocalAuth: true}
	handlers := NewStudentHandlers(svc, cfg)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Set("user", instructor)
	c.Request = httptest.NewRequest("GET", "/student/assignments", nil)

	handlers.ListAssignments(c)

	assert.Equal(t, http.StatusForbidden, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	assert.Equal(t, "Access denied", response["error"])
}
