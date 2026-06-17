# Variables
APP_NAME=devboard
MAIN_PATH=./cmd/api
MIGRATE_PATH=./migrations
DB_URL=postgresql://postgres:password@localhost:5432/devboard?sslmode=disable

.PHONY: run build test lint migrate-up migrate-down generate tidy help

## run: correr la aplicación
run:
	go run $(MAIN_PATH)/main.go

## build: compilar el binario
build:
	go build -o bin/$(APP_NAME) $(MAIN_PATH)

## test: correr todos los tests con race detector
test:
	go test -race -cover ./...

## lint: analizar el código con golangci-lint
lint:
	golangci-lint run ./...

## migrate-up: aplicar todas las migraciones pendientes
migrate-up:
	migrate -path $(MIGRATE_PATH) -database "$(DB_URL)" up

## migrate-down: revertir la última migración
migrate-down:
	migrate -path $(MIGRATE_PATH) -database "$(DB_URL)" down 1

## generate: correr go generate en todo el proyecto
generate:
	go generate ./...

## tidy: limpiar y verificar dependencias
tidy:
	go mod tidy
	go mod verify

## help: mostrar este menú
help:
	@grep -E '^##' Makefile | sed 's/## //'

## docker-up: levantar los servicios de desarrollo
docker-up:
	docker compose up -d

## docker-down: detener los servicios de desarrollo
docker-down: 
	docker compose down

## docker-logs: ver logs de todos los servicios
docker-logs:
	docker compose logs -f