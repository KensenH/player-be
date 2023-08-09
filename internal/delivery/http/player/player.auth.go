package player

import (
	"encoding/base64"
	"net/http"
	"strings"
	"time"

	"github.com/labstack/echo/v4"

	inerr "player-be/internal/entity/errors"
	playerEntity "player-be/internal/entity/player"

	"player-be/internal/entity/response"
)

// player sign up,
// api/v1/player/signup
func (h PlayerHandler) SignUp(c echo.Context) error {
	var (
		err  error
		form playerEntity.PlayerSignUpForm
		resp playerEntity.PlayerIdentity
	)

	err = c.Bind(&form)
	if err != nil {
		return echo.NewHTTPError(501, err.Error())
	}

	//validate form
	err = h.validator.Struct(&form)
	if err != nil {
		return err
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
		err error
	)

	authHeader := c.Request().Header.Get("Authorization")

	if len(authHeader) < 6 || authHeader[:6] != "Basic " {
		return echo.NewHTTPError(401, "Unauthorized")
	}

	decoded, err := base64.StdEncoding.DecodeString(authHeader[6:])
	if err != nil {
		return echo.NewHTTPError(401, "Unauthorized")
	}

	credential := strings.Split(string(decoded), ":")

	form := playerEntity.PlayerUserPass{
		Username: credential[0],
		Password: credential[1],
	}

	//validate form
	err = h.validator.Struct(form)
	if err != nil {
		return err
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
	)

	//get token
	cookie, err := c.Cookie("token")
	if err != nil {
		//no token found
		if strings.Contains(err.Error(), "named cookie not present") {
			return c.JSON(200, response.Response{Data: "OK"})
		}

		return echo.NewHTTPError(501, err)
	}

	//check to token from claims and redis
	err = h.Service.SignOut(c.Request().Context(), cookie.Value)
	if err != nil {
		return echo.NewHTTPError(501, err)
	}

	//delete token from client's side
	cookie = &http.Cookie{
		Name:   "token",
		Value:  "",
		MaxAge: -1,
	}
	c.SetCookie(cookie)

	return c.JSON(200, response.Response{Data: "OK"})
}
