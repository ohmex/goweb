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
	Server  *server.Server
	service *services.RoleService
}

func NewRoleHandler(server *server.Server) *RoleHandler {
	return &RoleHandler{
		Server:  server,
		service: services.NewRoleService(server.DB),
	}
}

// getDomain extracts and validates domain from context
func (h *RoleHandler) getDomain(e echo.Context) (*models.Domain, error) {
	domain, ok := e.Get("domain").(*models.Domain)
	if !ok || domain == nil {
		return nil, api.FIELD_VALIDATION_ERROR("Missing domain information")
	}
	return domain, nil
}

// validateRoleRequest validates and binds role request
func (h *RoleHandler) validateRoleRequest(e echo.Context) (*requests.RoleRequest, error) {
	roleRequest := new(requests.RoleRequest)

	if err := e.Bind(roleRequest); err != nil {
		return nil, api.FIELD_VALIDATION_ERROR("Invalid request format")
	}

	if err := roleRequest.Validate(); err != nil {
		return nil, api.FIELD_VALIDATION_ERROR("Validation failed: " + err.Error())
	}

	return roleRequest, nil
}

// getRoleByUUID retrieves a role by UUID within a domain
func (h *RoleHandler) getRoleByUUID(uuid string, domain *models.Domain) (*models.Role, error) {
	var role models.Role
	if err := h.service.GetRoleByUuidInDomain(&role, uuid, domain); err != nil {
		return nil, err
	}

	if role.ID == 0 {
		return nil, api.RESOURCE_NOT_FOUND("Role not found")
	}

	return &role, nil
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
	if err := h.service.GetRolesInDomain(&roles, domain); err != nil {
		return api.WebResponse(e, http.StatusInternalServerError, api.INTERNAL_SERVICE_ERROR("Failed to fetch roles"))
	}

	return api.WebResponse(e, http.StatusOK, roles)
}

func (h *RoleHandler) Create(e echo.Context) error {
	roleRequest, err := h.validateRoleRequest(e)
	if err != nil {
		return api.WebResponse(e, http.StatusBadRequest, err)
	}

	domain, err := h.getDomain(e)
	if err != nil {
		return api.WebResponse(e, http.StatusBadRequest, err)
	}

	if err := h.service.Create(e, roleRequest, domain); err != nil {
		return api.WebResponse(e, http.StatusInternalServerError, err)
	}

	return api.WebResponse(e, http.StatusCreated, api.RESOURCE_CREATED("Role created successfully"))
}

func (h *RoleHandler) Read(e echo.Context) error {
	domain, err := h.getDomain(e)
	if err != nil {
		return api.WebResponse(e, http.StatusBadRequest, err)
	}

	uuid := e.Param("uuid")
	if uuid == "" {
		return api.WebResponse(e, http.StatusBadRequest, api.FIELD_VALIDATION_ERROR("Role UUID is required"))
	}

	role, err := h.getRoleByUUID(uuid, domain)
	if err != nil {
		if err.Error() == "Role not found" {
			return api.WebResponse(e, http.StatusNotFound, err)
		}
		return api.WebResponse(e, http.StatusInternalServerError, api.INTERNAL_SERVICE_ERROR("Failed to fetch role"))
	}

	return api.WebResponse(e, http.StatusOK, role)
}

func (h *RoleHandler) Update(e echo.Context) error {
	roleRequest, err := h.validateRoleRequest(e)
	if err != nil {
		return api.WebResponse(e, http.StatusBadRequest, err)
	}

	domain, err := h.getDomain(e)
	if err != nil {
		return api.WebResponse(e, http.StatusBadRequest, err)
	}

	uuid := e.Param("uuid")
	if uuid == "" {
		return api.WebResponse(e, http.StatusBadRequest, api.FIELD_VALIDATION_ERROR("Role UUID is required"))
	}

	role, err := h.getRoleByUUID(uuid, domain)
	if err != nil {
		if err.Error() == "Role not found" {
			return api.WebResponse(e, http.StatusNotFound, err)
		}
		return api.WebResponse(e, http.StatusInternalServerError, api.INTERNAL_SERVICE_ERROR("Failed to fetch role"))
	}

	// Update role name
	role.Name = roleRequest.Name

	if err := h.service.UpdateRole(role); err != nil {
		return api.WebResponse(e, http.StatusInternalServerError, api.INTERNAL_SERVICE_ERROR("Failed to update role"))
	}

	return api.WebResponse(e, http.StatusOK, role)
}

func (h *RoleHandler) Delete(e echo.Context) error {
	domain, err := h.getDomain(e)
	if err != nil {
		return api.WebResponse(e, http.StatusBadRequest, err)
	}

	uuid := e.Param("uuid")
	if uuid == "" {
		return api.WebResponse(e, http.StatusBadRequest, api.FIELD_VALIDATION_ERROR("Role UUID is required"))
	}

	var role models.Role
	if err := h.service.DeleteRoleByUuidInDomain(&role, uuid, domain); err != nil {
		return api.WebResponse(e, http.StatusNotFound, api.RESOURCE_NOT_FOUND("Role not found"))
	}

	return api.WebResponse(e, http.StatusOK, api.RESOURCE_DELETED("Role deleted successfully"))
}
