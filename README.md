# scaffold

A project generator for Go + ConnectRPC + React Router apps. `scaffold new`
stamps out a fresh project from an always-on core plus optional features, and
`scaffold add` applies further features to a project after the fact.

## Usage

```bash
go run ./cmd/scaffold new ../my-app                 # interactive feature select
go run ./cmd/scaffold new ../my-app --with air,dockerfile   # non-interactive

# add features to an already-generated project (run inside it, or pass --dir)
go run ./cmd/scaffold add posthog
go run ./cmd/scaffold add sentry --dir ../my-app
```

`new` records the chosen features in a `scaffold.toml` manifest at the project
root. `add` reads that manifest to know what is already applied: it copies only
the new feature's `files/` payload (leaving your edits to existing files
intact), re-renders every managed region from the union of all active features,
re-runs the core finalizers (`go mod tidy`, `bun install`) plus the new
feature's own post-steps, and updates the manifest. Re-adding an applied feature
is a no-op.

## Catalog

- **Always:** `backend` (Go/Echo/connect-go + proto), `frontend` (React Router SPA), `infra` (Nix/just/direnv).
- **Optional:** `air` (Go hot-reload), `dockerfile` (production Dockerfiles), `postgres` (PostgreSQL via pgx + sqlc + goose), `sqlite` (SQLite via modernc.org/sqlite + sqlc + goose), `posthog` (React and Go).

`postgres` and `sqlite` are mutually exclusive — selecting both is rejected.

Each feature is a directory under `features/` with a `feature.toml` and a
`files/` payload copied verbatim into the project. Shared files (`flake.nix`,
`justfile`, `.gitignore`) carry named managed regions filled from the active
features. A feature may declare relationships in its `feature.toml`:

```toml
conflicts = ["sqlite"]    # cannot be active alongside these features
requires  = ["postgres"]  # these features must also be active
```

`Resolve()` enforces both (for `new` and `add`) and fails with a clear error
when a selection is contradictory.