package server

import (
	"gowebmvc/model"
	"net/http"

	"github.com/labstack/echo/v4"
)

func (app *App) InitializeRoutes() {
	app.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	app.AddResource("/user", model.User{})
}
