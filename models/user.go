package models

import (
	"goweb/api"
	"goweb/server"
	"net/http"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/labstack/echo/v4"
)

type User struct {
	Base
	Email    string    `json:"email" gorm:"type:varchar(200);"`
	Name     string    `json:"name" gorm:"type:varchar(200);"`
	Password string    `json:"password" gorm:"type:varchar(200);"`
	Tenants  []*Tenant `json:"tenants" gorm:"many2many:tenant_user;"`
	Posts    []*Post   `json:"posts"`
}

func (u User) Type() string {
	return "User"
}

func (u User) List(server *server.Server) func(echo.Context) error {
	return func(e echo.Context) error {
		return api.WebResponse(e, http.StatusOK, api.STATUS_OK(gofakeit.Sentence(100)))
	}
}

func (u User) Create(server *server.Server) func(echo.Context) error {
	return func(e echo.Context) error {
		return api.WebResponse(e, http.StatusNotFound, api.RESOURCE_NOT_FOUND("Create user not implemented"))
	}
}

func (u User) Read(server *server.Server) func(echo.Context) error {
	return func(e echo.Context) error {
		return api.WebResponse(e, http.StatusNotFound, api.RESOURCE_NOT_FOUND("Read user not implemented"))
	}
}

func (u User) Update(server *server.Server) func(echo.Context) error {
	return func(e echo.Context) error {
		return api.WebResponse(e, http.StatusNotFound, api.RESOURCE_NOT_FOUND("Update user not implemented"))
	}
}

func (u User) Delete(server *server.Server) func(echo.Context) error {
	return func(e echo.Context) error {
		return api.WebResponse(e, http.StatusNotFound, api.RESOURCE_NOT_FOUND("Delete user not implemented"))
	}
}
