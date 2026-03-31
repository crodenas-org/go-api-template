# hello-world-go

A production-template Go REST API built incrementally. See [PLAN.md](PLAN.md) for the full build roadmap.

## Requirements

- Go 1.25+
- PostgreSQL (see setup below)
- `goose` — `go install github.com/pressly/goose/v3/cmd/goose@latest`
- Docker or Podman (for containerized local dev)

## Setup

**1. Database**

Create the database, users, and schema (one-time):

```bash
# Create DB and users (requires postgres superuser)
psql -h <host> -U postgres -f db/setup.sql

# Run migrations
make migrate-up
```

**2. Environment**

Create a `.env` file in the project root:

```env
DATABASE_URL=postgres://hw_app:<password>@<host>:5432/hello_world_go
MIGRATION_URL=postgres://hw_admin:<password>@<host>:5432/hello_world_go
AZURE_TENANT_ID=<tenant-id>
AZURE_CLIENT_ID=<api-client-id>      # API app registration
MCP_ADDR=:8081                        # optional, default :8081
```

**3. Run**

```bash
make run
```

## Docker / Podman

Runs the API in a container, connecting to the database configured in `.env`.

```bash
make docker-up    # build image and start in background
make docker-down  # stop

make podman-up    # same, using podman
make podman-down
```

## API Docs

Swagger UI is available at `http://localhost:8080/docs/index.html` when the server is running.

Regenerate after changing handler annotations:
```bash
make docs
```

## Make Targets

| Target | Description |
|---|---|
| `make run` | Run the API server |
| `make build` | Compile binary to `bin/api` |
| `make clean` | Remove compiled binaries |
| `make test` | Run all tests |
| `make test-verbose` | Run tests with verbose output |
| `make migrate-up` | Apply pending migrations (homelab DB) |
| `make migrate-down` | Roll back last migration (homelab DB) |
| `make migrate-status` | Show migration status (homelab DB) |
| `make docs` | Regenerate OpenAPI spec from annotations |
| `make docker-build` | Build Docker image |
| `make docker-up` | Start app via docker compose |
| `make docker-down` | Stop docker compose services |
| `make podman-build` | Build image with podman |
| `make podman-up` | Start app via podman compose |
| `make podman-down` | Stop podman compose services |

## API

### Health

```bash
# Hello world
curl http://localhost:8080/hello
```

### Items

```bash
# List all items
curl http://localhost:8080/items

# Create an item
curl -X POST http://localhost:8080/items \
  -H "Content-Type: application/json" \
  -d '{"name": "my item"}'

# Get item by ID
curl http://localhost:8080/items/1

# Get a non-existent item (404)
curl http://localhost:8080/items/999
```

## Project Structure

```
cmd/
  api/            # REST API entrypoint
  mcp/            # MCP SSE server entrypoint
db/
  migrations/     # goose SQL migrations
  setup.sql       # one-time DB and user setup
internal/
  handler/        # HTTP handlers
  middleware/     # auth and role enforcement
  model/          # shared data types
  repository/     # database queries
  server/         # router setup
infra/
  entra/          # Azure app registration scripts
```
