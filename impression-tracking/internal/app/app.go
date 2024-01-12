package app

import (
	"github.com/SleepNFire/mediakeys/impression-tracking/config"
	"github.com/SleepNFire/mediakeys/impression-tracking/internal/data"
	"github.com/SleepNFire/mediakeys/impression-tracking/internal/rest"
	"go.uber.org/fx"
)

func Init() *fx.App {

	app := fx.New(
		fx.Options(
			fx.Provide(config.Init),
			fx.Provide(data.NewRedisAccessor),
			fx.Invoke(rest.Init),
		),
	)

	return app
}
