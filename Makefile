.PHONY: build run test

build:
	go build server.go

run:
	go run server.go

test:
	go test ./...
