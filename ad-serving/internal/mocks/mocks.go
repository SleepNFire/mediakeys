package mocks

//go:generate mockgen -destination=./advert/mock_cache_accessor.go -package=advert_mock -source=../advert/rest.go CacheAccessor
