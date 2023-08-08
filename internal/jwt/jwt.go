// package to help create and parse jwt token
// internal use only
package jwt

import (
	"player-be/internal/pkg/random"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/pkg/errors"

	e "player-be/internal/entity/player"
)

type jwtData struct {
	Secret string
}

type jwtPack interface {
	CreateJWT(playerId e.PlayerIdentity, expirationTime time.Time) (string, error)
	ParseJWT(tokenStr string) (token *jwt.Token, claims e.JwtClaims, err error)
}

func New(secret string) jwtPack {
	return jwtData{
		Secret: secret,
	}
}

func (d jwtData) CreateJWT(playerId e.PlayerIdentity, expirationTime time.Time) (string, error) {
	var (
		err      error
		tokenStr string
	)
	claims := &e.JwtClaims{
		PlayerID: playerId.PlayerID,
		Username: playerId.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        random.RandString(20),
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenStr, err = token.SignedString([]byte(d.Secret))
	if err != nil {
		return tokenStr, errors.Wrap(err, "error while creating jwt")
	}

	return tokenStr, err
}

func (d jwtData) ParseJWT(tokenStr string) (token *jwt.Token, claims e.JwtClaims, err error) {
	token, err = jwt.ParseWithClaims(tokenStr, &claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(d.Secret), nil
	})
	if err != nil {
		return token, claims, errors.Wrap(err, "error while parsing jwt token")
	}

	return token, claims, err
}
