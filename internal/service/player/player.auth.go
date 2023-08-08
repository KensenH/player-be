package player

import (
	"context"
	"encoding/json"
	"time"

	"github.com/pkg/errors"
	"gorm.io/gorm"

	inerr "player-be/internal/entity/errors"
	e "player-be/internal/entity/player"
	pwd "player-be/internal/pkg/password"
)

// validate form, and register new player to db
func (s *PlayerService) SignUp(ctx context.Context, playerForm e.PlayerSignUpForm) (e.PlayerIdentity, error) {
	var (
		err       error
		resp      e.PlayerIdentity
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

// signing in player by checking creds and creating new jwt token
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

	playerId := e.PlayerIdentity{
		PlayerID: stored.ID,
		Username: stored.Username,
	}

	tokenStr, err = s.JwtTool.CreateJWT(playerId, expirationTime)
	if err != nil {
		return tokenStr, errors.Wrap(err, "error while creating jwt token")
	}

	return tokenStr, err
}

// signing out player by adding token to blocklist in redis
func (s *PlayerService) SignOut(ctx context.Context, tokenStr string) error {
	var (
		err error
	)

	token, claims, err := s.JwtTool.ParseJWT(tokenStr)
	if err != nil {
		return errors.Wrap(err, "error while parsing jwt token")
	}

	expirationTime, err := token.Claims.GetExpirationTime()
	if err != nil {
		return errors.Wrap(err, "error while fetching jwt expire time")
	}

	valid, err := s.Data.TokenIsValid(ctx, claims.ID)
	if err != nil {
		return errors.Wrap(err, "error validating token")
	}

	if valid && token.Valid {
		err = s.Data.InvalidateToken(ctx, claims.ID, expirationTime.Time)
		if err != nil {
			return errors.Wrap(err, "error invalidating token")
		}
	}

	return err
}

// check jwt token still valid or not, by parsing and check blocklist from redis
func (s *PlayerService) JWTTokenValid(ctx context.Context, tokenStr string) (valid bool, playerId e.PlayerIdentity, err error) {
	token, claims, err := s.JwtTool.ParseJWT(tokenStr)
	if err != nil {
		return false, e.PlayerIdentity{}, errors.Wrap(err, "error while parsing jwt token")
	}

	valid, err = s.Data.TokenIsValid(ctx, claims.ID)
	if err != nil {
		return false, e.PlayerIdentity{}, errors.Wrap(err, "error validating token")
	}

	if valid && token.Valid {
		return true, e.PlayerIdentity{
			PlayerID: claims.PlayerID,
			Username: claims.Username,
		}, nil
	}

	return false, e.PlayerIdentity{}, err
}
