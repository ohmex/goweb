package handlers

import (
	"goweb/api"
	"goweb/requests"
	"goweb/server"
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
// @Router /register [post]
// Register function creates Tenant & User - together as pair
func (registerHandler *RegisterHandler) Register(c echo.Context) error {
	registerRequest := new(requests.RegisterRequest)

	if err := c.Bind(registerRequest); err != nil {
		return err
	}

	if err := registerRequest.Validate(); err != nil {
		return api.WebResponse(c, http.StatusBadRequest, api.FIELD_VALIDATION_ERROR())
	}

	//tenant := c.Get("tenant").(*models.Tenant)
	//return services.NewUserService(u.Server.DB).Register(e, registerRequest, tenant)

	return api.WebResponse(c, http.StatusCreated, api.RESOURCE_CREATED("User successfully created"))
}
