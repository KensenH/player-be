package config

import (
	"os"

	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Database Database `yaml:"database"`
}

type Database struct {
	Postgres  string `yaml:"postgres"`
	Redis     Redis  `yaml:"redis"`
	JWTSecret string `yaml:"jwt_secret"`
}

type Redis struct {
	Address  string `yaml:"address"`
	Password string `yaml:"password"`
}

func New() Config {
	var cfg Config

	cfgBytes, err := os.ReadFile("./config.yaml")
	if err != nil {
		log.Fatalf("[Config] error reading config file: %v", err)
	}

	err = yaml.Unmarshal(cfgBytes, &cfg)
	if err != nil {
		log.Fatalf("[Config] error parsing config: %v", err)
	}

	return cfg
}
