package engine

import (
	"path/filepath"
	"testing"
)

// AddFeatures applies newly-selected features to an existing project without
// clobbering the user's edits to already-applied files: it copies only the
// payloads named in `add`, re-renders managed regions from the full plan, and
// rewrites scaffold.toml.
func TestAddFeaturesCopiesOnlyNewPayloadAndReRendersRegions(t *testing.T) {
	src := t.TempDir()
	// "base" owns the shared flake.nix and ships an app file the user may edit.
	mustWrite(t, filepath.Join(src, "base", "files", "flake.nix"),
		"packages = [\n# scaffold:region:pkgs:start\ngo\n# scaffold:region:pkgs:end\n]\n")
	mustWrite(t, filepath.Join(src, "base", "files", "app.go"), "package app // pristine\n")
	// "extra" is the newly added optional feature: a new file + a region contribution.
	mustWrite(t, filepath.Join(src, "extra", "files", "extra.txt"), "extra payload\n")

	// Simulate an existing project that already has base applied, with a
	// user-edited app.go that must survive the add.
	dst := t.TempDir()
	mustWrite(t, filepath.Join(dst, "flake.nix"),
		"packages = [\n# scaffold:region:pkgs:start\ngo\n# scaffold:region:pkgs:end\n]\n")
	mustWrite(t, filepath.Join(dst, "app.go"), "package app // USER EDIT\n")
	mustWrite(t, filepath.Join(dst, "scaffold.toml"), "features = [\"base\"]\n")

	plan := Plan{
		Features:  []string{"base", "extra"},
		FileRoots: []string{filepath.Join(src, "base", "files"), filepath.Join(src, "extra", "files")},
		Regions: map[RegionKey][]string{
			{Target: "flake.nix", Region: "pkgs"}: {"go", "air"},
		},
	}

	if err := AddFeatures(plan, dst, map[string]bool{"extra": true}); err != nil {
		t.Fatalf("AddFeatures: %v", err)
	}

	// new feature's payload is copied in
	if got := read(t, filepath.Join(dst, "extra.txt")); got != "extra payload\n" {
		t.Errorf("extra.txt = %q, want it copied", got)
	}
	// the user's edit to an already-applied file is preserved (base payload NOT re-copied)
	if got := read(t, filepath.Join(dst, "app.go")); got != "package app // USER EDIT\n" {
		t.Errorf("app.go = %q, want the user edit preserved", got)
	}
	// the shared region is re-rendered from the full union
	wantFlake := "packages = [\n# scaffold:region:pkgs:start\ngo\nair\n# scaffold:region:pkgs:end\n]\n"
	if got := read(t, filepath.Join(dst, "flake.nix")); got != wantFlake {
		t.Errorf("flake.nix:\n%q\nwant:\n%q", got, wantFlake)
	}
	// manifest updated to the full feature set
	if got := read(t, filepath.Join(dst, "scaffold.toml")); got != "features = [\"base\", \"extra\"]\n" {
		t.Errorf("scaffold.toml = %q", got)
	}
}
