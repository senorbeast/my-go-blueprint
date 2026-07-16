# Demo-Ready Generated Applications PRD

## Problem Statement

The generator can select auth, RBAC, jobs, cron, CMS, and CRM, but generated applications do not consistently provide a complete, sign-in-ready workflow. Some feature combinations render a sparse landing page, demo credentials are not reliably usable, and the CMS, CRM, workspace administration, and operations surfaces do not provide the workflows that a prospective customer expects to evaluate. The generator must produce a useful, tenant-safe application for every valid resolved feature graph without hand-editing generated projects.

## Solution

Enhance generator-owned templates to produce a capability-driven, shadcn-based application shell and a deep Workspace module whenever Auth is selected. The Workspace module owns sessions, active organization context, tenant administration, platform administration, role/permission management, invitations, audit history, and session revocation. CMS, CRM, jobs, and cron add complete, permission-aware workflows to that shell. Deterministic minimal and demo seed profiles create documented credentials and representative records so a newly generated local application can be evaluated immediately.

The generator will verify all 18 valid resolved feature graphs across PostgreSQL and MySQL through public HTTP, worker, and browser behavior. Planning artifacts are local Markdown under `docs/plans`; `full-stack-app` remains an untouched generated fixture.

## User Stories

1. As a generator user, I want every valid feature selection to produce a navigable and usable application, so that I can evaluate the template without writing missing glue code.
2. As a demo evaluator, I want a documented local email and password, so that I can sign in immediately.
3. As a demo evaluator, I want representative content, CRM, operations, notification, and audit data, so that dashboards and workflows are meaningful on first launch.
4. As an application operator, I want seed profiles to be deterministic and replay-safe, so that local resets and CI remain reliable.
5. As a first-time deployer using no seeds, I want clear initial-user setup instructions, so that the application is still bootable and secure.
6. As a signed-in user, I want my actor, active organization, organizations, capabilities, and platform-admin state exposed in my session context, so that the UI reflects my authorized workspace.
7. As a user with membership in several organizations, I want to switch active organizations safely, so that I can work in the correct tenant.
8. As a user, I want server logout to revoke my session, so that signing out ends access on the server as well as in the browser.
9. As a platform super-admin, I want to create and administer organizations and users, so that I can operate a multi-tenant installation.
10. As a tenant owner or administrator, I want to manage my organization’s members, invitations, roles, and feature access, so that I can delegate work safely.
11. As a tenant administrator, I want protection against removing the final owner or granting privileges I do not hold, so that tenant administration cannot lock out or escalate users.
12. As a security reviewer, I want immutable audit history for administrative and feature lifecycle actions, so that I can investigate changes.
13. As a user, I want navigation and empty states to reflect selected features and my permissions, so that I understand what I can do without seeing forbidden controls.
14. As a CMS editor, I want to create, edit, filter, preview, schedule, publish, archive, restore, and revise content, so that I can manage an editorial lifecycle.
15. As a CMS administrator, I want categories, tags, and media-library workflows, so that content is organized and reusable.
16. As a CMS publisher, I want scheduled publishing to execute through the generated worker and scheduler, so that future content goes live reliably.
17. As a CRM user, I want companies, contacts, pipelines, stages, deals, activities, notes, and assignments, so that I can manage customer work end to end.
18. As a sales user, I want a deal board and pipeline metrics, so that I can prioritize and understand opportunities.
19. As a CRM administrator, I want CSV import/export with clear validation failures, so that I can exchange data safely.
20. As a CRM user, I want follow-up reminders and task automation, so that opportunities do not go unattended.
21. As an operations user, I want job, schedule, run, retry, dead-letter, and notification visibility, so that I can diagnose automation.
22. As a backend-only generator user, I want equivalent documented authenticated HTTP and OpenAPI workflows, so that a frontend is not required for operability.
23. As a maintainer, I want generated frontend clients to originate from OpenAPI, so that browser/API contracts do not drift.
24. As a maintainer, I want behavior-based proof across both database dialects and all valid feature graphs, so that template rendering alone cannot hide a broken combination.

## Implementation Decisions

- Generate Workspace with Auth. Its deep module boundary owns session context, session revocation and rotation, organization switching, users, organizations, memberships, invitations, roles, permission grants, and audit emission. Feature modules register a capability catalog at bootstrap; Workspace owns grants and administrative access.
- Keep platform authority separate from tenant authority. Platform super-admin is a distinct user claim and can administer platform organizations/users. Tenant roles grant capabilities only within a membership’s organization; they cannot confer platform authority.
- Require membership before organization switching, rotate or reissue the server session on a valid switch, reject cross-tenant enumeration, and protect final-owner removal and privilege escalation.
- Generate authenticated contracts for session/logout/switching; workspace and platform administration; CMS content/taxonomy/media/revisions; CRM records/import-export/reporting; and job/schedule observability. Regenerate hey-api from OpenAPI; do not maintain duplicate frontend API shapes.
- Use shadcn UI as the generated React baseline. Generate a responsive authenticated shell, sidebar/navigation, breadcrumbs, forms, tables, dialogs/sheets, tabs, alerts, skeletons, badges, toasts, charts, and permission-aware states. Backend-only projects instead receive documented OpenAPI/HTTP paths.
- Make `none` bootable with first-user instructions. Make `minimal` create a deterministic credential, organization, owner membership, relevant capability grants, and minimal feature records. Make `demo` create an environment-gated documented credential plus representative selected-feature, automation, notification, and audit records. All seeds are idempotent.
- CMS includes taxonomy, search/filter/pagination, edit/archive/restore, preview, revisions/restore, media selection/upload, and scheduled publication when jobs/cron are available. CMS lifecycle actions require explicit capabilities and create audit events.
- CRM includes companies, contacts, pipelines/stages, deals, activities/tasks, notes/associations, ownership/assignment, searchable lists, a moveable deal board, reporting, CSV import/export validation reports, and follow-up automation when jobs/cron are available. CRM mutations are tenant-scoped, capability-checked, and audited.
- Jobs and cron have an operational UI/API for queued/running/retry/dead-letter jobs, schedules, run history, failures, and notifications. Jobs/cron-only selections generate a representative runnable operational workflow.
- Preserve generator ownership: modify only generator source/templates and generator documentation, add new disposable generated fixtures for verification, and never modify `full-stack-app`.

## Testing Decisions

- Test observable public behavior first, at the highest available seam: generated HTTP contracts, composed application behavior, worker/scheduler behavior, and browser flows. Avoid coupling to internal helpers.
- Every slice follows feature-slice verification: define allow, deny, and error behavior; add the public-seam test; implement the smallest end-to-end change; verify generator, fresh fixtures, persistence, contracts, frontend, and browser behavior where applicable.
- Verify all 18 resolved feature graphs on PostgreSQL and MySQL. Exercise migrations up/reset/up and seed replay; validate manifests, routes, navigation, docs, workers, and schedules for each graph.
- Run generated backend tests, OpenAPI export and hey-api regeneration, frontend typecheck/unit tests/production build, and Playwright flows for frontend-valid graphs. Pull requests use focused and alternate-dialect coverage; nightly/manual verification runs the complete matrix.
- Cover documented sign-in, logout revocation, organization switching, tenant denial, invite expiry/acceptance, role-escalation denial, final-owner protection, super-admin auditing, permission-aware navigation, CMS workflows, CRM workflows, queue retry/idempotency, and scheduled execution.
- Record exact commands, configurations, outcomes, failures, and unresolved gaps in the implementation ledger after each completed slice.

## Out of Scope

- Billing, subscriptions, payment processing, SSO, SCIM, custom-object/schema builders, email delivery-provider configuration, and enterprise compliance certification.
- Editing or repairing `full-stack-app` directly.
- Publishing a PRD or issues to an external tracker.

## Further Notes

- The local-development demo credential is intentionally visible in generated README documentation and must be environment-gated; production deployments must replace or disable it.
- The generator’s six selectable features resolve to 18 valid feature graphs. Dependency resolution is necessary but is not evidence that a graph is usable.
