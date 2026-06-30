package engine

import (
	"path/filepath"
	"reflect"
	"testing"
)

func TestReadManifest(t *testing.T) {
	dir := t.TempDir()
	mustWrite(t, filepath.Join(dir, "scaffold.toml"), "features = [\"backend\", \"posthog\"]\n")
	got, err := ReadManifest(dir)
	if err != nil {
		t.Fatalf("ReadManifest: %v", err)
	}
	if want := []string{"backend", "posthog"}; !reflect.DeepEqual(got, want) {
		t.Fatalf("ReadManifest = %v, want %v", got, want)
	}
}

func TestReadManifestMissing(t *testing.T) {
	if _, err := ReadManifest(t.TempDir()); err == nil {
		t.Fatal("expected error reading a missing scaffold.toml")
	}
}

func TestReadManifestRoundTripsWrite(t *testing.T) {
	// the manifest ReadManifest parses is exactly what Write produces.
	dir := t.TempDir()
	if err := writeScaffoldToml([]string{"infra", "air", "backend"}, dir); err != nil {
		t.Fatalf("writeScaffoldToml: %v", err)
	}
	got, err := ReadManifest(dir)
	if err != nil {
		t.Fatalf("ReadManifest: %v", err)
	}
	// Write sorts the recorded features.
	if want := []string{"air", "backend", "infra"}; !reflect.DeepEqual(got, want) {
		t.Fatalf("round-trip = %v, want %v", got, want)
	}
}
