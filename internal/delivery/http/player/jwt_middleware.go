package player

import (
	"github.com/labstack/echo/v4"
)

func (h PlayerHandler) JwtMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		var (
			err error
			ok  bool

			//skip path below
			skip = []string{
				"/player-be/api/v1/player/signup",
				"/player-be/api/v1/player/signin",
				"/player-be/api/v1/player/signout",
				"/",
			}
		)

		for _, path := range skip {
			if path == c.Request().URL.Path {
				return next(c)
			}
		}

		//take cookie from client
		tokenCookie, err := c.Cookie("token")
		if err != nil {
			return echo.NewHTTPError(401, "Unauthorized")
		}

		tokenStr := tokenCookie.Value

		if ok, err = h.Service.JWTTokenValid(c.Request().Context(), tokenStr); ok {
			return next(c)
		}
		if err != nil {
			echo.NewHTTPError(501, err.Error())
		}

		return echo.NewHTTPError(401, "Unauthorized")
	}
}
