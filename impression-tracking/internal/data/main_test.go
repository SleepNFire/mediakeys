package data_test

import (
	"context"
	"os"
	"testing"

	"github.com/SleepNFire/mediakeys/impression-tracking/config"
	"github.com/SleepNFire/mediakeys/impression-tracking/internal/data"
	"go.uber.org/fx"
)

var Redis data.RedisAccessor

func TestMain(m *testing.M) {
	if os.Getenv("FUNCTIONNEL_TEST") != "1" {
		return
	}

	ctx := context.Background()

	testApp := fx.New(
		fx.Options(
			fx.Provide(config.Init),
			fx.Provide(data.NewRedisAccessor),
			fx.Invoke(func(redis *data.RedisAccessor) {
				Redis = *redis
			}),
		),
	)

	if err := testApp.Start(ctx); err != nil {
		os.Exit(1)
	}
	defer testApp.Stop(ctx)

	code := m.Run()

	os.Exit(code)
}
