package player

import (
	"context"
	"encoding/json"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/pkg/errors"
	"gorm.io/gorm"

	inerr "player-be/internal/entity/errors"
	e "player-be/internal/entity/player"
	pwd "player-be/internal/pkg/password"
)

type Option func(*PlayerService)

type PlayerService struct {
	Data PlayerData
}

type PlayerData interface {
	AddNewPlayer(ctx context.Context, newUser e.Player) (e.PlayerSignUpSuccess, error)
	UsernameExist(ctx context.Context, username string) bool
	EmailRegistered(ctx context.Context, email string) bool
	GetHashedPassword(ctx context.Context, username string) (e.PlayerUserPass, error)
	InvalidateToken(ctx context.Context, tokenID string, expiredTime time.Time) error
	TokenIsValid(ctx context.Context, tokenID string) (bool, error)
}

func New(playerData PlayerData, opts ...Option) *PlayerService {
	playerService := &PlayerService{
		Data: playerData,
	}

	for _, opt := range opts {
		opt(playerService)
	}

	return playerService
}

// validate form, and register new player to db
func (s *PlayerService) SignUp(ctx context.Context, playerForm e.PlayerSignUpForm) (e.PlayerSignUpSuccess, error) {
	var (
		err       error
		resp      e.PlayerSignUpSuccess
		newPlayer e.Player
	)

	//check username
	if exist := s.Data.UsernameExist(ctx, playerForm.Username); exist {
		return resp, inerr.ErrUsernameIsTaken
	}

	//check email
	if exist := s.Data.EmailRegistered(ctx, playerForm.Email); exist {
		return resp, inerr.ErrEmailIsTaken
	}

	//hash password
	hashed, err := pwd.HashPassword(playerForm.Password)
	if err != nil {
		return resp, errors.Wrap(err, "error while hashing password")
	}

	playerForm.Password = hashed

	//insert player
	//convert from form to player struct
	playerFormJson, err := json.Marshal(playerForm)
	if err != nil {
		return resp, inerr.ErrMarshal
	}

	err = json.Unmarshal(playerFormJson, &newPlayer)
	if err != nil {
		return resp, inerr.ErrUnMarshal
	}

	//add new player to db
	resp, err = s.Data.AddNewPlayer(ctx, newPlayer)
	if err != nil {
		return resp, errors.Wrap(err, "error while signing up new player")
	}

	return resp, err
}

func (s *PlayerService) SignIn(ctx context.Context, expirationTime time.Time, playerForm e.PlayerUserPass) (tokenStr string, err error) {
	stored, err := s.Data.GetHashedPassword(ctx, playerForm.Username)
	if err != nil {
		if errors.Unwrap(err) == gorm.ErrRecordNotFound {
			return tokenStr, inerr.ErrIncorrectUsernamePassword
		}
		return tokenStr, errors.Wrap(err, "error fetching player data")
	}

	if !pwd.CheckPasswordHash(playerForm.Password, stored.Password) && stored.Username == playerForm.Username {
		return tokenStr, inerr.ErrIncorrectUsernamePassword
	}

	claims := &e.JwtClaims{
		Username: stored.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenStr, err = token.SignedString([]byte("mysecretjwtkey"))
	if err != nil {
		return tokenStr, errors.Wrap(err, "error while creating jwt")
	}

	return tokenStr, err
}

func (s *PlayerService) SignOut(ctx context.Context, tokenStr string) error {
	var (
		err error
	)

	claims := e.JwtClaims{}
	token, err := jwt.ParseWithClaims(tokenStr, &claims, func(token *jwt.Token) (interface{}, error) {
		return "mysecretjwtkey", nil
	})
	if err != nil {
		return errors.Wrap(err, "error while parsing jwt token")
	}

	valid, err := s.Data.TokenIsValid(ctx, claims.ID)
	if err != nil {
		return errors.Wrap(err, "error validating token")
	}

	if valid && token.Valid {
		err = s.Data.InvalidateToken(ctx, claims.ID, claims.ExpiresAt.Time)
		if err != nil {
			return errors.Wrap(err, "error invalidating token")
		}
	}

	return err
}

func (s *PlayerService) JWTTokenValid(ctx context.Context, tokenStr string) (bool, error) {
	var (
		err    error
		claims = e.JwtClaims{}
	)

	token, err := jwt.ParseWithClaims(tokenStr, &claims, func(token *jwt.Token) (interface{}, error) {
		return "mysecretjwtkey", nil
	})
	if err != nil {
		return false, errors.Wrap(err, "error while parsing jwt token")
	}

	valid, err := s.Data.TokenIsValid(ctx, claims.ID)
	if err != nil {
		return false, errors.Wrap(err, "error validating token")
	}

	if valid && token.Valid {
		return true, err
	}

	return false, err
}
