package routes

import (
	"goweb/handlers"
	"goweb/interceptor"
	"goweb/server"
	"goweb/services"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
)

func ConfigureRoutes(server *server.Server) {
	postHandler := handlers.NewPostHandlers(server)
	authHandler := handlers.NewAuthHandler(server)
	registerHandler := handlers.NewRegisterHandler(server)

	server.Echo.Use(middleware.Logger())

	server.Echo.GET("", func(ctx echo.Context) error {
		return ctx.String(http.StatusOK, "Welcome to HOME!")
	})
	server.Echo.GET("/hello", func(ctx echo.Context) error {
		return ctx.String(http.StatusOK, "Hello, World!")
	})
	server.Echo.GET("/swagger/*", echoSwagger.WrapHandler)

	server.Echo.POST("/login", authHandler.Login)
	server.Echo.POST("/register", registerHandler.Register)
	server.Echo.POST("/refresh", authHandler.RefreshToken)

	// Configure middleware with the custom claims type
	config := echojwt.Config{
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(services.JwtCustomClaims)
		},
		SigningKey: []byte(server.Config.Auth.AccessSecret),
	}

	protectedGroup := server.Echo.Group("")
	protectedGroup.Use(echojwt.WithConfig(config))

	protectedGroup.Use(interceptor.ValidateJWT(server))
	protectedGroup.POST("/logout", authHandler.Logout)

	apiGroup := server.Echo.Group("/api")
	apiGroup.Use(echojwt.WithConfig(config))
	apiGroup.Use(interceptor.ValidateJWT(server))
	apiGroup.Use(interceptor.CasbinAuthorizer)

	apiGroup.GET("/posts", postHandler.GetPosts)
	apiGroup.POST("/posts", postHandler.CreatePost)
	apiGroup.DELETE("/posts/:id", postHandler.DeletePost)
	apiGroup.PUT("/posts/:id", postHandler.UpdatePost)
}
