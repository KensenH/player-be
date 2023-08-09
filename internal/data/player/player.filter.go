package player

import (
	"time"

	"gorm.io/gorm"
)

// scopes
func (d *PlayerData) PlayerId(playerId uint) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("id = ?", playerId)
	}
}

func (d *PlayerData) MinInGameCurrency(min int64) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("in_game_currency >= ?", min)
	}
}

func (d *PlayerData) MaxInGameCurrency(max int64) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("in_game_currency <= ?", max)
	}
}

func (d *PlayerData) UsernameLike(input string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("username ILIKE ('%?%')", input)
	}
}

func (d *PlayerData) JoinAfter(date time.Time) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("sign_up >= ?", date)
	}
}

func JoinBefore(date time.Time) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("sign_up <= ?", date)
	}
}
