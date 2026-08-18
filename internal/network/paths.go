package network

import (
	"sort"
)

type AllPathsResult struct {
	Paths [][]string
	Count int
}

func EnumeratePaths(g *Graph) AllPathsResult {
	paths := [][]string{}
	seen := make(map[string]bool)
	var walk func(id string, cur []string)
	walk = func(id string, cur []string) {
		if id == g.Sink {
			path := append([]string(nil), cur...)
			path = append(path, id)
			paths = append(paths, path)
			return
		}
		if seen[id] {
			return
		}
		seen[id] = true
		next := append([]string(nil), cur...)
		next = append(next, id)
		succ := g.Successors[id]
		sort.Strings(succ)
		for _, s := range succ {
			walk(s, next)
		}
		seen[id] = false
	}
	walk(g.Source, []string{})
	return AllPathsResult{Paths: paths, Count: len(paths)}
}

type PathLength struct {
	Path   []string
	Length float64
}

func LongestPaths(g *Graph, project float64) []PathLength {
	all := EnumeratePaths(g)
	out := []PathLength{}
	for _, p := range all.Paths {
		length := 0.0
		for _, id := range p {
			length += g.activityDuration(id)
		}
		out = append(out, PathLength{Path: p, Length: length})
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].Length == out[j].Length {
			return len(out[i].Path) > len(out[j].Path)
		}
		return out[i].Length > out[j].Length
	})
	return out
}

func ParallelCriticalPaths(g *Graph, project float64) [][]string {
	longest := LongestPaths(g, project)
	if len(longest) == 0 {
		return nil
	}
	target := longest[0].Length
	out := [][]string{}
	for _, l := range longest {
		if abs(l.Length-target) < 1e-9 {
			out = append(out, l.Path)
		}
	}
	return out
}

func CriticalActivityCount(g *Graph) int {
	n := 0
	for id, node := range g.Nodes {
		if id == g.Source || id == g.Sink {
			continue
		}
		if node.Critical {
			n++
		}
	}
	return n
}

func TotalActivities(g *Graph) int {
	return len(g.Nodes) - 2
}

func CriticalRatio(g *Graph) float64 {
	total := TotalActivities(g)
	if total <= 0 {
		return 0
	}
	return float64(CriticalActivityCount(g)) / float64(total)
}

func abs(v float64) float64 {
	if v < 0 {
		return -v
	}
	return v
}
