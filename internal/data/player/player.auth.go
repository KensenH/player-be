package player

import (
	"context"
	"fmt"
	"time"

	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	inerr "player-be/internal/entity/errors"
	e "player-be/internal/entity/player"
)

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
	var (
		err error
	)

	//store invalid token to redis
	err = d.Redis.Set(ctx, fmt.Sprintf("player:token:expired:%s", tokenID), expiredTime.String(), expiredTime.Sub(time.Now())).Err()
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
