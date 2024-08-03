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
	domain, ok := e.Get("domain").(*models.Domain)
	if !ok {
		return api.WebResponse(e, http.StatusBadRequest, api.FIELD_VALIDATION_ERROR("Missing domain information"))
	}

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

	domain, ok := e.Get("domain").(*models.Domain)
	if !ok {
		return api.WebResponse(e, http.StatusBadRequest, api.FIELD_VALIDATION_ERROR("Missing domain information"))
	}

	return services.NewRoleService(h.Server.DB).Create(e, roleRequest, domain)
}

func (h RoleHandler) Read(e echo.Context) error {
	var role models.Role
	domain := e.Get("domain").(*models.Domain)
	uuid := e.Param("uuid")
	services.NewRoleService(h.Server.DB).GetRoleByUuidInDomain(&role, uuid, domain)
	if role.ID == 0 {
		return api.WebResponse(e, http.StatusNotFound, api.RESOURCE_NOT_FOUND("Role not found"))
	}
	return api.WebResponse(e, http.StatusOK, role)
}

func (h RoleHandler) Update(e echo.Context) error {
	var role models.Role
	service := services.NewRoleService(h.Server.DB)

	// Extract and validate the role request
	roleRequest := new(requests.RoleRequest)
	if err := e.Bind(roleRequest); err != nil {
		return api.WebResponse(e, http.StatusBadRequest, api.FIELD_VALIDATION_ERROR())
	}

	if err := roleRequest.Validate(); err != nil {
		return api.WebResponse(e, http.StatusBadRequest, api.FIELD_VALIDATION_ERROR())
	}

	// Extract UUID and domain from the request context
	uuid := e.Param("uuid")
	domain, ok := e.Get("domain").(*models.Domain)
	if !ok {
		return api.WebResponse(e, http.StatusBadRequest, api.FIELD_VALIDATION_ERROR("Missing domain information"))
	}

	service.GetRoleByUuidInDomain(&role, uuid, domain)
	if role.ID == 0 {
		return api.WebResponse(e, http.StatusNotFound, api.RESOURCE_NOT_FOUND("Role not found"))
	}
	role.Name = roleRequest.Name
	service.UpdateRole(&role)

	return api.WebResponse(e, http.StatusOK, role)
}

func (h RoleHandler) Delete(e echo.Context) error {
	var role models.Role
	uuid := e.Param("uuid")
	domain, OK := e.Get("domain").(*models.Domain)
	if !OK {
		return api.WebResponse(e, http.StatusBadRequest, api.FIELD_VALIDATION_ERROR("Missing domain information"))
	}

	ok := services.NewRoleService(h.Server.DB).DeleteRoleByUuidInDomain(&role, uuid, domain)
	if ok != nil {
		return api.WebResponse(e, http.StatusNotFound, api.RESOURCE_NOT_FOUND("Role not found"))
	}
	return api.WebResponse(e, http.StatusOK, api.RESOURCE_DELETED("Role deleted"))
}
