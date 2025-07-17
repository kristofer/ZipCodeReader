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

	// Future models will be added here in subsequent phases
	return nil
}
