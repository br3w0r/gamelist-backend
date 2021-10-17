.PHONY: build run test migrate migrate-undo fmt .build-lint

dbopt = "user=postgres password=pgpass sslmode=disable dbname=gamelist"

build:
	go build server.go

run:
	go run server.go

test:
	go test ./...

migrate:
	goose -dir ./migrations postgres $(dbopt) up

migrate-undo:
	goose -dir ./migrations postgres $(dbopt) down

fmt:
	go fmt ./...

build-lint:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

lint:
	golangci-lint run
