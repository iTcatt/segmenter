.PHONY: all build clean

all: build

build: 
	go build -o segmenter cmd/segmenter/main.go

docker-build:
	docker compose build

docker-run:
	docker compose up

clean:
	rm segmenter