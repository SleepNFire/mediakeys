package app

import (
	"github.com/SleepNFire/mediakeys/ad-serving/config"
	"github.com/SleepNFire/mediakeys/ad-serving/internal/advert"
	"github.com/SleepNFire/mediakeys/ad-serving/internal/data"
	"github.com/SleepNFire/mediakeys/ad-serving/internal/impression"
	"github.com/SleepNFire/mediakeys/ad-serving/internal/rest"
	"go.uber.org/fx"
)

func Init() *fx.App {
	app := fx.New(
		fx.Options(
			fx.Provide(config.Init),
			fx.Provide(
				fx.Annotate(
					data.NewRedisAccessor,
					fx.As(new(advert.CacheAccessor)),
					fx.As(new(rest.ApiTechnical)),
				),
			),
			fx.Provide(
				fx.Annotate(
					impression.NewImpressionGrpc,
					fx.As(new(advert.ImpressionAccessor)),
				),
			),
			fx.Provide(advert.NewAdvertEndpoint),
			fx.Invoke(rest.Init),
		),
	)

	return app
}
