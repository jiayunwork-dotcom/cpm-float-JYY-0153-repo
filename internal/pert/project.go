package pert

import (
	"math"
	"sort"

	"cpm-float/internal/network"
)

type Project struct {
	Network    *network.Graph
	Triples    map[string]Triple
	Critical   PathStat
	Duration   float64
	StdDev     float64
}

func BuildProject(in network.Input, triples map[string]Triple) (*Project, error) {
	net, err := network.Build(in)
	if err != nil {
		return nil, err
	}
	path, err := CriticalPathStats(net, triples)
	if err != nil {
		return nil, err
	}
	return &Project{
		Network:  net,
		Triples:  triples,
		Critical: path,
		Duration: path.Mean,
		StdDev:   math.Sqrt(path.Variance),
	}, nil
}

func (p *Project) Probability(target float64) float64 {
	return CompletionProbability(p.Duration, p.StdDev, target)
}

func (p *Project) Target(probability float64) float64 {
	return TargetForProbability(p.Duration, p.StdDev, probability)
}

func (p *Project) SortedCriticalActivities() []string {
	out := append([]string(nil), p.Critical.Path...)
	sort.Strings(out)
	return out
}

func (p *Project) ContainsActivity(id string) bool {
	for _, a := range p.Critical.Path {
		if a == id {
			return true
		}
	}
	return false
}

func (p *Project) CriticalWithoutDummies() []string {
	out := []string{}
	for _, a := range p.Critical.Path {
		if a == p.Network.Source || a == p.Network.Sink {
			continue
		}
		out = append(out, a)
	}
	return out
}

func (p *Project) IsCritical(id string) bool {
	for _, a := range p.Critical.Path {
		if a == id {
			return true
		}
	}
	return false
}

func (p *Project) VarianceAlongCritical() float64 {
	return p.Critical.Variance
}

func (p *Project) MeanAlongCritical() float64 {
	return p.Critical.Mean
}

func ReplaceDuration(in network.Input, id string, t Triple) network.Input {
	out := append([]network.Activity(nil), in.Activities...)
	for i := range out {
		if out[i].ID == id {
			out[i].Opt = &t.A
			out[i].Most = &t.M
			out[i].Pess = &t.B
		}
	}
	return network.Input{Activities: out}
}
