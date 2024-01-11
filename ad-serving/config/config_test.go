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
				{"ADSERVING_DEBUG", "some_debug"},
				{"ADSERVING_REDIS_HOST", "custom_localhost"},
				{"ADSERVING_REDIS_PORT", "some_port"},
				{"ADSERVING_REDIS_EXPIRATION", "5s"},
			},
			want: Config{
				Debug: "some_debug",
				Redis: Redis{
					Host:       "custom_localhost",
					Port:       "some_port",
					Expiration: 5 * time.Second,
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
