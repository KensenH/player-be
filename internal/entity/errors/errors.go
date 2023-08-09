package errors

import "errors"

var (
	ErrIncorrectUsernamePassword = errors.New("username/password is incorrect")
	ErrMarshal                   = errors.New("error while mashaling data")
	ErrUnMarshal                 = errors.New("error while unmashaling data")
	ErrUsernameIsTaken           = errors.New("username is already taken")
	ErrEmailIsTaken              = errors.New("email is already used by another account")
	ErrPlayerNotFound            = errors.New("player not found")
	ErrInsufficientInGameMoney   = errors.New("insufficient ingame money")
)
