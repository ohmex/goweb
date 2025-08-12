package routes

import (
	// Standard library
	"context"
	"net/http"
	"time"

	// Third-party
	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	// Local
	"goweb/handlers"
	"goweb/interceptor"
	"goweb/server"
	"goweb/services"
	"goweb/util"
)

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
		Skipper:   util.GzipSkipper,
		Level:     2,
		MinLength: 128,
	}))
	server.Echo.Use(middleware.Logger())
	server.Echo.Use(middleware.RequestID())
	server.Echo.Use(interceptor.PerformanceMonitoringMw(server))

	// Middleware assignments for later use
	server.JwtAuthenticationMw = echojwt.WithConfig(jwtConfig)
	server.JwtClaimsAuthorizationMw = interceptor.JwtClaimsAuthorizationMw(server)
	server.CasbinAuthorizationMw = interceptor.CasbinAuthorization(server)

	// Global Handlers
	authHandler := handlers.NewAuthHandler(server)
	registerHandler := handlers.NewRegisterHandler(server)
	socialHandler := handlers.NewSocialHandler(server)

	// Public routes
	server.Echo.GET("/", func(ctx echo.Context) error {
		return ctx.String(http.StatusOK, "Ohmex welcomes you to paradise!")
	})

	// Health check endpoint for monitoring
	server.Echo.GET("/health", func(ctx echo.Context) error {
		// Check database connection
		sqlDB, err := server.DB.DB()
		if err != nil {
			return ctx.JSON(http.StatusServiceUnavailable, map[string]interface{}{
				"status": "unhealthy",
				"error":  "database connection failed",
			})
		}

		if err := sqlDB.Ping(); err != nil {
			return ctx.JSON(http.StatusServiceUnavailable, map[string]interface{}{
				"status": "unhealthy",
				"error":  "database ping failed",
			})
		}

		// Check Redis connection
		ctxRedis := context.Background()
		if err := server.Redis.Ping(ctxRedis).Err(); err != nil {
			return ctx.JSON(http.StatusServiceUnavailable, map[string]interface{}{
				"status": "unhealthy",
				"error":  "redis connection failed",
			})
		}

		return ctx.JSON(http.StatusOK, map[string]interface{}{
			"status":    "healthy",
			"timestamp": time.Now().Unix(),
		})
	})

	server.Echo.Static("/swagger", "docs")
	server.Echo.POST("/login", authHandler.Login)
	server.Echo.POST("/register", registerHandler.Register)
	server.Echo.POST("/refresh", authHandler.RefreshToken)

	// Social login routes
	server.Echo.GET("/auth/google", socialHandler.GoogleLogin)
	server.Echo.GET("/auth/google/callback", socialHandler.GoogleCallback)
	server.Echo.GET("/auth/github", socialHandler.GitHubLogin)
	server.Echo.GET("/auth/github/callback", socialHandler.GitHubCallback)

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
