FROM golang:1.24-alpine AS builder

COPY go.mod go.sum /github.com/totorialman/avito-backend-autumn-2025/
WORKDIR /github.com/totorialman/avito-backend-autumn-2025

RUN go mod download
COPY . /github.com/totorialman/avito-backend-autumn-2025

RUN CGO_ENABLED=0 GOOS=linux go build -mod=readonly -o ./.bin ./cmd/main.go

FROM scratch AS runner

WORKDIR /build_v1/

COPY --from=builder /github.com/totorialman/avito-backend-autumn-2025/.bin .

COPY --from=builder /usr/local/go/lib/time/zoneinfo.zip /
ENV TZ="Europe/Moscow"
ENV ZONEINFO=/zoneinfo.zip

EXPOSE 5458

ENTRYPOINT ["./.bin"]