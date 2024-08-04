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

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	server *server.Server
}

func NewAuthHandler(server *server.Server) *AuthHandler {
	return &AuthHandler{server: server}
}

// Login godoc
// @Summary Authenticate a user
// @Description Perform user login
// @ID user-login
// @Tags User Actions
// @Accept json
// @Produce json
// @Param params body requests.LoginRequest true "User's credentials"
// @Router /login [post]
func (authHandler *AuthHandler) Login(c echo.Context) error {
	loginRequest := new(requests.LoginRequest)

	if err := c.Bind(loginRequest); err != nil {
		return err
	}

	if err := loginRequest.Validate(); err != nil {
		return api.WebResponse(c, http.StatusBadRequest, api.FIELD_VALIDATION_ERROR())
	}

	user := models.User{}
	userRepository := services.NewUserService(authHandler.server.DB)
	userRepository.GetUserByEmail(&user, loginRequest.Email)

	if user.ID == 0 || (bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginRequest.Password)) != nil) {
		return api.WebResponse(c, http.StatusUnauthorized, api.INVALID_CREDENTIALS())
	}

	tokenService := services.NewTokenService(authHandler.server)

	accessToken, refreshToken, exp, _ := tokenService.GenerateTokenPair(&user)

	res := responses.NewLoginResponse(accessToken, refreshToken, exp)

	return api.WebResponse(c, http.StatusOK, res)
}

// Refresh godoc
// @Summary Refresh access token
// @Description Perform refresh access token
// @ID user-refresh
// @Tags User Actions
// @Accept json
// @Produce json
// @Param params body requests.RefreshRequest true "Refresh token"
// @Router /refresh [post]
func (authHandler *AuthHandler) RefreshToken(c echo.Context) error {
	refreshRequest := new(requests.RefreshRequest)
	if err := c.Bind(refreshRequest); err != nil {
		return err
	}

	tokenService := services.NewTokenService(authHandler.server)
	claims, _ := tokenService.ParseToken(refreshRequest.Token, authHandler.server.Config.Auth.RefreshSecret)
	user, err := services.NewTokenService(authHandler.server).ValidateToken(claims, true)
	if err != nil {
		return api.WebResponse(c, http.StatusUnauthorized, err)
	}

	if user.ID == 0 {
		return api.WebResponse(c, http.StatusUnauthorized, api.USER_NOT_FOUND())
	}

	accessToken, refreshToken, exp, _ := tokenService.GenerateTokenPair(user)

	res := responses.NewLoginResponse(accessToken, refreshToken, exp)

	return api.WebResponse(c, http.StatusOK, res)
}

// Logout godoc
// @Summary Logout
// @Description Perform the user's logout
// @ID user-logout
// @Tags User Actions
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Router /logout [post]
func (authHandler *AuthHandler) Logout(c echo.Context) error {
	token := c.Get("token").(*jwt.Token)
	claims := token.Claims.(*services.JwtCustomClaims)

	authHandler.server.Redis.Del(context.Background(), fmt.Sprintf("token-%d", claims.ID))

	return api.WebResponse(c, http.StatusOK, api.USER_LOGGED_OUT())
}
