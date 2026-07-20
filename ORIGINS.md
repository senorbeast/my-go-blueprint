# Project Origins

## Summary

Opinionated Go Blueprint is an independently implemented, full-stack project generator inspired by [melkeydev/go-blueprint](https://github.com/melkeydev/go-blueprint). It retains the familiar interactive and flag-based project-creation workflow while expanding the generated application into a production-oriented Go and React foundation.

## Upstream relationship

- **Upstream project:** [melkeydev/go-blueprint](https://github.com/melkeydev/go-blueprint)
- **Upstream repository:** `https://github.com/melkeydev/go-blueprint.git`
- **Intended reference point:** upstream `main`, recorded 2026-07-13
- **License relationship:** the upstream project is MIT licensed; this repository preserves attribution and remains compatible with that licensing approach.

The upstream history was not imported during initial development because the environment could not reach GitHub over HTTPS. This repository therefore documents its relationship to upstream without claiming that it is a copied or history-preserving fork.

## What this project adds

The generator provides a feature-composed, clean-architecture application baseline with:

- Go services using Huma/OpenAPI, sqlc, Goose, PostgreSQL or MySQL, Zap, and Docker.
- React/Vite frontend templates using TanStack Router and Query, hey-api, Tailwind, and Zustand.
- Optional auth, multi-tenant RBAC, SQL-backed jobs and cron, CMS, and CRM feature packs.
- Dependency-aware feature selection, generated manifests, safe feature additions, verification, and diagnostics.
- Dialect-specific migrations, deterministic seeds, tests, and generated-project documentation.

For implementation status and known limitations, see [the progress ledger](./docs/plans/IMPLEMENTATION-PROGRESS.md). The original upstream note remains in [UPSTREAM.md](./UPSTREAM.md).
