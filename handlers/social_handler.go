package handlers

import (
	"goweb/api"
	"goweb/models"
	"goweb/responses"
	"goweb/server"
	"goweb/services"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"golang.org/x/oauth2"
)

// SocialHandler provides endpoints for social login operations
type SocialHandler struct {
	server        *server.Server
	socialService *services.SocialService
	tokenService  *services.TokenService
}

// NewSocialHandler initializes the SocialHandler with the provided server and its dependencies
func NewSocialHandler(server *server.Server) *SocialHandler {
	return &SocialHandler{
		server:        server,
		socialService: services.NewSocialService(server),
		tokenService:  services.NewTokenService(server),
	}
}

// Helper to generate token pair and return response
func (h *SocialHandler) respondWithTokenPair(c echo.Context, user *models.User) error {
	accessToken, refreshToken, exp, err := h.tokenService.GenerateTokenPair(user)
	if err != nil {
		log.Error().Str("event", "token_generation_failed").Err(err).Uint64("user_id", uint64(user.ID)).Msg("Failed to generate authentication tokens")
		return api.WebResponse(c, http.StatusInternalServerError, api.INTERNAL_SERVICE_ERROR("Failed to generate authentication tokens"))
	}
	res := responses.NewLoginResponse(accessToken, refreshToken, exp)
	return api.WebResponse(c, http.StatusOK, res)
}

// GoogleLogin godoc
// @Summary Initiate Google OAuth login
// @Description Redirects user to Google OAuth for authentication
// @ID google-login
// @Tags Social Login
// @Accept json
// @Produce json
// @Success 200 {object} api.Response
// @Failure 400 {object} api.Response
// @Failure 500 {object} api.Response
// @Router /auth/google [get]
func (h *SocialHandler) GoogleLogin(c echo.Context) error {
	config := h.socialService.GetGoogleOAuthConfig()
	url := config.AuthCodeURL("state", oauth2.AccessTypeOffline)
	return c.Redirect(http.StatusTemporaryRedirect, url)
}

// GoogleCallback godoc
// @Summary Handle Google OAuth callback
// @Description Processes the OAuth callback from Google and authenticates the user
// @ID google-callback
// @Tags Social Login
// @Accept json
// @Produce json
// @Param code query string true "Authorization code from Google"
// @Param state query string false "State parameter"
// @Success 200 {object} responses.LoginResponse
// @Failure 400 {object} api.Response
// @Failure 401 {object} api.Response
// @Failure 500 {object} api.Response
// @Router /auth/google/callback [get]
func (h *SocialHandler) GoogleCallback(c echo.Context) error {
	start := time.Now()
	code := c.QueryParam("code")
	state := c.QueryParam("state")

	if code == "" {
		log.Info().Str("event", "google_callback_failed").Str("error", "missing_code").Msg("Google callback failed: missing authorization code")
		return api.WebResponse(c, http.StatusBadRequest, api.FIELD_VALIDATION_ERROR("Missing authorization code"))
	}

	// TODO: Validate state parameter for security
	_ = state

	user, err := h.socialService.HandleGoogleLogin(code)
	if err != nil {
		log.Error().Str("event", "google_login_failed").Err(err).Msg("Google login failed")
		return api.WebResponse(c, http.StatusInternalServerError, api.INTERNAL_SERVICE_ERROR("Google login failed"))
	}

	log.Info().Str("event", "google_login_success").Uint64("user_id", uint64(user.ID)).Str("email", user.Email).Int64("duration_ms", time.Since(start).Milliseconds()).Msg("Google login successful")
	return h.respondWithTokenPair(c, user)
}

// GitHubLogin godoc
// @Summary Initiate GitHub OAuth login
// @Description Redirects user to GitHub OAuth for authentication
// @ID github-login
// @Tags Social Login
// @Accept json
// @Produce json
// @Success 200 {object} api.Response
// @Failure 400 {object} api.Response
// @Failure 500 {object} api.Response
// @Router /auth/github [get]
func (h *SocialHandler) GitHubLogin(c echo.Context) error {
	config := h.socialService.GetGitHubOAuthConfig()
	url := config.AuthCodeURL("state")
	return c.Redirect(http.StatusTemporaryRedirect, url)
}

// GitHubCallback godoc
// @Summary Handle GitHub OAuth callback
// @Description Processes the OAuth callback from GitHub and authenticates the user
// @ID github-callback
// @Tags Social Login
// @Accept json
// @Produce json
// @Param code query string true "Authorization code from GitHub"
// @Param state query string false "State parameter"
// @Success 200 {object} responses.LoginResponse
// @Failure 400 {object} api.Response
// @Failure 401 {object} api.Response
// @Failure 500 {object} api.Response
// @Router /auth/github/callback [get]
func (h *SocialHandler) GitHubCallback(c echo.Context) error {
	start := time.Now()
	code := c.QueryParam("code")
	state := c.QueryParam("state")

	if code == "" {
		log.Info().Str("event", "github_callback_failed").Str("error", "missing_code").Msg("GitHub callback failed: missing authorization code")
		return api.WebResponse(c, http.StatusBadRequest, api.FIELD_VALIDATION_ERROR("Missing authorization code"))
	}

	// TODO: Validate state parameter for security
	_ = state

	user, err := h.socialService.HandleGitHubLogin(code)
	if err != nil {
		log.Error().Str("event", "github_login_failed").Err(err).Msg("GitHub login failed")
		return api.WebResponse(c, http.StatusInternalServerError, api.INTERNAL_SERVICE_ERROR("GitHub login failed"))
	}

	log.Info().Str("event", "github_login_success").Uint64("user_id", uint64(user.ID)).Str("email", user.Email).Int64("duration_ms", time.Since(start).Milliseconds()).Msg("GitHub login successful")
	return h.respondWithTokenPair(c, user)
}
