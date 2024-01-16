protoc_gen:
	protoc --go_out=grpcgen --go-grpc_out=grpcgen grpcgen/print.proto

build_serving: 
	go build -o ./ad-serving/bin/serving ./ad-serving/cmd/adserving/ad-serving.go

serving_app_test:
	docker compose up -d
	FUNCTIONNEL_TEST=1 ADSERVING_IMPRESSION_CERTPATH=../../../certificat/ ADSERVING_REDIS_EXPIRATION=1s go test -cover -v ./ad-serving/internal/app/

serving_functionnel_test: 
	docker compose up -d
	FUNCTIONNEL_TEST=1 ADSERVING_IMPRESSION_CERTPATH=../../../certificat/ ADSERVING_REDIS_EXPIRATION=1s go test -cover -v ./ad-serving/...

build_impression: 
	go build -o ./impression-tracking/bin/impression ./impression-tracking/cmd/impression/impression.go

impression_app_test:
	docker compose up -d
	FUNCTIONNEL_TEST=1 IMPRESSION_GRPC_CERTPATH=../../../certificat/ IMPRESSION_REDIS_EXPIRATION=1s go test -cover -v ./impression-tracking/internal/app/

impression_functionnel_test: 
	docker compose up -d
	FUNCTIONNEL_TEST=1 IMPRESSION_GRPC_CERTPATH=../../../certificat/ IMPRESSION_REDIS_EXPIRATION=2s go test -cover -v ./impression-tracking/...


start: build_serving build_impression
	docker compose down
	docker compose up -d