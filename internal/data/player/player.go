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
		err error
	)

	if d.DB.Model(&bankAcc).Where("player_id = ?", bankAcc.PlayerID).Updates(&bankAcc).RowsAffected == 0 {
		result := d.DB.Create(&bankAcc)
		if result.Error != nil {
			return errors.Wrap(err, "[Data]AddBankAccount")
		}
	}

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

// input transaction to history
func (d *PlayerData) InputTopUpHistory(ctx context.Context, topUp *e.TopUpHistory) error {
	result := d.DB.Create(&topUp)
	if result.Error != nil {
		return errors.Wrap(result.Error, "[Data]CreateTopUpHistory")
	}

	return nil
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
