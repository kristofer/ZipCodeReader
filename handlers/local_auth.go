package handlers

import (
	"net/http"

	"zipcodereader/models"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// LocalAuthHandler handles local authentication requests
type LocalAuthHandler struct {
	db *gorm.DB
}

// NewLocalAuthHandler creates a new local authentication handler
func NewLocalAuthHandler(db *gorm.DB) *LocalAuthHandler {
	return &LocalAuthHandler{
		db: db,
	}
}

// ShowLogin shows the local login form
func (h *LocalAuthHandler) ShowLogin(c *gin.Context) {
	c.HTML(http.StatusOK, "local_login.html", gin.H{
		"title":          "Login",
		"use_local_auth": true,
	})
}

// Login handles local login form submission
func (h *LocalAuthHandler) Login(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")

	if username == "" || password == "" {
		c.HTML(http.StatusBadRequest, "local_login.html", gin.H{
			"title":          "Login",
			"error":          "Username and password are required",
			"use_local_auth": true,
		})
		return
	}

	// Authenticate user
	user, err := models.AuthenticateLocalUser(h.db, username, password)
	if err != nil {
		c.HTML(http.StatusUnauthorized, "local_login.html", gin.H{
			"title":          "Login",
			"error":          "Invalid credentials",
			"use_local_auth": true,
		})
		return
	}

	// Create session
	session := sessions.Default(c)
	session.Set("user_id", user.ID)
	session.Set("user_role", user.Role)
	session.Save()

	c.Redirect(http.StatusSeeOther, "/dashboard")
}

// ShowRegister shows the local registration form
func (h *LocalAuthHandler) ShowRegister(c *gin.Context) {
	c.HTML(http.StatusOK, "local_register.html", gin.H{
		"title":          "Register",
		"use_local_auth": true,
	})
}

// Register handles local registration form submission
func (h *LocalAuthHandler) Register(c *gin.Context) {
	username := c.PostForm("username")
	email := c.PostForm("email")
	password := c.PostForm("password")
	confirmPassword := c.PostForm("confirm_password")
	role := c.PostForm("role")

	// Validation
	if username == "" || email == "" || password == "" {
		c.HTML(http.StatusBadRequest, "local_register.html", gin.H{
			"title":          "Register",
			"error":          "All fields are required",
			"use_local_auth": true,
		})
		return
	}

	if password != confirmPassword {
		c.HTML(http.StatusBadRequest, "local_register.html", gin.H{
			"title":          "Register",
			"error":          "Passwords do not match",
			"use_local_auth": true,
		})
		return
	}

	if len(password) < 6 {
		c.HTML(http.StatusBadRequest, "local_register.html", gin.H{
			"title":          "Register",
			"error":          "Password must be at least 6 characters long",
			"use_local_auth": true,
		})
		return
	}

	// Create user with specified role
	user, err := models.CreateLocalUserWithRole(h.db, username, email, password, role)
	if err != nil {
		c.HTML(http.StatusConflict, "local_register.html", gin.H{
			"title":          "Register",
			"error":          "Username already exists",
			"use_local_auth": true,
		})
		return
	}

	// Create session
	session := sessions.Default(c)
	session.Set("user_id", user.ID)
	session.Set("user_role", user.Role)
	session.Save()

	c.Redirect(http.StatusSeeOther, "/dashboard")
}

// Logout clears the user session
func (h *LocalAuthHandler) Logout(c *gin.Context) {
	session := sessions.Default(c)
	session.Clear()
	session.Save()

	c.Redirect(http.StatusSeeOther, "/")
}
