# Generated Project Behavioral Testing

## Objective

Prove generated projects work as applications, not only that templates render. Tests exercise public module interfaces and public HTTP/browser behavior while avoiding coupling to internal helpers.

## Generated-Project Tests

- Domain/application tests cover value invariants, policies, commands, queries, validation, typed errors, and cursor behavior.
- Testcontainers repository tests apply real Goose migrations and run actual sqlc queries for the selected database.
- HTTP tests run the composed Echo/Huma app with real repositories and verify cookies, middleware, validation, errors, auth, RBAC, tenants, and rate limits.
- Worker tests verify enqueue, claim, lease, heartbeat, retry, idempotency, dead letter, restart recovery, cron leadership, and graceful shutdown.
- Frontend tests use Vitest/Testing Library; key live-stack behavior uses Playwright.
- Architecture tests reject dependency violations and incomplete route/permission declarations.
- Contract tests export OpenAPI, regenerate hey-api, typecheck it, and reject stale generated output.

## Acceptance Harness

For each declared scenario:

1. Generate into an isolated temporary directory.
2. Verify its manifest, feature graph, files, migrations, and docs.
3. Compile Go, run sqlc, export OpenAPI, regenerate hey-api, and build React.
4. Start uniquely named Docker Compose resources with dynamic ports.
5. Wait on readiness checks, migrate, and seed.
6. Start API, worker, scheduler, and frontend.
7. Execute HTTP and Playwright scenarios.
8. Restart API and workers and verify persistence, sessions, leases, and recovery.
9. Capture logs, diagnostics, OpenAPI, reports, traces, and screenshots on failure.
10. Always remove containers, networks, volumes, credentials, and temporary projects.

## Seeds

Seeds are deterministic, versioned, idempotent, and environment-gated. `minimal` creates a usable user, organization, owner membership, roles, and permissions. `demo` adds representative CMS, CRM, and job records. Run seeds twice and verify no duplication. Verify primarily through authenticated public behavior; use direct SQL only for constraints, token hashes, queue state, and migration invariants.

## Behavioral Coverage

- Health, readiness, version, OpenAPI, and Swagger UI.
- Signup, login, refresh rotation/replay, logout, reset, verification, and session revocation.
- Cookie attributes, cross-subdomain CORS, and CSRF rejection/acceptance.
- Organization switching, tenant isolation, and RBAC allow/deny behavior.
- Cursor traversal, invalid/tampered cursors, validation, conflicts, safe errors, and log redaction.
- Rolling-window enforcement, reset, headers, concurrency, and tenant isolation.
- CMS workflow/taxonomy and CRM pipeline/activity behavior using seeded data.
- Transactional jobs, retries, leases, dead letters, idempotency, restart recovery, and cron deduplication.
- Migration up, rollback to zero, remigration, and seed replay.
- Browser auth, tenant switching, permission-aware navigation, CRUD, pagination, seeded data, and refresh recovery.

## Matrix and Cadence

Pull requests run PostgreSQL default, MySQL complete/demo, minimal backend-only, auth/RBAC, jobs/cron, focused CMS, focused CRM, and generator/architecture/contract tests.

Nightly/manual CI runs pairwise combinations across database, frontend, auth, RBAC, jobs, cron, CMS, CRM, and seeds; canonical minimal and complete stacks on both databases; repeated concurrency-sensitive tests; and Playwright key flows on canonical frontend stacks.

Commit scenario definitions and a coverage report. A new CLI option must fail CI until included in matrix coverage.

## Diagnostics and Acceptance

Failures preserve the generator command, feature graph, manifest, container/application logs, migration state, database diagnostics, OpenAPI diff, SDK output, test reports, and Playwright artifacts with secrets redacted.

Every README command is executed by automation, both databases pass, valid feature combinations work, seeds remain idempotent, restarts recover correctly, generated contracts are reproducible, and no Docker resources or credentials remain after tests.
