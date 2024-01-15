package config

import (
	"time"

	"github.com/SleepNFire/mediakeys/impression-tracking/pkg"
	"github.com/kelseyhightower/envconfig"
)

type Redis struct {
	Host       string
	Port       string
	Expiration time.Duration
}

type Config struct {
	Level string // not used
	Redis Redis
}

func Default() Config {
	return Config{
		Level: "info",
		Redis: Redis{
			Host:       "redis",
			Port:       "6379",
			Expiration: 60 * time.Second,
		},
	}
}

func Init() (*Config, error) {
	config := Default()

	err := envconfig.Process("impression", &config)
	if err != nil {
		return nil, pkg.ErrEnvironnementVariable
	}

	return &config, nil
}
