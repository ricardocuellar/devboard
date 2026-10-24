FROM golang:1.26.4-alpine AS builder

WORKDIR /src

COPY contracts/go.mod contracts/go.sum ./contracts/
COPY services/notification/go.mod services/notification/go.sum ./services/notification/

WORKDIR /src/services/notification

RUN go mod download && go mod verify

WORKDIR /src

COPY contracts ./contracts
COPY services/notification ./services/notification

WORKDIR /src/services/notification

RUN CGO_ENABLED=0 \
    GOOS=linux \
    go build \
    -trimpath \
    -ldflags="-s -w" \
    -o /out/consumer \
    ./cmd/consumer


FROM alpine:3.24

RUN apk add --no-cache ca-certificates \
    && addgroup -S app \
    && adduser -S app -G app

WORKDIR /app

COPY --from=builder /out/consumer /app/consumer

USER app

ENTRYPOINT ["/app/consumer"]