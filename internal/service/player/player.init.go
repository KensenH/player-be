package player

import (
	"context"
	"time"

	"github.com/golang-jwt/jwt/v5"

	e "player-be/internal/entity/player"
)

type Option func(*PlayerService)

type PlayerService struct {
	Data    PlayerData
	JwtTool JwtTool
}

type PlayerData interface {
	AddNewPlayer(ctx context.Context, newUser e.Player) (e.PlayerIdentity, error)

	UsernameExist(ctx context.Context, username string) bool
	EmailRegistered(ctx context.Context, email string) bool
	
	GetHashedPassword(ctx context.Context, username string) (e.PlayerUserPass, error)
	
	InvalidateToken(ctx context.Context, tokenID string, expiredTime time.Time) error
	TokenIsValid(ctx context.Context, tokenID string) (bool, error)
	
	GetPlayerDetail(ctx context.Context, playerId uint) (e.PlayerDetail, error)
	
	AddBankAccount(ctx context.Context, bankAcc e.BankAccount) error
}

type JwtTool interface {
	CreateJWT(playerId e.PlayerIdentity, expirationTime time.Time) (string, error)
	ParseJWT(tokenStr string) (token *jwt.Token, claims e.JwtClaims, err error)
}

func New(playerData PlayerData, jwtTool JwtTool, opts ...Option) *PlayerService {
	playerService := &PlayerService{
		Data:    playerData,
		JwtTool: jwtTool,
	}

	for _, opt := range opts {
		opt(playerService)
	}

	return playerService
}
