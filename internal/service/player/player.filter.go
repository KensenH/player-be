package player

import (
	"context"
	e "player-be/internal/entity/player"

	"github.com/pkg/errors"
)

func (s *PlayerService) SearchPlayer(ctx context.Context, filter e.PlayerFilter) ([]e.Player, error) {
	var (
		err     error
		players []e.Player
		scopes  []Scope
	)

	if filter.PlayerId > 0 {
		scopes = append(scopes, s.Data.PlayerId(filter.PlayerId))
	}

	players, err = s.Data.SearchPlayer(ctx, scopes)
	if err != nil {
		return players, errors.Wrap(err, "[Service]SearchPlayer")
	}

	return players, err
}
