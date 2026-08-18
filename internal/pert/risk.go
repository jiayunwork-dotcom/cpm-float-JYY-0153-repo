package pert

import (
	"math"
	"sort"
)

type ActivityRisk struct {
	ID       string
	Mean     float64
	Variance float64
	StdDev   float64
	Critical bool
}

type RiskReport struct {
	Activities []ActivityRisk
	TotalMean  float64
	TotalStd   float64
}

func RiskReportFor(p *Project) RiskReport {
	report := RiskReport{TotalMean: p.Duration, TotalStd: p.StdDev}
	ids := p.CriticalWithoutDummies()
	for _, id := range ids {
		t := p.Triples[id]
		e, err := FromTriple(t)
		if err != nil {
			continue
		}
		report.Activities = append(report.Activities, ActivityRisk{
			ID:       id,
			Mean:     e.Mean,
			Variance: e.Variance,
			StdDev:   math.Sqrt(e.Variance),
			Critical: true,
		})
	}
	return report
}

func MostRisky(p *Project) string {
	report := RiskReportFor(p)
	if len(report.Activities) == 0 {
		return ""
	}
	sort.Slice(report.Activities, func(i, j int) bool {
		return report.Activities[i].Variance > report.Activities[j].Variance
	})
	return report.Activities[0].ID
}

func ContributionOf(p *Project, id string) float64 {
	t, ok := p.Triples[id]
	if !ok {
		return 0
	}
	if !p.IsCritical(id) {
		return 0
	}
	e, err := FromTriple(t)
	if err != nil {
		return 0
	}
	if p.Critical.Variance <= 0 {
		return 0
	}
	return e.Variance / p.Critical.Variance
}

func CriticalSigma(p *Project) float64 {
	return math.Sqrt(p.Critical.Variance)
}

func ProbabilityRange(p *Project, lo, hi float64) float64 {
	return p.Probability(hi) - p.Probability(lo)
}

func MedianCompletion(p *Project) float64 {
	return p.Target(0.5)
}

func NinetyPercent(p *Project) float64 {
	return p.Target(0.9)
}
