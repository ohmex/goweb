package interceptor

import (
	"context"
	"fmt"
	"goweb/responses"
	"goweb/server"
	"net/http"
	"time"

	tokenService "goweb/services/token"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

// Middleware for additional steps:
// 1. Check the user exists in DB
// 2. Check the token info exists in Redis
// 3. Add the user DB data to Context
// 4. Prolong the Redis TTL of the current token pair
func ValidateJWT(server *server.Server) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			token := c.Get("user").(*jwt.Token)
			claims := token.Claims.(*tokenService.JwtCustomClaims)

			user, err := tokenService.NewTokenService(server).ValidateToken(claims, false)
			if err != nil {
				return responses.ErrorResponse(c, http.StatusUnauthorized, err.Error())
			}

			c.Set("currentUser", user)

			go func() {
				server.Redis.Expire(context.Background(), fmt.Sprintf("token-%d", claims.ID), time.Minute*tokenService.AutoLogoffMinutes)
			}()

			return next(c)
		}
	}
}
