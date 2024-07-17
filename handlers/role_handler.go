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

type RoleHandler struct {
	Server *server.Server
}

func NewRoleHandler(server *server.Server) *RoleHandler {
	return &RoleHandler{Server: server}
}

func (h RoleHandler) Type() string {
	return "Role"
}

func (h RoleHandler) List(e echo.Context) error {
	domain := e.Get("domain").(*models.Domain)
	var roles []*models.Role
	services.NewRoleService(h.Server.DB).GetRolesInDomain(&roles, domain)
	return api.WebResponse(e, http.StatusOK, roles)
}

func (h RoleHandler) Create(e echo.Context) error {
	roleRequest := new(requests.RoleRequest)

	if err := e.Bind(roleRequest); err != nil {
		return api.WebResponse(e, http.StatusBadRequest, api.FIELD_VALIDATION_ERROR())
	}

	if err := roleRequest.Validate(); err != nil {
		return api.WebResponse(e, http.StatusBadRequest, api.FIELD_VALIDATION_ERROR())
	}

	domain := e.Get("domain").(*models.Domain)
	return services.NewRoleService(h.Server.DB).Create(e, roleRequest, domain)
}

func (h RoleHandler) Read(e echo.Context) error {
	return api.WebResponse(e, http.StatusNotFound, api.RESOURCE_NOT_FOUND("Read Role not implemented"))
}

func (h RoleHandler) Update(e echo.Context) error {
	return api.WebResponse(e, http.StatusNotFound, api.RESOURCE_NOT_FOUND("Update Role not implemented"))
}

func (h RoleHandler) Delete(e echo.Context) error {
	return api.WebResponse(e, http.StatusNotFound, api.RESOURCE_NOT_FOUND("Delete Role not implemented"))
}
