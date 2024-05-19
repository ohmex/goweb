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
	Echo                  *echo.Echo
	DB                    *gorm.DB
	Redis                 *redis.Client
	Config                *config.Config
	Casbin                *casbin.Enforcer
	JwtAuthenticationMw   echo.MiddlewareFunc
	JwtAuthorizationMw    echo.MiddlewareFunc
	CasbinAuthorizationMw echo.MiddlewareFunc
}

func NewServer(cfg *config.Config) *Server {
	adaptor, _ := gormadapter.NewAdapter("sqlite3", "casbin.db")
	enforcer, _ := casbin.NewEnforcer("casbin/model.conf", adaptor)
	enforcer.EnableLog(true)
	enforcer.LoadPolicy()
	//CreateDefaultPolicy(enforcer)

	e := echo.New()
	e.Pre(middleware.RemoveTrailingSlash())

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

func CreateDefaultPolicy(casbin *casbin.Enforcer) {
	casbin.AddRoleForUserInDomain("Sachin", "Admin", "Reliance")
	casbin.AddRoleForUserInDomain("Sachin", "Admin", "More2Store")

	casbin.AddRoleForUserInDomain("Anant", "Operator", "Reliance")
	casbin.AddRoleForUserInDomain("Aditya", "Operator", "More2Store")

	casbin.AddPolicy("Admin", "Reliance", "User", "Read")
	casbin.AddPolicy("Admin", "Reliance", "User", "Write")
	casbin.AddPolicy("Admin", "More2Store", "User", "Read")
	casbin.AddPolicy("Admin", "More2Store", "User", "Write")

	casbin.AddPolicy("Operator", "Reliance", "User", "Read")
	casbin.AddPolicy("Operator", "More2Store", "User", "Read")

	casbin.SavePolicy()
}
