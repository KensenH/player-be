package player

import (
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
		// response resp.Response
	)

	err = c.Bind(&form)
	if err != nil {

	}

	return c.String(200, "OK")
}
