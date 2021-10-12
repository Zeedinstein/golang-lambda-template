.PHONY: build run deploy tidy

build: tidy
	sam build

run: 
	sam local start-api --docker-network channel-api_dynamodb-net

deploy: build
	sls deploy --verbose

tidy:
	go mod tidy
