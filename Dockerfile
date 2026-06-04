# Stage 1: Build the Go binary
FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o task-api ./cmd/api


# Stage 2: Run the binary
FROM alpine:3.20

WORKDIR /app

COPY --from=builder /app/task-api .
COPY --from=builder /app/api ./api

EXPOSE 8080

CMD ["./task-api"]