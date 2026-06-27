package engine

import (
	"os"
	"path/filepath"
	"testing"
)

func TestWriteCopiesRendersAndRecords(t *testing.T) {
	// Build a two-feature source layout in a temp dir.
	src := t.TempDir()
	// feature "infra": owns the shared flake.nix skeleton
	mustWrite(t, filepath.Join(src, "infra", "files", "flake.nix"),
		"packages = [\n# scaffold:region:pkgs:start\n# scaffold:region:pkgs:end\n]\n")
	// feature "backend": payload file + contribution
	mustWrite(t, filepath.Join(src, "backend", "files", "backend", "main.go"), "package main\n")

	plan := Plan{
		Features:  []string{"backend", "infra"},
		FileRoots: []string{filepath.Join(src, "backend", "files"), filepath.Join(src, "infra", "files")},
		Regions: map[RegionKey][]string{
			{Target: "flake.nix", Region: "pkgs"}: {"go", "just"},
		},
	}

	dst := t.TempDir()
	if err := Write(plan, dst); err != nil {
		t.Fatalf("Write: %v", err)
	}

	if got := read(t, filepath.Join(dst, "backend", "main.go")); got != "package main\n" {
		t.Errorf("payload not copied: %q", got)
	}
	wantFlake := "packages = [\n# scaffold:region:pkgs:start\ngo\njust\n# scaffold:region:pkgs:end\n]\n"
	if got := read(t, filepath.Join(dst, "flake.nix")); got != wantFlake {
		t.Errorf("flake render:\n%q\nwant:\n%q", got, wantFlake)
	}
	if got := read(t, filepath.Join(dst, "scaffold.toml")); got != "features = [\"backend\", \"infra\"]\n" {
		t.Errorf("scaffold.toml = %q", got)
	}
}

func TestWriteLaterFeatureOverwritesEarlier(t *testing.T) {
	src := t.TempDir()
	mustWrite(t, filepath.Join(src, "a", "files", "shared.txt"), "from-a\n")
	mustWrite(t, filepath.Join(src, "b", "files", "shared.txt"), "from-b\n")
	plan := Plan{
		Features:  []string{"a", "b"},
		FileRoots: []string{filepath.Join(src, "a", "files"), filepath.Join(src, "b", "files")},
		Regions:   map[RegionKey][]string{},
	}
	dst := t.TempDir()
	if err := Write(plan, dst); err != nil {
		t.Fatalf("Write: %v", err)
	}
	if got := read(t, filepath.Join(dst, "shared.txt")); got != "from-b\n" {
		t.Errorf("shared.txt = %q, want from-b (later feature wins)", got)
	}
}

func TestWriteSkipsMissingFileRoot(t *testing.T) {
	src := t.TempDir()
	mustWrite(t, filepath.Join(src, "a", "files", "x.txt"), "x\n")
	plan := Plan{
		Features:  []string{"a", "ghost"},
		FileRoots: []string{filepath.Join(src, "a", "files"), filepath.Join(src, "ghost", "files")}, // second does not exist
		Regions:   map[RegionKey][]string{},
	}
	dst := t.TempDir()
	if err := Write(plan, dst); err != nil {
		t.Fatalf("Write should skip a missing FileRoot, got: %v", err)
	}
	if got := read(t, filepath.Join(dst, "x.txt")); got != "x\n" {
		t.Errorf("x.txt = %q, want x", got)
	}
}

func mustWrite(t *testing.T, path, content string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}

func read(t *testing.T, path string) string {
	t.Helper()
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return string(b)
}
