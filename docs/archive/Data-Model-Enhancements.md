# ZipCode Project Data Model Enhancements

## User

scan this project and suggest changes to the data model so that students can track assigned reading. also analyze the data model and make suggestions about more assignment tracking of programming lab assignments.

## Agent

I'll scan the ZipCodeReader project to analyze the current data model and suggest improvements for tracking assigned reading and programming lab assignments.

**Use Tool: list_directory (toolu_01MYLdbTdWQDqEjja4QHwMC9)**
```json
{
  "path": "ZipCodeReader"
}
```

**Tool Results: toolu_01MYLdbTdWQDqEjja4QHwMC9**

# Folders:
ZipCodeReader/config
ZipCodeReader/database
ZipCodeReader/handlers
ZipCodeReader/middleware
ZipCodeReader/models
ZipCodeReader/services
ZipCodeReader/static
ZipCodeReader/templates

# Files:
ZipCodeReader/.gitignore
ZipCodeReader/CLAUDE.md
ZipCodeReader/LICENSE
ZipCodeReader/PHASE2_TESTING.md
ZipCodeReader/PHASE3_COMPLETE.md
ZipCodeReader/README.md
ZipCodeReader/cookies.txt
ZipCodeReader/demo_assignment_flow.sh
ZipCodeReader/go.mod
ZipCodeReader/go.sum
ZipCodeReader/main.go
ZipCodeReader/student_cookies.txt
ZipCodeReader/test_cookies.txt
ZipCodeReader/test_phase3_complete.sh
ZipCodeReader/test_task5_progress_tracking.sh
ZipCodeReader/verify_task5_integration.sh
ZipCodeReader/zipcodereader
ZipCodeReader/zipcodereader.db

## Agent

Let me examine the current data models first:

**Use Tool: list_directory (toolu_0146iPenDxw6sWrwC15iW7jB)**
```json
{
  "path": "ZipCodeReader/models"
}
```

**Tool Results: toolu_0146iPenDxw6sWrwC15iW7jB**


# Files:
ZipCodeReader/models/assignment.go
ZipCodeReader/models/assignment_test.go
ZipCodeReader/models/models.go
ZipCodeReader/models/student_assignment.go
ZipCodeReader/models/student_assignment_test.go
ZipCodeReader/models/user.go

## Agent



**Use Tool: read_file (toolu_01HtxLAYwKyDvZoNMTr2esLY)**
```json
{
  "path": "ZipCodeReader/models/assignment.go"
}
```

**Tool Results: toolu_01HtxLAYwKyDvZoNMTr2esLY**

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

## Agent



**Use Tool: read_file (toolu_01G9hq44gurzRVwH6PgQ6sWj)**
```json
{
  "path": "ZipCodeReader/models/student_assignment.go"
}
```

**Tool Results: toolu_01G9hq44gurzRVwH6PgQ6sWj**

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

// GetStudentAssignmentsByStudent retrieves all assignments for a specific student
func GetStudentAssignmentsByStudent(db *gorm.DB, studentID uint) ([]StudentAssignment, error) {
	var studentAssignments []StudentAssignment
	result := db.Preload("Assignment").Preload("Assignment.CreatedBy").Where("student_id = ?", studentID).Find(&studentAssignments)
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

## Agent



**Use Tool: read_file (toolu_01ErfUhqk3TFSNiyqsYRgk7d)**
```json
{
  "path": "ZipCodeReader/models/user.go"
}
```

**Tool Results: toolu_01ErfUhqk3TFSNiyqsYRgk7d**

package models

import (
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// User represents a user in the system
type User struct {
	ID           uint           `json:"id" gorm:"primaryKey"`
	GitHubID     *int64         `json:"github_id" gorm:"uniqueIndex"` // Made nullable for local auth
	Username     string         `json:"username" gorm:"uniqueIndex;not null"`
	Email        string         `json:"email"`
	AvatarURL    string         `json:"avatar_url"`
	PasswordHash string         `json:"-" gorm:"column:password_hash"` // Hidden from JSON
	Role         string         `json:"role" gorm:"default:student"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}

// IsInstructor checks if the user has instructor role
func (u *User) IsInstructor() bool {
	return u.Role == "instructor"
}

// IsStudent checks if the user has student role
func (u *User) IsStudent() bool {
	return u.Role == "student"
}

// CreateUser creates a new user from GitHub data
func CreateUser(db *gorm.DB, githubID int64, username, email, avatarURL string) (*User, error) {
	user := &User{
		GitHubID:  &githubID,
		Username:  username,
		Email:     email,
		AvatarURL: avatarURL,
		Role:      "student", // Default role
	}

	result := db.Create(user)
	if result.Error != nil {
		return nil, result.Error
	}

	return user, nil
}

// GetUserByGitHubID retrieves a user by their GitHub ID
func GetUserByGitHubID(db *gorm.DB, githubID int64) (*User, error) {
	var user User
	result := db.Where("github_id = ?", githubID).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

// GetUserByID retrieves a user by their ID
func GetUserByID(db *gorm.DB, id uint) (*User, error) {
	var user User
	result := db.First(&user, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

// UpdateUser updates user information
func (u *User) Update(db *gorm.DB) error {
	return db.Save(u).Error
}

// Local Authentication Methods

// SetPassword hashes and sets the password for local authentication
func (u *User) SetPassword(password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.PasswordHash = string(hashedPassword)
	return nil
}

// CheckPassword verifies the password for local authentication
func (u *User) CheckPassword(password string) error {
	return bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password))
}

// IsLocalUser checks if user was created with local authentication
func (u *User) IsLocalUser() bool {
	return u.GitHubID == nil
}

// CreateLocalUser creates a new user with local authentication
func CreateLocalUser(db *gorm.DB, username, email, password string) (*User, error) {
	return CreateLocalUserWithRole(db, username, email, password, "student")
}

// CreateLocalUserWithRole creates a new user with local authentication and specified role
func CreateLocalUserWithRole(db *gorm.DB, username, email, password, role string) (*User, error) {
	// Check if user already exists
	var existingUser User
	if err := db.Where("username = ?", username).First(&existingUser).Error; err == nil {
		return nil, errors.New("user already exists")
	}

	// Validate role
	if role != "student" && role != "instructor" {
		role = "student" // Default to student if invalid role
	}

	user := &User{
		Username: username,
		Email:    email,
		Role:     role,
	}

	// Set password
	if err := user.SetPassword(password); err != nil {
		return nil, err
	}

	result := db.Create(user)
	if result.Error != nil {
		return nil, result.Error
	}

	return user, nil
}

// GetUserByUsername retrieves a user by their username (for local auth)
func GetUserByUsername(db *gorm.DB, username string) (*User, error) {
	var user User
	result := db.Where("username = ?", username).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

// AuthenticateLocalUser authenticates a user with username and password
func AuthenticateLocalUser(db *gorm.DB, username, password string) (*User, error) {
	user, err := GetUserByUsername(db, username)
	if err != nil {
		return nil, err
	}

	if err := user.CheckPassword(password); err != nil {
		return nil, errors.New("invalid credentials")
	}

	return user, nil
}

// GetAllStudents retrieves all users with student role
func GetAllStudents(db *gorm.DB) ([]User, error) {
	var students []User
	result := db.Where("role = ?", "student").Find(&students)
	if result.Error != nil {
		return nil, result.Error
	}
	return students, nil
}

## Agent

Now let me check the services to understand how the data is being used:

**Use Tool: list_directory (toolu_01MvGrFCpamRGL9JP6JiXJS2)**
```json
{
  "path": "ZipCodeReader/services"
}
```

**Tool Results: toolu_01MvGrFCpamRGL9JP6JiXJS2**


# Files:
ZipCodeReader/services/assignment.go
ZipCodeReader/services/assignment_test.go
ZipCodeReader/services/auth.go
ZipCodeReader/services/due_date_notifications.go
ZipCodeReader/services/progress_tracking.go
ZipCodeReader/services/progress_tracking_test.go
ZipCodeReader/services/student_assignment.go
ZipCodeReader/services/student_assignment_test.go

## Agent



**Use Tool: read_file (toolu_01RD5nmZWmkxwaKvPkL9rsr9)**
```json
{
  "path": "ZipCodeReader/services/student_assignment.go"
}
```

**Tool Results: toolu_01RD5nmZWmkxwaKvPkL9rsr9**

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
		"total_assignments":       0,
		"completed_assignments":   0,
		"in_progress_assignments": 0,
		"overdue_assignments":     0,
	}

	// Get all assignments
	assignments, err := models.GetStudentAssignmentsByStudent(s.db, studentID)
	if err != nil {
		return nil, err
	}

	stats["total_assignments"] = len(assignments)

	// Count by status
	for _, assignment := range assignments {
		switch assignment.Status {
		case models.StatusCompleted:
			stats["completed_assignments"]++
		case models.StatusInProgress:
			stats["in_progress_assignments"]++
		}
	}

	// Count overdue assignments
	overdueAssignments, err := models.GetOverdueAssignments(s.db, studentID)
	if err != nil {
		return nil, err
	}

	stats["overdue_assignments"] = len(overdueAssignments)

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

## Agent

Let me also check the database migrations to see the current schema:

**Use Tool: read_file (toolu_01Ae1HnqKa4LMtP3rv2aq6PM)**
```json
{
  "path": "ZipCodeReader/database/migrations.go"
}
```

**Tool Results: toolu_01Ae1HnqKa4LMtP3rv2aq6PM**

package database

import (
	"zipcodereader/models"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Initialize creates a new database connection
func Initialize(databaseURL string) (*gorm.DB, error) {
	// Configure GORM
	config := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	}

	// Connect to SQLite database
	db, err := gorm.Open(sqlite.Open(databaseURL), config)
	if err != nil {
		return nil, err
	}

	// Auto-migrate schemas (will be expanded in later phases)
	err = autoMigrate(db)
	if err != nil {
		return nil, err
	}

	return db, nil
}

// autoMigrate runs database migrations
func autoMigrate(db *gorm.DB) error {
	// Auto-migrate the User model
	err := db.AutoMigrate(&models.User{})
	if err != nil {
		return err
	}

	// Auto-migrate the Assignment model
	err = db.AutoMigrate(&models.Assignment{})
	if err != nil {
		return err
	}

	// Auto-migrate the StudentAssignment model
	err = db.AutoMigrate(&models.StudentAssignment{})
	if err != nil {
		return err
	}

	// Create indexes for better performance
	err = createIndexes(db)
	if err != nil {
		return err
	}

	return nil
}

// createIndexes creates database indexes for better performance
func createIndexes(db *gorm.DB) error {
	// Index on assignments.created_by_id for instructor queries
	err := db.Exec("CREATE INDEX IF NOT EXISTS idx_assignments_created_by ON assignments(created_by_id)").Error
	if err != nil {
		return err
	}

	// Index on assignments.category for category filtering
	err = db.Exec("CREATE INDEX IF NOT EXISTS idx_assignments_category ON assignments(category)").Error
	if err != nil {
		return err
	}

	// Index on assignments.due_date for due date queries
	err = db.Exec("CREATE INDEX IF NOT EXISTS idx_assignments_due_date ON assignments(due_date)").Error
	if err != nil {
		return err
	}

	// Index on student_assignments.student_id for student queries
	err = db.Exec("CREATE INDEX IF NOT EXISTS idx_student_assignments_student ON student_assignments(student_id)").Error
	if err != nil {
		return err
	}

	// Index on student_assignments.assignment_id for assignment queries
	err = db.Exec("CREATE INDEX IF NOT EXISTS idx_student_assignments_assignment ON student_assignments(assignment_id)").Error
	if err != nil {
		return err
	}

	// Index on student_assignments.status for status filtering
	err = db.Exec("CREATE INDEX IF NOT EXISTS idx_student_assignments_status ON student_assignments(status)").Error
	if err != nil {
		return err
	}

	// Composite index for student assignment lookups
	err = db.Exec("CREATE INDEX IF NOT EXISTS idx_student_assignments_composite ON student_assignments(student_id, assignment_id)").Error
	if err != nil {
		return err
	}

	return nil
}

## Agent

Now let me check some handlers to understand how the system is being used:

**Use Tool: grep (toolu_019iFMumtHrCgPBNkYk5JHoq)**
```json
{
  "regex": "type.*Assignment",
  "include_pattern": "**/*.go"
}
```

**Tool Results: toolu_019iFMumtHrCgPBNkYk5JHoq**

Found 10 matches:

## Matches in ZipCodeReader/handlers/instructor_assignments.go

### L15-17
```
type InstructorAssignmentHandlers struct {
	assignmentService *services.AssignmentService
}
```

### L68-74
```
type CreateAssignmentRequest struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description"`
	URL         string `json:"url" binding:"required"`
	Category    string `json:"category"`
	DueDate     string `json:"due_date"` // ISO 8601 format
}
```

### L186-192
```
type UpdateAssignmentRequest struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description"`
	URL         string `json:"url" binding:"required"`
	Category    string `json:"category"`
	DueDate     string `json:"due_date"` // ISO 8601 format
}
```

## Matches in ZipCodeReader/handlers/student_assignments.go

### L14-16
```
type StudentAssignmentHandlers struct {
	studentService *services.StudentAssignmentService
}
```

## Matches in ZipCodeReader/models/assignment.go

### L10-20
```
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
```

2 lines remaining in ancestor node. Read the file to see all.

## Matches in ZipCodeReader/models/student_assignment.go

### L10-20
```
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
```

1 lines remaining in ancestor node. Read the file to see all.

## Matches in ZipCodeReader/services/assignment.go

### L12-14
```
type AssignmentService struct {
	db *gorm.DB
}
```

### L27-33
```
type CreateAssignmentInput struct {
	Title       string
	Description string
	URL         string
	Category    string
	DueDate     *time.Time
}
```

### L115-121
```
type UpdateAssignmentInput struct {
	Title       string
	Description string
	URL         string
	Category    string
	DueDate     *time.Time
}
```

## Matches in ZipCodeReader/services/student_assignment.go

### L11-13
```
type StudentAssignmentService struct {
	db *gorm.DB
}
```

## Agent

<thinking>
Now I have a good understanding of the current data model. Let me analyze what we have and what can be improved for better assignment tracking, especially for reading assignments and programming lab assignments.

## Current Data Model Analysis

### Current Models:
1. **User**: Basic user model with GitHub OAuth2 and local auth support, roles (student/instructor)
2. **Assignment**: Basic assignment model with title, description, URL, category, due date, created by instructor
3. **StudentAssignment**: Junction table tracking student-assignment relationships with status (assigned, in_progress, completed)

### Current Features:
- Basic assignment creation and management
- Student assignment tracking with status
- Due date management
- Category-based organization
- Progress tracking at a basic level
- Search and filtering capabilities

## Issues and Improvement Areas:

### 1. Reading Assignment Tracking Issues:
- **No Reading Progress Tracking**: No way to track how much of a reading assignment a student has completed
- **No Time Tracking**: No tracking of time spent on assignments
- **No Bookmark/Session Management**: Students can't save their reading progress or bookmarks
- **Limited Feedback Mechanism**: No way to record notes, highlights, or comprehension checks
- **No Reading Analytics**: No insights into reading patterns, speed, or engagement

### 2. Programming Lab Assignment Issues:
- **No Assignment Type Differentiation**: All assignments are treated the same, but reading and programming assignments need different tracking
- **No Code Submission**: No way to submit or track programming work
- **No Test Results**: No integration with automated testing or grading
- **No Repository Management**: No Git repository tracking for programming assignments
- **No Build/Compilation Status**: No tracking of whether code compiles or runs successfully
- **No Rubric Support**: No structured grading criteria for programming assignments

### 3. General Data Model Issues:
- **Limited Assignment Metadata**: Missing fields like estimated time, difficulty level, prerequisites
- **No Assignment Dependencies**: Can't track prerequisite assignments
- **No Submission Management**: No formal submission system with versioning
- **Limited Analytics**: Basic progress tracking without detailed insights
- **No Collaborative Features**: No group assignments or peer review capabilities
- **No Learning Objectives**: No way to track learning outcomes or competency mapping

## Suggestions:

### For Reading Assignments:
1. Add reading progress tracking (percentage read, bookmarks)
2. Add time tracking and reading analytics
3. Add note-taking and highlight capabilities
4. Add comprehension check integration
5. Add reading session management

### For Programming Lab Assignments:
1. Add assignment type differentiation
2. Add code submission and versioning
3. Add automated testing integration
4. Add repository management
5. Add build/compilation tracking
6. Add rubric-based grading

### General Improvements:
1. Enhanced assignment metadata
2. Assignment dependency tracking
3. Improved analytics and reporting
4. Learning objective mapping
5. Collaborative assignment support
</thinking>

Based on my analysis of the ZipCodeReader project's data model, I can see that while it has a solid foundation for basic assignment management, there are significant opportunities to enhance it for better tracking of both reading assignments and programming lab assignments. Let me provide detailed recommendations:

## Current Data Model Analysis

The current system has three main models:
- **User**: Basic user management with role-based access (student/instructor)
- **Assignment**: Simple assignment with title, description, URL, category, and due date
- **StudentAssignment**: Junction table tracking assignment status (assigned, in_progress, completed)

## Recommendations for Reading Assignment Tracking

### 1. Enhanced Reading Progress Tracking

**Use Tool: thinking (toolu_019JAsZRPHfD8V313ewhVWTR)**
```json
{
  "content": "I should create specific model suggestions for tracking reading assignments. Let me think about what fields and tables would be most useful:\n\nFor Reading Assignments:\n1. Reading progress tracking (percentage, bookmarks, reading time)\n2. Reading sessions (when they started/stopped reading)\n3. Notes and highlights\n4. Comprehension checks/quizzes\n5. Reading analytics\n\nI'll suggest specific database schema changes and new models."
}
```

**Tool Results: toolu_019JAsZRPHfD8V313ewhVWTR**

Finished thinking.
## Agent

```ZipCodeReader/models/reading_progress.go#L1-35
// ReadingProgress tracks detailed reading progress for assignments
type ReadingProgress struct {
    ID                  uint           `json:"id" gorm:"primaryKey"`
    StudentAssignmentID uint           `json:"student_assignment_id" gorm:"not null"`
    StudentAssignment   StudentAssignment `json:"student_assignment" gorm:"foreignKey:StudentAssignmentID"`
    
    // Progress tracking
    ProgressPercentage  float64        `json:"progress_percentage" gorm:"default:0.0"` // 0.0 to 100.0
    CurrentPage         int            `json:"current_page"`
    TotalPages          int            `json:"total_pages"`
    BookmarkURL         string         `json:"bookmark_url"`
    LastReadPosition    string         `json:"last_read_position"`
    
    // Time tracking
    TotalReadingTime    int            `json:"total_reading_time"` // in minutes
    EstimatedTimeLeft   int            `json:"estimated_time_left"` // in minutes
    AverageReadingSpeed float64        `json:"average_reading_speed"` // words per minute
    
    // Session tracking
    LastReadingSession  *time.Time     `json:"last_reading_session"`
    SessionCount        int            `json:"session_count" gorm:"default:0"`
    
    CreatedAt           time.Time      `json:"created_at"`
    UpdatedAt           time.Time      `json:"updated_at"`
    DeletedAt           gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}

// ReadingSession tracks individual reading sessions
type ReadingSession struct {
    ID                  uint           `json:"id" gorm:"primaryKey"`
    StudentAssignmentID uint           `json:"student_assignment_id" gorm:"not null"`
    StudentAssignment   StudentAssignment `json:"student_assignment" gorm:"foreignKey:StudentAssignmentID"`
    
    StartTime           time.Time      `json:"start_time"`
    EndTime             *time.Time     `json:"end_time"`
    Duration            int            `json:"duration"` // in minutes
    PagesRead           int            `json:"pages_read"`
    StartPosition       string         `json:"start_position"`
    EndPosition         string         `json:"end_position"`
    
    CreatedAt           time.Time      `json:"created_at"`
    UpdatedAt           time.Time      `json:"updated_at"`
}
```

### 2. Note-Taking and Annotation System

```ZipCodeReader/models/reading_annotation.go#L1-25
// ReadingAnnotation stores student notes and highlights
type ReadingAnnotation struct {
    ID                  uint           `json:"id" gorm:"primaryKey"`
    StudentAssignmentID uint           `json:"student_assignment_id" gorm:"not null"`
    StudentAssignment   StudentAssignment `json:"student_assignment" gorm:"foreignKey:StudentAssignmentID"`
    
    Type                string         `json:"type"` // "note", "highlight", "bookmark"
    Content             string         `json:"content"`
    Position            string         `json:"position"` // URL fragment, page number, etc.
    Color               string         `json:"color"` // for highlights
    IsPrivate           bool           `json:"is_private" gorm:"default:true"`
    
    CreatedAt           time.Time      `json:"created_at"`
    UpdatedAt           time.Time      `json:"updated_at"`
    DeletedAt           gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}
```

## Recommendations for Programming Lab Assignment Tracking

### 1. Assignment Type Differentiation

```ZipCodeReader/models/assignment.go#L10-25
// Enhanced Assignment model with type differentiation
type Assignment struct {
    ID          uint           `json:"id" gorm:"primaryKey"`
    Title       string         `json:"title" gorm:"not null"`
    Description string         `json:"description"`
    URL         string         `json:"url" gorm:"not null"`
    Category    string         `json:"category"`
    
    // New fields for better assignment management
    Type                string         `json:"type" gorm:"default:reading"` // "reading", "programming", "quiz", "project"
    DueDate             *time.Time     `json:"due_date"`
    EstimatedDuration   int            `json:"estimated_duration"` // in minutes
    DifficultyLevel     string         `json:"difficulty_level"` // "beginner", "intermediate", "advanced"
    MaxAttempts         int            `json:"max_attempts"`
    
    // Programming-specific fields
    RepositoryURL       string         `json:"repository_url"`
    StarterCodeURL      string         `json:"starter_code_url"`
    TestSuiteURL        string         `json:"test_suite_url"`
    RequiredLanguage    string         `json:"required_language"` // "java", "python", "javascript", etc.
    
    CreatedByID uint      `json:"created_by_id"`
    CreatedBy   User      `json:"created_by" gorm:"foreignKey:CreatedByID"`
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
    DeletedAt   gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}
```

### 2. Programming Assignment Tracking

```ZipCodeReader/models/programming_submission.go#L1-45
// ProgrammingSubmission tracks code submissions and results
type ProgrammingSubmission struct {
    ID                  uint           `json:"id" gorm:"primaryKey"`
    StudentAssignmentID uint           `json:"student_assignment_id" gorm:"not null"`
    StudentAssignment   StudentAssignment `json:"student_assignment" gorm:"foreignKey:StudentAssignmentID"`
    
    // Submission details
    SubmissionNumber    int            `json:"submission_number"`
    GitCommitHash       string         `json:"git_commit_hash"`
    RepositoryURL       string         `json:"repository_url"`
    SubmittedFiles      string         `json:"submitted_files"` // JSON array of file paths
    
    // Build and test results
    BuildStatus         string         `json:"build_status"` // "pending", "success", "failed"
    BuildOutput         string         `json:"build_output"`
    TestStatus          string         `json:"test_status"` // "pending", "passed", "failed", "partial"
    TestResults         string         `json:"test_results"` // JSON with detailed test results
    TestScore           float64        `json:"test_score"` // percentage of tests passed
    
    // Code quality metrics
    LinesOfCode         int            `json:"lines_of_code"`
    Complexity          float64        `json:"complexity"` // cyclomatic complexity
    CodeQualityScore    float64        `json:"code_quality_score"`
    
    // Performance metrics
    ExecutionTime       float64        `json:"execution_time"` // in milliseconds
    MemoryUsage         int            `json:"memory_usage"` // in KB
    
    // Grading
    AutoGradeScore      float64        `json:"auto_grade_score"`
    ManualGradeScore    *float64       `json:"manual_grade_score"`
    FeedbackComments    string         `json:"feedback_comments"`
    GradedBy            *uint          `json:"graded_by"`
    GradedAt            *time.Time     `json:"graded_at"`
    
    CreatedAt           time.Time      `json:"created_at"`
    UpdatedAt           time.Time      `json:"updated_at"`
    DeletedAt           gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}
```

### 3. Assignment Dependencies and Learning Objectives

```ZipCodeReader/models/assignment_dependency.go#L1-20
// AssignmentDependency tracks prerequisite relationships
type AssignmentDependency struct {
    ID                  uint       `json:"id" gorm:"primaryKey"`
    AssignmentID        uint       `json:"assignment_id" gorm:"not null"`
    Assignment          Assignment `json:"assignment" gorm:"foreignKey:AssignmentID"`
    PrerequisiteID      uint       `json:"prerequisite_id" gorm:"not null"`
    Prerequisite        Assignment `json:"prerequisite" gorm:"foreignKey:PrerequisiteID"`
    IsRequired          bool       `json:"is_required" gorm:"default:true"`
    MinimumScore        float64    `json:"minimum_score"` // minimum score needed in prerequisite
    
    CreatedAt           time.Time  `json:"created_at"`
    UpdatedAt           time.Time  `json:"updated_at"`
}

// LearningObjective tracks learning goals
type LearningObjective struct {
    ID          uint       `json:"id" gorm:"primaryKey"`
    Title       string     `json:"title" gorm:"not null"`
    Description string     `json:"description"`
    Category    string     `json:"category"` // "knowledge", "skill", "competency"
    
    CreatedAt   time.Time  `json:"created_at"`
    UpdatedAt   time.Time  `json:"updated_at"`
}

// AssignmentObjective links assignments to learning objectives
type AssignmentObjective struct {
    ID                  uint              `json:"id" gorm:"primaryKey"`
    AssignmentID        uint              `json:"assignment_id" gorm:"not null"`
    Assignment          Assignment        `json:"assignment" gorm:"foreignKey:AssignmentID"`
    LearningObjectiveID uint              `json:"learning_objective_id" gorm:"not null"`
    LearningObjective   LearningObjective `json:"learning_objective" gorm:"foreignKey:LearningObjectiveID"`
    
    CreatedAt           time.Time         `json:"created_at"`
    UpdatedAt           time.Time         `json:"updated_at"`
}
```

## Enhanced StudentAssignment Model

```ZipCodeReader/models/student_assignment.go#L10-35
// Enhanced StudentAssignment with better tracking
type StudentAssignment struct {
    ID           uint           `json:"id" gorm:"primaryKey"`
    AssignmentID uint           `json:"assignment_id" gorm:"not null"`
    Assignment   Assignment     `json:"assignment" gorm:"foreignKey:AssignmentID"`
    StudentID    uint           `json:"student_id" gorm:"not null"`
    Student      User           `json:"student" gorm:"foreignKey:StudentID"`
    
    // Enhanced status tracking
    Status       string         `json:"status" gorm:"default:assigned"` // assigned, started, in_progress, submitted, completed, graded
    StartedAt    *time.Time     `json:"started_at"`
    SubmittedAt  *time.Time     `json:"submitted_at"`
    CompletedAt  *time.Time     `json:"completed_at"`
    
    // Attempt tracking
    AttemptCount int            `json:"attempt_count" gorm:"default:0"`
    BestScore    *float64       `json:"best_score"`
    LastScore    *float64       `json:"last_score"`
    
    // Time tracking
    TimeSpent    int            `json:"time_spent"` // in minutes
    
    // Extensions and special considerations
    DueDateExtension *time.Time  `json:"due_date_extension"`
    ExtensionReason  string      `json:"extension_reason"`
    
    CreatedAt    time.Time      `json:"created_at"`
    UpdatedAt    time.Time      `json:"updated_at"`
    DeletedAt    gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}
```

## Database Migration Updates

```ZipCodeReader/database/migrations.go#L35-60
// Enhanced migrations with new tables
func autoMigrate(db *gorm.DB) error {
    // Existing migrations
    models := []interface{}{
        &models.User{},
        &models.Assignment{},
        &models.StudentAssignment{},
        
        // New models for enhanced tracking
        &models.ReadingProgress{},
        &models.ReadingSession{},
        &models.ReadingAnnotation{},
        &models.ProgrammingSubmission{},
        &models.AssignmentDependency{},
        &models.LearningObjective{},
        &models.AssignmentObjective{},
    }
    
    for _, model := range models {
        if err := db.AutoMigrate(model); err != nil {
            return err
        }
    }
    
    return createEnhancedIndexes(db)
}

func createEnhancedIndexes(db *gorm.DB) error {
    indexes := []string{
        "CREATE INDEX IF NOT EXISTS idx_reading_progress_student_assignment ON reading_progresses(student_assignment_id)",
        "CREATE INDEX IF NOT EXISTS idx_reading_sessions_student_assignment ON reading_sessions(student_assignment_id)",
        "CREATE INDEX IF NOT EXISTS idx_reading_annotations_student_assignment ON reading_annotations(student_assignment_id)",
        "CREATE INDEX IF NOT EXISTS idx_programming_submissions_student_assignment ON programming_submissions(student_assignment_id)",
        "CREATE INDEX IF NOT EXISTS idx_assignment_dependencies_assignment ON assignment_dependencies(assignment_id)",
        "CREATE INDEX IF NOT EXISTS idx_assignment_dependencies_prerequisite ON assignment_dependencies(prerequisite_id)",
        "CREATE INDEX IF NOT EXISTS idx_assignments_type ON assignments(type)",
        "CREATE INDEX IF NOT EXISTS idx_student_assignments_status ON student_assignments(status)",
        "CREATE INDEX IF NOT EXISTS idx_programming_submissions_build_status ON programming_submissions(build_status)",
        "CREATE INDEX IF NOT EXISTS idx_programming_submissions_test_status ON programming_submissions(test_status)",
    }
    
    for _, index := range indexes {
        if err := db.Exec(index).Error; err != nil {
            return err
        }
    }
    
    return nil
}
```

## Key Benefits of These Changes

### For Reading Assignments:
1. **Detailed Progress Tracking**: Track exactly where students are in their reading
2. **Time Analytics**: Understand reading patterns and time investment
3. **Note Integration**: Support student learning with built-in note-taking
4. **Session Management**: Track reading habits and engagement patterns

### For Programming Assignments:
1. **Automated Testing**: Integration with CI/CD pipelines for automatic grading
2. **Version Control**: Track code evolution and submission history
3. **Code Quality**: Monitor code complexity and quality metrics
4. **Performance Tracking**: Measure execution time and memory usage

### General Improvements:
1. **Assignment Dependencies**: Ensure prerequisite completion before new assignments
2. **Learning Objectives**: Map assignments to educational outcomes
3. **Enhanced Analytics**: Better insights for both students and instructors
4. **Flexible Assignment Types**: Support different kinds of learning activities

These changes would transform ZipCodeReader from a basic assignment tracker into a comprehensive learning management system tailored for technical education.

