package config

import (
	"time"

	"github.com/SleepNFire/mediakeys/ad-serving/pkg"
	"github.com/kelseyhightower/envconfig"
)

type Redis struct {
	Host       string
	Port       string
	Expiration time.Duration
}

type Grpc struct {
	Addr     string
	Port     string
	CertPath string
}

type Config struct {
	Level      string // not used
	Redis      Redis
	Impression Grpc
}

func Default() Config {
	return Config{
		Level: "info",
		Redis: Redis{
			Host:       "redis",
			Port:       "6379",
			Expiration: 600 * time.Second,
		},
		Impression: Grpc{
			Addr:     "impression-tracking",
			Port:     "8503",
			CertPath: "/app/certificat",
		},
	}
}

func Init() (*Config, error) {
	config := Default()

	err := envconfig.Process("adserving", &config)
	if err != nil {
		return nil, pkg.ErrEnvironnementVariable
	}

	return &config, nil
}
