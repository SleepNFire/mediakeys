package printing

import (
	"context"

	"github.com/SleepNFire/mediakeys/ad-serving/pkg"
	print_grpc "github.com/SleepNFire/mediakeys/grpcgen/print.grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type CacheAccessor interface {
	Inc(id string) error
	Find(id string) (uint64, error)
}

type PrintingGrpc struct {
	Redis CacheAccessor
	print_grpc.UnimplementedImpressionServer
}

func NewPrintingGrpc(redis CacheAccessor) (*PrintingGrpc, error) {
	return &PrintingGrpc{
		Redis: redis,
	}, nil
}

func (imp *PrintingGrpc) GetNumber(ctx context.Context, key *print_grpc.AdvertID) (*print_grpc.AdvertPrint, error) {
	value, err := imp.Redis.Find(key.Id)
	if err != nil {
		switch err {
		case pkg.ErrNotFound:
			return nil, status.Error(codes.NotFound, pkg.ErrNotFound.Error())
		default:
			return nil, status.Error(codes.Internal, pkg.ErrInternalError.Error())
		}
	}
	return &print_grpc.AdvertPrint{Print: value}, nil
}

func (imp *PrintingGrpc) Inc(ctx context.Context, key *print_grpc.AdvertID) (*print_grpc.Error, error) {
	err := imp.Redis.Inc(key.Id)
	if err != nil {
		switch err {
		case pkg.ErrNotFound:
			return nil, status.Error(codes.NotFound, pkg.ErrNotFound.Error())
		default:
			return nil, status.Error(codes.Internal, pkg.ErrInternalError.Error())
		}
	}
	return &print_grpc.Error{Error: "SUCCESS"}, nil
}
