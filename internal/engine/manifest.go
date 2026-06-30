package engine

import (
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

// manifestName is the file written into a generated project recording which
// features it was stamped from, so `scaffold add` knows what is already applied.
const manifestName = "scaffold.toml"

type manifest struct {
	Features []string `toml:"features"`
}

// ReadManifest returns the feature names recorded in <projectRoot>/scaffold.toml.
func ReadManifest(projectRoot string) ([]string, error) {
	data, err := os.ReadFile(filepath.Join(projectRoot, manifestName))
	if err != nil {
		return nil, err
	}
	var m manifest
	if err := toml.Unmarshal(data, &m); err != nil {
		return nil, err
	}
	return m.Features, nil
}
