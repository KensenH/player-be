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

	runEnv := os.Getenv("RUN_ENV")
	if runEnv != "" {
		cfg.Database.Postgres = "postgresql://postgres:mysecretpassword@postgres:5432/player"
		cfg.Database.Redis.Address = "redis:6379"
	}

	return cfg
}
