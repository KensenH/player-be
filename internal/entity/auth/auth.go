package auth

type Auth struct {
	Authorization string `query:"token" json:"Authorization"`
}
