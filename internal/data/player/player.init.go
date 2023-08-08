package player

import (
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"

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
