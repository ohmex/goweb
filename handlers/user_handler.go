package handlers

import (
	"goweb/api"
	"goweb/models"
	"goweb/requests"
	"goweb/server"
	"goweb/services"
	"goweb/util"
	"net/http"

	"github.com/labstack/echo/v4"
)

// UserHandler provides endpoints for managing users within a domain, including CRUD operations.
type UserHandler struct {
	BaseHandler
	UserService *services.UserService
}

// NewUserHandler initializes the UserHandler with the provided server and its dependencies.
func NewUserHandler(server *server.Server) *UserHandler {
	return &UserHandler{
		BaseHandler: BaseHandler{
			Server: server,
		},
		UserService: services.NewUserService(server.DB),
	}
}

// Type returns the string identifier for the UserHandler.
func (u *UserHandler) Type() string {
	return "User"
}

// Helper to fetch user by UUID in domain
func findUserByUUID(e echo.Context, userService *services.UserService, domain *models.Domain) (*models.User, error) {
	uuid, err := util.GetUUIDParam(e)
	if err != nil {
		return nil, err
	}
	var user models.User
	userService.GetUserByUuidInDomain(&user, uuid, domain)
	if user.ID == 0 {
		return nil, echo.NewHTTPError(http.StatusNotFound, "User not found")
	}
	return &user, nil
}

// List godoc
// @Summary List users
// @Description Returns a list of users for the specified domain.
// @ID user-list
// @Tags User Management
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {array} models.User
// @Failure 400 {object} api.Response
// @Router /api/user [get]
func (u *UserHandler) List(e echo.Context) error {
	d, err := util.ExtractDomain(e)
	if err != nil {
		return api.WebResponse(e, http.StatusBadRequest, api.FIELD_VALIDATION_ERROR())
	}
	domain, _ := d.(*models.Domain)
	var users []*models.User
	u.UserService.GetUsersInDomain(&users, domain)
	return api.WebResponse(e, http.StatusOK, users)
}

// Create godoc
// @Summary Create user
// @Description Creates a new user in the specified domain.
// @ID user-create
// @Tags User Management
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param params body requests.RegisterRequest true "User registration data"
// @Success 201 {object} api.Response
// @Failure 400 {object} api.Response
// @Router /api/user [post]
func (u *UserHandler) Create(e echo.Context) error {
	registerRequest, err := util.BindAndValidate[requests.RegisterRequest](e)
	if err != nil {
		return api.WebResponse(e, http.StatusBadRequest, api.FIELD_VALIDATION_ERROR())
	}
	d, err := util.ExtractDomain(e)
	if err != nil {
		return api.WebResponse(e, http.StatusBadRequest, api.FIELD_VALIDATION_ERROR())
	}
	domain, _ := d.(*models.Domain)
	return u.UserService.Register(e, registerRequest, domain)
}

// Read godoc
// @Summary Get user
// @Description Returns the details of a user by UUID within the specified domain.
// @ID user-read
// @Tags User Management
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param uuid path string true "User UUID"
// @Success 200 {object} models.User
// @Failure 400 {object} api.Response
// @Failure 404 {object} api.Response
// @Router /api/user/{uuid} [get]
func (u *UserHandler) Read(e echo.Context) error {
	d, err := util.ExtractDomain(e)
	if err != nil {
		return api.WebResponse(e, http.StatusBadRequest, api.FIELD_VALIDATION_ERROR())
	}
	domain, _ := d.(*models.Domain)
	user, err := findUserByUUID(e, u.UserService, domain)
	if err != nil {
		return api.WebResponse(e, http.StatusNotFound, api.RESOURCE_NOT_FOUND("User not found"))
	}
	return api.WebResponse(e, http.StatusOK, user)
}

// Update godoc
// @Summary Update user
// @Description Modifies the details of a user by UUID within the specified domain.
// @ID user-update
// @Tags User Management
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param uuid path string true "User UUID"
// @Param params body requests.UpdateRequest true "User update data"
// @Success 200 {object} models.User
// @Failure 400 {object} api.Response
// @Failure 404 {object} api.Response
// @Router /api/user/{uuid} [put]
func (u *UserHandler) Update(e echo.Context) error {
	updateRequest, err := util.BindAndValidate[requests.UpdateRequest](e)
	if err != nil {
		return api.WebResponse(e, http.StatusBadRequest, api.FIELD_VALIDATION_ERROR())
	}
	d, err := util.ExtractDomain(e)
	if err != nil {
		return api.WebResponse(e, http.StatusBadRequest, api.FIELD_VALIDATION_ERROR())
	}
	domain, _ := d.(*models.Domain)
	user, err := findUserByUUID(e, u.UserService, domain)
	if err != nil {
		return api.WebResponse(e, http.StatusNotFound, api.RESOURCE_NOT_FOUND("User not found"))
	}
	user.Name = updateRequest.Name
	u.UserService.UpdateUser(user)

	return api.WebResponse(e, http.StatusOK, user)
}

// Delete godoc
// @Summary Delete user
// @Description Removes a user by UUID from the specified domain.
// @ID user-delete
// @Tags User Management
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param uuid path string true "User UUID"
// @Success 200 {object} api.Response
// @Failure 400 {object} api.Response
// @Failure 404 {object} api.Response
// @Router /api/user/{uuid} [delete]
func (u *UserHandler) Delete(e echo.Context) error {
	d, err := util.ExtractDomain(e)
	if err != nil {
		return api.WebResponse(e, http.StatusBadRequest, api.FIELD_VALIDATION_ERROR())
	}
	domain, _ := d.(*models.Domain)
	uuid, err := util.GetUUIDParam(e)
	if err != nil {
		return api.WebResponse(e, http.StatusBadRequest, err)
	}
	var user models.User
	err = u.UserService.DeleteUserByUuidInDomain(&user, uuid, domain)
	if err != nil {
		return api.WebResponse(e, http.StatusNotFound, api.RESOURCE_NOT_FOUND("User not found"))
	}
	return api.WebResponse(e, http.StatusOK, api.RESOURCE_DELETED("User deleted"))
}
