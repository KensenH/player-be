package player

import (
	"context"
	"encoding/json"

	"github.com/pkg/errors"
	"gorm.io/gorm"

	inerr "player-be/internal/entity/errors"
	e "player-be/internal/entity/player"
)

// register new player to db
func (d *PlayerData) AddNewPlayer(ctx context.Context, newUser e.Player) (e.PlayerIdentity, error) {
	var (
		err  error
		resp e.PlayerIdentity
	)

	result := d.DB.Create(&newUser)
	if result.Error != nil {
		return resp, errors.Wrap(result.Error, "[DATA][AddNewPlayer] ")
	}

	//convert to e.Player to e.PlayerSignUpSuccess struct
	playerJson, err := json.Marshal(newUser)
	if err != nil {
		return resp, errors.Wrap(err, "[DATA][AddNewPlayer] ")
	}

	err = json.Unmarshal(playerJson, &resp)
	if err != nil {
		return resp, errors.Wrap(err, "[DATA][AddNewPlayer] ")
	}

	return resp, err
}

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
		player = &e.Player{
			ID:          bankAcc.PlayerID,
			BankAccount: bankAcc,
		}
	)

	result := d.DB.Save(&player)
	if result.Error != nil {
		return errors.Wrap(result.Error, "[Data] AddBankAccount")
	}

	return err
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
	if result.Error != nil {
		return false
	}

	return true
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
	if result.Error != nil {
		return false
	}

	return true
}
