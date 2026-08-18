package pert

import (
	"errors"
	"math"
	"testing"

	"cpm-float/internal/network"
)

func TestTripleMeanVariance(t *testing.T) {
	e, err := FromTriple(Triple{A: 2, M: 4, B: 6})
	if err != nil {
		t.Fatalf("FromTriple: %v", err)
	}
	if math.Abs(e.Mean-4) > 1e-12 {
		t.Fatalf("mean want 4, got %v", e.Mean)
	}
	if math.Abs(e.Variance-(16.0/36.0)) > 1e-12 {
		t.Fatalf("variance want 16/36, got %v", e.Variance)
	}
}

func TestTripleValidation(t *testing.T) {
	if _, err := FromTriple(Triple{A: 5, M: 4, B: 6}); !errors.Is(err, ErrBadTriple) {
		t.Fatalf("a>m should fail, got %v", err)
	}
	if _, err := FromTriple(Triple{A: 2, M: 4, B: 3}); !errors.Is(err, ErrBadTriple) {
		t.Fatalf("m>b should fail, got %v", err)
	}
	if _, err := FromTriple(Triple{A: -1, M: 2, B: 4}); !errors.Is(err, ErrBadTriple) {
		t.Fatalf("negative a should fail, got %v", err)
	}
}

func TestPhiStandard(t *testing.T) {
	if math.Abs(Phi(0)-0.5) > 1e-9 {
		t.Fatalf("Phi(0) want 0.5, got %v", Phi(0))
	}
	if math.Abs(Phi(1.96)-0.9750021) > 1e-4 {
		t.Fatalf("Phi(1.96) want ~0.975, got %v", Phi(1.96))
	}
	if math.Abs(Phi(-1.96)+0.9750021-1) > 1e-4 {
		t.Fatalf("Phi(-1.96) symmetry broken: %v", Phi(-1.96))
	}
}

func TestCompletionProbability(t *testing.T) {
	p := CompletionProbability(20, 4, 20)
	if math.Abs(p-0.5) > 1e-9 {
		t.Fatalf("P(mean) want 0.5, got %v", p)
	}
	high := CompletionProbability(20, 4, 28)
	if high <= 0.95 || high >= 1 {
		t.Fatalf("P(mean+2sigma) want ~0.977, got %v", high)
	}
}

func TestCriticalPathVarianceDominant(t *testing.T) {
	in := network.Input{Activities: []network.Activity{
		{ID: "a", Duration: 1},
		{ID: "b", Duration: 1, Predecessors: []string{"a"}},
		{ID: "c", Duration: 3, Predecessors: []string{"a"}},
		{ID: "d", Duration: 1, Predecessors: []string{"b", "c"}},
	}}
	triples := map[string]Triple{
		"a": {A: 1, M: 1, B: 1},
		"b": {A: 1, M: 1, B: 1},
		"c": {A: 3, M: 3, B: 3},
		"d": {A: 1, M: 1, B: 1},
	}
	proj, err := BuildProject(in, triples)
	if err != nil {
		t.Fatalf("BuildProject: %v", err)
	}
	if !proj.IsCritical("c") {
		t.Fatal("c should be on the critical path")
	}
	if proj.IsCritical("b") {
		t.Fatal("b should not be on the critical path")
	}
}

func TestVarianceOnlyAlongCritical(t *testing.T) {
	in := network.Input{Activities: []network.Activity{
		{ID: "a", Duration: 1},
		{ID: "b", Duration: 1, Predecessors: []string{"a"}},
		{ID: "c", Duration: 1, Predecessors: []string{"a"}},
		{ID: "d", Duration: 1, Predecessors: []string{"b", "c"}},
	}}
	base := map[string]Triple{
		"a": {A: 1, M: 1, B: 1},
		"b": {A: 1, M: 1, B: 1},
		"c": {A: 1, M: 1, B: 1},
		"d": {A: 1, M: 1, B: 1},
	}
	p1, err := BuildProject(in, base)
	if err != nil {
		t.Fatalf("BuildProject: %v", err)
	}
	// both b and c are critical; widen b: project variance must change
	triplesB := map[string]Triple{
		"a": {A: 1, M: 1, B: 1},
		"b": {A: 0, M: 1, B: 6},
		"c": {A: 1, M: 1, B: 1},
		"d": {A: 1, M: 1, B: 1},
	}
	p2, err := BuildProject(in, triplesB)
	if err != nil {
		t.Fatalf("BuildProject: %v", err)
	}
	if p2.VarianceAlongCritical() <= p1.VarianceAlongCritical() {
		t.Fatalf("widening a critical activity must raise variance: %v -> %v",
			p1.VarianceAlongCritical(), p2.VarianceAlongCritical())
	}
}

func TestTargetProbabilityRoundTrip(t *testing.T) {
	mean, std := 20.0, 4.0
	target := TargetForProbability(mean, std, 0.9)
	p := CompletionProbability(mean, std, target)
	if math.Abs(p-0.9) > 1e-4 {
		t.Fatalf("round trip want 0.9, got %v", p)
	}
}
