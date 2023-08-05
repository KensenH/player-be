package boot

import (
	config "player-be/internal/config"

	playerData "player-be/internal/data/player"
	playerHandler "player-be/internal/delivery/http/player"
	playerService "player-be/internal/service/player"

	s "player-be/internal/delivery/http"

	"github.com/pkg/errors"
	// "github.com/redis/go-redis/v9"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func HTTP() error {
	var err error
	cfg := config.New()

	// rdb := redis.NewClient(&redis.Options{
	// 	Addr:     cfg.Database.Redis.Address,
	// 	Password: cfg.Database.Redis.Password,
	// 	DB:       0,
	// })

	db, err := gorm.Open(postgres.Open(cfg.Database.Postgres), &gorm.Config{})
	if err != nil {
		return errors.Wrap(err, "db open")
	}

	db.AutoMigrate()

	playerD := playerData.New(db)
	playerS := playerService.New(playerD)
	playerH := playerHandler.New(playerS)

	server := s.New(playerH)

	err = server.Serve()
	if err != nil {
		return errors.Wrap(err, "server serve")
	}

	return err
}
