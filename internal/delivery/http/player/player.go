package player

import (
	"time"

	"github.com/labstack/echo/v4"

	playerEntity "player-be/internal/entity/player"
)

type PlayerHandler struct {
	Service PlayerService
}

type PlayerService interface{}

func New(playerService PlayerService) PlayerHandler {
	playerHandler := PlayerHandler{
		Service: playerService,
	}

	return playerHandler
}

func (h PlayerHandler) SignUp(c echo.Context) error {
	var (
		err  error
		form playerEntity.PlayerSignUpForm
	)

	err = c.Bind(&form)
	if err != nil {
		return echo.NewHTTPError(501, err.Error())
	}

	return c.JSON(200,
		playerEntity.PlayerSignUpSuccess{
			PlayerID: 0,
			Username: form.Username,
			CreateAt: time.Now(),
		})
}
