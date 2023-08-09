package player

type PlayerFilter struct {
	PlayerId          uint   `query:"player_id" json:"player_id"`
	JoinBefore        string `query:"join_before" json:"join_before" validate:"datetime=02-01-2006"`
	JoinAfter         string `query:"join_after" json:"join_after" validate:"datetime=02-01-2006"`
	MinInGameCurrency int64  `query:"min_ingame_currency" json:"min_ingame_currency"`
	MaxInGameCurrency int64  `query:"max_ingame_currency" json:"max_ingame_currency"`
	UsernameLike      string `query:"username" json:"username"`
}
