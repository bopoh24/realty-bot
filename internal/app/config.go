package app

import (
	"github.com/ilyakaznacheev/cleanenv"
	"log"
)

type Config struct {
	Token     string `env:"TOKEN" env-required"`
	Query     string `env:"QUERY" env-required`
	FileUsers string `env:"FILE_USERS" env-default:"users.json"  env-required`
	FileAds   string `env:"FILE_ADS" env-default:"ads.json"  env-required`
}

// MustConfig returns app configuration
func MustConfig() *Config {
	config := Config{}
	if err := cleanenv.ReadEnv(&config); err != nil {
		log.Fatalf("Unable to load config: %s", err)
	}
	return &config
}
