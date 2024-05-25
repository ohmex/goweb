package handlers

import (
	"goweb/api"
	"goweb/server"
	"net/http"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/labstack/echo/v4"
)

type DomainHandler struct {
	server *server.Server
}

func NewDomainHandler(server *server.Server) *DomainHandler {
	return &DomainHandler{server: server}
}

func (u DomainHandler) Type() string {
	return "Domain"
}

func (u DomainHandler) List(e echo.Context) error {
	return api.WebResponse(e, http.StatusOK, api.STATUS_OK(gofakeit.Paragraph(100, 10, 10, " ")))
}

func (u DomainHandler) Create(e echo.Context) error {
	return api.WebResponse(e, http.StatusNotFound, api.RESOURCE_NOT_FOUND("Create Domain not implemented"))
}

func (u DomainHandler) Read(e echo.Context) error {
	return api.WebResponse(e, http.StatusNotFound, api.RESOURCE_NOT_FOUND("Read Domain not implemented"))
}

func (u DomainHandler) Update(e echo.Context) error {
	return api.WebResponse(e, http.StatusNotFound, api.RESOURCE_NOT_FOUND("Update Domain not implemented"))
}

func (u DomainHandler) Delete(e echo.Context) error {
	return api.WebResponse(e, http.StatusNotFound, api.RESOURCE_NOT_FOUND("Delete Domain not implemented"))
}
