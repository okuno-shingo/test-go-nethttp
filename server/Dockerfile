FROM golang:1.22-alpine

WORKDIR /app
COPY server.go .

RUN apk update && apk add git

RUN go install github.com/cosmtrek/air@v1.45.0 && \
  chmod +x ${GOPATH}/bin/air

COPY go.mod go.sum .air.toml ./

ENV TZ="Asia/Tokyo"

RUN go mod download

RUN go build -o server server.go
CMD air -c .air.toml
