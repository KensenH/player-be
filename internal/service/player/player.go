package player

type Option func(*PlayerService)

type PlayerService struct {
	Data PlayerData
}

type PlayerData interface {
}

func New(playerData PlayerData, opts ...Option) PlayerService {
	playerService := PlayerService{
		Data: playerData,
	}

	for _, opt := range opts {
		opt(&playerService)
	}

	return playerService
}
