# ConnectRPC Boilerplate

A copy-to-start monorepo: **Go + Echo + ConnectRPC** backend and a **React Router (SPA) + Vite + Tailwind** frontend, wired end-to-end via buf-generated types, all inside a **Nix** devshell.

## Stack

- Nix flake devshell + direnv (`go`, `bun`, `buf`, `just`, `air`)
- buf v2 with remote codegen plugins (Go + TS)
- Backend: Go 1.25, Echo, connect-go, plain-struct DI container (`internal/app`)
- Frontend: React Router v7 (SPA), Vite, Tailwind v4, connect-es-web

## Quick start

```bash
cp -r boilerplate my-project && cd my-project
direnv allow            # or: nix develop
just install            # install frontend deps
just dev                # run backend (:8000) + frontend together
```

Open the frontend dev URL, type a name, click **Greet** — it round-trips through the backend.

## Common tasks

```bash
just            # list tasks
just gen        # regenerate Go + TS from proto/
just dev        # backend + frontend
just build      # build both
```

## Layout

```
proto/        # app.v1 service definitions + buf config
backend/      # cmd/server entrypoint + internal/{config,app,server,service}
frontend/     # React Router app; connect client in app/lib/client.ts
```

Generated code (`backend/gen`, `frontend/gen`) is committed so a fresh copy builds immediately. Run `just gen` after editing `.proto` files.

## Rename (optional)

The Go module is `backend` and the proto package is `app.v1` — both generic, so a copy compiles unchanged. To rebrand:

- Proto package: edit `package` and `go_package` in `proto/app/v1/app.proto`, adjust `out` paths if you move them, then `just gen`.
- Go module: `cd backend && go mod edit -module <name>` and update imports.

## Add an RPC

1. Add the `rpc` (and messages) to `proto/app/v1/app.proto`.
2. `just gen`.
3. Implement the method on a service in `backend/internal/service/` (see `greeter.go`).
4. If it's a new service, mount it with one line in `backend/internal/server/server.go`:
   `service.NewMyService(a).Mount(e)`.
5. Call it from the frontend via `client` in `app/lib/client.ts`.

## Add a feature via the DI container

The container is `backend/internal/app/app.go`. To add, say, a database:

1. Add a field: `DB *gorm.DB` to `App`.
2. Construct it in `New(cfg)` from config and assign it.
3. Services already receive `*app.App`, so they can use `a.DB`.

## Deployment

Both apps have a `Dockerfile`. Build with `docker build ./backend` and `docker build ./frontend`. Set the frontend's `VITE_API_URL` to the backend's public URL.
