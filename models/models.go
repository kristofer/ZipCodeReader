package models

import (
	"time"

	"gorm.io/gorm"
)

// BaseModel contains common fields for all models
type BaseModel struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}

// User model will be implemented in Phase 2
// Assignment model will be implemented in Phase 3
// Group model will be implemented in Phase 5

// This file serves as a placeholder for the models package
// Actual models will be added in their respective phases
