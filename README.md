# FastCRUD CLI

The FastCRUD CLI generates strongly-typed clients (Go, C#, Java, TypeScript) for the REST and GraphQL APIs FastCRUD exposes on top of your existing **SQL Server, Oracle, PostgreSQL, MySQL, or MongoDB** databases.

It's the same code-gen pipeline our enterprise customers use to drop a typed client into their internal services without hand-rolling DTOs against a 400-table legacy schema.

## What FastCRUD is

FastCRUD ([fastcrud.dev](https://fastcrud.dev)) is the audit-ready data API for the databases regulated teams already run on. Connect an existing SQL Server, Oracle, PostgreSQL, MySQL, or MongoDB instance, and you get:

- A documented REST + GraphQL API in about a minute.
- Full request-level audit logs (actor, IP, query, latency, outcome) — exportable for SOC 2 and HIPAA evidence.
- Per-key row-level access policies, IP allowlists, KMS-encrypted credentials.
- SSO, data residency, and SLAs on the Business plan.
- BYO-VPC / on-prem on Enterprise.

This CLI consumes the same `/authenticate/crud/:key` flow and the same introspected schema metadata your dashboard uses, so the generated client stays in sync with whatever audit rules and row-level policies you've defined on the access key.

## Features

- **Multi-language code generation** — Go, C#, Java, TypeScript out of the box.
- **Built for legacy schemas** — composite keys, ugly column names, denormalized tables, mixed casing — handled.
- **Scoped access keys** — generated clients honor the row-level filters and table scopes attached to the key.
- **CI-friendly** — single binary, deterministic output, no network calls outside of your FastCRUD instance.

## Installation

```bash
go install github.com/Skyvko6607/fastcrud/cli@latest
```

Or build from source:

```bash
git clone https://github.com/Skyvko6607/fastcrud-cli
cd fastcrud-cli
go build -o fastcrud-cli
```

## Usage

```bash
fastcrud-cli --key <your-access-key-id> --lang <language> [--output <dir>] [--url <base-url>]
```

### Options

| Flag       | Description                                                              | Default                       |
|------------|--------------------------------------------------------------------------|-------------------------------|
| `--key`    | **Required.** Access key ID (UUID) issued from the FastCRUD dashboard.   |                               |
| `--lang`   | **Required.** Target language: `go`, `csharp`, `typescript`, `java`.     |                               |
| `--output` | Output directory for generated code.                                     | `./generated`                 |
| `--url`    | Base URL of the FastCRUD API. Override for self-hosted / BYO-VPC.        | `https://crud.fastcrud.dev`   |

### Example — typed client against an internal SQL Server

```bash
fastcrud-cli \
  --key 123e4567-e89b-12d3-a456-426614174000 \
  --lang csharp \
  --output ./src/Internal.FinanceLedger.Client
```

### Self-hosted / Enterprise

If you're on the Enterprise plan with a BYO-VPC or on-prem deployment, point `--url` at your dedicated API host. The generated client is identical — only the audit trail and credential store change.

```bash
fastcrud-cli \
  --key <key> \
  --lang go \
  --url https://fastcrud.corp.internal
```

## Supported Languages

- **Go** — Native structs and `database/sql`-shaped helpers.
- **C# / .NET** — Strongly-typed records, DI-friendly client.
- **Java** — POJOs and a repository-pattern client.
- **TypeScript** — Interfaces and a fetch-based client suitable for Node, Bun, and the browser.

## Audit & access notes

Generated clients send requests over HTTPS to the FastCRUD API. Every request is logged server-side against the access key's audit trail — including the queries the generated client issues during development. If you don't want this, generate against a key scoped to a test schema, not a production tenant.

## Licensing

Generated code is yours to use, modify, and redistribute. The CLI itself is licensed under the FastCRUD CLI License — see `LICENSE` for details.

## Contact

- Sales / Enterprise: [sales@fastcrud.dev](mailto:sales@fastcrud.dev)
- Support: [support@fastcrud.dev](mailto:support@fastcrud.dev)
- Dashboard: [dashboard.fastcrud.dev](https://dashboard.fastcrud.dev)
