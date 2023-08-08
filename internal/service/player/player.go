package player

import (
	"context"
	inerr "player-be/internal/entity/errors"
	e "player-be/internal/entity/player"

	"github.com/pkg/errors"
)

func (s *PlayerService) GetPlayerDetail(ctx context.Context, playerId uint) (e.PlayerDetail, error) {
	var (
		err    error
		player e.PlayerDetail
	)

	player, err = s.Data.GetPlayerDetail(ctx, playerId)
	if err != nil {
		if errors.Is(err, inerr.ErrPlayerNotFound) {
			return player, err
		}
		return player, errors.Wrap(err, "[Service]GetPlayerDetail")
	}

	return player, err
}

func (s *PlayerService) AddBankAccount(ctx context.Context, bankAcc e.BankAccount) error {
	var (
		err error
	)

	err = s.Data.AddBankAccount(ctx, bankAcc)
	if err != nil {
		return errors.Wrap(err, "[Service]AddBankAccount")
	}

	return err
}
