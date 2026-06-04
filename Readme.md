# Task Management API

![CI](https://github.com/NguyenNH36/task-management-api/actions/workflows/ci.yml/badge.svg)

A production-style Task Management API built with Go, Gin, OpenAPI 3.0, and oapi-codegen.

## Features

- Create, read, update, and delete tasks
- Health check endpoint
- OpenAPI 3.0 spec with request validation
- PostgreSQL persistence via pgx
- Database migrations with golang-migrate
- Structured JSON logging with request IDs
- Layered architecture (handler/service/repository)

## Tech Stack

- Go
- Gin framework
- OpenAPI 3.0 + oapi-codegen
- PostgreSQL (pgx)
- golang-migrate
- slog for structured logging
- go test + Testify
- go.uber.org/mock (mockgen)

## API Endpoints

| Method | Endpoint | Description |
|---|---|---|
| GET | `/health` | Health check |
| POST | `/tasks` | Create task |
| GET | `/tasks` | Get all tasks |
| GET | `/tasks/{id}` | Get task by ID |
| PUT | `/tasks/{id}` | Update task |
| DELETE | `/tasks/{id}` | Delete task |

## Database Migrations

This project uses `golang-migrate` for PostgreSQL schema migrations.

Migration files are stored in:

```text
migrations/
```

## Configuration

This project uses environment variables for configuration.

Create a local `.env` file:

```env
PORT=8080
DATABASE_URL=postgres://task_user:task_password@localhost:5432/task_db?sslmode=disable
GIN_MODE=debug
```

## OpenAPI

The OpenAPI spec lives in [api/openapi.yaml](api/openapi.yaml) and generated handlers live in [internal/api/task_api.gen.go](internal/api/task_api.gen.go).

## Run Locally

1. Start PostgreSQL (Docker Compose):

```bash
docker compose up -d
```

2. Apply migrations:

```bash
migrate -path ./migrations -database "$DATABASE_URL" up
```

3. Run the API:

```bash
go run ./cmd/api
```

## Run with Docker

Build the image:

```bash
docker build -t task-api .
```

Run the container (make sure PostgreSQL is reachable):

```bash
docker run --rm -p 8080:8080 -e DATABASE_URL="$DATABASE_URL" -e GIN_MODE=release task-api
```

## Testing

This project uses:

- `go test` for unit testing
- `testify` for assertions
- `mockgen` / `gomock` for generating mocks

Generate mocks:

```bash
mockgen -source="./internal/repository/task_repository.go" -destination="./internal/mocks/mock_task_repository.go" -package=mocks
```

Run tests:

```bash
go test ./...
```


## Linting

This project uses `golangci-lint` for static analysis and code quality checks.

Run lint:

```bash
golangci-lint run
```

## Logging and Request Tracing

This project uses structured JSON logging with Go `slog`.

Each request receives an `X-Request-ID`.

If the client sends `X-Request-ID`, the API reuses it.  
If not, the API generates a new UUID request ID.

Example:

```bash
curl -i http://localhost:8080/health -H "X-Request-ID: test-request-123"
```