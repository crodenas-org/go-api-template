# hello-world-go

A production-template Go REST API built incrementally. See [PLAN.md](PLAN.md) for the full build roadmap.

## Requirements

- Go 1.22+
- PostgreSQL (see setup below)
- `goose` — `go install github.com/pressly/goose/v3/cmd/goose@latest`

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
```

**3. Run**

```bash
make run
```

## Make Targets

| Target | Description |
|---|---|
| `make run` | Run the API server |
| `make build` | Compile binary to `bin/api` |
| `make clean` | Remove compiled binaries |
| `make test` | Run all tests |
| `make test-verbose` | Run tests with verbose output |
| `make migrate-up` | Apply pending migrations |
| `make migrate-down` | Roll back last migration |
| `make migrate-status` | Show migration status |

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
  api/            # main entrypoint
db/
  migrations/     # goose SQL migrations
  setup.sql       # one-time DB and user setup
internal/
  handler/        # HTTP handlers
  model/          # shared data types
  repository/     # database queries
  server/         # router and middleware setup
```
