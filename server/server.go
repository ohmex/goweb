package server

import (
	"goweb/config"
	"goweb/db"
	"goweb/models"

	"github.com/casbin/casbin/v2"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
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
	//enforcer.EnableLog(true)
	enforcer.LoadPolicy()

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
	//CreateDefaultPolicy(enforcer)
	//TestPolicy(enforcer)
	TestTableData(server.DB)

	return server.Echo.Start(":" + addr)
}

func CreateDefaultPolicy(casbin *casbin.Enforcer) {
	// Admin on a domain can have all rights
	// These are added automatically as a tenant is added
	// Also as a tenant is added, its Admin has to be added with these permissions
	casbin.AddPolicy("Admin", "DMart", "User", "Read")
	casbin.AddPolicy("Admin", "DMart", "User", "Create")
	casbin.AddPolicy("Admin", "DMart", "User", "Update")
	casbin.AddPolicy("Admin", "DMart", "User", "Delete")
	casbin.AddPolicy("Admin", "Reliance", "User", "Read")
	casbin.AddPolicy("Admin", "Reliance", "User", "Create")
	casbin.AddPolicy("Admin", "Reliance", "User", "Update")
	casbin.AddPolicy("Admin", "Reliance", "User", "Delete")

	// This gets added as default for Manager of Tenant
	casbin.AddPolicy("Manager", "Reliance", "User", "Read")
	casbin.AddPolicy("Manager", "Reliance", "User", "Update")
	casbin.AddPolicy("Manager", "DMart", "User", "Read")
	casbin.AddPolicy("Manager", "DMart", "User", "Update")

	casbin.AddPolicy("Operator", "Reliance", "User", "Read")
	casbin.AddPolicy("Operator", "DMart", "User", "Read")

	// These will have to be added when a user for a tenant is added
	casbin.AddRoleForUserInDomain("Sachin", "Admin", "Reliance")
	casbin.AddRoleForUserInDomain("Sachin", "Admin", "DMart")

	casbin.AddRoleForUserInDomain("UserA", "Admin", "Reliance")
	casbin.AddRoleForUserInDomain("UserB", "Manager", "Reliance")
	casbin.AddRoleForUserInDomain("UserC", "Operator", "Reliance")

	casbin.AddRoleForUserInDomain("UserD", "Admin", "DMart")
	casbin.AddRoleForUserInDomain("UserE", "Manager", "DMart")
	casbin.AddRoleForUserInDomain("UserF", "Operator", "DMart")

	casbin.SavePolicy()
}

func TestPolicy(casbin *casbin.Enforcer) {
	permissions := casbin.GetPermissionsForUserInDomain("Sachin", "Reliance")
	log.Info().Interface("Sachin", permissions).Send()
	permissions = casbin.GetPermissionsForUserInDomain("Sachin", "DMart")
	log.Info().Interface("Sachin", permissions).Send()

	permissions = casbin.GetPermissionsForUserInDomain("UserA", "Reliance")
	log.Info().Interface("UserA", permissions).Send()
	permissions = casbin.GetPermissionsForUserInDomain("UserA", "DMart")
	log.Info().Interface("UserA", permissions).Send()

	permissions = casbin.GetPermissionsForUserInDomain("UserB", "Reliance")
	log.Info().Interface("UserB", permissions).Send()
	permissions = casbin.GetPermissionsForUserInDomain("UserB", "DMart")
	log.Info().Interface("UserB", permissions).Send()

	permissions = casbin.GetPermissionsForUserInDomain("UserC", "Reliance")
	log.Info().Interface("UserC", permissions).Send()
	permissions = casbin.GetPermissionsForUserInDomain("UserC", "DMart")
	log.Info().Interface("UserC", permissions).Send()

	permissions = casbin.GetPermissionsForUserInDomain("UserD", "Reliance")
	log.Info().Interface("UserD", permissions).Send()
	permissions = casbin.GetPermissionsForUserInDomain("UserD", "DMart")
	log.Info().Interface("UserD", permissions).Send()

	permissions = casbin.GetPermissionsForUserInDomain("UserE", "Reliance")
	log.Info().Interface("UserE", permissions).Send()
	permissions = casbin.GetPermissionsForUserInDomain("UserE", "DMart")
	log.Info().Interface("UserE", permissions).Send()

	permissions = casbin.GetPermissionsForUserInDomain("UserF", "Reliance")
	log.Info().Interface("UserF", permissions).Send()
	permissions = casbin.GetPermissionsForUserInDomain("UserF", "DMart")
	log.Info().Interface("UserF", permissions).Send()
}

func TestTableData(db *gorm.DB) {
	var users []*models.User

	//db.Find(&users, "name = ?", "Sachin")
	//db.Where("name = ?", "Sachin").Find(&users)
	//db.Preload("Tenants").Where("name = ?", "Sachin").Find(&users)

	db.Preload("Tenants").Where("name = ?", "Sachin").Find(&users)

	log.Info().Interface("USERS", users).Send()
}
