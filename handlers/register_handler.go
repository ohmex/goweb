package handlers

import (
	"goweb/api"
	"goweb/models"
	"goweb/requests"
	"goweb/server"
	"goweb/services"
	"net/http"

	"github.com/labstack/echo/v4"
)

type RegisterHandler struct {
	server *server.Server
}

func NewRegisterHandler(server *server.Server) *RegisterHandler {
	return &RegisterHandler{server: server}
}

// Register godoc
// @Summary Register
// @Description New user registration
// @ID user-register
// @Tags User Actions
// @Accept json
// @Produce json
// @Param params body requests.RegisterRequest true "User's email, user's password"
// @Success 201 {object} responses.Data
// @Failure 400 {object} responses.Error
// @Router /register [post]
func (registerHandler *RegisterHandler) Register(c echo.Context) error {
	registerRequest := new(requests.RegisterRequest)

	if err := c.Bind(registerRequest); err != nil {
		return err
	}

	if err := registerRequest.Validate(); err != nil {
		return api.WebResponse(c, http.StatusBadRequest, api.FIELD_VALIDATION_ERROR())
	}

	existUser := models.User{}
	userRepository := services.NewUserService(registerHandler.server.DB)
	userRepository.GetUserByEmail(&existUser, registerRequest.Email)

	if existUser.ID != 0 {
		return api.WebResponse(c, http.StatusBadRequest, api.USER_EXISTS())
	}

	userService := services.NewUserService(registerHandler.server.DB)
	if err := userService.Register(registerRequest); err != nil {
		return api.WebResponse(c, http.StatusInternalServerError, api.INTERNAL_SERVICE_ERROR())
	}

	return api.WebResponse(c, http.StatusCreated, api.RESOURCE_CREATED("User successfully created"))
}
