package player

import (
	"github.com/labstack/echo/v4"

	e "player-be/internal/entity/player"
)

func (h PlayerHandler) JwtMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		var (
			err      error
			ok       bool
			playerId e.PlayerIdentity

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

		// if valid create custom context so the service layer know who is accessing
		if ok, playerId, err = h.Service.JWTTokenValid(c.Request().Context(), tokenStr); ok {
			c.Set("playerId", playerId)
			return next(c)
		}
		if err != nil {
			return echo.NewHTTPError(501, err.Error())
		}

		return echo.NewHTTPError(401, "Unauthorized")
	}
}
