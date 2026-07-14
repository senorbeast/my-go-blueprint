# Opinionated Go Full-Stack Blueprint — Master Plan

## Goal

Enhance `melkeydev/go-blueprint` into an opinionated generator for production-oriented, multi-tenant Go and React applications. Generated projects use Echo, Huma/OpenAPI, sqlc, Goose, PostgreSQL or MySQL, Zap, cookie authentication, RBAC, SQL jobs and cron, rolling-window rate limiting, React/Vite, TanStack Router and Query, hey-api, Tailwind, Zustand, Docker, and Air.

The architecture contract is defined in [001-clean-feature-module-architecture.md](./001-clean-feature-module-architecture.md). Generated-project verification is defined in [002-generated-project-behavioral-testing.md](./002-generated-project-behavioral-testing.md). Delivery status is recorded in [IMPLEMENTATION-PROGRESS.md](./IMPLEMENTATION-PROGRESS.md).

## Product and CLI

- Preserve interactive `create` and non-interactive flags.
- Select exactly one database: PostgreSQL (default) or MySQL.
- Support repeatable feature selection for auth, RBAC, jobs, cron, CMS, and CRM.
- Support `none`, `minimal`, and `demo` seed profiles.
- Generate React by default, with a backend-only option.
- Support idempotent `add feature`, dry-run conflict reporting, `doctor`, and `verify`.
- Store a generation manifest and never silently overwrite user-owned files.
- Resolve dependencies: RBAC requires auth; CMS/CRM require auth and RBAC; cron requires jobs.

## Delivery Slices

1. Import upstream and establish repository, documentation, lint, and test baselines.
2. Build deterministic feature-pack and CLI foundations.
3. Generate a runnable Echo/Huma backend and React/Vite frontend.
4. Add dialect-specific sqlc, Goose, Docker, migrations, and seeds.
5. Add users, organizations, authentication, refresh rotation, and RBAC.
6. Add SQL jobs, workers, scheduler, cron leadership, and rate limiting.
7. Add CMS as a complete backend/frontend vertical slice.
8. Add CRM as a complete backend/frontend vertical slice.
9. Add post-generation feature installation and conflict safety.
10. Complete generated-project Docker and browser acceptance coverage.

Each slice must generate, compile, and behaviorally verify a usable project before the next slice begins.

## Feature Completion Gate

Every feature slice must name its user-visible capability and observable acceptance behaviors before implementation. A feature is complete only after its real adapters and routes are composed and all applicable checks pass: repository-wide verification, fresh PostgreSQL and MySQL generated fixtures, migration/reset/seed replay for persistence changes, composed HTTP authorization/tenant/cookie behavior, worker or scheduler behavior, frontend contract/type/unit/build checks, and a live browser flow for user-facing full-stack behavior.

Feature graph changes must extend the canonical acceptance scenarios and pairwise-coverage proof. Verify each feature in an independent valid configuration where possible and in one configuration with its required dependencies or the complete pack. Dependency resolution alone is not runtime evidence. Record the commands, configurations, outcomes, failures, and remaining gaps in the implementation ledger.

## Generated Documentation and Commands

Generate a root README, folder/architecture guide, feature-authoring guide, auth/RBAC guide, jobs/cron guide, testing guide, `.env.example`, Air configuration, Docker Compose, sqlc, Goose, hey-api, and CI configuration.

Document and automate commands to start the database, migrate and roll back, seed, run the backend with Air, run workers and scheduler, run Vite, export OpenAPI, regenerate the SDK, run all test layers, and start/stop the complete stack.

## Completion

- PostgreSQL and MySQL variants compile, migrate, seed, start, and pass behavioral tests.
- Valid feature combinations work independently and in the complete stack.
- Later feature installation is idempotent and protects user code.
- OpenAPI and hey-api generation are reproducible.
- Every documented command is exercised by automation.
- The progress ledger contains evidence for every completed slice.
