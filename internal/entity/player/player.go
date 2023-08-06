package player

import (
	"time"

	"gorm.io/gorm"

	"github.com/golang-jwt/jwt/v5"
)

type Player struct {
	gorm.Model
	ID             uint           `json:"player_id" gorm:"primaryKey"`
	Username       string         `json:"username" gorm:"type:varchar(15);uniqueIndex;not null"`
	Password       string         `json:"password" gorm:"type:varchar(255);not null"`
	FirstName      string         `json:"first_name" gorm:"type:varchar(255);not null"`
	LastName       string         `json:"last_name" gorm:"type:varchar(255)"`
	Email          string         `json:"email" gorm:"uniqueIndex;not null"`
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

type PlayerUserPass struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type JwtClaims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

type PlayerID struct {
	ID uint
}
