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
