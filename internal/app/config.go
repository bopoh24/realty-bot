package app

import (
	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Token     string `env:"TOKEN" env-required:"true"`
	Query     string `env:"QUERY" env-required:"true"`
	FileUsers string `env:"FILE_USERS" env-default:"users.json"  env-required:"true"`
	FileAds   string `env:"FILE_ADS" env-default:"ads.json"  env-required:"true"`
}

// NewConfig returns app configuration
func NewConfig() (*Config, error) {
	config := Config{}
	if err := cleanenv.ReadEnv(&config); err != nil {
		return nil, err
	}
	return &config, nil
}
