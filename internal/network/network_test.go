package network

import (
	"errors"
	"testing"
)

func TestPlanHouseCriticalPath(t *testing.T) {
	in := Input{Activities: []Activity{
		{ID: "site-prep", Duration: 2},
		{ID: "foundation", Duration: 4, Predecessors: []string{"site-prep"}},
		{ID: "framing", Duration: 6, Predecessors: []string{"foundation"}},
		{ID: "roofing", Duration: 3, Predecessors: []string{"framing"}},
		{ID: "electrical", Duration: 4, Predecessors: []string{"framing"}},
		{ID: "plumbing", Duration: 3, Predecessors: []string{"framing"}},
		{ID: "drywall", Duration: 5, Predecessors: []string{"electrical", "plumbing"}},
		{ID: "painting", Duration: 2, Predecessors: []string{"drywall"}},
		{ID: "flooring", Duration: 2, Predecessors: []string{"drywall"}},
		{ID: "finish", Duration: 1, Predecessors: []string{"painting", "flooring"}},
	}}
	plan, err := Plan(in)
	if err != nil {
		t.Fatalf("Plan: %v", err)
	}
	if plan.MaxDuration() != 24 {
		t.Fatalf("project duration want 24, got %.2f", plan.MaxDuration())
	}
	if !CriticalSpansSourceToSink(plan.Schedule, plan.Graph) {
		t.Fatal("critical path should span source to sink")
	}
	for _, ref := range plan.Activities {
		if !TFIdentityHolds(plan.Graph, ref.ID, 1e-9) {
			t.Fatalf("TF identity broken for %s", ref.ID)
		}
	}
}

func TestPlanSimpleFloats(t *testing.T) {
	in := Input{Activities: []Activity{
		{ID: "design", Duration: 3},
		{ID: "probe", Duration: 1, Predecessors: []string{"design"}},
		{ID: "cut", Duration: 2, Predecessors: []string{"design"}},
		{ID: "assemble", Duration: 2, Predecessors: []string{"probe", "cut"}},
		{ID: "test", Duration: 1, Predecessors: []string{"assemble"}},
	}}
	plan, err := Plan(in)
	if err != nil {
		t.Fatalf("Plan: %v", err)
	}
	if plan.MaxDuration() != 8 {
		t.Fatalf("project duration want 8, got %.2f", plan.MaxDuration())
	}
	tf, ok := plan.TotalFloat("probe")
	if !ok {
		t.Fatal("probe not found")
	}
	if tf != 1 {
		t.Fatalf("probe TF want 1, got %.2f", tf)
	}
	ff, ok := plan.FreeFloats["probe"]
	if !ok || ff != 1 {
		t.Fatalf("probe FF want 1, got %v", ff)
	}
}

func TestRejectUnknownPredecessor(t *testing.T) {
	in := Input{Activities: []Activity{
		{ID: "a", Duration: 1},
		{ID: "b", Duration: 1, Predecessors: []string{"zzz"}},
	}}
	_, err := Plan(in)
	if !errors.Is(err, ErrUnknownPredecessor) {
		t.Fatalf("want ErrUnknownPredecessor, got %v", err)
	}
}

func TestRejectSelfLoop(t *testing.T) {
	in := Input{Activities: []Activity{
		{ID: "a", Duration: 1, Predecessors: []string{"a"}},
	}}
	_, err := Plan(in)
	if !errors.Is(err, ErrSelfLoop) {
		t.Fatalf("want ErrSelfLoop, got %v", err)
	}
}

func TestRejectCycle(t *testing.T) {
	in := Input{Activities: []Activity{
		{ID: "a", Duration: 1, Predecessors: []string{"c"}},
		{ID: "b", Duration: 1, Predecessors: []string{"a"}},
		{ID: "c", Duration: 1, Predecessors: []string{"b"}},
	}}
	_, err := Plan(in)
	if !errors.Is(err, ErrCycle) {
		t.Fatalf("want ErrCycle, got %v", err)
	}
}

func TestRejectZeroDuration(t *testing.T) {
	in := Input{Activities: []Activity{
		{ID: "a", Duration: 0},
	}}
	_, err := Plan(in)
	if !errors.Is(err, ErrNonPositiveDuration) {
		t.Fatalf("want ErrNonPositiveDuration, got %v", err)
	}
}

func TestRejectDuplicateID(t *testing.T) {
	in := Input{Activities: []Activity{
		{ID: "a", Duration: 1},
		{ID: "a", Duration: 2},
	}}
	_, err := Plan(in)
	if !errors.Is(err, ErrDuplicateID) {
		t.Fatalf("want ErrDuplicateID, got %v", err)
	}
}

func TestSingleActivity(t *testing.T) {
	in := Input{Activities: []Activity{
		{ID: "only", Duration: 5},
	}}
	plan, err := Plan(in)
	if err != nil {
		t.Fatalf("Plan: %v", err)
	}
	if plan.MaxDuration() != 5 {
		t.Fatalf("duration want 5, got %.2f", plan.MaxDuration())
	}
	if len(plan.PathActivities) != 1 || plan.PathActivities[0] != "only" {
		t.Fatalf("critical path should be [only], got %v", plan.PathActivities)
	}
}
