package network

type Dependency struct {
	From string
	To   string
	Lag  float64
}

func ListDependencies(g *Graph) []Dependency {
	out := []Dependency{}
	for _, e := range g.Edges {
		from := g.Nodes[e.From]
		to := g.Nodes[e.To]
		lag := 0.0
		if from != nil && to != nil {
			lag = to.ES - from.EF
		}
		out = append(out, Dependency{From: e.From, To: e.To, Lag: lag})
	}
	return out
}

func ZeroLagEdges(g *Graph) []Dependency {
	all := ListDependencies(g)
	out := []Dependency{}
	for _, d := range all {
		if d.Lag == 0 {
			out = append(out, d)
		}
	}
	return out
}

func FreeFloatEdges(g *Graph) []Dependency {
	all := ListDependencies(g)
	out := []Dependency{}
	for _, d := range all {
		if d.Lag > 0 {
			out = append(out, d)
		}
	}
	return out
}

func ActivityChain(g *Graph, from, to string) []string {
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
		for _, s := range g.Successors[id] {
			chain := walk(s)
			if chain != nil {
				return append([]string{id}, chain...)
			}
		}
		seen[id] = false
		return nil
	}
	return walk(from)
}

func HasPath(g *Graph, from, to string) bool {
	return ActivityChain(g, from, to) != nil
}

func ImmediateSuccessors(g *Graph, id string) []string {
	return g.Successors[id]
}

func ImmediatePredecessors(g *Graph, id string) []string {
	return g.Predecessors[id]
}

func MergeSort(a, b []string) []string {
	seen := make(map[string]bool)
	out := []string{}
	for _, s := range append(a, b...) {
		if !seen[s] {
			seen[s] = true
			out = append(out, s)
		}
	}
	return out
}
