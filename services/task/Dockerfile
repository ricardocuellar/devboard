FROM golang:1.26.4-alpine AS builder

WORKDIR /src

COPY contracts/go.mod contracts/go.sum ./contracts/
COPY services/task/go.mod services/task/go.sum ./services/task/

WORKDIR /src/services/task/

RUN go mod download \
    && go mod verify

WORKDIR /src

COPY contracts ./contracts
COPY services/task ./services/task

WORKDIR /src/services/task/


RUN CGO_ENABLED=0 \
    GOOS=linux \
    go build \
    -trimpath \
    -ldflags="-s -w" \
    -o /out/service \
    ./cmd/api


FROM alpine:3.24

RUN apk add --no-cache \
    ca-certificates \
    curl \
    && addgroup -S app \
    && adduser -S app -G app

WORKDIR /app

COPY --from=builder \
    /out/service \
    /app/service

USER app

ENTRYPOINT ["/app/service"]