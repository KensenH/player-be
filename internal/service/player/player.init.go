package player

import (
	"context"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"

	e "player-be/internal/entity/player"
)

type Option func(*PlayerService)
type Scope []func(db *gorm.DB) *gorm.DB

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

	SearchPlayer(ctx context.Context, scopes []Scope) ([]e.Player, error)

	AddBankAccount(ctx context.Context, bankAcc e.BankAccount) error
	AddInGameCurrency(ctx context.Context, playerId uint, sum int64) error
	InputTopUpHistory(ctx context.Context, topUp *e.TopUpHistory) error

	//scopes
	PlayerId(playerId uint) Scope
	MinInGameCurrency(min int64) Scope
	MaxInGameCurrency(max int64) Scope
	UsernameLike(input string) Scope
	JoinAfter(date time.Time) Scope
	JoinBefore(date time.Time) Scope
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
