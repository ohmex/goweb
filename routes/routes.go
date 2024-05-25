package routes

import (
	"goweb/handlers"
	"goweb/interceptor"
	"goweb/server"
	"goweb/services"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
)

func ConfigureRoutes(server *server.Server) {
	// Configure middleware with the custom claims type
	config := echojwt.Config{
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(services.JwtCustomClaims)
		},
		ContextKey:    "token",
		SigningMethod: jwt.SigningMethodHS512.Name,
		SigningKey:    []byte(server.Config.Auth.AccessSecret),
	}
	server.JwtAuthenticationMw = echojwt.WithConfig(config)
	server.JwtAuthorizationMw = interceptor.JwtAuthorization(server)
	server.CasbinAuthorizationMw = interceptor.CasbinAuthorization(server)

	authHandler := handlers.NewAuthHandler(server)
	registerHandler := handlers.NewRegisterHandler(server)

	server.Echo.Use(middleware.GzipWithConfig(middleware.GzipConfig{
		Skipper: func(c echo.Context) bool {
			return strings.Contains(c.Request().URL.Path, "swagger")
		},
		Level:     5,
		MinLength: 64,
	}))
	server.Echo.Use(middleware.Logger())

	server.Echo.GET("", func(ctx echo.Context) error {
		return ctx.String(http.StatusOK, "Welcome to HOME!")
	})
	server.Echo.GET("/swagger/*", echoSwagger.WrapHandler)

	server.Echo.POST("/login", authHandler.Login)
	server.Echo.POST("/register", registerHandler.Register)
	server.Echo.POST("/refresh", authHandler.RefreshToken)

	protectedGroup := server.Echo.Group("")
	protectedGroup.Use(server.JwtAuthenticationMw)
	protectedGroup.Use(server.JwtAuthorizationMw)
	protectedGroup.POST("/logout", authHandler.Logout)

	AddResource(server, "/user", handlers.NewUserHandler(server))
	AddResource(server, "/post", handlers.NewPostHandler(server))
	AddResource(server, "/tenant", handlers.NewTenantHandler(server))
}

func AddResource(server *server.Server, p string, h handlers.BaseInterface) {
	group := server.Echo.Group("/api" + p)
	group.Use(server.JwtAuthenticationMw)
	group.Use(server.JwtAuthorizationMw)
	group.Use(server.CasbinAuthorizationMw)
	group.GET("", h.List, interceptor.ResourceAuthorization(server, h.Type(), "List"))              // Respond back with a the List of Resource
	group.GET("/:uuid", h.Read, interceptor.ResourceAuthorization(server, h.Type(), "Read"))        // Read a single Resource identified by id
	group.POST("", h.Create, interceptor.ResourceAuthorization(server, h.Type(), "Create"))         // Create a new Resource
	group.PUT("/:uuid", h.Update, interceptor.ResourceAuthorization(server, h.Type(), "Update"))    // Update an existing Resource identified by id
	group.DELETE("/:uuid", h.Delete, interceptor.ResourceAuthorization(server, h.Type(), "Delete")) // Delete a single Resource identified by id
}
