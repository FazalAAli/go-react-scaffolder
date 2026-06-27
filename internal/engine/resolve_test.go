package engine

import "testing"

func cat() []Feature {
	return []Feature{
		{Name: "infra", Always: true, Dir: "f/infra",
			Regions: []RegionContribution{{Target: "flake.nix", Region: "pkgs", Content: "just"}}},
		{Name: "backend", Always: true, Dir: "f/backend",
			Regions: []RegionContribution{{Target: "flake.nix", Region: "pkgs", Content: "go"}}},
		{Name: "air", Always: false, Dir: "f/air",
			Regions:   []RegionContribution{{Target: "flake.nix", Region: "pkgs", Content: "air"}},
			PostSteps: []string{"tidy"}},
	}
}

func TestResolveAlwaysOnly(t *testing.T) {
	p, err := Resolve(cat(), nil)
	if err != nil {
		t.Fatal(err)
	}
	// always features sorted alpha: backend, infra
	if got := p.Features; len(got) != 2 || got[0] != "backend" || got[1] != "infra" {
		t.Fatalf("Features = %v", got)
	}
	key := RegionKey{Target: "flake.nix", Region: "pkgs"}
	if got := p.Regions[key]; len(got) != 2 || got[0] != "go" || got[1] != "just" {
		t.Fatalf("contribs = %v", got)
	}
	if len(p.FileRoots) != 2 || p.FileRoots[0] != "f/backend/files" {
		t.Fatalf("FileRoots = %v", p.FileRoots)
	}
}

func TestResolveWithOptional(t *testing.T) {
	p, err := Resolve(cat(), []string{"air"})
	if err != nil {
		t.Fatal(err)
	}
	// always (backend, infra) then optional (air)
	if got := p.Features; len(got) != 3 || got[2] != "air" {
		t.Fatalf("Features = %v", got)
	}
	key := RegionKey{Target: "flake.nix", Region: "pkgs"}
	if got := p.Regions[key]; len(got) != 3 || got[2] != "air" {
		t.Fatalf("contribs = %v", got)
	}
	if len(p.PostSteps) != 1 || p.PostSteps[0] != "tidy" {
		t.Fatalf("PostSteps = %v", p.PostSteps)
	}
}

func TestResolveUnknownFeature(t *testing.T) {
	if _, err := Resolve(cat(), []string{"nope"}); err == nil {
		t.Fatal("expected error for unknown feature")
	}
}

func TestResolveRejectsAlwaysSelection(t *testing.T) {
	if _, err := Resolve(cat(), []string{"infra"}); err == nil {
		t.Fatal("expected error selecting an always feature")
	}
}

func TestResolveMultiOptionalOrder(t *testing.T) {
	c := []Feature{
		{Name: "core", Always: true, Dir: "f/core"},
		{Name: "zeta", Always: false, Dir: "f/zeta"},
		{Name: "alpha", Always: false, Dir: "f/alpha"},
	}
	p, err := Resolve(c, []string{"zeta", "alpha"})
	if err != nil {
		t.Fatal(err)
	}
	// always (core) then optional alpha-sorted (alpha, zeta)
	if got := p.Features; len(got) != 3 || got[0] != "core" || got[1] != "alpha" || got[2] != "zeta" {
		t.Fatalf("Features = %v, want [core alpha zeta]", got)
	}
}

func TestResolveDuplicateSelection(t *testing.T) {
	p, err := Resolve(cat(), []string{"air", "air"})
	if err != nil {
		t.Fatalf("duplicate selection should not error: %v", err)
	}
	// air appears exactly once (always backend, infra + optional air)
	count := 0
	for _, f := range p.Features {
		if f == "air" {
			count++
		}
	}
	if count != 1 {
		t.Fatalf("air appears %d times, want 1 (dedup); Features = %v", count, p.Features)
	}
}
