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
// BaseHandler implements the BaseInterface and provides default not-implemented responses for CRUD operations.
type BaseHandler struct {
	Server *server.Server
}

// Type returns the type of the handler.
// Type returns the string identifier for the BaseHandler.
func (u BaseHandler) Type() string {
	return "Undefined"
}

// List returns a not implemented response for the List operation.
// List returns a not implemented response for listing resources.
func (u BaseHandler) List(e echo.Context) error {
	return api.WebResponse(e, http.StatusNotFound, api.RESOURCE_NOT_FOUND("BaseHandler: List not implemented"))
}

// Create returns a not implemented response for the Create operation.
// Create returns a not implemented response for creating a resource.
func (u BaseHandler) Create(e echo.Context) error {
	return api.WebResponse(e, http.StatusNotFound, api.RESOURCE_NOT_FOUND("BaseHandler: Create not implemented"))
}

// Read returns a not implemented response for the Read operation.
// Read returns a not implemented response for reading a resource.
func (u BaseHandler) Read(e echo.Context) error {
	return api.WebResponse(e, http.StatusNotFound, api.RESOURCE_NOT_FOUND("BaseHandler: Read not implemented"))
}

// Update returns a not implemented response for the Update operation.
// Update returns a not implemented response for updating a resource.
func (u BaseHandler) Update(e echo.Context) error {
	return api.WebResponse(e, http.StatusNotFound, api.RESOURCE_NOT_FOUND("BaseHandler: Update not implemented"))
}

// Delete returns a not implemented response for the Delete operation.
// Delete returns a not implemented response for deleting a resource.
func (u BaseHandler) Delete(e echo.Context) error {
	return api.WebResponse(e, http.StatusNotFound, api.RESOURCE_NOT_FOUND("BaseHandler: Delete not implemented"))
}
