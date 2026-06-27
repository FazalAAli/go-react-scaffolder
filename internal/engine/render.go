package engine

import (
	"fmt"
	"strings"
)

// RenderRegions replaces the body of each managed region in content with the
// joined contributions for that region name. A region is delimited by lines
// containing "scaffold:region:<name>:start" and "scaffold:region:<name>:end";
// the comment prefix is irrelevant. Markers are preserved. Re-rendering is
// idempotent: existing body content between markers is discarded and replaced.
func RenderRegions(content string, contributions map[string][]string) (string, error) {
	lines := strings.Split(content, "\n")
	var out []string
	for i := 0; i < len(lines); i++ {
		name, ok := regionMarker(lines[i], "start")
		if !ok {
			out = append(out, lines[i])
			continue
		}
		// find matching end marker
		end := -1
		for j := i + 1; j < len(lines); j++ {
			if n, ok := regionMarker(lines[j], "end"); ok && n == name {
				end = j
				break
			}
		}
		if end < 0 {
			return "", fmt.Errorf("region %q: missing end marker", name)
		}
		out = append(out, lines[i]) // start marker
		for _, c := range contributions[name] {
			out = append(out, strings.Split(strings.TrimRight(c, "\n"), "\n")...)
		}
		out = append(out, lines[end]) // end marker
		i = end
	}
	return strings.Join(out, "\n"), nil
}

func regionMarker(line, kind string) (string, bool) {
	const prefix = "scaffold:region:"
	idx := strings.Index(line, prefix)
	if idx < 0 {
		return "", false
	}
	rest := strings.TrimSpace(line[idx+len(prefix):])
	suffix := ":" + kind
	if !strings.HasSuffix(rest, suffix) {
		return "", false
	}
	return strings.TrimSuffix(rest, suffix), true
}
