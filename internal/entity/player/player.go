package player

import (
	"time"

	"gorm.io/gorm"
)

type Player struct {
	ID             uint           `gorm:"primaryKey"`
	Username       string         `gorm:"type:varchar(15);uniqueIndex;not null"`
	Password       string         `gorm:"type:varchar(255);not null"`
	FirstName      string         `gorm:"type:varchar(255);not null"`
	LastName       string         `gorm:"type:varchar(255)"`
	Email          string         `gorm:"uniqueIndex;not null"`
	InGameCurrency int64          `gorm:"default:0"`
	SignUp         time.Time      `gorm:"not null"`
	BankAccount    []BankAccount  `gorm:"foreignKey:PlayerID;references:ID"`
	TopUpHistory   []TopUpHistory `gorm:"foreignKey:PlayerID;references:ID"`
}

type BankAccount struct {
	gorm.Model
	PlayerID      uint   `json:"player_ids"`
	BankName      string `json:"bank_name"`
	AccountNumber string `json:"account_number"`
}

type TopUpHistory struct {
	gorm.Model
	PlayerID       uint  `json:"player_id" gorm:"column:player_id"`
	InGameCurrency int64 `json:"ingame_cucrrency"`
	Price          int64 `json:"price"`
}

type PlayerSignUpForm struct {
	Username  string `json:"username" validate:"required"`
	Password  string `json:"password" validate:"required"`
	FirstName string `json:"first_name" validate:"required"`
	LastName  string `json:"last_name"`
	Email     string `json:"email" validate:"required,email"`
}

type PlayerSignUpSuccess struct {
	PlayerID uint      `json:"player_id"`
	Username string    `json:"username"`
	CreateAt time.Time `json:"created_at"`
}
