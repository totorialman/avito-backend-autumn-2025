FROM golang:1.24-alpine

RUN apk add --no-cache git bash ca-certificates

RUN go install github.com/pressly/goose/v3/cmd/goose@latest

WORKDIR /migrations
COPY ./migrations /migrations

ENTRYPOINT ["goose"]