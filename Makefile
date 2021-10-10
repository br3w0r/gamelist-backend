.PHONY: build run test migrate

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
