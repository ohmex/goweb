package handlers

import (
	"goweb/api"
	"goweb/server"
	"net/http"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/labstack/echo/v4"
)

type TenantHandler struct {
	server *server.Server
}

func NewTenantHandler(server *server.Server) *TenantHandler {
	return &TenantHandler{server: server}
}

func (u TenantHandler) Type() string {
	return "Tenant"
}

func (u TenantHandler) List(e echo.Context) error {
	return api.WebResponse(e, http.StatusOK, api.STATUS_OK(gofakeit.Paragraph(100, 10, 10, " ")))
}

func (u TenantHandler) Create(e echo.Context) error {
	return api.WebResponse(e, http.StatusNotFound, api.RESOURCE_NOT_FOUND("Create Tenant not implemented"))
}

func (u TenantHandler) Read(e echo.Context) error {
	return api.WebResponse(e, http.StatusNotFound, api.RESOURCE_NOT_FOUND("Read Tenant not implemented"))
}

func (u TenantHandler) Update(e echo.Context) error {
	return api.WebResponse(e, http.StatusNotFound, api.RESOURCE_NOT_FOUND("Update Tenant not implemented"))
}

func (u TenantHandler) Delete(e echo.Context) error {
	return api.WebResponse(e, http.StatusNotFound, api.RESOURCE_NOT_FOUND("Delete Tenant not implemented"))
}
