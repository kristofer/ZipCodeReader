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
