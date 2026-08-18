package network

import (
	"sort"
)

type CriticalSummary struct {
	IDs       []string
	Duration  float64
	Count     int
	AllCritical bool
}

func SummarizeCritical(g *Graph) CriticalSummary {
	proj := ForwardPass(g)
	BackwardPass(g, proj)
	ids := []string{}
	total := 0.0
	for id, n := range g.Nodes {
		if id == g.Source || id == g.Sink {
			continue
		}
		if n.Critical {
			ids = append(ids, id)
			total += g.activityDuration(id)
		}
	}
	sort.Strings(ids)
	return CriticalSummary{
		IDs:      ids,
		Duration: proj,
		Count:    len(ids),
	}
}

func IsCriticalActivity(g *Graph, id string) bool {
	n := g.Nodes[id]
	return n != nil && n.Critical
}

func FloatGap(a, b float64) float64 {
	if a < b {
		return b - a
	}
	return a - b
}

func AllZeroFloat(g *Graph) bool {
	for _, n := range g.Nodes {
		if n.Critical != (n.TF == 0) {
			return false
		}
	}
	return true
}

func EarliestStart(g *Graph, id string) float64 {
	n := g.Nodes[id]
	if n == nil {
		return 0
	}
	return n.ES
}

func LatestStart(g *Graph, id string) float64 {
	n := g.Nodes[id]
	if n == nil {
		return 0
	}
	return n.LS
}

func DurationAlong(g *Graph, ids []string) float64 {
	sum := 0.0
	for _, id := range ids {
		if id == g.Source || id == g.Sink {
			continue
		}
		sum += g.activityDuration(id)
	}
	return sum
}
