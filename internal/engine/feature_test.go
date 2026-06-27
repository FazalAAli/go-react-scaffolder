package engine

import (
	"path/filepath"
	"testing"
)

func TestLoadFeature(t *testing.T) {
	dir := filepath.Join("testdata", "features", "alpha")
	f, err := LoadFeature(dir)
	if err != nil {
		t.Fatalf("LoadFeature: %v", err)
	}
	if f.Name != "alpha" {
		t.Errorf("Name = %q, want alpha", f.Name)
	}
	if !f.Always {
		t.Errorf("Always = false, want true")
	}
	if f.Dir != dir {
		t.Errorf("Dir = %q, want %q", f.Dir, dir)
	}
	if len(f.Regions) != 1 || f.Regions[0].Target != "flake.nix" || f.Regions[0].Region != "devshell-packages" || f.Regions[0].Content != "pkgs.alpha" {
		t.Errorf("Regions = %+v", f.Regions)
	}
	if len(f.PostSteps) != 1 || f.PostSteps[0] != "echo hi" {
		t.Errorf("PostSteps = %+v", f.PostSteps)
	}
}

func TestLoadCatalog(t *testing.T) {
	cat, err := LoadCatalog(filepath.Join("testdata", "features"))
	if err != nil {
		t.Fatalf("LoadCatalog: %v", err)
	}
	if len(cat) != 1 || cat[0].Name != "alpha" {
		t.Errorf("catalog = %+v", cat)
	}
}
