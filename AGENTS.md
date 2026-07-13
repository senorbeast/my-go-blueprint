# Repository Guidance

## Architecture

- Build one generated-project vertical slice at a time.
- Domain code must not import HTTP, database, logging, JWT, or queue implementations.
- Application code depends on domain and platform interfaces; adapters point inward.
- Define interfaces only at real seams. Keep generic CRUD/listing machinery internal.
- Bootstrap is the only composition root; do not add mutable package globals.
- Keep generated feature code grouped by feature and split by behavior.

## Generated Files

- Never hand-edit generated sqlc or hey-api output.
- Generator-owned files must be declared in the generated manifest.
- `add feature` must refuse to overwrite changed user-owned files.
- PostgreSQL and MySQL assets stay dialect-specific.
- Use one table per reversible Goose migration.

## Verification

- Run `go test ./...` for generator changes.
- Generate at least one fixture for template changes.
- Run generated backend/frontend checks when their templates change.
- Update `docs/plans/IMPLEMENTATION-PROGRESS.md` after each completed slice, design deviation, failed check, or handoff.

## Collaboration

- Respect the ownership table in the progress ledger.
- Do not edit another agent's reserved paths.
- Propose shared-contract changes to the primary agent.
