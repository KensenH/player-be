package player

import (
	"context"
	"fmt"
	"time"

	"encoding/json"

	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	inerr "player-be/internal/entity/errors"
	e "player-be/internal/entity/player"
)

// register new player to db
func (d *PlayerData) AddNewPlayer(ctx context.Context, newUser e.Player) (e.PlayerIdentity, error) {
	var (
		err  error
		resp e.PlayerIdentity
	)

	result := d.DB.Create(&newUser)
	if result.Error != nil {
		return resp, errors.Wrap(result.Error, "[DATA][AddNewPlayer] ")
	}

	//convert to e.Player to e.PlayerSignUpSuccess struct
	playerJson, err := json.Marshal(newUser)
	if err != nil {
		return resp, errors.Wrap(err, "[DATA][AddNewPlayer] ")
	}

	err = json.Unmarshal(playerJson, &resp)
	if err != nil {
		return resp, errors.Wrap(err, "[DATA][AddNewPlayer] ")
	}

	return resp, err
}

// get user's hashed password from db
func (d *PlayerData) GetHashedPassword(ctx context.Context, username string) (e.PlayerUserPass, error) {
	var (
		err  error
		resp e.PlayerUserPass

		player = e.Player{
			Username: username,
		}
	)

	result := d.DB.Model(&player).Where(&player).First(&resp)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return resp, inerr.ErrIncorrectUsernamePassword
		}
		return resp, errors.Wrap(result.Error, "[DATA][GetHashedPassword]")
	}

	return resp, err
}

// store invalid tokenID to redis
func (d *PlayerData) InvalidateToken(ctx context.Context, tokenID string, expiredTime time.Time) error {
	//store invalid token to redis
	err := d.Redis.Set(ctx, fmt.Sprintf("player:token:expired:%s", tokenID), expiredTime.String(), time.Until(expiredTime)).Err()
	if err != nil {
		return errors.Wrap(err, "error while registering invalid token")
	}

	return err
}

// validate wether token is valid or not from redis
func (d *PlayerData) TokenIsValid(ctx context.Context, tokenID string) (bool, error) {
	var (
		err error
	)

	key := fmt.Sprintf("player:token:expired:%s", tokenID)

	//check if token is black listed in redis
	_, err = d.Redis.Get(ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return true, nil
		}
		return false, errors.Wrap(err, "error while fetching data from redis")
	}

	return false, err
}
