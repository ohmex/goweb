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

type UserHandler struct {
	BaseHandler
}

func NewUserHandler(server *server.Server) *UserHandler {
	return &UserHandler{
		BaseHandler: BaseHandler{
			Server: server,
		},
	}
}

func (u UserHandler) Type() string {
	return "User"
}

func (u UserHandler) List(e echo.Context) error {
	domain := e.Get("domain").(*models.Domain)
	var users []*models.User
	services.NewUserService(u.Server.DB).GetUsersByDomain(&users, domain)
	return api.WebResponse(e, http.StatusOK, users)
}

func (u UserHandler) Create(e echo.Context) error {
	registerRequest := new(requests.RegisterRequest)

	if err := e.Bind(registerRequest); err != nil {
		return api.WebResponse(e, http.StatusBadRequest, api.FIELD_VALIDATION_ERROR())
	}

	if err := registerRequest.Validate(); err != nil {
		return api.WebResponse(e, http.StatusBadRequest, api.FIELD_VALIDATION_ERROR())
	}

	domain := e.Get("domain").(*models.Domain)
	return services.NewUserService(u.Server.DB).Register(e, registerRequest, domain)
}

func (u UserHandler) Read(e echo.Context) error {
	var user models.User
	domain := e.Get("domain").(*models.Domain)
	uuid := e.Param("uuid")
	services.NewUserService(u.Server.DB).GetUserByDomainAndUUID(&user, domain, uuid)
	if user.ID == 0 {
		return api.WebResponse(e, http.StatusNotFound, api.RESOURCE_NOT_FOUND("User not found"))
	}
	return api.WebResponse(e, http.StatusOK, user)
}

// TODO: Check requested user belongs to current domain
func (u UserHandler) Update(e echo.Context) error {
	return api.WebResponse(e, http.StatusNotFound, api.RESOURCE_NOT_FOUND("Update User not implemented"))
}

// TODO: Check requested user belongs to current domain
func (u UserHandler) Delete(e echo.Context) error {
	return api.WebResponse(e, http.StatusNotFound, api.RESOURCE_NOT_FOUND("Delete User not implemented"))
}
