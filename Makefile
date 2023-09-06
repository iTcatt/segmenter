.PHONY: all build clean

all: build

build: 
	go build -o ./build/user-sergments ./cmd/user-segments/main.go

clean:
	rm user-segments