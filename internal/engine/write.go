package engine

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// Write executes a plan into projectRoot: copy each feature's files/ tree (in
// order, later features overwrite earlier on conflict), re-render managed
// regions in any shared target file, then write scaffold.toml.
func Write(plan Plan, projectRoot string) error {
	if err := copyPayloads(plan.FileRoots, projectRoot); err != nil {
		return err
	}
	if err := renderPlanRegions(plan, projectRoot); err != nil {
		return err
	}
	return writeScaffoldToml(plan.Features, projectRoot)
}

// AddFeatures applies newly-selected features to an existing project. Unlike
// Write it copies only the payloads of the features named in `add` (so it does
// not clobber the user's edits to files owned by already-applied features), but
// it re-renders every managed region from the full plan, so shared files end up
// with the union of all active features' contributions. It then rewrites
// scaffold.toml to the plan's full feature set.
func AddFeatures(plan Plan, projectRoot string, add map[string]bool) error {
	var roots []string
	for i, name := range plan.Features {
		if add[name] {
			roots = append(roots, plan.FileRoots[i])
		}
	}
	if err := copyPayloads(roots, projectRoot); err != nil {
		return err
	}
	if err := renderPlanRegions(plan, projectRoot); err != nil {
		return err
	}
	return writeScaffoldToml(plan.Features, projectRoot)
}

// copyPayloads copies each existing files/ tree into projectRoot in order;
// later roots overwrite earlier ones on conflict. Missing roots are skipped.
func copyPayloads(roots []string, projectRoot string) error {
	for _, root := range roots {
		if _, err := os.Stat(root); err != nil {
			if os.IsNotExist(err) {
				continue // a feature may have no payload
			}
			return err
		}
		if err := copyTree(root, projectRoot); err != nil {
			return err
		}
	}
	return nil
}

// renderPlanRegions re-renders every managed region target in projectRoot from
// the plan's contributions, replacing each region body in place.
func renderPlanRegions(plan Plan, projectRoot string) error {
	byTarget := map[string]map[string][]string{}
	for key, contribs := range plan.Regions {
		if byTarget[key.Target] == nil {
			byTarget[key.Target] = map[string][]string{}
		}
		byTarget[key.Target][key.Region] = contribs
	}
	for target, regions := range byTarget {
		path := filepath.Join(projectRoot, target)
		info, err := os.Stat(path)
		if err != nil {
			return err
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		rendered, err := RenderRegions(string(data), regions)
		if err != nil {
			return err
		}
		if err := os.WriteFile(path, []byte(rendered), info.Mode()); err != nil {
			return err
		}
	}
	return nil
}

func writeScaffoldToml(features []string, projectRoot string) error {
	sorted := append([]string(nil), features...)
	sort.Strings(sorted)
	var b strings.Builder
	b.WriteString("features = [")
	for i, f := range sorted {
		if i > 0 {
			b.WriteString(", ")
		}
		fmt.Fprintf(&b, "%q", f)
	}
	b.WriteString("]\n")
	return os.WriteFile(filepath.Join(projectRoot, "scaffold.toml"), []byte(b.String()), 0o644)
}

func copyTree(srcRoot, dstRoot string) error {
	return filepath.WalkDir(srcRoot, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(srcRoot, path)
		if err != nil {
			return err
		}
		dst := filepath.Join(dstRoot, rel)
		if d.IsDir() {
			return os.MkdirAll(dst, 0o755)
		}
		return copyFile(path, dst)
	})
}

func copyFile(src, dst string) error {
	if err := os.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
		return err
	}
	info, err := os.Stat(src)
	if err != nil {
		return err
	}
	data, err := os.ReadFile(src)
	if err != nil {
		return err
	}
	return os.WriteFile(dst, data, info.Mode())
}
