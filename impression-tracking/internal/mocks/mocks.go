package mocks

//go:generate mockgen -destination=./printing/mock_cache_accessor.go -package=printing_mock -source=../printing/grpc.go CacheAccessor
