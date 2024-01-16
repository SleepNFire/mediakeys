package mocks

//go:generate mockgen -destination=./advert/mock_accessors.go -package=advert_mock -source=../advert/rest.go Accessors
