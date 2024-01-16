package grpc

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"net"
	"os"
	"path/filepath"

	"github.com/rs/zerolog/log"

	print_grpc "github.com/SleepNFire/mediakeys/grpcgen/print.grpc"
	"github.com/SleepNFire/mediakeys/impression-tracking/config"
	"github.com/SleepNFire/mediakeys/impression-tracking/internal/printing"
	"go.uber.org/fx"
	"google.golang.org/grpc"
	rpc "google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func Init(lc fx.Lifecycle, conf *config.Config, impression *printing.PrintingGrpc) (*rpc.Server, error) {
	absPath, err := filepath.Abs(conf.Grpc.CertPath)
	if err != nil {
		log.Error().Err(err).Msg("failed to find certificat folder")
		return nil, err
	}

	caPem, err := os.ReadFile(absPath + "/ca-cert.pem")
	if err != nil {
		log.Error().Err(err).Str("path", absPath).Msg("failed to read ca-cert")
		return nil, err
	}

	certPool := x509.NewCertPool()
	if !certPool.AppendCertsFromPEM(caPem) {
		log.Error().Err(err).Msg("failed to read ca-cert")
		return nil, err
	}

	serverCert, err := tls.LoadX509KeyPair(absPath+"/server-cert.pem", absPath+"/server-key.pem")
	if err != nil {
		log.Error().Err(err).Str("path", absPath).Msg("failed to read ca-cert")
		return nil, err
	}

	tlsCert := &tls.Config{
		Certificates: []tls.Certificate{serverCert},
		ClientAuth:   tls.RequireAndVerifyClientCert,
		ClientCAs:    certPool,
	}

	tlsCredential := credentials.NewTLS(tlsCert)

	srv := rpc.NewServer(grpc.Creds(tlsCredential))

	print_grpc.RegisterImpressionServer(srv, impression)

	lis, err := net.Listen("tcp", ":"+conf.Grpc.Port)
	if err != nil {
		log.Error().Err(err).Msg("failed to listen")
		return nil, err
	}

	lc.Append(
		fx.Hook{
			OnStart: func(ctx context.Context) error {
				go startServer(srv, lis)
				return nil
			},
			OnStop: func(ctx context.Context) error {
				srv.Stop()
				lis.Close()
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
