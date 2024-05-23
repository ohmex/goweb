package models

import (
	"goweb/api"
	"net/http"

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

func (u User) List(c echo.Context) error {
	return api.WebResponse(c, http.StatusNotFound, api.RESOURCE_NOT_FOUND("List users not implemented"))
}

func (u User) Create(c echo.Context) error {
	// do work
	return api.WebResponse(c, http.StatusNotFound, api.RESOURCE_NOT_FOUND("Create user not implemented"))
}

func (u User) Read(c echo.Context) error {
	// do work
	return api.WebResponse(c, http.StatusNotFound, api.RESOURCE_NOT_FOUND("Read user not implemented"))
}

func (u User) Update(c echo.Context) error {
	// do work
	return api.WebResponse(c, http.StatusNotFound, api.RESOURCE_NOT_FOUND("Update user not implemented"))
}

func (u User) Delete(c echo.Context) error {
	// do work
	return api.WebResponse(c, http.StatusNotFound, api.RESOURCE_NOT_FOUND("Delete user not implemented"))
}
