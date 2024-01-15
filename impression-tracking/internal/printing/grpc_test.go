package printing

import (
	"context"
	"net"
	"testing"

	"github.com/SleepNFire/mediakeys/ad-serving/pkg"
	print_grpc "github.com/SleepNFire/mediakeys/grpcgen/print.grpc"
	printing_mock "github.com/SleepNFire/mediakeys/impression-tracking/internal/mocks/printing"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestPrintingGrpc_GetNumber(t *testing.T) {

	ctrl := gomock.NewController(t)
	mockRedis := printing_mock.NewMockCacheAccessor(ctrl)

	imp, err := NewPrintingGrpc(mockRedis)
	assert.NoError(t, err)

	lis, err := net.Listen("tcp", ":50051")
	assert.NoError(t, err)
	defer lis.Close()

	srv := grpc.NewServer()

	print_grpc.RegisterImpressionServer(srv, imp)

	go func() {
		if err := srv.Serve(lis); err != nil {
			assert.FailNow(t, "Failed to serve gRPC: "+err.Error())
		}
	}()
	defer srv.Stop()

	tests := []struct {
		name             string
		prepareMock      func(mr *printing_mock.MockCacheAccessor)
		body             *print_grpc.AdvertID
		expectedResponse *print_grpc.AdvertPrint
		expectedErr      error
	}{
		{
			name: "success",
			prepareMock: func(mr *printing_mock.MockCacheAccessor) {
				mr.EXPECT().Find("some_id").Return(uint64(10), nil)
			},
			body:             &print_grpc.AdvertID{Id: "some_id"},
			expectedResponse: &print_grpc.AdvertPrint{Print: uint64(10)},
			expectedErr:      nil,
		},
		{
			name: "not found",
			prepareMock: func(mr *printing_mock.MockCacheAccessor) {
				mr.EXPECT().Find("some_id").Return(uint64(0), pkg.ErrNotFound)
			},
			body:             &print_grpc.AdvertID{Id: "some_id"},
			expectedResponse: nil,
			expectedErr:      status.Error(codes.NotFound, pkg.ErrNotFound.Error()),
		},
		{
			name: "internal error",
			prepareMock: func(mr *printing_mock.MockCacheAccessor) {
				mr.EXPECT().Find("some_id").Return(uint64(0), pkg.ErrInternalError)
			},
			body:             &print_grpc.AdvertID{Id: "some_id"},
			expectedResponse: nil,
			expectedErr:      status.Error(codes.Internal, pkg.ErrInternalError.Error()),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.prepareMock(mockRedis)
			conn, err := grpc.Dial(lis.Addr().String(), grpc.WithInsecure())
			assert.NoError(t, err)
			defer conn.Close()

			client := print_grpc.NewImpressionClient(conn)
			response, err := client.GetNumber(context.Background(), tt.body)
			assert.Equal(t, tt.expectedErr, err)
			if tt.expectedResponse != nil && response != nil {
				assert.Equal(t, tt.expectedResponse.Print, response.Print)
			} else if tt.expectedResponse == nil {
				assert.Nil(t, response)
			} else {
				assert.FailNow(t, "result not expected")
			}
		})
	}
}

func TestPrintingGrpc_Inc(t *testing.T) {

	ctrl := gomock.NewController(t)
	mockRedis := printing_mock.NewMockCacheAccessor(ctrl)

	imp, err := NewPrintingGrpc(mockRedis)
	assert.NoError(t, err)

	lis, err := net.Listen("tcp", ":50051")
	assert.NoError(t, err)
	defer lis.Close()

	srv := grpc.NewServer()

	print_grpc.RegisterImpressionServer(srv, imp)

	go func() {
		if err := srv.Serve(lis); err != nil {
			assert.FailNow(t, "Failed to serve gRPC: "+err.Error())
		}
	}()
	defer srv.Stop()

	tests := []struct {
		name             string
		prepareMock      func(mr *printing_mock.MockCacheAccessor)
		body             *print_grpc.AdvertID
		expectedResponse *print_grpc.Error
		expectedErr      error
	}{
		{
			name: "success",
			prepareMock: func(mr *printing_mock.MockCacheAccessor) {
				mr.EXPECT().Inc("some_id").Return(nil)
			},
			body:             &print_grpc.AdvertID{Id: "some_id"},
			expectedResponse: &print_grpc.Error{Error: "SUCCESS"},
			expectedErr:      nil,
		},
		{
			name: "not found",
			prepareMock: func(mr *printing_mock.MockCacheAccessor) {
				mr.EXPECT().Inc("some_id").Return(pkg.ErrNotFound)
			},
			body:             &print_grpc.AdvertID{Id: "some_id"},
			expectedResponse: nil,
			expectedErr:      status.Error(codes.NotFound, pkg.ErrNotFound.Error()),
		},
		{
			name: "internal error",
			prepareMock: func(mr *printing_mock.MockCacheAccessor) {
				mr.EXPECT().Inc("some_id").Return(pkg.ErrInternalError)
			},
			body:             &print_grpc.AdvertID{Id: "some_id"},
			expectedResponse: nil,
			expectedErr:      status.Error(codes.Internal, pkg.ErrInternalError.Error()),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.prepareMock(mockRedis)
			conn, err := grpc.Dial(lis.Addr().String(), grpc.WithInsecure())
			assert.NoError(t, err)
			defer conn.Close()

			client := print_grpc.NewImpressionClient(conn)
			response, err := client.Inc(context.Background(), tt.body)
			assert.Equal(t, tt.expectedErr, err)
			if tt.expectedResponse != nil && response != nil {
				assert.Equal(t, tt.expectedResponse.Error, response.Error)
			} else if tt.expectedResponse == nil {
				assert.Nil(t, response)
			} else {
				assert.FailNow(t, "result not expected")
			}
		})
	}
}
