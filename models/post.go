package models

import (
	"goweb/api"
	"net/http"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type Post struct {
	gorm.Model
	Title   string `json:"title" gorm:"type:text"`
	Content string `json:"content" gorm:"type:text"`
	UserID  uint
	User    User `gorm:"foreignkey:UserID"`
}

func (u Post) List(c echo.Context) error {
	return c.String(http.StatusOK, gofakeit.Paragraph(100, 10, 10, " "))
}

func (u Post) Create(c echo.Context) error {
	// do work
	return api.WebResponse(c, http.StatusNotFound, api.RESOURCE_NOT_FOUND("Create post not implemented"))
}

func (u Post) Read(c echo.Context) error {
	// do work
	return api.WebResponse(c, http.StatusNotFound, api.RESOURCE_NOT_FOUND("Read post not implemented"))
}

func (u Post) Update(c echo.Context) error {
	// do work
	return api.WebResponse(c, http.StatusNotFound, api.RESOURCE_NOT_FOUND("Update post not implemented"))
}

func (u Post) Delete(c echo.Context) error {
	// do work
	return api.WebResponse(c, http.StatusNotFound, api.RESOURCE_NOT_FOUND("Delete post not implemented"))
}
