FROM golang:1.22-alpine

WORKDIR /app
COPY client.go .

RUN go install github.com/cosmtrek/air@v1.45.0 && \
  chmod +x ${GOPATH}/bin/air

COPY go.mod go.sum .air.toml ./

ENV TZ="Asia/Tokyo"

RUN go mod download

RUN go build -o client client.go
CMD air -c .air.toml