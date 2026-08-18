package pert

import (
	"math"
	"math/rand"
	"testing"
)

func TestSimulateCloseToMean(t *testing.T) {
	triples := map[string]Triple{
		"a": {A: 1, M: 2, B: 3},
		"b": {A: 1, M: 2, B: 3},
	}
	p := &Project{
		Triples: triples,
		Critical: PathStat{
			Path:     []string{"START", "a", "b", "END"},
			Mean:     4,
			Variance: 2 * (2.0/6.0)*(2.0/6.0),
		},
		Duration: 4,
		StdDev:   math.Sqrt(2 * (2.0/6.0)*(2.0/6.0)),
	}
	sim := SimulateDuration(p, 2000, 42)
	if math.Abs(sim.Mean-4) > 0.2 {
		t.Fatalf("simulated mean %.3f too far from 4", sim.Mean)
	}
	if sim.StdDev <= 0 {
		t.Fatal("simulated std dev should be positive")
	}
	if sim.P25 > sim.P50 || sim.P50 > sim.P75 {
		t.Fatalf("quantiles out of order: %.3f %.3f %.3f", sim.P25, sim.P50, sim.P75)
	}
}

func TestSimulatedProbability(t *testing.T) {
	s := SimResult{Samples: []float64{1, 2, 3, 4, 5}}
	if got := SimulatedProbability(s, 3); math.Abs(got-0.6) > 1e-12 {
		t.Fatalf("P(<=3) want 0.6, got %v", got)
	}
	if got := SimulatedProbability(s, 10); got != 1 {
		t.Fatalf("P(<=10) want 1, got %v", got)
	}
}

func TestConfidenceInterval(t *testing.T) {
	s := SimResult{Mean: 20, StdDev: 2}
	lo, hi := ConfidenceInterval(s)
	if math.Abs(lo-(20-1.96*2)) > 1e-9 || math.Abs(hi-(20+1.96*2)) > 1e-9 {
		t.Fatalf("interval wrong: %v %v", lo, hi)
	}
}

func TestPercentile(t *testing.T) {
	got := percentile([]float64{3, 1, 2}, 0.5)
	if got != 2 {
		t.Fatalf("median want 2, got %v", got)
	}
}

func TestSampleTripleDeterministic(t *testing.T) {
	triple := Triple{A: 2, M: 4, B: 6}
	a := SampleTriple(triple, rand.New(rand.NewSource(1)))
	b := SampleTriple(triple, rand.New(rand.NewSource(1)))
	if a != b {
		t.Fatalf("same seed should give same sample: %v %v", a, b)
	}
}
