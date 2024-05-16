package middleware

// Middleware for additional steps:
// 1. Check the user exists in DB
// 2. Check the token info exists in Redis
// 3. Add the user DB data to Context
// 4. Prolong the Redis TTL of the current token pair
/*
func ValidateJWT(server *s.Server) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			token := c.Get("user").(*jwtGo.Token)
			claims := token.Claims.(*tokenService.JwtCustomClaims)

			user, err := tokenService.NewTokenService(server).ValidateToken(claims, false)
			if err != nil {
				return responses.MessageResponse(c, http.StatusUnauthorized, "Not authorized")
			}

			c.Set("currentUser", user)

			go func() {
				server.Redis.Expire(fmt.Sprintf("token-%d", claims.ID),
					time.Minute*tokenService.AutoLogoffMinutes)
			}()

			return next(c)
		}
	}
}
*/
