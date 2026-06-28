package engine

import (
	"path/filepath"
	"strings"
	"testing"
)

// catalogDir is the real feature catalog at the repo root, relative to this
// package directory (internal/engine).
const catalogDir = "../../features"

// writeRealCatalog loads the real catalog, resolves it with the given optional
// features, writes the project into a temp dir, and returns that dir.
func writeRealCatalog(t *testing.T, optional ...string) string {
	t.Helper()
	cat, err := LoadCatalog(catalogDir)
	if err != nil {
		t.Fatalf("LoadCatalog: %v", err)
	}
	plan, err := Resolve(cat, optional)
	if err != nil {
		t.Fatalf("Resolve: %v", err)
	}
	dst := t.TempDir()
	if err := Write(plan, dst); err != nil {
		t.Fatalf("Write: %v", err)
	}
	return dst
}

func TestCatalogEnvFilesManaged(t *testing.T) {
	dst := writeRealCatalog(t)

	backendEnv := read(t, filepath.Join(dst, "backend", ".env"))
	if !strings.Contains(backendEnv, "scaffold:region:env:start") {
		t.Errorf("backend/.env is not managed (no env region marker):\n%s", backendEnv)
	}
	for _, want := range []string{"PORT=8000", "FRONTEND_BASE_URL=", "ENV=development"} {
		if !strings.Contains(backendEnv, want) {
			t.Errorf("backend/.env missing %q:\n%s", want, backendEnv)
		}
	}
	if strings.Contains(backendEnv, "DATABASE_URL") {
		t.Errorf("backend/.env should not contain DATABASE_URL without a db feature:\n%s", backendEnv)
	}

	frontendEnv := read(t, filepath.Join(dst, "frontend", ".env"))
	if !strings.Contains(frontendEnv, "scaffold:region:env:start") {
		t.Errorf("frontend/.env is not managed (no env region marker):\n%s", frontendEnv)
	}
	if !strings.Contains(frontendEnv, "VITE_API_URL=") {
		t.Errorf("frontend/.env missing VITE_API_URL:\n%s", frontendEnv)
	}
}

func TestCatalogAppGoSeams(t *testing.T) {
	dst := writeRealCatalog(t)
	appGo := read(t, filepath.Join(dst, "backend", "internal", "app", "app.go"))

	for _, want := range []string{
		"scaffold:region:app-imports:start",
		"scaffold:region:app-fields:start",
		"scaffold:region:app-init:start",
		"a := &App{Config: cfg}",
		"return a, nil",
	} {
		if !strings.Contains(appGo, want) {
			t.Errorf("app.go missing %q:\n%s", want, appGo)
		}
	}
	if strings.Contains(appGo, "db.Pool") {
		t.Errorf("app.go should not reference db.Pool without a db feature:\n%s", appGo)
	}
}

func TestCatalogPostgresFeature(t *testing.T) {
	dst := writeRealCatalog(t, "postgres")

	backendEnv := read(t, filepath.Join(dst, "backend", ".env"))
	if !strings.Contains(backendEnv, "DATABASE_URL=postgres://") {
		t.Errorf("backend/.env missing postgres DATABASE_URL:\n%s", backendEnv)
	}

	appGo := read(t, filepath.Join(dst, "backend", "internal", "app", "app.go"))
	for _, want := range []string{
		"DB *db.Pool",
		`"backend/internal/db"`,
		"db.New(context.Background())",
	} {
		if !strings.Contains(appGo, want) {
			t.Errorf("app.go missing %q:\n%s", want, appGo)
		}
	}

	dbGo := read(t, filepath.Join(dst, "backend", "internal", "db", "db.go"))
	if !strings.Contains(dbGo, "pgxpool") {
		t.Errorf("db.go should use pgxpool:\n%s", dbGo)
	}

	// sqlc + migrations scaffolding present (read() fails the test if absent)
	read(t, filepath.Join(dst, "backend", "sqlc.yaml"))
	read(t, filepath.Join(dst, "backend", "migrations", ".gitkeep"))
	read(t, filepath.Join(dst, "backend", "queries", ".gitkeep"))

	justfile := read(t, filepath.Join(dst, "justfile"))
	for _, want := range []string{"db-gen:", "db-migrate", "sqlc generate", "goose"} {
		if !strings.Contains(justfile, want) {
			t.Errorf("justfile missing %q:\n%s", want, justfile)
		}
	}

	flake := read(t, filepath.Join(dst, "flake.nix"))
	for _, want := range []string{"pkgs.sqlc", "pkgs.goose"} {
		if !strings.Contains(flake, want) {
			t.Errorf("flake.nix missing %q:\n%s", want, flake)
		}
	}
}

func TestCatalogSqliteFeature(t *testing.T) {
	dst := writeRealCatalog(t, "sqlite")

	backendEnv := read(t, filepath.Join(dst, "backend", ".env"))
	if !strings.Contains(backendEnv, "DATABASE_URL=app.db") {
		t.Errorf("backend/.env missing sqlite DATABASE_URL:\n%s", backendEnv)
	}

	appGo := read(t, filepath.Join(dst, "backend", "internal", "app", "app.go"))
	for _, want := range []string{
		"DB *db.Pool",
		`"backend/internal/db"`,
		"db.New(context.Background())",
	} {
		if !strings.Contains(appGo, want) {
			t.Errorf("app.go missing %q:\n%s", want, appGo)
		}
	}

	dbGo := read(t, filepath.Join(dst, "backend", "internal", "db", "db.go"))
	if !strings.Contains(dbGo, "modernc.org/sqlite") {
		t.Errorf("db.go should use modernc.org/sqlite:\n%s", dbGo)
	}

	read(t, filepath.Join(dst, "backend", "sqlc.yaml"))
	read(t, filepath.Join(dst, "backend", "migrations", ".gitkeep"))
	read(t, filepath.Join(dst, "backend", "queries", ".gitkeep"))
	read(t, filepath.Join(dst, "backend", "go.mod"))

	justfile := read(t, filepath.Join(dst, "justfile"))
	for _, want := range []string{"db-gen:", "db-migrate", "sqlc generate", "goose"} {
		if !strings.Contains(justfile, want) {
			t.Errorf("justfile missing %q:\n%s", want, justfile)
		}
	}

	flake := read(t, filepath.Join(dst, "flake.nix"))
	for _, want := range []string{"pkgs.sqlc", "pkgs.goose"} {
		if !strings.Contains(flake, want) {
			t.Errorf("flake.nix missing %q:\n%s", want, flake)
		}
	}

	gitignore := read(t, filepath.Join(dst, ".gitignore"))
	if !strings.Contains(gitignore, "/backend/*.db") {
		t.Errorf(".gitignore missing sqlite db ignore:\n%s", gitignore)
	}
}
