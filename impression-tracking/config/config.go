package config

import (
	"github.com/SleepNFire/mediakeys/impression-tracking/pkg"
	"github.com/kelseyhightower/envconfig"
)

type Redis struct {
	Host string
	Port string
}

type Config struct {
	Level string // not used
	Redis Redis
	Grpc  Grpc
}

type Grpc struct {
	CertPath string
	Port     string
}

func Default() Config {
	return Config{
		Level: "info",
		Redis: Redis{
			Host: "redis",
			Port: "6379",
		},
		Grpc: Grpc{
			CertPath: "/app/certificat",
			Port:     "8503",
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
