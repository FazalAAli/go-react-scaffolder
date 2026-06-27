# scaffold

A project generator for Go + ConnectRPC + React Router apps. `scaffold new`
stamps out a fresh project from an always-on core plus optional features.

## Usage

```bash
go run ./cmd/scaffold new ../my-app                 # interactive feature select
go run ./cmd/scaffold new ../my-app --with air,dockerfile   # non-interactive
```

## Catalog

- **Always:** `backend` (Go/Echo/connect-go + proto), `frontend` (React Router SPA), `infra` (Nix/just/direnv).
- **Optional:** `air` (Go hot-reload), `dockerfile` (production Dockerfiles).

Each feature is a directory under `features/` with a `feature.toml` and a
`files/` payload copied verbatim into the project. Shared files (`flake.nix`,
`justfile`, `.gitignore`) carry named managed regions filled from the active
features. See `docs/superpowers/specs/` for the design.
