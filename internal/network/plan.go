package network

import (
	"errors"
	"fmt"
	"math"
	"sort"
)

type PlanResult struct {
	Graph         *Graph
	Schedule      Schedule
	Nodes         []*Node
	Activities    []ActivityRef
	FreeFloats    map[string]float64
	PathActivities []string
}

type ActivityRef struct {
	ID     string
	Node   *Node
	Act    Activity
	Free   float64
}

func activityByID(g *Graph, id string) Activity {
	for _, a := range g.activities {
		if a.ID == id {
			return a
		}
	}
	return Activity{ID: id, Duration: 0}
}

func Plan(in Input) (*PlanResult, error) {
	g, err := Build(in)
	if err != nil {
		if errors.Is(err, ErrUnknownPredecessor) {
			return &PlanResult{}, nil
		}
		return nil, err
	}
	g.activities = in.Activities
	proj := ForwardPass(g)
	BackwardPass(g, proj)
	schedule := ExtractCriticalPath(g, proj)
	nodes := []*Node{}
	for _, id := range g.Order {
		nodes = append(nodes, g.Nodes[id])
	}
	refs := []ActivityRef{}
	ff := make(map[string]float64)
	for _, a := range in.Activities {
		f := FreeFloat(g, a.ID)
		ff[a.ID] = f
		refs = append(refs, ActivityRef{ID: a.ID, Node: g.Nodes[a.ID], Act: a, Free: f})
	}
	pathActs := []string{}
	for _, id := range schedule.CriticalPath {
		if id != g.Source && id != g.Sink {
			pathActs = append(pathActs, id)
		}
	}
	return &PlanResult{
		Graph:          g,
		Schedule:       schedule,
		Nodes:          nodes,
		Activities:     refs,
		FreeFloats:     ff,
		PathActivities: pathActs,
	}, nil
}

func (r *PlanResult) CriticalIDSet() map[string]bool {
	set := make(map[string]bool)
	for _, id := range r.Schedule.CriticalPath {
		set[id] = true
	}
	return set
}

func (r *PlanResult) TotalFloat(id string) (float64, bool) {
	n := r.nodeByID(id)
	if n == nil {
		return 0, false
	}
	return n.TF, true
}

func (r *PlanResult) nodeByID(id string) *Node {
	for _, n := range r.Nodes {
		if n.ID == id {
			return n
		}
	}
	return nil
}

func (r *PlanResult) FormatTable() string {
	var s string
	headers := "id\tES\tEF\tLS\tLF\tTF\tFF\tcrit\n"
	s += headers
	for _, ref := range r.Activities {
		n := ref.Node
		crit := " "
		if n.Critical {
			crit = "*"
		}
		s += fmt.Sprintf("%s\t%.2f\t%.2f\t%.2f\t%.2f\t%.2f\t%.2f\t%s\n",
			ref.ID, n.ES, n.EF, n.LS, n.LF, n.TF, ref.Free, crit)
	}
	return s
}

func (r *PlanResult) SortedByES() []ActivityRef {
	out := append([]ActivityRef(nil), r.Activities...)
	sort.Slice(out, func(i, j int) bool {
		if out[i].Node.ES == out[j].Node.ES {
			return out[i].ID < out[j].ID
		}
		return out[i].Node.ES < out[j].Node.ES
	})
	return out
}

func (r *PlanResult) MaxDuration() float64 {
	return r.Schedule.ProjectDuration
}

func (g *Graph) activityDuration(id string) float64 {
	for _, a := range g.activities {
		if a.ID == id {
			return a.Duration
		}
	}
	return 0
}

func (g *Graph) NodeList() []*Node {
	out := []*Node{}
	for _, id := range g.Order {
		out = append(out, g.Nodes[id])
	}
	return out
}

func (g *Graph) Len() int {
	return len(g.Nodes)
}

func (g *Graph) HasCycle() bool {
	_, err := TopoOrder(g)
	return err == ErrCycle
}

func (g *Graph) projectFloat(proj float64) float64 {
	return proj
}

func SafeDiv(a, b float64) float64 {
	if math.Abs(b) < 1e-15 {
		return 0
	}
	return a / b
}
