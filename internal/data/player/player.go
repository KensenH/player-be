package player

import (
	"context"

	"gorm.io/gorm"

	e "player-be/internal/entity/player"

	log "github.com/sirupsen/logrus"
)

type Option func(*PlayerData)

type PlayerData struct {
	DB *gorm.DB `validate:"required"`
}

func New(db *gorm.DB, opts ...Option) PlayerData {
	playerData := PlayerData{
		DB: db,
	}

	for _, opt := range opts {
		opt(&playerData)
	}

	err := db.AutoMigrate(&e.Player{}, &e.BankAccount{}, &e.TopUpHistory{})
	if err != nil {
		log.Fatalf("[Player Data] %s", err.Error())
	}

	return playerData
}

// register new player to db
func (d *PlayerData) AddNewPlayer(ctx context.Context) error {
	var err error

	return err
}

// username exist in db
func (d *PlayerData) UsernameExist(ctx context.Context) (bool, error) {
	var (
		err   error
		exist bool
	)

	return exist, err
}

// email already registered in db
func (d *PlayerData) EmailRegistered(ctx context.Context) (bool, error) {
	var (
		err        error
		registered bool
	)

	return registered, err
}
