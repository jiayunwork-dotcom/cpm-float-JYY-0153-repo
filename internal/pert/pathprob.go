package pert

import (
	"math"

	"cpm-float/internal/network"
)

type PathProbability struct {
	Path       []string
	Mean       float64
	StdDev     float64
	Probability float64
}

func AllPathProbabilities(p *Project, target float64) []PathProbability {
	all := network.EnumeratePaths(p.Network)
	out := []PathProbability{}
	for _, path := range all.Paths {
		mean := 0.0
		variance := 0.0
		for _, id := range path {
			t, ok := p.Triples[id]
			if !ok {
				continue
			}
			e, err := FromTriple(t)
			if err != nil {
				continue
			}
			mean += e.Mean
			variance += e.Variance
		}
		std := math.Sqrt(variance)
		out = append(out, PathProbability{
			Path:        path,
			Mean:        mean,
			StdDev:      std,
			Probability: CompletionProbability(mean, std, target),
		})
	}
	return out
}

func MostLikelyDelayingPath(p *Project, target float64) string {
	probs := AllPathProbabilities(p, target)
	if len(probs) == 0 {
		return ""
	}
	best := probs[0]
	for _, pr := range probs {
		if pr.Mean > best.Mean {
			best = pr
		}
	}
	return formatPath(best.Path)
}

func formatPath(ids []string) string {
	out := ""
	for i, id := range ids {
		if i > 0 {
			out += " -> "
		}
		out += id
	}
	return out
}

func PathCount(p *Project) int {
	return network.EnumeratePaths(p.Network).Count
}

func WeakestPath(p *Project, target float64) PathProbability {
	probs := AllPathProbabilities(p, target)
	if len(probs) == 0 {
		return PathProbability{}
	}
	best := probs[0]
	for _, pr := range probs {
		if pr.Probability < best.Probability {
			best = pr
		}
	}
	return best
}
