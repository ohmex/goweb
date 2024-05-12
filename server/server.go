package server

import (
	"gowebmvc/auth"
	"gowebmvc/model"
	"gowebmvc/mware"
	"gowebmvc/service"
	"os"

	"github.com/casbin/casbin/v2"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Server struct {
	*echo.Echo
	*casbin.Enforcer
	*service.RedisClient
	jwtValidator echo.MiddlewareFunc
}

func New() *Server {
	adaptor, _ := gormadapter.NewAdapter("sqlite3", "casbin.db")
	enforcer, _ := casbin.NewEnforcer("config/model.conf", adaptor)
	enforcer.EnableLog(true)

	enforcer.LoadPolicy()

	// c.AddRoleForUserInDomain("Sachin", "Admin", "Reliance")
	// c.AddRoleForUserInDomain("Sachin", "Admin", "More2Store")

	// c.AddRoleForUserInDomain("Anant", "Operator", "Reliance")
	// c.AddRoleForUserInDomain("Aditya", "Operator", "More2Store")

	// c.AddPolicy("Admin", "Reliance", "User", "Read")
	// c.AddPolicy("Admin", "Reliance", "User", "Write")
	// c.AddPolicy("Admin", "More2Store", "User", "Read")
	// c.AddPolicy("Admin", "More2Store", "User", "Write")

	// c.AddPolicy("Operator", "Reliance", "User", "Read")
	// c.AddPolicy("Operator", "More2Store", "User", "Read")

	// c.SavePolicy()

	e := echo.New()

	// Configure middleware with the custom claims type
	config := echojwt.Config{
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(auth.CustomClaims)
		},
		SigningKey: []byte(os.Getenv("ACCESS_SECRET")),
	}

	server := &Server{e, enforcer, service.NewRedisClient(), echojwt.WithConfig(config)}
	server.Pre(middleware.AddTrailingSlash())
	server.InitializeRoutes()
	return server
}

func (server *Server) AddResource(p string, r model.Resource) {
	group := server.Group(p)
	group.Use(server.jwtValidator)
	group.Use(mware.CasbinAuthorizer)
	group.GET("/", r.List)         // Respond back with a the List of Resource
	group.GET("/:id", r.Read)      // Read a single Resource identified by id
	group.POST("/", r.Create)      // Create a new Resource
	group.PUT("/:id", r.Update)    // Update an existing Resource identified by id
	group.DELETE("/:id", r.Delete) // Delete a single Resource identified by id
}
