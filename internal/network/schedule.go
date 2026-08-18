package network

import "math"

type Schedule struct {
	ProjectDuration float64
	Critical        []string
	CriticalPath    []string
	ParallelCritical int
}

func ForwardPass(g *Graph) float64 {
	for _, id := range g.Order {
		n := g.Nodes[id]
		es := 0.0
		for _, p := range g.Predecessors[id] {
			pn := g.Nodes[p]
			if pn.EF > es {
				es = pn.EF
			}
		}
		n.ES = es
		n.EF = es + durationOf(g, id)
	}
	return g.Nodes[g.Sink].EF
}

func durationOf(g *Graph, id string) float64 {
	if id == g.Source || id == g.Sink {
		return 0
	}
	act := activityByID(g, id)
	return act.Duration
}

func BackwardPass(g *Graph, project float64) {
	for i := len(g.Order) - 1; i >= 0; i-- {
		id := g.Order[i]
		n := g.Nodes[id]
		lf := project
		succ := g.Successors[id]
		if len(succ) > 0 {
			lf = math.Inf(1)
			for _, s := range succ {
				sn := g.Nodes[s]
				if sn.LS < lf {
					lf = sn.LS
				}
			}
		}
		n.LF = lf
		n.LS = lf - durationOf(g, id)
		n.TF = n.LS - n.ES
		n.Critical = math.Abs(n.TF) < 1e-9
	}
}

func ComputeTF(g *Graph) {
	for _, n := range g.Nodes {
		if n.ES == 0 && n.LS == 0 {
			continue
		}
		n.TF = n.LS - n.ES
	}
}

func FreeFloat(g *Graph, id string) float64 {
	n := g.Nodes[id]
	if n == nil {
		return 0
	}
	succ := g.Successors[id]
	if len(succ) == 0 {
		return 0
	}
	earliest := math.Inf(1)
	for _, s := range succ {
		if g.Nodes[s].ES < earliest {
			earliest = g.Nodes[s].ES
		}
	}
	ff := earliest - n.EF
	if ff < 0 {
		return 0
	}
	return ff
}

func ExtractCriticalPath(g *Graph, project float64) Schedule {
	BackwardPass(g, project)
	crit := []string{}
	for _, n := range g.Nodes {
		if n.Critical {
			crit = append(crit, n.ID)
		}
	}
	path := pathThrough(g, g.Source, g.Sink)
	return Schedule{
		ProjectDuration: project,
		Critical:        crit,
		CriticalPath:    path,
	}
}

func pathThrough(g *Graph, from, to string) []string {
	seen := make(map[string]bool)
	var walk func(id string) []string
	walk = func(id string) []string {
		if id == to {
			return []string{id}
		}
		if seen[id] {
			return nil
		}
		seen[id] = true
		best := []string(nil)
		for _, s := range g.Successors[id] {
			if !g.Nodes[s].Critical {
				continue
			}
			tail := walk(s)
			if tail != nil && len(tail) > len(best) {
				best = append([]string{id}, tail...)
			}
		}
		seen[id] = false
		return best
	}
	return walk(from)
}

func TFIdentityHolds(g *Graph, id string, tol float64) bool {
	n := g.Nodes[id]
	return math.Abs(n.TF-(n.LS-n.ES)) < tol &&
		math.Abs(n.TF-(n.LF-n.EF)) < tol
}

func CriticalSpansSourceToSink(s Schedule, g *Graph) bool {
	if len(s.CriticalPath) == 0 {
		return false
	}
	return s.CriticalPath[0] == g.Source &&
		s.CriticalPath[len(s.CriticalPath)-1] == g.Sink
}
