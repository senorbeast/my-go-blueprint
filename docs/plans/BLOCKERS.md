# Blockers Register

This register separates external blockers from the concrete work they prevent.
Update it whenever a blocker changes state; keep completed implementation work in [IMPLEMENTATION-PROGRESS.md](./IMPLEMENTATION-PROGRESS.md).

## Active Blockers

### B-001 — Docker engine unavailable

- Status: active
- Evidence: Docker cannot reach `//./pipe/docker_engine` in this workspace.
- Owner: environment
- Unblock condition: Docker Desktop is running and the current user can access the Docker engine.

#### Blocked todos

- Run PostgreSQL and MySQL migration `up -> reset/rollback -> up` for the current generated fixtures.
- Run each seed profile twice against both databases and assert no duplicate records.
- Add and execute Testcontainers coverage for auth, organization/RBAC, CMS, CRM, jobs, and rate-limiter adapters.
- Run composed HTTP tests with real repositories: auth cookies/CSRF, organization and CMS/CRM tenant/RBAC allow-deny behavior, and rate-limit headers.
- Verify durable SQL worker claim, lease, retry, dead-letter, idempotency, and restart recovery.
- Verify shared-database limiter atomicity and multi-instance tenant isolation.
- Start complete generated stacks and run authenticated Playwright flows.

### B-002 — GitHub HTTPS access unavailable

- Status: active
- Evidence: upstream import failed with `Failed to connect to github.com:443`.
- Owner: environment/network
- Unblock condition: HTTPS access to GitHub is available.

#### Blocked todos

- Import and attach the upstream `melkeydev/go-blueprint` history while preserving the recorded compatible implementation and attribution.

## Resolved / No Longer Blocking

### sqlc generation

Generated projects pin `github.com/sqlc-dev/sqlc v1.30.0` and run `go tool sqlc generate` before backend tests. The initial tool download required network access, but sqlc generation is no longer a design blocker for adapter work.
