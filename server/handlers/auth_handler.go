package handlers

import (
	"context"
	"fmt"
	"goweb/models"
	"goweb/repositories"
	"goweb/requests"
	"goweb/responses"
	"goweb/server"
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
		return responses.ErrorResponse(c, http.StatusBadRequest, "Required fields are empty or not valid")
	}

	user := models.User{}
	userRepository := repositories.NewUserRepository(authHandler.server.DB)
	userRepository.GetUserByEmail(&user, loginRequest.Email)

	if user.ID == 0 || (bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginRequest.Password)) != nil) {
		return responses.ErrorResponse(c, http.StatusUnauthorized, "Invalid credentials")
	}

	tokenService := tokenservice.NewTokenService(authHandler.server)

	accessToken, refreshToken, exp, _ := tokenService.GenerateTokenPair(&user)

	res := responses.NewLoginResponse(accessToken, refreshToken, exp)

	return responses.Response(c, http.StatusOK, res)
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
		return responses.ErrorResponse(c, http.StatusUnauthorized, err.Error())
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok && !token.Valid {
		return responses.ErrorResponse(c, http.StatusUnauthorized, "Invalid token")
	}

	user := new(models.User)
	authHandler.server.DB.First(&user, int(claims["id"].(float64)))

	if user.ID == 0 {
		return responses.ErrorResponse(c, http.StatusUnauthorized, "User not found")
	}

	tokenService := tokenservice.NewTokenService(authHandler.server)

	accessToken, refreshToken, exp, _ := tokenService.GenerateTokenPair(user)

	res := responses.NewLoginResponse(accessToken, refreshToken, exp)

	return responses.Response(c, http.StatusOK, res)
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

	return responses.MessageResponse(c, http.StatusOK, "User logged out")
}
