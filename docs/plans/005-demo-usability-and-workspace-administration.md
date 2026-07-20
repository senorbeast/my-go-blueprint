# Demo Usability and Workspace Administration Plan

## Goal

Make the generated all-feature demo demonstrably usable: a demo owner can create CRM records, manage organization members and their roles, and observe a real jobs/cron workflow. Keep customer-facing sales work in CRM and workspace access control in the organizations/RBAC slice.

## Boundaries

- CRM owns companies, contacts, pipelines, stages, deals, and activities.
- Workspace administration owns organization members, roles, and permissions. It is IAM/workspace administration, not ERP.
- Domain packages remain independent of HTTP, database, JWT, and queue implementations.
- Generator changes are made only in source templates; sqlc and hey-api generated output is never hand-edited.

## Delivery slices

1. **CRM create workflows (P0)**
   - Surface the existing company and contact create APIs in the CRM UI, then progressively expose pipelines, stages, deals, and activities with permission-aware forms.
   - Extend contacts with an optional, tenant-validated company link before presenting the company selector.
   - Prove each write through frontend request/UI tests and generated backend authorization tests.

2. **Workspace administration (P0)**
   - Add organization-scoped member and role listing, then add member provisioning/invitation, membership removal, role creation, permission grants, and replacement role assignments as small vertical slices.
   - Use explicit `organizations.members.*` and `organizations.roles.*` permissions. Preserve an immutable owner policy so the UI cannot accidentally remove the only administrator or manufacture platform-wide access.
   - Add a `/workspace` admin UI with member and role tabs only for authorized callers.

3. **Demo data (P0)**
   - Seed a second deterministic organization member and useful roles such as sales manager and editor, along with their least-privilege grants.
   - Document all development-only demo credentials and prove idempotence for PostgreSQL and MySQL.

4. **Jobs and cron demonstration (P0)**
   - Register at least one real, idempotent demo handler and schedule it through the existing queue seam.
   - Make worker startup report a nonzero handler count and prove scheduled work reaches completion; do not present empty queue scaffolding as an operational demo.

5. **Runtime and generated-client hygiene (P1)**
   - Remove the deprecated direct `@hey-api/client-fetch` dependency by updating the hey-api configuration and generated client contract together.
   - Refresh compatible direct test-tool dependencies where they are the owner of transitive deprecation warnings; do not suppress warnings.
   - Expand platform tests for `/healthz`, `/readyz`, `/docs`, and OpenAPI completeness, and keep README URLs accurate for custom listener addresses.

## Verification gate

After every slice: root `go test ./...`, `git diff --check`, and a fresh generated fixture. For changed backend/frontend templates, run the generated backend tests and `pnpm api:generate`, `pnpm typecheck`, `pnpm test`, and `pnpm build`. The final all-feature PostgreSQL and MySQL demo must migrate and seed successfully; the final browser flow is owner login → CRM record creation → workspace member/role action → scheduled job evidence.

## Current evidence

The all-feature PostgreSQL demo migrated to version 35, seeded successfully, and served the frontend, readiness, owner login, CMS list, and CRM company list. Startup also revealed two remediation targets: the frontend dependency deprecation warning and a worker with zero registered handlers.
