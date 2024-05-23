package interceptor

import (
	"goweb/api"
	"goweb/models"
	"goweb/server"
	"goweb/util"
	"net/http"

	"github.com/labstack/echo/v4"
)

func CasbinAuthorizer(server *server.Server) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Get tenant from header, LATER this must be UUID
			tenant := c.Request().Header.Get("tenant")

			// Get user name, LATER this must UUID
			user := c.Get("user").(*models.User).Name

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
