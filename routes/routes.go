package routes

import (
	// Standard library
	"net/http"
	"strings"

	// Third-party
	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"

	// Local
	"goweb/handlers"
	"goweb/interceptor"
	"goweb/server"
	"goweb/services"
)

// gzipSkipper skips Gzip middleware for health and swagger endpoints
func gzipSkipper(c echo.Context) bool {
	p := c.Path()
	return p == "/health" || strings.HasPrefix(p, "/swagger")
}

// ConfigureRoutes sets up all routes and middleware for the server
func ConfigureRoutes(server *server.Server) {
	// JWT middleware configuration
	jwtConfig := echojwt.Config{
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(services.JwtCustomClaims)
		},
		ContextKey:    "token",
		SigningMethod: jwt.SigningMethodHS512.Name,
		SigningKey:    []byte(server.Config.Auth.AccessSecret),
	}

	// Global middleware
	server.Echo.Use(middleware.Recover())
	server.Echo.Use(middleware.GzipWithConfig(middleware.GzipConfig{
		Skipper:   gzipSkipper,
		Level:     2,
		MinLength: 128,
	}))
	server.Echo.Use(middleware.Logger())
	server.Echo.Use(middleware.RequestID())

	// Middleware assignments for later use
	server.JwtAuthenticationMw = echojwt.WithConfig(jwtConfig)
	server.JwtClaimsAuthorizationMw = interceptor.JwtClaimsAuthorizationMw(server)
	server.CasbinAuthorizationMw = interceptor.CasbinAuthorization(server)

	// Global Handlers
	authHandler := handlers.NewAuthHandler(server)
	registerHandler := handlers.NewRegisterHandler(server)

	// Public routes
	server.Echo.GET("/", func(ctx echo.Context) error {
		return ctx.String(http.StatusOK, "Ohmex welcomes you to paradise!")
	})
	server.Echo.GET("/swagger/*", echoSwagger.WrapHandler)
	server.Echo.POST("/login", authHandler.Login)
	server.Echo.POST("/register", registerHandler.Register)
	server.Echo.POST("/refresh", authHandler.RefreshToken)

	// Protected routes group
	protected := server.Echo.Group("")
	protected.Use(server.JwtAuthenticationMw, server.JwtClaimsAuthorizationMw)
	protected.POST("/logout", authHandler.Logout)

	// Resource Handlers
	roleHandler := handlers.NewRoleHandler(server)
	userHandler := handlers.NewUserHandler(server)
	postHandler := handlers.NewPostHandler(server)

	// API resource routes grouped under /api
	api := server.Echo.Group("/api")
	api.Use(server.JwtAuthenticationMw, server.JwtClaimsAuthorizationMw, server.CasbinAuthorizationMw)
	addResource(api, "/role", roleHandler, server)
	addResource(api, "/user", userHandler, server)
	addResource(api, "/post", postHandler, server)
}

// addResource adds RESTful resource routes to the given group
func addResource(group *echo.Group, p string, h handlers.BaseInterface, server *server.Server) {
	sub := group.Group(p)
	sub.GET("", h.List, interceptor.ResourceAuthorization(server, h.Type(), "List"))              // List resources
	sub.GET("/:uuid", h.Read, interceptor.ResourceAuthorization(server, h.Type(), "Read"))        // Read resource
	sub.POST("", h.Create, interceptor.ResourceAuthorization(server, h.Type(), "Create"))         // Create resource
	sub.PUT("/:uuid", h.Update, interceptor.ResourceAuthorization(server, h.Type(), "Update"))    // Update resource
	sub.DELETE("/:uuid", h.Delete, interceptor.ResourceAuthorization(server, h.Type(), "Delete")) // Delete resource
}
