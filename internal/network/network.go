package network

import (
	"errors"
	"fmt"
	"sort"
)

type Activity struct {
	ID        string  `json:"id"`
	Duration  float64 `json:"duration"`
	Predecessors []string `json:"predecessors,omitempty"`
	Opt    *float64 `json:"a,omitempty"`
	Most   *float64 `json:"m,omitempty"`
	Pess   *float64 `json:"b,omitempty"`
}

type Input struct {
	Activities []Activity `json:"activities"`
}

var (
	ErrUnknownPredecessor = errors.New("network: unknown predecessor")
	ErrSelfLoop           = errors.New("network: self loop")
	ErrCycle              = errors.New("network: directed cycle")
	ErrNonPositiveDuration = errors.New("network: duration must be positive")
	ErrDuplicateID        = errors.New("network: duplicate id")
	ErrBadTriple          = errors.New("network: pert triple invalid")
)

type Node struct {
	ID     string
	Index  int
	ES, EF float64
	LS, LF float64
	TF     float64
	Critical bool
}

type Edge struct {
	From string
	To   string
}

type Graph struct {
	Nodes    map[string]*Node
	Edges    []Edge
	Predecessors map[string][]string
	Successors   map[string][]string
	Source   string
	Sink     string
	Order    []string
	activities []Activity
}

func (g *Graph) node(id string) *Node {
	if n, ok := g.Nodes[id]; ok {
		return n
	}
	n := &Node{ID: id}
	g.Nodes[id] = n
	return n
}

func (g *Graph) AddEdge(from, to string) {
	g.Edges = append(g.Edges, Edge{From: from, To: to})
	g.Predecessors[to] = append(g.Predecessors[to], from)
	g.Successors[from] = append(g.Successors[from], to)
}

func Build(in Input) (*Graph, error) {
	if err := Validate(in); err != nil {
		return nil, err
	}
	g := &Graph{
		Nodes:        make(map[string]*Node),
		Predecessors: make(map[string][]string),
		Successors:   make(map[string][]string),
	}
	idx := 0
	for _, a := range in.Activities {
		n := g.node(a.ID)
		n.Index = idx
		idx++
	}
	hasPred := make(map[string]bool)
	for _, a := range in.Activities {
		for _, p := range a.Predecessors {
			g.AddEdge(p, a.ID)
			hasPred[a.ID] = true
		}
	}
	g.Source = "START"
	g.node(g.Source)
	for _, a := range in.Activities {
		if len(a.Predecessors) == 0 {
			g.AddEdge(g.Source, a.ID)
		}
	}
	g.Sink = "END"
	g.node(g.Sink)
	for _, a := range in.Activities {
		if len(g.Successors[a.ID]) == 0 {
			g.AddEdge(a.ID, g.Sink)
		}
	}
	g.Order, _ = TopoOrder(g)
	return g, nil
}

func Validate(in Input) error {
	seen := make(map[string]bool)
	for _, a := range in.Activities {
		if seen[a.ID] {
			return fmt.Errorf("%w: %s", ErrDuplicateID, a.ID)
		}
		seen[a.ID] = true
		if a.Duration <= 0 {
			return fmt.Errorf("%w: %s", ErrNonPositiveDuration, a.ID)
		}
		for _, p := range a.Predecessors {
			if p == a.ID {
				return fmt.Errorf("%w: %s", ErrSelfLoop, a.ID)
			}
			if !seen[p] && !seenInLater(in, p, a.ID) {
				return fmt.Errorf("%w: %s", ErrUnknownPredecessor, p)
			}
		}
	}
	if err := checkCycle(in); err != nil {
		return err
	}
	return nil
}

func seenInLater(in Input, pred, id string) bool {
	for _, a := range in.Activities {
		if a.ID == pred {
			return true
		}
	}
	return false
}

func checkCycle(in Input) error {
	graph := make(map[string][]string)
	for _, a := range in.Activities {
		for _, p := range a.Predecessors {
			graph[p] = append(graph[p], a.ID)
		}
	}
	state := make(map[string]int)
	var visit func(id string) error
	visit = func(id string) error {
		if state[id] == 2 {
			return nil
		}
		if state[id] == 1 {
			return fmt.Errorf("%w: %s", ErrCycle, id)
		}
		state[id] = 1
		for _, n := range graph[id] {
			if err := visit(n); err != nil {
				return err
			}
		}
		state[id] = 2
		return nil
	}
	for _, a := range in.Activities {
		if err := visit(a.ID); err != nil {
			return err
		}
	}
	return nil
}

func TopoOrder(g *Graph) ([]string, error) {
	indeg := make(map[string]int)
	for _, n := range g.Nodes {
		indeg[n.ID] = len(g.Predecessors[n.ID])
	}
	queue := []string{}
	for id, d := range indeg {
		if d == 0 {
			queue = append(queue, id)
		}
	}
	sort.Strings(queue)
	order := []string{}
	for len(queue) > 0 {
		id := queue[0]
		queue = queue[1:]
		order = append(order, id)
		for _, n := range g.Successors[id] {
			indeg[n]--
			if indeg[n] == 0 {
				queue = append(queue, n)
			}
		}
		sort.Strings(queue)
	}
	if len(order) != len(g.Nodes) {
		return nil, ErrCycle
	}
	return order, nil
}
