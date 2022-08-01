# syntax=docker/dockerfile:1
FROM golang:1.18.5-alpine

WORKDIR /app
COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY * ./
RUN mkdir -p ./data/orbitdb
RUN go build ./cmd/astroid-api/main.go -o asteroid-api

EXPOSE 3000
CMD ["asteroid-api"]