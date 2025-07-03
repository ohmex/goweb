package handlers

import (
	"goweb/api"
	"goweb/server"
	"net/http"

	"github.com/labstack/echo/v4"
)

// BaseInterface defines the standard CRUD interface for handlers.
type BaseInterface interface {
	Type() string
	List(echo.Context) error
	Create(echo.Context) error
	Read(echo.Context) error
	Update(echo.Context) error
	Delete(echo.Context) error
}

// BaseHandler provides a base implementation for handler methods.
type BaseHandler struct {
	Server *server.Server
}

// Type returns the type of the handler.
func (u BaseHandler) Type() string {
	return "Undefined"
}

// List returns a not implemented response for the List operation.
func (u BaseHandler) List(e echo.Context) error {
	return api.WebResponse(e, http.StatusNotFound, api.RESOURCE_NOT_FOUND("BaseHandler: List not implemented"))
}

// Create returns a not implemented response for the Create operation.
func (u BaseHandler) Create(e echo.Context) error {
	return api.WebResponse(e, http.StatusNotFound, api.RESOURCE_NOT_FOUND("BaseHandler: Create not implemented"))
}

// Read returns a not implemented response for the Read operation.
func (u BaseHandler) Read(e echo.Context) error {
	return api.WebResponse(e, http.StatusNotFound, api.RESOURCE_NOT_FOUND("BaseHandler: Read not implemented"))
}

// Update returns a not implemented response for the Update operation.
func (u BaseHandler) Update(e echo.Context) error {
	return api.WebResponse(e, http.StatusNotFound, api.RESOURCE_NOT_FOUND("BaseHandler: Update not implemented"))
}

// Delete returns a not implemented response for the Delete operation.
func (u BaseHandler) Delete(e echo.Context) error {
	return api.WebResponse(e, http.StatusNotFound, api.RESOURCE_NOT_FOUND("BaseHandler: Delete not implemented"))
}
