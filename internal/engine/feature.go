package engine

import (
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

// RegionContribution is a flat chunk a feature wants placed inside a named
// managed region of a shared file.
type RegionContribution struct {
	Target  string `toml:"target"`
	Region  string `toml:"region"`
	Content string `toml:"content"`
}

// Feature is one entry in the catalog, loaded from its feature.toml.
type Feature struct {
	Name        string               `toml:"name"`
	Always      bool                 `toml:"always"`
	Description string               `toml:"description"`
	Regions     []RegionContribution `toml:"regions"`
	PostSteps   []string             `toml:"post_steps"`

	Dir string `toml:"-"` // path to the feature directory, set by the loader
}

// LoadFeature reads <dir>/feature.toml.
func LoadFeature(dir string) (Feature, error) {
	var f Feature
	data, err := os.ReadFile(filepath.Join(dir, "feature.toml"))
	if err != nil {
		return f, err
	}
	if err := toml.Unmarshal(data, &f); err != nil {
		return f, err
	}
	f.Dir = dir
	return f, nil
}

// LoadCatalog loads every immediate subdirectory of featuresDir as a feature.
func LoadCatalog(featuresDir string) ([]Feature, error) {
	entries, err := os.ReadDir(featuresDir)
	if err != nil {
		return nil, err
	}
	var features []Feature
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		f, err := LoadFeature(filepath.Join(featuresDir, e.Name()))
		if err != nil {
			return nil, err
		}
		features = append(features, f)
	}
	return features, nil
}
