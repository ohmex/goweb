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

func (h *RoleHandler) getDomain(e echo.Context) (*models.Domain, error) {
	domain, ok := e.Get("domain").(*models.Domain)
	if !ok || domain == nil {
		return nil, api.FIELD_VALIDATION_ERROR("Missing domain information")
	}
	return domain, nil
}

func (h *RoleHandler) Type() string {
	return "Role"
}

func (h *RoleHandler) List(e echo.Context) error {
	domain, err := h.getDomain(e)
	if err != nil {
		return api.WebResponse(e, http.StatusBadRequest, err)
	}

	var roles []*models.Role
	service := services.NewRoleService(h.Server.DB)
	if err := service.GetRolesInDomain(&roles, domain); err != nil {
		return api.WebResponse(e, http.StatusInternalServerError, err)
	}
	return api.WebResponse(e, http.StatusOK, roles)
}

func (h *RoleHandler) Create(e echo.Context) error {
	roleRequest := new(requests.RoleRequest)

	if err := e.Bind(roleRequest); err != nil {
		return api.WebResponse(e, http.StatusBadRequest, api.FIELD_VALIDATION_ERROR())
	}

	if err := roleRequest.Validate(); err != nil {
		return api.WebResponse(e, http.StatusBadRequest, api.FIELD_VALIDATION_ERROR())
	}

	domain, err := h.getDomain(e)
	if err != nil {
		return api.WebResponse(e, http.StatusBadRequest, err)
	}

	service := services.NewRoleService(h.Server.DB)
	if err := service.Create(e, roleRequest, domain); err != nil {
		return api.WebResponse(e, http.StatusInternalServerError, err)
	}
	return api.WebResponse(e, http.StatusCreated, api.RESOURCE_CREATED("Role created"))
}

func (h *RoleHandler) Read(e echo.Context) error {
	var role models.Role
	domain, err := h.getDomain(e)
	if err != nil {
		return api.WebResponse(e, http.StatusBadRequest, err)
	}
	uuid := e.Param("uuid")
	service := services.NewRoleService(h.Server.DB)
	if err := service.GetRoleByUuidInDomain(&role, uuid, domain); err != nil {
		return api.WebResponse(e, http.StatusInternalServerError, err)
	}
	if role.ID == 0 {
		return api.WebResponse(e, http.StatusNotFound, api.RESOURCE_NOT_FOUND("Role not found"))
	}
	return api.WebResponse(e, http.StatusOK, role)
}

func (h *RoleHandler) Update(e echo.Context) error {
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
	domain, err := h.getDomain(e)
	if err != nil {
		return api.WebResponse(e, http.StatusBadRequest, err)
	}

	if err := service.GetRoleByUuidInDomain(&role, uuid, domain); err != nil {
		return api.WebResponse(e, http.StatusInternalServerError, err)
	}
	if role.ID == 0 {
		return api.WebResponse(e, http.StatusNotFound, api.RESOURCE_NOT_FOUND("Role not found"))
	}
	role.Name = roleRequest.Name
	if err := service.UpdateRole(&role); err != nil {
		return api.WebResponse(e, http.StatusInternalServerError, err)
	}

	return api.WebResponse(e, http.StatusOK, role)
}

func (h *RoleHandler) Delete(e echo.Context) error {
	var role models.Role
	uuid := e.Param("uuid")
	domain, err := h.getDomain(e)
	if err != nil {
		return api.WebResponse(e, http.StatusBadRequest, err)
	}

	service := services.NewRoleService(h.Server.DB)
	if err := service.DeleteRoleByUuidInDomain(&role, uuid, domain); err != nil {
		return api.WebResponse(e, http.StatusNotFound, api.RESOURCE_NOT_FOUND("Role not found"))
	}
	return api.WebResponse(e, http.StatusOK, api.RESOURCE_DELETED("Role deleted"))
}
