package player

import (
	"context"

	"github.com/labstack/echo/v4"

	"github.com/go-playground/validator/v10"

	playerEntity "player-be/internal/entity/player"
)

type Option func(h *PlayerHandler)

type PlayerHandler struct {
	Service   PlayerService
	validator *validator.Validate
}

type PlayerService interface {
	SignUp(ctx context.Context, playerForm playerEntity.PlayerSignUpForm) (playerEntity.PlayerSignUpSuccess, error)
}

// new player handler
func New(playerService PlayerService, opts ...Option) *PlayerHandler {
	playerHandler := &PlayerHandler{
		Service: playerService,
	}

	for _, opt := range opts {
		opt(playerHandler)
	}

	if playerHandler.validator == nil {
		playerHandler.validator = validator.New()
	}

	return playerHandler
}

// add playground validator
func WithValidator(v *validator.Validate) Option {
	return func(h *PlayerHandler) {
		h.validator = v
	}
}

// player sign up,
// api/v1/player/signup
func (h PlayerHandler) SignUp(c echo.Context) error {
	var (
		err  error
		form playerEntity.PlayerSignUpForm
		resp playerEntity.PlayerSignUpSuccess
	)

	err = c.Bind(&form)
	if err != nil {
		return echo.NewHTTPError(501, err.Error())
	}

	err = h.validator.Struct(form)
	if err != nil {
		return echo.NewHTTPError(400, err.Error())
	}

	resp, err = h.Service.SignUp(c.Request().Context(), form)
	if err != nil {
		return echo.NewHTTPError(501, err.Error())
	}

	return c.JSON(200, resp)
}
