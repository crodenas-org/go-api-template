# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Commands

```bash
# Run
make run              # Start API server on :8080
make run-mcp          # Start MCP SSE server on :8081

# Build
make build            # Compile to bin/api

# Test
make test             # go test ./...
make test-verbose     # go test -v ./...

# Database migrations (requires MIGRATION_URL env var)
make migrate-up       # Apply pending migrations
make migrate-down     # Roll back last migration
make migrate-status   # Show migration state

# Docs (after modifying Swagger annotations in handlers)
make docs             # Regenerate docs/ from swag annotations

# Docker
make docker-up        # docker compose up --build -d
make docker-down      # docker compose down
```

To run a single test:
```bash
go test ./internal/handler/... -run TestHelloHandler
go test ./internal/server/... -run TestRouter
```

## Architecture

**Entry points:** `cmd/api/main.go` (REST API) and `cmd/mcp/main.go` (MCP SSE server). `main.go` wires together: pgxpool DB connection, OIDC provider (Azure AD), and the chi router from `internal/server`.

**Request flow:**
1. `internal/server/server.go` — Chi router; routes are grouped by auth requirement. Public routes get no middleware. Authenticated routes run `Authenticate()`, then role-specific sub-groups run `RequireRole(role)`.
2. `internal/middleware/auth.go` — Validates the Bearer JWT against the Azure AD OIDC provider, extracts claims, stores them in request context.
3. `internal/handler/` — Handlers call the repository, marshal JSON responses.
4. `internal/repository/item.go` — All SQL lives here; uses pgx directly (no ORM).

**Two app roles:** `items.read` (required for GET /items, GET /items/{id}) and `items.write` (required for POST /items). Roles come from Azure Entra ID app registration (`infra/entra/setup.sh`).

**Two DB users:** `hw_admin` (DDL, used only by goose migrations via `MIGRATION_URL`) and `hw_app` (DML only, used by the running app via `DATABASE_URL`).

## Environment

Copy `.env.compose` and fill in passwords for local Docker dev. For local non-Docker dev, provide `.env` with:
```
DATABASE_URL=postgres://hw_app:<pass>@localhost:5432/hello_world_go
MIGRATION_URL=postgres://hw_admin:<pass>@localhost:5432/hello_world_go
AZURE_TENANT_ID=<tenant>
AZURE_CLIENT_ID=<client>
```

The DB schema is in `db/setup.sql` (one-time provisioning) + `db/migrations/` (goose-managed). Swagger UI is available at `http://localhost:8080/docs/index.html` when running.

## CORS

Configured in `internal/server/server.go` via `github.com/go-chi/cors`. Currently allows `http://localhost:5173` (Vite dev server). When adding a production frontend origin, add it to the `AllowedOrigins` slice — do not use wildcards since requests carry `Authorization` headers.

## Adding new endpoints

1. Add handler method in `internal/handler/` with swag annotations.
2. Add repository method in `internal/repository/` if DB access is needed.
3. Register the route in `internal/server/server.go` under the correct auth group.
4. Run `make docs` to regenerate Swagger docs.
