# hello-world-go — Build Plan

A production-template Go REST API built incrementally. Each step is a self-contained learning milestone.

## Goals
- MS OAuth2 authentication + Azure AD app roles
- Postgres database
- OpenAPI/Swagger docs served from the app
- MCP server companion
- GitHub Actions CI/CD pipeline
- Containerized deployment on AWS ECS
- Tests where they add real value
- Makefile targets grown alongside the app

---

## Steps

### ✅ Step 0 — Hello World
- `main.go` with a single `GET /hello` endpoint (stdlib `net/http`)
- `go.mod`, `.gitignore`, `Makefile` (run, build, clean)

### ✅ Step 1 — Routing & Package Structure
- Introduce [chi](https://github.com/go-chi/chi) router
- Restructure into a proper package layout:
  ```
  cmd/api/        # main entrypoint
  internal/
    handler/      # HTTP handlers
    server/       # server setup and routing
  ```
- Makefile: no new targets needed yet

### ✅ Step 2 — Middleware
- Request logging (method, path, status, latency)
- Request ID generation
- Centralized error handling
- Makefile: no new targets needed yet

### ✅ Step 3 — Tests
- Handler tests using stdlib `net/http/httptest`
- Middleware tests
- Makefile: `make test`, `make test-verbose`

### ✅ Step 4 — Database
- Postgres with [`pgx`](https://github.com/jackc/pgx)
- Migrations with [`goose`](https://github.com/pressly/goose)
- Repository pattern to keep DB logic out of handlers
- Makefile: `make migrate-up`, `make migrate-down`, `make migrate-status`

### ✅ Step 5 — MS OAuth2 / Azure AD Roles
- OIDC authentication via Azure AD
- App roles mapped to middleware-enforced access control
- Same conceptual flow as FastAPI + `fastapi-azure-auth`
- Makefile: no new targets needed yet

### Step 6 — API Docs
- OpenAPI spec (hand-authored or generated)
- Swagger UI served at `/docs`
- Makefile: `make docs` if generation is needed

### Step 7 — MCP Server
- Companion MCP server that exposes app functionality as tools
- Co-located in this repo under `cmd/mcp/`
- Makefile: `make run-mcp`

### Step 8 — Containerize
- `Dockerfile` (multi-stage build)
- `docker-compose.yml` for local dev (app + Postgres)
- Makefile: `make docker-build`, `make docker-up`, `make docker-down`

### Step 9 — GitHub Actions + AWS ECS
- CI: lint, test, build, push image to ECR
- CD: deploy to ECS on merge to main
- Self-hosted runner containers (consistent with other projects)
- Makefile: `make lint`
