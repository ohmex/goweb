package server

import (
	"goweb/config"
	"goweb/db"

	"github.com/casbin/casbin/v2"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type Server struct {
	Echo   *echo.Echo
	DB     *gorm.DB
	Redis  *redis.Client
	Config *config.Config
	Casbin *casbin.Enforcer
}

func NewServer(cfg *config.Config) *Server {
	adaptor, _ := gormadapter.NewAdapter("sqlite3", "casbin.db")
	enforcer, _ := casbin.NewEnforcer("casbin/model.conf", adaptor)
	enforcer.EnableLog(true)

	e := echo.New()
	e.Pre(middleware.RemoveTrailingSlash())

	enforcer.LoadPolicy()
	return &Server{
		Echo:   e,
		DB:     db.InitDB(cfg),
		Redis:  db.InitRedis(cfg),
		Config: cfg,
		Casbin: enforcer,
	}
}

func (server *Server) Start(addr string) error {
	return server.Echo.Start(":" + addr)
}
