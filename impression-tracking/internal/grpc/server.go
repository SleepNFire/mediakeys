package grpc

import (
	"context"
	"net"

	"github.com/rs/zerolog/log"

	print_grpc "github.com/SleepNFire/mediakeys/grpcgen/print.grpc"
	"github.com/SleepNFire/mediakeys/impression-tracking/internal/printing"
	"go.uber.org/fx"
	rpc "google.golang.org/grpc"
)

func Init(lc fx.Lifecycle, impression *printing.PrintingGrpc) (*rpc.Server, error) {

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Error().Err(err).Msg("failed to listen")
		return nil, err
	}
	srv := rpc.NewServer()

	print_grpc.RegisterImpressionServer(srv, impression)

	lc.Append(
		fx.Hook{
			OnStart: func(ctx context.Context) error {
				go startServer(srv, lis)
				return nil
			},
			OnStop: func(ctx context.Context) error {
				srv.Stop()
				return nil
			},
		},
	)

	return srv, nil
}

func startServer(srv *rpc.Server, lis net.Listener) {
	log.Info().Interface("listener: ", lis.Addr()).Msg("RPC's server is starting")
	if err := srv.Serve(lis); err != nil {
		log.Error().Err(err).Msg("Server fail to start")
	}
}
