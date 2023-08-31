all: build

build: 
	go build -v ./cmd/user-segments 

clean:
	rm user-segments