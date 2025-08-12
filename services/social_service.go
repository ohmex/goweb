package services

import (
	"context"
	"encoding/json"
	"fmt"
	"goweb/config"
	"goweb/models"
	"goweb/server"
	"io"

	"github.com/rs/zerolog/log"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"golang.org/x/oauth2/github"
	"gorm.io/gorm"
)

type SocialService struct {
	server *server.Server
	db     *gorm.DB
	config *config.Config
}

type GoogleUserInfo struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	Name          string `json:"name"`
	Picture       string `json:"picture"`
	VerifiedEmail bool   `json:"verified_email"`
}

type GitHubUserInfo struct {
	ID    int    `json:"id"`
	Login string `json:"login"`
	Email string `json:"email"`
	Name  string `json:"name"`
	AvatarURL string `json:"avatar_url"`
}

func NewSocialService(server *server.Server) *SocialService {
	return &SocialService{
		server: server,
		db:     server.DB,
		config: server.Config,
	}
}

// GetGoogleOAuthConfig returns the OAuth2 configuration for Google
func (s *SocialService) GetGoogleOAuthConfig() *oauth2.Config {
	return &oauth2.Config{
		ClientID:     s.config.Auth.Social.Google.ClientID,
		ClientSecret: s.config.Auth.Social.Google.ClientSecret,
		RedirectURL:  s.config.Auth.Social.Google.RedirectURL,
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: google.Endpoint,
	}
}

// GetGitHubOAuthConfig returns the OAuth2 configuration for GitHub
func (s *SocialService) GetGitHubOAuthConfig() *oauth2.Config {
	return &oauth2.Config{
		ClientID:     s.config.Auth.Social.GitHub.ClientID,
		ClientSecret: s.config.Auth.Social.GitHub.ClientSecret,
		RedirectURL:  s.config.Auth.Social.GitHub.RedirectURL,
		Scopes: []string{
			"user:email",
			"read:user",
		},
		Endpoint: github.Endpoint,
	}
}

// GetGoogleUserInfo fetches user information from Google OAuth
func (s *SocialService) GetGoogleUserInfo(token *oauth2.Token) (*GoogleUserInfo, error) {
	client := s.GetGoogleOAuthConfig().Client(context.Background(), token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		return nil, fmt.Errorf("failed to get user info: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var userInfo GoogleUserInfo
	if err := json.Unmarshal(body, &userInfo); err != nil {
		return nil, fmt.Errorf("failed to unmarshal user info: %w", err)
	}

	return &userInfo, nil
}

// GetGitHubUserInfo fetches user information from GitHub OAuth
func (s *SocialService) GetGitHubUserInfo(token *oauth2.Token) (*GitHubUserInfo, error) {
	client := s.GetGitHubOAuthConfig().Client(context.Background(), token)
	resp, err := client.Get("https://api.github.com/user")
	if err != nil {
		return nil, fmt.Errorf("failed to get user info: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var userInfo GitHubUserInfo
	if err := json.Unmarshal(body, &userInfo); err != nil {
		return nil, fmt.Errorf("failed to unmarshal user info: %w", err)
	}

	// Get email if not provided in user info
	if userInfo.Email == "" {
		emailResp, err := client.Get("https://api.github.com/user/emails")
		if err == nil {
			defer emailResp.Body.Close()
			var emails []struct {
				Email   string `json:"email"`
				Primary bool   `json:"primary"`
			}
			if emailBody, err := io.ReadAll(emailResp.Body); err == nil {
				if json.Unmarshal(emailBody, &emails) == nil {
					for _, email := range emails {
						if email.Primary {
							userInfo.Email = email.Email
							break
						}
					}
				}
			}
		}
	}

	return &userInfo, nil
}

// FindOrCreateUser finds an existing user or creates a new one from social login
func (s *SocialService) FindOrCreateUser(provider string, providerID string, email string, name string, avatar string, isVerified bool) (*models.User, error) {
	var user models.User

	// First, try to find by provider ID
	err := s.db.Where("provider = ? AND provider_id = ?", provider, providerID).First(&user).Error
	if err == nil {
		// User found, update information
		user.Name = name
		user.Email = email
		user.Avatar = avatar
		user.IsVerified = isVerified
		s.db.Save(&user)
		return &user, nil
	}

	// If not found by provider ID, try to find by email
	err = s.db.Where("email = ?", email).First(&user).Error
	if err == nil {
		// User exists with this email, link the social account
		user.Provider = provider
		user.ProviderID = providerID
		user.Avatar = avatar
		user.IsVerified = isVerified
		s.db.Save(&user)
		return &user, nil
	}

	// Create new user
	user = models.User{
		Email:        email,
		Name:         name,
		Avatar:       avatar,
		Provider:     provider,
		ProviderID:   providerID,
		IsVerified:   isVerified,
		Password:     "", // No password for social login users
	}

	if err := s.db.Create(&user).Error; err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	log.Info().Str("event", "social_user_created").Str("provider", provider).Str("email", email).Uint64("user_id", uint64(user.ID)).Msg("New user created via social login")
	return &user, nil
}

// HandleGoogleLogin processes Google OAuth login
func (s *SocialService) HandleGoogleLogin(code string) (*models.User, error) {
	config := s.GetGoogleOAuthConfig()
	
	// Exchange code for token
	token, err := config.Exchange(context.Background(), code)
	if err != nil {
		return nil, fmt.Errorf("failed to exchange code for token: %w", err)
	}

	// Get user info
	userInfo, err := s.GetGoogleUserInfo(token)
	if err != nil {
		return nil, fmt.Errorf("failed to get Google user info: %w", err)
	}

	// Find or create user
	user, err := s.FindOrCreateUser(
		"google",
		userInfo.ID,
		userInfo.Email,
		userInfo.Name,
		userInfo.Picture,
		userInfo.VerifiedEmail,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to find or create user: %w", err)
	}

	return user, nil
}

// HandleGitHubLogin processes GitHub OAuth login
func (s *SocialService) HandleGitHubLogin(code string) (*models.User, error) {
	config := s.GetGitHubOAuthConfig()
	
	// Exchange code for token
	token, err := config.Exchange(context.Background(), code)
	if err != nil {
		return nil, fmt.Errorf("failed to exchange code for token: %w", err)
	}

	// Get user info
	userInfo, err := s.GetGitHubUserInfo(token)
	if err != nil {
		return nil, fmt.Errorf("failed to get GitHub user info: %w", err)
	}

	// Find or create user
	user, err := s.FindOrCreateUser(
		"github",
		fmt.Sprintf("%d", userInfo.ID),
		userInfo.Email,
		userInfo.Name,
		userInfo.AvatarURL,
		true, // GitHub emails are typically verified
	)
	if err != nil {
		return nil, fmt.Errorf("failed to find or create user: %w", err)
	}

	return user, nil
}
