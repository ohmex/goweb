package handlers

import (
	"context"
	"fmt"
	"goweb/api"
	"goweb/models"
	"goweb/repositories"
	"goweb/server"
	"goweb/server/requests"
	"goweb/server/responses"
	tokenservice "goweb/services/token"
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/golang-jwt/jwt/v5"
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
// @Success 200 {object} responses.LoginResponse
// @Failure 401 {object} responses.Error
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
	userRepository := repositories.NewUserRepository(authHandler.server.DB)
	userRepository.GetUserByEmail(&user, loginRequest.Email)

	if user.ID == 0 || (bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginRequest.Password)) != nil) {
		return api.WebResponse(c, http.StatusUnauthorized, api.INVALID_CREDENTIALS())
	}

	tokenService := tokenservice.NewTokenService(authHandler.server)

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
// @Success 200 {object} responses.LoginResponse
// @Failure 401 {object} responses.Error
// @Router /refresh [post]
func (authHandler *AuthHandler) RefreshToken(c echo.Context) error {
	refreshRequest := new(requests.RefreshRequest)
	if err := c.Bind(refreshRequest); err != nil {
		return err
	}

	token, err := jwt.Parse(refreshRequest.Token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(authHandler.server.Config.Auth.RefreshSecret), nil
	})

	if err != nil {
		return api.WebResponse(c, http.StatusUnauthorized, api.INVALID_SIGNING_METHOD(err.Error()))
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok && !token.Valid {
		return api.WebResponse(c, http.StatusUnauthorized, api.INVALID_TOKEN())
	}

	user := new(models.User)
	authHandler.server.DB.First(&user, int(claims["id"].(float64)))

	if user.ID == 0 {
		return api.WebResponse(c, http.StatusUnauthorized, api.USER_NOT_FOUND())
	}

	tokenService := tokenservice.NewTokenService(authHandler.server)

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
// @Success 200 {object} responses.Data
// @Failure 401 {object} responses.Data
// @Security ApiKeyAuth
// @Router /logout [post]
func (authHandler *AuthHandler) Logout(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*tokenservice.JwtCustomClaims)

	authHandler.server.Redis.Del(context.Background(), fmt.Sprintf("token-%d", claims.ID))

	return api.WebResponse(c, http.StatusOK, api.USER_LOGGED_OUT())
}
