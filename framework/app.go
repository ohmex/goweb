package framework

import (
	"github.com/casbin/casbin/v2"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type App struct {
	*echo.Echo
	*casbin.Enforcer
}

func New() *App {
	g, _ := gormadapter.NewAdapter("sqlite3", "casbin.db")
	c, _ := casbin.NewEnforcer("model.conf", g)

	c.LoadPolicy()

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
	a := &App{e, c}
	a.Pre(middleware.AddTrailingSlash())
	return a
}

func (a *App) Resource(p string, r Resource) {
	g := a.Group(p)
	g.GET("/", r.List)         // GET /users => r.List
	g.GET("/:id", r.Read)      // GET /users/:id => r.Read
	g.POST("/", r.Create)      // POST /users => r.Create
	g.PUT("/:id", r.Update)    // PUT /users/:id => r.Update
	g.DELETE("/:id", r.Delete) // DELETE /users/:id => r.Delete
}
