package controller

import (
	"gowebmvc/auth"
	"gowebmvc/model"
	"net/http"

	"github.com/labstack/echo/v4"
)

func HandleLogin(c echo.Context) error {
	username := c.FormValue("username")
	password := c.FormValue("password")

	// Throws unauthorized error
	if username != "user" || password != "password" {
		return echo.ErrUnauthorized
	}

	u := model.User{
		ID:     1,
		Name:   "Sachin",
		Role:   "Admin",
		Domain: "Reliance",
	}

	token, _ := auth.CreateToken(u)

	return c.JSON(http.StatusOK, echo.Map{
		"token": token,
	})
}
