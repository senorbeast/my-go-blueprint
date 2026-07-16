# Demo-Ready Generated Applications — Approved Tracer-Bullet Issues

Local Markdown issue set for the approved PRD: [003-demo-ready-generated-applications-prd.md](./003-demo-ready-generated-applications-prd.md).

## 1. Generate a usable application shell and seed contract

### What to build

Generate a shadcn-based public/authenticated application shell that is capability-driven, never blank for a valid feature graph, and includes backend-only documentation where no frontend is selected. Establish the reusable `none`, `minimal`, and `demo` seed contract with a working, documented local demo credential and idempotent representative data.

### Acceptance criteria

- [ ] Every valid graph renders either a usable authenticated shell or an intentional public/first-user experience; no graph reaches a generic empty landing page.
- [ ] Minimal and demo profiles sign in with the generated README credential; replaying seed data adds no duplicates.
- [ ] Generated frontend navigation reflects selected features and current capabilities; backend-only README/OpenAPI documents equivalent flows.
- [ ] Fresh generated PostgreSQL and MySQL fixtures pass seed and browser/HTTP proof appropriate to their frontend selection.

### Blocked by

None — can start immediately.

## 2. Deliver Workspace identity, tenancy, and administration

### What to build

Generate the deep Workspace module with server session context/revocation, safe organization switching, tenant user/member/invitation/role/permission administration, platform super-admin organization/user administration, and audit history. Feature capability registration remains at bootstrap while Workspace owns assignment and access administration.

### Acceptance criteria

- [ ] Session context exposes actor, active organization, available organizations, capabilities, and platform-admin status; logout revokes server access.
- [ ] A valid membership can switch organizations with a rotated/reissued session; invalid or cross-tenant switching is denied.
- [ ] Tenant admin cannot enumerate other tenants, remove the final owner, or assign a privilege they cannot administer; platform super-admin is not obtainable through tenant roles.
- [ ] Generated UI/API supports workspace and platform administration and records auditable changes.

### Blocked by

- Issue 1 — generated shell and seed contract.

## 3. Deliver a complete CMS editorial workflow

### What to build

Generate a tenant-scoped CMS workflow for content dashboard/listing/search/filter/pagination, categories/tags, create/edit/preview/publish/schedule/archive/restore, revision restore, media-library use, capability enforcement, and audit history.

### Acceptance criteria

- [ ] A seeded authorized editor can manage taxonomy and content through the generated UI/API, including filter, preview, archive/restore, and revision restore.
- [ ] Publishing and media/taxonomy/revision actions enforce their explicit capabilities and are tenant-safe and audited.
- [ ] When jobs/cron are selected, scheduled publishing reaches the real queue/scheduler and exposes outcome/notification status.
- [ ] HTTP and browser tests prove allowed, denied, and lifecycle behavior on fresh dialect fixtures.

### Blocked by

- Issue 1 — generated shell and seed contract.
- Issue 2 — Workspace authorization and audit behavior.

## 4. Deliver a complete CRM sales workflow

### What to build

Generate a tenant-scoped CRM workflow for companies, contacts, pipelines, stages, deals, activities/tasks, assignments, notes/associations, deal board movements, searchable lists, reports, dashboard metrics, and capability-aware audit history.

### Acceptance criteria

- [ ] A seeded authorized user can complete a company/contact/deal/activity workflow and view the resulting activity timeline.
- [ ] Pipeline/stage configuration and deal-board movement are capability-checked, tenant-safe, and reflected in reports/metrics.
- [ ] Unauthorized entity access and cross-tenant operations are denied and mutations are audited.
- [ ] Generated API, UI, OpenAPI client, and browser tests prove core CRM behavior.

### Blocked by

- Issue 1 — generated shell and seed contract.
- Issue 2 — Workspace authorization and audit behavior.

## 5. Deliver operational jobs, schedules, and notifications

### What to build

Generate an operations workflow for jobs and cron: queues, active/retrying/dead-letter work, schedules, run history, failures, notifications, and a meaningful jobs/cron-only sample. Add CMS scheduled-publishing and CRM follow-up automations when the corresponding feature packs are selected.

### Acceptance criteria

- [ ] Jobs/cron-only generated applications expose a runnable sample workflow and operational status instead of an empty page.
- [ ] Operators can inspect queue/schedule/run/failure state through authorized API/UI views.
- [ ] CMS scheduling and CRM follow-up jobs enqueue, execute, retry, and surface outcome notifications when their dependencies are selected.
- [ ] Tests prove claim, retry, idempotency, scheduling, restart recovery, and access control behavior.

### Blocked by

- Issue 1 — generated shell and seed contract.
- Issue 2 — Workspace authorization and audit behavior.
- Issue 3 — CMS scheduled-publication producer.
- Issue 4 — CRM follow-up producer.

## 6. Add CRM CSV exchange and import diagnostics

### What to build

Generate CRM CSV export and validated import workflows for supported CRM records, including row-level validation/error reports, authorization, tenant safety, audit events, and demo data that lets the workflow be evaluated immediately.

### Acceptance criteria

- [ ] An authorized user can export filtered CRM records and import a valid CSV into the active organization.
- [ ] Invalid rows receive actionable errors without partial unsafe writes; successful rows remain idempotent according to documented identity rules.
- [ ] Imports/exports reject unauthorized and cross-tenant access and create audit entries.
- [ ] Public HTTP and browser tests cover success, validation failure, authorization denial, and seeded demo behavior.

### Blocked by

- Issue 4 — complete CRM sales workflow.

## 7. Prove the resolved feature-graph matrix

### What to build

Extend generator scenarios and CI-ready verification so all 18 valid resolved feature graphs are generated into disposable fixtures and proven usable across PostgreSQL/MySQL, frontend/backend-only variants, and applicable seed profiles. Maintain acceptance evidence in the progress ledger.

### Acceptance criteria

- [ ] Every resolved graph has a declared acceptance scenario covering manifest, files, migrations, routes, navigation/documentation, and selected workers/schedules.
- [ ] PostgreSQL and MySQL run migrations up/reset/up and idempotent seed replay for the matrix; generated backend checks pass.
- [ ] Frontend-valid scenarios export OpenAPI, regenerate hey-api, typecheck, test/build, and execute representative browser flows.
- [ ] The ledger records commands, graph/dialect/seed coverage, outcomes, diagnostics, and residual gaps; `full-stack-app` is never changed.

### Blocked by

- Issue 2 — Workspace identity and administration.
- Issue 3 — CMS workflow.
- Issue 5 — operations and automation.
- Issue 6 — CRM CSV exchange.
