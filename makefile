serving_app_test:
	FUNCTIONNEL_TEST=1 ADSERVING_REDIS_EXPIRATION=1s go test -cover -v ./ad-serving/internal/app/

serving_functionnel_test: 
	FUNCTIONNEL_TEST=1 ADSERVING_REDIS_EXPIRATION=1s go test -cover -v ./ad-serving/...