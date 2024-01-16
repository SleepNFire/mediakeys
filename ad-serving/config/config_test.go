package config

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestInit(t *testing.T) {
	tests := []struct {
		name          string
		environnement [][]string
		want          Config
	}{
		{
			name: "default",
			want: Default(),
		},
		{
			name: "custom_conf",
			environnement: [][]string{
				{"ADSERVING_LEVEL", "some_debug"},
				{"ADSERVING_REDIS_HOST", "custom_localhost"},
				{"ADSERVING_REDIS_PORT", "some_port"},
				{"ADSERVING_REDIS_EXPIRATION", "5s"},
				{"ADSERVING_IMPRESSION_ADDR", "some_addr"},
				{"ADSERVING_IMPRESSION_CERTPATH", "some_path"},
				{"ADSERVING_IMPRESSION_PORT", "some_port_2"},
			},
			want: Config{
				Level: "some_debug",
				Redis: Redis{
					Host:       "custom_localhost",
					Port:       "some_port",
					Expiration: 5 * time.Second,
				},
				Impression: Grpc{
					Addr:     "some_addr",
					Port:     "some_port_2",
					CertPath: "some_path",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Clearenv()
			if len(tt.environnement) > 0 {
				for _, env := range tt.environnement {
					os.Setenv(env[0], env[1])
				}
			}

			conf, err := Init()

			assert.NoError(t, err)
			assert.Equal(t, tt.want, *conf)

		})
	}
}
