package pert

import (
	"fmt"
	"strings"

	"cpm-float/internal/network"
)

func ExtractTriples(in network.Input) (map[string]Triple, error) {
	triples := make(map[string]Triple)
	for _, a := range in.Activities {
		if a.Opt == nil || a.Most == nil || a.Pess == nil {
			continue
		}
		t := Triple{A: *a.Opt, M: *a.Most, B: *a.Pess}
		if _, err := FromTriple(t); err != nil {
			return nil, fmt.Errorf("%w: %s", ErrBadTriple, a.ID)
		}
		triples[a.ID] = t
	}
	return triples, nil
}

func FillDurations(in network.Input, triples map[string]Triple) network.Input {
	out := in
	out.Activities = append([]network.Activity(nil), in.Activities...)
	for i := range out.Activities {
		if t, ok := triples[out.Activities[i].ID]; ok {
			out.Activities[i].Duration = DurationFromTriple(t)
		}
	}
	return out
}

func HasTriples(in network.Input) bool {
	for _, a := range in.Activities {
		if a.Opt != nil && a.Most != nil && a.Pess != nil {
			return true
		}
	}
	return false
}

func FormatProbability(p *Project, target float64) string {
	prob := p.Probability(target)
	return fmt.Sprintf("P(duration <= %.2f) = %.4f (mean=%.2f std=%.2f)",
		target, prob, p.Duration, p.StdDev)
}

func FormatCriticalPath(p *Project) string {
	ids := p.CriticalWithoutDummies()
	return strings.Join(ids, " -> ")
}

func FormatStats(p *Project) string {
	var b strings.Builder
	b.WriteString(fmt.Sprintf("project_mean    %.4f\n", p.Duration))
	b.WriteString(fmt.Sprintf("project_sigma   %.4f\n", p.StdDev))
	b.WriteString(fmt.Sprintf("variance_along %s\n", FormatCriticalPath(p)))
	return b.String()
}

func MeanFromActivity(a network.Activity) (float64, error) {
	if a.Opt == nil || a.Most == nil || a.Pess == nil {
		return a.Duration, nil
	}
	e, err := FromTriple(Triple{A: *a.Opt, M: *a.Most, B: *a.Pess})
	if err != nil {
		return 0, err
	}
	return e.Mean, nil
}

func VarianceFromActivity(a network.Activity) (float64, error) {
	if a.Opt == nil || a.Most == nil || a.Pess == nil {
		return 0, nil
	}
	e, err := FromTriple(Triple{A: *a.Opt, M: *a.Most, B: *a.Pess})
	if err != nil {
		return 0, err
	}
	return e.Variance, nil
}
