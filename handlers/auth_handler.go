package handlers

import (
	"context"
	"fmt"
	"goweb/api"
	"goweb/models"
	"goweb/requests"
	"goweb/responses"
	"goweb/server"
	"goweb/services"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	server       *server.Server
	userService  *services.UserService
	tokenService *services.TokenService
}

func NewAuthHandler(server *server.Server) *AuthHandler {
	return &AuthHandler{
		server:       server,
		userService:  services.NewUserService(server.DB),
		tokenService: services.NewTokenService(server),
	}
}

// Login godoc
// @Summary Authenticate a user
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

	// Parse and validate request
	loginRequest := new(requests.LoginRequest)
	if err := c.Bind(loginRequest); err != nil {
		return api.WebResponse(c, http.StatusBadRequest, api.FIELD_VALIDATION_ERROR("Invalid request format"))
	}

	if err := loginRequest.Validate(); err != nil {
		return api.WebResponse(c, http.StatusBadRequest, api.FIELD_VALIDATION_ERROR(err.Error()))
	}

	// Get user by email
	user := &models.User{}
	if err := h.userService.GetUserByEmail(user, loginRequest.Email); err != nil {
		// Log failed login attempt for security monitoring
		h.logSecurityEvent("login_failed", map[string]interface{}{
			"email": loginRequest.Email,
			"error": "user_not_found",
		})
		return api.WebResponse(c, http.StatusUnauthorized, api.INVALID_CREDENTIALS())
	}

	if user.ID == 0 {
		return api.WebResponse(c, http.StatusUnauthorized, api.INVALID_CREDENTIALS())
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginRequest.Password)); err != nil {
		// Log failed login attempt for security monitoring
		h.logSecurityEvent("login_failed", map[string]interface{}{
			"email":   loginRequest.Email,
			"user_id": user.ID,
			"error":   "invalid_password",
		})
		return api.WebResponse(c, http.StatusUnauthorized, api.INVALID_CREDENTIALS())
	}

	// Generate token pair
	accessToken, refreshToken, exp, err := h.tokenService.GenerateTokenPair(user)
	if err != nil {
		h.logError("token_generation_failed", err, map[string]interface{}{
			"user_id": user.ID,
		})
		return api.WebResponse(c, http.StatusInternalServerError, api.INTERNAL_SERVICE_ERROR("Failed to generate authentication tokens"))
	}

	// Log successful login
	h.logSecurityEvent("login_success", map[string]interface{}{
		"user_id":     user.ID,
		"email":       user.Email,
		"duration_ms": time.Since(start).Milliseconds(),
	})

	res := responses.NewLoginResponse(accessToken, refreshToken, exp)
	return api.WebResponse(c, http.StatusOK, res)
}

// RefreshToken godoc
// @Summary Refresh access token
// @Description Refresh access token using a valid refresh token
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

	// Parse and validate request
	refreshRequest := new(requests.RefreshRequest)
	if err := c.Bind(refreshRequest); err != nil {
		return api.WebResponse(c, http.StatusBadRequest, api.FIELD_VALIDATION_ERROR("Invalid request format"))
	}

	if err := refreshRequest.Validate(); err != nil {
		return api.WebResponse(c, http.StatusBadRequest, api.FIELD_VALIDATION_ERROR(err.Error()))
	}

	// Parse and validate refresh token
	claims, err := h.tokenService.ParseToken(refreshRequest.Token, h.server.Config.Auth.RefreshSecret)
	if err != nil {
		h.logSecurityEvent("refresh_token_invalid", map[string]interface{}{
			"error": err.Error(),
		})
		return api.WebResponse(c, http.StatusUnauthorized, api.INVALID_TOKEN("Invalid or expired refresh token"))
	}

	// Validate token and get user
	user, err := h.tokenService.ValidateToken(claims, true)
	if err != nil {
		h.logSecurityEvent("refresh_token_validation_failed", map[string]interface{}{
			"user_id": claims.UserID,
			"error":   err.Error(),
		})
		return api.WebResponse(c, http.StatusUnauthorized, api.INVALID_TOKEN(err.Error()))
	}

	if user.ID == 0 {
		return api.WebResponse(c, http.StatusUnauthorized, api.USER_NOT_FOUND())
	}

	// Generate new token pair
	accessToken, refreshToken, exp, err := h.tokenService.GenerateTokenPair(user)
	if err != nil {
		h.logError("refresh_token_generation_failed", err, map[string]interface{}{
			"user_id": user.ID,
		})
		return api.WebResponse(c, http.StatusInternalServerError, api.INTERNAL_SERVICE_ERROR("Failed to generate new authentication tokens"))
	}

	// Log successful token refresh
	h.logSecurityEvent("token_refresh_success", map[string]interface{}{
		"user_id":     user.ID,
		"duration_ms": time.Since(start).Milliseconds(),
	})

	res := responses.NewLoginResponse(accessToken, refreshToken, exp)
	return api.WebResponse(c, http.StatusOK, res)
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
	// Extract token from context
	tokenInterface := c.Get("token")
	token, ok := tokenInterface.(*jwt.Token)
	if !ok {
		return api.WebResponse(c, http.StatusUnauthorized, api.INVALID_TOKEN("Invalid authentication token"))
	}

	claims, ok := token.Claims.(*services.JwtCustomClaims)
	if !ok {
		return api.WebResponse(c, http.StatusUnauthorized, api.INVALID_TOKEN("Invalid token claims"))
	}

	// Invalidate token in Redis
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := h.server.Redis.Del(ctx, fmt.Sprintf("token-%d", claims.UserID)).Err()
	if err != nil {
		h.logError("logout_redis_failed", err, map[string]interface{}{
			"user_id": claims.UserID,
		})
		return api.WebResponse(c, http.StatusInternalServerError, api.INTERNAL_SERVICE_ERROR("Failed to logout user"))
	}

	// Log successful logout
	h.logSecurityEvent("logout_success", map[string]interface{}{
		"user_id": claims.UserID,
	})

	return api.WebResponse(c, http.StatusOK, api.USER_LOGGED_OUT())
}

// Helper methods for logging and monitoring
func (h *AuthHandler) logSecurityEvent(event string, fields map[string]interface{}) {
	// In a production environment, you would use a proper logging library
	// like logrus, zap, or structured logging
	// For now, we'll use a simple approach
	fmt.Printf("[SECURITY] %s: %+v\n", event, fields)
}

func (h *AuthHandler) logError(event string, err error, fields map[string]interface{}) {
	// In a production environment, you would use a proper logging library
	fmt.Printf("[ERROR] %s: %v, fields: %+v\n", event, err, fields)
}
