package server

import (
	"gowebmvc/model"
	"net/http"

	"github.com/casbin/casbin/v2"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Server struct {
	*echo.Echo
	*casbin.Enforcer
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
	server := &Server{e, enforcer}
	server.Pre(middleware.AddTrailingSlash())
	server.InitializeRoutes()
	return server
}

func (server *Server) Authorize(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		// path := c.Request().URL.Path
		// method := c.Request().Method

		// userIDStr := c.QueryParam("user_id")
		// userID, err := strconv.Atoi(userIDStr)
		// if err != nil {
		// 	return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		// }

		// objectIDStr := c.Param("patient_id")
		// objectID, err := strconv.Atoi(objectIDStr)
		// if err != nil {
		// 	return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		// }

		// obj := Object{ID: objectID, Path: path}

		// ok, err := enforcer.Enforce(User{ID: userID}, obj, method)
		// if err != nil {
		// 	return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		// }

		ok := false

		if ok {
			return next(ctx)
		}
		return ctx.String(http.StatusForbidden, "Access denied\n\n")
	}
}

func (server *Server) AddResource(p string, r model.Resource) {
	group := server.Group(p, server.Authorize)
	group.GET("/", r.List)         // Respond back with a the List of Resource
	group.GET("/:id", r.Read)      // Read a single Resource identified by id
	group.POST("/", r.Create)      // Create a new Resource
	group.PUT("/:id", r.Update)    // Update an existing Resource identified by id
	group.DELETE("/:id", r.Delete) // Delete a single Resource identified by id
}
