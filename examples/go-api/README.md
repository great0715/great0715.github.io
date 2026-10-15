# Release Lab — Go API reference

A small, **real Go HTTP service** included as an educational portfolio code sample. This is separate from the browser-based Release Lab simulator. It is **not connected to any production system**.

## Why this exists

It demonstrates API routing, input validation, JSON responses, thread-safe in-memory state, health checks, graceful shutdown, unit tests, containerization, and a sample CI workflow — the backend/DevOps side of a full-stack portfolio.

## Run locally

Requires Go 1.23 or newer.

```bash
go run .
```

Then:

```bash
curl http://localhost:8080/healthz
curl -X POST http://localhost:8080/api/releases \
  -H 'Content-Type: application/json' \
  -d '{"service":"website","environment":"staging","version":"v1.2.0"}'
curl http://localhost:8080/api/releases
```

### Tests

```bash
go test -race ./...
go vet ./...
```

### Docker

```bash
docker build -t release-lab-api .
docker run --rm -p 8080:8080 release-lab-api
```

The `.github/workflows/ci.yml` file is a **template for this example as a standalone GitHub repository**. To activate the workflow, put the contents of `examples/go-api/` at the root of its own repository. It does not automatically run when nested inside a static portfolio repository.

## Routes

| Endpoint | Action |
|---|---|
| `GET /healthz` | Returns `{"status":"ok"}` |
| `GET /api/releases` | Lists in-memory release records |
| `POST /api/releases` | Accepts service, environment (`staging` or `production`) and version; creates a queued record |

## Deliberate limitations

This is a learning/reference service, not a production deployment system. Data disappears on restart, there is no authentication, no database, and nothing is deployed to cloud infrastructure. The POST route creates a record **only**; it never triggers an actual release. Before real-world use it would need authorization, durable storage, rate limiting, audit logging, and deployment safeguards.
