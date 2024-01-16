package impression

import (
	"context"
	"fmt"
	"net"
	"testing"

	"github.com/SleepNFire/mediakeys/ad-serving/config"
	print_grpc "github.com/SleepNFire/mediakeys/grpcgen/print.grpc"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
)

type impTest struct {
	print_grpc.UnimplementedImpressionServer
}

func (test *impTest) GetNumber(ctx context.Context, AdId *print_grpc.AdvertID) (*print_grpc.AdvertPrint, error) {
	if AdId.Id == "error" {
		return nil, fmt.Errorf("some_err")
	}
	return &print_grpc.AdvertPrint{Print: 10}, nil
}
func (test *impTest) Inc(ctx context.Context, AdId *print_grpc.AdvertID) (*print_grpc.Message, error) {
	if AdId.Id == "error" {
		return nil, fmt.Errorf("some_err")
	}
	return &print_grpc.Message{Message: "some_message"}, nil
}

func NewTestImpressionGrpc(global config.Config) (*ImpressionGrpc, error) {
	conn, err := grpc.Dial(global.Impression.Addr, grpc.WithInsecure())
	if err != nil {
		log.Error().Err(err).Msg("error during connection on impression grpc's server")
		return nil, err
	}

	client := print_grpc.NewImpressionClient(conn)

	return &ImpressionGrpc{Client: client}, nil
}

func TestImpressionGrpc_GetNumber(t *testing.T) {

	lis, err := net.Listen("tcp", ":50051")
	assert.NoError(t, err)
	defer lis.Close()

	srv := grpc.NewServer()
	imp := &impTest{}
	print_grpc.RegisterImpressionServer(srv, imp)

	go func() {
		if err := srv.Serve(lis); err != nil {
			assert.FailNow(t, "Failed to serve gRPC: "+err.Error())
		}
	}()
	defer srv.Stop()

	testClient, err := NewTestImpressionGrpc(config.Config{
		Impression: config.Grpc{
			Addr: lis.Addr().String(),
		},
	})
	assert.NoError(t, err)

	response, err := testClient.GetNumber("some_id")
	assert.NoError(t, err)
	assert.Equal(t, "10", response)

	response, err = testClient.GetNumber("error")
	assert.ErrorContains(t, err, "unknown")
	assert.Equal(t, "", response)

}

func TestImpressionGrpc_Inc(t *testing.T) {

	lis, err := net.Listen("tcp", ":50051")
	assert.NoError(t, err)
	defer lis.Close()

	srv := grpc.NewServer()
	imp := &impTest{}
	print_grpc.RegisterImpressionServer(srv, imp)

	go func() {
		if err := srv.Serve(lis); err != nil {
			assert.FailNow(t, "Failed to serve gRPC: "+err.Error())
		}
	}()
	defer srv.Stop()

	testClient, err := NewTestImpressionGrpc(config.Config{
		Impression: config.Grpc{
			Addr: lis.Addr().String(),
		},
	})
	assert.NoError(t, err)

	response, err := testClient.Inc("some_id")
	assert.NoError(t, err)
	assert.Equal(t, "some_message", response)

	response, err = testClient.Inc("error")
	assert.ErrorContains(t, err, "impossible to print an ad")
	assert.Equal(t, "", response)

}
