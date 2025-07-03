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

// List returns a not implemented response for listing posts.
func (u PostHandler) List(e echo.Context) error {
	return api.WebResponse(e, http.StatusNotFound, api.RESOURCE_NOT_FOUND("List Posts not implemented"))
}

// Create returns a not implemented response for creating a post.
func (u PostHandler) Create(e echo.Context) error {
	return api.WebResponse(e, http.StatusNotFound, api.RESOURCE_NOT_FOUND("Create Post not implemented"))
}

// Read returns a not implemented response for reading a post.
func (u PostHandler) Read(e echo.Context) error {
	return api.WebResponse(e, http.StatusNotFound, api.RESOURCE_NOT_FOUND("Read Post not implemented"))
}

// Update returns a not implemented response for updating a post.
func (u PostHandler) Update(e echo.Context) error {
	return api.WebResponse(e, http.StatusNotFound, api.RESOURCE_NOT_FOUND("Update Post not implemented"))
}

// Delete returns a not implemented response for deleting a post.
func (u PostHandler) Delete(e echo.Context) error {
	return api.WebResponse(e, http.StatusNotFound, api.RESOURCE_NOT_FOUND("Delete Post not implemented"))
}
