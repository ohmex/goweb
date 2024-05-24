package handlers

import (
	"goweb/api"
	"goweb/server"
	"net/http"

	"github.com/labstack/echo/v4"
)

type PostHandler struct {
	server *server.Server
}

func NewPostHandler(server *server.Server) *PostHandler {
	return &PostHandler{server: server}
}

func (u PostHandler) Type() string {
	return "Post"
}

func (u PostHandler) List(e echo.Context) error {
	return api.WebResponse(e, http.StatusNotFound, api.RESOURCE_NOT_FOUND("List Posts not implemented"))
}

func (u PostHandler) Create(e echo.Context) error {
	return api.WebResponse(e, http.StatusNotFound, api.RESOURCE_NOT_FOUND("Create Post not implemented"))
}

func (u PostHandler) Read(e echo.Context) error {
	return api.WebResponse(e, http.StatusNotFound, api.RESOURCE_NOT_FOUND("Read Post not implemented"))
}

func (u PostHandler) Update(e echo.Context) error {
	return api.WebResponse(e, http.StatusNotFound, api.RESOURCE_NOT_FOUND("Update Post not implemented"))
}

func (u PostHandler) Delete(e echo.Context) error {
	return api.WebResponse(e, http.StatusNotFound, api.RESOURCE_NOT_FOUND("Delete Post not implemented"))
}
