# Feature Catalog

The generator composes a project from a small set of product-facing packs. Select packs with repeated `--feature` flags; required packs are added automatically.

| Feature | What it provides | Requires |
| --- | --- | --- |
| `auth` | Accounts, sign-up/sign-in, cookie sessions, refresh rotation, and organization bootstrap | — |
| `rbac` | Tenant-scoped roles, permissions, and authorization | `auth` |
| `workspace` | Member provisioning, role assignment, and workspace administration | `auth`, `rbac` |
| `content` | Posts, categories, tags, and publishing workflows | `auth`, `rbac` |
| `customers` | Companies and contacts for customer operations | `auth`, `rbac` |
| `sales` | Pipelines, stages, deals, and activities | `customers` |
| `jobs` | Durable background-job execution | — |
| `cron` | Recurring scheduling through the jobs queue | `jobs` |
| `audit` | Tenant-authorized audit domain and viewer foundation | `auth`, `rbac` |
| `files` | Tenant-authorized S3-compatible file domain and upload foundation | `auth`, `rbac` |
| `email` | Provider-neutral transactional message domain and safe log adapter | `jobs` |

`content` is the clear product name for the publishing pack previously called CMS. `customers` and `sales` replace the former broad CRM label: companies/contacts are useful without a deal pipeline, while sales always has customer records available.

`workspace` manages application users, memberships, and roles. CRM contacts are external business records, not login accounts. A later customer-portal pack will explicitly link selected contacts to accounts when customer login is needed.

The generator rejects the retired `cms` and `crm` feature names. Replace them with `content`, and either `customers` or `sales` as appropriate.

## Initial owner

For a secure first production deployment, set all of `BOOTSTRAP_OWNER_EMAIL`, `BOOTSTRAP_OWNER_PASSWORD`, and `BOOTSTRAP_ORGANIZATION_NAME`. On an empty database the application creates one organization owner with `platform.admin` exactly once. Supplying only part of this configuration fails startup.

The `minimal` and `demo` seed profiles also create the local-only account `owner@example.test` with password `DemoPass123!`. Never use those public development credentials outside local development.

## Examples

```sh
# Editorial workspace and customer records.
go run . create --name acme --module example.com/acme --feature content --feature customers

# Full customer-operations workspace; sales resolves customers, auth, and RBAC.
go run . create --name acme --module example.com/acme --feature workspace --feature sales --feature jobs --feature cron
```
