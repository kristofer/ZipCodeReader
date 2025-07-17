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
		"title": "Login",
	})
}

// Login handles local login form submission
func (h *LocalAuthHandler) Login(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")

	if username == "" || password == "" {
		c.HTML(http.StatusBadRequest, "local_login.html", gin.H{
			"title": "Login",
			"error": "Username and password are required",
		})
		return
	}

	// Authenticate user
	user, err := models.AuthenticateLocalUser(h.db, username, password)
	if err != nil {
		c.HTML(http.StatusUnauthorized, "local_login.html", gin.H{
			"title": "Login",
			"error": "Invalid credentials",
		})
		return
	}

	// Create session
	session := sessions.Default(c)
	session.Set("user_id", user.ID)
	session.Set("user_role", user.Role)
	session.Save()

	c.Redirect(http.StatusTemporaryRedirect, "/dashboard")
}

// ShowRegister shows the local registration form
func (h *LocalAuthHandler) ShowRegister(c *gin.Context) {
	c.HTML(http.StatusOK, "local_register.html", gin.H{
		"title": "Register",
	})
}

// Register handles local registration form submission
func (h *LocalAuthHandler) Register(c *gin.Context) {
	username := c.PostForm("username")
	email := c.PostForm("email")
	password := c.PostForm("password")
	confirmPassword := c.PostForm("confirm_password")

	// Validation
	if username == "" || email == "" || password == "" {
		c.HTML(http.StatusBadRequest, "local_register.html", gin.H{
			"title": "Register",
			"error": "All fields are required",
		})
		return
	}

	if password != confirmPassword {
		c.HTML(http.StatusBadRequest, "local_register.html", gin.H{
			"title": "Register",
			"error": "Passwords do not match",
		})
		return
	}

	if len(password) < 6 {
		c.HTML(http.StatusBadRequest, "local_register.html", gin.H{
			"title": "Register",
			"error": "Password must be at least 6 characters long",
		})
		return
	}

	// Create user
	user, err := models.CreateLocalUser(h.db, username, email, password)
	if err != nil {
		c.HTML(http.StatusConflict, "local_register.html", gin.H{
			"title": "Register",
			"error": "Username already exists",
		})
		return
	}

	// Create session
	session := sessions.Default(c)
	session.Set("user_id", user.ID)
	session.Set("user_role", user.Role)
	session.Save()

	c.Redirect(http.StatusTemporaryRedirect, "/dashboard")
}

// Logout clears the user session
func (h *LocalAuthHandler) Logout(c *gin.Context) {
	session := sessions.Default(c)
	session.Clear()
	session.Save()

	c.Redirect(http.StatusTemporaryRedirect, "/")
}
