.PHONY: all build clean

all: build

build: 
	go build -o ./build/user-segments ./cmd/user-segments/main.go

docker-build:
	docker-compose build

docker-run:
	docker-compose up

clean:
	rm ./build/user-segments