set shell := ["bash", "-uc"]

# List available tasks
default:
    @just --list

# Regenerate Go + TS code from proto
gen:
    cd proto && buf generate

# Install frontend dependencies
install:
    cd frontend && bun install

# Run the Go backend with hot reload
dev-backend:
    cd backend && air

# Run the frontend dev server
dev-frontend:
    cd frontend && bun dev

# Run backend and frontend together
dev:
    just dev-backend & just dev-frontend

# Build both backend and frontend
build:
    cd backend && go build ./...
    cd frontend && bun run build
