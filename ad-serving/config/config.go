package config

import (
	"errors"
	"time"

	"github.com/kelseyhightower/envconfig"
)

type Redis struct {
	Host       string
	Port       string
	Expiration time.Duration
}

type Config struct {
	Debug string
	Redis Redis
}

func Default() Config {
	return Config{
		Debug: "info",
		Redis: Redis{
			Host:       "127.0.0.1",
			Port:       "6379",
			Expiration: 60 * time.Second,
		},
	}
}

func Init() (*Config, error) {
	config := Default()

	err := envconfig.Process("adserving", &config)
	if err != nil {
		return nil, errors.New("failed to read environnement variable")
	}

	return &config, nil
}
