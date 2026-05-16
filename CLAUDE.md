# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Commands

```bash
# Start the server (PostgreSQL must be running first)
go run cmd/auth/main.go

# Start local PostgreSQL + Adminer via Docker
docker compose -f docker/docker-compose.yml up -d

# Run all tests
go test ./...

# Run a single package's tests
go test ./handlers/auth/...

# Run code generators (models, queries, inputs — run after changing data-model specs)
go generate ./cmd/auth/...
```

## Architecture

This is a Go REST API with a custom HTTP router, layered into handlers → models → database.

**Entry point:** `cmd/auth/main.go` — connects to PostgreSQL and starts the HTTP server on `localhost:8080`. Routes are registered in `web/http-server.go`.

**Router:** `handlers_utils/router.go` implements a custom router (not `net/http` mux). Handlers are registered via `router.AddAuthenticated` or `router.AddNoAuthenticated`. Authentication uses JWT middleware in `handlers_utils/request-authorizer.go`.

**Handler pattern:** Each handler is a struct that satisfies an interface with `IsValid`, `IsEnabled`, `GetQuery`, and `ServeHTTP`. Handlers receive a `database.IQuery` for all DB access.

**Database layer:** `database/interfaces.go` defines `IQuery` (QueryOne, QueryMany, QueryManyPaged, Exec). The PostgreSQL implementation is in `database/postgres.go`. All handlers depend on this interface — tests mock it to avoid hitting a real DB.

**Code generation:** A significant portion of the codebase under `generated/` is auto-generated from JSON specs in `data-model/`. The generators live in `generators/` and are invoked via `go generate`. Do not edit files under `generated/` directly — edit the specs or templates instead. Generators produce: `generated/models/` (struct types), `generated/queries/` (SQL queries), `generated/inputs/` (request input types), and `generated/handlers/users/` (CRUD handlers).

**Custom types:** `types/` defines nullable wrappers (`types.String`, `types.Int`, etc.) and domain enums used throughout models, inputs, and scanners. `scanners/` contains corresponding database row scanners for each type.

**Auth flow:** `handlers/auth/` contains hand-written auth handlers (login, signup, verify-email, refresh-token, change-password, resend-code). JWT helpers are in `helpers/`. Email is sent via Resend (`API_RESEND_KEY` env var).

**Workflow system:** `pkg/workflow/` contains a workflow execution engine with D2 diagram rendering. `handlers/social/` exposes the `POST /social/add-step/{nameId}` endpoint.

## Database

Default connection: `host=localhost port=5432 dbname=stories user=postgres password=stories`. Docker Compose uses `user=stories password=stories`. The `cmd/auth/main.go` hardcodes the postgres user — adjust `DatabaseConfig` there or set `DATABASE_URL`.

PostgreSQL must have the PostGIS extension enabled.

## Environment Variables

| Variable | Purpose |
|---|---|
| `API_RESEND_KEY` | Resend email service API key |
| `URL_VERIFY_EMAIL` | Link sent in verification emails |
| `URL_CHANGE_PASSWORD` | Link sent in password-reset emails |
| `JWT_SECRET` | JWT signing secret (has a fallback default in constants) |
