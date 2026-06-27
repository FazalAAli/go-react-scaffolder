package engine

import "testing"

func TestRenderRegionsFillsUnion(t *testing.T) {
	in := "before\n# scaffold:region:pkgs:start\n# scaffold:region:pkgs:end\nafter\n"
	out, err := RenderRegions(in, map[string][]string{"pkgs": {"go", "bun\nair"}})
	if err != nil {
		t.Fatal(err)
	}
	want := "before\n# scaffold:region:pkgs:start\ngo\nbun\nair\n# scaffold:region:pkgs:end\nafter\n"
	if out != want {
		t.Errorf("got:\n%q\nwant:\n%q", out, want)
	}
}

func TestRenderRegionsReplacesExisting(t *testing.T) {
	// re-rendering an already-filled region replaces, not appends (idempotent)
	in := "# scaffold:region:pkgs:start\nstale\n# scaffold:region:pkgs:end\n"
	out, err := RenderRegions(in, map[string][]string{"pkgs": {"fresh"}})
	if err != nil {
		t.Fatal(err)
	}
	want := "# scaffold:region:pkgs:start\nfresh\n# scaffold:region:pkgs:end\n"
	if out != want {
		t.Errorf("got %q want %q", out, want)
	}
}

func TestRenderRegionsEmptyContribution(t *testing.T) {
	in := "# scaffold:region:pkgs:start\n# scaffold:region:pkgs:end\n"
	out, err := RenderRegions(in, map[string][]string{})
	if err != nil {
		t.Fatal(err)
	}
	if out != in {
		t.Errorf("got %q want unchanged", out)
	}
}

func TestRenderRegionsMultiple(t *testing.T) {
	in := "# scaffold:region:pkgs:start\n# scaffold:region:pkgs:end\nmid\n# scaffold:region:recipes:start\n# scaffold:region:recipes:end\n"
	out, err := RenderRegions(in, map[string][]string{
		"pkgs":    {"go"},
		"recipes": {"build"},
	})
	if err != nil {
		t.Fatal(err)
	}
	want := "# scaffold:region:pkgs:start\ngo\n# scaffold:region:pkgs:end\nmid\n# scaffold:region:recipes:start\nbuild\n# scaffold:region:recipes:end\n"
	if out != want {
		t.Errorf("got:\n%q\nwant:\n%q", out, want)
	}
}

func TestRenderRegionsUnreferencedContribution(t *testing.T) {
	in := "# scaffold:region:pkgs:start\n# scaffold:region:pkgs:end\n"
	out, err := RenderRegions(in, map[string][]string{"other": {"x"}})
	if err != nil {
		t.Fatal(err)
	}
	if out != in {
		t.Errorf("got %q want unchanged %q", out, in)
	}
}

func TestRenderRegionsMissingEnd(t *testing.T) {
	in := "# scaffold:region:pkgs:start\n"
	if _, err := RenderRegions(in, nil); err == nil {
		t.Fatal("expected error for missing end marker")
	}
}
