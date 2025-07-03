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

// Middleware for additional steps:
// 1. Check the token info exists in Redis
// 2. Check the user exists in DB
// 3. Add the user data to Echo Context
// 4. Prolong the Redis TTL of the current token pair
func JwtClaimsAuthorizationMw(server *server.Server) echo.MiddlewareFunc {
	tokenService := services.NewTokenService(server)
	domainService := services.NewDomainService(server.DB)
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			tokenVal := c.Get("token")
			token, ok := tokenVal.(*jwt.Token)
			if !ok || token == nil {
				return api.WebResponse(c, http.StatusUnauthorized, api.FIELD_VALIDATION_ERROR("Invalid token in context"))
			}

			claims, ok := token.Claims.(*services.JwtCustomClaims)
			if !ok || claims == nil {
				return api.WebResponse(c, http.StatusUnauthorized, api.FIELD_VALIDATION_ERROR("Invalid claims in token"))
			}

			domainuuid := c.Request().Header.Get("domain")

			user, err := tokenService.ValidateToken(claims, false)
			if err != nil {
				return api.WebResponse(c, http.StatusUnauthorized, err) // TODO: Change this return statement
			}

			c.Set("user", user)

			domain := new(models.Domain)
			err = domainService.GetDomainByUUID(domain, domainuuid)
			if err != nil || domain.ID == 0 {
				return api.WebResponse(c, http.StatusBadRequest, api.FIELD_VALIDATION_ERROR("Missing or incorrect domain"))
			}

			c.Set("domain", domain)

			// Asynchronously update Redis TTL for the token in a goroutine
			go func(userID uint64) {
				ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
				defer cancel()
				key := fmt.Sprintf("token-%d", userID)
				if err := server.Redis.Expire(ctx, key, time.Minute*services.AutoLogoffMinutes).Err(); err != nil {
					// TODO: Use a proper logger here for production
					fmt.Printf("Failed to update Redis TTL for %s: %v\n", key, err)
				}
			}(claims.UserID)

			return next(c)
		}
	}
}

// Domain Authorization: Check if user belongs to domain
func CasbinAuthorization(server *server.Server) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Get domain from header identified by UUID
			domain := c.Request().Header.Get("domain")
			userVal := c.Get("user")
			userModel, ok := userVal.(*models.User)
			if !ok || userModel == nil {
				return api.WebResponse(c, http.StatusUnauthorized, api.FIELD_VALIDATION_ERROR("User not found in context"))
			}
			// Get user name as UUID
			user := userModel.UUID.String()

			// Check, user though maybe associated with a domain in DB
			// Does he has casbin domain assigned to him or not
			domains, err := server.Casbin.GetDomainsForUser(user)
			if err != nil {
				return api.WebResponse(c, http.StatusInternalServerError, api.FIELD_VALIDATION_ERROR("Failed to get user domains"))
			}
			if util.Contains(domains, domain) {
				return next(c)
			}
			return api.WebResponse(c, http.StatusForbidden, api.CASBIN_UNAUTHORIZED("Access denied - domain authorization failed"))
		}
	}
}

// Resource Authorization: Check if resource action can be done
func ResourceAuthorization(server *server.Server, resource string, action string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(e echo.Context) error {
			userVal := e.Get("user")
			userModel, ok := userVal.(*models.User)
			if !ok || userModel == nil {
				return api.WebResponse(e, http.StatusUnauthorized, api.FIELD_VALIDATION_ERROR("User not found in context"))
			}
			user := userModel.UUID.String()
			domain := e.Request().Header.Get("domain")

			// Enforce used nomenclature (sub, dom, obj, act)
			ok, err := server.Casbin.Enforce(user, domain, resource, action)
			if err != nil {
				return api.WebResponse(e, http.StatusInternalServerError, api.FIELD_VALIDATION_ERROR("Casbin enforcement error"))
			}
			if !ok {
				return api.WebResponse(e, http.StatusForbidden, api.CASBIN_UNAUTHORIZED())
			}
			return next(e)
		}
	}
}
