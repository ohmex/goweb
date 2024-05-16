package server

import (
	"goweb/config"
	"goweb/db"

	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type Server struct {
	Echo   *echo.Echo
	DB     *gorm.DB
	Redis  *redis.Client
	Config *config.Config
}

func NewServer(cfg *config.Config) *Server {
	return &Server{
		Echo:   echo.New(),
		DB:     db.InitDB(cfg),
		Redis:  db.InitRedis(cfg),
		Config: cfg,
	}
}

func (server *Server) Start(addr string) error {
	return server.Echo.Start(":" + addr)
}
