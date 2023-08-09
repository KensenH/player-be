package player

import (
	"errors"
	"strconv"

	inerr "player-be/internal/entity/errors"
	e "player-be/internal/entity/player"
	resp "player-be/internal/entity/response"

	"github.com/labstack/echo/v4"
)

// get player detail
func (h PlayerHandler) GetPlayerDetail(c echo.Context) error {
	var (
		err    error
		player e.PlayerDetail
	)

	idStr := c.Param("id")

	playerId, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		return echo.NewHTTPError(500, err.Error())
	}

	player, err = h.Service.GetPlayerDetail(c.Request().Context(), uint(playerId))
	if err != nil {
		if errors.Is(err, inerr.ErrPlayerNotFound) {
			return c.JSON(200, map[string]string{"message": err.Error()})
		}
		return echo.NewHTTPError(500, err.Error())
	}

	return c.JSON(200, player)
}

// get player's profile
func (h PlayerHandler) GetProfile(c echo.Context) error {
	var (
		err      error
		player   e.PlayerDetail
		playerId = c.Get("playerId").(e.PlayerIdentity)
	)

	player, err = h.Service.GetPlayerDetail(c.Request().Context(), uint(playerId.PlayerID))
	if err != nil {
		if errors.Is(err, inerr.ErrPlayerNotFound) {
			return c.JSON(200, map[string]string{"message": err.Error()})
		}
		return echo.NewHTTPError(500, err.Error())
	}

	return c.JSON(200, player)
}

// add or update player's bank account
func (h PlayerHandler) AddBankAccount(c echo.Context) error {
	var (
		err      error
		bankAcc  e.BankAccount
		playerId = c.Get("playerId").(e.PlayerIdentity)
	)

	if playerId.Username == "" {
		return echo.NewHTTPError(500, "username kosong")
	}

	err = c.Bind(&bankAcc)
	if err != nil {
		return err
	}

	bankAcc.PlayerID = playerId.PlayerID

	err = h.validator.Struct(&bankAcc)
	if err != nil {
		return err
	}

	err = h.Service.AddBankAccount(c.Request().Context(), bankAcc)
	if err != nil {
		return echo.NewHTTPError(500, err.Error())
	}

	return c.JSON(200, map[string]string{"message": "success adding bank account"})
}

// top up / buy ingame currency
func (h PlayerHandler) TopUp(c echo.Context) error {
	var (
		err      error
		playerId = c.Get("playerId").(e.PlayerIdentity)
		topUp    e.TopUpRequest
	)

	err = c.Bind(&topUp)
	if err != nil {
		return echo.NewHTTPError(500, err.Error())
	}

	receipt, err := h.Service.TopUp(c.Request().Context(), playerId.PlayerID, topUp.TopUpAmount)
	if err != nil {
		return echo.NewHTTPError(500, err.Error())
	}

	return c.JSON(200, map[string]interface{}{
		"receipt": receipt,
	})
}

// search player
func (h PlayerHandler) SearchPlayer(c echo.Context) error {
	var (
		err     error
		filter  e.PlayerFilter
		players []e.PlayerDetail
	)
	err = c.Bind(&filter)
	if err != nil {
		return echo.NewHTTPError(401)
	}

	err = h.validator.Struct(&filter)
	if err != nil {
		return err
	}

	players, err = h.Service.SearchPlayer(c.Request().Context(), filter)
	if err != nil {
		return echo.NewHTTPError(500, err.Error())
	}

	return c.JSON(200, resp.Response{
		Data: players,
	})
}

// show player's topup histories
func (h PlayerHandler) Receipts(c echo.Context) error {
	var (
		err       error
		playerId  = c.Get("playerId").(e.PlayerIdentity)
		histories []e.TopUpHistory
	)

	histories, err = h.Service.GetTopUpHistory(c.Request().Context(), playerId.PlayerID)
	if err != nil {
		return echo.NewHTTPError(500, err.Error())
	}

	return c.JSON(200, resp.Response{
		Data: histories,
	})
}
