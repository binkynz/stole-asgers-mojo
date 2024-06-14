package config

import (
	"github.com/caarlos0/env/v11"
)

type Config struct {
	apiConfig
	dbConfig
}

type apiConfig struct {
	API_KEY string `env:"API_KEY,required"`
}

type dbConfig struct {
	DATABASE_URL string `env:"DATABASE_URL,required"`
}

func NewConfig() (Config, error) {
	cfg := Config{}
	if err := env.Parse(&cfg.apiConfig); err != nil {
		return cfg, err
	}

	if err := env.Parse(&cfg.dbConfig); err != nil {
		return cfg, err
	}

	return cfg, nil
}
