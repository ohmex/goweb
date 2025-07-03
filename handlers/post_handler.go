package handlers

import (
	"goweb/api"
	"goweb/server"
	"net/http"

	"github.com/labstack/echo/v4"
)

// PostHandler provides endpoints for managing posts, including CRUD operations (not implemented).
type PostHandler struct {
	server *server.Server
}

// NewPostHandler initializes the PostHandler with the provided server.
func NewPostHandler(server *server.Server) *PostHandler {
	return &PostHandler{server: server}
}

// Type returns the string identifier for the PostHandler.
func (u PostHandler) Type() string {
	return "Post"
}

// List godoc
// @Summary List posts
// @Description Not implemented.
// @ID post-list
// @Tags Post Management
// @Accept json
// @Produce json
// @Failure 404 {object} api.Response
// @Router /posts [get]
func (u PostHandler) List(e echo.Context) error {
	return api.WebResponse(e, http.StatusNotFound, api.RESOURCE_NOT_FOUND("List Posts not implemented"))
}

// Create godoc
// @Summary Create post
// @Description Not implemented.
// @ID post-create
// @Tags Post Management
// @Accept json
// @Produce json
// @Failure 404 {object} api.Response
// @Router /posts [post]
func (u PostHandler) Create(e echo.Context) error {
	return api.WebResponse(e, http.StatusNotFound, api.RESOURCE_NOT_FOUND("Create Post not implemented"))
}

// Read godoc
// @Summary Get post
// @Description Not implemented.
// @ID post-read
// @Tags Post Management
// @Accept json
// @Produce json
// @Param uuid path string true "Post UUID"
// @Failure 404 {object} api.Response
// @Router /posts/{uuid} [get]
func (u PostHandler) Read(e echo.Context) error {
	return api.WebResponse(e, http.StatusNotFound, api.RESOURCE_NOT_FOUND("Read Post not implemented"))
}

// Update godoc
// @Summary Update post
// @Description Not implemented.
// @ID post-update
// @Tags Post Management
// @Accept json
// @Produce json
// @Param uuid path string true "Post UUID"
// @Failure 404 {object} api.Response
// @Router /posts/{uuid} [put]
func (u PostHandler) Update(e echo.Context) error {
	return api.WebResponse(e, http.StatusNotFound, api.RESOURCE_NOT_FOUND("Update Post not implemented"))
}

// Delete godoc
// @Summary Delete post
// @Description Not implemented.
// @ID post-delete
// @Tags Post Management
// @Accept json
// @Produce json
// @Param uuid path string true "Post UUID"
// @Failure 404 {object} api.Response
// @Router /posts/{uuid} [delete]
func (u PostHandler) Delete(e echo.Context) error {
	return api.WebResponse(e, http.StatusNotFound, api.RESOURCE_NOT_FOUND("Delete Post not implemented"))
}
