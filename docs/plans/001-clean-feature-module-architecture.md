# Clean Feature Module Architecture

## Principles

- Build deep modules: small interfaces with substantial behavior behind them.
- Put interfaces at real seams and define them in the domain/application code that consumes them.
- Treat HTTP, SQL, JWT, logging, queues, clocks, IDs, and external effects as adapters.
- Keep deterministic validators, mappers, cursor codecs, and pure helpers concrete.
- Avoid mutable package globals and service locators; wire typed process-lifetime instances in one composition root.
- Test modules through the same interfaces used by production callers.

## Dependency Direction

```text
HTTP/Huma ---------┐
SQL/sqlc ----------┼--> application --> domain
jobs/cron ---------┤
external adapters -┘

bootstrap constructs and connects adapters
```

Domain imports no infrastructure. Application imports domain and platform interfaces, never concrete adapters. Features may depend on another feature only through an explicit domain/application interface, never its HTTP or repository adapter. Architecture tests reject forbidden imports and cycles.

## Generated Layout

```text
backend/
  cmd/{api,worker,scheduler,migrate,seed}/
  internal/bootstrap/
  internal/platform/{auth,authorization,database,errors,http,jobs,logging,pagination,ratelimit,validation}/
  internal/features/<feature>/
    domain/
    application/
    adapters/http/
    adapters/repository/
    queries/
    jobs/
    seeds/
    feature.go
  migrations/{postgres,mysql}/
```

Split files by behavior and avoid generic `utils`, `helpers`, or `common` packages.

## Feature Contract

Each feature exposes one constructor and a compact module interface for route registration, job registration, and seeding. Its descriptor declares dependencies, route groups, permissions, jobs, migrations, seeds, OpenAPI tags, and frontend contributions. Generator validation rejects duplicate routes/operation IDs, cycles, and incomplete declarations.

Domain repositories expose meaningful typed operations rather than a public universal CRUD interface. Reusable CRUD, mapping, filtering, sorting, and pagination mechanics remain internal implementation primitives.

Use typed IDs, email, slug, permission, status, cursor, page-limit, and optional patch values. Constructors validate values so invalid domain objects cannot be created through public paths.

## Routes and HTTP

- Provide nested logical route groups with prefixes, versions, OpenAPI tags, middleware, authentication mode, tenant requirements, permissions, and rate policies.
- Huma typed inputs/outputs are the OpenAPI source.
- Handlers receive validated input, obtain a transport-independent principal, invoke an application interface, and return typed output/errors.
- Startup validates route uniqueness, operation IDs, permissions, and middleware order.
- Middleware order covers request IDs, recovery, security/body limits, access logging, CORS, authentication, CSRF, rate limiting, tenant resolution, authorization, and timeouts.

## Common Contracts

- Cursor pagination uses opaque, versioned, authenticated tokens containing ordered values, stable ID tie-breaker, sort direction, query fingerprint, and tenant context.
- Filters and sorts use typed enums; raw column names or SQL fragments never cross application interfaces.
- Central errors classify validation, unauthenticated, forbidden, not found, conflict, rate limited, dependency unavailable, and internal failures.
- HTTP maps errors to RFC 9457 responses with stable codes, safe messages, field violations, and request IDs.
- Zap logging uses a narrow contextual interface and redacts cookies, credentials, tokens, and configured sensitive fields.

## Composition, Database, and Transactions

The composition root owns configuration, logger, database pool, transaction runner, clock, ID generator, session manager, authorizer, limiter, queue, job registry, route registrar, and feature modules. Constructors do not perform hidden environment reads.

Generate PostgreSQL with `pgx/v5` or MySQL with `database/sql`. Keep sqlc records inside repository adapters. Translate driver errors into typed errors. Application modules never receive raw pools, transactions, or sqlc query objects. A transaction runner supplies transaction-scoped repositories and supports atomic state change plus job enqueueing.

Use one reversible Goose migration per table, including join and platform tables. Tenant ID is mandatory in every tenant-scoped query and database constraints reinforce isolation.

## Authentication and RBAC

- Access JWTs are short-lived Secure HttpOnly cookies.
- Refresh tokens are opaque high-entropy values; only hashes are stored.
- Rotation is atomic, reuse revokes the family, and current/all-session logout is supported.
- Include signup, login, refresh, current user, organization switching, session management, reset, and verification hooks.
- Support cross-subdomain credentials with explicit origins, cookie domain, SameSite policy, and CSRF protection.
- Typed principal extraction is independent of Echo context.
- Organization-scoped memberships, roles, permissions, and assignments support route and resource policies.
- Application behavior rechecks authorization when invoked outside HTTP.

## Jobs, Cron, and Rate Limiting

SQL jobs provide versioned payloads, transactional enqueueing, scheduling, skip-locked claims, leases, heartbeats, retries/backoff, idempotency, cancellation, dead letters, and graceful shutdown. API, worker, and scheduler share application modules but run as separate commands. Cron obtains a SQL leader lease and only enqueues jobs.

SQL rolling-window limiting uses atomic buckets shared by instances. Keys may include IP, user, organization, operation, or sensitive auth action. Decisions expose remaining, limit, retry, and reset values. Sensitive authentication routes fail closed; ordinary routes fail open with structured logging. Cleanup runs as a scheduled job and storage stays replaceable behind the limiter interface.

## Frontend

Generate feature-local React code with TanStack Router, TanStack Query, Tailwind, Zustand, Vitest, Testing Library, and Playwright. hey-api output is generated from Huma OpenAPI and never edited manually. TanStack Query owns server state; Zustand owns client-only UI state. Protected routes, tenant switching, permission-aware navigation, CSRF, problem responses, and single-flight refresh are built in.

CMS supplies tenant-scoped posts, categories, tags, workflow, permissions, seeds, and admin UI. CRM supplies companies, contacts, pipelines, stages, deals, activities, permissions, seeds, and admin UI. Each table has its own migration and each pack owns its complete vertical slice.
