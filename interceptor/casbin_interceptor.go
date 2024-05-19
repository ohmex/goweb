package interceptor

import (
	"goweb/api"
	"goweb/services"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

func CasbinAuthorizer(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
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

		// ok, err := server.Enforce(User{ID: userID}, obj, method)
		// if err != nil {
		// 	return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		// }

		user := c.Get("user").(*jwt.Token)
		claims := user.Claims.(*services.JwtCustomClaims)

		log.Info().Interface("Claims", claims).Send()

		ok := false
		if ok {
			return next(c)
		}

		return api.WebResponse(c, http.StatusForbidden, api.CASBIN_UNAUTHORIZED())
	}
}
