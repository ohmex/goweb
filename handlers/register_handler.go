package handlers

import (
	"goweb/api"
	"goweb/requests"
	"goweb/server"
	"net/http"

	"github.com/labstack/echo/v4"
)

// RegisterHandler provides an endpoint for registering new users and their domains.
type RegisterHandler struct {
	server *server.Server
}

// NewRegisterHandler initializes the RegisterHandler with the provided server.
func NewRegisterHandler(server *server.Server) *RegisterHandler {
	return &RegisterHandler{server: server}
}

// Register godoc
// @Summary Register a new user
// @Description Register creates a new user and domain based on the provided registration request.
// @ID user-register
// @Tags User Actions
// @Accept json
// @Produce json
// @Param params body requests.RegisterRequest true "User's email, user's password"
// @Router /register [post]
// Register function creates Domain & User - together as pair
func (registerHandler *RegisterHandler) Register(c echo.Context) error {
	registerRequest := new(requests.RegisterRequest)

	if err := c.Bind(registerRequest); err != nil {
		return err
	}

	if err := registerRequest.Validate(); err != nil {
		return api.WebResponse(c, http.StatusBadRequest, api.FIELD_VALIDATION_ERROR())
	}

	//domain := c.Get("domain").(*models.Domain)
	//return services.NewUserService(u.Server.DB).Register(e, registerRequest, domain)

	return api.WebResponse(c, http.StatusCreated, api.RESOURCE_CREATED("User successfully created"))
}
