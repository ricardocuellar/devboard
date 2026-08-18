# Devboard

API REST para gestionar equipos de trabajo, proyectos, tareas y comentarios. Devboard está construida en Go y utiliza PostgreSQL para la persistencia, Redis para el manejo de tokens de actualización y JWT para autenticación.

## Características

- Registro, inicio y cierre de sesión con access tokens y refresh tokens.
- Gestión de usuarios y miembros de un workspace.
- Creación y consulta de workspaces y proyectos.
- Creación, asignación y actualización de estado de tareas.
- Comentarios asociados a tareas.
- Persistencia en PostgreSQL mediante `pgx` y consultas generadas con `sqlc`.
- Almacenamiento de tokens en Redis.
- Validación de solicitudes y respuestas de error consistentes.
- Middlewares de autenticación, logging, recuperación ante panics y seguridad HTTP.
- Documentación OpenAPI/Swagger disponible desde la aplicación.
- Apagado ordenado del servidor.

## Tecnologías

- [Go](https://go.dev/) 1.26.4
- [PostgreSQL](https://www.postgresql.org/) 16
- [Redis](https://redis.io/) 7
- [Docker Compose](https://docs.docker.com/compose/)
- [pgx](https://github.com/jackc/pgx)
- [sqlc](https://sqlc.dev/)
- [golang-migrate](https://github.com/golang-migrate/migrate)
- [Swaggo](https://github.com/swaggo/swag)

## Arquitectura

El proyecto separa las responsabilidades por capas:

```text
cmd/api/                 Punto de entrada y composición de dependencias
internal/domain/         Entidades y reglas del dominio
internal/handler/        Adaptadores HTTP
internal/usecase/        Casos de uso y lógica de aplicación
internal/repository/     Persistencia en PostgreSQL y memoria
internal/auth/           JWT y almacenamiento de refresh tokens
internal/cache/          Cliente de Redis
internal/middlewares/    Autenticación, logging, recovery y seguridad
internal/server/         Configuración del servidor HTTP
migrations/              Migraciones de PostgreSQL
docs/                    Especificación OpenAPI generada
configs/                 Configuración de referencia
```

## Requisitos

Para ejecutar el proyecto necesitas:

- Go 1.26.4 o una versión compatible con la indicada en `go.mod`.
- Docker con Docker Compose.
- [`golang-migrate`](https://github.com/golang-migrate/migrate/tree/master/cmd/migrate) para aplicar las migraciones.
- Opcionalmente, `sqlc`, `swag` y `golangci-lint` para generación de código, documentación y análisis estático.

## Puesta en marcha

1. Clona el repositorio:

   ```bash
   git clone https://github.com/ricardocuellar/devboard.git
   cd devboard
   ```

2. Descarga y verifica las dependencias:

   ```bash
   go mod download
   go mod verify
   ```

3. Inicia PostgreSQL y Redis:

   ```bash
   make docker-up
   ```

4. Aplica las migraciones:

   ```bash
   make migrate-up
   ```

5. Inicia la API:

   ```bash
   make run
   ```

La API quedará disponible en `http://localhost:8080`.

## Variables de entorno

El comando `make run` incluye valores adecuados para desarrollo local. Si ejecutas el binario directamente, configura estas variables:

| Variable | Descripción | Ejemplo local |
| --- | --- | --- |
| `GO_ENV` | Entorno de ejecución | `development` |
| `DATABASE_URL` | Cadena de conexión a PostgreSQL | `postgresql://postgres:password@localhost:5432/devboard?sslmode=disable` |
| `REDIS_ADDR` | Dirección de Redis | `localhost:6379` |
| `JWT_SECRET` | Secreto usado para firmar los JWT | Una cadena segura de al menos 32 caracteres |

> No reutilices las credenciales de desarrollo ni guardes secretos reales en el repositorio.

## Documentación y endpoints

Con el servidor en ejecución puedes consultar:

- Estado del servicio: `GET http://localhost:8080/health`
- Swagger UI: `http://localhost:8080/docs/`
- Base de la API: `http://localhost:8080/api/v1`

Los recursos principales son:

| Recurso | Operaciones disponibles |
| --- | --- |
| Autenticación | Registro, login, refresh y logout |
| Usuarios | Creación y consulta por ID |
| Workspaces | Creación, listado, consulta y administración de miembros |
| Proyectos | Creación, consulta y listado por workspace |
| Tareas | Creación, consulta, asignación, cambio de estado y listado por proyecto |
| Comentarios | Creación, listado por tarea y eliminación |

Las rutas de workspaces, proyectos, tareas y comentarios requieren un access token:

```http
Authorization: Bearer <access_token>
```

## Comandos útiles

```bash
make help          # Muestra los comandos disponibles
make run           # Ejecuta la API en modo desarrollo
make build         # Compila el binario en bin/devboard
make test          # Ejecuta pruebas con race detector y cobertura
make lint          # Ejecuta golangci-lint
make docker-up     # Inicia PostgreSQL y Redis
make docker-down   # Detiene los servicios
make docker-logs   # Sigue los logs de los servicios
make migrate-up    # Aplica las migraciones pendientes
make migrate-down  # Revierte la última migración
make generate      # Regenera consultas con sqlc y ejecuta go generate
make docs          # Regenera la documentación OpenAPI
make tidy          # Ordena y verifica las dependencias
```

Para crear una migración nueva:

```bash
make migrate-create name=nombre_de_la_migracion
```

## Pruebas

Ejecuta la suite completa con detección de condiciones de carrera y reporte de cobertura:

```bash
make test
```

## Desarrollo

Antes de proponer cambios:

1. Crea una rama a partir de `main`.
2. Mantén las responsabilidades dentro de la capa correspondiente.
3. Agrega o actualiza pruebas para el comportamiento modificado.
4. Ejecuta `make test` y `make lint`.
5. Si modificas anotaciones de la API, regenera la documentación con `make docs`.

## Autor

Desarrollado por [Ricardo Cuéllar](https://github.com/ricardocuellar).

