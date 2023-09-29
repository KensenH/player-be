package player

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	e "player-be/internal/entity/player"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

type PlayerData struct {
	players map[int]e.PlayerDetail
}

type PlayerServiceStubs struct {
	data PlayerData
}

func (s PlayerServiceStubs) SignUp(ctx context.Context, playerForm e.PlayerSignUpForm) (e.PlayerIdentity, error) {
	return e.PlayerIdentity{}, nil
}
func (s PlayerServiceStubs) SignIn(ctx context.Context, expirationTime time.Time, playerForm e.PlayerUserPass) (tokenStr string, err error) {
	return "", nil
}
func (s PlayerServiceStubs) SignOut(ctx context.Context, tokenStr string) error {
	return nil
}
func (s PlayerServiceStubs) JWTTokenValid(ctx context.Context, tokenStr string) (bool, e.PlayerIdentity, error) {
	return false, e.PlayerIdentity{}, nil
}
func (s PlayerServiceStubs) GetPlayerDetail(ctx context.Context, playerId uint) (e.PlayerDetail, error) {
	var err error
	return s.data.players[int(playerId)], err
}
func (s PlayerServiceStubs) SearchPlayer(ctx context.Context, filter e.PlayerFilter) ([]e.PlayerDetail, error) {
	return []e.PlayerDetail{}, nil
}
func (s PlayerServiceStubs) AddBankAccount(ctx context.Context, bankAcc e.BankAccount) error {
	return nil
}
func (s PlayerServiceStubs) GetTopUpHistory(ctx context.Context, playerId uint) ([]e.TopUpHistory, error) {
	return []e.TopUpHistory{}, nil
}
func (s PlayerServiceStubs) TopUp(ctx context.Context, playerId uint, sum int64) (e.TopUpHistory, error) {
	return e.TopUpHistory{}, nil
}

func TestGetPlayerDetail(t *testing.T) {
	data := PlayerData{
		players: map[int]e.PlayerDetail{
			1: {
				Username:       "kensen",
				FirstName:      "Kensen",
				LastName:       "Huang",
				PhoneNumber:    "+6281348595422",
				Email:          "kensen.huang@gmail.com",
				InGameCurrency: 1000,
			},
		},
	}

	service := PlayerServiceStubs{
		data: data,
	}

	handler := New(service)

	ec := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := ec.NewContext(req, rec)
	c.SetPath("/detail/:id")
	c.SetParamNames("id")
	c.SetParamValues("1")

	if assert.NoError(t, handler.GetPlayerDetail(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)

		var got e.PlayerDetail
		err := json.Unmarshal(rec.Body.Bytes(), &got)
		if err != nil {
			t.Fail()
		}

		assert.Equal(t, data.players[1], got)
	}

}
