package player

import (
	"context"
	"encoding/json"

	"github.com/pkg/errors"
	"gorm.io/gorm"

	inerr "player-be/internal/entity/errors"
	e "player-be/internal/entity/player"

	log "github.com/sirupsen/logrus"
)

type Option func(*PlayerData)

type PlayerData struct {
	DB *gorm.DB `validate:"required"`
}

func New(db *gorm.DB, opts ...Option) *PlayerData {
	playerData := &PlayerData{
		DB: db,
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
