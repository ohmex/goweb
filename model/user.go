package model

import (
	"net/http"

	"github.com/go-faker/faker/v4"
	"github.com/labstack/echo/v4"
)

type User struct{}

func (u User) List(c echo.Context) error {
	return c.String(http.StatusOK, faker.Paragraph())
}

func (u User) Create(c echo.Context) error {
	// do work
	return echo.NewHTTPError(http.StatusNotFound, "Create resource not implemented")
}

func (u User) Read(c echo.Context) error {
	// do work
	return echo.NewHTTPError(http.StatusNotFound, "Read resource not implemented")
}

func (u User) Update(c echo.Context) error {
	// do work
	return echo.NewHTTPError(http.StatusNotFound, "Update resource not implemented")
}

func (u User) Delete(c echo.Context) error {
	// do work
	return echo.NewHTTPError(http.StatusNotFound, "Delete resource not implemented")
}
