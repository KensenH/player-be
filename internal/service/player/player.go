package player

import (
	"context"
	"encoding/json"
	e "player-be/internal/entity/player"
	pwd "player-be/internal/pkg/password"

	"github.com/pkg/errors"
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
		return resp, errors.New("username is already taken")
	}

	//check email
	if exist := s.Data.EmailRegistered(ctx, playerForm.Email); exist {
		return resp, errors.New("email is already used by another account")
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
		return resp, errors.Wrap(err, "error while mashaling data")
	}

	err = json.Unmarshal(playerFormJson, &newPlayer)
	if err != nil {
		return resp, errors.Wrap(err, "error while unmashaling data")
	}

	//add new player to db
	resp, err = s.Data.AddNewPlayer(ctx, newPlayer)
	if err != nil {
		return resp, errors.Wrap(err, "error while signing up new player")
	}

	return resp, err
}
