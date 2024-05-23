package models

import (
	"goweb/server"

	"github.com/labstack/echo/v4"
)

type Resource interface {
	Type() string
	List(server *server.Server) func(echo.Context) error
	Create(server *server.Server) func(echo.Context) error
	Read(server *server.Server) func(echo.Context) error
	Update(server *server.Server) func(echo.Context) error
	Delete(server *server.Server) func(echo.Context) error
}
