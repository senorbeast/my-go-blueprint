# Opinionated Go Blueprint

An opinionated full-stack project generator inspired by and designed as an enhancement of [melkeydev/go-blueprint](https://github.com/melkeydev/go-blueprint).

The generated stack targets Echo, Huma/OpenAPI, sqlc, Goose, PostgreSQL or MySQL, Zap, cookie authentication, multi-tenant RBAC, SQL-backed jobs and cron, React/Vite, TanStack Router and Query, hey-api, Tailwind, Zustand, Docker, and Air.

See [the progress ledger](./docs/plans/IMPLEMENTATION-PROGRESS.md), [architecture plan](./docs/plans/001-clean-feature-module-architecture.md), and [behavioral-testing plan](./docs/plans/002-generated-project-behavioral-testing.md).

## Create a project

Interactive:

```powershell
go run . create
```

Non-interactive full stack:

```powershell
go run . create --name acme --module example.com/acme --database postgres --seed demo --feature cms --feature crm --feature jobs --feature cron
```

Backend-only MySQL project:

```powershell
go run . create --name acme-api --module example.com/acme-api --database mysql --no-frontend --feature auth --feature rbac
```

The CLI resolves feature dependencies, writes `.blueprint/manifest.json`, supports `--dry-run`, and protects customized files during later additions:

```powershell
go run . add feature cms --dir ./acme
go run . verify --dir ./acme
go run . doctor --dir ./acme
```

Generated projects include their own README, architecture guide, Makefile, Air configuration, Docker Compose, migrations, seeds, tests, and GitHub Actions workflow.

## Development

This repository currently requires Go 1.25 or newer.

```powershell
$env:GOCACHE = "$PWD/.cache/go-build"
go test ./...
```

To verify the generator:

```powershell
./scripts/verify.ps1
```

## Upstream and License

See [UPSTREAM.md](./UPSTREAM.md). The upstream project is MIT licensed; this enhancement preserves attribution and will retain compatible licensing.
