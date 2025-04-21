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
	Echo                     *echo.Echo
	DB                       *gorm.DB
	Redis                    *redis.Client
	Config                   *config.Config
	Casbin                   *casbin.Enforcer
	JwtAuthenticationMw      echo.MiddlewareFunc
	JwtClaimsAuthorizationMw echo.MiddlewareFunc
	CasbinAuthorizationMw    echo.MiddlewareFunc
}

func NewServer(cfg *config.Config) *Server {
	database := db.InitDB(cfg)
	//adaptor, _ := gormadapter.NewAdapter("sqlite3", "casbin.db")
	adaptor, _ := gormadapter.NewAdapterByDBUseTableName(database, "", "casbin")
	enforcer, _ := casbin.NewEnforcer("casbin/model.conf", adaptor)
	//enforcer.EnableLog(true)
	enforcer.LoadPolicy()

	e := echo.New()
	//e.HideBanner = true
	//e.HidePort = true
	e.Pre(middleware.RemoveTrailingSlash())

	return &Server{
		Echo:   e,
		DB:     database,
		Redis:  db.InitRedis(cfg),
		Config: cfg,
		Casbin: enforcer,
	}
}

func (server *Server) Start(addr string) error {
	return server.Echo.Start(":" + addr)
}
