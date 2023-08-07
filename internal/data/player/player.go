package player

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	inerr "player-be/internal/entity/errors"
	e "player-be/internal/entity/player"

	log "github.com/sirupsen/logrus"
)

type Option func(*PlayerData)

type PlayerData struct {
	DB    *gorm.DB
	Redis *redis.Client
}

func New(db *gorm.DB, redis *redis.Client, opts ...Option) *PlayerData {
	playerData := &PlayerData{
		DB:    db,
		Redis: redis,
	}

	for _, opt := range opts {
		opt(playerData)
	}

	err := db.AutoMigrate(&e.Player{}, &e.BankAccount{}, &e.TopUpHistory{})
	if err != nil {
		log.Fatalf("[Player Data] %s", err.Error())
	}

	return playerData
}

// register new player to db
func (d *PlayerData) AddNewPlayer(ctx context.Context, newUser e.Player) (e.PlayerSignUpSuccess, error) {
	var (
		err  error
		resp e.PlayerSignUpSuccess
	)

	result := d.DB.Create(&newUser)
	if result.Error != nil {
		return resp, errors.Wrap(err, "[DATA][AddNewPlayer] ")
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

// username exist in db
func (d *PlayerData) UsernameExist(ctx context.Context, username string) bool {
	var (
		player = e.Player{
			Username: username,
		}

		dbResult e.PlayerID
	)

	result := d.DB.Model(&player).Where(&player).First(&dbResult)
	if result.Error != nil {
		return false
	}

	return true
}

// email already registered in db
func (d *PlayerData) EmailRegistered(ctx context.Context, email string) bool {
	var (
		player = e.Player{
			Email: email,
		}

		dbResult e.PlayerID
	)

	result := d.DB.Model(&player).Where(&player).First(&dbResult)
	if result.Error != nil {
		return false
	}

	return true
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
