package handlers

import (
	"goweb/api"
	"goweb/server"
	"net/http"

	"github.com/labstack/echo/v4"
)

type BaseInterface interface {
	Type() string
	List(echo.Context) error
	Create(echo.Context) error
	Read(echo.Context) error
	Update(echo.Context) error
	Delete(echo.Context) error
}

type BaseHandler struct {
	Server *server.Server
}

func (u BaseHandler) Type() string {
	return "Undefined"
}

func (u BaseHandler) List(e echo.Context) error {
	return api.WebResponse(e, http.StatusNotFound, api.RESOURCE_NOT_FOUND("BaseHandler: List not implemented"))
}

func (u BaseHandler) Create(e echo.Context) error {
	return api.WebResponse(e, http.StatusNotFound, api.RESOURCE_NOT_FOUND("BaseHandler: Create not implemented"))
}

func (u BaseHandler) Read(e echo.Context) error {
	return api.WebResponse(e, http.StatusNotFound, api.RESOURCE_NOT_FOUND("BaseHandler: Read not implemented"))
}

func (u BaseHandler) Update(e echo.Context) error {
	return api.WebResponse(e, http.StatusNotFound, api.RESOURCE_NOT_FOUND("BaseHandler: Update not implemented"))
}

func (u BaseHandler) Delete(e echo.Context) error {
	return api.WebResponse(e, http.StatusNotFound, api.RESOURCE_NOT_FOUND("BaseHandler: Delete not implemented"))
}
