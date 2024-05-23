package models

import (
	"goweb/api"
	"goweb/server"
	"net/http"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/labstack/echo/v4"
)

type Post struct {
	Base
	Title   string `json:"title" gorm:"type:text"`
	Content string `json:"content" gorm:"type:text"`
	UserID  uint
	User    User `gorm:"foreignkey:UserID"`
}

func (p Post) List(server *server.Server) func(echo.Context) error {
	return func(e echo.Context) error {
		return api.WebResponse(e, http.StatusOK, api.STATUS_OK(gofakeit.Paragraph(100, 10, 10, " ")))
	}
}

func (p Post) Create(server *server.Server) func(echo.Context) error {
	return func(e echo.Context) error {
		return api.WebResponse(e, http.StatusNotFound, api.RESOURCE_NOT_FOUND("Create post not implemented"))
	}
}

func (p Post) Read(server *server.Server) func(echo.Context) error {
	return func(e echo.Context) error {
		return api.WebResponse(e, http.StatusNotFound, api.RESOURCE_NOT_FOUND("Read post not implemented"))
	}
}

func (p Post) Update(server *server.Server) func(echo.Context) error {
	return func(e echo.Context) error {
		return api.WebResponse(e, http.StatusNotFound, api.RESOURCE_NOT_FOUND("Update post not implemented"))
	}
}

func (p Post) Delete(server *server.Server) func(echo.Context) error {
	return func(e echo.Context) error {
		return api.WebResponse(e, http.StatusNotFound, api.RESOURCE_NOT_FOUND("Delete post not implemented"))
	}
}
