package handlers

import (
	"net/http"

	"zipcodereader/services"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

// AuthHandler handles authentication-related requests
type AuthHandler struct {
	authService *services.AuthService
}

// NewAuthHandler creates a new authentication handler
func NewAuthHandler(authService *services.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

// Login initiates GitHub OAuth2 flow
func (h *AuthHandler) Login(c *gin.Context) {
	session := sessions.Default(c)

	// Generate state token
	state, err := h.authService.GenerateStateToken()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate state token"})
		return
	}

	// Store state in session
	session.Set("oauth_state", state)
	session.Save()

	// Redirect to GitHub OAuth2 URL
	authURL := h.authService.GetAuthURL(state)
	c.Redirect(http.StatusTemporaryRedirect, authURL)
}

// Callback handles GitHub OAuth2 callback
func (h *AuthHandler) Callback(c *gin.Context) {
	session := sessions.Default(c)

	// Verify state parameter
	storedState := session.Get("oauth_state")
	if storedState == nil || storedState != c.Query("state") {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid state parameter"})
		return
	}

	// Clear state from session
	session.Delete("oauth_state")

	// Exchange code for token
	code := c.Query("code")
	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Authorization code not found"})
		return
	}

	token, err := h.authService.ExchangeCodeForToken(c.Request.Context(), code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to exchange code for token"})
		return
	}

	// Get user information from GitHub
	githubUser, err := h.authService.GetGitHubUser(c.Request.Context(), token)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user information"})
		return
	}

	// Create or update user in database
	user, err := h.authService.CreateOrUpdateUser(githubUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create/update user"})
		return
	}

	// Store user ID in session
	session.Set("user_id", user.ID)
	session.Set("user_role", user.Role)
	session.Save()

	// Redirect to dashboard
	c.Redirect(http.StatusTemporaryRedirect, "/dashboard")
}

// Logout clears the user session
func (h *AuthHandler) Logout(c *gin.Context) {
	session := sessions.Default(c)
	session.Clear()
	session.Save()

	c.Redirect(http.StatusTemporaryRedirect, "/")
}

// Dashboard shows the user dashboard
func (h *AuthHandler) Dashboard(c *gin.Context) {
	session := sessions.Default(c)
	userID := session.Get("user_id")

	if userID == nil {
		c.Redirect(http.StatusTemporaryRedirect, "/")
		return
	}

	// Get user information
	user, err := h.authService.ValidateUser(userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user information"})
		return
	}

	c.HTML(http.StatusOK, "dashboard.html", gin.H{
		"title": "Dashboard",
		"user":  user,
	})
}
