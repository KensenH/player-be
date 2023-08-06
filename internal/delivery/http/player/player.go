package player

import (
	"context"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"

	"github.com/go-playground/validator/v10"

	inerr "player-be/internal/entity/errors"
	playerEntity "player-be/internal/entity/player"

	"player-be/internal/entity/response"
)

type Option func(h *PlayerHandler)

type PlayerHandler struct {
	Service   PlayerService
	validator *validator.Validate
}

type PlayerService interface {
	SignUp(ctx context.Context, playerForm playerEntity.PlayerSignUpForm) (playerEntity.PlayerSignUpSuccess, error)
	SignIn(ctx context.Context, expirationTime time.Time, playerForm playerEntity.PlayerUserPass) (tokenStr string, err error)
	SignOut(ctx context.Context, tokenStr string) error
	JWTTokenValid(ctx context.Context, tokenStr string) (bool, error)
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

	//validate form
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

// signin,
// api/v1/player/signin
func (h PlayerHandler) SignIn(c echo.Context) error {
	var (
		err  error
		form playerEntity.PlayerUserPass
	)

	//bind body
	err = c.Bind(&form)
	if err != nil {
		return echo.NewHTTPError(501, err.Error())
	}

	//validate form
	err = h.validator.Struct(form)
	if err != nil {
		return echo.NewHTTPError(400, err.Error())
	}

	//jwt TTL
	expirationTime := time.Now().Add(24 * time.Hour)

	tokenStr, err := h.Service.SignIn(c.Request().Context(), expirationTime, form)
	if err != nil {
		if err == inerr.ErrIncorrectUsernamePassword {
			return echo.NewHTTPError(401, err.Error())
		}

		return echo.NewHTTPError(501, err.Error())
	}

	//create cookie
	cookie := &http.Cookie{
		Name:    "token",
		Value:   tokenStr,
		Expires: expirationTime,
	}

	//set cookie
	c.SetCookie(cookie)

	return c.JSON(200, response.Response{
		Data: map[string]string{
			"jwt_token": tokenStr,
		},
	})
}

// signout,
// api/v1/player/signout
func (h PlayerHandler) SignOut(c echo.Context) error {
	var (
		err error
		// query authEntity.Auth
	)

	cookie, err := c.Cookie("token")
	if err != nil {
		return echo.NewHTTPError(501, err)
	}

	cookie = &http.Cookie{
		Name:   "token",
		Value:  "",
		MaxAge: -1,
	}
	c.SetCookie(cookie)

	return c.JSON(200, response.Response{
		Data: "OK",
	})
}
