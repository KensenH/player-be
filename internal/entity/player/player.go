package player

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

type Player struct {
	ID             uint           `json:"player_id" gorm:"primaryKey"`
	Username       string         `json:"username" gorm:"type:varchar(15);uniqueIndex;not null"`
	Password       string         `json:"password" gorm:"type:varchar(255);not null"`
	FirstName      string         `json:"first_name" gorm:"type:varchar(255);not null"`
	LastName       string         `json:"last_name" gorm:"type:varchar(255)"`
	Email          string         `json:"email" gorm:"uniqueIndex;not null"`
	InGameCurrency int64          `json:"in_game_currency" gorm:"default:0"`
	SignUp         time.Time      `json:"join_date" gorm:"not null;default:now()"`
	BankAccount    BankAccount    `json:"bank_accounts" gorm:"foreignKey:PlayerID;references:ID"`
	TopUpHistory   []TopUpHistory `json:"top_up_history" gorm:"foreignKey:PlayerID;references:ID"`
}

type PlayerDetail struct {
	ID             uint   `json:"player_id"`
	Username       string `json:"username"`
	FirstName      string `json:"first_name"`
	LastName       string `json:"last_name"`
	Email          string `json:"email"`
	InGameCurrency int64  `json:"in_game_currency"`
}

type BankAccount struct {
	gorm.Model
	PlayerID         uint   `json:"player_id" validate:"required" gorm:"uniqueIndex"`
	BankName         string `json:"bank_name" validate:"required"`
	AccountOwnerName string `json:"account_owner_name" validate:"required"`
	AccountNumber    int64  `json:"account_number" validate:"required"`
}

type TopUpHistory struct {
	gorm.Model
	PlayerID       uint  `json:"player_id" gorm:"column:player_id"`
	InGameCurrency int64 `json:"ingame_cucrrency"`
	Price          int64 `json:"price"`
}

type TopUpRequest struct {
	TopUpAmount int64 `json:"top_up_amount"`
}

type PlayerSignUpForm struct {
	Username    string `json:"username" validate:"required"`
	Password    string `json:"password" validate:"required"`
	FirstName   string `json:"first_name" validate:"required,alpha"`
	LastName    string `json:"last_name" validate:"alpha"`
	PhoneNumber string `json:"phone_number" validate:"required,e164"`
	Email       string `json:"email" validate:"required,email"`
}

type PlayerIdentity struct {
	PlayerID uint      `json:"player_id"`
	Username string    `json:"username"`
	CreateAt time.Time `json:"created_at,omitempty"`
}

type PlayerUserPass struct {
	ID       uint   `json:"player_id"`
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type JwtClaims struct {
	PlayerID uint   `json:"player_id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

type PlayerID struct {
	ID uint
}
