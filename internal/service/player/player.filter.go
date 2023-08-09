package player

import (
	"context"
	e "player-be/internal/entity/player"
	"time"

	"github.com/pkg/errors"
)

// search player
func (s *PlayerService) SearchPlayer(ctx context.Context, filter e.PlayerFilter) ([]e.PlayerDetail, error) {
	var (
		err            error
		players        []e.PlayerDetail
		feedJoinAfter  time.Time
		feedJoinBefore time.Time
	)

	if filter.JoinAfter != "" {
		feedJoinAfter, err = time.Parse("02-01-2006", filter.JoinAfter)
	}

	if filter.JoinBefore != "" {
		feedJoinBefore, err = time.Parse("02-01-2006", filter.JoinBefore)
	}

	feed := e.PlayerFilterFeed{
		PlayerId:          filter.PlayerId,
		MinInGameCurrency: filter.MinInGameCurrency,
		MaxInGameCurrency: filter.MaxInGameCurrency,
		UsernameLike:      filter.UsernameLike,
		JoinAfter:         feedJoinAfter,
		JoinBefore:        feedJoinBefore,
		BankName:          filter.BankName,
		BankAccountName:   filter.BankAccountName,
		BankAccountNumber: filter.BankAccountNumber,
	}

	players, err = s.Data.SearchPlayer(ctx, feed)
	if err != nil {
		return players, errors.Wrap(err, "[Service]SearchPlayer")
	}

	return players, err
}
