# MyGameList Backend

Backend written in Go.

## MyGameList family

- [Backend](https://github.com/br3w0r/gamelist-backend)
- [Frontend](https://github.com/br3w0r/gamelist-frontend)
- [Site Scraper](https://github.com/br3w0r/gamelist-scraper)
- [Proto file](https://github.com/br3w0r/gamelist-proto)

## Try it yourself

[mygamelist.tk](https://mygamelist.tk/) is available to try how MyGameList works. No legit email address is required, but all data is subject to change or complete removal.

## Technologies

- Gin
- GORM
- gRPC
- REST
- JWT

## Running in development mode

The easiest way to run the backend is to use `go get` and `go run server.go` commands. It will serve static files and allow addition of new games via api, but if you want to use scraper, you should set `FORCE_SCRAPE` environment variable to `1`. For example,

```bash
FORCE_SCRAPE=1 go run server.go
```

This will connect to scraper gRPC server on localhost and fetch the games.

## Docker building and running

### Build

```bash
# docker build -t backend .
```

### Common run command

```bash
# docker run -p 8080:8080 \
    --network gamelist \
    -v <path_to_static>/ \
    -v gamelist-data:/data \
    -e SERVE_STATIC=0 \
    -e PRODUCTION_MODE=1 \
    -e STATIC_FOLDER=/static \
    -e DATABASE_DIST=/data/gamelist.db \
    -e FORCE_SCRAPE=0 \
    -e SCRAPER_GRPC_ADDRESS=scraper \
    gamelist-backend
```

Replace `<path_to_static>` with your path

If you want to add games by yourself, set `PRODUCTION_MODE=0` and add them with api requests (for example, with Postman).

If you want to use scraper and used its tutorial to build and run it, just set `FORCE_SCRAPE=1`

## API and data structure

Check [api_desctiption.md](/api_desctiption.md) and [data_structure.md](/data_structure.md) if you want to make api requests to the backend.
