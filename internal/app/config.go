package app

import "github.com/ilyakaznacheev/cleanenv"

type Config struct {
	Token     string `env:"TOKEN" env-required"`
	Query     string `env:"QUERY" env-required`
	FileUsers string `env:"FILE_USERS" env-default:"users.json"  env-required`
	FileAds   string `env:"FILE_ADS" env-default:"ads.json"  env-required`
}

// NewConfig returns app configuration
func NewConfig() (*Config, error) {
	config := Config{}
	if err := cleanenv.ReadEnv(&config); err != nil {
		return nil, err
	}
	return &config, nil
}
