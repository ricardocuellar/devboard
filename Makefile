# Variables
APP_NAME=devboard
MAIN_PATH=./cmd/api
MIGRATE_PATH=./migrations
DB_URL=postgresql://postgres:password@localhost:5432/devboard?sslmode=disable
JWT=dev-secret-key-change-in-production-min-32-chars$$
IDENTITY_PATH=./services/identity
IDENTITY_MAIN_PATH=${IDENTITY_PATH}/cmd/api
IDENTITY_MIGRATE_PATH=$(IDENTITY_PATH)/migrations
IDENTITY_DB_URL=postgresql://identity:identity@localhost:5433/identity?sslmode=disable
WORKSPACE_PATH=./services/workspace
WORKSPACE_MAIN_PATH=${WORKSPACE_PATH}/cmd/api
WORKSPACE_MIGRATE_PATH=${WORKSPACE_PATH}/migrations
WORKSPACE_DB_URL=postgresql://workspace:workspace@localhost:5434/workspace?sslmode=disable

.PHONY: run build test lint migrate-up migrate-down generate tidy help docs run-identity identity-generate identity-migrate-create identity-migrate-up identity-migrate-down run-workspace workspace-generate

## run: correr la aplicación
run:
	GO_ENV=development DATABASE_URL=$(DB_URL) JWT_SECRET=$(JWT) REDIS_ADDR=localhost:6379 go run $(MAIN_PATH)/main.go

## build: compilar el binario
build:
	go build -o bin/$(APP_NAME) $(MAIN_PATH)

## test: correr todos los tests con race detector
test:
	go test -race -cover ./...

## lint: analizar el código con golangci-lint
lint:
	golangci-lint run ./...

## migrate-create: crea una nueva migración (OJO para usarlo su uso es: make migrate-create name=create_users)
migrate-create:
	migrate create -ext sql -dir $(MIGRATE_PATH) -seq $(name)

## migrate-up: aplicar todas las migraciones pendientes
migrate-up:
	migrate -path $(MIGRATE_PATH) -database "$(DB_URL)" up

## migrate-down: revertir la última migración
migrate-down:
	migrate -path $(MIGRATE_PATH) -database "$(DB_URL)" down 1

## generate: correr go generate en todo el proyecto
generate:
	sqlc generate
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

## docs: generar documentación OpenAPI con swaggo
docs: 
	swag init -g cmd/api/main.go -o docs

## run-identity: Inicia servidor de identity service
run-identity: 
	IDENTITY_HTTP_ADDR=:8081 \
	IDENTITY_DATABASE_URL=$(IDENTITY_DB_URL) \
	JWT_SECRET=$(JWT) \
	REDIS_ADDR=localhost:6379 \
	go run $(IDENTITY_MAIN_PATH)

## identity-generate: Genera código Go con sqlc
identity-generate:
	cd $(IDENTITY_PATH) && sqlc generate

## identity-migrate-create: crea una nueva migración de identity (OJO para usarlo su uso es: make migrate-create name=create_users)
identity-migrate-create:
	migrate create -ext sql -dir $(IDENTITY_MIGRATE_PATH) -seq $(name)

## identity-migrate-up: aplicar todas las migraciones pendientes para identity
identity-migrate-up:
	migrate -path $(IDENTITY_MIGRATE_PATH) -database "$(IDENTITY_DB_URL)" up

## identity-migrate-down: revertir la última migración de identity
identity-migrate-down:
	migrate -path $(IDENTITY_MIGRATE_PATH) -database "$(IDENTITY_DB_URL)" down 1


## run-workspace: Inicia servidor de workspace service
run-workspace: 
	WORKSPACE_HTTP_ADDR=:8082 \
	WORKSPACE_DATABASE_URL=$(WORKSPACE_DB_URL) \
	JWT_SECRET=$(JWT) \
	go run $(WORKSPACE_MAIN_PATH)

## workspace-generate: Genera código Go con sqlc
workspace-generate:
	cd $(WORKSPACE_PATH) && sqlc generate


## workspace-migrate-create: crea una nueva migración de workspace (OJO para usarlo su uso es: make migrate-create name=create_users)
workspace-migrate-create:
	migrate create -ext sql -dir $(WORKSPACE_MIGRATE_PATH) -seq $(name)

## workspace-migrate-up: aplicar todas las migraciones pendientes para workspace
workspace-migrate-up:
	migrate -path $(WORKSPACE_MIGRATE_PATH) -database "$(WORKSPACE_DB_URL)" up

## workspace-migrate-down: revertir la última migración de workspace
workspace-migrate-down:
	migrate -path $(WORKSPACE_MIGRATE_PATH) -database "$(WORKSPACE_DB_URL)" down 1
