package services

import (
	"context"
	"crypto/rand"
	"encoding/base64"

	"zipcodereader/config"
	"zipcodereader/models"

	"github.com/google/go-github/v45/github"
	"golang.org/x/oauth2"
	githuboauth "golang.org/x/oauth2/github"
	"gorm.io/gorm"
)

// AuthService handles authentication operations
type AuthService struct {
	db          *gorm.DB
	oauthConfig *oauth2.Config
}

// NewAuthService creates a new authentication service
func NewAuthService(db *gorm.DB, cfg *config.Config) *AuthService {
	oauthConfig := &oauth2.Config{
		ClientID:     cfg.GitHubClientID,
		ClientSecret: cfg.GitHubClientSecret,
		RedirectURL:  cfg.BaseURL + "/auth/callback",
		Scopes:       []string{"user:email"},
		Endpoint:     githuboauth.Endpoint,
	}

	return &AuthService{
		db:          db,
		oauthConfig: oauthConfig,
	}
}

// GenerateStateToken generates a random state token for OAuth2
func (s *AuthService) GenerateStateToken() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(b), nil
}

// GetAuthURL returns the GitHub OAuth2 authorization URL
func (s *AuthService) GetAuthURL(state string) string {
	return s.oauthConfig.AuthCodeURL(state, oauth2.AccessTypeOnline)
}

// ExchangeCodeForToken exchanges authorization code for access token
func (s *AuthService) ExchangeCodeForToken(ctx context.Context, code string) (*oauth2.Token, error) {
	return s.oauthConfig.Exchange(ctx, code)
}

// GetGitHubUser retrieves user information from GitHub API
func (s *AuthService) GetGitHubUser(ctx context.Context, token *oauth2.Token) (*github.User, error) {
	client := github.NewClient(s.oauthConfig.Client(ctx, token))
	user, _, err := client.Users.Get(ctx, "")
	return user, err
}

// CreateOrUpdateUser creates a new user or updates existing user from GitHub data
func (s *AuthService) CreateOrUpdateUser(githubUser *github.User) (*models.User, error) {
	// Check if user already exists
	existingUser, err := models.GetUserByGitHubID(s.db, githubUser.GetID())
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	// If user exists, update their information
	if existingUser != nil {
		existingUser.Username = githubUser.GetLogin()
		existingUser.Email = githubUser.GetEmail()
		existingUser.AvatarURL = githubUser.GetAvatarURL()

		err = existingUser.Update(s.db)
		if err != nil {
			return nil, err
		}
		return existingUser, nil
	}

	// Create new user
	user, err := models.CreateUser(
		s.db,
		githubUser.GetID(),
		githubUser.GetLogin(),
		githubUser.GetEmail(),
		githubUser.GetAvatarURL(),
	)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// ValidateUser ensures user exists and returns user information
func (s *AuthService) ValidateUser(userID uint) (*models.User, error) {
	return models.GetUserByID(s.db, userID)
}
