package routes

import (
	"goweb/interceptor"
	"goweb/server"
	"goweb/server/handlers"
	"goweb/services/token"

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

	server.Echo.GET("/swagger/*", echoSwagger.WrapHandler)

	server.Echo.POST("/login", authHandler.Login)
	server.Echo.POST("/register", registerHandler.Register)
	server.Echo.POST("/refresh", authHandler.RefreshToken)

	// Configure middleware with the custom claims type
	config := echojwt.Config{
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(token.JwtCustomClaims)
		},
		SigningKey: []byte(server.Config.Auth.AccessSecret),
	}

	protectedGroup := server.Echo.Group("")
	protectedGroup.Use(echojwt.WithConfig(config))
	protectedGroup.POST("/logout", authHandler.Logout)

	// Below APIs must validtae the JWT @ REDIS & DB
	protectedGroup.Use(interceptor.ValidateJWT(server))
	protectedGroup.GET("/posts", postHandler.GetPosts)
	protectedGroup.POST("/posts", postHandler.CreatePost)
	protectedGroup.DELETE("/posts/:id", postHandler.DeletePost)
	protectedGroup.PUT("/posts/:id", postHandler.UpdatePost)
}
