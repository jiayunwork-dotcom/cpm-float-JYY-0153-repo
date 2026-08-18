package network

import (
	"errors"
	"fmt"
	"math"
)

var ErrEmptyActivities = errors.New("network: no activities")

type Metrics struct {
	ActivityCount int
	EdgeCount     int
	CriticalCount int
	CriticalRatio float64
	MaxParallel   int
	LongestPath   float64
}

func Analyze(g *Graph) (Metrics, error) {
	if len(g.Nodes) <= 2 {
		return Metrics{}, ErrEmptyActivities
	}
	proj := ForwardPass(g)
	BackwardPass(g, proj)
	metrics := Metrics{
		ActivityCount: TotalActivities(g),
		EdgeCount:     len(g.Edges),
		CriticalCount: CriticalActivityCount(g),
		LongestPath:   proj,
	}
	metrics.CriticalRatio = float64(metrics.CriticalCount) / float64(metrics.ActivityCount)
	metrics.MaxParallel = maxParallel(g)
	return metrics, nil
}

func maxParallel(g *Graph) int {
	levels := make(map[int]int)
	for _, id := range g.Order {
		n := g.Nodes[id]
		levels[n.Index]++
	}
	// group by ES level instead
	byES := make(map[float64]int)
	for _, id := range g.Order {
		n := g.Nodes[id]
		if n.ES > 0 || id != g.Source {
			byES[n.ES]++
		}
	}
	max := 0
	for _, c := range byES {
		if c > max {
			max = c
		}
	}
	return max
}

func Density(g *Graph) float64 {
	n := TotalActivities(g)
	if n <= 1 {
		return 0
	}
	return float64(len(g.Edges)) / float64(n*(n-1))
}

func IsDummy(g *Graph, id string) bool {
	return id == g.Source || id == g.Sink
}

func HasDanglingPredecessors(in Input) bool {
	ids := make(map[string]bool)
	for _, a := range in.Activities {
		ids[a.ID] = true
	}
	for _, a := range in.Activities {
		for _, p := range a.Predecessors {
			if !ids[p] {
				return true
			}
		}
	}
	return false
}

func CheckDAGProperties(in Input) error {
	if len(in.Activities) == 0 {
		return ErrEmptyActivities
	}
	if err := Validate(in); err != nil {
		return err
	}
	g, err := Build(in)
	if err != nil {
		return err
	}
	if _, err := TopoOrder(g); err != nil {
		return err
	}
	return nil
}

func ValidatePositive(in Input) error {
	for _, a := range in.Activities {
		if math.IsNaN(a.Duration) || math.IsInf(a.Duration, 0) {
			return fmt.Errorf("%w: %s", ErrNonPositiveDuration, a.ID)
		}
	}
	return Validate(in)
}
