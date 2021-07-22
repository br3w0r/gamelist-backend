FROM golang:1.16-alpine

RUN apk add g++
WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY server.go .
COPY service/ service/
COPY repository/ repository/
COPY proto/ proto/
COPY helpers/ helpers/
COPY entity/ entity/
COPY controller/ controller/
COPY server/ server/
COPY test/stress/ test/stress

RUN go build -o /server

EXPOSE 8080

CMD ["/server"]
