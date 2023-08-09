package player

import (
	"context"

	"github.com/pkg/errors"
	"gorm.io/gorm"

	inerr "player-be/internal/entity/errors"
	e "player-be/internal/entity/player"
)

// get player detail
func (d *PlayerData) GetPlayerDetail(ctx context.Context, playerId uint) (e.PlayerDetail, error) {
	var (
		err    error
		player = e.Player{
			ID: playerId,
		}
		playerDetail e.PlayerDetail
	)

	result := d.DB.Model(&player).Where(&player).First(&playerDetail)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return playerDetail, inerr.ErrPlayerNotFound
		}
		return playerDetail, errors.Wrap(result.Error, "[DATA][GetPlayerDetail] ")
	}

	return playerDetail, err
}

// add or update player's bank account
func (d *PlayerData) AddBankAccount(ctx context.Context, bankAcc e.BankAccount) error {
	var (
		err    error
		player = e.Player{
			ID:                bankAcc.PlayerID,
			BankName:          bankAcc.BankName,
			BankAccountName:   bankAcc.AccountOwnerName,
			BankAccountNumber: bankAcc.AccountNumber,
		}
	)

	d.DB.Model(&player).Updates(&player)

	return err
}

// add in-game currency
func (d *PlayerData) AddInGameCurrency(ctx context.Context, playerId uint, sum int64) error {
	var (
		err    error
		player = &e.Player{
			ID: playerId,
		}
	)

	err = d.DB.First(&player).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return inerr.ErrPlayerNotFound
		}
		return errors.Wrap(err, "[Data]AddInGameCurrency")
	}

	player.InGameCurrency += sum

	err = d.DB.Save(&player).Error
	if err != nil {
		return errors.Wrap(err, "[Data]AddInGameCurrency")
	}

	return err
}

// subtract in-game currency
func (d *PlayerData) SubInGameCurrency(ctx context.Context, playerId uint, sum int64) error {
	var (
		err    error
		player = &e.Player{
			ID: playerId,
		}
	)

	err = d.DB.First(&player).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return inerr.ErrPlayerNotFound
		}
		return errors.Wrap(err, "[Data]AddInGameCurrency")
	}

	if player.InGameCurrency < sum {
		return inerr.ErrInsufficientInGameMoney
	}

	player.InGameCurrency -= sum

	err = d.DB.Save(&player).Error
	if err != nil {
		return errors.Wrap(err, "[Data]AddInGameCurrency")
	}

	return err
}

// input receipt to db
func (d *PlayerData) InputTopUpHistory(ctx context.Context, topUp *e.TopUpHistory) error {
	result := d.DB.Create(&topUp)
	if result.Error != nil {
		return errors.Wrap(result.Error, "[Data]CreateTopUpHistory")
	}

	return nil
}

// get top up histories by player id
func (d *PlayerData) GetTopUpHistory(ctx context.Context, playerId uint) ([]e.TopUpHistory, error) {
	var (
		err     error
		history = e.TopUpHistory{
			PlayerID: playerId,
		}
		histories []e.TopUpHistory
	)

	result := d.DB.Where(&history).Find(&histories)
	if result.Error != nil {
		return histories, errors.Wrap(err, "[Data]GetTopUpHistory")
	}

	return histories, err
}

// search player
func (d *PlayerData) SearchPlayer(ctx context.Context, filter e.PlayerFilterFeed) ([]e.PlayerDetail, error) {
	var (
		err     error
		player  e.Player
		players []e.PlayerDetail
		scopes  []func(db *gorm.DB) *gorm.DB
	)

	if filter.PlayerId > 0 {
		scopes = append(scopes, scopePlayerId(filter.PlayerId))
	}

	if filter.UsernameLike != "" {
		scopes = append(scopes, scopeUsernameLike(filter.UsernameLike))
	}

	if !filter.JoinAfter.IsZero() {
		scopes = append(scopes, scopeJoinAfter(filter.JoinAfter))
	}

	if !filter.JoinBefore.IsZero() {
		scopes = append(scopes, scopeJoinBefore(filter.JoinBefore))
	}

	if filter.MinInGameCurrency > 0 {
		scopes = append(scopes, scopeMinInGameCurrency(filter.MinInGameCurrency))
	}

	if filter.MaxInGameCurrency > 0 {
		scopes = append(scopes, scopeMaxInGameCurrency(filter.MaxInGameCurrency))
	}

	if filter.BankName != "" {
		scopes = append(scopes, scopeBankNameLike(filter.BankName))
	}

	if filter.BankAccountName != "" {
		scopes = append(scopes, scopeBankAccountNameLike(filter.BankAccountName))
	}

	if filter.BankAccountNumber > 0 {
		scopes = append(scopes, scopeBankAccuntNumberLike(filter.BankAccountNumber))
	}

	result := d.DB.Model(&player).Scopes(scopes...).Find(&players)
	if result.Error != nil {
		return players, errors.Wrap(err, "[Data]SearchPlayer")
	}

	return players, err
}

// username exist in db
func (d *PlayerData) UsernameExist(ctx context.Context, username string) bool {
	var (
		player = e.Player{
			Username: username,
		}

		dbResult e.PlayerID
	)

	result := d.DB.Model(&player).Where(&player).First(&dbResult)

	return result.Error == nil
}

// email already registered in db
func (d *PlayerData) EmailRegistered(ctx context.Context, email string) bool {
	var (
		player = e.Player{
			Email: email,
		}

		dbResult e.PlayerID
	)

	result := d.DB.Model(&player).Where(&player).First(&dbResult)

	return result.Error == nil
}
