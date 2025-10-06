package services

import (
	"errors"
	"time"
	"zipcodereader/models"

	"gorm.io/gorm"
)

// AssignmentService handles business logic for assignments
type AssignmentService struct {
	db *gorm.DB
}

// NewAssignmentService creates a new assignment service
func NewAssignmentService(db *gorm.DB) *AssignmentService {
	return &AssignmentService{db: db}
}

// GetDB returns the database instance
func (s *AssignmentService) GetDB() *gorm.DB {
	return s.db
}

// CreateAssignmentInput represents input for creating an assignment
type CreateAssignmentInput struct {
	Title       string
	Description string
	URL         string
	Category    string
	DueDate     *time.Time
}

// CreateAssignment creates a new assignment with validation
func (s *AssignmentService) CreateAssignment(instructorID uint, input CreateAssignmentInput) (*models.Assignment, error) {
	// Validate instructor exists and has instructor role
	var instructor models.User
	if err := s.db.First(&instructor, instructorID).Error; err != nil {
		return nil, errors.New("instructor not found")
	}

	if !instructor.IsInstructor() {
		return nil, errors.New("user is not an instructor")
	}

	// Validate input
	if input.Title == "" {
		return nil, errors.New("title is required")
	}

	if input.URL == "" {
		return nil, errors.New("URL is required")
	}

	// Create assignment
	assignment, err := models.CreateAssignment(s.db, input.Title, input.Description, input.URL, input.Category, input.DueDate, instructorID)
	if err != nil {
		return nil, err
	}

	return assignment, nil
}

// GetAssignmentByID retrieves an assignment by ID with authorization check
func (s *AssignmentService) GetAssignmentByID(assignmentID uint, userID uint) (*models.Assignment, error) {
	assignment, err := models.GetAssignmentByID(s.db, assignmentID)
	if err != nil {
		return nil, err
	}

	// Check if user is the creator or has access to this assignment
	if assignment.CreatedByID != userID {
		// Check if user is a student assigned to this assignment
		var user models.User
		if err := s.db.First(&user, userID).Error; err != nil {
			return nil, errors.New("user not found")
		}

		if user.IsStudent() {
			// Check if student is assigned to this assignment
			_, err := models.GetStudentAssignment(s.db, assignmentID, userID)
			if err != nil {
				return nil, errors.New("assignment not found or access denied")
			}
		} else {
			return nil, errors.New("access denied")
		}
	}

	return assignment, nil
}

// GetAssignmentsByInstructor retrieves all assignments for an instructor
func (s *AssignmentService) GetAssignmentsByInstructor(instructorID uint) ([]models.Assignment, error) {
	// Validate instructor exists and has instructor role
	var instructor models.User
	if err := s.db.First(&instructor, instructorID).Error; err != nil {
		return nil, errors.New("instructor not found")
	}

	if !instructor.IsInstructor() {
		return nil, errors.New("user is not an instructor")
	}

	assignments, err := models.GetAssignmentsByInstructor(s.db, instructorID)
	if err != nil {
		return nil, err
	}

	return assignments, nil
}

// UpdateAssignmentInput represents input for updating an assignment
type UpdateAssignmentInput struct {
	Title       string
	Description string
	URL         string
	Category    string
	DueDate     *time.Time
}

// UpdateAssignment updates an existing assignment
func (s *AssignmentService) UpdateAssignment(assignmentID uint, instructorID uint, input UpdateAssignmentInput) error {
	// Get assignment and validate ownership
	assignment, err := models.GetAssignmentByID(s.db, assignmentID)
	if err != nil {
		return err
	}

	if assignment.CreatedByID != instructorID {
		return errors.New("access denied")
	}

	// Validate input
	if input.Title == "" {
		return errors.New("title is required")
	}

	if input.URL == "" {
		return errors.New("URL is required")
	}

	// Update assignment
	return assignment.UpdateAssignment(s.db, input.Title, input.Description, input.URL, input.Category, input.DueDate)
}

// DeleteAssignment deletes an assignment
func (s *AssignmentService) DeleteAssignment(assignmentID uint, instructorID uint) error {
	// Get assignment and validate ownership
	assignment, err := models.GetAssignmentByID(s.db, assignmentID)
	if err != nil {
		return err
	}

	if assignment.CreatedByID != instructorID {
		return errors.New("access denied")
	}

	// Delete assignment (soft delete)
	return assignment.DeleteAssignment(s.db)
}

// AssignToStudent assigns an assignment to a student
func (s *AssignmentService) AssignToStudent(assignmentID uint, studentID uint, instructorID uint) error {
	// Validate assignment exists and instructor owns it
	assignment, err := models.GetAssignmentByID(s.db, assignmentID)
	if err != nil {
		return err
	}

	if assignment.CreatedByID != instructorID {
		return errors.New("access denied")
	}

	// Validate student exists and has student role
	var student models.User
	if err := s.db.First(&student, studentID).Error; err != nil {
		return errors.New("student not found")
	}

	if !student.IsStudent() {
		return errors.New("user is not a student")
	}

	// Check if already assigned
	_, err = models.GetStudentAssignment(s.db, assignmentID, studentID)
	if err == nil {
		return errors.New("assignment already assigned to this student")
	}

	// Create student assignment
	_, err = models.CreateStudentAssignment(s.db, assignmentID, studentID)
	return err
}

// AssignToMultipleStudents assigns an assignment to multiple students
func (s *AssignmentService) AssignToMultipleStudents(assignmentID uint, studentIDs []uint, instructorID uint) error {
	// Validate assignment exists and instructor owns it
	assignment, err := models.GetAssignmentByID(s.db, assignmentID)
	if err != nil {
		return err
	}

	if assignment.CreatedByID != instructorID {
		return errors.New("access denied")
	}

	// Validate all students exist and have student role
	var students []models.User
	if err := s.db.Where("id IN ? AND role = ?", studentIDs, "student").Find(&students).Error; err != nil {
		return err
	}

	if len(students) != len(studentIDs) {
		return errors.New("some students not found or not valid students")
	}

	// Filter out already assigned students
	var validStudentIDs []uint
	for _, studentID := range studentIDs {
		_, err := models.GetStudentAssignment(s.db, assignmentID, studentID)
		if err != nil { // Student not assigned yet
			validStudentIDs = append(validStudentIDs, studentID)
		}
	}

	if len(validStudentIDs) == 0 {
		return errors.New("all students are already assigned to this assignment")
	}

	// Bulk create student assignments
	return models.BulkCreateStudentAssignments(s.db, assignmentID, validStudentIDs)
}

// RemoveStudentAssignment removes a student assignment
func (s *AssignmentService) RemoveStudentAssignment(assignmentID uint, studentID uint, instructorID uint) error {
	// Validate assignment exists and instructor owns it
	assignment, err := models.GetAssignmentByID(s.db, assignmentID)
	if err != nil {
		return err
	}

	if assignment.CreatedByID != instructorID {
		return errors.New("access denied")
	}

	// Remove student assignment
	return models.RemoveStudentAssignment(s.db, assignmentID, studentID)
}

// GetAssignmentProgress gets progress statistics for an assignment
func (s *AssignmentService) GetAssignmentProgress(assignmentID uint, instructorID uint) (map[string]int, error) {
	// Validate assignment exists and instructor owns it
	assignment, err := models.GetAssignmentByID(s.db, assignmentID)
	if err != nil {
		return nil, err
	}

	if assignment.CreatedByID != instructorID {
		return nil, errors.New("access denied")
	}

	// Get progress statistics
	return models.GetAssignmentProgress(s.db, assignmentID)
}

// GetAssignmentStudents gets all students assigned to an assignment
func (s *AssignmentService) GetAssignmentStudents(assignmentID uint, instructorID uint) ([]models.StudentAssignment, error) {
	// Validate assignment exists and instructor owns it
	assignment, err := models.GetAssignmentByID(s.db, assignmentID)
	if err != nil {
		return nil, err
	}

	if assignment.CreatedByID != instructorID {
		return nil, errors.New("access denied")
	}

	// Get assigned students
	return models.GetStudentAssignmentsByAssignment(s.db, assignmentID)
}

// SearchAssignments searches assignments by query
func (s *AssignmentService) SearchAssignments(query string, instructorID uint) ([]models.Assignment, error) {
	// Validate instructor exists and has instructor role
	var instructor models.User
	if err := s.db.First(&instructor, instructorID).Error; err != nil {
		return nil, errors.New("instructor not found")
	}

	if !instructor.IsInstructor() {
		return nil, errors.New("user is not an instructor")
	}

	return models.SearchAssignments(s.db, query, instructorID)
}

// GetAssignmentsByCategory gets assignments by category
func (s *AssignmentService) GetAssignmentsByCategory(category string, instructorID uint) ([]models.Assignment, error) {
	// Validate instructor exists and has instructor role
	var instructor models.User
	if err := s.db.First(&instructor, instructorID).Error; err != nil {
		return nil, errors.New("instructor not found")
	}

	if !instructor.IsInstructor() {
		return nil, errors.New("user is not an instructor")
	}

	return models.GetAssignmentsByCategory(s.db, category, instructorID)
}

// GetAllStudents gets all students for assignment purposes
func (s *AssignmentService) GetAllStudents(instructorID uint) ([]models.User, error) {
	// Validate instructor exists and has instructor role
	var instructor models.User
	if err := s.db.First(&instructor, instructorID).Error; err != nil {
		return nil, errors.New("instructor not found")
	}

	if !instructor.IsInstructor() {
		return nil, errors.New("user is not an instructor")
	}

	var students []models.User
	err := s.db.Where("role = ?", "student").Find(&students).Error
	if err != nil {
		return nil, err
	}

	return students, nil
}
