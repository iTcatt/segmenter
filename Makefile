.PHONY: all build clean test

all: build

build: 
	go build -o segmenter cmd/segmenter/main.go

docker-build:
	docker compose build

docker-run:
	docker compose up

swag:
	swag init -g ./cmd/segmenter/main.go

clean:
	rm segmenter

test:
	go test -v ./... -coverprofile=coverage.out
	go tool cover -html=coverage.out -o coverage.html
	rm coverage.out