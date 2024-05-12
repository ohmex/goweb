package server

import (
	"gowebmvc/controller"
	"gowebmvc/model"
	"net/http"

	"github.com/labstack/echo/v4"
)

func (server *Server) InitializeRoutes() {
	server.GET("/", func(ctx echo.Context) error {
		return ctx.String(http.StatusOK, "Hello, World!")
	})

	server.POST("/login/", controller.HandleLogin)

	server.AddResource("/user", model.User{})
}
