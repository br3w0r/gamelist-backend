goose -dir ./migrations postgres "user=postgres password=pgpass sslmode=disable dbname=gamelist" create $1 sql
