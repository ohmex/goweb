package interceptor

import (
	"context"
	"fmt"
	"goweb/api"
	"goweb/models"
	"goweb/server"
	"goweb/services"
	"goweb/util"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

func CasbinAuthorization(server *server.Server) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Get tenant from header identified by UUID
			tenant := c.Request().Header.Get("tenant")

			// Get user name as UUID
			user := c.Get("user").(*models.User).UUID.String()

			// Check, user though mabe assiciated with a tenant in DB
			// Does he has casbin domain assigned to him or not
			domains, _ := server.Casbin.GetDomainsForUser(user)
			ok := util.Contains[string](domains, tenant)

			if ok {
				return next(c)
			}

			return api.WebResponse(c, http.StatusForbidden, api.CASBIN_UNAUTHORIZED())
		}
	}
}

// Middleware for additional steps:
// 1. Check the token info exists in Redis
// 2. Check the user exists in DB
// 3. Add the user data to Echo Context
// 4. Prolong the Redis TTL of the current token pair
func JwtAuthorization(server *server.Server) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			token := c.Get("token").(*jwt.Token)
			claims := token.Claims.(*services.JwtCustomClaims)

			user, err := services.NewTokenService(server).ValidateToken(claims, false)
			if err != nil {
				return api.WebResponse(c, http.StatusUnauthorized, err)
			}

			c.Set("user", user)

			go func() {
				server.Redis.Expire(context.Background(), fmt.Sprintf("token-%d", claims.ID), time.Minute*services.AutoLogoffMinutes)
			}()

			return next(c)
		}
	}
}

// Check if resource action can be done
func ResourceAuthorization(server *server.Server, resource string, action string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(e echo.Context) error {
			user := e.Get("user").(*models.User).UUID.String()
			tenant := e.Request().Header.Get("tenant")

			// Enforce used nomentlature (sub, dom, obj, act)
			ok, _ := server.Casbin.Enforce(user, tenant, resource, action)

			if !ok {
				return api.WebResponse(e, http.StatusForbidden, api.CASBIN_UNAUTHORIZED())
			}

			return next(e)
		}
	}
}
