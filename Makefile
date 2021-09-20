.PHONY: build clean deploy tidy

build: tidy
	export GO111MODULE=on
	env GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o bin/api api/main.go
	env GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o bin/cron cron/main.go

clean:
	rm -rf ./bin ./vendor go.sum

deploy: clean build
	sls deploy --verbose

tidy:
	go mod tidy
