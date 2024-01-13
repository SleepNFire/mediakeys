build_serving: 
	go build -o ./ad-serving/bin/serving ./ad-serving/cmd/adserving/ad-serving.go

serving_app_test:
	docker compose up -d
	FUNCTIONNEL_TEST=1 ADSERVING_REDIS_EXPIRATION=1s go test -cover -v ./ad-serving/internal/app/

serving_functionnel_test: 
	docker compose up -d
	FUNCTIONNEL_TEST=1 ADSERVING_REDIS_EXPIRATION=1s go test -cover -v ./ad-serving/...