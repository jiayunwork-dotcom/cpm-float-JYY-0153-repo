package network

import "sort"

type SlackProfile struct {
	ID      string
	ES, EF  float64
	LS, LF  float64
	TF, FF  float64
}

func FloatTable(g *Graph, project float64) []SlackProfile {
	profiles := []SlackProfile{}
	for _, id := range g.Order {
		n := g.Nodes[id]
		ff := 0.0
		succ := g.Successors[id]
		if len(succ) > 0 {
			earliest := inf()
			for _, s := range succ {
				if g.Nodes[s].ES < earliest {
					earliest = g.Nodes[s].ES
				}
			}
			ff = earliest - n.EF
			if ff < 0 {
				ff = 0
			}
		}
		profiles = append(profiles, SlackProfile{
			ID: id, ES: n.ES, EF: n.EF, LS: n.LS, LF: n.LF, TF: n.TF, FF: ff,
		})
	}
	return profiles
}

func inf() float64 {
	return 1e18
}

func SortedCriticalIDs(g *Graph) []string {
	ids := []string{}
	for id, n := range g.Nodes {
		if n.Critical {
			ids = append(ids, id)
		}
	}
	sort.Strings(ids)
	return ids
}

func TotalFloatOf(g *Graph, id string) (float64, bool) {
	n := g.Nodes[id]
	if n == nil {
		return 0, false
	}
	return n.TF, true
}

func LateStartOf(g *Graph, id string) (float64, bool) {
	n := g.Nodes[id]
	if n == nil {
		return 0, false
	}
	return n.LS, true
}

func EarlyFinishOf(g *Graph, id string) (float64, bool) {
	n := g.Nodes[id]
	if n == nil {
		return 0, false
	}
	return n.EF, true
}

func CriticalEdgeCount(g *Graph) int {
	n := 0
	for _, e := range g.Edges {
		from := g.Nodes[e.From]
		to := g.Nodes[e.To]
		if from != nil && to != nil && from.Critical && to.Critical {
			n++
		}
	}
	return n
}

func TotalEdgeCount(g *Graph) int {
	return len(g.Edges)
}
