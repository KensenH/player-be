package player

import (
	"context"
	"time"

	"github.com/go-playground/validator/v10"

	e "player-be/internal/entity/player"
)

type Option func(h *PlayerHandler)

type PlayerHandler struct {
	Service   PlayerService
	validator *validator.Validate
}

type PlayerService interface {
	SignUp(ctx context.Context, playerForm e.PlayerSignUpForm) (e.PlayerIdentity, error)
	SignIn(ctx context.Context, expirationTime time.Time, playerForm e.PlayerUserPass) (tokenStr string, err error)
	SignOut(ctx context.Context, tokenStr string) error

	JWTTokenValid(ctx context.Context, tokenStr string) (bool, e.PlayerIdentity, error)

	GetPlayerDetail(ctx context.Context, playerId uint) (e.PlayerDetail, error)

	AddBankAccount(ctx context.Context, bankAcc e.BankAccount) error

	TopUp(ctx context.Context, playerId uint, sum int64) (e.TopUpHistory, error)
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
