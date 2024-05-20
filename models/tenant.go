package models

import (
	"goweb/api"
	"net/http"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/labstack/echo/v4"
)

type Tenant struct {
	Base
	Name  string  `json:"name" gorm:"type:text"`
	Users []*User `json:"users" gorm:"many2many:tenant_user;"`
}

func (u Tenant) List(c echo.Context) error {
	return api.WebResponse(c, http.StatusOK, api.STATUS_OK(gofakeit.Paragraph(100, 10, 10, " ")))
}

func (u Tenant) Create(c echo.Context) error {
	// do work
	return api.WebResponse(c, http.StatusNotFound, api.RESOURCE_NOT_FOUND("Create post not implemented"))
}

func (u Tenant) Read(c echo.Context) error {
	// do work
	return api.WebResponse(c, http.StatusNotFound, api.RESOURCE_NOT_FOUND("Read post not implemented"))
}

func (u Tenant) Update(c echo.Context) error {
	// do work
	return api.WebResponse(c, http.StatusNotFound, api.RESOURCE_NOT_FOUND("Update post not implemented"))
}

func (u Tenant) Delete(c echo.Context) error {
	// do work
	return api.WebResponse(c, http.StatusNotFound, api.RESOURCE_NOT_FOUND("Delete post not implemented"))
}
