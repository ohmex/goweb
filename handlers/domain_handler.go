package handlers

import (
	"goweb/api"
	"goweb/models"
	"goweb/requests"
	"goweb/server"
	"goweb/services"
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

// List godoc
// @Summary List domains
// @Description Returns a sample response for listing domains.
// @ID domain-list
// @Tags Domain Management
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} api.Response
// @Router /api/domain [get]
func (u DomainHandler) List(e echo.Context) error {
	return api.WebResponse(e, http.StatusOK, api.STATUS_OK(gofakeit.Paragraph(100, 10, 10, " ")))
}

// Create godoc
// @Summary Create domain
// @Description Creates a new domain with automatic partition creation.
// @ID domain-create
// @Tags Domain Management
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param domain body requests.CreateDomainRequest true "Domain creation request"
// @Success 201 {object} api.Response
// @Failure 400 {object} api.Response
// @Failure 500 {object} api.Response
// @Router /api/domain [post]
func (u DomainHandler) Create(e echo.Context) error {
	var request requests.CreateDomainRequest

	if err := e.Bind(&request); err != nil {
		return api.WebResponse(e, http.StatusBadRequest, api.FIELD_VALIDATION_ERROR("Invalid request format"))
	}

	if err := request.Validate(); err != nil {
		return api.WebResponse(e, http.StatusBadRequest, api.FIELD_VALIDATION_ERROR(err.Error()))
	}

	// Create domain service
	domainService := services.NewDomainService(u.server.DB)

	// Create domain model
	domain := &models.Domain{
		Name: request.Name,
	}

	// Create domain with automatic partition creation
	if err := domainService.CreateDomain(domain); err != nil {
		return api.WebResponse(e, http.StatusInternalServerError, api.RESOURCE_CREATION_FAILED("Failed to create domain: "+err.Error()))
	}

	return api.WebResponse(e, http.StatusCreated, api.RESOURCE_CREATED("Domain created successfully"))
}

// Read godoc
// @Summary Get domain
// @Description Not implemented.
// @ID domain-read
// @Tags Domain Management
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param uuid path string true "Domain UUID"
// @Failure 404 {object} api.Response
// @Router /api/domain/{uuid} [get]
func (u DomainHandler) Read(e echo.Context) error {
	return api.WebResponse(e, http.StatusNotFound, api.RESOURCE_NOT_FOUND("Read Domain not implemented"))
}

// Update godoc
// @Summary Update domain
// @Description Not implemented.
// @ID domain-update
// @Tags Domain Management
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param uuid path string true "Domain UUID"
// @Failure 404 {object} api.Response
// @Router /api/domain/{uuid} [put]
func (u DomainHandler) Update(e echo.Context) error {
	return api.WebResponse(e, http.StatusNotFound, api.RESOURCE_NOT_FOUND("Update Domain not implemented"))
}

// Delete godoc
// @Summary Delete domain
// @Description Not implemented.
// @ID domain-delete
// @Tags Domain Management
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param uuid path string true "Domain UUID"
// @Failure 404 {object} api.Response
// @Router /api/domain/{uuid} [delete]
func (u DomainHandler) Delete(e echo.Context) error {
	return api.WebResponse(e, http.StatusNotFound, api.RESOURCE_NOT_FOUND("Delete Domain not implemented"))
}
