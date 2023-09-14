FROM golang:1.19-alpine

WORKDIR /usr/src/app

COPY ./ ./
RUN go mod download

RUN go build -o ./build/user-segments ./cmd/user-segments/main.go

CMD ["./build/user-segments"]