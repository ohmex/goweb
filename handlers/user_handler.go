package handlers

import (
	"goweb/api"
	"goweb/models"
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
	tenant := e.Get("tenant").(*models.Tenant)
	var users []*models.User
	services.NewUserService(u.Server.DB).GetUsersByTenant(&users, tenant)
	return api.WebResponse(e, http.StatusOK, users)
}

func (u UserHandler) Create(e echo.Context) error {
	return api.WebResponse(e, http.StatusNotFound, api.RESOURCE_NOT_FOUND("Create User not implemented"))
}

func (u UserHandler) Read(e echo.Context) error {
	return api.WebResponse(e, http.StatusNotFound, api.RESOURCE_NOT_FOUND("Read User not implemented"))
}

func (u UserHandler) Update(e echo.Context) error {
	return api.WebResponse(e, http.StatusNotFound, api.RESOURCE_NOT_FOUND("Update User not implemented"))
}

func (u UserHandler) Delete(e echo.Context) error {
	return api.WebResponse(e, http.StatusNotFound, api.RESOURCE_NOT_FOUND("Delete User not implemented"))
}
