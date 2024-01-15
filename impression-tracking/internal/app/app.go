package app

import (
	"github.com/SleepNFire/mediakeys/impression-tracking/config"
	"github.com/SleepNFire/mediakeys/impression-tracking/internal/data"
	"github.com/SleepNFire/mediakeys/impression-tracking/internal/grpc"
	"github.com/SleepNFire/mediakeys/impression-tracking/internal/printing"
	"github.com/SleepNFire/mediakeys/impression-tracking/internal/rest"
	"go.uber.org/fx"
)

func Init() *fx.App {

	app := fx.New(
		fx.Options(
			fx.Provide(config.Init),
			fx.Provide(printing.NewPrintingGrpc),
			fx.Provide(
				fx.Annotate(
					data.NewRedisAccessor,
					fx.As(new(printing.CacheAccessor)),
					fx.As(new(rest.ApiTechnical)),
				),
			),
			fx.Invoke(grpc.Init),
			fx.Invoke(rest.Init),
		),
	)

	return app
}
