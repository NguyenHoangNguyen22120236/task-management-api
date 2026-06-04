# Task Management API

![CI](https://github.com/YOUR_GITHUB_USERNAME/task-management-api/actions/workflows/ci.yml/badge.svg)

A production-style Task Management API built with Go, Gin, OpenAPI 3.0, and oapi-codegen.

## Current Scope

This project currently supports:

- Create task
- Get all tasks
- Get task by ID
- Update task
- Delete task
- Health check
- OpenAPI-generated routes
- OpenAPI request validation
- In-memory repository
- Service/repository/handler layering
- Unit tests with go test and Testify

## Tech Stack

### Services & Frameworks

- Go
- Gin framework
- OpenAPI 3.0
- JSON Schema
- oapi-codegen v2

### Testing

- go test
- Testify

### Planned Next

- Docker
- GitHub Actions
- DynamoDB or RDS
- AWS SDK for Go
- SQS + DLQ
- ECS or Lambda
- Terraform
- Kong API Gateway
- CloudWatch / Datadog

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

## Testing

This project uses:

- `go test` for unit testing
- `testify` for assertions
- `mockgen` / `gomock` for generating mocks

Generate mocks:

```bash
mockgen -source="./internal/repository/task_repository.go" -destination="./internal/mocks/mock_task_repository.go" -package=mocks
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

## Run Locally

```bash
go run ./cmd/api