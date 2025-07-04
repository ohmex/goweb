package handlers

import (
	"fmt"
	"goweb/api"
	"goweb/models"
	"goweb/requests"
	"goweb/responses"
	"goweb/server"
	"goweb/services"
	"goweb/util"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/bcrypt"
)

// AuthHandler provides endpoints for user authentication, token refresh, and logout operations.
type AuthHandler struct {
	server       *server.Server
	userService  *services.UserService
	tokenService *services.TokenService
}

// NewAuthHandler initializes the AuthHandler with the provided server and its dependencies.
func NewAuthHandler(server *server.Server) *AuthHandler {
	return &AuthHandler{
		server:       server,
		userService:  services.NewUserService(server.DB),
		tokenService: services.NewTokenService(server),
	}
}

// Helper to generate token pair and return response
func (h *AuthHandler) respondWithTokenPair(c echo.Context, user *models.User) error {
	accessToken, refreshToken, exp, err := h.tokenService.GenerateTokenPair(user)
	if err != nil {
		log.Error().Str("event", "token_generation_failed").Err(err).Uint64("user_id", uint64(user.ID)).Msg("Failed to generate authentication tokens")
		return api.WebResponse(c, http.StatusInternalServerError, api.INTERNAL_SERVICE_ERROR("Failed to generate authentication tokens"))
	}
	res := responses.NewLoginResponse(accessToken, refreshToken, exp)
	return api.WebResponse(c, http.StatusOK, res)
}

// Login godoc
// @Summary Authenticates a user using email and password, and returns a token pair if successful.
// @Description Perform user login with email and password
// @ID user-login
// @Tags User Actions
// @Accept json
// @Produce json
// @Param params body requests.LoginRequest true "User's credentials"
// @Success 200 {object} responses.LoginResponse
// @Failure 400 {object} api.Response
// @Failure 401 {object} api.Response
// @Failure 500 {object} api.Response
// @Router /login [post]
func (h *AuthHandler) Login(c echo.Context) error {
	start := time.Now()

	loginRequest, err := util.BindAndValidate[requests.LoginRequest](c)
	if err != nil {
		return api.WebResponse(c, http.StatusBadRequest, err)
	}

	user := &models.User{}
	if err := h.userService.GetUserByEmail(user, loginRequest.Email); err != nil || user.ID == 0 {
		log.Info().Str("event", "login_failed").Str("email", loginRequest.Email).Str("error", "user_not_found").Msg("Login failed: user not found")
		return api.WebResponse(c, http.StatusUnauthorized, api.INVALID_CREDENTIALS())
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginRequest.Password)); err != nil {
		log.Info().Str("event", "login_failed").Str("email", loginRequest.Email).Uint64("user_id", uint64(user.ID)).Str("error", "invalid_password").Msg("Login failed: invalid password")
		return api.WebResponse(c, http.StatusUnauthorized, api.INVALID_CREDENTIALS())
	}

	log.Info().Str("event", "login_success").Uint64("user_id", uint64(user.ID)).Str("email", user.Email).Int64("duration_ms", time.Since(start).Milliseconds()).Msg("Login successful")
	return h.respondWithTokenPair(c, user)
}

// RefreshToken godoc
// @Summary Refresh access token
// @Description RefreshToken issues a new access token using a valid refresh token.
// @ID user-refresh
// @Tags User Actions
// @Accept json
// @Produce json
// @Param params body requests.RefreshRequest true "Refresh token"
// @Success 200 {object} responses.LoginResponse
// @Failure 400 {object} api.Response
// @Failure 401 {object} api.Response
// @Failure 500 {object} api.Response
// @Router /refresh [post]
func (h *AuthHandler) RefreshToken(c echo.Context) error {
	start := time.Now()

	refreshRequest, err := util.BindAndValidate[requests.RefreshRequest](c)
	if err != nil {
		return api.WebResponse(c, http.StatusBadRequest, err)
	}

	claims, err := h.tokenService.ParseToken(refreshRequest.Token, h.server.Config.Auth.RefreshSecret)
	if err != nil {
		log.Info().Str("event", "refresh_token_invalid").Str("error", err.Error()).Msg("Invalid or expired refresh token")
		return api.WebResponse(c, http.StatusUnauthorized, api.INVALID_TOKEN("Invalid or expired refresh token"))
	}

	user, err := h.tokenService.ValidateToken(claims, true)
	if err != nil || user.ID == 0 {
		log.Info().Str("event", "refresh_token_validation_failed").Uint64("user_id", claims.UserID).Str("error", err.Error()).Msg("Refresh token validation failed")
		return api.WebResponse(c, http.StatusUnauthorized, api.INVALID_TOKEN(err.Error()))
	}

	log.Info().Str("event", "token_refresh_success").Uint64("user_id", uint64(user.ID)).Int64("duration_ms", time.Since(start).Milliseconds()).Msg("Token refresh successful")
	return h.respondWithTokenPair(c, user)
}

// Logout godoc
// @Summary Logout user
// @Description Invalidate user session and logout
// @ID user-logout
// @Tags User Actions
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} api.Response
// @Failure 401 {object} api.Response
// @Failure 500 {object} api.Response
// @Router /logout [post]
func (h *AuthHandler) Logout(c echo.Context) error {
	tokenInterface := c.Get("token")
	token, ok := tokenInterface.(*jwt.Token)
	if !ok {
		return api.WebResponse(c, http.StatusUnauthorized, api.INVALID_TOKEN("Invalid authentication token"))
	}

	claims, ok := token.Claims.(*services.JwtCustomClaims)
	if !ok {
		return api.WebResponse(c, http.StatusUnauthorized, api.INVALID_TOKEN("Invalid token claims"))
	}

	err := h.server.Redis.Del(c.Request().Context(), fmt.Sprintf("token-%d", claims.UserID)).Err()
	if err != nil {
		log.Error().Str("event", "logout_redis_failed").Err(err).Uint64("user_id", claims.UserID).Msg("Failed to logout user (redis error)")
		return api.WebResponse(c, http.StatusInternalServerError, api.INTERNAL_SERVICE_ERROR("Failed to logout user"))
	}

	log.Info().Str("event", "logout_success").Uint64("user_id", claims.UserID).Msg("Logout successful")
	return api.WebResponse(c, http.StatusOK, api.USER_LOGGED_OUT())
}
