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
	UserService *services.UserService
}

func NewUserHandler(server *server.Server) *UserHandler {
	return &UserHandler{
		BaseHandler: BaseHandler{
			Server: server,
		},
		UserService: services.NewUserService(server.DB),
	}
}

func (u *UserHandler) Type() string {
	return "User"
}

// getDomain safely extracts the domain from context
func getDomain(e echo.Context) (*models.Domain, error) {
	domain, ok := e.Get("domain").(*models.Domain)
	if !ok || domain == nil {
		return nil, echo.NewHTTPError(http.StatusBadRequest, "Invalid domain context")
	}
	return domain, nil
}

func (u *UserHandler) List(e echo.Context) error {
	domain, err := getDomain(e)
	if err != nil {
		return api.WebResponse(e, http.StatusBadRequest, api.FIELD_VALIDATION_ERROR())
	}
	var users []*models.User
	u.UserService.GetUsersInDomain(&users, domain)
	return api.WebResponse(e, http.StatusOK, users)
}

func (u *UserHandler) Create(e echo.Context) error {
	registerRequest := new(requests.RegisterRequest)

	if err := e.Bind(registerRequest); err != nil {
		return api.WebResponse(e, http.StatusBadRequest, api.FIELD_VALIDATION_ERROR())
	}

	if err := registerRequest.Validate(); err != nil {
		return api.WebResponse(e, http.StatusBadRequest, api.FIELD_VALIDATION_ERROR())
	}
	domain, err := getDomain(e)
	if err != nil {
		return api.WebResponse(e, http.StatusBadRequest, api.FIELD_VALIDATION_ERROR())
	}
	return u.UserService.Register(e, registerRequest, domain)
}

func (u *UserHandler) Read(e echo.Context) error {
	var user models.User
	domain, err := getDomain(e)
	if err != nil {
		return api.WebResponse(e, http.StatusBadRequest, api.FIELD_VALIDATION_ERROR())
	}
	uuid := e.Param("uuid")
	u.UserService.GetUserByUuidInDomain(&user, uuid, domain)
	if user.ID == 0 {
		return api.WebResponse(e, http.StatusNotFound, api.RESOURCE_NOT_FOUND("User not found"))
	}
	return api.WebResponse(e, http.StatusOK, user)
}

// TODO: Check the resource being accessed belongs to the domain that user have access to
func (u *UserHandler) Update(e echo.Context) error {
	var user models.User
	updateRequest := new(requests.UpdateRequest)

	if err := e.Bind(updateRequest); err != nil {
		return api.WebResponse(e, http.StatusBadRequest, api.FIELD_VALIDATION_ERROR())
	}

	if err := updateRequest.Validate(); err != nil {
		return api.WebResponse(e, http.StatusBadRequest, api.FIELD_VALIDATION_ERROR())
	}
	domain, err := getDomain(e)
	if err != nil {
		return api.WebResponse(e, http.StatusBadRequest, api.FIELD_VALIDATION_ERROR())
	}
	uuid := e.Param("uuid")
	u.UserService.GetUserByUuidInDomain(&user, uuid, domain)
	if user.ID == 0 {
		return api.WebResponse(e, http.StatusNotFound, api.RESOURCE_NOT_FOUND("User not found"))
	}
	user.Name = updateRequest.Name
	u.UserService.UpdateUser(&user)

	return api.WebResponse(e, http.StatusOK, user)
}

// TODO: Check the resource being accessed belongs to the domain that user have access to
func (u *UserHandler) Delete(e echo.Context) error {
	var user models.User
	domain, err := getDomain(e)
	if err != nil {
		return api.WebResponse(e, http.StatusBadRequest, api.FIELD_VALIDATION_ERROR())
	}
	uuid := e.Param("uuid")
	err = u.UserService.DeleteUserByUuidInDomain(&user, uuid, domain)
	if err != nil {
		return api.WebResponse(e, http.StatusNotFound, api.RESOURCE_NOT_FOUND("User not found"))
	}
	return api.WebResponse(e, http.StatusOK, api.RESOURCE_DELETED("User deleted"))
}
