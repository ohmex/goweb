package handlers

import (
	"goweb/api"
	"goweb/server"
	"net/http"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/labstack/echo/v4"
)

// DomainHandler provides endpoints for managing domains, including CRUD operations (some not implemented).
type DomainHandler struct {
	server *server.Server
}

// NewDomainHandler initializes the DomainHandler with the provided server.
func NewDomainHandler(server *server.Server) *DomainHandler {
	return &DomainHandler{server: server}
}

// Type returns the string identifier for the DomainHandler.
func (u DomainHandler) Type() string {
	return "Domain"
}

// List returns a sample response for listing domains.
func (u DomainHandler) List(e echo.Context) error {
	return api.WebResponse(e, http.StatusOK, api.STATUS_OK(gofakeit.Paragraph(100, 10, 10, " ")))
}

// Create returns a not implemented response for creating a domain.
func (u DomainHandler) Create(e echo.Context) error {
	return api.WebResponse(e, http.StatusNotFound, api.RESOURCE_NOT_FOUND("Create Domain not implemented"))
}

// Read returns a not implemented response for reading a domain.
func (u DomainHandler) Read(e echo.Context) error {
	return api.WebResponse(e, http.StatusNotFound, api.RESOURCE_NOT_FOUND("Read Domain not implemented"))
}

// Update returns a not implemented response for updating a domain.
func (u DomainHandler) Update(e echo.Context) error {
	return api.WebResponse(e, http.StatusNotFound, api.RESOURCE_NOT_FOUND("Update Domain not implemented"))
}

// Delete returns a not implemented response for deleting a domain.
func (u DomainHandler) Delete(e echo.Context) error {
	return api.WebResponse(e, http.StatusNotFound, api.RESOURCE_NOT_FOUND("Delete Domain not implemented"))
}
