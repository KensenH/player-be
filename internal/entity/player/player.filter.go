package player

import "time"

type PlayerFilter struct {
	PlayerId          uint   `query:"player_id" json:"player_id"`
	JoinBefore        string `query:"join_before" json:"join_before" validate:"omitempty,datetime=02-01-2006"`
	JoinAfter         string `query:"join_after" json:"join_after" validate:"omitempty,datetime=02-01-2006"`
	MinInGameCurrency int64  `query:"min_ingame_currency" json:"min_ingame_currency"`
	MaxInGameCurrency int64  `query:"max_ingame_currency" json:"max_ingame_currency"`
	UsernameLike      string `query:"username" json:"username"`
	BankName          string `query:"bank_name" json:"bank_name"`
	BankAccountName   string `query:"bank_account_name" json:"bank_account_name"`
	BankAccountNumber int64  `query:"bank_account_number" json:"bank_account_number"`
}

type PlayerFilterFeed struct {
	PlayerId          uint      `query:"player_id" json:"player_id"`
	JoinBefore        time.Time `query:"join_before" json:"join_before"`
	JoinAfter         time.Time `query:"join_after" json:"join_after"`
	MinInGameCurrency int64     `query:"min_ingame_currency"`
	MaxInGameCurrency int64     `query:"max_ingame_currency"`
	UsernameLike      string    `query:"username" json:"username"`
	BankName          string    `query:"bank_name" json:"bank_name"`
	BankAccountName   string    `query:"bank_account_name" json:"bank_account_name"`
	BankAccountNumber int64     `query:"bank_account_number" json:"bank_account_number"`
}
