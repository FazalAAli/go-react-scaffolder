package engine

import (
	"fmt"
	"path/filepath"
	"sort"
)

// RegionKey identifies a managed region within a specific shared file.
type RegionKey struct {
	Target string
	Region string
}

// Plan is the resolved, ordered work to produce a project.
type Plan struct {
	Features  []string               // applied order: always (alpha) then optional (alpha)
	FileRoots []string               // each feature's files/ dir, same order
	Regions   map[RegionKey][]string // region -> ordered contributions
	PostSteps []string
}

// Resolve includes all always-features plus the selected optional ones and
// produces a deterministic Plan.
func Resolve(catalog []Feature, selectedOptional []string) (Plan, error) {
	byName := map[string]Feature{}
	for _, f := range catalog {
		byName[f.Name] = f
	}
	for _, name := range selectedOptional {
		f, ok := byName[name]
		if !ok {
			return Plan{}, fmt.Errorf("unknown feature %q", name)
		}
		if f.Always {
			return Plan{}, fmt.Errorf("feature %q is always-on and cannot be selected", name)
		}
	}
	selected := map[string]bool{}
	for _, n := range selectedOptional {
		selected[n] = true
	}

	var always, optional []Feature
	for _, f := range catalog {
		switch {
		case f.Always:
			always = append(always, f)
		case selected[f.Name]:
			optional = append(optional, f)
		}
	}
	sort.Slice(always, func(i, j int) bool { return always[i].Name < always[j].Name })
	sort.Slice(optional, func(i, j int) bool { return optional[i].Name < optional[j].Name })
	ordered := append(always, optional...)

	plan := Plan{Regions: map[RegionKey][]string{}}
	for _, f := range ordered {
		plan.Features = append(plan.Features, f.Name)
		plan.FileRoots = append(plan.FileRoots, filepath.Join(f.Dir, "files"))
		for _, r := range f.Regions {
			key := RegionKey{Target: r.Target, Region: r.Region}
			plan.Regions[key] = append(plan.Regions[key], r.Content)
		}
		plan.PostSteps = append(plan.PostSteps, f.PostSteps...)
	}
	return plan, nil
}
