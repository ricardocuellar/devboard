FROM golang:1.26.4-alpine AS builder

WORKDIR /src

COPY go.mod ./
COPY cmd ./cmd
COPY internal ./internal


RUN CGO_ENABLED=0 \
    GOOS=linux \
    go build \
    -trimpath \
    -ldflags="-s -w" \
    -o /out/gateway \
    ./cmd/api


FROM alpine:3.24

RUN apk add --no-cache \
    ca-certificates \
    curl \
    && addgroup -S app \
    && adduser -S app -G app

WORKDIR /app

COPY --from=builder \
    /out/gateway \
    /app/gateway

USER app

ENTRYPOINT ["/app/gateway"]