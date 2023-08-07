package jwt

import (
	"player-be/internal/pkg/random"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/pkg/errors"
)

type jwtData struct {
	Secret string
}

type jwtPack interface {
	CreateJWT(username string, expirationTime time.Time) (string, error)
	ParseJWT(tokenStr string) (token *jwt.Token, claimsID string, err error)
}

type JwtClaims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func New(secret string) jwtPack {
	return jwtData{
		Secret: secret,
	}
}

func (d jwtData) CreateJWT(username string, expirationTime time.Time) (string, error) {
	var (
		err      error
		tokenStr string
	)
	claims := &JwtClaims{
		Username: username,
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

func (d jwtData) ParseJWT(tokenStr string) (token *jwt.Token, claimsID string, err error) {
	var claims JwtClaims
	token, err = jwt.ParseWithClaims(tokenStr, &claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(d.Secret), nil
	})
	if err != nil {
		return token, claimsID, errors.Wrap(err, "error while parsing jwt token")
	}

	return token, claims.ID, err
}
