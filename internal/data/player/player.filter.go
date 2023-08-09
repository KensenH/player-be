package player

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

// scopes
func scopePlayerId(playerId uint) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("id = ?", playerId)
	}
}

func scopeMinInGameCurrency(min int64) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("in_game_currency >= ?", min)
	}
}

func scopeMaxInGameCurrency(max int64) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("in_game_currency <= ?", max)
	}
}

func scopeUsernameLike(input string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("username ILIKE (?)", fmt.Sprintf("%%%s%%", input))
	}
}

func scopeJoinAfter(date time.Time) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("sign_up >= ?", date)
	}
}

func scopeJoinBefore(date time.Time) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("sign_up <= ?", date)
	}
}

func scopeBankNameLike(input string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("bank_name ILIKE (?)", fmt.Sprintf("%%%s%%", input))
	}
}

func scopeBankAccountNameLike(input string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("bank_account_name ILIKE (?)", fmt.Sprintf("%%%s%%", input))
	}
}

func scopeBankAccuntNumberLike(input int64) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("bank_name = ?", input)
	}
}