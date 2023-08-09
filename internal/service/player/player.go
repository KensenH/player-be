package player

import (
	"context"
	inerr "player-be/internal/entity/errors"
	e "player-be/internal/entity/player"

	"github.com/pkg/errors"
)

func (s *PlayerService) GetPlayerDetail(ctx context.Context, playerId uint) (e.PlayerDetail, error) {
	player, err := s.Data.GetPlayerDetail(ctx, playerId)
	if err != nil {
		if errors.Is(err, inerr.ErrPlayerNotFound) {
			return player, err
		}
		return player, errors.Wrap(err, "[Service]GetPlayerDetail")
	}

	return player, err
}

func (s *PlayerService) AddBankAccount(ctx context.Context, bankAcc e.BankAccount) error {
	err := s.Data.AddBankAccount(ctx, bankAcc)
	if err != nil {
		return errors.Wrap(err, "[Service]AddBankAccount")
	}

	return err
}

func (s *PlayerService) TopUp(ctx context.Context, playerId uint, sum int64) (e.TopUpHistory, error) {
	//create virtual account or anything

	//verify payment

	//add game currency to player's account
	err := s.Data.AddInGameCurrency(ctx, playerId, sum)
	if err != nil {
		if errors.Is(err, inerr.ErrPlayerNotFound) {
			return e.TopUpHistory{}, err
		}

		return e.TopUpHistory{}, errors.Wrap(err, "[Service]TopUp")
	}

	topUp := e.TopUpHistory{
		PlayerID:       playerId,
		InGameCurrency: sum,
		Price:          sum,
	}

	//create transaction history
	err = s.Data.InputTopUpHistory(ctx, &topUp)
	if err != nil {
		return e.TopUpHistory{}, errors.Wrap(err, "[Service]TopUp")
	}

	return topUp, err
}
