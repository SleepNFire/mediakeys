package impression

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"github.com/SleepNFire/mediakeys/ad-serving/config"
	print_grpc "github.com/SleepNFire/mediakeys/grpcgen/print.grpc"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type ImpressionGrpc struct {
	Client print_grpc.ImpressionClient
}

func NewImpressionGrpc(global *config.Config) (*ImpressionGrpc, error) {
	absPath, err := filepath.Abs(global.Impression.CertPath)
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

	serverCert, err := tls.LoadX509KeyPair(absPath+"/client-cert.pem", absPath+"/client-key.pem")
	if err != nil {
		log.Error().Err(err).Str("path", absPath).Msg("failed to read client cert and key")
		return nil, err
	}

	tlsCert := &tls.Config{
		Certificates: []tls.Certificate{serverCert},
		RootCAs:      certPool,
	}

	tlsCredential := credentials.NewTLS(tlsCert)

	conn, err := grpc.Dial(fmt.Sprintf("%s:%s", global.Impression.Addr, global.Impression.Port), grpc.WithTransportCredentials(tlsCredential))
	if err != nil {
		log.Error().Err(err).Msg("error during connection on impression grpc's server")
		return nil, err
	}

	client := print_grpc.NewImpressionClient(conn)

	return &ImpressionGrpc{Client: client}, nil
}

func (imp *ImpressionGrpc) GetNumber(id string) (string, error) {
	response, err := imp.Client.GetNumber(context.Background(), &print_grpc.AdvertID{Id: id})
	if err != nil {
		return "", fmt.Errorf("unknown")
	}
	valueStr := strconv.FormatUint(response.Print, 10)
	return valueStr, nil
}

func (imp *ImpressionGrpc) Inc(id string) (string, error) {
	response, err := imp.Client.Inc(context.Background(), &print_grpc.AdvertID{Id: id})
	if err != nil {
		log.Error().Err(err).Interface("response", response).Msg("impossible to print an ad")
		return "", fmt.Errorf("impossible to print an ad")
	}
	return response.Message, nil
}
