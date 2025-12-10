# Multi-Tenant Hotel Management System (Go + Gin + PostgreSQL + Redis)

A minimal multi-tenant backend service built using Go, Gin, PostgreSQL (pgx), and Redis.
Each request must include a X-Business-ID header to ensure strict tenant-level data isolation.

## Project Structure

`cmd/server/`        → Application entrypoint

`internal/app/`       → DI, router setup

`internal/tenant/`    → Tenant validation middleware

`internal/hotels/`    → Hotel logic (handler/service/repo)

`internal/rooms/`     → Room logic

`internal/db/`        → Postgres connection (pgx)

`internal/cache/`     → Redis client

`migrations/`         → SQL schema + seed data


## Requirements

- Go 1.22+
- PostgreSQL
- Redis
- Environment variables for configuration (or config file via Viper)

## Database Setup

1. Create a PostgreSQL database.
2. Run the migration files (includes seeding 100 business ids and names)


## Running the Application

`go run ./cmd/server`

The server will start at `http://localhost:8080`

### Multi-Tenancy Notes

- Every API request must include: `X-Business-ID: <business-uuid>`

The middleware:
- Checks Redis for the tenant.
- Falls back to PostgreSQL if not cached.
- Injects the tenant ID into gin.Context so all logic is scoped properly.