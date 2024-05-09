package framework

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type Resource interface {
	List(echo.Context) error
	Create(echo.Context) error
	Read(echo.Context) error
	Update(echo.Context) error
	Delete(echo.Context) error
}

type BaseResource struct{}

// List default implementation. Returns a 404
func (v BaseResource) List(c echo.Context) error {
	return echo.NewHTTPError(http.StatusNotFound, "resource not implemented")
}

// Create default implementation. Returns a 404
func (v BaseResource) Create(c echo.Context) error {
	return echo.NewHTTPError(http.StatusNotFound, "resource not implemented")
}

// Show default implementation. Returns a 404
func (v BaseResource) Read(c echo.Context) error {
	return echo.NewHTTPError(http.StatusNotFound, "resource not implemented")
}

// Update default implementation. Returns a 404
func (v BaseResource) Update(c echo.Context) error {
	return echo.NewHTTPError(http.StatusNotFound, "resource not implemented")
}

// Delete default implementation. Returns a 404
func (v BaseResource) Delete(c echo.Context) error {
	return echo.NewHTTPError(http.StatusNotFound, "resource not implemented")
}
